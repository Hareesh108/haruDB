// internal/storage/persist.go
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type onDiskTable struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

func (db *Database) tablePath(name string) string {
	name = strings.ToLower(name)
	return filepath.Join(db.DataDir, name+".harudb")
}

func (db *Database) saveTable(t *Table) error {
	payload := onDiskTable{
		Name:    t.Name,
		Columns: t.Columns,
		Rows:    t.Rows,
	}
	data, err := json.MarshalIndent(&payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal table %s: %w", t.Name, err)
	}
	if err := os.WriteFile(db.tablePath(t.Name), data, 0644); err != nil {
		return fmt.Errorf("write table %s: %w", t.Name, err)
	}
	return nil
}

func (db *Database) loadTables() error {
	entries, err := os.ReadDir(db.DataDir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".harudb" {
			continue
		}
		path := filepath.Join(db.DataDir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			// best-effort: skip unreadable
			continue
		}
		var disk onDiskTable
		if err := json.Unmarshal(raw, &disk); err != nil {
			continue
		}
		name := strings.TrimSuffix(strings.ToLower(e.Name()), ".harudb")
		if disk.Name != "" {
			name = strings.ToLower(disk.Name)
		}
		db.Tables[name] = &Table{
			Name:    name,
			Columns: disk.Columns,
			Rows:    disk.Rows,
		}
	}
	return nil
}
