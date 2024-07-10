package connect

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"tinyIM/config"
)

type Connect struct {
	ServerId string
}

func New() *Connect {
	return &Connect{}
}

func (c *Connect) Run() {
	cfg := config.Conf.Connect

	runtime.GOMAXPROCS(cfg.ConnectBucket.CpuNum)

	/**
	TODO:
	init logic layer Rpc client
	init connect layer rpc server
	*/
	//	init logic layer Rpc client
	if err := c.InitLogicRpcClient(); err != nil {
		logrus.Panicf("init logic layer Rpc client failed, err: %v", err)
	}
	
}
