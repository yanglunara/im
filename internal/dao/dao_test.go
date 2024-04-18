package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/yanglunara/im/cmd/logic/conf"
	"github.com/yanglunara/im/internal/model"
	"github.com/yanglunara/im/internal/plugin"
)

func TestNewDao(t *testing.T) {
	plugin.InitService(plugin.WithKafka(&conf.Kafka{
		Topic:   "im-push-topic",
		Group:   "im-push-topic-test-group",
		Brokers: []string{"127.0.0.1:9092"},
	}))
	ctx := context.Background()
	// fmt.Printf("TestNewDao %v", server.Servrice.GetKafkaServer().GetKafkaTopic())
	dao := NewDao(plugin.Servrice)
	fmt.Println(dao.PushMessage(ctx, &model.PushMessage{
		Op:       1,
		Servrice: "im-push-topic",
		Keys:     []string{"Test"},
		Message:  []byte("Test"),
	}))
}
