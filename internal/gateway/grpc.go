package gateway

import (
	"context"
	"net"
	"sync"
	"time"

	gatewaypb "github.com/yanglunara/im/api/gateway"
	"github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yanglunara/im/internal/model"
	imgrpc "github.com/yanglunara/im/pkg/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

func NewPushServer(c *conf.GRPCServer, gps model.GatewayPushService) *grpc.Server {
	keepalive := grpc.KeepaliveParams(keepalive.ServerParameters{
		// MaxConnectionIdle是一个持续时间，如果一个客户端连接在这段时间内没有任何活动，gRPC服务器将关闭该连接。
		MaxConnectionIdle: time.Duration(c.IdleTimeout) * time.Second,
		// MaxConnectionAgeGrace是一个持续时间，在一个连接达到MaxConnectionAge后，gRPC服务器将给这个连接这么长的时间来完成正在进行的RPC调用。
		MaxConnectionAgeGrace: time.Duration(c.ForceCloseWait) * time.Second,
		// Time是gRPC服务器在两次keepalive探测之间等待的时间。
		Time: time.Duration(c.KeepaliveInterval) * time.Second,
		// Timeout是等待keepalive探测响应的时间。如果在这段时间内没有收到响应，gRPC服务器将关闭连接。
		Timeout: time.Duration(c.KeepaliveTimeout) * time.Second,
		// MaxConnectionAge是一个持续时间，如果一个连接的寿命超过这个时间，gRPC服务器将关闭该连接。
		MaxConnectionAge: time.Duration(c.MaxLiftTime) * time.Second,
	})
	srv := grpc.NewServer(keepalive)
	// 注册服务
	gatewaypb.RegisterGatewayServer(srv, &gatewayPushServer{
		gps: gps,
	})
	var (
		lis net.Listener
		err error
	)
	if lis, err = net.Listen(c.Network, c.Address); err != nil {
		panic(err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
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
