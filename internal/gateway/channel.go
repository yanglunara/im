package gateway

import (
	"bufio"
	"errors"
	"sync"
	"time"

	pb "github.com/yanglunara/im/api/protocol"
)

var (
	ProtoRead   = &ProtoRing{Proto: &pb.Proto{Op: int32(pb.Operation_ProtoReady)}}
	ProtoFinish = &ProtoRing{Proto: &pb.Proto{Op: int32(pb.Operation_ProtoFinish)}}
)

type Channel struct {
	Room     *Room
	CliProto Ring
	signal   chan *ProtoRing
	Writer   bufio.Writer
	Reader   bufio.Reader
	Next     *Channel
	Prev     *Channel

	Mid       int64
	Key       string
	IP        string
	watchOps  map[int32]struct{}
	mutex     sync.RWMutex
	ProtoRing *ProtoRing
	LastHB    time.Time
}

func NewChannel(cli, srv int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *ProtoRing, srv)
	c.watchOps = make(map[int32]struct{})
	return c
}

func (c *Channel) SetLastHB() {
	c.LastHB = time.Now()
}

func (c *Channel) SetProtoRing(pb *ProtoRing) {
	c.ProtoRing = pb
}

// Watch 监听
func (c *Channel) Watch(ops ...int32) {
	c.mutex.Lock()
	for _, op := range ops {
		c.watchOps[op] = struct{}{}
	}
	c.mutex.Unlock()
}

func (c *Channel) UnWatch(ops ...int32) {
	c.mutex.Lock()
	for _, op := range ops {
		delete(c.watchOps, op)
	}
	c.mutex.Unlock()
}

func (c *Channel) NeedPush(op int32) bool {
	c.mutex.RLock() // 使用 RLock 替代 Lock
	defer c.mutex.RUnlock()
	if _, ok := c.watchOps[op]; ok {
		return ok
	}
	return false
}

func (c *Channel) Push(p *ProtoRing) (err error) {
	select {
	case c.signal <- p:
	case <-time.After(time.Second): // 添加超时机制
		err = errors.New("signal channel full, msg dropped")
	}
	return
}

func (c *Channel) Ready() *ProtoRing {
	return <-c.signal
}

func (c *Channel) Signal() {
	c.signal <- ProtoRead
}

func (c *Channel) Close() {
	c.signal <- ProtoFinish
}
