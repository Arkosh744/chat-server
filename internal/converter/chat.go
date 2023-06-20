package converter

import (
	"github.com/Arkosh744/chat-server/internal/models"
	chatV1 "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

func ToMessage(msg *chatV1.Message) *models.Message {
	return &models.Message{
		From:      msg.GetFrom(),
		Text:      msg.GetText(),
		Timestamp: msg.GetCreatedAt().AsTime(),
	}
}

func ToDescChatResponse(chat *models.Chat) *chatV1.GetChatResponse {
	usernames := make([]string, 0, len(chat.Usernames))

	for k := range chat.Usernames {
		usernames = append(usernames, k)
	}

	return &chatV1.GetChatResponse{
		Id:          chat.ID,
		Usernames:   usernames,
		SaveHistory: chat.SaveHistory,
	}
}
