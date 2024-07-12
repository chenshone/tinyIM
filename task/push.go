package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math/rand/v2"
	"tinyIM/config"
	"tinyIM/proto"
)

type PushParams struct {
	ServerId string
	UserId   int
	Msg      []byte
	RoomId   int
}

var pushChannel []chan *PushParams

func init() {
	pushChannel = make([]chan *PushParams, config.Conf.Task.TaskBase.PushChan)
}

func (task *Task) GoPush() {
	for i := 0; i < len(pushChannel); i++ {
		pushChannel[i] = make(chan *PushParams, config.Conf.Task.TaskBase.PushChanSize)
		go task.processSinglePush(pushChannel[i])
	}
}

func (task *Task) processSinglePush(ch chan *PushParams) {
	var arg *PushParams
	for {
		arg = <-ch
		// // 因为用户会绑定到唯一的connect server， 所以当这个server关闭时，该用户的聊天信息会丢失(因为这个serverid失效)
		ConnectRpc.pushSingleToConnect(arg.ServerId, arg.UserId, arg.Msg)
	}
}

func (task *Task) Push(msg string) {
	m := &proto.RedisMsg{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		logrus.Infof("task server, unmarshal msg err:%v", err)
	}
	logrus.Infof("push msg info room: %d,op is: %d", m.RoomId, m.Op)
	switch m.Op {
	case config.OpSingleSend:
		pushChannel[rand.Int()%config.Conf.Task.TaskBase.PushChan] <- &PushParams{
			ServerId: m.ServerId,
			UserId:   m.UserId,
			Msg:      m.Msg,
		}
	}
}
