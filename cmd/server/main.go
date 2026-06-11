package main

import (
	"bytes"
	"lan-chat/internal"
	"log"
	"net"
)

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

		// add to the connected users
		hub.Connected(conn)
		// check if it is in connected user
		if hub.HasConnection(conn) {
			welcome := "Connected to the Server successfully! \n\n Register your Username[A-Z,a-z,0-9]: "
			_, err := conn.Write([]byte(welcome))
			if err != nil {
				log.Println("Error Connecting to the server!", err)
				conn.Close()
				continue
			}
			username_buffer := make([]byte, 1024)
			n, err := conn.Read(username_buffer)
			// username := string(username_buffer[:n])
			username := string(bytes.TrimSpace(username_buffer[:n]))

			if err != nil {
				log.Printf("Error: Reading the username: %s", err)
			}
			if hub.ValidateUsername(username) == true {
				hub.Register(conn, username)
			}

		}
		if hub.IsRegister(conn) {
			conn.Write([]byte(
				"\nWelcome to LAN Chat!\n" +
					"Type your messages below.\n" +
					"-------------------------\n",
			))
			conn.Write([]byte("Message:"))

			registeredName := hub.GetUsername(conn)
			client := internal.NewClient(conn, hub, registeredName)
			go client.Read()
		}

	}

}
