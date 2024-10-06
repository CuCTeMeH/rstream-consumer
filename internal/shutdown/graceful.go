package shutdown

import (
	"context"
	"github.com/palantir/stacktrace"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownHandler struct {
	shutdownDeadline time.Duration
}

func NewShutdownHandler(shutdownDeadline time.Duration) *ShutdownHandler {
	return &ShutdownHandler{
		shutdownDeadline: shutdownDeadline,
	}
}

func (s *ShutdownHandler) Run(ctx context.Context) error {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-osSignalChan:

		go func() {
			time.Sleep(s.shutdownDeadline)
			os.Exit(1)
		}()

		go func() {
			sig = <-osSignalChan
			os.Exit(1)
		}()

		return stacktrace.NewError("shutdown by sig %+v", sig)
	case <-ctx.Done():
		return nil
	}
}
