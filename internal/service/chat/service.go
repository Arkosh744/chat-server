package chat

import (
	"context"
	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/model"
	"github.com/google/uuid"
)

var _ Service = (*service)(nil)

type Service interface {
	ListUsers(ctx context.Context) ([]*model.User, error)
	CreateChat(ctx context.Context, usernames []string) (uuid.UUID, error)
}

type service struct {
	auth auth.Client
}

func NewService(a auth.Client) *service {
	return &service{auth: a}
}

func (s *service) ListUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.auth.List(ctx)
	if err != nil {
		log.Errorf("failed to list users: %v", err)
		return nil, err
	}

	var res []*model.User
	for _, u := range users {
		res = append(res, &model.User{
			UserIdentifier: model.UserIdentifier{
				Username: u.GetUsername(),
				Email:    u.GetEmail()},
			Role: model.ToRole(u.Role),
		})
	}

	return res, nil
}

func (s *service) CreateChat(_ context.Context, _ []string) (uuid.UUID, error) {
	return uuid.New(), nil
}
