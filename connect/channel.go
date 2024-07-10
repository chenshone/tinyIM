package connect

import (
	"github.com/gorilla/websocket"
	"tinyIM/proto"
)

type Channel struct {
	Room      *Room
	Prev      *Channel
	Next      *Channel
	broadcast chan *proto.Msg
	userId    int
	conn      *websocket.Conn
}

func NewChannel(size int) *Channel {
	c := &Channel{}
	c.broadcast = make(chan *proto.Msg, size)
	return c
}

func (ch *Channel) Push(msg *proto.Msg) error {
	// 尝试将消息 msg 推送到通道 ch.broadcast 中。如果通道已经满了，消息将被丢弃
	select {
	case ch.broadcast <- msg:
	default:
	}
	return nil
}
