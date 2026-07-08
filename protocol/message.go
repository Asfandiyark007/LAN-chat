package protocol

import (
	"time"
)

const (
	SystemMessage = "system"
	ChatMessage   = "chat"
)

type WireMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

func NewSystemMessage(content string) WireMessage {
	return WireMessage{
		Type:      "system",
		Sender:    "Server",
		Timestamp: time.Now(),
		Content:   content,
	}
}

func NewChatMessage(sender, content string) WireMessage {
	return WireMessage{
		Type:      "chat",
		Sender:    sender,
		Timestamp: time.Now(),
		Content:   content,
	}
}
