package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

var _ LogConfig = (*logConfig)(nil)

const (
	loggerEnvPrefix = "LOG"
)

type LogConfig interface {
	GetPreset() string
}

type logConfig struct {
	Preset string `default:"dev"`
}

func NewLogConfig() (*logConfig, error) {
	var cfg logConfig
	err := envconfig.Process(loggerEnvPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *logConfig) GetPreset() string {
	return c.Preset
}

func getLogConfig() LogConfig {
	cfg, err := NewLogConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %s", err.Error())
	}

	return cfg
}

func SelectLogger() (*zap.Logger, error) {
	switch getLogConfig().GetPreset() {
	case "prod":
		return zap.NewProduction()
	case "dev":
		return zap.NewDevelopment()
	default:
		log.Println("unknown logger preset, using development preset")
		return zap.NewDevelopment()
	}
}
