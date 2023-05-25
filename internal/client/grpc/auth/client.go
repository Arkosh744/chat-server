package auth

import (
	"context"
	accessV1 "github.com/Arkosh744/auth-service-api/pkg/access_v1"
	userV1 "github.com/Arkosh744/auth-service-api/pkg/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ Client = (*client)(nil)

type Client interface {
	List(ctx context.Context) ([]*userV1.UserInfo, error)
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

func (c *client) List(ctx context.Context) ([]*userV1.UserInfo, error) {
	res, err := c.userClient.List(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	return res.GetUsers(), nil
}

func (c *client) Check(ctx context.Context, endpoint string) error {
	if _, err := c.accessClient.CheckAccess(ctx, &accessV1.CheckAccessRequest{Endpoint: endpoint}); err != nil {
		return err
	}

	return nil
}
