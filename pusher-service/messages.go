package main

import (
	"time"
)

const (
	KindNotifyCreated = iota + 1
)

type NotifyCreatedMessage struct {
	Kind      uint32    `json:"kind"`
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Service   string    `json:"service"`
	CreatedAt time.Time `json:"created_at"`
}

func newNotifyCreatedMessage(id string, title string, body string, service string, createdAt time.Time) *NotifyCreatedMessage {
	return &NotifyCreatedMessage{
		Kind:      KindNotifyCreated,
		ID:        id,
		Title:     title,
		Body:      body,
		Service:   service,
		CreatedAt: createdAt,
	}
}
