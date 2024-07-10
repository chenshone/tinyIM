package tools

import (
	"time"
	"tinyIM/config"

	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
)

func AddRegistryPlugin(s *server.Server, network string, addr string) {
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
