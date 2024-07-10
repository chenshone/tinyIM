package connect

import (
	"context"
	"github.com/sirupsen/logrus"
	"tinyIM/proto"
)

type RpcConnect struct {
}

func (rpc *RpcConnect) Connect(req *proto.ConnectRequest) (uid int, err error) {
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

func (rpc *RpcConnect) DisConnect(req *proto.DisConnectRequest) (err error) {
	reply := &proto.DisConnectReply{}
	err = logicRpcClient.Call(context.Background(), "DisConnect", req, reply)
	if err != nil {
		logrus.Errorf("failed to call: %v", err)
	}
	return
}
