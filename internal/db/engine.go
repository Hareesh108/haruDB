package db

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Table struct {
	Cols []string
	Rows [][]string
	mu   sync.RWMutex
}

type Database struct {
	tables map[string]*Table
	mu     sync.RWMutex
	wal    *os.File
	walMu  sync.Mutex
}

// NewDatabase opens/creates a WAL file and replays it
func NewDatabase(walpath string) (*Database, error) {
	f, err := os.OpenFile(walpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	db := &Database{
		tables: make(map[string]*Table),
		wal:    f,
	}
	if err := db.replayWAL(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) Close() error {
	db.walMu.Lock()
	defer db.walMu.Unlock()
	return db.wal.Close()
}

func (db *Database) appendWALLine(line string) error {
	db.walMu.Lock()
	defer db.walMu.Unlock()
	if _, err := db.wal.WriteString(line + "\n"); err != nil {
		return err
	}
	return db.wal.Sync()
}

func (db *Database) replayWAL() error {
	db.walMu.Lock()
	defer db.walMu.Unlock()
	if _, err := db.wal.Seek(0, io.SeekStart); err != nil {
		return err
	}
	scanner := bufio.NewScanner(db.wal)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if err := db.execLineNoWAL(line); err != nil {
			return fmt.Errorf("replay error on %q: %w", line, err)
		}
	}
	if _, err := db.wal.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	return scanner.Err()
}

// Exec executes a command and logs to WAL
func (db *Database) Exec(line string) (string, error) {
	trim := strings.TrimSpace(line)
	if trim == "" {
		return "", nil
	}
	if err := db.appendWALLine(trim); err != nil {
		return "", err
	}
	if err := db.execLineNoWAL(trim); err != nil {
		return "", err
	}
	return "OK", nil
}

// execLineNoWAL executes without WAL logging (used in replay)
func (db *Database) execLineNoWAL(line string) error {
	upper := strings.ToUpper(line)
	switch {
	case strings.HasPrefix(upper, "CREATE TABLE"):
		return db.handleCreate(line)
	case strings.HasPrefix(upper, "INSERT INTO"):
		return db.handleInsert(line)
	default:
		return fmt.Errorf("unknown command: %s", line)
	}
}

func (db *Database) handleCreate(line string) error {
	parts := strings.SplitN(line, "(", 2)
	if len(parts) != 2 {
		return errors.New("invalid CREATE TABLE syntax")
	}
	before := strings.TrimSpace(parts[0])
	colsPart := strings.TrimSpace(parts[1])
	if !strings.HasSuffix(colsPart, ")") {
		return errors.New("missing ) in CREATE")
	}
	colsPart = colsPart[:len(colsPart)-1]
	beforeParts := strings.Fields(before)
	if len(beforeParts) < 3 {
		return errors.New("invalid CREATE TABLE syntax")
	}
	name := beforeParts[2]
	cols := parseCSV(colsPart)
	if len(cols) == 0 {
		return errors.New("no columns specified")
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, exists := db.tables[name]; exists {
		return fmt.Errorf("table %s exists", name)
	}
	t := &Table{Cols: cols, Rows: make([][]string, 0)}
	db.tables[name] = t
	return nil
}

func (db *Database) handleInsert(line string) error {
	up := strings.ToUpper(line)
	valIdx := strings.Index(up, "VALUES")
	if valIdx == -1 {
		return errors.New("missing VALUES in INSERT")
	}
	before := strings.TrimSpace(line[:valIdx])
	after := strings.TrimSpace(line[valIdx+len("VALUES"):])
	beforeParts := strings.Fields(before)
	if len(beforeParts) < 3 {
		return errors.New("invalid INSERT syntax")
	}
	name := beforeParts[2]
	if !strings.HasPrefix(after, "(") || !strings.HasSuffix(after, ")") {
		return errors.New("values must be in (...)")
	}
	valsPart := after[1 : len(after)-1]
	vals := parseCSV(valsPart)

	db.mu.RLock()
	table, ok := db.tables[name]
	db.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %s not found", name)
	}
	if len(vals) != len(table.Cols) {
		return fmt.Errorf("column count mismatch: expected %d", len(table.Cols))
	}
	table.mu.Lock()
	table.Rows = append(table.Rows, vals)
	table.mu.Unlock()
	return nil
}

func (db *Database) SelectAll(name string) ([][]string, []string, error) {
	db.mu.RLock()
	table, ok := db.tables[name]
	db.mu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("table %s not found", name)
	}
	table.mu.RLock()
	defer table.mu.RUnlock()
	rowsCopy := make([][]string, len(table.Rows))
	for i := range table.Rows {
		r := make([]string, len(table.Rows[i]))
		copy(r, table.Rows[i])
		rowsCopy[i] = r
	}
	cols := make([]string, len(table.Cols))
	copy(cols, table.Cols)
	return rowsCopy, cols, nil
}

func parseCSV(s string) []string {
	var out []string
	cur := strings.Builder{}
	inQuotes := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			inQuotes = !inQuotes
			continue
		}
		if c == ',' && !inQuotes {
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		out = append(out, strings.TrimSpace(cur.String()))
	}
	for i := range out {
		out[i] = strings.Trim(out[i], "\"")
	}
	return out
}
