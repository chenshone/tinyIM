package handler

import (
	"tinyIM/api/rpc"
	"tinyIM/proto"
	"tinyIM/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

	code, authToken, msg := rpc.LogicInstance.Login(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "login success", authToken)
}

type UserRegister struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var userRegister UserRegister
	if err := c.ShouldBindBodyWith(&userRegister, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	req := &proto.RegisterRequest{
		Name:     userRegister.Username,
		Password: tools.HashWithSalt(userRegister.Password),
	}
	code, authToken, msg := rpc.LogicInstance.Register(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "register success", authToken)
}

type UserCheckAuth struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckAuth(c *gin.Context) {
	var userCheckAuth UserCheckAuth
	if err := c.ShouldBindBodyWith(&userCheckAuth, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := userCheckAuth.AuthToken
	req := &proto.CheckAuthRequest{
		AuthToken: authToken,
	}
	code, userId, userName := rpc.LogicInstance.CheckAuth(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "auth fail")
		return
	}
	var jsonData = map[string]interface{}{
		"userId":   userId,
		"userName": userName,
	}
	tools.SuccessWithMsg(c, "auth success", jsonData)
}

type UserLogout struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Logout(c *gin.Context) {
	var userLogout UserLogout
	if err := c.ShouldBindBodyWith(&userLogout, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := userLogout.AuthToken
	req := &proto.LogoutRequest{
		AuthToken: authToken,
	}
	code := rpc.LogicInstance.Logout(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "logout fail!")
		return
	}
	tools.SuccessWithMsg(c, "logout ok!", nil)
}
