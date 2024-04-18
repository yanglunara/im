package model

import (
	"context"
)

type PushKafka interface {
	PushMessage(ctx context.Context, pm *PushMessage) (err error)

	BraodcastRoomMsg(ctx context.Context, brm *BraodcastMsg) (err error)

	BraodcastMsg(ctx context.Context, brm *BraodcastMsg) (err error)
}

type Stoper interface {
	Close() error
	Ping(ctx context.Context) error
}
