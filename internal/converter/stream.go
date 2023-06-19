package converter

import (
	"github.com/Arkosh744/chat-server/internal/models"
	chatV1 "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

//type ChatStream struct {
//	GrpcStream chatV1.ChatV1_ConnectToChatServer
//}
//
//func (cs *ChatStream) Send(msg *models.Message) error {
//	return cs.GrpcStream.Send(&chatV1.Message{
//		From:      msg.From,
//		Text:      msg.Text,
//		CreatedAt: timestamppb.New(msg.Timestamp),
//	})
//}

func ToMessage(msg *chatV1.Message) *models.Message {
	return &models.Message{
		From:      msg.GetFrom(),
		Text:      msg.GetText(),
		Timestamp: msg.GetCreatedAt().AsTime(),
	}
}
