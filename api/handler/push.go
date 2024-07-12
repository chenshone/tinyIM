package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"tinyIM/api/rpc"
	"tinyIM/config"
	"tinyIM/proto"
	"tinyIM/tools"
)

type FormPush struct {
	Msg       string `form:"msg" json:"msg" binding:"required"`
	ToUserId  string `form:"toUserId" json:"toUserId" binding:"required"`
	RoomId    int    `form:"roomId" json:"roomId" binding:"required"`
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Push(c *gin.Context) {
	var formPush FormPush
	if err := c.ShouldBindBodyWith(&formPush, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}

	authToken := formPush.AuthToken
	msg := formPush.Msg
	toUserId, _ := strconv.Atoi(formPush.ToUserId)

	getUsernameReq := &proto.GetUserInfoRequest{
		UserId: toUserId,
	}

	code, toUsername := rpc.LogicInstance.GetUserNameByUserId(getUsernameReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc fail get friend userName")
		return
	}

	checkAuthReq := &proto.CheckAuthRequest{AuthToken: authToken}
	code, fromUserId, fromUserName := rpc.LogicInstance.CheckAuth(checkAuthReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc fail get self info")
		return
	}
	roomId := formPush.RoomId
	req := &proto.Send{
		Msg:          msg,
		FromUserId:   fromUserId,
		FromUserName: fromUserName,
		ToUserId:     toUserId,
		ToUserName:   toUsername,
		RoomId:       roomId,
		Op:           config.OpSingleSend,
	}
	code, replyMsg := rpc.LogicInstance.PushSingle(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, replyMsg)
		return
	}
	tools.SuccessWithMsg(c, "ok", nil)
}
