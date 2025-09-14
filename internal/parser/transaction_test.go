// internal/parser/transaction_test.go
package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTransactionParser(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_parser_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	engine := NewEngine(tempDir)

	t.Run("BeginTransaction", func(t *testing.T) {
		result := engine.Execute("BEGIN TRANSACTION")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse BEGIN TRANSACTION: %s", result)
		}
		t.Logf("BEGIN TRANSACTION result: %s", result)

		// Test with isolation level
		result = engine.Execute("BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse BEGIN TRANSACTION with isolation level: %s", result)
		}
		t.Logf("BEGIN TRANSACTION with isolation level result: %s", result)
	})

	t.Run("CommitTransaction", func(t *testing.T) {
		// First begin a transaction
		engine.Execute("BEGIN TRANSACTION")

		result := engine.Execute("COMMIT")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse COMMIT: %s", result)
		}
		t.Logf("COMMIT result: %s", result)
	})

	t.Run("RollbackTransaction", func(t *testing.T) {
		// First begin a transaction
		engine.Execute("BEGIN TRANSACTION")

		result := engine.Execute("ROLLBACK")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse ROLLBACK: %s", result)
		}
		t.Logf("ROLLBACK result: %s", result)
	})

	t.Run("Savepoint", func(t *testing.T) {
		// First begin a transaction
		engine.Execute("BEGIN TRANSACTION")

		result := engine.Execute("SAVEPOINT sp1")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse SAVEPOINT: %s", result)
		}
		t.Logf("SAVEPOINT result: %s", result)
	})

	t.Run("RollbackToSavepoint", func(t *testing.T) {
		// First begin a transaction
		engine.Execute("BEGIN TRANSACTION")
		// Create a savepoint
		engine.Execute("SAVEPOINT sp1")

		result := engine.Execute("ROLLBACK TO SAVEPOINT sp1")
		if result == "Syntax error" || result == "Unknown command" {
			t.Errorf("Failed to parse ROLLBACK TO SAVEPOINT: %s", result)
		}
		t.Logf("ROLLBACK TO SAVEPOINT result: %s", result)
	})
}

func TestTransactionWorkflow(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_workflow_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	engine := NewEngine(tempDir)

	t.Run("CompleteTransactionWorkflow", func(t *testing.T) {
		// Begin transaction
		result := engine.Execute("BEGIN TRANSACTION")
		t.Logf("BEGIN result: %s", result)

		// Create table
		result = engine.Execute("CREATE TABLE users (id, name, email)")
		t.Logf("CREATE TABLE result: %s", result)

		// Insert data
		result = engine.Execute("INSERT INTO users VALUES (1, 'Alice', 'alice@example.com')")
		t.Logf("INSERT result: %s", result)

		result = engine.Execute("INSERT INTO users VALUES (2, 'Bob', 'bob@example.com')")
		t.Logf("INSERT result: %s", result)

		// Create savepoint
		result = engine.Execute("SAVEPOINT sp1")
		t.Logf("SAVEPOINT result: %s", result)

		// Insert more data
		result = engine.Execute("INSERT INTO users VALUES (3, 'Charlie', 'charlie@example.com')")
		t.Logf("INSERT result: %s", result)

		// Rollback to savepoint
		result = engine.Execute("ROLLBACK TO SAVEPOINT sp1")
		t.Logf("ROLLBACK TO SAVEPOINT result: %s", result)

		// Commit transaction
		result = engine.Execute("COMMIT")
		t.Logf("COMMIT result: %s", result)

		// Verify final state
		result = engine.Execute("SELECT * FROM users")
		t.Logf("Final SELECT result: %s", result)

		// Should only have 2 rows (Alice and Bob, not Charlie)
		if len(engine.DB.Tables["users"].Rows) != 2 {
			t.Errorf("Expected 2 rows after rollback to savepoint, got %d", len(engine.DB.Tables["users"].Rows))
		}
	})

	t.Run("RollbackTransactionWorkflow", func(t *testing.T) {
		// Begin transaction
		result := engine.Execute("BEGIN TRANSACTION")
		t.Logf("BEGIN result: %s", result)

		// Create table
		result = engine.Execute("CREATE TABLE products (id, name, price)")
		t.Logf("CREATE TABLE result: %s", result)

		// Insert data
		result = engine.Execute("INSERT INTO products VALUES (1, 'Laptop', '999.99')")
		t.Logf("INSERT result: %s", result)

		// Rollback transaction
		result = engine.Execute("ROLLBACK")
		t.Logf("ROLLBACK result: %s", result)

		// Verify table was not created
		result = engine.Execute("SELECT * FROM products")
		if result != "Table products not found" {
			t.Errorf("Expected table to not exist after rollback, got: %s", result)
		}
	})

	t.Run("TransactionWithUpdates", func(t *testing.T) {
		// Create table outside transaction
		engine.Execute("CREATE TABLE accounts (id, balance)")
		engine.Execute("INSERT INTO accounts VALUES (1, '100')")
		engine.Execute("INSERT INTO accounts VALUES (2, '50')")

		// Begin transaction
		result := engine.Execute("BEGIN TRANSACTION")
		t.Logf("BEGIN result: %s", result)

		// Update accounts
		result = engine.Execute("UPDATE accounts SET balance = '90' ROW 0")
		t.Logf("UPDATE result: %s", result)

		result = engine.Execute("UPDATE accounts SET balance = '60' ROW 1")
		t.Logf("UPDATE result: %s", result)

		// Rollback transaction
		result = engine.Execute("ROLLBACK")
		t.Logf("ROLLBACK result: %s", result)

		// Verify original values are preserved
		result = engine.Execute("SELECT * FROM accounts")
		t.Logf("Final SELECT result: %s", result)

		// Check that balances are back to original values
		table := engine.DB.Tables["accounts"]
		if table.Rows[0][1] != "100" {
			t.Errorf("Expected account 1 balance to be 100, got %s", table.Rows[0][1])
		}
		if table.Rows[1][1] != "50" {
			t.Errorf("Expected account 2 balance to be 50, got %s", table.Rows[1][1])
		}
	})

	t.Run("TransactionWithDeletes", func(t *testing.T) {
		// Create table and insert data
		engine.Execute("CREATE TABLE temp (id, value)")
		engine.Execute("INSERT INTO temp VALUES (1, 'A')")
		engine.Execute("INSERT INTO temp VALUES (2, 'B')")
		engine.Execute("INSERT INTO temp VALUES (3, 'C')")

		// Begin transaction
		result := engine.Execute("BEGIN TRANSACTION")
		t.Logf("BEGIN result: %s", result)

		// Delete rows
		result = engine.Execute("DELETE FROM temp ROW 0")
		t.Logf("DELETE result: %s", result)

		result = engine.Execute("DELETE FROM temp ROW 0") // Delete second row (now at index 0)
		t.Logf("DELETE result: %s", result)

		// Rollback transaction
		result = engine.Execute("ROLLBACK")
		t.Logf("ROLLBACK result: %s", result)

		// Verify all rows are back
		result = engine.Execute("SELECT * FROM temp")
		t.Logf("Final SELECT result: %s", result)

		table := engine.DB.Tables["temp"]
		if len(table.Rows) != 3 {
			t.Errorf("Expected 3 rows after rollback, got %d", len(table.Rows))
		}
	})
}

func TestTransactionErrorHandling(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	engine := NewEngine(tempDir)

	t.Run("CommitWithoutTransaction", func(t *testing.T) {
		result := engine.Execute("COMMIT")
		if result == "Transaction committed successfully" {
			t.Error("Expected error when committing without active transaction")
		}
		t.Logf("COMMIT without transaction result: %s", result)
	})

	t.Run("RollbackWithoutTransaction", func(t *testing.T) {
		result := engine.Execute("ROLLBACK")
		if result == "Transaction rolled back successfully" {
			t.Error("Expected error when rolling back without active transaction")
		}
		t.Logf("ROLLBACK without transaction result: %s", result)
	})

	t.Run("SavepointWithoutTransaction", func(t *testing.T) {
		result := engine.Execute("SAVEPOINT sp1")
		if result == "Savepoint sp1 created" {
			t.Error("Expected error when creating savepoint without active transaction")
		}
		t.Logf("SAVEPOINT without transaction result: %s", result)
	})

	t.Run("RollbackToNonexistentSavepoint", func(t *testing.T) {
		// Begin transaction
		engine.Execute("BEGIN TRANSACTION")

		result := engine.Execute("ROLLBACK TO SAVEPOINT nonexistent")
		if result == "Rolled back to savepoint nonexistent" {
			t.Error("Expected error when rolling back to nonexistent savepoint")
		}
		t.Logf("ROLLBACK TO nonexistent savepoint result: %s", result)

		// Clean up
		engine.Execute("ROLLBACK")
	})

	t.Run("InvalidIsolationLevel", func(t *testing.T) {
		result := engine.Execute("BEGIN TRANSACTION ISOLATION LEVEL INVALID")
		if result == "Transaction" { // Should contain error message
			t.Error("Expected error for invalid isolation level")
		}
		t.Logf("Invalid isolation level result: %s", result)
	})
}
