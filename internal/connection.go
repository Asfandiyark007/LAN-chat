package internal

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"time"

	"lan-chat/protocol"
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
		var wireMsg protocol.WireMessage
		if err := json.Unmarshal([]byte(line), &wireMsg); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}
		// not trusting whatever client is sending
		wireMsg.Sender = c.Username
		wireMsg.Timestamp = time.Now()

		c.Hub.Broadcast(wireMsg)
		log.Printf("[%s][%s][%s] Received: %s", wireMsg.Sender, wireMsg.Timestamp.Format("15:04:05"), c.Conn.RemoteAddr(), wireMsg.Content)

	}
}

func (c *Client) Write(msg protocol.WireMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error: Mashaling the message. %v", err)
		c.Hub.Unregister(c.Conn)
		return
	}

	data = append(data, '\n')

	_, err = c.Conn.Write(data)
	if err != nil {
		log.Printf("Error writing message: %v", err)
		c.Hub.Unregister(c.Conn)
	}

}
