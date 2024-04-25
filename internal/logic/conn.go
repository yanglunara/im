package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/yanglunara/im/api/logic"
	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/model"
	"github.com/yanglunara/im/internal/plugin"
)

var (
	_ model.LogicConnService = (*connLogic)(nil)
)

type connLogic struct {
	c         *conf.Config
	server    model.ServerInter
	roomCount map[string]int32
	replicant model.Replicant
}

func NewLogic(c *conf.Config) model.LogicConnService {
	server := plugin.InitService(plugin.WithKafka(&conf.Kafka{
		Topic:   c.Kafka.Topic,
		Group:   c.Kafka.Group,
		Brokers: c.Kafka.Brokers,
	}),
		plugin.WithRedis(
			&conf.Redis{
				Network:      c.Redis.Network,
				Addr:         c.Redis.Addr,
				Auth:         c.Redis.Auth,
				Active:       c.Redis.Active,
				Idle:         c.Redis.Idle,
				DialTimeout:  c.Redis.DialTimeout,
				ReadTimeout:  c.Redis.ReadTimeout,
				WriteTimeout: c.Redis.WriteTimeout,
				IdleTimeout:  c.Redis.IdleTimeout,
				Expire:       c.Redis.Expire,
			},
		),
	)
	cl := &connLogic{
		c:         c,
		server:    server,
		roomCount: make(map[string]int32),
	}
	return cl
}
func (l *connLogic) SetReplicant(_ context.Context, serverName string) {
	l.replicant = newReplicant(l.server, conf.Conf.Consul.Address, serverName)
}

func (l *connLogic) Connect(ctx context.Context, req *logic.ConnectReq) (resp *logic.ConnectResp, err error) {
	var params struct {
		Mid      int64   `json:"mid"`
		Key      string  `json:"key"`
		RoomID   string  `json:"room_id"`
		Platform string  `json:"platform"`
		Accepts  []int32 `json:"accepts"`
	}
	if err = json.Unmarshal(req.Token, &params); err != nil {
		return
	}
	resp = &logic.ConnectResp{
		Mid:       params.Mid,
		RoomID:    params.RoomID,
		Accepts:   params.Accepts,
		Heartbeat: int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax),
	}
	if params.Key == "" {
		resp.Key = uuid.New().String()
	}
	if err = l.server.
		GetRedis().
		AddMapping(ctx, &model.Mapping{
			Mid:    resp.Mid,
			Key:    resp.Key,
			Server: req.Server,
		}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *connLogic) Disconnect(ctx context.Context, req *logic.DisconnectReq) (resp *logic.DisconnectResp, err error) {
	m := &model.Mapping{
		Mid:    req.Mid,
		Key:    req.Key,
		Server: req.Server,
	}
	resp = &logic.DisconnectResp{
		Has: false,
	}
	if resp.Has, err = l.server.
		GetRedis().
		DelMapping(ctx, m); err != nil {
		return nil, err
	}
	return resp, nil
}
func (l *connLogic) Heartbeat(ctx context.Context, req *logic.HeartbeatReq) (hbr *logic.HeartbeatResp, err error) {
	hbr = &logic.HeartbeatResp{
		Has: false,
	}
	if hbr.Has, err = l.server.GetRedis().ExpireMapping(ctx, &model.Mapping{
		Mid: req.Mid,
		Key: req.Key,
	}); err != nil {
		return nil, err
	}
	if !hbr.Has {
		if err = l.server.GetRedis().AddMapping(ctx, &model.Mapping{
			Key:    req.Key,
			Server: req.Server,
		}); err != nil {
			return nil, err
		}
	}
	return
}
func (l *connLogic) RenewOnline(ctx context.Context, req *logic.OnlineReq) (*logic.OnlineResp, error) {
	online := &model.Online{
		Server:    req.Server,
		RoomCount: req.RoomCount,
		Updated:   time.Now().Unix(),
	}
	if err := l.server.GetRedis().AddServerOnline(ctx, &model.Mapping{
		Server: req.Server,
		Online: online,
	}); err != nil {
		return nil, err
	}
	return &logic.OnlineResp{
		AllRoomCount: l.replicant.GetRoomCount(),
	}, nil
}
func (l *connLogic) Receive(ctx context.Context, req *logic.ReceiveReq) (*logic.ReceiveResp, error) {
	return nil, nil

}
func (l *connLogic) Nodes(ctx context.Context, req *logic.NodesReq) (*logic.NodesResp, error) {
	return nil, nil
}

func (l *connLogic) Ping(c context.Context) error {
	return l.server.Ping(c)
}
