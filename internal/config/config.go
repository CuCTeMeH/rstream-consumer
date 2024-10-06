package config

import (
	"github.com/palantir/stacktrace"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ProcessedMessagesStreamName string        `mapstructure:"processed_messages_stream_name"`
	ShutdownDeadline            time.Duration `mapstructure:"shutdown_deadline"`

	Redis      RedisConfig      `mapstructure:"redis"`
	Consumer   ConsumerConfig   `mapstructure:"consumer"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
}

func NewConfig() (*Config, error) {
	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("monitoring.interval", "1s")
	viper.SetDefault("consumer.group_size", 1)
	viper.SetDefault("consumer.consumer_ids_list_name", "consumer:ids")
	viper.SetDefault("consumer.published_messages_stream_name", "messages:published")
	viper.SetDefault("processed_messages_stream_name", "messages:processed")
	viper.SetDefault("shutdown_deadline", "1s")

	// Bind environment variables
	viper.BindEnv("redis.address", "REDIS_ADDRESS")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("monitoring.interval", "MONITORING_INTERVAL")
	viper.BindEnv("consumer.group_size", "CONSUMER_GROUP_SIZE")
	viper.BindEnv("consumer.consumer_ids_list_name", "CONSUMER_IDS_LIST_NAME")
	viper.BindEnv(
		"consumer.published_messages_stream_name",
		"CONSUMER_PUBLISHED_MESSAGES_STREAM_NAME",
	)
	viper.BindEnv("processed_messages_stream_name", "PROCESSED_MESSAGES_STREAM_NAME")
	viper.BindEnv("shutdown_deadline", "SHUTDOWN_DEADLINE")

	viper.AutomaticEnv()

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to unmarshal config")
	}

	return &cfg, nil
}
