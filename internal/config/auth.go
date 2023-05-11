package config

import (
	"github.com/kelseyhightower/envconfig"
)

var _ AuthConfig = (*authConfig)(nil)

const (
	authEnvPrefix = "AUTH"
)

type AuthConfig interface {
	GetPort() string
}

type authConfig struct {
	Port string `required:"true"`
}

func NewAuthConfig() (*authConfig, error) {
	var cfg authConfig
	err := envconfig.Process(authEnvPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *authConfig) GetPort() string {
	return c.Port
}
