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
	return r.getChat(chatID)
}

func (r *repository) ConnectToChat(_ context.Context, chatID string, username string) (*models.Chat, error) {
	chat, err := r.getChat(chatID)
	if err != nil {
		return nil, err
	}

	if err = checkUserExists(chat, username, chatID); err != nil {
		return nil, err
	}

	chat.MuStreams.Lock()
	defer chat.MuStreams.Unlock()

	if _, ok := chat.Streams[username]; ok {
		return nil, fmt.Errorf("user %s already connected to chat %s", username, chatID)
	}

	return chat, nil
}

func (r *repository) SaveMessage(_ context.Context, chatID string, message models.Message) error {
	chat, err := r.getChat(chatID)
	if err != nil {
		return err
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
	chat, err := r.getChat(chatID)
	if err != nil {
		return err
	}

	chat.MuUsers.Lock()
	defer chat.MuUsers.Unlock()

	chat.Usernames[username] = struct{}{}

	return nil
}

func checkUserExists(chat *models.Chat, username string, chatID string) error {
	chat.MuUsers.RLock()

	if _, ok := chat.Usernames[username]; !ok {
		chat.MuUsers.RUnlock()

		return fmt.Errorf("user %s not allowed to be in chat %s", username, chatID)
	}
	chat.MuUsers.RUnlock()
	return nil
}

func (r *repository) getChat(chatID string) (*models.Chat, error) {
	r.muChat.RLock()
	defer r.muChat.RUnlock()

	chat, ok := r.chats[chatID]
	if !ok {
		return nil, ErrChatNotFound
	}

	return chat, nil
}
