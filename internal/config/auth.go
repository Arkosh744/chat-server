package config

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

var _ AuthConfig = (*authConfig)(nil)

const (
	authEnvPrefix = "AUTH"
)

type AuthConfig interface {
	GetHost() string
}

type authConfig struct {
	Port string `required:"true"`
	Host string `required:"true"`
}

func NewAuthConfig() (*authConfig, error) {
	var cfg authConfig
	err := envconfig.Process(authEnvPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *authConfig) GetHost() string {
	return net.JoinHostPort(c.Host, c.Port)
}
