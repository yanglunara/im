// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.18.0
// source: api/logic/logic.proto

package logic

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Logic_Connect_FullMethodName     = "/im.logic.Logic/Connect"
	Logic_Disconnect_FullMethodName  = "/im.logic.Logic/Disconnect"
	Logic_Heartbeat_FullMethodName   = "/im.logic.Logic/Heartbeat"
	Logic_RenewOnline_FullMethodName = "/im.logic.Logic/RenewOnline"
	Logic_Receive_FullMethodName     = "/im.logic.Logic/Receive"
	Logic_Nodes_FullMethodName       = "/im.logic.Logic/Nodes"
)

// LogicClient is the client API for Logic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogicClient interface {
	// Connect
	Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectResp, error)
	// Disconnect
	Disconnect(ctx context.Context, in *DisconnectReq, opts ...grpc.CallOption) (*DisconnectResp, error)
	// Heartbeat
	Heartbeat(ctx context.Context, in *HeartbeatReq, opts ...grpc.CallOption) (*HeartbeatResp, error)
	// RenewOnline
	RenewOnline(ctx context.Context, in *OnlineReq, opts ...grpc.CallOption) (*OnlineResp, error)
	// Receive
	Receive(ctx context.Context, in *ReceiveReq, opts ...grpc.CallOption) (*ReceiveResp, error)
	// ServerList
	Nodes(ctx context.Context, in *NodesReq, opts ...grpc.CallOption) (*NodesResp, error)
}

type logicClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicClient(cc grpc.ClientConnInterface) LogicClient {
	return &logicClient{cc}
}

func (c *logicClient) Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectResp, error) {
	out := new(ConnectResp)
	err := c.cc.Invoke(ctx, Logic_Connect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Disconnect(ctx context.Context, in *DisconnectReq, opts ...grpc.CallOption) (*DisconnectResp, error) {
	out := new(DisconnectResp)
	err := c.cc.Invoke(ctx, Logic_Disconnect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Heartbeat(ctx context.Context, in *HeartbeatReq, opts ...grpc.CallOption) (*HeartbeatResp, error) {
	out := new(HeartbeatResp)
	err := c.cc.Invoke(ctx, Logic_Heartbeat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) RenewOnline(ctx context.Context, in *OnlineReq, opts ...grpc.CallOption) (*OnlineResp, error) {
	out := new(OnlineResp)
	err := c.cc.Invoke(ctx, Logic_RenewOnline_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Receive(ctx context.Context, in *ReceiveReq, opts ...grpc.CallOption) (*ReceiveResp, error) {
	out := new(ReceiveResp)
	err := c.cc.Invoke(ctx, Logic_Receive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Nodes(ctx context.Context, in *NodesReq, opts ...grpc.CallOption) (*NodesResp, error) {
	out := new(NodesResp)
	err := c.cc.Invoke(ctx, Logic_Nodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicServer is the server API for Logic service.
// All implementations must embed UnimplementedLogicServer
// for forward compatibility
type LogicServer interface {
	// Connect
	Connect(context.Context, *ConnectReq) (*ConnectResp, error)
	// Disconnect
	Disconnect(context.Context, *DisconnectReq) (*DisconnectResp, error)
	// Heartbeat
	Heartbeat(context.Context, *HeartbeatReq) (*HeartbeatResp, error)
	// RenewOnline
	RenewOnline(context.Context, *OnlineReq) (*OnlineResp, error)
	// Receive
	Receive(context.Context, *ReceiveReq) (*ReceiveResp, error)
	// ServerList
	Nodes(context.Context, *NodesReq) (*NodesResp, error)
	mustEmbedUnimplementedLogicServer()
}

// UnimplementedLogicServer must be embedded to have forward compatible implementations.
type UnimplementedLogicServer struct {
}

func (UnimplementedLogicServer) Connect(context.Context, *ConnectReq) (*ConnectResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedLogicServer) Disconnect(context.Context, *DisconnectReq) (*DisconnectResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (UnimplementedLogicServer) Heartbeat(context.Context, *HeartbeatReq) (*HeartbeatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedLogicServer) RenewOnline(context.Context, *OnlineReq) (*OnlineResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewOnline not implemented")
}
func (UnimplementedLogicServer) Receive(context.Context, *ReceiveReq) (*ReceiveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Receive not implemented")
}
func (UnimplementedLogicServer) Nodes(context.Context, *NodesReq) (*NodesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nodes not implemented")
}
func (UnimplementedLogicServer) mustEmbedUnimplementedLogicServer() {}

// UnsafeLogicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogicServer will
// result in compilation errors.
type UnsafeLogicServer interface {
	mustEmbedUnimplementedLogicServer()
}

func RegisterLogicServer(s grpc.ServiceRegistrar, srv LogicServer) {
	s.RegisterService(&Logic_ServiceDesc, srv)
}

func _Logic_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Connect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Connect(ctx, req.(*ConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisconnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Disconnect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Disconnect(ctx, req.(*DisconnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Heartbeat(ctx, req.(*HeartbeatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_RenewOnline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).RenewOnline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_RenewOnline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).RenewOnline(ctx, req.(*OnlineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Receive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Receive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Receive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Receive(ctx, req.(*ReceiveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Nodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Nodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Nodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Nodes(ctx, req.(*NodesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Logic_ServiceDesc is the grpc.ServiceDesc for Logic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Logic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "im.logic.Logic",
	HandlerType: (*LogicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _Logic_Connect_Handler,
		},
		{
			MethodName: "Disconnect",
			Handler:    _Logic_Disconnect_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Logic_Heartbeat_Handler,
		},
		{
			MethodName: "RenewOnline",
			Handler:    _Logic_RenewOnline_Handler,
		},
		{
			MethodName: "Receive",
			Handler:    _Logic_Receive_Handler,
		},
		{
			MethodName: "Nodes",
			Handler:    _Logic_Nodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/logic/logic.proto",
}
