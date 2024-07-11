package connect

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
	"tinyIM/proto"
	"tinyIM/tools"
)

type Server struct {
	Buckets   []*Bucket
	Option    ServerOption
	bucketIdx uint32
	op        Operator
}

type ServerOption struct {
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
}

func NewServer(b []*Bucket, op Operator, opt ServerOption) *Server {
	return &Server{
		Buckets:   b,
		Option:    opt,
		bucketIdx: uint32(len(b)),
		op:        op,
	}
}

// reduce lock competition, use google city hash insert to different bucket
func (s *Server) GetBucket(userId int) *Bucket {
	userIdStr := fmt.Sprintf("%d", userId)
	idx := tools.CityHash32([]byte(userIdStr), uint32(len(userIdStr))) % s.bucketIdx
	return s.Buckets[idx]
}

func (s *Server) writePump(ch *Channel) {
	ticker := time.NewTicker(s.Option.PingPeriod)
	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ch.broadcast:
			ch.conn.SetWriteDeadline(time.Now().Add(s.Option.WriteWait))
			if !ok {
				logrus.Warn("channel's broadcast chan closed")
				ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			/**
			当你调用 conn.NextWriter(messageType) 时，WebSocket 库会开始一个新的消息帧，返回的 writer 对象会将所有写入的数据添加到该帧中。这个 writer 对象在你调用 Close() 方法之前不会将数据实际发送出去。调用 Close() 后，数据帧会被完整地发送到对端。
			*/
			writer, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logrus.Warnf(" ch.conn.NextWriter err :%s", err.Error())
				return
			}
			logrus.Infof("writePump write message :%s", message.Body)
			writer.Write(message.Body)
			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			//heartbeat，if ping error will exit and close current websocket conn
			ch.conn.SetWriteDeadline(time.Now().Add(s.Option.WriteWait))
			logrus.Infof("websocket.PingMessage :%v", websocket.PingMessage)
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) readPump(ch *Channel, c *Connect) {
	defer func() {
		defer ch.conn.Close()
		logrus.Info("start exec disConnect ...")
		if ch.Room == nil || ch.userId == 0 {
			logrus.Info("ch.Room == nil || ch.userId == 0")
			return
		}
		logrus.Info("exec disConnect ...")
		disConnReq := &proto.DisConnectRequest{
			RoomId: ch.Room.Id,
			UserId: ch.userId,
		}
		s.GetBucket(ch.userId).DelChannel(ch)
		if err := s.op.Disconnect(disConnReq); err != nil {
			logrus.Warnf("disconnect err :%s", err.Error())
		}
	}()

	ch.conn.SetReadLimit(s.Option.MaxMessageSize)
	ch.conn.SetReadDeadline(time.Now().Add(s.Option.PongWait))
	// 因为connect layer 的 ws write中是主动发生ping的，所以接受pong的响应
	ch.conn.SetPongHandler(func(string) error {
		// 忽略pong携带的数据, 更新读取截至时间
		ch.conn.SetReadDeadline(time.Now().Add(s.Option.PongWait))
		return nil
	})

	for {
		_, msg, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("readPump ReadMessage err:%s", err.Error())
				return
			}
		}
		if msg == nil {
			return
		}

		connReq := &proto.ConnectRequest{}
		logrus.Infof("readPump read message :%s", msg)
		if err = json.Unmarshal(msg, connReq); err != nil {
			logrus.Errorf("message struct %+v", connReq)
		}
		if connReq.AuthToken == "" {
			logrus.Errorf("s.operator.Connect no authToken")
			return
		}
		connReq.ServerId = c.ServerId
		userId, err := s.op.Connect(connReq)
		if err != nil {
			logrus.Errorf("s.operator.Connect err:%s", err.Error())
			return
		}
		if userId == 0 {
			logrus.Error("Invalid AuthToken ,userId empty")
			return
		}

		logrus.Infof("websocket rpc call return userId:%d,RoomId:%d", userId, connReq.RoomId)
		b := s.GetBucket(userId)
		//	insert to bucket
		if err = b.PutChannel(userId, connReq.RoomId, ch); err != nil {
			logrus.Errorf("conn will close...\n err: %s", err.Error())
			// 这里是否可以直接return
			ch.conn.Close()
		}
	}
}
