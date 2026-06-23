package internal

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Client struct {
	Conn     net.Conn
	Hub      *Hub
	Username string
	Reader   *bufio.Reader
}

func NewClient(conn net.Conn, Hub *Hub, username string, reader *bufio.Reader) *Client {
	return &Client{
		Conn:     conn,
		Hub:      Hub,
		Username: username,
		Reader:   reader,
	}
}

func (c *Client) Read() {
	for {

		line, err := c.Reader.ReadString('\n')
		if err != nil {
			c.Hub.Unregister(c.Conn)
			log.Printf("Could not read: using unregister()")
			return
		}

		line = strings.TrimSpace(line)

		message := NewMessage(c.Conn, []byte(line))
		c.Hub.Broadcast(message.Content, c.Conn, c.Username)
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
