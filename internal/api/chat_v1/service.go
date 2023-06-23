package chat_v1

import (
	"github.com/Arkosh744/chat-server/internal/service/chat"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

var _ desc.ChatV1Server = (*Implementation)(nil)

type Implementation struct {
	desc.UnimplementedChatV1Server

	chatService chat.Service
}

func NewImplementation(chatService chat.Service) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
