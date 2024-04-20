package gateway

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	pt "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/yanglunara/im/api/logic"
	"github.com/yanglunara/im/api/protocol"
	pb "github.com/yanglunara/im/api/protocol"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
)

type WebSocketInter interface {
	Upgrade(w http.ResponseWriter, r *http.Request) (*Channel, string, error)
	DispatchWebsocket()
}

var (
	WebSocket WebSocketInter
	mutex     sync.Mutex
)

type imSocket struct {
	keyGUID  []byte
	conf     *conf.Config
	serverID string
	gs       *GatewayServer
	ch       *Channel
}

func NewWebSocket(conf *conf.Config, gs *GatewayServer, serverID string) WebSocketInter {
	if WebSocket == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if WebSocket == nil {
			WebSocket = &imSocket{
				keyGUID:  []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B22"),
				conf:     conf,
				serverID: serverID,
				gs:       gs,
				ch: NewChannel(
					conf.Protocol.CliProto,
					conf.Protocol.SvrProto,
				),
			}
		}
	}
	return WebSocket
}

func (i *imSocket) acceptKey(challengeKey string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(challengeKey))
	_, _ = h.Write(i.keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type AuthWebSocket struct {
	MID      int64
	Key, RID string
	Accepts  []int32
	Hb       time.Duration
}

func (i *imSocket) Upgrade(w http.ResponseWriter, r *http.Request) (ch *Channel, k string, err error) {
	// 设置启动时间
	i.ch.SetLastHB()
	ws := new(websocket.Upgrader)
	ws.CheckOrigin = func(req *http.Request) bool {
		check := true
		switch {
		case req.Method != "GET":
			check = false
		case req.Header.Get("Sec-Websocket-Version") != "13":
			check = false
		case !strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade"):
			check = false
		}
		if challengeKey := req.Header.Get("Sec-Websocket-Key"); challengeKey == "" {
			check = false
		} else {
			req.Header.Set("Sec-Websocket-Key", i.acceptKey(challengeKey))
		}
		return check
	}
	conn := new(websocket.Conn)
	if conn, err = ws.Upgrade(w, r, nil); err != nil {
		return nil, "", err
	}
	var (
		p *ProtoRing
	)
	var (
		a *AuthWebSocket
	)
	//设置IP
	i.ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	// 注册链接
	if p, err = i.ch.CliProto.Put(); err == nil {
		if a, err = i.authWebSocket(r.Context(), conn, p, r.Header.Get("Cookie")); err == nil {
			i.ch.Watch(a.Accepts...)
			i.ch.Mid, i.ch.Key = a.MID, a.Key
			if err = i.gs.Bucket(i.ch.Key).Put(a.RID, i.ch); err == nil {
				// 成功后绑定连接
				p.Conn = conn
			}
		}
	}
	if err != nil {
		_ = conn.Close()
		return nil, "", err
	}
	conn = nil
	// 设置当前的 protoRing
	i.ch.SetProtoRing(p)
	return i.ch, r.Header.Get("Cookie"), err
}

func (i *imSocket) authWebSocket(ctx context.Context, ws *websocket.Conn, p *ProtoRing, cookie string) (a *AuthWebSocket, err error) {
	var (
		buf []byte
		op  int
	)
	if op, buf, err = ws.ReadMessage(); err != nil {
		return
	}

	switch op {
	case websocket.BinaryMessage, websocket.TextMessage:
		if err = pt.Unmarshal(buf, p.Proto); err != nil {
			return
		}
		// 鉴权失败
		var (
			res *logic.ConnectResp
		)
		if protocol.Operation_Type(p.Op) == protocol.Operation_Auth {
			res, err = LogicGrpcClient.Connect(ctx, &logic.ConnectReq{
				Server: i.serverID,
				Cookie: cookie,
				Token:  p.Body,
			})
			if err != nil {
				return
			}
			a = &AuthWebSocket{
				MID:     res.Mid,
				Key:     res.Key,
				RID:     res.RoomID,
				Accepts: res.Accepts,
				Hb:      time.Duration(res.Heartbeat),
			}
		}
		p.Op = int32(pb.Operation_AuthResp)
		p.Body = nil
		if buf, err = pt.Marshal(p.Proto); err != nil {
			return
		}
		if err = ws.WriteMessage(websocket.BinaryMessage, buf); err != nil {
			return
		}
	case websocket.PongMessage, websocket.CloseMessage:
		return nil, errors.New("messageType not supported")
	}
	return
}

func (i *imSocket) DispatchWebsocket() {
	var (
		err    error
		online int32
		pp     *ProtoRing
	)
	for {
		var (
			p = i.ch.Ready()
		)
		pp = p
		switch protocol.Operation_Type(p.Op) {
		case protocol.Operation_ProtoFinish:
			goto failed
		case protocol.Operation_ProtoReady:
			for {
				var (
					buf []byte
				)
				if p, err = i.ch.CliProto.Get(); err != nil {
					goto failed
				}
				if p.Op == int32(pb.Operation_HeartbeatResp) {
					if i.ch.Room != nil {
						online = i.ch.Room.OnlineNum()
					}
					p.Op = int32(pb.Operation_AuthResp)
					onlineStr := map[string]int32{
						"online": online,
					}
					if p.Body, err = json.Marshal(onlineStr); err != nil {
						goto failed
					}
					if buf, err = pt.Marshal(p); err != nil {
						goto failed
					}
					if err = p.Conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
						goto failed
					}
				} else {
					if err = p.Conn.WriteMessage(websocket.BinaryMessage, p.Body); err != nil {
						goto failed
					}
				}
				p.Body = nil
				i.ch.CliProto.IncrementReadPointer()
			}
		default:
			// 默认转发
			if err = p.Conn.WriteMessage(websocket.BinaryMessage, p.Body); err != nil {
				goto failed
			}
		}
	}
failed:
	// 关闭链接
	if pp.Conn != nil {
		_ = pp.Conn.Close()
	}
	logger.Logger.Error("DispatchWebsocket failed err %v p:%v", zap.Error(err), zap.Any("point", pp))
	pp = nil
}
