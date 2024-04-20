package gateway

import (
	"context"
	"sync"

	"github.com/yanglunara/discovery/recovery"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	gatewaypb "github.com/yanglunara/im/api/gateway"
	"github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yanglunara/im/internal/model"
	imgrpc "github.com/yanglunara/im/pkg/transport/grpc"
)

var (
	LogicGrpcClient logic.LogicClient
	once            sync.Once
)

var (
	_ gatewaypb.GatewayServer = (*gatewayPushServer)(nil)
)

type gatewayPushServer struct {
	gatewaypb.UnimplementedGatewayServer
	gps model.GatewayPushService
}

func NewLogicGrpc(ctx context.Context, c *conf.Config) logic.LogicClient {
	once.Do(func() {
		g := imgrpc.NewGrpcConnService(
			imgrpc.WithDiscoveryAddress(c.Consul.Address),
			imgrpc.WithEendpoint(c.GrpcClient.Addr),
		)
		client, err := g.Dial(ctx)
		if err != nil {
			panic(err)
		}
		LogicGrpcClient = logic.NewLogicClient(client)
	})
	return LogicGrpcClient
}

func NewGatewayGrpcServer(c *conf.GrpcServer, gps model.GatewayPushService) *bgrpc.Service {
	var opts = []bgrpc.ServiceOption{
		bgrpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Network != "" {
		opts = append(opts, bgrpc.Network(c.Network))
	}
	if c.Addr != "" {
		opts = append(opts, bgrpc.Address(c.Addr))
	}
	if c.Timeout != 0 {
		opts = append(opts, bgrpc.Timeout(c.Timeout))
	}
	opts = append(opts, bgrpc.OpenHealth())
	srv := bgrpc.NewGrpcServer(opts...)
	// 注册服务
	gatewaypb.RegisterGatewayServer(srv, &gatewayPushServer{
		gps: gps,
	})
	return srv

}

func (g *gatewayPushServer) Broadcast(ctx context.Context, req *gatewaypb.BroadcastReq) (*gatewaypb.BroadcastResp, error) {
	return g.gps.Broadcast(ctx, req)
}

func (g *gatewayPushServer) PushMessage(ctx context.Context, req *gatewaypb.PushMessageReq) (*gatewaypb.PushMessageResp, error) {
	return g.gps.PushMessage(ctx, req)
}

func (g *gatewayPushServer) BroadcastRoom(ctx context.Context, req *gatewaypb.BroadcastRoomReq) (*gatewaypb.BroadcastRoomResp, error) {
	return g.gps.BroadcastRoom(ctx, req)
}

func (g *gatewayPushServer) Rooms(ctx context.Context, req *gatewaypb.RoomsReq) (*gatewaypb.RoomsResp, error) {
	return g.gps.Rooms(ctx, req)
}
