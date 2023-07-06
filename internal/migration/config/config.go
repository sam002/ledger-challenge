package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"go.uber.org/zap"
)

type MigrationConfig struct {
	DSN   string `env:"DSN,unset"`
	Debug bool   `env:"DEBUG" envDefault:"true"`
}

func GetConfig(logger *zap.Logger) MigrationConfig {
	cfg := MigrationConfig{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		logger.Error("Cannot parse log", zap.Error(err))
	}
	if cfg.Debug {
		nLogger, err := zap.NewDevelopment()
		*logger = *nLogger
		if err != nil {
			fmt.Printf("%+v\n", err)
			logger.Error("Cannot parse log", zap.Error(err))
		}
	}

	logger.Info("Successful load config")
	logger.Debug("Config environment", zap.Any("Config", cfg))

	return cfg
}
