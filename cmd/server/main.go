package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()
	// reading buffer of 1024
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error Reading buffer:", err)
		return
	}
	// Write and responding
	data := buffer[:n]
	conversion := string(data)
	fmt.Println("Received:", conversion)

	response := "Message received, thank you!"
	a, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Error sending Message Receviced Acception", err)
		return
	}
	log.Printf("Sent %d bytes to client", a)

}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	} else {
		fmt.Printf("Server is listening on port 8080:")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error Acceping connection:", err)
			continue
		}

		log.Println("New connection accepted!", conn.RemoteAddr())
		go handleConnection(conn)

	}
}
