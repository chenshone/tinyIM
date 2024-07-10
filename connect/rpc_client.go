package connect

import (
	"context"
	"tinyIM/proto"

	"github.com/sirupsen/logrus"
)

type RpcLogicClient struct {
}

func (rpc *RpcLogicClient) Connect(req *proto.ConnectRequest) (uid int, err error) {
	reply := &proto.ConnectReply{}
	err = logicRpcClient.Call(context.Background(), "Connect", req, reply)
	if err != nil {
		logrus.Errorf("logic rpc call Connect failed, err: %v", err)
		return
	}
	uid = reply.UserId
	logrus.Infof("logic rpc call Connect success, uid: %d", uid)
	return
}

func (rpc *RpcLogicClient) DisConnect(req *proto.DisConnectRequest) (err error) {
	reply := &proto.DisConnectReply{}
	err = logicRpcClient.Call(context.Background(), "DisConnect", req, reply)
	if err != nil {
		logrus.Errorf("failed to call: %v", err)
	}
	return
}
