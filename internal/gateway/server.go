package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	ttime "github.com/yanglunara/im/pkg/time"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
	"github.com/yanglunara/im/api/logic"
	"github.com/yanglunara/im/api/protocol"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	m "github.com/yanglunara/im/internal/model"
	imgrpc "github.com/yanglunara/im/pkg/transport/grpc"
	"github.com/zhenjl/cityhash"
)

type RegisterChan struct {
	*websocket.Conn
	*http.Request
}
type WriteChan struct {
	*protocol.Proto
	*websocket.Conn
}

var (
	WebSocket *Websocket
	mu        sync.Mutex
)

type Websocket struct {
	c            *conf.Config
	rpcClient    logic.LogicClient
	buckets      []*Bucket
	bucketIdx    uint32
	serverID     string
	RegisterChan chan *RegisterChan
	WriteChan    chan *WriteChan
	round        *Round
	Count        int
	MaxInt       int
}

func NewServer(ctx context.Context, c *conf.Config, serverID string) *Websocket {
	if WebSocket == nil {
		mu.Lock()
		defer mu.Unlock()
		if WebSocket == nil {
			g := imgrpc.NewGrpcConnService(
				imgrpc.WithDiscoveryAddress(c.Consul.Address),
				imgrpc.WithEendpoint(c.GrpcClient.Addr),
			)
			client, err := g.Dial(ctx)
			if err != nil {
				panic(err)
			}
			gs := &Websocket{
				c:            c,
				rpcClient:    logic.NewLogicClient(client),
				serverID:     serverID,
				round:        NewRound(c),
				RegisterChan: make(chan *RegisterChan, 1024),
				WriteChan:    make(chan *WriteChan, 1024),
				MaxInt:       1<<31 - 1,
				Count:        0,
			}
			gs.buckets = make([]*Bucket, c.Bucket.Size)
			gs.bucketIdx = uint32(c.Bucket.Size)
			for i := 0; i < c.Bucket.Size; i++ {
				gs.buckets[i] = NewBucket(&c.Bucket)
			}
			WebSocket = gs
			// 启动申报当前链接
			go gs.runProc(ctx)
		}
	}
	return WebSocket
}

func (ws *Websocket) Bucket(key string) *Bucket {
	return ws.buckets[cityhash.CityHash32([]byte(key), uint32(len(key)))%ws.bucketIdx]
}

func (ws *Websocket) runProc(ctx context.Context) {
	timeTicker := time.NewTicker(time.Second * 10)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			roomCount := make(map[string]int32)
			for _, bucket := range ws.buckets {
				for roomID, count := range bucket.CountRooms() {
					roomCount[roomID] += count
				}
			}
			res, err := ws.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
				Server:    ws.serverID,
				RoomCount: roomCount,
			})
			if err != nil || res.AllRoomCount == nil {
				continue
			}
			for _, bucket := range ws.buckets {
				bucket.UpRoomsCount(res.AllRoomCount)
			}
		}
	}
}

func (w *Websocket) Buckets() []*Bucket {
	return w.buckets
}

func (w *Websocket) RandServerHeartbeat() time.Duration {
	return m.MinServerHeartbeat + time.Duration(rand.Int63n(int64(m.MaxServerHeartbeat-m.MinServerHeartbeat)))
}

func (w *Websocket) ServerWebsocket(conn *websocket.Conn, req *http.Request, tr *ttime.Timer) {
	var (
		err     error
		rid     string
		accepts []int32
		hb      time.Duration
		p       *protocol.Proto
		b       *Bucket
		trd     *ttime.TimerData
		//lastHB  = time.Now()
		ch = NewChannel(w.c.Protocol.CliProto, w.c.Protocol.SvrProto)
	)
	// reader
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// handshake
	step := 0
	trd = tr.Add(w.c.Protocol.HandShakeTimeout, func() {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		_ = conn.Close()
	})
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	//step = 1
	if p, err = ch.CliProto.Put(); err == nil {
		if ch.Mid, ch.Key, rid, accepts, hb, err = w.authWebsocket(ctx, conn, p, req.Header.Get("Cookie")); err == nil {
			ch.Watch(accepts...)
			b = w.Bucket(ch.Key)
			err = b.Put(rid, ch)
		}
	}
	step = 2
	if err != nil {
		_ = conn.Close()
		tr.Del(trd)
		return
	}
	// trd.Key = ch.Key
	// tr.Set(trd, time.Duration(1*time.Second))
	// step = 3
	// go s.dispatchWebsocket(conn, ch)
	//serverHeartbeat := s.RandServerHearbeat()

	// for {
	// 	if p, err = ch.CliProto.Put(); err != nil {
	// 		break
	// 	}
	// 	if p.Op == protocol.Op_Heartbeat {
	// 		tr.Set(trd, hb)
	// 		p.Op = protocol.Op_HeartbeatResp
	// 		p.Body = nil
	// 		if now := time.Now(); now.Sub(lastHB) > serverHeartbeat {
	// 			if err1 := s.Hearbeat(ctx, ch.Mid, ch.Key); err1 == nil {
	// 				lastHB = now
	// 			}
	// 		}
	// 		step++
	// 	} else {
	// 		if err = s.Broadcast(ctx, p, ch, b); err != nil {
	// 			break
	// 		}
	// 	}
	// 	ch.CliProto.IncrementReadPointer()
	// 	ch.Signal()
	// }
	logger.Logger.Info("websocket message", zap.Int("step", step), zap.Duration("hb", hb), zap.String("rid", rid), zap.Any("accepts", accepts))
	// b.Del(ch)
	// tr.Del(trd)
	// _ = conn.Close()
	// ch.Close()
	// if _, err = s.rpcClient.Disconnect(ctx, &logic.DisconnectReq{
	// 	Mid: ch.Mid,
	// 	Key: ch.Key,
	// }); err != nil {
	// 	logger.Logger.Error("disconnect error", zap.Error(err))
	// }
	w.WriteChan <- &WriteChan{
		Proto: p,
		Conn:  conn,
	}
}

func SplitInt32s(s, p string) ([]int32, error) {
	if s == "" {
		return nil, nil
	}
	ss := strings.Split(s, p)
	res := make([]int32, 0, len(ss))
	for _, sc := range ss {
		i, err := strconv.ParseInt(sc, 10, 32)
		if err != nil {
			return nil, err
		}
		res = append(res, int32(i))
	}
	return res, nil
}

func (s *Websocket) Broadcast(ctx context.Context, p *protocol.Proto, ch *Channel, b *Bucket) error {
	switch p.Op {
	case protocol.Op_ChangeRoom:
		if err := b.ChangeRoom(string(p.Body), ch); err != nil {
			logger.Logger.Error("change room error", zap.Error(err))
		}
		p.Op = protocol.Op_ChangeRoomResp
	case protocol.Op_Sub:
		if ops, err := SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = protocol.Op_SubReply
	case protocol.Op_Unsub:
		if ops, err := SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = protocol.Op_UnsubResp
	default:
		if _, err := s.rpcClient.Receive(ctx, &logic.ReceiveReq{
			Mid:   ch.Mid,
			Proto: p,
		}); err != nil {
			logger.Logger.Error("receive error", zap.Error(err))
		}
		p.Body = nil
	}
	return nil
}

func (s *Websocket) Hearbeat(ctx context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatReq{
		Mid: mid,
		Key: key,
	})
	return

}

func (s *Websocket) authWebsocket(ctx context.Context, ws *websocket.Conn, p *protocol.Proto, cookie string) (
	mid int64, key, rid string, accepts []int32, hb time.Duration, err error) {
	var (
		buf []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if err = proto.Unmarshal(buf, p); err != nil {
		return
	}
	var (
		res *logic.ConnectResp
	)
	if res, err = s.rpcClient.Connect(ctx, &logic.ConnectReq{
		Server: s.serverID,
		Cookie: cookie,
		Token:  p.Body,
	}); err != nil {
		return
	}
	mid, key, rid, accepts, hb = res.Mid, res.Key, res.RoomID, res.Accepts, time.Duration(res.Heartbeat)
	p.Op = protocol.Op_AuthResp
	p.Body = nil
	//通知前端
	s.WriteChan <- &WriteChan{
		Proto: p,
		Conn:  ws,
	}
	logger.Logger.Info("auth success", zap.Int64("mid", mid), zap.String("key", key), zap.String("rid", rid), zap.Duration("hb", hb))
	return
}

func (wws *Websocket) dispatchWebsocket(ws *websocket.Conn, ch *Channel) {
	var (
		err    error
		finish bool
		online int32
	)
	for {
		var p = ch.Ready()
		switch p {
		case ProtoFinish:
			finish = true
			goto failed
		case ProtoRead:
			if p, err = ch.CliProto.Get(); err != nil {
				break
			}
			if p.Op == protocol.Op_HeartbeatResp {
				if ch.Room != nil {
					online = ch.Room.OnlineNum()
				}
				p.Body = []byte(fmt.Sprintf(`{"online":%d}`, online))
			}
			ch.CliProto.IncrementReadPointer()
			wws.WriteChan <- &WriteChan{
				Proto: p,
				Conn:  ws,
			}
			finish = true
			goto failed
		}
	}
failed:
	ws.Close()
	for !finish {
		finish = (ch.Ready() == ProtoFinish)
	}
}

func (wws *Websocket) Upgrade(w http.ResponseWriter, r *http.Request) error {
	var (
		err error
	)
	ws := new(websocket.Upgrader)
	ws.CheckOrigin = func(req *http.Request) bool {
		return true
	}
	conn := new(websocket.Conn)
	if conn, err = ws.Upgrade(w, r, nil); err != nil {
		return err
	}
	wws.RegisterChan <- &RegisterChan{
		Conn:    conn,
		Request: r,
	}
	return nil
}

func (w *Websocket) Start() {
	for {
		select {
		case conn := <-w.RegisterChan:
			// 无序处理
			w.ServerWebsocket(
				conn.Conn,
				conn.Request,
				w.round.Timer(w.Count),
			)
			if w.Count++; w.Count == w.MaxInt {
				w.Count = 0
			}
		case p := <-w.WriteChan:
			bytes, _ := json.Marshal(p.Proto)
			if err := p.Conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				logger.Logger.Error("write message error", zap.Any("proto", p.Proto))
			}
		}
	}
}

func (w *Websocket) Close() error {
	close(w.RegisterChan)
	close(w.WriteChan)
	return nil
}
