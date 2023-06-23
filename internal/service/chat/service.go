package chat

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/internal/repo"
)

var _ Service = (*service)(nil)

type Service interface {
	CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error)
	GetChat(ctx context.Context, chatID string) (*models.Chat, error)
	ConnectToChat(ctx context.Context, chatID string, username string) (*models.Chat, error)
	AddUserToChat(ctx context.Context, chatID string, username string) error
	SaveMessage(ctx context.Context, chatID string, message models.Message) error
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

func (s *service) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
	chat, err := s.repo.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	log.Infof("got chat: %s", chatID)

	return chat, nil
}

func (s *service) ConnectToChat(ctx context.Context, chatID string, username string) (*models.Chat, error) {
	chat, err := s.repo.ConnectToChat(ctx, chatID, username)
	if err != nil {
		log.Warnf("failed to get chat to connect user: %s to chat: %s", username, chatID)

		return nil, err
	}

	return chat, nil
}

func (s *service) SaveMessage(ctx context.Context, chatID string, message models.Message) error {
	if err := s.repo.SaveMessage(ctx, chatID, message); err != nil {
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
