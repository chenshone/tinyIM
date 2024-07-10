package connect

import (
	"fmt"
	"runtime"
	"time"
	"tinyIM/config"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var DefaultServer *Server

type Connect struct {
	ServerId string
}

func New() *Connect {
	return &Connect{}
}

func (c *Connect) Run() {
	cfg := config.Conf.Connect

	runtime.GOMAXPROCS(cfg.ConnectBucket.CpuNum)

	//	init logic layer Rpc client
	if err := c.InitLogicRpcClient(); err != nil {
		logrus.Panicf("init logic layer Rpc client failed, err: %v", err)
	}

	//	init connect layer rpc server
	buckets := make([]*Bucket, cfg.ConnectBucket.CpuNum)
	for i := 0; i < cfg.ConnectBucket.CpuNum; i++ {
		buckets[i] = NewBucket(&BucketOption{
			ChannelSize:   cfg.ConnectBucket.Channel,
			RoomSize:      cfg.ConnectBucket.Room,
			RoutineAmount: cfg.ConnectBucket.RoutineAmount,
			RoutineSize:   cfg.ConnectBucket.RoutineSize,
		})
	}

	operator := &DefaultOperator{}

	DefaultServer = NewServer(buckets, operator, ServerOption{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      54 * time.Second,
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		BroadcastSize:   512,
	})

	c.ServerId = fmt.Sprintf("%s:%s", "ws", uuid.New().String())

	if err := c.InitWebSocketRpcServer(); err != nil {
		logrus.Panicf("InitConnectWebsocketRpcServer Fatal error: %s \n", err.Error())
	}

	//start Connect layer server handler persistent connection
	c.InitWebSocket()
}
