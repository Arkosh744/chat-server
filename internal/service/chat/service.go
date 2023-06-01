package chat

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/google/uuid"
)

var _ Service = (*service)(nil)

type Service interface {
	CreateChat(ctx context.Context, usernames []string) (uuid.UUID, error)
}

type service struct {
	auth auth.Client
}

func NewService(a auth.Client) *service {
	return &service{auth: a}
}

func (s *service) CreateChat(_ context.Context, _ []string) (uuid.UUID, error) {
	return uuid.New(), nil
}
