package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	// Prioritize reading from environment variables
	viper.AutomaticEnv()

	// Optional: load .env for local dev
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	return nil
}
