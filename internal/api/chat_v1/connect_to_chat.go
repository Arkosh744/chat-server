package chat_v1

import (
	"github.com/Arkosh744/chat-server/internal/converter"
	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/Arkosh744/chat-server/internal/repo"
	desc "github.com/Arkosh744/chat-server/pkg/chat_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ConnectToChat(req *desc.ConnectChatRequest, grpcStream desc.ChatV1_ConnectToChatServer) error {
	var (
		ctx          = grpcStream.Context()
		username     = req.GetUsername()
		messagesChan = make(chan models.Message, models.ChanMessageCapacity)
	)

	chat, err := i.chatService.ConnectToChat(ctx, req.GetChatId(), req.GetUsername())
	if err != nil {
		if errors.Is(err, repo.ErrChatNotFound) {
			return status.Errorf(codes.NotFound, "chat not found")
		}

		if errors.Is(err, repo.ErrUserNotAllowed) {
			return status.Errorf(codes.PermissionDenied, "user not allowed")
		}

		return status.Errorf(codes.Internal, "error connecting to chat: %v", err)
	}

	chat.MuStreams.Lock()
	chat.Streams[username] = messagesChan
	chat.MuStreams.Unlock()

	chat.MuMessages.RLock()
	if chat.SaveHistory {
		log.Infof("%s", chat.Messages)
		for _, message := range chat.Messages {
			messagesChan <- message
		}
	}
	chat.MuMessages.RUnlock()

	for {
		select {
		case <-ctx.Done():
			close(messagesChan)

			chat.MuStreams.Lock()
			delete(chat.Streams, username)
			chat.MuStreams.Unlock()

			return nil

		case msg, ok := <-messagesChan:
			if !ok {
				return nil
			}

			if err = grpcStream.Send(converter.ToDescMessage(msg)); err != nil {
				log.Errorf("Error sending message: %v", err)
			}
		}
	}
}
