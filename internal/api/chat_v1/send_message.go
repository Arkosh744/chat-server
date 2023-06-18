package chat_v1

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/converter"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if err := i.chatService.SendMessage(ctx, req.GetChatId(), converter.ToMessage(req.GetMessage())); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
