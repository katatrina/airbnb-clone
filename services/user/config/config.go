package config

import (
	"errors"
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

// Validate checks that all required configuration is present.
// This is called after loading to fail fast if config is missing.
func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return errors.New("SERVER_PORT is required")
	}
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		return errors.New("REDIS_URL is required")
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long")
	}
	if c.JWTExpiry <= 0 {
		return errors.New("JWT_EXPIRY must be greater than 0")
	}

	return nil
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

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
