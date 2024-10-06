package cmd

import (
	"context"
	"github.com/cuctemeh/rstream-consumer/internal/shutdown"
	"github.com/cuctemeh/rstream-consumer/internal/storage"
	"github.com/sumup-oss/go-pkgs/task"
	"log/slog"
	"os"
	"strconv"

	"github.com/palantir/stacktrace"
	"github.com/spf13/cobra"

	"github.com/cuctemeh/rstream-consumer/internal/config"
	"github.com/cuctemeh/rstream-consumer/internal/consumer"
	"github.com/cuctemeh/rstream-consumer/internal/monitor"
)

func NewConsumerCMD() *cobra.Command {
	CMDInstance := &cobra.Command{
		Use:   "consumer",
		Short: "start redis stream consumer",
		Long:  `start redis stream consumer`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
			loggerInstance := slog.New(logHandler)

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			redisClient, err := storage.NewClient(
				cfg.Redis,
			)
			if err != nil {
				return stacktrace.Propagate(err, "failed to connect to redis")
			}
			defer redisClient.Close()
			ctx := context.Background()

			// check if consumer count is not exceeded
			consumerCount, err := redisClient.SCard(
				ctx,
				cfg.Consumer.ConsumerIDsListName,
			)
			if err != nil {
				return stacktrace.Propagate(err, "failed to get consumer count")
			}

			if consumerCount >= cfg.Consumer.GroupSize {
				return stacktrace.NewError("Max consumer count exceeded.")
			}

			// generate consumer id based on the process id currently running.
			consumerID := strconv.Itoa(os.Getpid())
			// add consumer id to list
			_, err = redisClient.SAdd(
				ctx,
				cfg.Consumer.ConsumerIDsListName,
				consumerID,
			)
			if err != nil {
				return stacktrace.Propagate(err, "failed to add consumer id to list")
			}

			// remove consumer id from list
			defer redisClient.SRem(
				ctx,
				cfg.Consumer.ConsumerIDsListName,
				consumerID,
			)

			mon := monitor.NewMonitor(
				cfg.Monitoring,
				redisClient,
				cfg.ProcessedMessagesStreamName,
				loggerInstance,
			)

			pubsub := redisClient.Subscribe(ctx, cfg.Consumer.PublishedMessagesStreamName)
			defer pubsub.Close()

			redisStreamConsumer := consumer.NewRedisStreamConsumer(
				cfg.Consumer,
				redisClient,
				pubsub,
				cfg.ProcessedMessagesStreamName,
				consumerID,
				loggerInstance,
			)

			taskRunner := task.NewGroup()

			shutdownHandler := shutdown.NewShutdownHandler(cfg.ShutdownDeadline)

			taskRunner.Go(
				task.Task(redisStreamConsumer),
				task.Task(mon),
				task.Task(shutdownHandler),
			)

			err = taskRunner.Wait(ctx)
			if err != nil {
				return stacktrace.Propagate(err, "shutdown with error")
			}

			loggerInstance.Info("Shutdown successfully.")

			return nil
		},
	}
	return CMDInstance
}
