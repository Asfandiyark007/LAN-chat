package internal

import (
	"log"
	"net"
	"sync"
)

type Hub struct {
	connections map[net.Conn]bool
	mu          sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[net.Conn]bool),
	}
}

// Register client
func (h *Hub) Register(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[conn] = true
}

// unregister client
func (h *Hub) Unregister(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.connections[conn]
	if ok {
		delete(h.connections, conn)
		conn.Close()
	}
}

// Count

func (h *Hub) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.connections)
}

// Broadcast

func (h *Hub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for key := range h.connections {
		_, err := key.Write(message)

		if err != nil {
			log.Print("Error writing message to: ", err)
			key.Close()
			delete(h.connections, key)
		}

	}
}

// Send To
func (h *Hub) Sendto(conn net.Conn, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.connections[conn]
	if ok {
		_, err := conn.Write([]byte(message))
		if err != nil {
			conn.Close()
			delete(h.connections, conn)
			log.Printf("Writing Fail, closed the connection and deleted as well")
		}
		log.Println("Message send Successful!")

	} else {
		log.Println("Log not found!")
	}
}
