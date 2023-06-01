package chat_v1

import (
	"context"

	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

func (i *Implementation) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	chatID, err := i.chatService.CreateChat(ctx, req.GetUsernames())
	if err != nil {
		return nil, err
	}

	return &desc.CreateChatResponse{ChatId: chatID.String()}, nil
}
