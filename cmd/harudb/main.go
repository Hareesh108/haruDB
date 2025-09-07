package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := "54321"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("🚀 HaruDB server started on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Fprintf(conn, "Welcome to HaruDB v0.1 🎉\n")
}
