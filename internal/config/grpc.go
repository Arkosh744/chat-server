package config

import (
	"fmt"
	"net"
	"os"
)

var _ GRPCConfig = (*grpcConfig)(nil)

const (
	grpcEnvHost = "GRPC_HOST"
	grpcEnvPort = "GRPC_PORT"
)

type GRPCConfig interface {
	GetHost() string
}

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcEnvHost)
	port := os.Getenv(grpcEnvPort)
	if port == "" || host == "" {
		return nil, fmt.Errorf("grpc addr is not set")
	}

	return &grpcConfig{host: host, port: port}, nil
}

func (c *grpcConfig) GetHost() string {
	return net.JoinHostPort(c.host, c.port)
}
