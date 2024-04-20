package gateway

import (
	"fmt"
	"sync"
	"sync/atomic"

	pb "github.com/yanglunara/im/api/gateway"
	conf "github.com/yanglunara/im/internal/conf/gateway"
)

// Bucket 结构体
type Bucket struct {
	c           *conf.Bucket                // 配置信息
	cLock       sync.RWMutex                // 读写锁
	channels    map[string]*Channel         // 通道集合
	rooms       map[string]*Room            // 房间集合
	routines    []chan *pb.BroadcastRoomReq // 广播房间请求集合
	routinesNum uint64                      // 广播房间数量

	ipCnts map[string]int32 // IP计数器
}

// NewBucket 创建新的Bucket
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.channels = make(map[string]*Channel, c.Channel)               // 初始化通道集合
	b.ipCnts = make(map[string]int32)                               // 初始化IP计数器
	b.c = c                                                         // 设置配置信息
	b.rooms = make(map[string]*Room, c.Room)                        // 初始化房间集合
	b.routines = make([]chan *pb.BroadcastRoomReq, c.RoutineAmount) // 初始化广播房间请求集合
	for i := uint64(0); i < c.RoutineAmount; i++ {                  // 创建广播房间请求
		c := make(chan *pb.BroadcastRoomReq, c.RoutineSize)
		b.routines[i] = c
		go b.roomProc(c) // 启动房间处理协程
	}
	return
}

// CountChannel 计算通道数量
func (b *Bucket) CountChannel() int {
	return len(b.channels)
}

// CountRoom 计算房间数量
func (b *Bucket) CountRoom() int {
	return len(b.rooms)
}

// CountRooms 计算在线的房间数量
func (b *Bucket) CountRooms() (res map[string]int32) {
	b.cLock.RLock()         // 加读锁
	defer b.cLock.RUnlock() // 函数结束时释放读锁
	res = make(map[string]int32)
	for roomID, room := range b.rooms {
		if room.Online > 0 {
			res[roomID] = room.Online
		}
	}
	return
}

// DeleteRoom 删除房间
func (b *Bucket) DeleteRoom(r *Room) {
	b.cLock.Lock()        // 加写锁
	delete(b.rooms, r.ID) // 从房间集合中删除房间
	b.cLock.Unlock()      // 释放写锁
	r.Close()             // 关闭房间
}

// ChangeRoom 更改房间
func (b *Bucket) ChangeRoom(nRID string, ch *Channel) (err error) {
	oroom := ch.Room // 获取原房间
	if nRID == "" {
		if oroom != nil && oroom.Del(ch) {
			b.DeleteRoom(oroom) // 如果原房间存在且删除成功，则删除原房间
		}
		ch.Room = nil // 清空通道的房间
	}
	b.cLock.Lock() // 加写锁
	var (
		nm *Room
		ok bool
	)
	if nm, ok = b.rooms[nRID]; !ok {
		if nm = NewRoom(nRID); nm != nil {
			b.rooms[nRID] = nm // 如果新房间不存在，则创建新房间并添加到房间集合中
		}
	}
	b.cLock.Unlock() // 释放写锁
	if oroom != nil && oroom.Del(ch) {
		b.DeleteRoom(oroom) // 如果原房间存在且删除成功，则删除原房间
	}
	if err = nm.Put(ch); err != nil {
		return // 如果添加通道到新房间失败，则返回错误
	}
	ch.Room = nm // 设置通道的房间为新房间
	return
}

// Put 方法将一个新的通道添加到 Bucket 中，并可能将其添加到一个指定的房间中
func (b *Bucket) Put(rID string, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock() // 获取 Bucket 的写锁，防止其他协程同时修改 Bucket
	if och := b.channels[ch.Key]; och != nil {
		och.Close() // 如果 Bucket 中已存在一个与 ch 具有相同 Key 的通道，就关闭这个通道
	}
	b.channels[ch.Key] = ch // 将新的通道 ch 添加到 Bucket 的通道集合中
	if rID != "" {
		// 如果提供了房间 ID
		if room, ok = b.rooms[rID]; !ok {
			// 如果不存在与 rID 对应的房间，就创建一个新的房间，并将其添加到房间集合中
			room = NewRoom(rID)
			b.rooms[rID] = room
		}
		ch.Room = room // 将 ch 的 Room 字段设置为这个房间
	}
	b.ipCnts[ch.IP]++ // 增加与 ch 的 IP 对应的计数器的值
	b.cLock.Unlock()
	if room != nil {
		// 如果 room 不为 nil，就尝试将 ch 添加到 room 中
		if err = room.Put(ch); err != nil {
			// 如果添加失败，就返回错误
			return
		}
	}
	return // 返回 nil 错误
}

func (b *Bucket) Del(ch *Channel) {
	room := ch.Room
	b.cLock.Lock()
	if dch, ok := b.channels[ch.Key]; ok {
		if dch == ch {
			delete(b.channels, dch.Key)
		}
		// 如果与 ch 的 IP 对应的计数器的值大于 1，就将其减 1
		// 否则，就从 Bucket 的 IP 计数器集合中删除这个 IP
		if b.ipCnts[ch.IP] > 1 {
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.cLock.Unlock()
	if room != nil && room.Del(ch) {
		// 如果 ch 所在的房间不为 nil，
		// 且从房间中删除 ch 成功，就从 Bucket 的房间集合中删除这个房间
		b.DeleteRoom(room)
	}
}

func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.channels[key]
	b.cLock.RUnlock()
	return
}

// Broadcast 广播
func (b *Bucket) Broadcast(p *ProtoRing, op int32) {
	b.cLock.RLock()
	defer b.cLock.RUnlock()
	for _, ch := range b.channels {
		if !ch.NeedPush(op) {
			continue
		}
		if err := ch.Push(p); err != nil {
			fmt.Printf("broadcast room push error(%v)\n", err)
		}
	}
}

func (b *Bucket) DelRoom(room *Room) {
	b.cLock.Lock()
	delete(b.rooms, room.ID)
	b.cLock.Unlock()
	room.Close()
}

// BroadcastRoom 方法向一个指定的房间广播一个消息
func (b *Bucket) BroadcastRoom(arg *pb.BroadcastRoomReq) {
	// 使用 atomic.AddUint64 和取模操作来选择一个协程
	b.routines[atomic.AddUint64(&b.routinesNum, 1)%b.c.RoutineAmount] <- arg
}
func (b *Bucket) Rooms() (res map[string]struct{}) {
	res = make(map[string]struct{})
	b.cLock.RLock()
	defer b.cLock.RUnlock()
	for roomID, room := range b.rooms {
		if room.Online > 0 {
			res[roomID] = struct{}{}
		}
	}
	return
}

// Room 获取房间
func (b *Bucket) Room(roomID string) (r *Room) {
	b.cLock.RLock()     // 加读锁
	r = b.rooms[roomID] // 获取房间
	b.cLock.RUnlock()   // 释放读锁
	return
}

func (b *Bucket) IpCount() (res map[string]struct{}) {
	b.cLock.RLock()
	defer b.cLock.RUnlock()
	res = make(map[string]struct{})
	for ip := range b.ipCnts {
		res[ip] = struct{}{}
	}
	return
}

func (b *Bucket) UpRoomsCount(rcMap map[string]int32) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	for roomID, room := range b.rooms {
		room.AllOnline = rcMap[roomID]
	}
	b.cLock.RUnlock()
}

// roomProc 房间处理协程
func (b *Bucket) roomProc(c chan *pb.BroadcastRoomReq) {
	for {
		if arg, ok := <-c; ok { // 从通道中获取广播房间请求
			if room := b.Room(arg.RoomID); room != nil { // 如果房间存在
				// 创建一个新的 ProtoRing
				room.Push(&ProtoRing{Proto: arg.Proto}) // 推送协议
			}
		}
	}
}
