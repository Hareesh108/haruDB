package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

	conn.Write([]byte("Welcome to HaruDB v0.1 🎉\n"))
	conn.Write([]byte("Type 'exit' to quit.\n\n"))

	scanner := bufio.NewScanner(conn)
	for {
		conn.Write([]byte("haruDB> "))

		if !scanner.Scan() {
			break // client disconnected
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			conn.Write([]byte("Goodbye 👋\n"))
			break
		}

		// For now, just echo back
		conn.Write([]byte("You said: " + input + "\n"))
	}
}
