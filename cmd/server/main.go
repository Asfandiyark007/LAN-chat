package main

import (
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
		client := internal.NewClient(conn, hub, "anonymous")

		go client.Read()

	}

}
