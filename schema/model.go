package schema

import (
	"time"
)

type Notify struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Service   string    `json:"service"`
	CreatedAt time.Time `json:"created_at"`
}
