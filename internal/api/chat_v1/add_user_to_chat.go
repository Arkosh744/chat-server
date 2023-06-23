package chat_v1

import (
	"context"

	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) AddUserToChat(ctx context.Context, req *desc.AddUserToChatRequest) (*emptypb.Empty, error) {
	if err := i.chatService.AddUserToChat(ctx, req.GetChatId(), req.GetUsername()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
