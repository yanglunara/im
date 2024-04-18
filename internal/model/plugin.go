package model

import (
	"context"

	kafka "github.com/IBM/sarama"
	"github.com/gomodule/redigo/redis"
)

type ServerInter interface {
	GetKafka() KafkaInter
	GetRedis() RedisInter
	Stoper
}

type KafkaInter interface {
	GetSyncProducer() kafka.SyncProducer
	GetTopic() string
}

type ResKeysByMids struct {
	Res    map[string]string
	OlMids []int64
}

type RedisInter interface {
	AddMapping(c context.Context, m *Mapping) (err error)
	ExpireMapping(c context.Context, m *Mapping) (check bool, err error)
	DelMapping(c context.Context, m *Mapping) (has bool, err error)
	ServersByKeys(c context.Context, keys []string) (res []string, err error)
	KeysByMids(c context.Context, mids []int64) (res *ResKeysByMids, err error)
	AddServerOnline(c context.Context, m *Mapping) (err error)
	ServerOnline(c context.Context, server string) (online *Online, err error)
	DelServerOnline(c context.Context, server string) (err error)
	GetRedis() *redis.Pool
	Ping(ctx context.Context) (err error)
}

type RedisSendPlugin interface {
	DoubleSet(conn redis.Conn, spo *SendPluginOption) (maxTry int, err error)
	Expire(conn redis.Conn, key string, expire int32) (err error)
	DoubleDel(conn redis.Conn, spo *SendPluginOption) (maxTry int, err error)
}
