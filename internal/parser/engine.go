// internal/parser/engine.go
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Hareesh108/haruDB/internal/storage"
)

const (
	ErrSyntaxError = "Syntax error"
)

type Engine struct {
	DB *storage.Database
}

func NewEngine(dataDir string) *Engine {
	return &Engine{DB: storage.NewDatabase(dataDir)}
}

func (e *Engine) Execute(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimSuffix(input, ";") // remove trailing semicolon

	upper := strings.ToUpper(input)

	switch {
	case strings.HasPrefix(upper, "CREATE TABLE"):
		// CREATE TABLE users (id, name)
		parts := strings.SplitN(input, "(", 2)
		if len(parts) < 2 {
			return ErrSyntaxError
		}
		header := strings.TrimSpace(parts[0])
		fields := strings.Fields(header)
		if len(fields) < 3 {
			return ErrSyntaxError
		}
		tableName := fields[2]

		colsRaw := strings.TrimSuffix(parts[1], ")")
		columns := strings.Split(colsRaw, ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}
		return e.DB.CreateTable(tableName, columns)

	case strings.HasPrefix(upper, "INSERT INTO"):
		// INSERT INTO users VALUES (1, 'Hareesh')
		parts := strings.SplitN(input, "VALUES", 2)
		if len(parts) < 2 {
			return ErrSyntaxError
		}
		tableName := strings.Fields(parts[0])[2]
		tableName = strings.ToLower(tableName)

		valRaw := strings.Trim(parts[1], " ();")
		values := strings.Split(valRaw, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
			values[i] = strings.Trim(values[i], "'")
		}
		return e.DB.Insert(tableName, values)

	case strings.HasPrefix(upper, "SELECT * FROM"):
		// SELECT * FROM users
		parts := strings.Fields(input)
		if len(parts) < 4 {
			return ErrSyntaxError
		}
		tableName := strings.ToLower(parts[3])
		return e.DB.SelectAll(tableName)

	case strings.HasPrefix(upper, "UPDATE"):
		// UPDATE users SET name = 'NewName' WHERE id = 1
		// For simplicity, we'll use row index instead of WHERE clause
		// UPDATE users SET name = 'NewName' ROW 0
		parts := strings.Fields(input)
		if len(parts) < 6 {
			return "Syntax error: UPDATE table SET column = value ROW index"
		}
		tableName := strings.ToLower(parts[1])

		// Find SET clause
		setIndex := -1
		for i, part := range parts {
			if strings.ToUpper(part) == "SET" {
				setIndex = i
				break
			}
		}
		if setIndex == -1 || setIndex+3 >= len(parts) {
			return "Syntax error: missing SET clause"
		}

		// Find ROW clause
		rowIndex := -1
		for i, part := range parts {
			if strings.ToUpper(part) == "ROW" && i+1 < len(parts) {
				// Parse row index
				if idx, err := strconv.Atoi(parts[i+1]); err == nil {
					rowIndex = idx
					break
				}
			}
		}
		if rowIndex == -1 {
			return "Syntax error: missing ROW index"
		}

		// Get table to find column index
		table, exists := e.DB.Tables[tableName]
		if !exists {
			return fmt.Sprintf("Table %s not found", tableName)
		}

		// Find column index
		columnName := parts[setIndex+1]
		columnIndex := -1
		for i, col := range table.Columns {
			if col == columnName {
				columnIndex = i
				break
			}
		}
		if columnIndex == -1 {
			return fmt.Sprintf("Column %s not found", columnName)
		}

		// Get new value
		newValue := strings.Trim(parts[setIndex+3], "'")

		// Create new row with updated value
		if rowIndex >= len(table.Rows) {
			return "Row index out of bounds"
		}
		newRow := make([]string, len(table.Rows[rowIndex]))
		copy(newRow, table.Rows[rowIndex])
		newRow[columnIndex] = newValue

		return e.DB.Update(tableName, rowIndex, newRow)

	case strings.HasPrefix(upper, "DELETE FROM"):
		// DELETE FROM users ROW 0
		parts := strings.Fields(input)
		if len(parts) < 4 {
			return "Syntax error: DELETE FROM table ROW index"
		}
		tableName := strings.ToLower(parts[2])

		// Find ROW clause
		rowIndex := -1
		for i, part := range parts {
			if strings.ToUpper(part) == "ROW" && i+1 < len(parts) {
				// Parse row index
				if idx, err := strconv.Atoi(parts[i+1]); err == nil {
					rowIndex = idx
					break
				}
			}
		}
		if rowIndex == -1 {
			return "Syntax error: missing ROW index"
		}

		return e.DB.Delete(tableName, rowIndex)

	case strings.HasPrefix(upper, "DROP TABLE"):
		// DROP TABLE users
		parts := strings.Fields(input)
		if len(parts) < 3 {
			return "Syntax error: DROP TABLE table_name"
		}
		tableName := strings.ToLower(parts[2])
		return e.DB.DropTable(tableName)

	default:
		return "Unknown command"
	}
}
