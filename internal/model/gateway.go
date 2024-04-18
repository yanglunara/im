package model

import (
	"context"

	"github.com/yanglunara/im/api/gateway"
)

type GatewayPushService interface {
	PushMessage(ctx context.Context, req *gateway.PushMessageReq) (*gateway.PushMessageResp, error)
	Broadcast(ctx context.Context, req *gateway.BroadcastReq) (*gateway.BroadcastResp, error)
	BroadcastRoom(ctx context.Context, req *gateway.BroadcastRoomReq) (*gateway.BroadcastRoomResp, error)
	Rooms(ctx context.Context, req *gateway.RoomsReq) (*gateway.RoomsResp, error)
}
