package internal

import (
	"net"
	"time"
)

type Message struct {
	From      net.Conn
	Content   []byte
	Timestamp time.Time
}

func NewMessage(conn net.Conn, content []byte) *Message {
	return &Message{
		From:      conn,
		Content:   content,
		Timestamp: time.Now(),
	}
}
