package chat_v1

import (
	"log"

	"github.com/Arkosh744/chat-server/internal/converter"
	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/internal/repo"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ConnectToChat(req *desc.ConnectChatRequest, grpcStream desc.ChatV1_ConnectToChatServer) error {
	ctx := grpcStream.Context()

	messages := make(chan models.Message)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			case msg := <-messages:
				if err := grpcStream.Send(converter.ToDesckMessage(msg)); err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
		}
	}()

	if err := i.chatService.ConnectToChat(ctx, req.GetChatId(), req.GetUsername(), messages); err != nil {
		if errors.Is(err, repo.ErrChatNotFound) {
			return status.Errorf(codes.NotFound, "chat not found")
		}

		if errors.Is(err, repo.ErrUserNotAllowed) {
			return status.Errorf(codes.PermissionDenied, "user not allowed")
		}

		return status.Errorf(codes.Internal, "error connecting to chat: %v", err)
	}

	return nil
}
