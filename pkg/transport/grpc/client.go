package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/yanglunara/discovery/builder"
	"github.com/yanglunara/discovery/register"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type (
	GrpcConnService interface {
		Dial(ctx context.Context) (*grpc.ClientConn, error)
	}
)
type ClientOption func(o *grpcConnService)

func WithDiscoveryAddress(addres string) ClientOption {
	return func(o *grpcConnService) {
		o.address = addres
	}
}
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *grpcConnService) {
		o.grpcOpts = opts
	}
}
func WithEendpoint(endpoint string) ClientOption {
	return func(o *grpcConnService) {
		o.endpoint = fmt.Sprintf("discovery:///%s", endpoint)
	}
}

var (
	_ GrpcConnService = (*grpcConnService)(nil)
)

type grpcConnService struct {
	timeout                time.Duration
	balancerName           string
	subsetSize             int
	printDiscoveryDebugLog bool
	healthCheckConfig      string
	discovery              register.Discovery // 服务发现
	address                string
	grpcOpts               []grpc.DialOption
	endpoint               string
	WindowSize             int32
	aliveTime              time.Duration
}

func NewGrpcConnService(opt ...ClientOption) GrpcConnService {
	gcs := grpcConnService{
		timeout:                3 * time.Second,
		balancerName:           "round_robin",
		subsetSize:             25,
		printDiscoveryDebugLog: true,
		healthCheckConfig:      `,"healthCheckConfig":{"serviceName":""}`,
		WindowSize:             1 << 24,
		aliveTime:              10 * time.Second,
	}
	for _, o := range opt {
		o(&gcs)
	}
	gcs.discovery = builder.NewConsulDiscovery(gcs.address)
	return &gcs
}

func (g *grpcConnService) Dial(ctx context.Context) (*grpc.ClientConn, error) {
	return g.dial(ctx, true)
}

func (g *grpcConnService) dial(ctx context.Context, insecure bool) (*grpc.ClientConn, error) {
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]%s}`,
			g.balancerName, g.healthCheckConfig)),
		grpc.WithInitialWindowSize(g.WindowSize),
		grpc.WithInitialConnWindowSize(g.WindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(g.WindowSize))),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(g.WindowSize))),
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                g.aliveTime,
		//	Timeout:             g.timeout,
		//	PermitWithoutStream: true,
		//}),
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if g.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			builder.NewBuilder(g.discovery),
		))
	}
	if len(g.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, g.grpcOpts...)
	}
	return grpc.DialContext(ctx, g.endpoint, grpcOpts...)
}
