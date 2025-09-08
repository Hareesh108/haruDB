// internal/parser/simple.go
package parser

import (
	"strings"

	"github.com/Hareesh108/haruDB/internal/storage"
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
			return "Syntax error"
		}
		header := strings.TrimSpace(parts[0])
		fields := strings.Fields(header)
		if len(fields) < 3 {
			return "Syntax error"
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
			return "Syntax error"
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
			return "Syntax error"
		}
		tableName := strings.ToLower(parts[3])
		return e.DB.SelectAll(tableName)

	default:
		return "Unknown command"
	}
}
