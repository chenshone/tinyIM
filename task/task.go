package task

import (
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

	/**
		  TODO:
		  1. init mq (use redis now)
		  2. init connect rpc client
	      3. init task server
	*/
}
