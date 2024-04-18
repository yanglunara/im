package plugin

import (
	kafka "github.com/IBM/sarama"
	"github.com/yanglunara/im/cmd/logic/conf"
	"github.com/yanglunara/im/internal/model"
)

type serverKafka struct {
	server *conf.Kafka
}

func newKafka(kafka *conf.Kafka) model.KafkaInter {
	return &serverKafka{
		server: kafka,
	}
}

func (k *serverKafka) GetSyncProducer() kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll
	kc.Producer.Retry.Max = 10
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(k.server.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

func (k *serverKafka) GetTopic() string {
	return k.server.Topic
}
