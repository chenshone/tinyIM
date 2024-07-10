package connect

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"tinyIM/proto"
)

const NoRoom = -1

type Room struct {
	Id          int
	OnlineCount int // 在线人数
	lock        sync.RWMutex
	drop        bool     // 该房间是否删除
	next        *Channel // channel是一个双向链表
}

func NewRoom(roomId int) *Room {
	return &Room{
		Id: roomId,
	}
}

func (r *Room) Put(ch *Channel) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.OnlineCount++
	} else {
		return errors.New("room has been closed")
	}
	return nil
}

func (r *Room) Del(ch *Channel) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	if ch.Next != nil {
		// 当前channel不是最后一个
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// 当前channel不是第一个
		ch.Prev.Next = ch.Next
	} else {
		// 当前channel是第一个
		r.next = ch.Next
	}
	r.OnlineCount--
	r.drop = false
	if r.OnlineCount <= 0 {
		r.drop = true
	}
	return r.drop
}

func (r *Room) PushMsg(msg *proto.Msg) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for ch := r.next; ch != nil; ch = ch.Next {
		if err := ch.Push(msg); err != nil {
			logrus.Infof("push msg to channel failed: %v", err)
		}
	}
}
