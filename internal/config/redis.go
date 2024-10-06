package config

type RedisConfig struct {
	Addr     string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
