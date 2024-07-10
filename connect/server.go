package connect

import (
	"fmt"
	"time"
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
