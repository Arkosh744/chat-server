package models

import (
	"sync"
	"time"
)

const ChanMessageCapacity = 100

type Chat struct {
	ID        string
	Usernames map[string]struct{}

	Messages []Message
	Streams  map[string]chan<- Message

	SaveHistory bool

	MuUsers    sync.RWMutex
	MuMessages sync.RWMutex
	MuStreams  sync.Mutex
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
