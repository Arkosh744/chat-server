package auth

import (
	"context"
	userV1 "github.com/Arkosh744/auth-service-api/pkg/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

var _ Client = (*client)(nil)

type Client interface {
	List(ctx context.Context) ([]*userV1.UserInfo, error)
}

type client struct {
	userClient userV1.UserClient
}

func NewClient(c userV1.UserClient) *client {
	return &client{
		userClient: c,
	}
}

func (c *client) List(ctx context.Context) ([]*userV1.UserInfo, error) {
	res, err := c.userClient.List(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	return res.GetUsers(), nil
}
