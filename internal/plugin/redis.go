package plugin

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/model"
	"github.com/zhenjl/cityhash"
)

var (
	_ model.RedisInter      = (*serverRedis)(nil)
	_ model.RedisSendPlugin = (*redisSendPlugin)(nil)
)

type serverRedis struct {
	rp     *redis.Pool
	expire int32
	rsp    model.RedisSendPlugin
	maxTry int
}

func newPluginRedis(c *conf.Redis, rsp model.RedisSendPlugin) model.RedisInter {
	return &serverRedis{
		rp: &redis.Pool{
			MaxIdle:     c.Idle,
			MaxActive:   c.Active,
			IdleTimeout: c.IdleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial(c.Network, c.Addr,
					redis.DialPassword(c.Auth),
					redis.DialReadTimeout(c.ReadTimeout),
					redis.DialWriteTimeout(c.WriteTimeout),
					redis.DialConnectTimeout(c.DialTimeout),
				)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
		expire: int32(c.Expire / time.Second),
		rsp:    rsp,
		maxTry: 1,
	}
}

func (r *serverRedis) AddMapping(c context.Context, m *model.Mapping) (err error) {
	conn := r.rp.Get()
	defer conn.Close()
	rpo := &model.SendPluginOption{
		MaxTry:    r.maxTry + 1,
		RedisKey1: Key[string, int64]("mid", m.Mid).Sprintf(),
		RedisKey2: Key[string, string]("key", m.Key).Sprintf(),
		Key:       m.Key,
		Field:     m.Server,
		Expire:    r.expire,
		Mid:       int(m.Mid),
	}
	var (
		n int
	)
	if n, err = r.rsp.DoubleSet(conn, rpo); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	// 重试
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return nil
}

func (r *serverRedis) ExpireMapping(c context.Context, m *model.Mapping) (check bool, err error) {
	conn := r.rp.Get()
	defer conn.Close()
	n := r.maxTry
	if m.Mid > 0 {
		if err = r.rsp.Expire(
			conn,
			Key[string, int64]("mid", m.Mid).Sprintf(),
			r.expire,
		); err != nil {
			return
		}
		n++
	}
	if err = r.rsp.Expire(
		conn,
		Key[string, string]("key", m.Key).Sprintf(),
		r.expire,
	); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	// 重试
	for i := 0; i < n; i++ {
		if check, err = redis.Bool(conn.Receive()); err != nil {
			return
		}
	}
	return
}
func (r *serverRedis) DelMapping(c context.Context, m *model.Mapping) (has bool, err error) {
	conn := r.rp.Get()
	defer conn.Close()
	rpo := &model.SendPluginOption{
		MaxTry:    r.maxTry,
		RedisKey1: Key[string, int64]("mid", m.Mid).Sprintf(),
		RedisKey2: Key[string, string]("key", m.Key).Sprintf(),
		Key:       m.Key,
		Expire:    r.expire,
		Mid:       int(m.Mid),
	}
	var (
		n int
	)
	if n, err = r.rsp.DoubleDel(conn, rpo); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			return
		}
	}
	return false, nil
}
func (r *serverRedis) ServersByKeys(c context.Context, keys []string) (res []string, err error) {
	conn := r.rp.Get()
	defer conn.Close()
	var (
		args = make([]interface{}, 0, len(keys))
	)
	for _, key := range keys {
		args = append(args, Key[string, string]("key", key).Sprintf())
	}
	if res, err = redis.Strings(conn.Do("MGET", args...)); err != nil {
		return
	}
	return
}
func (r *serverRedis) KeysByMids(c context.Context, mids []int64) (res *model.ResKeysByMids, err error) {
	conn := r.rp.Get()
	defer conn.Close()
	res.Res = make(map[string]string)
	for _, mid := range mids {
		if err = conn.Send("HGETALL",
			Key[string, int64]("mid", mid).Sprintf(), "key"); err != nil {
			return
		}
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for idx := 0; idx < len(mids); idx++ {
		var (
			re map[string]string
		)
		if re, err = redis.StringMap(conn.Receive()); err != nil {
			return
		}
		if len(re) > 0 {
			res.OlMids = append(res.OlMids, mids[idx])
		}
		for k, v := range re {
			res.Res[k] = v
		}
	}
	return
}
func (r *serverRedis) AddServerOnline(c context.Context, m *model.Mapping) (err error) {
	roomsMap := map[uint32]map[string]int32{}
	for room, count := range m.Online.RoomCount {
		// 哈希取模 分配房间
		rMap := roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64]
		if rMap == nil {
			rMap = make(map[string]int32)
			roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64] = rMap
		}
		rMap[room] = count
	}
	key := Key[string, string]("ol", m.Server).Sprintf()
	for hashKey, roomMap := range roomsMap {
		if err = r.addServerOnline(c, key, strconv.FormatInt(int64(hashKey), 10), &model.Online{
			RoomCount: roomMap,
			Server:    m.Online.Server,
			Updated:   m.Online.Updated,
		}); err != nil {
			return
		}
	}
	return nil
}

func (r *serverRedis) addServerOnline(_ context.Context, key, hashKey string, online *model.Online) (err error) {
	conn := r.rp.Get()
	defer conn.Close()
	b, _ := json.Marshal(online)
	if err = conn.Send("HSET", key, hashKey, b); err != nil {
		return
	}
	if err = r.rsp.Expire(conn, key, r.expire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return

}
func (r *serverRedis) ServerOnline(c context.Context, server string) (online *model.Online, err error) {
	online = &model.Online{
		RoomCount: make(map[string]int32),
	}
	key := Key[string, string]("ol", server).Sprintf()
	for i := 0; i < 64; i++ {
		if ol, err := r.serverOnline(c, key, strconv.FormatInt(int64(i), 10)); err == nil && ol != nil {
			online.Server = ol.Server
			if ol.Updated > online.Updated {
				online.Updated = ol.Updated
			}
			for room, count := range ol.RoomCount {
				online.RoomCount[room] = count
			}
		}
	}
	return
}

func (r *serverRedis) serverOnline(_ context.Context, key, hashKey string) (online *model.Online, err error) {
	conn := r.rp.Get()
	defer conn.Close()
	var (
		b []byte
	)
	if b, err = redis.Bytes(conn.Do("HGET", key, hashKey)); err != nil {
		return
	}
	online = new(model.Online)
	if err = json.Unmarshal(b, online); err != nil {
		return
	}
	return
}

func (r *serverRedis) DelServerOnline(c context.Context, server string) (err error) {
	conn := r.rp.Get()
	defer conn.Close()
	if _, err = conn.Do("DEL", Key[string, string]("ol", server).Sprintf()); err != nil {
		return
	}
	return nil
}

func (r *serverRedis) GetRedis() *redis.Pool {
	return r.rp
}
func (r *serverRedis) Ping(ctx context.Context) (err error) {
	conn := r.rp.Get()
	defer conn.Close()
	_, err = conn.Do("SET", "PING", "PONG")
	return err
}
