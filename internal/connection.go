package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"lan-chat/protocol"
)

type Client struct {
	Conn       net.Conn
	Hub        *Hub
	Username   string
	Reader     *bufio.Reader
	tokens     float64
	lastRefill time.Time
}

func NewClient(conn net.Conn, Hub *Hub, username string, reader *bufio.Reader) *Client {
	return &Client{
		Conn:       conn,
		Hub:        Hub,
		Username:   username,
		Reader:     reader,
		tokens:     3,
		lastRefill: time.Now(),
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

		switch wireMsg.Type {
		case protocol.ChatMessage:
			wireMsg.Sender = c.Username
			wireMsg.Timestamp = time.Now()

			if allowed, wait := c.Allow(); allowed {
				c.Hub.Broadcast(wireMsg)
			} else {
				reply := protocol.NewSystemMessage(
					fmt.Sprintf("Slow down! Try again in %.1fs.", wait.Seconds()),
				)
				c.Hub.Send(c.Conn, reply)
			}

			log.Printf("[%s][%s][%s] Received: %s",
				wireMsg.Sender,
				wireMsg.Timestamp.Format("15:04:05"),
				c.Conn.RemoteAddr(),
				wireMsg.Content,
			)

		case protocol.CommandMessage:
			c.handleCommand(wireMsg)
		default:
			log.Printf("Unknown message type: %s", wireMsg.Type)
		}

	}
}

func (c *Client) handleCommand(msg protocol.WireMessage) {
	parts := strings.SplitN(msg.Content, " ", 3)
	if len(parts) == 0 {
		return
	}
	command := parts[0]
	switch command {

	case "who":
		users := c.Hub.Who()

		reply := protocol.NewSystemMessage(
			"Connected users: " + strings.Join(users, ", "),
		)

		c.Hub.Send(c.Conn, reply)
	case "msg":
		if len(parts) < 3 {
			reply := protocol.NewSystemMessage("Usage: /msg @username message")
			c.Hub.Send(c.Conn, reply)
			return
		}

		target := strings.TrimPrefix(parts[1], "@")

		targetConn, ok := c.Hub.GetConnectionByUsername(target)
		if !ok {
			reply := protocol.NewSystemMessage("User not found.")
			c.Hub.Send(c.Conn, reply)
			return
		}
		// send to the recipient
		private := protocol.NewPrivateMessage(c.Username, target, parts[2])
		c.Hub.Send(targetConn, private)

		// echo back to the sends so they can still see their own DM they sent
		c.Hub.Send(c.Conn, private)

	default:
		reply := protocol.NewSystemMessage("Unknown command")
		c.Hub.Send(c.Conn, reply)
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

// rate limit and token cal
func (c *Client) Allow() (bool, time.Duration) {
	now := time.Now()
	elapsed := now.Sub(c.lastRefill)

	earned := elapsed.Seconds() / 2.0
	c.tokens += earned

	if c.tokens > 3 {
		c.tokens = 3
	}

	c.lastRefill = now

	if c.tokens >= 1 {
		c.tokens -= 1
		return true, 0
	}

	waitSeconds := (1 - c.tokens) * 2.0
	return false, time.Duration(waitSeconds * float64(time.Second))
}
