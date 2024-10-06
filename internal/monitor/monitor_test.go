package monitor_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cuctemeh/rstream-consumer/internal/config"
	"github.com/cuctemeh/rstream-consumer/internal/monitor"
	"github.com/cuctemeh/rstream-consumer/internal/testing/mocks"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
)

func TestMonitor_Run(t *testing.T) {
	mockRedisClient := new(mocks.MockRedisClient)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	cfg := config.MonitoringConfig{
		Interval: 1 * time.Second,
	}

	monitorInstance := monitor.NewMonitor(cfg, mockRedisClient, "processed_messages_stream", logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	// Set up expectations
	mockRedisClient.On("XLen", mock.Anything, "processed_messages_stream").Return(
		int64(10),
		nil,
	).Once().Run(
		func(args mock.Arguments) {
			done <- struct{}{}
		},
	)
	mockRedisClient.On("XLen", mock.Anything, "processed_messages_stream").Return(
		int64(20),
		nil,
	).Once().Run(
		func(args mock.Arguments) {
			done <- struct{}{}
		},
	)

	// Run the monitor in a separate goroutine
	go func() {
		defer wg.Done()
		monitorInstance.Run(ctx)
	}()

	// Wait for the monitor to process the expected calls
	<-done
	<-done

	// Cancel the context to stop the monitor
	cancel()

	// Wait for the monitor to finish
	wg.Wait()

	// Assert expectations
	mockRedisClient.AssertExpectations(t)
}
