package internal

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"
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

	for _, existing := range h.register {
		if existing == username {
			return false
		}
	}
	h.register[conn] = username
	return true
}

// unregister client from Register or connected client
func (h *Hub) Unregister(conn net.Conn) {
	h.mu.Lock()
	_, inConnections := h.connections[conn]
	_, inRegister := h.register[conn]

	if inConnections {
		delete(h.connections, conn)
		log.Printf("closed the Connected client connection successfully")
	}
	if inRegister {
		delete(h.register, conn)
		log.Printf("User removed from registered users")
	}

	h.mu.Unlock()

	if inConnections || inRegister {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

}

// Count

func (h *Hub) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.connections)
}

// Broadcast

func (h *Hub) Broadcast(message []byte, conn net.Conn, username string) {
	h.mu.Lock()
	targets := make([]net.Conn, 0, len(h.connections))
	for key := range h.connections {
		if key != conn {
			targets = append(targets, key)
		}
	}
	h.mu.Unlock()

	formatted := fmt.Sprintf(
		"[%s][%s]: %s\n",
		username,
		time.Now().Format("15:04:05"),
		string(message),
	)

	for _, key := range targets {
		if _, err := key.Write([]byte(formatted)); err != nil {
			log.Print("Error writing message to: ", err)
			h.Unregister(key)
		}
	}
}

func (h *Hub) Sendto(conn net.Conn, message []byte) {
	h.mu.Lock()
	_, ok := h.connections[conn]
	h.mu.Unlock()

	if !ok {
		log.Println("Log not found!")
		return
	}

	if _, err := conn.Write(message); err != nil {
		log.Printf("Writing Fail, closed the connection and deleted as well")
		h.Unregister(conn)
		return
	}
	log.Println("Message send Successful!")
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
