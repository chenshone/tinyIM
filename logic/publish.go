package logic

import (
	"bytes"
	"strings"
	"tinyIM/config"
	"tinyIM/tools"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"golang.org/x/net/context"
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
			logrus.Errorf("InitLogicRpc ParseNetwork error : %s", err.Error())
			return
		}
		logrus.Infof("logic start run at network: %s, addr: %s", network, addr)
		go logic.createRpcServer(network, addr)
	}
	return
}

func (logic *Logic) createRpcServer(network, addr string) {
	s := server.NewServer()
	tools.AddRegistryPlugin(s, network, addr)

	if err := s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathLogic, new(RpcLogic), logic.ServerId); err != nil {
		logrus.Errorf("register error:%s", err.Error())
	}

	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func (logic *Logic) getUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getRoomUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getRoomOnlineCountKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomOnlinePrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}
