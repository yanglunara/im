package logic

import (
	"context"

	"github.com/yanglunara/discovery/recovery"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	logicpb "github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type logicServer struct {
	logicpb.UnimplementedLogicServer
	srv model.LogicConnService
}

var (
	_ logicpb.LogicServer = (*logicServer)(nil)
)

func NewLogicService(c *conf.GrpcServer, lgs model.LogicConnService) *bgrpc.Service {
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
	logicpb.RegisterLogicServer(srv, &logicServer{
		srv: lgs,
	})
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
