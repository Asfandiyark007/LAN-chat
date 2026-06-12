package main

import (
	"bytes"
	"lan-chat/internal"
	"log"
	"net"
)

func handleConnection(conn net.Conn, hub *internal.Hub) {
	hub.Connected(conn)

	welcome := "Connected to the Server successfully!\n\nRegister your Username[A-Z,a-z,0-9]: "

	_, err := conn.Write([]byte(welcome))
	if err != nil {
		conn.Close()
		return
	}

	usernameBuffer := make([]byte, 1024)

	n, err := conn.Read(usernameBuffer)
	if err != nil {
		conn.Close()
		return
	}

	username := string(bytes.TrimSpace(usernameBuffer[:n]))

	if !hub.ValidateUsername(username) {
		conn.Write([]byte("Invalid username\n"))
		conn.Close()
		return
	}

	hub.Register(conn, username)

	_, err = conn.Write([]byte(
		"\nWelcome to LAN Chat!\n" +
			"Type your messages below.\n" +
			"-------------------------\n",
	))
	if err != nil {
		hub.Unregister(conn)
		return
	}

	client := internal.NewClient(conn, hub, username)

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
