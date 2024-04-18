package plugin

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/yanglunara/im/internal/model"
)

type redisSendPlugin struct {
}

func newRedisSendPlugin() model.RedisSendPlugin {
	return &redisSendPlugin{}
}
func (r *redisSendPlugin) Set(conn redis.Conn, key, value string, expire int32) (err error) {
	if err = conn.Send("SET", key, value); err != nil {
		return
	}
	return r.Expire(conn, key, expire)
}

func (r *redisSendPlugin) Expire(conn redis.Conn, key string, expire int32) (err error) {
	return conn.Send("EXPIRE", key, expire)
}

func (r *redisSendPlugin) HSet(conn redis.Conn, key, field, value string, expire int32) (err error) {
	if err = conn.Send("HSET", key, field, value); err != nil {
		return
	}
	return r.Expire(conn, key, expire)
}

func (r *redisSendPlugin) DoubleSet(conn redis.Conn, spo *model.SendPluginOption) (maxTry int, err error) {
	var (
		key1   string = spo.RedisKey1
		key2   string = spo.RedisKey2
		value  string = spo.Key
		field  string = spo.Field
		expire int32  = spo.Expire
	)
	n := spo.MaxTry
	if spo.Mid > 0 {
		if err = conn.Send("HSET", key1, value, spo.Field); err != nil {
			return
		}
		if err = conn.Send("EXPIRE", key1, expire); err != nil {
			return
		}
		n++
	}
	if err = conn.Send("SET", key2, field); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", key2, expire); err != nil {
		return
	}
	return n, nil
}

func (r *redisSendPlugin) DoubleDel(conn redis.Conn, spo *model.SendPluginOption) (maxTry int, err error) {
	var (
		key1  = spo.RedisKey1
		key2  = spo.RedisKey2
		field = spo.Field
	)
	n := spo.MaxTry
	if spo.Mid > 0 {
		if err = conn.Send("HDEL", key1, field); err != nil {
			return
		}
		n++
	}
	if err = conn.Send("DEL", key2); err != nil {
		return
	}
	return n, nil
}
func (r *redisSendPlugin) Del(conn redis.Conn, key string) (err error) {
	return conn.Send("DEL", key)
}

type keys[K comparable, V any] struct {
	Key   K
	Value V
}

func Key[K comparable, V any](key K, value V) keys[K, V] {
	return keys[K, V]{Key: key, Value: value}
}

func (r keys[K, V]) Sprintf() string {
	return fmt.Sprintf("%v_%v", r.Key, r.Value)
}
