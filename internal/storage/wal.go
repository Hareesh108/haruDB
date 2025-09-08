// internal/storage/wal.go
package storage

import (
	"os"
	"path/filepath"
	"strings"
)

// WAL struct
type WAL struct {
	file *os.File
	path string
}

// Open or create WAL
func NewWAL(dataDir string) (*WAL, error) {
	walPath := filepath.Join(dataDir, "harudb.wal")
	f, err := os.OpenFile(walPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f, path: walPath}, nil
}

// Append entry to WAL
func (w *WAL) WriteEntry(cmd string) error {
	_, err := w.file.WriteString(cmd + "\n")
	if err != nil {
		return err
	}
	return w.file.Sync() // flush to disk
}

// Replay WAL for recovery
func (w *WAL) Replay() ([]string, error) {
	data, err := os.ReadFile(w.path)
	if err != nil {
		// If file doesn't exist, just return empty (first run)
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Split into lines (commands)
	content := string(data)
	lines := []string{}
	for _, l := range strings.Split(content, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		lines = append(lines, l)
	}
	return lines, nil
}

// Close WAL
func (w *WAL) Close() error {
	return w.file.Close()
}
