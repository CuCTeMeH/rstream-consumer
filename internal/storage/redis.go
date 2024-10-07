package storage

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"

	"github.com/cuctemeh/rstream-consumer/internal/config"
)

type RedisClient interface {
	Subscribe(ctx context.Context, channel string) PubSub

	SIsMember(ctx context.Context, set string, member interface{}) (bool, error)

	XAdd(ctx context.Context, args *redis.XAddArgs) (string, error)
	XLen(ctx context.Context, key string) (int64, error)
	XRead(ctx context.Context, args *redis.XReadArgs) ([]redis.XStream, error)

	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SCard(ctx context.Context, key string) (int64, error)

	Close() error
}

type client struct {
	*redis.Client
}

func NewClient(cfg config.RedisConfig) (RedisClient, error) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
	)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to connect to redis")
	}
	return &client{rdb}, nil
}

func (c *client) Subscribe(ctx context.Context, channel string) PubSub {
	redisPubSub := c.Client.Subscribe(ctx, channel)
	return &PubSubClient{redisPubSub}
}

func (c *client) XRead(ctx context.Context, args *redis.XReadArgs) ([]redis.XStream, error) {
	return c.Client.XRead(ctx, args).Result()
}

func (c *client) SIsMember(ctx context.Context, set string, member interface{}) (bool, error) {
	return c.Client.SIsMember(ctx, set, member).Result()
}

func (c *client) XAdd(ctx context.Context, args *redis.XAddArgs) (string, error) {
	return c.Client.XAdd(ctx, args).Result()
}

func (c *client) XLen(ctx context.Context, key string) (int64, error) {
	return c.Client.XLen(ctx, key).Result()
}

func (c *client) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.Client.SAdd(ctx, key, members...).Result()
}

func (c *client) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.Client.SRem(ctx, key, members...).Result()
}

func (c *client) SCard(ctx context.Context, key string) (int64, error) {
	return c.Client.SCard(ctx, key).Result()
}

func (c *client) Close() error {
	return c.Client.Close()
}
