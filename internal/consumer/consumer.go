package consumer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/cuctemeh/rstream-consumer/internal/config"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
	"log/slog"

	"github.com/cuctemeh/rstream-consumer/internal/storage"
)

type Consumer struct {
	publishedMessagesStreamName string
	processedMessagesStreamName string
	consumerID                  string
	redisClient                 storage.RedisClient
	pubsub                      storage.PubSub
	logger                      *slog.Logger
}

func NewRedisStreamConsumer(
	cfg config.ConsumerConfig,
	redisClient storage.RedisClient,
	pubsub storage.PubSub,
	processedMessagesStreamName string,
	consumerID string,
	logger *slog.Logger,
) *Consumer {
	return &Consumer{
		publishedMessagesStreamName: cfg.PublishedMessagesStreamName,
		processedMessagesStreamName: processedMessagesStreamName,
		consumerID:                  consumerID,
		redisClient:                 redisClient,
		pubsub:                      pubsub,
		logger:                      logger,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	msgCh := make(chan string)
	errCh := make(chan error)

	go func(ctx context.Context) {
		for {
			msg, errReceive := c.pubsub.ReceiveMessage(ctx)
			if errReceive != nil {
				errCh <- errReceive
				return
			}
			msgCh <- msg
		}
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context done")
			return stacktrace.Propagate(ctx.Err(), "context done")
		case errMsg := <-errCh:
			c.logger.InfoContext(ctx, "error reading group", "err", errMsg)
			return stacktrace.Propagate(errMsg, "error reading group")
		case msg := <-msgCh:
			errMsg := c.processMessage(ctx, msg, c.consumerID)
			if errMsg != nil {
				return stacktrace.Propagate(errMsg, "failed to process message")
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg string, consumerID string) error {
	var message Message
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal message")
	}

	message.ConsumerID = consumerID

	// Check if the message has already been consumed
	isMember, err := c.redisClient.SIsMember(ctx, "consumed_messages", message.MessageID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to check if message has been consumed")
	}
	if isMember {
		// The message has already been consumed, so ignore it
		return nil
	}

	// The message has not been consumed, so add it to the set of consumed messages
	_, err = c.redisClient.SAdd(ctx, "consumed_messages", message.MessageID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to add message to set of consumed messages")
	}

	// cleanup - remove message from set of consumed messages once processed
	defer func(ctx context.Context, messageID string) {
		_, errSRem := c.redisClient.SRem(ctx, "consumed_messages", messageID)
		if errSRem != nil {
			c.logger.InfoContext(
				ctx,
				"error removing message from set of consumed messages",
				"err",
				errSRem,
			)
		}
	}(ctx, message.MessageID)

	// add hash value of consumer ID and message ID to random property
	h := sha256.New()
	h.Write([]byte(consumerID + message.MessageID))
	hash := hex.EncodeToString(h.Sum(nil))

	// add random property and value to message object
	message.RandomProperty = hash

	// store message data in Redis stream
	_, err = c.redisClient.XAdd(
		ctx, &redis.XAddArgs{
			Stream: c.processedMessagesStreamName,
			Values: map[string]interface{}{
				"message_id":      message.MessageID,
				"random_property": message.RandomProperty,
				"consumer_id":     consumerID,
			},
		},
	)
	if err != nil {
		return stacktrace.Propagate(err, "failed to store message data in Redis stream")
	}

	c.logger.InfoContext(
		ctx,
		"message processed",
		"message_id", message.MessageID,
	)
	return nil
}
