package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	ReceiveMessage(ctx context.Context) (string, error)
	Close() error
}
type PubSubClient struct {
	*redis.PubSub
}

func (p *PubSubClient) ReceiveMessage(ctx context.Context) (string, error) {
	msg, err := p.PubSub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	return msg.Payload, nil
}

func (p *PubSubClient) Close() error {
	return p.PubSub.Close()
}
