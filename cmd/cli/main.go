// cmd/cli/main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

func main() {
	port := flag.String("port", "54321", "Port to connect to")
	host := flag.String("host", "localhost", "Host to connect to")
	flag.Parse()

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// history file
	historyFile := filepath.Join(os.TempDir(), ".harudb_history")
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// connect to server
	serverAddr := *host + ":" + *port
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("❌ Failed to connect:", err)
		return
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)

	// read server welcome banner until first prompt
	for {
		lineStr, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("❌ Connection closed")
			return
		}
		fmt.Print(lineStr)
		if strings.HasPrefix(lineStr, "haruDB> ") {
			break
		}
	}

	// Show initial help
	fmt.Println("\n💡 Type 'HELP' for available commands")
	fmt.Println("🔐 You need to login first: LOGIN admin admin123")

	for {
		// show CLI prompt
		input, err := line.Prompt("haruDB> ")
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		line.AppendHistory(input)

		// send command to server
		fmt.Fprintln(conn, input)

		// exit immediately if user typed exit
		if input == "exit" {
			break
		}

		// read server response line by line until next prompt
		for {
			respLine, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Connection closed")
				return
			}
			if strings.HasPrefix(respLine, "haruDB> ") {
				// prompt detected → break to show CLI prompt
				break
			}
			fmt.Print(respLine)
		}
	}

	// save history
	if f, err := os.Create(historyFile); err == nil {
		line.WriteHistory(f)
		f.Close()
	}
}
