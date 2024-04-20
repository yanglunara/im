package model

import (
	"context"

	logic "github.com/yanglunara/im/api/logic"
)

type LogicConnService interface {
	Connect(ctx context.Context, req *logic.ConnectReq) (*logic.ConnectResp, error)
	Disconnect(ctx context.Context, req *logic.DisconnectReq) (*logic.DisconnectResp, error)
	Heartbeat(ctx context.Context, req *logic.HeartbeatReq) (*logic.HeartbeatResp, error)
	RenewOnline(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error)
	Receive(ctx context.Context, req *logic.ReceiveReq) (*logic.ReceiveResp, error)
	Nodes(ctx context.Context, req *logic.NodesReq) (*logic.NodesResp, error)
	Ping(ctx context.Context) error
	SetReplicant(ctx context.Context, serverName string)
}
