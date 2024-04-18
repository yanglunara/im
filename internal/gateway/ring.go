package gateway

import (
	"errors"
	"math"

	"github.com/yanglunara/im/api/protocol"
)

var (
	ErrRingEmpty            = errors.New("ring buffer empty")
	ErrRingFull             = errors.New("ring buffer full")
	ErrSignalFullMsgDropped = errors.New("signal channel full, msg dropped")
	ErrRoomDroped           = errors.New("room droped")
)

type Ring struct {
	rp   uint64
	num  uint64
	mask uint64
	wp   uint64
	data []protocol.Proto
}

func NewRing(num uint64) *Ring {
	r := new(Ring)
	r.init(num)
	return r
}

func (r *Ring) init(num uint64) {
	if num&(num-1) != 0 {
		num = 1 << uint(math.Ceil(math.Log2(float64(num))))
	}
	r.data = make([]protocol.Proto, num)
	r.num = num
	r.mask = r.num - 1
}

func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

func (r *Ring) Get() (p *protocol.Proto, err error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	return &r.data[r.rp&r.mask], nil
}

func (r *Ring) Put() (p *protocol.Proto, err error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	return &r.data[r.wp&r.mask], nil
}

func (r *Ring) IncrementReadPointer() {
	r.rp++
}

func (r *Ring) IncrementWritePointer() {
	r.wp++
}

func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
}
