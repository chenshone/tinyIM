package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSuccess      = 0
	CodeFail         = 1
	CodeUnknownError = -1
	CodeSessionError = 40000
)

var MsgCodeMap = map[int]string{
	CodeUnknownError: "unKnow error",
	CodeSuccess:      "success",
	CodeFail:         "fail",
	CodeSessionError: "Session error",
}

func SuccessWithMsg(c *gin.Context, msg interface{}, data interface{}) {
	ResponseWithCode(c, CodeSuccess, msg, data)
}

func FailWithMsg(c *gin.Context, msg interface{}) {
	ResponseWithCode(c, CodeFail, msg, nil)
}

func ResponseWithCode(c *gin.Context, msgCode int, msg any, data any) {
	if msg == nil {
		if val, ok := MsgCodeMap[msgCode]; ok {
			msg = val
		} else {
			msg = MsgCodeMap[CodeUnknownError]
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    msgCode,
		"message": msg,
		"data":    data,
	})
}
