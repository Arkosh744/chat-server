package chat

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/internal/repo"
)

var _ Service = (*service)(nil)

type Service interface {
	CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error)
	ConnectToChat(ctx context.Context, chatID string, username string, stream models.Stream) error
	AddUserToChat(_ context.Context, chatID string, username string) error
	SendMessage(ctx context.Context, chatID string, message *models.Message) error
}

type service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) *service {
	return &service{
		repo: r}
}

func (s *service) CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error) {
	chatID, err := s.repo.CreateChat(ctx, usernames, saveHistory)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (s *service) ConnectToChat(ctx context.Context, chatID string, username string, stream models.Stream) error {
	return s.repo.ConnectToChat(ctx, chatID, username, stream)
}

func (s *service) SendMessage(ctx context.Context, chatID string, message *models.Message) error {
	return s.repo.SendMessage(ctx, chatID, message)
}

func (s *service) AddUserToChat(ctx context.Context, chatID string, username string) error {
	return s.repo.AddUserToChat(ctx, chatID, username)
}
