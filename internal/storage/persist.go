// internal/storage/persist.go
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// onDiskTable is the JSON layout stored in .harudb files
type onDiskTable struct {
	Name           string     `json:"name"`
	Columns        []string   `json:"columns"`
	Rows           [][]string `json:"rows"`
	IndexedColumns []string   `json:"indexed_columns,omitempty"`
}

// tablePath returns the target .harudb file path for a table
func (db *Database) tablePath(name string) string {
	name = strings.ToLower(name)
	return filepath.Join(db.DataDir, name+".harudb")
}

// saveTable writes a table atomically to disk using a temp file + rename.
// It writes the temp file in the same directory (required for atomic rename),
// fsyncs the file, closes it, renames to the final path, and fsyncs the directory.
func (db *Database) saveTable(t *Table) error {
	// Prepare serialized payload
	payload := onDiskTable{
		Name:           t.Name,
		Columns:        t.Columns,
		Rows:           t.Rows,
		IndexedColumns: t.IndexedColumns,
	}
	data, err := json.MarshalIndent(&payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal table %s: %w", t.Name, err)
	}

	dir := db.DataDir
	finalPath := db.tablePath(t.Name)

	// Create a temp file in the same directory. Pattern ensures readable name.
	tempFile, err := os.CreateTemp(dir, t.Name+".harudb.tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file for %s: %w", t.Name, err)
	}
	tempPath := tempFile.Name()

	// Write the bytes
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("write temp file %s: %w", tempPath, err)
	}

	// Ensure data is flushed to disk for the temp file
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("fsync temp file %s: %w", tempPath, err)
	}

	// Close the temp file
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("close temp file %s: %w", tempPath, err)
	}

	// Try atomic rename. On some platforms (Windows) rename may fail if destination exists:
	if err := os.Rename(tempPath, finalPath); err != nil {
		// Try removing existing final path and retry rename (best-effort for Windows)
		if removeErr := os.Remove(finalPath); removeErr == nil {
			// try again
			if err2 := os.Rename(tempPath, finalPath); err2 == nil {
				// success
				goto afterRename
			}
		}
		// cleanup and return original error
		os.Remove(tempPath)
		return fmt.Errorf("rename temp %s -> final %s: %w", tempPath, finalPath, err)
	}

afterRename:
	// Best-effort fsync of containing directory so the rename is durable.
	// If this fails, we still return an error so callers know persistence may be weaker.
	if err := syncDir(dir); err != nil {
		return fmt.Errorf("sync dir %s: %w", dir, err)
	}

	return nil
}

// loadTables loads all .harudb files from DataDir into db.Tables (best-effort).
func (db *Database) loadTables() error {
	entries, err := os.ReadDir(db.DataDir)
	if err != nil {
		// If data dir doesn't exist or unreadable, return error (caller can ignore)
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
			// skip unreadable files
			continue
		}
		var disk onDiskTable
		if err := json.Unmarshal(raw, &disk); err != nil {
			// skip invalid JSON (do not stop loading other tables)
			continue
		}
		name := strings.TrimSuffix(strings.ToLower(e.Name()), ".harudb")
		if disk.Name != "" {
			name = strings.ToLower(disk.Name)
		}
		t := &Table{
			Name:           name,
			Columns:        disk.Columns,
			Rows:           disk.Rows,
			IndexedColumns: disk.IndexedColumns,
			Indexes:        make(map[string]map[string][]int),
		}
		db.Tables[name] = t
		db.rebuildAllIndexes(t)
	}
	return nil
}

// syncDir opens the directory and calls Sync() so the rename is durable on disk.
// Best-effort: returns error if sync fails.
func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	// On some platforms Sync may not be supported; return whatever it gives.
	return d.Sync()
}
