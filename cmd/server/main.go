// cmd/server/main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Hareesh108/haruDB/internal/parser"
)

const DB_VERSION string = "v0.0.2"

func main() {
	port := "54321"

	dataDir := flag.String("data-dir", "./data", "Directory to store .harudb files")
	flag.Parse()

	// Make sure the data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data dir %s: %v", *dataDir, err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("🚀 HaruDB server started on port %s (data dir: %s)\n", port, *dataDir)

	engine := parser.NewEngine(*dataDir)

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

	fmt.Fprintf(conn, "\nWelcome to HaruDB %s 🎉\n", DB_VERSION)
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
