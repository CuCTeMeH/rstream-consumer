package consumer_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"

	"github.com/cuctemeh/rstream-consumer/internal/config"
	"github.com/cuctemeh/rstream-consumer/internal/consumer"
	"github.com/cuctemeh/rstream-consumer/internal/testing/mocks"
)

func TestConsumer_Run_Success(t *testing.T) {
	mockDB := new(mocks.MockRedisClient)

	ConsumerIDInt := os.Getpid()
	consumerID := strconv.Itoa(ConsumerIDInt)

	mockDB.On("SAdd", mock.Anything, "messages:processing", []interface{}{"1"}).Return(
		int64(1),
		nil,
	)
	mockDB.On("SRem", mock.Anything, "messages:processing", []interface{}{"1"}).Return(
		int64(1),
		nil,
	)
	mockDB.On("SIsMember", mock.Anything, "messages:processing", "1").Return(false, nil)
	mockDB.On("XAdd", mock.Anything, mock.Anything).Return("1", nil)

	exampleMsg := &consumer.Message{
		MessageID:  "1",
		ConsumerID: consumerID,
	}

	exampleMsgStr, _ := json.Marshal(exampleMsg)

	pubsubMock := new(mocks.MockRedisPubSub)
	pubsubMock.On(
		"ReceiveMessage",
		mock.Anything,
	).Return(string(exampleMsgStr), nil)

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() = %v, want nil", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	loggerInstance := slog.New(logHandler)

	consumer := consumer.NewRedisStreamConsumer(
		cfg.Consumer,
		mockDB,
		pubsubMock,
		cfg.ProcessedMessagesStreamName,
		consumerID,
		loggerInstance,
	)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	go consumer.Run(ctx)

	<-ctx.Done()

	mockDB.AssertExpectations(t)
	pubsubMock.AssertExpectations(t)
}

func TestConsumer_Run_E2E(t *testing.T) {
	mockDB := new(mocks.MockRedisClient)

	ConsumerIDInt := os.Getpid()
	consumerID := strconv.Itoa(ConsumerIDInt)

	mockDB.On("SAdd", mock.Anything, "messages:processing", []interface{}{"1"}).Return(
		int64(1),
		nil,
	)
	mockDB.On("SRem", mock.Anything, "messages:processing", []interface{}{"1"}).Return(
		int64(1),
		nil,
	)
	mockDB.On("SIsMember", mock.Anything, "messages:processing", "1").Return(false, nil)
	mockDB.On("XAdd", mock.Anything, mock.Anything).Return("1", nil)

	exampleMsg := &consumer.Message{
		MessageID:  "1",
		ConsumerID: consumerID,
	}

	exampleMsgStr, _ := json.Marshal(exampleMsg)

	pubsubMock := new(mocks.MockRedisPubSub)
	pubsubMock.On(
		"ReceiveMessage",
		mock.Anything,
	).Return(string(exampleMsgStr), nil)

	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() = %v, want nil", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	loggerInstance := slog.New(logHandler)

	consumer := consumer.NewRedisStreamConsumer(
		cfg.Consumer,
		mockDB,
		pubsubMock,
		cfg.ProcessedMessagesStreamName,
		consumerID,
		loggerInstance,
	)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	go consumer.Run(ctx)

	<-ctx.Done()

	mockDB.AssertExpectations(t)
	pubsubMock.AssertExpectations(t)

	// Mock XRead to verify the processed message
	mockDB.On("XRead", mock.Anything, mock.Anything).Return(
		[]redis.XStream{
			{
				Stream: cfg.ProcessedMessagesStreamName,
				Messages: []redis.XMessage{
					{
						ID: "1",
						Values: map[string]interface{}{
							"message_id":      "1",
							"random_property": exampleMsg.RandomProperty,
							"consumer_id":     consumerID,
						},
					},
				},
			},
		}, nil,
	)

	// Read the processed message from the stream
	streams, err := mockDB.XRead(
		ctx, &redis.XReadArgs{
			Streams: []string{cfg.ProcessedMessagesStreamName, "0"},
			Count:   1,
			Block:   0,
		},
	)
	if err != nil {
		t.Fatalf("XRead() = %v, want nil", err)
	}

	if len(streams) != 1 || len(streams[0].Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(streams[0].Messages))
	}

	processedMsg := streams[0].Messages[0]
	if processedMsg.Values["message_id"] != "1" {
		t.Errorf("Expected message_id to be '1', got %v", processedMsg.Values["message_id"])
	}
	if processedMsg.Values["random_property"] != exampleMsg.RandomProperty {
		t.Errorf(
			"Expected random_property to be %v, got %v",
			exampleMsg.RandomProperty,
			processedMsg.Values["random_property"],
		)
	}
	if processedMsg.Values["consumer_id"] != consumerID {
		t.Errorf(
			"Expected consumer_id to be %v, got %v",
			consumerID,
			processedMsg.Values["consumer_id"],
		)
	}
}
