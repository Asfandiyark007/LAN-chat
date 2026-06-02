package main

import (
	"fmt"
	"log"
	"net"
)

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
		//defer conn.Close() // cant use (defer) inside for loop

		log.Println("New connection accepted!", conn.RemoteAddr())
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error Reading buffer:", err)
			conn.Close()
			continue
		}

		data := buffer[:n]
		conversion := string(data)
		fmt.Println("Received:", conversion)

		response := "Message received, thank you!"
		a, err := conn.Write([]byte(response))
		if err != nil {
			log.Println("Error sending Message Receviced Acception", err)
			conn.Close()
			continue
		}
		log.Printf("Sent %d bytes to client", a)
		conn.Close()
	}

}
