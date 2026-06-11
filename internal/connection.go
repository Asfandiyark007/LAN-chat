package internal

import (
	"log"
	"net"
)

type Client struct {
	Conn     net.Conn
	Hub      *Hub
	Username string
}

func NewClient(conn net.Conn, Hub *Hub, username string) *Client {
	return &Client{
		Conn:     conn,
		Hub:      Hub,
		Username: username,
	}
}

func (c *Client) Read() {
	for {
		buffer := make([]byte, 1024)
		n, err := c.Conn.Read(buffer)

		if err != nil {
			c.Hub.Unregister(c.Conn)
			log.Printf("Could not read: using unregister()")
			return
		}
		message := NewMessage(c.Conn, buffer[:n])
		c.Conn.Write([]byte("\n"))
		c.Hub.Broadcast(message.Content, c.Conn, c.Username)
		c.Conn.Write([]byte("Message:"))
		log.Printf("[%s][%s][%s] Received: %s", c.Username, message.Timestamp.Format("15:04:05"), c.Conn.RemoteAddr(), message.Content)

	}
}

func (c *Client) Write(message []byte) {
	_, err := c.Conn.Write(message)
	if err != nil {
		log.Printf("Error: writing the message. %s", err)
		c.Hub.Unregister(c.Conn)
		return
	}

}
