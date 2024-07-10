package connect

import (
	"sync"
	"sync/atomic"
	"tinyIM/proto"
)

type Bucket struct {
	lock        sync.RWMutex // protect the channels for chs
	chs         map[int]*Channel
	bucketOpt   *BucketOption
	rooms       map[int]*Room // bucket room channels
	routines    []chan *proto.PushRoomMsgRequest
	routinesNum uint64
	broadcast   chan []byte
}

type BucketOption struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}

func NewBucket(opt *BucketOption) *Bucket {
	b := &Bucket{
		chs:       make(map[int]*Channel, opt.ChannelSize),
		bucketOpt: opt,
		routines:  make([]chan *proto.PushRoomMsgRequest, opt.RoutineAmount),
		rooms:     make(map[int]*Room, opt.RoomSize),
	}

	for i := uint64(0); i < opt.RoutineAmount; i++ {
		c := make(chan *proto.PushRoomMsgRequest, opt.RoutineSize)
		b.routines[i] = c
		go b.PushMsg2Room(c)
	}

	return b

}

func (b *Bucket) PushMsg2Room(ch chan *proto.PushRoomMsgRequest) {
	for {
		var (
			arg  *proto.PushRoomMsgRequest
			room *Room
		)
		arg = <-ch
		if room = b.GetRoom(arg.RoomId); room != nil {
			room.PushMsg(&arg.Msg)
		}
	}
}

func (b *Bucket) GetRoom(rid int) *Room {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.rooms[rid]
}

func (b *Bucket) PutChannel(uid, rid int, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.lock.Lock()
	if rid != NoRoom {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room
	}
	ch.userId = uid
	b.chs[uid] = ch
	b.lock.Unlock()

	if room != nil {
		err = room.Put(ch)
	}
	return
}

func (b *Bucket) DelChannel(ch *Channel) {
	var (
		ok   bool
		room *Room
	)
	b.lock.Lock()
	defer b.lock.Unlock()
	if ch, ok = b.chs[ch.userId]; ok {
		room = b.chs[ch.userId].Room
		//delete from bucket
		delete(b.chs, ch.userId)
	}
	if room != nil && room.Del(ch) {
		// if room empty delete,will mark room.drop is true
		if room.drop {
			delete(b.rooms, room.Id)
		}
	}
}

func (b *Bucket) GetChannel(uid int) *Channel {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.chs[uid]
}

func (b *Bucket) BroadcastRoom(req *proto.PushRoomMsgRequest) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.bucketOpt.RoutineAmount
	b.routines[num] <- req
}
