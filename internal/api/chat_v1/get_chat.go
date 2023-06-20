package chat_v1

import (
	"context"
	"errors"

	"github.com/Arkosh744/chat-server/internal/converter"
	"github.com/Arkosh744/chat-server/internal/repo"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetChat(ctx context.Context, req *desc.GetChatRequest) (*desc.GetChatResponse, error) {
	chat, err := i.chatService.GetChat(ctx, req.GetChatId())
	if err != nil {
		if errors.Is(err, repo.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "error getting chat: %v", err)
	}

	return converter.ToDescChatResponse(chat), nil
}
