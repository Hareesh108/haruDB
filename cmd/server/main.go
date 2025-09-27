// cmd/server/main.go
package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Hareesh108/haruDB/internal/auth"
	"github.com/Hareesh108/haruDB/internal/parser"
)

const DB_VERSION string = "v0.0.4"

func main() {
	port := "54321"

	dataDir := flag.String("data-dir", "./data", "Directory to store .harudb files")
	enableTLS := flag.Bool("tls", false, "Enable TLS encryption")
	flag.Parse()

	// Make sure the data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data dir %s: %v", *dataDir, err)
	}

	// Initialize TLS manager if enabled
	var tlsManager *auth.TLSManager
	if *enableTLS {
		tlsManager = auth.NewTLSManager(*dataDir)
		if !tlsManager.IsTLSEnabled() {
			log.Printf("Warning: TLS requested but not properly configured")
		} else {
			fmt.Printf("🔒 TLS encryption enabled\n")
		}
	}

	var listener net.Listener
	var err error

	if *enableTLS && tlsManager != nil && tlsManager.IsTLSEnabled() {
		// Create TLS listener
		tcpListener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", port, err)
		}
		listener = tls.NewListener(tcpListener, tlsManager.GetTLSConfig())
		fmt.Printf("🚀 HaruDB server started on port %s with TLS (data dir: %s)\n", port, *dataDir)
	} else {
		// Create regular TCP listener
		listener, err = net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", port, err)
		}
		fmt.Printf("🚀 HaruDB server started on port %s (data dir: %s)\n", port, *dataDir)
	}
	defer listener.Close()

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
	conn.Write([]byte("🔐 Authentication Required\n"))
	conn.Write([]byte("Default admin: admin / admin123\n"))
	conn.Write([]byte("Please change the default password after first login!\n\n"))

	scanner := bufio.NewScanner(conn)
	for {
		// send prompt with newline
		conn.Write([]byte("haruDB> \n"))

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			conn.Write([]byte("Goodbye 👋\n"))
			break
		}

		// Execute with timeout to prevent hanging
		resultChan := make(chan string, 1)
		go func() {
			result := engine.Execute(input)
			resultChan <- result
		}()

		var result string
		select {
		case result = <-resultChan:
			// Command completed successfully
		case <-time.After(10 * time.Second):
			// Command timed out
			result = "Error: Command timed out after 10 seconds"
		}

		if !strings.HasSuffix(result, "\n") {
			result += "\n"
		}

		// send result
		conn.Write([]byte(result))
	}
}
