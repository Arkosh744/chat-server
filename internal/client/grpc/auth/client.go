package auth

import (
	"context"

	accessV1 "github.com/Arkosh744/auth-service-api/pkg/access_v1"
	userV1 "github.com/Arkosh744/auth-service-api/pkg/user_v1"
)

var _ Client = (*client)(nil)

type Client interface {
	Check(ctx context.Context, endpoint string) error
}

type client struct {
	userClient   userV1.UserClient
	accessClient accessV1.AccessV1Client
}

func NewClient(c userV1.UserClient, a accessV1.AccessV1Client) *client {
	return &client{
		userClient:   c,
		accessClient: a,
	}
}

func (c *client) Check(ctx context.Context, endpoint string) error {
	if _, err := c.accessClient.CheckAccess(ctx, &accessV1.CheckAccessRequest{Endpoint: endpoint}); err != nil {
		return err
	}

	return nil
}
