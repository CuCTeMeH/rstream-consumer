package monitor

import (
	"context"
	"log/slog"
	"time"

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
) *Monitor {
	return &Monitor{
		Interval:                  cfg.Interval,
		redisClient:               redisClient,
		processMessagesStreamName: processMessagesStreamName,
		logger:                    logger,
		lastCount:                 0,
	}
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
			m.lastCount = count

			m.logger.Info(
				"messages processed",
				"message_count",
				messagesProcessed,
				"interval",
				m.Interval,
			)
		}
	}
}
