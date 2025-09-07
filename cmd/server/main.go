package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Hareesh108/haruDB/internal/parser"
)

func main() {
	port := "54321"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("🚀 HaruDB server started on port %s\n", port)

	engine := parser.NewEngine()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, engine)
	}
}

func handleConnection(conn net.Conn, engine *parser.Engine) {
	defer conn.Close()

	conn.Write([]byte("Welcome to HaruDB v0.1 🎉\n"))
	conn.Write([]byte("Type 'exit' to quit.\n\n"))

	scanner := bufio.NewScanner(conn)
	for {
		conn.Write([]byte("haruDB> "))

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			conn.Write([]byte("Goodbye 👋\n"))
			break
		}

		result := engine.Execute(input)
		conn.Write([]byte(result + "\n"))
	}
}
