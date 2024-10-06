package config

import "time"

type MonitoringConfig struct {
	Interval time.Duration `mapstructure:"interval"`
}
