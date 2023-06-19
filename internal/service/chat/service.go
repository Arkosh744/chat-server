package chat

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/internal/repo"
	chatV1 "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

var _ Service = (*service)(nil)

type Service interface {
	CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error)
	ConnectToChat(ctx context.Context, chatID string, username string, stream chatV1.ChatV1_ConnectToChatServer) error
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

	log.Infof("created chat: %s for users: %s", chatID, usernames)

	return chatID, nil
}

func (s *service) ConnectToChat(ctx context.Context, chatID string, username string, stream chatV1.ChatV1_ConnectToChatServer) error {
	if err := s.repo.ConnectToChat(ctx, chatID, username, stream); err != nil {
		log.Warnf("failed to connect user: %s to chat: %s", username, chatID)

		return err
	}

	log.Infof("connected user: %s to chat: %s", username, chatID)

	return nil
}

func (s *service) SendMessage(ctx context.Context, chatID string, message *models.Message) error {
	if err := s.repo.SendMessage(ctx, chatID, message); err != nil {
		return err
	}

	log.Infof("sent message: %s to chat: %s", message.Text, chatID)

	return nil
}

func (s *service) AddUserToChat(ctx context.Context, chatID string, username string) error {
	if err := s.repo.AddUserToChat(ctx, chatID, username); err != nil {
		return err
	}

	log.Infof("added user: %s to chat: %s", username, chatID)

	return nil
}
