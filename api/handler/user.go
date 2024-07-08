package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"tinyIM/proto"
	"tinyIM/tools"
)

type UserLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var user UserLogin
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.LoginRequest{
		Name:     user.Username,
		Password: user.Password,
	}
}
