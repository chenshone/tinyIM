package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"tinyIM/config"
	"tinyIM/data/db"
	"tinyIM/logic/dao"
	"tinyIM/proto"
	"tinyIM/tools"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
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

	txn := db.GetDb().Begin()
	userId, err := u.Add(txn)
	if err != nil {
		txn.Rollback()
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
	_, err = RedisSessionClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		RedisSessionClient.HSet(ctx, sessionId, userData)
		RedisSessionClient.Expire(ctx, sessionId, 86400*time.Second) // 24h
		return nil
	})
	if err != nil {
		txn.Rollback()
		logrus.Infof("register set redis token fail!")
		return err
	}

	txn.Commit()
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
	userData["username"] = data.Username

	_, err = RedisSessionClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		RedisSessionClient.HSet(ctx, sessionId, userData)
		RedisSessionClient.Expire(ctx, sessionId, 86400*time.Second)
		RedisSessionClient.Set(ctx, loginSessionId, randomToken, 86400*time.Second)
		return nil
	})
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
	reply.Code = config.SuccessReplyCode
	reply.Username = userDataMap["username"]
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

func (rpc *RpcLogic) Connect(ctx context.Context, args *proto.ConnectRequest, reply *proto.ConnectReply) (err error) {
	if args == nil {
		logrus.Errorf("logic,connect args empty")
		return errors.New("logic,connect args empty")
	}

	logrus.Infof("logic,authToken is:%s", args.AuthToken)
	logic := &Logic{}

	key := tools.GetSessionName(args.AuthToken)
	userInfo, err := RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		logrus.Infof("RedisCli HGetAll key :%s , err:%s", key, err.Error())
		return
	}
	if len(userInfo) == 0 {
		logrus.Infof("no this user session,authToken is:%s", args.AuthToken)
		return errors.New("no this user session")
	}

	reply.UserId, _ = strconv.Atoi(userInfo["userId"])
	roomUserKey := logic.getRoomUserKey(strconv.Itoa(args.RoomId))
	if reply.UserId != 0 {
		userKey := logic.getUserKey(strconv.Itoa(reply.UserId))
		logrus.Infof("logic redis set userKey:%s, serverId : %s", userKey, args.ServerId)
		validTime := config.RedisBaseValidTime * time.Second // 24h
		err = RedisClient.Set(ctx, userKey, args.ServerId, validTime).Err()
		if err != nil {
			logrus.Warnf("logic set err:%s", err)
		}
		if RedisClient.HGet(ctx, roomUserKey, strconv.Itoa(reply.UserId)).Val() == "" {
			// add curr user to room
			RedisClient.HSet(ctx, roomUserKey, strconv.Itoa(reply.UserId), userInfo["username"])
			RedisClient.Incr(ctx, logic.getRoomOnlineCountKey(strconv.Itoa(args.RoomId)))
		}
	}
	logrus.Infof("logic rpc userId:%d", reply.UserId)
	return
}

func (rpc *RpcLogic) DisConnect(ctx context.Context, args *proto.DisConnectRequest, reply *proto.DisConnectReply) (err error) {
	logic := &Logic{}
	roomUserKey := logic.getRoomUserKey(strconv.Itoa(args.RoomId))
	roomOnlineCountKey := logic.getRoomOnlineCountKey(strconv.Itoa(args.RoomId))

	if args.RoomId > 0 {
		// room user count - 1
		count, _ := RedisClient.Get(ctx, roomOnlineCountKey).Int()
		if count > 0 {
			_, _ = RedisClient.Decr(ctx, roomOnlineCountKey).Result()
		}
	}

	if args.UserId > 0 {
		// remove user from room
		err = RedisClient.HDel(ctx, roomUserKey, fmt.Sprintf("%d", args.UserId)).Err()
		if err != nil {
			logrus.Warnf("HDel getRoomUserKey err : %s", err)
		}
	}

	// TODO: 通知房间中其他user，该user退出
	return
}
