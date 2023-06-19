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
	Streams  map[string]chatV1.ChatV1_ConnectToChatServer

	SaveHistory bool
	Mu          sync.RWMutex
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
