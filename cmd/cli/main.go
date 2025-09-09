// cmd/cli/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

func main() {
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
	conn, err := net.Dial("tcp", "localhost:54321")
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

	fmt.Println("Connected to HaruDB 🎉 (history enabled with ↑↓)")

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
