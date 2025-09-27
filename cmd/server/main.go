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
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Hareesh108/haruDB/internal/auth"
	"github.com/Hareesh108/haruDB/internal/parser"
)

const DB_VERSION string = "v0.0.5"

// checkPortUsage checks what process is using the specified port
func checkPortUsage(port string) {
	// Try to connect to the port to see if something is listening
	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		return // Port is free
	}
	conn.Close()

	// Port is in use, try to identify what's using it
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("lsof", "-i", ":"+port)
	case "windows":
		cmd = exec.Command("netstat", "-ano", "-p", "TCP", "-f", "inet")
	default:
		fmt.Printf("⚠️  Port %s is already in use by another process\n", port)
		fmt.Printf("   Please stop the other service or use a different port\n")
		return
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Port %s is already in use by another process\n", port)
		fmt.Printf("   Please stop the other service or use a different port\n")
		return
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		// Parse the output to find the process
		for _, line := range lines[1:] { // Skip header
			if strings.Contains(line, ":"+port) {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					processName := parts[0]
					if strings.Contains(strings.ToLower(processName), "harudb") {
						fmt.Printf("ℹ️  Another HaruDB instance is already running on port %s\n", port)
						fmt.Printf("   Please stop the existing instance or use a different port\n")
					} else {
						fmt.Printf("⚠️  Port %s is already in use by: %s\n", port, processName)
						fmt.Printf("   Please stop the other service or use a different port\n")
						fmt.Printf("   Common solutions:\n")
						fmt.Printf("   - Stop the other service: sudo systemctl stop <service>\n")
						fmt.Printf("   - Kill the process: sudo kill -9 <PID>\n")
						fmt.Printf("   - Use a different port: ./harudb --port 54322\n")
					}
					return
				}
			}
		}
	}

	fmt.Printf("⚠️  Port %s is already in use by another process\n", port)
	fmt.Printf("   Please stop the other service or use a different port\n")
}

func main() {
	dataDir := flag.String("data-dir", "./data", "Directory to store .harudb files")
	enableTLS := flag.Bool("tls", false, "Enable TLS encryption")
	port := flag.String("port", "54321", "Port to listen on")
	flag.Parse()

	// Check if port is already in use
	checkPortUsage(*port)

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
		tcpListener, err := net.Listen("tcp", ":"+*port)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", *port, err)
		}
		listener = tls.NewListener(tcpListener, tlsManager.GetTLSConfig())
		fmt.Printf("🚀 HaruDB server started on port %s with TLS (data dir: %s)\n", *port, *dataDir)
	} else {
		// Create regular TCP listener
		listener, err = net.Listen("tcp", ":"+*port)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", *port, err)
		}
		fmt.Printf("🚀 HaruDB server started on port %s (data dir: %s)\n", *port, *dataDir)
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
