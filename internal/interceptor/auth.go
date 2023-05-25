package interceptor

import (
	"context"
	"github.com/Arkosh744/chat-server/internal/client/grpc/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthInterceptor struct {
	authClient auth.Client
}

func NewAuthInterceptor(authClient auth.Client) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("failed to get metadata from incoming context")
		}

		if len(md["authorization"]) == 0 {
			return nil, errors.New("authorization header is not provided")
		}

		if !strings.Contains(md["authorization"][0], "Bearer ") {
			return nil, errors.New("invalid authorization header format")
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		if err = i.authClient.Check(ctx, info.FullMethod); err != nil {
			if strings.Contains(err.Error(), "access denied") {
				return nil, status.Errorf(codes.PermissionDenied, "access denied")
			}

			return nil, errors.Wrap(err, "failed to check access")
		}

		return handler(ctx, req)
	}
}
