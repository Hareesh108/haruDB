// internal/storage/wal.go
package storage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WALEntryType represents the type of operation in WAL
type WALEntryType uint8

const (
	WAL_CREATE_TABLE WALEntryType = iota + 1
	WAL_INSERT
	WAL_UPDATE
	WAL_DELETE
	WAL_DROP_TABLE
	WAL_CHECKPOINT
)

// WALEntry represents a single entry in the WAL
type WALEntry struct {
	Timestamp time.Time    `json:"timestamp"`
	Type      WALEntryType `json:"type"`
	TableName string       `json:"table_name"`
	Data      interface{}  `json:"data"`
}

// WALManager handles Write-Ahead Logging
type WALManager struct {
	dataDir    string
	walFile    *os.File
	walPath    string
	mu         sync.Mutex
	checkpoint time.Time
}

// NewWALManager creates a new WAL manager
func NewWALManager(dataDir string) (*WALManager, error) {
	walPath := filepath.Join(dataDir, "wal.log")

	// Open or create WAL file
	walFile, err := os.OpenFile(walPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}

	wm := &WALManager{
		dataDir:    dataDir,
		walFile:    walFile,
		walPath:    walPath,
		checkpoint: time.Now(),
	}

	return wm, nil
}

// Close closes the WAL file
func (wm *WALManager) Close() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.walFile != nil {
		return wm.walFile.Close()
	}
	return nil
}

// WriteEntry writes an entry to the WAL
func (wm *WALManager) WriteEntry(entryType WALEntryType, tableName string, data interface{}) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	entry := WALEntry{
		Timestamp: time.Now(),
		Type:      entryType,
		TableName: tableName,
		Data:      data,
	}

	// Serialize entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal WAL entry: %w", err)
	}

	// Write entry length (4 bytes) + entry data
	length := uint32(len(jsonData))

	// Write length
	if err := binary.Write(wm.walFile, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("failed to write WAL entry length: %w", err)
	}

	// Write data
	if _, err := wm.walFile.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write WAL entry data: %w", err)
	}

	// Flush to ensure data is written to disk
	if err := wm.walFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync WAL file: %w", err)
	}

	return nil
}

// WriteCheckpoint writes a checkpoint entry
func (wm *WALManager) WriteCheckpoint() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.checkpoint = time.Now()

	entry := WALEntry{
		Timestamp: time.Now(),
		Type:      WAL_CHECKPOINT,
		TableName: "",
		Data:      nil,
	}

	// Serialize entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal WAL entry: %w", err)
	}

	// Write entry length (4 bytes) + entry data
	length := uint32(len(jsonData))

	// Write length
	if err := binary.Write(wm.walFile, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("failed to write WAL entry length: %w", err)
	}

	// Write data
	if _, err := wm.walFile.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write WAL entry data: %w", err)
	}

	// Flush to ensure data is written to disk
	if err := wm.walFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync WAL file: %w", err)
	}

	return nil
}

// ReplayWAL replays WAL entries since last checkpoint
func (wm *WALManager) ReplayWAL(db *Database) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Close current WAL file
	if wm.walFile != nil {
		wm.walFile.Close()
	}

	// Open WAL file for reading
	walFile, err := os.Open(wm.walPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No WAL file exists, nothing to replay
			return nil
		}
		return fmt.Errorf("failed to open WAL file for replay: %w", err)
	}
	defer walFile.Close()

	reader := bufio.NewReader(walFile)

	for {
		// Read entry length
		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			if err.Error() == "EOF" {
				break // End of file
			}
			return fmt.Errorf("failed to read WAL entry length: %w", err)
		}

		// Read entry data
		jsonData := make([]byte, length)
		if _, err := reader.Read(jsonData); err != nil {
			return fmt.Errorf("failed to read WAL entry data: %w", err)
		}

		// Deserialize entry
		var entry WALEntry
		if err := json.Unmarshal(jsonData, &entry); err != nil {
			return fmt.Errorf("failed to unmarshal WAL entry: %w", err)
		}

		// For now, replay all entries (we'll optimize this later)
		// Skip entries before checkpoint (but always process CHECKPOINT entries)
		// if entry.Type != WAL_CHECKPOINT && entry.Timestamp.Before(wm.checkpoint) {
		// 	continue
		// }

		// Replay entry
		if err := wm.replayEntry(db, &entry); err != nil {
			return fmt.Errorf("failed to replay WAL entry: %w", err)
		}
	}

	// Reopen WAL file for writing
	wm.walFile, err = os.OpenFile(wm.walPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen WAL file: %w", err)
	}

	return nil
}

// replayEntry replays a single WAL entry
func (wm *WALManager) replayEntry(db *Database, entry *WALEntry) error {
	switch entry.Type {
	case WAL_CREATE_TABLE:
		if data, ok := entry.Data.(map[string]interface{}); ok {
			if columns, ok := data["columns"].([]interface{}); ok {
				colStrs := make([]string, len(columns))
				for i, col := range columns {
					colStrs[i] = col.(string)
				}
				db.Tables[entry.TableName] = &Table{
					Name:    entry.TableName,
					Columns: colStrs,
					Rows:    [][]string{},
				}
				_ = db.saveTable(db.Tables[entry.TableName])
			}
		}

	case WAL_INSERT:
		if data, ok := entry.Data.(map[string]interface{}); ok {
			if values, ok := data["values"].([]interface{}); ok {
				valStrs := make([]string, len(values))
				for i, val := range values {
					valStrs[i] = val.(string)
				}
				if table, exists := db.Tables[entry.TableName]; exists {
					table.Rows = append(table.Rows, valStrs)
					_ = db.saveTable(table)
				}
			}
		}

	case WAL_UPDATE:
		if data, ok := entry.Data.(map[string]interface{}); ok {
			if rowIndex, ok := data["row_index"].(float64); ok {
				if values, ok := data["values"].([]interface{}); ok {
					valStrs := make([]string, len(values))
					for i, val := range values {
						valStrs[i] = val.(string)
					}
					if table, exists := db.Tables[entry.TableName]; exists {
						if int(rowIndex) < len(table.Rows) {
							table.Rows[int(rowIndex)] = valStrs
							_ = db.saveTable(table)
						}
					}
				}
			}
		}

	case WAL_DELETE:
		if data, ok := entry.Data.(map[string]interface{}); ok {
			if rowIndex, ok := data["row_index"].(float64); ok {
				if table, exists := db.Tables[entry.TableName]; exists {
					if int(rowIndex) < len(table.Rows) {
						// Remove row at index
						table.Rows = append(table.Rows[:int(rowIndex)], table.Rows[int(rowIndex)+1:]...)
						_ = db.saveTable(table)
					}
				}
			}
		}

	case WAL_DROP_TABLE:
		delete(db.Tables, entry.TableName)
		os.Remove(db.tablePath(entry.TableName))

	case WAL_CHECKPOINT:
		// Update checkpoint time
		wm.checkpoint = entry.Timestamp
	}

	return nil
}

// TruncateWAL truncates the WAL file after successful checkpoint
func (wm *WALManager) TruncateWAL() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Close current file
	if wm.walFile != nil {
		wm.walFile.Close()
	}

	// Truncate WAL file
	if err := os.Truncate(wm.walPath, 0); err != nil {
		return fmt.Errorf("failed to truncate WAL file: %w", err)
	}

	// Reopen for writing
	var err error
	wm.walFile, err = os.OpenFile(wm.walPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen WAL file after truncation: %w", err)
	}

	return nil
}
