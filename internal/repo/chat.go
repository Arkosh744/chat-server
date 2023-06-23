package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/Arkosh744/chat-server/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository interface {
	CreateChat(ctx context.Context, usernames []string, saveHistory bool) (string, error)
	GetChat(_ context.Context, chatID string) (*models.Chat, error)
	ConnectToChat(ctx context.Context, chatID string, username string) (*models.Chat, error)
	AddUserToChat(_ context.Context, chatID string, username string) error
	SaveMessage(ctx context.Context, chatID string, message models.Message) error
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

var (
	ErrChatNotFound   = errors.New("chat not found")
	ErrUserNotAllowed = errors.New("user not allowed")
)

func (r *repository) CreateChat(_ context.Context, usernames []string, saveHistory bool) (string, error) {
	chatID := uuid.New().String()

	r.muChat.Lock()

	r.chats[chatID] = &models.Chat{
		ID:          chatID,
		SaveHistory: saveHistory,
		Usernames:   make(map[string]struct{}, len(usernames)),
		Streams:     make(map[string]chan<- models.Message, len(usernames)),
	}

	for _, username := range usernames {
		r.chats[chatID].Usernames[username] = struct{}{}
	}

	r.muChat.Unlock()

	return chatID, nil
}

func (r *repository) GetChat(_ context.Context, chatID string) (*models.Chat, error) {
	r.muChat.RLock()
	defer r.muChat.RUnlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return nil, ErrChatNotFound
	}

	return chat, nil
}

func (r *repository) ConnectToChat(_ context.Context, chatID string, username string) (*models.Chat, error) {
	r.muChat.RLock()

	chat, ok := r.chats[chatID]
	if !ok {
		r.muChat.RUnlock()
		return nil, ErrChatNotFound
	}

	r.muChat.RUnlock()

	chat.MuUsers.RLock()

	if _, ok := chat.Usernames[username]; !ok {
		chat.MuUsers.RUnlock()

		return nil, fmt.Errorf("user %s not allowed to be in chat %s", username, chatID)
	}
	chat.MuUsers.RUnlock()

	chat.MuStreams.Lock()
	defer chat.MuStreams.Unlock()

	if _, ok := chat.Streams[username]; ok {
		return nil, fmt.Errorf("user %s already connected to chat %s", username, chatID)
	}

	return chat, nil
}

func (r *repository) SaveMessage(_ context.Context, chatID string, message models.Message) error {
	r.muChat.RLock()
	defer r.muChat.RUnlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return ErrChatNotFound
	}

	if !chat.SaveHistory {
		return nil
	}

	chat.MuMessages.Lock()
	chat.Messages = append(chat.Messages, message)
	chat.MuMessages.Unlock()

	return nil
}

func (r *repository) AddUserToChat(_ context.Context, chatID string, username string) error {
	r.muChat.Lock()
	defer r.muChat.Unlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return ErrChatNotFound
	}

	chat.MuUsers.Lock()
	defer chat.MuUsers.Unlock()

	chat.Usernames[username] = struct{}{}

	return nil
}
