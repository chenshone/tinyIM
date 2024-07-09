package logic

import (
	"fmt"
	"runtime"
	"tinyIM/config"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Logic struct {
	ServerId string
}

func New() *Logic {
	return &Logic{}
}

func (logic *Logic) Run() {
	logicCfg := config.Conf.Logic

	runtime.GOMAXPROCS(logicCfg.LogicBase.CpuNum)
	logic.ServerId = fmt.Sprintf("logic-%s", uuid.New().String())

	err := logic.InitRedisClient()
	if err != nil {
		logrus.Panicf("init redis client failed, err: %s", err.Error())
	}

	err = logic.InitRpcServer()
	if err != nil {
		logrus.Panicf("init rpc server failed, err: %s", err.Error())
	}

}
