package app

import (
	"context"
	"github.com/Arkosh744/auth-service-api/internal/api/user_v1"
	userService "github.com/Arkosh744/auth-service-api/internal/service/user"
	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/Arkosh744/chat-server/internal/closer"
	"github.com/Arkosh744/chat-server/internal/config"
	"github.com/Arkosh744/chat-server/internal/service/chat"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type serviceProvider struct {
	authConfig  config.AuthConfig
	authClient  auth.Client
	chatService *chat.Service

	log *zap.SugaredLogger
}

func (s *serviceProvider) newAuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			s.log.Fatal("failed to get auth config", zap.Error(err))
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) GetAuthClient(_ context.Context) auth.Client {
	if s.authClient == nil {
		conn, err := grpc.Dial(s.authConfig.GetPort(), grpc.WithDefaultCallOptions())
		if err != nil {
			s.log.Fatalf("failed to connect %s: %s", s.authConfig.GetPort(), err)
		}
		closer.Add(conn.Close)

		client := user_v1.NewUserClient(conn)
		s.authClient = authClient.NewClient(client)
	}

	return s.authClient
}

func newServiceProvider(log *zap.SugaredLogger) *serviceProvider {
	return &serviceProvider{log: log}
}

func (s *serviceProvider) GetUserService(ctx context.Context) userService.Service {
	if s.userService == nil {
		s.userService = userService.NewService(s.GetUserRepo(ctx), s.log)
	}

	return s.userService
}

func (s *serviceProvider) GetUserImpl(ctx context.Context) *userV1.Implementation {
	if s.userImpl == nil {
		s.userImpl = userV1.NewImplementation(s.GetUserService(ctx), s.log)
	}

	return s.userImpl
}
