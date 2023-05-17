package app

import (
	"context"
	"github.com/Arkosh744/auth-service-api/pkg/user_v1"
	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/Arkosh744/chat-server/internal/closer"
	"github.com/Arkosh744/chat-server/internal/config"
	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/service/chat"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	authConfig  config.AuthConfig
	authClient  auth.Client
	chatService chat.Service
}

func (s *serviceProvider) newAuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config", zap.Error(err))
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) GetAuthClient(_ context.Context) auth.Client {
	if s.authClient == nil {
		conn, err := grpc.Dial(s.newAuthConfig().GetPort(), grpc.WithDefaultCallOptions(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect %s: %s", s.authConfig.GetPort(), err)
		}
		closer.Add(conn.Close)

		client := user_v1.NewUserClient(conn)
		s.authClient = auth.NewClient(client)
	}

	return s.authClient
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetChatService(ctx context.Context) chat.Service {
	if s.chatService == nil {
		s.chatService = chat.NewService(s.GetAuthClient(ctx))
	}

	return s.chatService
}
