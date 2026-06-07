package main

import (
	"io"
	"lan-chat/internal"
	"log"
	"net"
)

func handleConnection(conn net.Conn, hub *internal.Hub) {

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", conn.RemoteAddr())
			} else {
				log.Printf("[%s] Read error: %v", conn.RemoteAddr(), err)
			}
			hub.Unregister(conn)
			return
		}

		// Write and responding
		log.Printf("[%s] Received: %s", conn.RemoteAddr(), string(buffer[:n]))

		hub.Broadcast(buffer[:n])
		log.Printf("%d bytes were broadcasted from %s", n, conn.RemoteAddr())
	}

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
		hub.Register(conn)

		welcome := "Connected to the Server successfully! \n\n Write your Message: "
		n, err := conn.Write([]byte(welcome))
		if err != nil {
			log.Println("Error Connecting to the server!", err)
			conn.Close()
			continue
		}
		log.Printf("[%s] Welcome message sent (%d bytes)", conn.RemoteAddr(), n)

		log.Printf("New connection accepted: %s", conn.RemoteAddr())
		go handleConnection(conn, hub)

	}

}
