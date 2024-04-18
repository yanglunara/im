package logic

import (
	"context"
	"net"
	"time"

	logicpb "github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type logicServer struct {
	logicpb.UnimplementedLogicServer
	srv model.LogicConnService
}

var (
	_ logicpb.LogicServer = (*logicServer)(nil)
)

func NewService(c *conf.GRPCServer, lgs model.LogicConnService) *grpc.Server {
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
	logicpb.RegisterLogicServer(srv, &logicServer{
		srv: lgs,
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

func (l *logicServer) Connect(ctx context.Context, req *logicpb.ConnectReq) (*logicpb.ConnectResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.Connect(ctx, req)
}

func (l *logicServer) Disconnect(ctx context.Context, req *logicpb.DisconnectReq) (*logicpb.DisconnectResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.Disconnect(ctx, req)
}

func (l *logicServer) Heartbeat(ctx context.Context, req *logicpb.HeartbeatReq) (*logicpb.HeartbeatResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.Heartbeat(ctx, req)
}

func (l *logicServer) RenewOnline(ctx context.Context, req *logicpb.OnlineReq) (*logicpb.OnlineResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.RenewOnline(ctx, req)
}

func (l *logicServer) Receive(ctx context.Context, req *logicpb.ReceiveReq) (*logicpb.ReceiveResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.Receive(ctx, req)
}

func (l *logicServer) Nodes(ctx context.Context, req *logicpb.NodesReq) (*logicpb.NodesResp, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return l.srv.Nodes(ctx, req)
}
