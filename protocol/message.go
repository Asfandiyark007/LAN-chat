package protocol

import (
	"time"
)

type WireMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}
