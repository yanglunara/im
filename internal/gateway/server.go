package gateway

import (
	"context"
	"time"

	"github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	imgrpc "github.com/yanglunara/im/pkg/transport/grpc"
	"github.com/zhenjl/cityhash"
)

type GatewayServer struct {
	c         *conf.Config
	rpcClient logic.LogicClient
	buckets   []*Bucket
	bucketIdx uint32
	serverID  string
}

func NewServer(ctx context.Context, c *conf.Config, serverID string) *GatewayServer {
	g := imgrpc.NewGrpcConnService(
		imgrpc.WithDiscoveryAddress(c.Consul.Address),
		imgrpc.WithEendpoint(c.GrpcClient.Addr),
	)
	client, err := g.Dial(ctx)
	if err != nil {
		panic(err)
	}
	gs := &GatewayServer{
		c:         c,
		rpcClient: logic.NewLogicClient(client),
	}
	gs.buckets = make([]*Bucket, c.Bucket.Size)
	gs.bucketIdx = uint32(c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		gs.buckets[i] = NewBucket(c.Bucket)
	}
	// 启动申报当前链接
	go gs.runProc(ctx)
	return gs
}

func (gs *GatewayServer) Bucket(key string) *Bucket {
	return gs.buckets[cityhash.CityHash32([]byte(key), uint32(len(key)))%gs.bucketIdx]
}

func (g *GatewayServer) runProc(ctx context.Context) {
	timeTicker := time.NewTicker(time.Second * 10)
	defer timeTicker.Stop()
	for {
		roomCount := make(map[string]int32)
		for _, bucket := range g.buckets {
			for roomID, count := range bucket.CountRooms() {
				roomCount[roomID] += count
			}
		}
		res, err := g.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
			Server:    g.serverID,
			RoomCount: roomCount,
		})
		if err != nil || (res != nil && len(res.AllRoomCount) > 0) {
			time.Sleep(1 * time.Second)
			continue
		}
		for _, bucket := range g.buckets {
			bucket.UpRoomsCount(res.AllRoomCount)
		}
	}
}
