// internal/parser/where_integration_test.go
package parser

import (
	"strings"
	"testing"

	"github.com/Hareesh108/haruDB/internal/storage"
)

func TestAdvancedWhereIntegration(t *testing.T) {
	dataDir := t.TempDir()
	db := storage.NewDatabase(dataDir)

	// Create test table
	if msg := db.CreateTable("employees", []string{"id", "name", "age", "salary", "department", "status"}); !strings.Contains(msg, "created") {
		t.Fatalf("create table failed: %s", msg)
	}

	// Insert test data
	testData := [][]string{
		{"1", "John Doe", "25", "50000", "Engineering", "active"},
		{"2", "Jane Smith", "30", "60000", "Marketing", "active"},
		{"3", "Bob Johnson", "35", "70000", "Engineering", "inactive"},
		{"4", "Alice Brown", "28", "55000", "Sales", "active"},
		{"5", "Charlie Wilson", "45", "80000", "Engineering", "active"},
		{"6", "Diana Prince", "22", "45000", "Marketing", "active"},
	}

	for _, row := range testData {
		_ = db.Insert("employees", row)
	}

	// Create indexes for better performance
	_ = db.CreateIndex("employees", "age")
	_ = db.CreateIndex("employees", "department")
	_ = db.CreateIndex("employees", "status")

	tests := []struct {
		name         string
		whereClause  string
		expectedRows int
		description  string
	}{
		{
			name:         "simple equality",
			whereClause:  "age = 25",
			expectedRows: 1,
			description:  "Find employees aged exactly 25",
		},
		{
			name:         "not equals",
			whereClause:  "department != 'Engineering'",
			expectedRows: 4,
			description:  "Find employees not in Engineering",
		},
		{
			name:         "less than",
			whereClause:  "age < 30",
			expectedRows: 4,
			description:  "Find employees younger than 30",
		},
		{
			name:         "greater than",
			whereClause:  "salary > 55000",
			expectedRows: 4,
			description:  "Find employees with salary > 55000",
		},
		{
			name:         "and condition",
			whereClause:  "age > 25 AND department = 'Engineering'",
			expectedRows: 2,
			description:  "Find Engineering employees older than 25",
		},
		{
			name:         "or condition",
			whereClause:  "age < 25 OR age > 40",
			expectedRows: 2,
			description:  "Find employees younger than 25 or older than 40",
		},
		{
			name:         "like pattern",
			whereClause:  "name LIKE 'J%'",
			expectedRows: 2,
			description:  "Find employees with names starting with J",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse WHERE clause
			whereExpr, err := ParseWhereClause(tt.whereClause)
			if err != nil {
				t.Fatalf("Failed to parse WHERE clause '%s': %v", tt.whereClause, err)
			}

			// Execute query
			result := db.SelectWhereAdvanced("employees", whereExpr)

			// Count result rows (excluding header and "(no rows)" message)
			lines := strings.Split(strings.TrimSpace(result), "\n")
			rowCount := 0
			for _, line := range lines {
				if line != "" && !strings.Contains(line, " | ") {
					// Skip header line
					continue
				}
				if strings.Contains(line, " | ") {
					rowCount++
				}
			}

			if rowCount != tt.expectedRows {
				t.Errorf("Expected %d rows, got %d for query: %s\nResult:\n%s",
					tt.expectedRows, rowCount, tt.whereClause, result)
			}

			t.Logf("✓ %s: %s -> %d rows", tt.name, tt.description, rowCount)
		})
	}
}
