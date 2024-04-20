package gateway

import (
	"context"

	gatewaypb "github.com/yanglunara/im/api/gateway"
	"github.com/yanglunara/im/internal/model"
)

var (
	_ model.GatewayPushService = (*pushService)(nil)
)

type pushService struct {
}

func NewGatewayService() model.GatewayPushService {
	return &pushService{}
}

func (p *pushService) Broadcast(ctx context.Context, req *gatewaypb.BroadcastReq) (*gatewaypb.BroadcastResp, error) {
	return nil, nil
}

func (p *pushService) PushMessage(ctx context.Context, req *gatewaypb.PushMessageReq) (*gatewaypb.PushMessageResp, error) {
	return nil, nil
}

func (p *pushService) BroadcastRoom(ctx context.Context, req *gatewaypb.BroadcastRoomReq) (*gatewaypb.BroadcastRoomResp, error) {
	return nil, nil
}

func (p *pushService) Rooms(ctx context.Context, req *gatewaypb.RoomsReq) (*gatewaypb.RoomsResp, error) {
	return nil, nil
}
