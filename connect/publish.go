package connect

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"tinyIM/config"
	"tinyIM/tools"

	"github.com/rpcxio/libkv/store"
	etcdV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
)

var logicRpcClient client.XClient
var once sync.Once

func (c *Connect) InitLogicRpcClient() error {
	etcdCfgOpt := &store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: time.Duration(config.Conf.Common.CommonEtcd.ConnectionTimeout) * time.Second,
		Bucket:            "",
		PersistConnection: true,
		Username:          config.Conf.Common.CommonEtcd.UserName,
		Password:          config.Conf.Common.CommonEtcd.Password,
	}

	once.Do(func() {
		d, e := etcdV3.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath, config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			true,
			etcdCfgOpt,
		)
		if e != nil {
			logrus.Fatalf("init connect rpc etcd discovery client fail:%s", e.Error())
		}
		logicRpcClient = client.NewXClient(config.Conf.Common.CommonEtcd.ServerPathLogic, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	})
	if logicRpcClient == nil {
		return errors.New("get rpc client nil")
	}
	return nil
}

func (c *Connect) InitWebSocketRpcServer() (err error) {
	var network, addr string
	connRpcAddr := strings.Split(config.Conf.Connect.ConnectRpcAddressWebSockts.Address, ",")
	for _, bind := range connRpcAddr {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Errorf("InitConnectWebsocketRpcServer ParseNetwork error : %s", err)
			return
		}
		logrus.Infof("Connect server start run at-->%s:%s", network, addr)

		go c.createWebSocketRpcServer(network, addr)
	}
	return
}

func (c *Connect) createWebSocketRpcServer(network, addr string) {
	s := server.NewServer()
	tools.AddRegistryPlugin(s, network, addr)

	s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathConnect, &RpcConnect{}, fmt.Sprintf("serverId=%s&serverType=ws", c.ServerId))

	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})

	s.Serve(network, addr)
}
