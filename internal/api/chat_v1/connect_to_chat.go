package chat_v1

import (
	"github.com/Arkosh744/chat-server/internal/converter"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

func (i *Implementation) ConnectToChat(req *desc.ConnectChatRequest, grpcStream desc.ChatV1_ConnectToChatServer) error {
	ctx := grpcStream.Context()
	stream := &converter.ChatStream{GrpcStream: grpcStream}

	return i.chatService.ConnectToChat(ctx, req.GetChatId(), req.GetUsername(), stream)
}
