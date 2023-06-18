package interceptor

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authClient auth.Client
}

func NewAuthInterceptor(authClient auth.Client) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient}
}

var errAccessDenied = status.Error(codes.PermissionDenied, "access denied")

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		if err = i.authClient.Check(ctx, info.FullMethod); err != nil {
			if errors.Is(err, errAccessDenied) {
				return nil, errAccessDenied
			}

			return nil, errors.Wrap(err, "failed to check access")
		}

		return handler(ctx, req)
	}
}
