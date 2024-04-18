package dao

import (
	"context"
	"strconv"
	"sync"

	"github.com/IBM/sarama"
	pb "github.com/yanglunara/im/api/logic"
	"github.com/yanglunara/im/internal/model"
	"google.golang.org/protobuf/proto"
)

var (
	_    model.PushKafka = (*dao)(nil)
	Dao  *dao
	once sync.Once
)

type dao struct {
	service model.ServerInter
}

func NewDao(si model.ServerInter) *dao {
	once.Do(func() {
		Dao = &dao{
			service: si,
		}
	})
	return Dao
}

func (d *dao) PushMessage(ctx context.Context, pm *model.PushMessage) (err error) {
	pushMsg := &pb.PushMessage{
		Type:      pb.PushMessage_Push,
		Operation: pm.Op,
		Servrice:  pm.Servrice,
		Keys:      pm.Keys,
		Msg:       pm.Message,
	}
	var (
		b []byte
	)
	// proto 序列化
	if b, err = proto.Marshal(pushMsg); err != nil {
		return
	}
	if _, _, err = d.service.
		GetKafka().
		GetSyncProducer().
		SendMessage(&sarama.ProducerMessage{
			Key:   sarama.StringEncoder(pm.Keys[0]),
			Topic: d.service.GetKafka().GetTopic(),
			Value: sarama.ByteEncoder(b),
		}); err != nil {
		return
	}
	return

}
func (d *dao) BraodcastMsg(ctx context.Context, brm *model.BraodcastMsg) (err error) {
	pushMsg := &pb.PushMessage{
		Type:      pb.PushMessage_Braodcast,
		Operation: brm.Op,
		Room:      brm.Room,
		Msg:       brm.Message,
	}
	var (
		b []byte
	)
	if b, err = proto.Marshal(pushMsg); err != nil {
		return
	}
	kafka := d.service.GetKafka()
	if _, _, err = kafka.GetSyncProducer().SendMessage(&sarama.ProducerMessage{
		Topic: kafka.GetTopic(),
		Value: sarama.ByteEncoder(b),
		// 以操作码作为key
		Key: sarama.StringEncoder(strconv.FormatInt(int64(brm.Op), 10)),
	}); err != nil {
		return
	}
	return
}

func (d *dao) BraodcastRoomMsg(ctx context.Context, brm *model.BraodcastMsg) (err error) {
	pushMsg := &pb.PushMessage{
		Type:      pb.PushMessage_Room,
		Operation: brm.Op,
		Room:      brm.Room,
		Msg:       brm.Message,
	}
	var (
		b []byte
	)
	if b, err = proto.Marshal(pushMsg); err != nil {
		return
	}
	kafka := d.service.GetKafka()
	if _, _, err = kafka.GetSyncProducer().SendMessage(&sarama.ProducerMessage{
		Topic: kafka.GetTopic(),
		Value: sarama.ByteEncoder(b),
		// 以操作码作为key
		Key: sarama.StringEncoder(strconv.FormatInt(int64(brm.Op), 10)),
	}); err != nil {
		return
	}
	return
}
