package models

import (
	"sync"
	"time"

	chatV1 "github.com/Arkosh744/chat-server/pkg/chat_v1"
)

type Chat struct {
	ID        string
	Usernames map[string]struct{}

	Messages []*Message
	Streams  map[string]Stream

	SaveHistory bool
	Mu          sync.RWMutex
}

type Stream interface {
	chatV1.ChatV1_ConnectToChatServer
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
