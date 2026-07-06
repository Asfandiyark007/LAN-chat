package main

import (
	"bufio"
	"lan-chat/internal"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, hub *internal.Hub) {
	hub.Connected(conn)

	_, err := conn.Write([]byte("Connected to the Server successfully!\n\n"))
	if err != nil {
		conn.Close()
		return
	}

	reader := bufio.NewReader(conn)

	var username string

	for {
		_, err := conn.Write([]byte("Register your Username[A-Z,a-z,0-9]:\n"))
		if err != nil {
			hub.Unregister(conn)
			return
		}

		raw, err := reader.ReadString('\n')
		if err != nil {
			hub.Unregister(conn)
			return
		}
		raw = strings.TrimSpace(raw)

		if !hub.ValidateUsername(raw) {
			_, err := conn.Write([]byte("Invalid username. Use only letters and numbers, 1-9 characters long.\n"))
			if err != nil {
				hub.Unregister(conn)
				return
			}
			continue
		}

		if !hub.Register(conn, raw) {

			_, err = conn.Write([]byte("Username already taken. Try another.\n"))
			if err != nil {
				hub.Unregister(conn)
				return
			}
			continue
		}
		username = raw
		break
	}

	_, err = conn.Write([]byte("REGISTERED_OK\n"))
	if err != nil {
		hub.Unregister(conn)
		return
	}

	_, err = conn.Write([]byte(
		"Welcome to LAN Chat!\n" +
			"Type your messages below.\n" +
			"-------------------------\n",
	))
	if err != nil {
		hub.Unregister(conn)
		return
	}

	client := internal.NewClient(conn, hub, username, reader)

	client.Read()
}

func main() {

	hub := internal.NewHub()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	} else {
		log.Println("Server is listening on port 8080:")
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, hub)
	}

}
