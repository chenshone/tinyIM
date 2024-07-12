package task

import (
	"context"
	"github.com/sirupsen/logrus"
	"tinyIM/config"
	"tinyIM/proto"
	"tinyIM/tools"
)

var ConnectRpc = &RpcConnect{}

type RpcConnect struct {
}

func (rpc *RpcConnect) pushSingleToConnect(serverId string, userId int, msg []byte) {
	logrus.Infof("pushSingleToConnect Body %s", string(msg))
	pushMsgReq := &proto.PushMsgRequest{
		UserId: userId,
		Msg: proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpSingleSend,
			SeqId:     tools.GetSnowflakeId(),
			Body:      msg,
		},
	}
	reply := &proto.SuccessReply{}
	connRpcClient, err := RClient.GetRpcClientByServerId(serverId)
	if err != nil {
		logrus.Infof("get rpc client err %v", err)
		return
	}

	err = connRpcClient.Call(context.Background(), "PushSingleMsg", pushMsgReq, reply)

	if err != nil {
		logrus.Infof("pushSingleToConnect Call err %v", err)
		return
	}
	logrus.Infof("reply %s", reply.Msg)
}
