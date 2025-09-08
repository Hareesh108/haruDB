// internal/storage/memory.go
package storage

import (
	"fmt"
	"strings"
)

type Table struct {
	Name    string
	Columns []string
	Rows    [][]string
}

type Database struct {
	DataDir string
	Tables  map[string]*Table
}

func NewDatabase(dataDir string) *Database {
	db := &Database{
		DataDir: dataDir,
		Tables:  make(map[string]*Table),
	}
	// Load any existing .harudb files
	_ = db.loadTables()
	return db
}

func (db *Database) CreateTable(name string, columns []string) string {
	name = strings.ToLower(name)
	if _, exists := db.Tables[name]; exists {
		return fmt.Sprintf("Table %s already exists", name)
	}
	db.Tables[name] = &Table{Name: name, Columns: columns, Rows: [][]string{}}
	if err := db.saveTable(db.Tables[name]); err != nil {
		return fmt.Sprintf("Table %s created (warning: failed to persist: %v)", name, err)
	}
	return fmt.Sprintf("Table %s created", name)
}

func (db *Database) Insert(tableName string, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf("Table %s not found", tableName)
	}
	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}
	table.Rows = append(table.Rows, values)
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("1 row inserted (warning: failed to persist: %v)", err)
	}
	return "1 row inserted"
}

func (db *Database) SelectAll(tableName string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf("Table %s not found", tableName)
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
