package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string        `mapstructure:"SERVER_PORT"`
	DatabaseURL string        `mapstructure:"DATABASE_URL"`
	RedisURL    string        `mapstructure:"REDIS_URL"`
	JWTSecret   string        `mapstructure:"JWT_SECRET"`
	JWTExpiry   time.Duration `mapstructure:"JWT_EXPIRY"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
