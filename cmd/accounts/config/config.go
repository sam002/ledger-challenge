package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"go.uber.org/zap"
)

type Config struct {
	Port  int    `env:"PORT" envDefault:"3000"`
	DSN   string `env:"DSN,unset"`
	Debug bool   `env:"DEBUG" envDefault:"true"`
}

func GetConfig(logger *zap.Logger) Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		logger.Error("Cannot parse log", zap.Error(err))
	}
	if cfg.Debug {
		if logger, err := zap.NewDevelopment(); err != nil {
			fmt.Printf("%+v\n", err)
			logger.Error("Cannot parse log", zap.Error(err))
		}
	}

	logger.Info("Successful load config")
	logger.Debug("Successful load config", zap.Any("Config", cfg))

	return cfg
}
