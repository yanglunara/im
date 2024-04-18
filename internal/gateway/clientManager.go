package gateway

import (
	"context"
	"time"

	pt "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/yanglunara/im/api/logic"
	pb "github.com/yanglunara/im/api/protocol"
)

type ClientManager struct {
	ch       *Channel
	serverID string
}

type AuthWebSocket struct {
	MID      int64
	Key, RID string
	Accepts  []int32
	Hb       time.Duration
}

func NewClientManager(serverID string) *ClientManager {
	return &ClientManager{
		serverID: serverID,
	}
}

func (c *ClientManager) Start(ctx context.Context, p *pb.Proto, cookie string, ws *websocket.Conn) (a *AuthWebSocket, err error) {
	if a, err = c.authWebSocket(ctx, ws, p, cookie); err != nil {
		return nil, err
	}
	return
}

func (c *ClientManager) authWebSocket(ctx context.Context, ws *websocket.Conn, p *pb.Proto, cookie string) (a *AuthWebSocket, err error) {
	var (
		buf []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) > 0 {
		p.Body = buf
	}
	res, err := GRPCClient.Connect(ctx, &logic.ConnectReq{
		Server: c.serverID,
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
	return
}

func (c *ClientManager) WriteWebsocket(ws *websocket.Conn, p *pb.Proto) (err error) {
	var (
		buf []byte
	)
	p.Op = int32(pb.Operation_AuthResp)
	p.Body = nil
	if buf, err = pt.Marshal(p); err != nil {
		return
	}
	if err = ws.WriteMessage(websocket.TextMessage, buf); err != nil {
		return
	}
	return nil
}
