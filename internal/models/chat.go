package models

import (
	"sync"
	"time"
)

type Chat struct {
	ID          string
	Usernames   map[string]struct{}
	SaveHistory bool
	Messages    []*Message
	Streams     map[string]Stream

	Mu sync.RWMutex
}

type Stream interface {
	Send(*Message) error
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
