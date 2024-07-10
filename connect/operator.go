package connect

import "tinyIM/proto"

type Operator interface {
	Connect(*proto.ConnectRequest) (int, error)
	Disconnect(*proto.DisConnectRequest) error
}

type DefaultOperator struct {
}

// rpc call logic layer
func (o *DefaultOperator) Connect(req *proto.ConnectRequest) (uid int, err error) {
	rpcConnect := &RpcConnect{}
	uid, err = rpcConnect.Connect(req)
	return
}

// rpc call logic layer
func (o *DefaultOperator) DisConnect(req *proto.DisConnectRequest) (err error) {
	rpcConnect := &RpcConnect{}
	err = rpcConnect.DisConnect(req)
	return
}
