package chat_v1

import (
	"context"

	"github.com/Arkosh744/chat-server/internal/converter"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	chat, err := i.chatService.GetChat(ctx, req.GetChatId())
	if err != nil {
		return nil, err
	}

	for _, messages := range chat.Streams {
		messages <- converter.ToMessage(req.GetMessage())
	}

	if err = i.chatService.SaveMessage(ctx, req.GetChatId(), converter.ToMessage(req.GetMessage())); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
