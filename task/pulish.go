package task

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rpcxio/libkv/store"
	etcdV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"strings"
	"sync"
	"time"
	"tinyIM/config"
	"tinyIM/tools"
)

// TODO: use rocketmq

var RedisClient *redis.Client

func (task *Task) InitQueueClient() error {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)

	//	check redis client
	if _, err := RedisClient.Ping(RedisClient.Context()).Result(); err != nil {
		return err
	}

	go func() {
		for {
			// 10s timeout
			result, err := RedisClient.BRPop(RedisClient.Context(), 10*time.Second, config.QueueName).Result()
			if err != nil {
				logrus.Infof("task queue block timeout,no msg err:%s", err.Error())
			} else {
				// result 为一个包含队列名称和弹出值的切片
				task.Push(result[1])
			}
		}
	}()

	return nil
}

var RClient = &RpcConnectClient{
	ServerInsMap: make(map[string][]Instance),
	IndexMap:     make(map[string]int),
}

type Instance struct {
	ServerType string
	ServerId   string
	Client     client.XClient
}

type RpcConnectClient struct {
	lock         sync.RWMutex
	ServerInsMap map[string][]Instance //serverId--[]ins
	IndexMap     map[string]int        //serverId--index
}

func (rc *RpcConnectClient) GetRpcClientByServerId(serverId string) (c client.XClient, err error) {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	if _, ok := rc.ServerInsMap[serverId]; !ok || len(rc.ServerInsMap[serverId]) <= 0 {
		return nil, errors.New("no connect layer ip:" + serverId)
	}
	if _, ok := rc.IndexMap[serverId]; !ok {
		rc.IndexMap[serverId] = 0
	}
	idx := rc.IndexMap[serverId] % len(rc.ServerInsMap[serverId])
	ins := rc.ServerInsMap[serverId][idx]
	rc.IndexMap[serverId] = (rc.IndexMap[serverId] + 1) % len(rc.ServerInsMap[serverId])
	return ins.Client, nil
}

func (rc *RpcConnectClient) GetAllConnectTypeRpcClient() (rpcClientList []client.XClient) {
	for serverId := range rc.ServerInsMap {
		c, err := rc.GetRpcClientByServerId(serverId)
		if err != nil {
			logrus.Infof("GetAllConnectTypeRpcClient err:%s", err.Error())
			continue
		}
		rpcClientList = append(rpcClientList, c)
	}
	return
}

func (task *Task) InitConnectRpcClient() error {
	etcdConfigOption := &store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: time.Duration(config.Conf.Common.CommonEtcd.ConnectionTimeout) * time.Second,
		Bucket:            "",
		PersistConnection: true,
		Username:          config.Conf.Common.CommonEtcd.UserName,
		Password:          config.Conf.Common.CommonEtcd.Password,
	}
	etcdConfig := config.Conf.Common.CommonEtcd
	d, err := etcdV3.NewEtcdV3Discovery(
		etcdConfig.BasePath,
		etcdConfig.ServerPathConnect,
		[]string{etcdConfig.Host},
		true,
		etcdConfigOption,
	)
	if err != nil {
		return err
	}
	if len(d.GetServices()) <= 0 {
		return errors.New("no etcd server find")
	}

	for _, connConf := range d.GetServices() {
		logrus.Infof("get etcd connect rpcx server, key is:%s,value is:%s", connConf.Key, connConf.Value)
		serverType := getParamByKey(connConf.Value, "serverType")
		serverId := getParamByKey(connConf.Value, "serverId")
		logrus.Infof("serverType is:%s,serverId is:%s", serverType, serverId)
		if serverType == "" || serverId == "" {
			continue
		}

		d, err := client.NewPeer2PeerDiscovery(connConf.Key, "")
		if err != nil {
			logrus.Errorf("init task client.NewPeer2PeerDiscovery client fail:%s", err.Error())
			continue
		}

		c := client.NewXClient(etcdConfig.ServerPathConnect, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		ins := Instance{
			ServerType: serverType,
			ServerId:   serverId,
			Client:     c,
		}
		if _, ok := RClient.ServerInsMap[serverId]; !ok {
			RClient.ServerInsMap[serverId] = []Instance{ins}
		} else {
			RClient.ServerInsMap[serverId] = append(RClient.ServerInsMap[serverId], ins)
		}
	}

	//	watch connect server && update RpcConnectClientList
	go task.watchServicesChange(d)
	return nil
}

func (task *Task) watchServicesChange(d client.ServiceDiscovery) {
	etcdConfig := config.Conf.Common.CommonEtcd
	for kvChan := range d.WatchService() {
		if len(kvChan) <= 0 {
			logrus.Errorf("connect services change, connect alarm, no abailable ip")
		}
		logrus.Infof("connect services change trigger...")
		insMap := make(map[string][]Instance)
		for _, kv := range kvChan {
			logrus.Infof("connect services change,key is:%s,value is:%s", kv.Key, kv.Value)
			serverType := getParamByKey(kv.Value, "serverType")
			serverId := getParamByKey(kv.Value, "serverId")
			logrus.Infof("serverType is:%s,serverId is:%s", serverType, serverId)
			if serverType == "" || serverId == "" {
				continue
			}
			d, e := client.NewPeer2PeerDiscovery(kv.Key, "")
			if e != nil {
				logrus.Errorf("init task client.NewPeer2PeerDiscovery watch client fail:%s", e.Error())
				continue
			}
			c := client.NewXClient(etcdConfig.ServerPathConnect, client.Failtry, client.RandomSelect, d, client.DefaultOption)
			ins := Instance{
				ServerType: serverType,
				ServerId:   serverId,
				Client:     c,
			}
			if _, ok := insMap[serverId]; !ok {
				insMap[serverId] = []Instance{ins}
			} else {
				insMap[serverId] = append(insMap[serverId], ins)
			}
		}
		RClient.lock.Lock()
		RClient.ServerInsMap = insMap
		RClient.lock.Unlock()
	}
}

func getParamByKey(s string, key string) string {
	params := strings.Split(s, "&")
	for _, p := range params {
		kv := strings.Split(p, "=")
		if len(kv) == 2 && kv[0] == key {
			return kv[1]
		}
	}
	return ""
}
