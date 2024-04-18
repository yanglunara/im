package gateway

import (
	"context"
	"sync"

	"github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	imgrpc "github.com/yanglunara/im/pkg/transport/grpc"
)

var (
	GRPCClient logic.LogicClient
	once       sync.Once
)

func NewGRPCClient(ctx context.Context, c *conf.Config) logic.LogicClient {
	once.Do(func() {
		g := imgrpc.NewGrpcConnService(
			imgrpc.WithDiscoveryAddress(c.Consul.Address),
			imgrpc.WithEendpoint(c.GrpcClient.Addr),
		)
		client, err := g.Dial(ctx)
		if err != nil {
			panic(err)
		}
		GRPCClient = logic.NewLogicClient(client)
	})
	return GRPCClient
}
