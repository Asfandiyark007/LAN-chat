package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"lan-chat/protocol"
)

type Hub struct {
	connections map[net.Conn]struct{}
	mu          sync.Mutex
	register    map[net.Conn]string
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[net.Conn]struct{}),
		register:    make(map[net.Conn]string),
	}
}

// returns if user is connected
func (h *Hub) HasConnection(conn net.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.connections[conn]
	return ok
}

// return username
func (h *Hub) GetUsername(conn net.Conn) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	username, ok := h.register[conn]
	if !ok {
		return ""
	}

	return username
}

// // // return the connected users with /who command
func (h *Hub) Who() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	users := make([]string, 0, len(h.register))
	for _, username := range h.register {
		users = append(users, username)
	}
	return users
}

// return if user is registered
func (h *Hub) IsRegister(conn net.Conn) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.register[conn]
	return ok
}

// connected client
func (h *Hub) Connected(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[conn] = struct{}{}
}

// Register client
func (h *Hub) Register(conn net.Conn, username string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	target := strings.ToLower(username)
	for _, existing := range h.register {
		if strings.ToLower(existing) == target {
			return false
		}
	}
	h.register[conn] = username
	return true
}

// unregister client from Register or connected client
func (h *Hub) Unregister(conn net.Conn) {
	h.mu.Lock()

	username := h.register[conn]
	_, inConnections := h.connections[conn]
	_, inRegister := h.register[conn]

	if inConnections {
		delete(h.connections, conn)
		log.Printf("closed the Connected client connection successfully")
	}
	if inRegister {
		delete(h.register, conn)
		log.Printf("User [%s] removed from registered users", username)
	}

	h.mu.Unlock()

	if !(inConnections || inRegister) {
		return
	}

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}

	msg := protocol.NewSystemMessage(
		fmt.Sprintf("User [%s] left the chat.", username),
	)

	h.BroadcastExcept(msg, conn)

}

// Count

func (h *Hub) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.connections)
}

// Announcements (Sends to everyone other Except the sender)
func (h *Hub) BroadcastExcept(msg protocol.WireMessage, except net.Conn) {
	h.mu.Lock()
	targets := make([]net.Conn, 0, len(h.connections))
	for conn := range h.connections {
		if conn != except {
			targets = append(targets, conn)
		}
	}
	h.mu.Unlock()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error Marshaling message: %v", err)
		return
	}
	data = append(data, '\n')

	for _, conn := range targets {
		if _, err := conn.Write(data); err != nil {
			log.Printf("Error writing message: %v", err)
			h.Unregister(conn)
		}
	}
}

// Message to a single connection
func (h *Hub) Send(conn net.Conn, msg protocol.WireMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = conn.Write(data)
	if err != nil {
		h.Unregister(conn)
		return err
	}

	return nil
}

// Find user by username (Inefficent- use direct map lookup but want to test)
func (h *Hub) GetConnectionByUsername(username string) (net.Conn, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	target := strings.ToLower(username)
	for conn, registeredUsername := range h.register {
		if strings.ToLower(registeredUsername) == target {
			return conn, true
		}
	}
	return nil, false
}

// Broadcast (Sends it to everyone)

func (h *Hub) Broadcast(msg protocol.WireMessage) {
	h.mu.Lock()
	targets := make([]net.Conn, 0, len(h.connections))
	for key := range h.connections {
		targets = append(targets, key)
	}
	h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	data = append(data, '\n')

	for _, key := range targets {
		if _, err := key.Write(data); err != nil {
			log.Printf("Error writing message to: %v ", err)
			h.Unregister(key)
		}
	}
}

// username validation logic
func (h *Hub) ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	if !matched {
		return false
	}
	numChars := utf8.RuneCountInString(username)
	if numChars < 1 || numChars >= 10 {
		return false
	}

	return true
}
