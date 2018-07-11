package event

import (
	"time"
)

type Message interface {
	Key() string
}

type NotifyCreatedMessage struct {
	ID        string
	Title     string
	Body      string
	Service   string
	CreatedAt time.Time
}

func (m *NotifyCreatedMessage) Key() string {
	return "notify.created"
}
