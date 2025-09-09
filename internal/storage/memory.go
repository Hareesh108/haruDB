// internal/storage/memory.go
package storage

import (
	"fmt"
	"os"
	"strings"
)

const (
	ErrTableNotFound = "Table %s not found"
	ErrWALCheckpoint = "Warning: Failed to write WAL checkpoint: %v\n"
)

type Table struct {
	Name    string
	Columns []string
	Rows    [][]string
}

type Database struct {
	DataDir string
	Tables  map[string]*Table
	WAL     *WALManager
}

func NewDatabase(dataDir string) *Database {
	db := &Database{
		DataDir: dataDir,
		Tables:  make(map[string]*Table),
	}

	// Initialize WAL manager
	var err error
	db.WAL, err = NewWALManager(dataDir)
	if err != nil {
		// If WAL initialization fails, continue without WAL (degraded mode)
		fmt.Printf("Warning: Failed to initialize WAL: %v\n", err)
	}

	// Load any existing .harudb files first
	_ = db.loadTables()

	// Replay WAL entries if WAL is available (this will apply any changes not yet persisted)
	if db.WAL != nil {
		if err := db.WAL.ReplayWAL(db); err != nil {
			fmt.Printf("Warning: Failed to replay WAL: %v\n", err)
		}
	}

	return db
}

func (db *Database) CreateTable(name string, columns []string) string {
	name = strings.ToLower(name)
	if _, exists := db.Tables[name]; exists {
		return fmt.Sprintf("Table %s already exists", name)
	}

	// Write to WAL (Write Ahead Logs) first
	if db.WAL != nil {
		data := map[string]interface{}{
			"columns": columns,
		}
		if err := db.WAL.WriteEntry(WAL_CREATE_TABLE, name, data); err != nil {
			return fmt.Sprintf("Table %s created (warning: failed to write to WAL: %v)", name, err)
		}
	}

	// Apply changes to memory
	db.Tables[name] = &Table{Name: name, Columns: columns, Rows: [][]string{}}

	// Persist to disk
	if err := db.saveTable(db.Tables[name]); err != nil {
		return fmt.Sprintf("Table %s created (warning: failed to persist: %v)", name, err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return fmt.Sprintf("Table %s created", name)
}

func (db *Database) Insert(tableName string, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}
	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"values": values,
		}
		if err := db.WAL.WriteEntry(WAL_INSERT, tableName, data); err != nil {
			return fmt.Sprintf("1 row inserted (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows = append(table.Rows, values)

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("1 row inserted (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row inserted"
}

func (db *Database) SelectAll(tableName string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	result := strings.Join(table.Columns, " | ") + "\n"
	for _, row := range table.Rows {
		result += strings.Join(row, " | ") + "\n"
	}
	if len(table.Rows) == 0 {
		result += "(no rows)\n"
	}
	return result
}

// Update updates a row in the specified table
func (db *Database) Update(tableName string, rowIndex int, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
			"values":    values,
		}
		if err := db.WAL.WriteEntry(WAL_UPDATE, tableName, data); err != nil {
			return fmt.Sprintf("Row updated (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows[rowIndex] = values

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("Row updated (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row updated"
}

// Delete deletes a row from the specified table
func (db *Database) Delete(tableName string, rowIndex int) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
		}
		if err := db.WAL.WriteEntry(WAL_DELETE, tableName, data); err != nil {
			return fmt.Sprintf("Row deleted (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows = append(table.Rows[:rowIndex], table.Rows[rowIndex+1:]...)

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("Row deleted (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row deleted"
}

// DropTable drops the specified table
func (db *Database) DropTable(tableName string) string {
	tableName = strings.ToLower(tableName)
	_, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// Write to WAL first
	if db.WAL != nil {
		if err := db.WAL.WriteEntry(WAL_DROP_TABLE, tableName, nil); err != nil {
			return fmt.Sprintf("Table dropped (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	delete(db.Tables, tableName)

	// Remove table file from disk
	tablePath := db.tablePath(tableName)
	if err := os.Remove(tablePath); err != nil && !os.IsNotExist(err) {
		return fmt.Sprintf("Table dropped (warning: failed to remove table file: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return fmt.Sprintf("Table %s dropped", tableName)
}
