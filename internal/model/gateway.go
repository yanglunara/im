package model

import (
	"context"
	"time"

	"github.com/yanglunara/im/api/gateway"
)

type GatewayPushService interface {
	PushMessage(ctx context.Context, req *gateway.PushMessageReq) (*gateway.PushMessageResp, error)
	Broadcast(ctx context.Context, req *gateway.BroadcastReq) (*gateway.BroadcastResp, error)
	BroadcastRoom(ctx context.Context, req *gateway.BroadcastRoomReq) (*gateway.BroadcastRoomResp, error)
	Rooms(ctx context.Context, req *gateway.RoomsReq) (*gateway.RoomsResp, error)
}

var (
	MinServerHeartbeat = time.Minute * 10
	MaxServerHeartbeat = time.Minute * 30
)
