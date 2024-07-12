package task

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"tinyIM/config"
)

type Task struct{}

func New() *Task {
	return &Task{}
}

func (task *Task) Run() {
	cfg := config.Conf.Task

	runtime.GOMAXPROCS(cfg.TaskBase.CpuNum)

	//	1. init mq (use redis now)
	if err := task.InitQueueClient(); err != nil {
		logrus.Panicf("task init queue client fail, err: %s", err.Error())
	}

	//	2. init connect rpc client
	if err := task.InitConnectRpcClient(); err != nil {
		logrus.Panicf("task init connect rpc client fail, err: %s", err.Error())
	}

	//3. init task server
	task.GoPush()

}
