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
	case strings.HasPrefix(upper, "CREATE INDEX"):
		// CREATE INDEX ON users (email)
		parts := strings.SplitN(input, "(", 2)
		if len(parts) < 2 {
			return ErrSyntaxError
		}
		header := strings.TrimSpace(parts[0])
		seg := strings.Fields(header)
		if len(seg) < 4 { // CREATE INDEX ON <table>
			return ErrSyntaxError
		}
		tableName := strings.ToLower(seg[3])
		col := strings.TrimSpace(parts[1])
		col = strings.TrimSuffix(col, ")")
		col = strings.TrimSpace(col)
		return e.DB.CreateIndex(tableName, col)

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
		// SELECT * FROM users [WHERE col = 'val']
		// Basic WHERE support
		parts := strings.Fields(input)
		if len(parts) < 4 {
			return ErrSyntaxError
		}
		tableName := strings.ToLower(parts[3])
		// Check for WHERE clause
		whereIdx := -1
		for i, p := range parts {
			if strings.ToUpper(p) == "WHERE" {
				whereIdx = i
				break
			}
		}
		if whereIdx == -1 {
			return e.DB.SelectAll(tableName)
		}
		// Expect: WHERE <col> = <value>
		if whereIdx+3 >= len(parts) {
			return ErrSyntaxError
		}
		col := parts[whereIdx+1]
		op := parts[whereIdx+2]
		val := strings.Trim(parts[whereIdx+3], "'\"")
		if op != "=" {
			return "Only equality WHERE supported"
		}
		return e.DB.SelectWhere(tableName, col, val)

	case strings.HasPrefix(upper, "UPDATE"):
		// Example: UPDATE users SET name = 'NewName', email = 'new@example.com' ROW 0
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
		if setIndex == -1 {
			return "Syntax error: missing SET clause"
		}

		// Find ROW clause
		rowIndex := -1
		for i, part := range parts {
			if strings.ToUpper(part) == "ROW" && i+1 < len(parts) {
				if idx, err := strconv.Atoi(parts[i+1]); err == nil {
					rowIndex = idx
					break
				}
			}
		}
		if rowIndex == -1 {
			return "Syntax error: missing ROW index"
		}

		// Get table
		table, exists := e.DB.Tables[tableName]
		if !exists {
			return fmt.Sprintf("Table %s not found", tableName)
		}
		if rowIndex < 0 || rowIndex >= len(table.Rows) {
			return "Row index out of bounds"
		}

		// Reconstruct SET clause (everything between SET and ROW)
		setClause := strings.Join(parts[setIndex+1:], " ")
		rowClauseIndex := strings.Index(strings.ToUpper(setClause), "ROW")
		if rowClauseIndex != -1 {
			setClause = setClause[:rowClauseIndex]
		}
		setClause = strings.TrimSpace(setClause)

		// Split multiple assignments by comma
		assignments := strings.Split(setClause, ",")
		newRow := make([]string, len(table.Rows[rowIndex]))
		copy(newRow, table.Rows[rowIndex])

		for _, assign := range assignments {
			assign = strings.TrimSpace(assign)
			if assign == "" {
				continue
			}
			kv := strings.SplitN(assign, "=", 2)
			if len(kv) != 2 {
				return fmt.Sprintf("Invalid assignment: %s", assign)
			}
			columnName := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			value = strings.Trim(value, "'")

			// Find column index
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

			// Apply update
			newRow[columnIndex] = value
		}

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
