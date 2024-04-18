package plugin

import (
	"os"
	"testing"
	"time"

	"github.com/yanglunara/im/cmd/logic/conf"
)

func TestMain(m *testing.M) {
	InitService(WithKafka(&conf.Kafka{
		Topic:   "im-push-topic",
		Group:   "im-push-topic-test-group",
		Brokers: []string{"127.0.0.1:9092"},
	}),
		WithRedis(
			&conf.Redis{
				Network:      "tcp",
				Addr:         "127.0.0.1:26379",
				Auth:         "",
				Active:       10,
				Idle:         10,
				DialTimeout:  10 * time.Second,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  10 * time.Second,
				Expire:       600 * time.Second,
			},
		),
	)

	os.Exit(m.Run())
}
