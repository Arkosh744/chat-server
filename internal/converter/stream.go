package converter

import (
	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatStream struct {
	GrpcStream chat_v1.ChatV1_ConnectChatServer
}

func (cs *ChatStream) Send(msg *models.Message) error {
	return cs.GrpcStream.Send(&chat_v1.Message{
		From:      msg.From,
		Text:      msg.Text,
		Timestamp: timestamppb.New(msg.Timestamp),
	})
}

func ToMessage(msg *chat_v1.Message) *models.Message {
	return &models.Message{
		From:      msg.GetFrom(),
		Text:      msg.GetText(),
		Timestamp: msg.GetTimestamp().AsTime(),
	}
}
