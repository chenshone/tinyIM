package connect

import "tinyIM/proto"

type Operator interface {
	Connect(req *proto.ConnectRequest) (int, error)
	Disconnect(req *proto.DisConnectRequest) error
}

type DefaultOperator struct {
}

// rpc call logic layer
func (o *DefaultOperator) Connect(req *proto.ConnectRequest) (uid int, err error) {
	rpcConnect := &RpcLogicClient{}
	uid, err = rpcConnect.Connect(req)
	return
}

// rpc call logic layer
func (o *DefaultOperator) Disconnect(req *proto.DisConnectRequest) error {
	rpcConnect := &RpcLogicClient{}
	return rpcConnect.DisConnect(req)
}
