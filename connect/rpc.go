package connect

import (
	"context"
	"errors"
	"tinyIM/config"
	"tinyIM/proto"

	"github.com/sirupsen/logrus"
)

type RpcConnect struct {
}

func (rpc *RpcConnect) PushSingleMsg(ctx context.Context, req *proto.PushMsgRequest, reply *proto.SuccessReply) error {
	logrus.Infof("rpc PushSingleMsg: %v", req)
	if req == nil {
		errMsg := "rpc PushSingleMsg() args is nil"
		logrus.Error(errMsg)
		return errors.New(errMsg)
	}
	bucket := DefaultServer.GetBucket(req.UserId)
	if channel := bucket.GetChannel(req.UserId); channel != nil {
		if err := channel.Push(&req.Msg); err != nil {
			logrus.Errorf("DefaultServer Channel Push err, args: %v", req)
			return err
		}
	}

	reply.Code = config.SuccessReplyCode
	reply.Msg = config.SuccessReplyMsg
	logrus.Infof("successReply: %v", reply)
	return nil
}

func (rpc *RpcConnect) PushRoomMsg(ctx context.Context, req *proto.PushRoomMsgRequest, reply *proto.SuccessReply) error {
	logrus.Infof("PushRoomMsg msg %+v", req)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(req)
	}
	reply.Code = config.SuccessReplyCode
	reply.Msg = config.SuccessReplyMsg
	return nil
}

func (rpc *RpcConnect) PushRoomCount(ctx context.Context, req *proto.PushRoomMsgRequest, reply *proto.SuccessReply) (err error) {
	logrus.Infof("PushRoomCount msg %v", req)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(req)
	}
	reply.Code = config.SuccessReplyCode
	reply.Msg = config.SuccessReplyMsg
	return
}

func (rpc *RpcConnect) PushRoomInfo(ctx context.Context, req *proto.PushRoomMsgRequest, reply *proto.SuccessReply) (err error) {
	logrus.Infof("connect,PushRoomInfo msg %+v", req)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(req)
	}
	reply.Code = config.SuccessReplyCode
	reply.Msg = config.SuccessReplyMsg
	return
}
