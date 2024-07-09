package logic

import (
	"bytes"
	"github.com/rcrowley/go-metrics"
	"github.com/redis/go-redis/v9"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"golang.org/x/net/context"
	"strings"
	"time"
	"tinyIM/config"
	"tinyIM/tools"
)

var RedisClient *redis.Client
var RedisSessionClient *redis.Client

func (logic *Logic) InitRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}

	// 可以使用不同的redis服务器
	RedisSessionClient = RedisClient
	return
}

func (logic *Logic) InitRpcServer() (err error) {
	var network, addr string
	rpcAddrList := strings.Split(config.Conf.Logic.LogicBase.RpcAddress, ",")
	for _, bind := range rpcAddrList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		logrus.Infof("logic start run at network: %s, addr: %s", network, addr)
		go logic.createRpcServer(network, addr)
	}
	return
}

func (logic *Logic) createRpcServer(network, addr string) {
	s := server.NewServer()
	logic.addRegistryPlugin(s, network, addr)

	if err := s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathLogic, new(RpcLogic), logic.ServerId); err != nil {
		logrus.Errorf("register error:%s", err.Error())
	}

	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func (logic *Logic) addRegistryPlugin(s *server.Server, network string, addr string) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: network + "@" + addr,
		EtcdServers:    []string{config.Conf.Common.CommonEtcd.Host},
		BasePath:       config.Conf.Common.CommonEtcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		logrus.Fatal(err)
	}
	s.Plugins.Add(r)
}

func (logic *Logic) getUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}
