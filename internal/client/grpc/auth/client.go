package auth

import (
	"github.com/Arkosh744/auth-service-api/pkg/user_v1"
	"google.golang.org/grpc"
)

var _ Client = (*client)(nil)

type Client interface {
}

type client struct {
	userClient user_v1.UserV1Client
}

func NewClient(conn *grpc.ClientConn) *client {
	return &client{
		userClient: user_v1.NewUserV1Client(conn),
	}
}
