// repo.go

package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Repository interface {
	CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error)
	ConnectToChat(ctx context.Context, chatID string, username string, stream models.Stream) error
	AddUserToChat(_ context.Context, chatID string, username string) error
	SendMessage(ctx context.Context, chatID string, message *models.Message) error
}

type repository struct {
	chats  map[string]*models.Chat
	muChat sync.RWMutex
}

func NewRepo() Repository {
	return &repository{
		chats: make(map[string]*models.Chat),
	}
}

func (r *repository) CreateChat(_ context.Context, usernames []string, saveHistory bool) (string, error) {
	chatID := uuid.New().String()

	r.muChat.Lock()

	r.chats[chatID] = &models.Chat{
		ID:          chatID,
		SaveHistory: saveHistory,
		Usernames:   make(map[string]struct{}, len(usernames)),
		Streams:     make(map[string]models.Stream, len(usernames)),
	}

	for _, username := range usernames {
		r.chats[chatID].Usernames[username] = struct{}{}
	}

	r.muChat.Unlock()

	return chatID, nil
}

func (r *repository) ConnectToChat(ctx context.Context, chatID string, username string, stream models.Stream) error {
	r.muChat.Lock()

	chat, ok := r.chats[chatID]
	if !ok {
		r.muChat.Unlock()
		return fmt.Errorf("chat %s not found", chatID)
	}

	if _, ok := chat.Usernames[username]; !ok {
		r.muChat.Unlock()
		return fmt.Errorf("user %s not allowed to be in chat %s", username, chatID)
	}

	uid := uuid.New().String()
	chat.Streams[uid] = stream

	if chat.SaveHistory {
		for _, message := range chat.Messages {
			if err := stream.Send(message); err != nil {
				r.muChat.Unlock()
				return err
			}
		}
	}
	r.muChat.Unlock()

	go func() {
		<-ctx.Done()
		chat.Mu.Lock()
		delete(chat.Streams, uid)
		chat.Mu.Unlock()
	}()

	return nil
}

func (r *repository) SendMessage(_ context.Context, chatID string, message *models.Message) error {
	r.muChat.Lock()

	chat, ok := r.chats[chatID]
	if !ok {
		r.muChat.Unlock()
		return fmt.Errorf("chat does not exist")
	}

	if chat.SaveHistory {
		chat.Messages = append(chat.Messages, message)
	}

	var resErr *multierror.Error
	for _, stream := range chat.Streams {
		if err := stream.Send(message); err != nil {
			resErr = multierror.Append(resErr, err)
		}

	}
	if resErr.ErrorOrNil() != nil {
		return errors.Wrap(resErr.ErrorOrNil(), "error sending message to streams")
	}

	r.muChat.Unlock()

	return nil
}

func (r *repository) AddUserToChat(_ context.Context, chatID string, username string) error {
	r.muChat.Lock()
	defer r.muChat.Unlock()

	chat, ok := r.chats[chatID]
	if !ok {
		r.muChat.Unlock()
		return fmt.Errorf("chat %s not found", chatID)
	}

	chat.Usernames[username] = struct{}{}

	return nil
}
