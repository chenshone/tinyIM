package connect

import (
	"sync"
	"tinyIM/proto"
)

type Bucket struct {
	lock        sync.RWMutex // protect the channels for chs
	chs         map[int]*Channel
	bucketOpt   BucketOption
	rooms       map[int]*Room // bucket room channels
	routines    []chan *proto.PushRoomCountRequest
	routinesNum uint64
	broadcast   chan []byte
}

type BucketOption struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}
