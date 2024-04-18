package gateway

import (
	"sync"

	pb "github.com/yanglunara/im/api/protocol"
)

type Room struct {
	ID        string
	rLock     sync.RWMutex
	next      *Channel
	drop      bool
	Online    int32
	AllOnline int32
}

func NewRoom(id string) *Room {
	return &Room{
		ID:     id,
		drop:   false,
		next:   nil,
		Online: 0,
	}
}

func (r *Room) Push(p *pb.Proto) {
	r.rLock.RLock()
	defer r.rLock.RUnlock()
	for c := r.next; c != nil; c = c.Next {
		_ = c.Push(p)
	}
}

func (r *Room) Del(ch *Channel) bool {
	r.rLock.Lock()
	defer r.rLock.Unlock()
	switch {
	case ch.Next != nil:
		ch.Next.Prev = ch.Prev
	case ch.Prev != nil:
		ch.Prev.Next = ch.Next
	default:
		r.next = ch.Next
	}
	ch.Next = nil
	ch.Prev = nil
	r.Online--
	r.drop = r.Online == 0
	return r.drop
}

func (r *Room) Close() {
	r.rLock.Lock()
	defer r.rLock.Unlock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
}

func (r *Room) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	defer r.rLock.Lock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.Online++
		return
	}
	return ErrRoomDroped
}

func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
