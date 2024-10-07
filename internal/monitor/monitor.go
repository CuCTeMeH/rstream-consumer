package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/palantir/stacktrace"

	"github.com/cuctemeh/rstream-consumer/internal/config"
	"github.com/cuctemeh/rstream-consumer/internal/storage"
)

type Monitor struct {
	Interval                  time.Duration
	redisClient               storage.RedisClient
	processMessagesStreamName string
	logger                    *slog.Logger
	lastCount                 int64
}

func NewMonitor(
	cfg config.MonitoringConfig,
	redisClient storage.RedisClient,
	processMessagesStreamName string,
	logger *slog.Logger,
) (*Monitor, error) {
	if cfg.Interval < time.Second {
		return nil, stacktrace.NewError("interval should be greater than 1 second")
	}

	return &Monitor{
		Interval:                  cfg.Interval,
		redisClient:               redisClient,
		processMessagesStreamName: processMessagesStreamName,
		logger:                    logger,
		lastCount:                 0,
	}, nil
}

func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			count, err := m.redisClient.XLen(ctx, m.processMessagesStreamName)
			if err != nil {
				return err
			}

			messagesProcessed := count - m.lastCount
			msgsPerSecond := float64(messagesProcessed) / m.Interval.Seconds()
			m.lastCount = count

			m.logger.Info(
				"messages processed per second",
				"message_count",
				msgsPerSecond,
				"interval",
				m.Interval,
			)
		}
	}
}
