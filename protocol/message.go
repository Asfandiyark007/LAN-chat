package protocol

import (
	"time"
)

const (
	SystemMessage  = "system"
	ChatMessage    = "chat"
	CommandMessage = "command"
	PrivateMessage = "private"
)

type WireMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

func NewSystemMessage(content string) WireMessage {
	return WireMessage{
		Type:      SystemMessage,
		Sender:    "Server",
		Timestamp: time.Now(),
		Content:   content,
	}
}

func NewChatMessage(sender, content string) WireMessage {
	return WireMessage{
		Type:      ChatMessage,
		Sender:    sender,
		Timestamp: time.Now(),
		Content:   content,
	}
}

func NewCommandMessage(command string) WireMessage {
	return WireMessage{
		Type:      CommandMessage,
		Timestamp: time.Now(),
		Content:   command,
	}
}

func NewPrivateMessage(sender, content string) WireMessage {
	return WireMessage{
		Type:      PrivateMessage,
		Sender:    sender,
		Timestamp: time.Now(),
		Content:   content,
	}
}
