package plugin

import (
	"context"
	"sync"

	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/model"
)

var (
	_        model.ServerInter = (*service)(nil)
	Servrice *service
	once     sync.Once
)

type service struct {
	kafka model.KafkaInter
	redis model.RedisInter
}

type Option func(*service)

func WithKafka(kafka *conf.Kafka) Option {
	return func(s *service) {
		s.kafka = newKafka(kafka)
	}
}

func WithRedis(redis *conf.Redis) Option {
	return func(s *service) {
		s.redis = newPluginRedis(
			redis,
			newRedisSendPlugin(),
		)
	}
}

func InitService(opt ...Option) model.ServerInter {
	once.Do(func() {
		s := &service{}
		for _, o := range opt {
			o(s)
		}
		Servrice = s
	})
	return Servrice
}

func (s *service) GetKafka() model.KafkaInter {
	return s.kafka
}
func (s *service) GetRedis() model.RedisInter {
	return s.redis
}

func (s *service) Close() error {
	return s.redis.GetRedis().Close()
}

func (s *service) Ping(ctx context.Context) (err error) {
	return s.redis.Ping(ctx)
}
