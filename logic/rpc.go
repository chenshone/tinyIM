package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"tinyIM/config"
	"tinyIM/logic/dao"
	"tinyIM/proto"
	"tinyIM/tools"
)

type RpcLogic struct {
}

func (rpc *RpcLogic) Register(ctx context.Context, args *proto.RegisterRequest, reply *proto.RegisterReply) (err error) {
	reply.Code = config.FailReplyCode
	u := &dao.User{}
	_, ok := u.CheckByUsername(args.Name)
	if ok {
		return errors.New("username already exists")
	}
	u.Username, u.Password = args.Name, args.Password
	userId, err := u.Add()
	if err != nil {
		logrus.Infof("register error: %s", err.Error())
		return err
	}
	if userId == 0 {
		return errors.New("register userId empty")
	}

	randomToken := tools.GetRandomToken(32)
	sessionId := tools.CreateSessionId(randomToken)
	userData := make(map[string]any)
	userData["userId"] = userId
	userData["username"] = args.Name

	// set session to redis
	// 执行redis事务，保证原子性
	RedisSessionClient.Do(ctx, "MULTI")
	RedisSessionClient.HSet(ctx, sessionId, userData)
	RedisSessionClient.Expire(ctx, sessionId, 86400*time.Second) // 24h
	err = RedisSessionClient.Do(ctx, "EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}

	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randomToken
	return
}

func (rpc *RpcLogic) Login(ctx context.Context, args *proto.LoginRequest, reply *proto.LoginResponse) (err error) {
	reply.Code = config.FailReplyCode

	// check username and password
	u := &dao.User{}
	username, password := args.Name, args.Password
	data, ok := u.CheckByUsername(username)
	if !ok || password != data.Password {
		return errors.New("username or password error")
	}

	// Check if you are logged in
	loginSessionId := tools.GetSessionIdByUserId(data.Id)
	token, _ := RedisSessionClient.Get(ctx, loginSessionId).Result()
	if token != "" {
		// del old login session
		oldSessionId := tools.CreateSessionId(token)
		if err := RedisSessionClient.Del(ctx, oldSessionId).Err(); err != nil {
			return errors.New("logout user fail! token is:" + token)
		}
	}

	// set session to redis
	randomToken := tools.GetRandomToken(32)
	sessionId := tools.CreateSessionId(randomToken)
	userData := make(map[string]any)
	userData["userId"] = data.Id
	userData["userName"] = data.Username
	RedisSessionClient.Do(ctx, "MULTI")
	RedisSessionClient.HSet(ctx, sessionId, userData)
	RedisSessionClient.Expire(ctx, sessionId, 86400*time.Second)
	RedisSessionClient.Set(ctx, loginSessionId, randomToken, 86400*time.Second)
	err = RedisSessionClient.Do(ctx, "EXEC").Err()
	if err != nil {
		logrus.Infof("login set redis token fail!")
		return
	}
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randomToken
	return
}

func (rpc *RpcLogic) GetUserInfoByUserId(ctx context.Context, args *proto.GetUserInfoRequest, reply *proto.GetUserInfoResponse) (err error) {
	reply.Code = config.FailReplyCode
	u := &dao.User{}
	username := u.GetUserNameByUserId(args.UserId)
	reply.UserId = args.UserId
	reply.Username = username
	reply.Code = config.SuccessReplyCode
	return
}

func (rpc *RpcLogic) CheckAuth(ctx context.Context, args *proto.CheckAuthRequest, reply *proto.CheckAuthResponse) (err error) {
	reply.Code = config.FailReplyCode
	session := tools.GetSessionName(args.AuthToken)
	userDataMap, err := RedisSessionClient.HGetAll(ctx, session).Result()
	if err != nil {
		logrus.Infof("check auth fail!,authToken is:%s", args.AuthToken)
		return err
	}
	if len(userDataMap) == 0 {
		logrus.Infof("no this user session,authToken is:%s", args.AuthToken)
		return
	}
	intUserId, _ := strconv.Atoi(userDataMap["userId"])
	reply.UserId = intUserId
	userName, _ := userDataMap["userName"]
	reply.Code = config.SuccessReplyCode
	reply.Username = userName
	return
}

func (rpc *RpcLogic) Logout(ctx context.Context, args *proto.LogoutRequest, reply *proto.LogoutResponse) (err error) {
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	session := tools.GetSessionName(authToken)

	userDataMap, err := RedisSessionClient.HGetAll(ctx, session).Result()
	if err != nil {
		logrus.Infof("check auth fail!,authToken is:%s", authToken)
		return err
	}
	if len(userDataMap) == 0 {
		logrus.Infof("no this user session,authToken is:%s", authToken)
		return
	}
	intUserId, _ := strconv.Atoi(userDataMap["userId"])
	sessIdMap := tools.GetSessionIdByUserId(intUserId)
	//del sess_map like sess_map_1
	err = RedisSessionClient.Del(ctx, sessIdMap).Err()
	if err != nil {
		logrus.Infof("logout del sess map error:%s", err.Error())
		return err
	}
	//del serverId
	logic := &Logic{}
	serverIdKey := logic.getUserKey(fmt.Sprintf("%d", intUserId))
	err = RedisSessionClient.Del(ctx, serverIdKey).Err()
	if err != nil {
		logrus.Infof("logout del server id error:%s", err.Error())
		return err
	}
	err = RedisSessionClient.Del(ctx, session).Err()
	if err != nil {
		logrus.Infof("logout error:%s", err.Error())
		return err
	}
	reply.Code = config.SuccessReplyCode
	return
}
