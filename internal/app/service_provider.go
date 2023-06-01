package app

import (
	"context"

	accessV1 "github.com/Arkosh744/auth-service-api/pkg/access_v1"
	userV1 "github.com/Arkosh744/auth-service-api/pkg/user_v1"
	chatV1 "github.com/Arkosh744/chat-server/internal/api/chat_v1"
	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/Arkosh744/chat-server/internal/closer"
	"github.com/Arkosh744/chat-server/internal/config"
	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/service/chat"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceProvider struct {
	authConfig config.AuthConfig
	grpcConfig config.GRPCConfig

	authClient  auth.Client
	chatService chat.Service

	chatImpl *chatV1.Implementation
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

func (s *serviceProvider) GetAuthClient(ctx context.Context) auth.Client {
	if s.authClient == nil {
		creds, err := credentials.NewClientTLSFromFile("./certs/ca.crt", "localhost")
		if err != nil {
			log.Fatalf("failed to load credentials: %s", err)
		}

		conn, err := grpc.DialContext(
			ctx,
			s.newAuthConfig().GetHost(),
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			log.Fatalf("failed to connect %s: %s", s.authConfig.GetHost(), err)
		}
		closer.Add(conn.Close)

		userClient := userV1.NewUserClient(conn)
		accessClient := accessV1.NewAccessV1Client(conn)
		s.authClient = auth.NewClient(userClient, accessClient)
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

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config", zap.Error(err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetChatImpl(ctx context.Context) *chatV1.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chatV1.NewImplementation(s.GetChatService(ctx))
	}

	return s.chatImpl
}
