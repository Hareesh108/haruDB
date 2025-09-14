// internal/storage/transaction_test.go
package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTransactionManager(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	tm := db.TransactionManager

	t.Run("BeginTransaction", func(t *testing.T) {
		tx, err := tm.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if tx.State != TransactionActive {
			t.Errorf("Expected transaction state to be Active, got %d", tx.State)
		}

		if tx.IsolationLevel != ReadCommitted {
			t.Errorf("Expected isolation level to be ReadCommitted, got %d", tx.IsolationLevel)
		}

		if tx.ID == "" {
			t.Error("Expected transaction ID to be set")
		}
	})

	t.Run("CommitTransaction", func(t *testing.T) {
		// Create a table first
		db.CreateTable("test_table", []string{"id", "name"})

		tx, err := tm.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Add some operations
		err = tm.AddOperation(tx.ID, WAL_INSERT, "test_table", map[string]interface{}{
			"values": []interface{}{"1", "Alice"},
		})
		if err != nil {
			t.Fatalf("Failed to add operation: %v", err)
		}

		err = tm.AddOperation(tx.ID, WAL_INSERT, "test_table", map[string]interface{}{
			"values": []interface{}{"2", "Bob"},
		})
		if err != nil {
			t.Fatalf("Failed to add operation: %v", err)
		}

		// Commit transaction
		err = tm.CommitTransaction(tx.ID)
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify operations were applied
		table := db.Tables["test_table"]
		if len(table.Rows) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(table.Rows))
		}

		// Verify transaction is cleaned up
		_, exists := tm.transactions[tx.ID]
		if exists {
			t.Error("Transaction should be cleaned up after commit")
		}
	})

	t.Run("RollbackTransaction", func(t *testing.T) {
		// Create a table first
		db.CreateTable("test_table2", []string{"id", "name"})

		tx, err := tm.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Add some operations
		err = tm.AddOperation(tx.ID, WAL_INSERT, "test_table2", map[string]interface{}{
			"values": []interface{}{"1", "Alice"},
		})
		if err != nil {
			t.Fatalf("Failed to add operation: %v", err)
		}

		// Rollback transaction
		err = tm.RollbackTransaction(tx.ID)
		if err != nil {
			t.Fatalf("Failed to rollback transaction: %v", err)
		}

		// Verify operations were not applied
		table := db.Tables["test_table2"]
		if len(table.Rows) != 0 {
			t.Errorf("Expected 0 rows after rollback, got %d", len(table.Rows))
		}

		// Verify transaction is cleaned up
		_, exists := tm.transactions[tx.ID]
		if exists {
			t.Error("Transaction should be cleaned up after rollback")
		}
	})

	t.Run("Savepoints", func(t *testing.T) {
		// Create a table first
		db.CreateTable("test_table3", []string{"id", "name"})

		tx, err := tm.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Add first operation
		err = tm.AddOperation(tx.ID, WAL_INSERT, "test_table3", map[string]interface{}{
			"values": []interface{}{"1", "Alice"},
		})
		if err != nil {
			t.Fatalf("Failed to add operation: %v", err)
		}

		// Create savepoint
		err = tm.CreateSavepoint(tx.ID, "sp1")
		if err != nil {
			t.Fatalf("Failed to create savepoint: %v", err)
		}

		// Add more operations
		err = tm.AddOperation(tx.ID, WAL_INSERT, "test_table3", map[string]interface{}{
			"values": []interface{}{"2", "Bob"},
		})
		if err != nil {
			t.Fatalf("Failed to add operation: %v", err)
		}

		// Rollback to savepoint
		err = tm.RollbackToSavepoint(tx.ID, "sp1")
		if err != nil {
			t.Fatalf("Failed to rollback to savepoint: %v", err)
		}

		// Verify only operations before savepoint remain
		if len(tx.Operations) != 1 {
			t.Errorf("Expected 1 operation after rollback to savepoint, got %d", len(tx.Operations))
		}

		// Commit transaction
		err = tm.CommitTransaction(tx.ID)
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify only first operation was applied
		table := db.Tables["test_table3"]
		if len(table.Rows) != 1 {
			t.Errorf("Expected 1 row after rollback to savepoint, got %d", len(table.Rows))
		}
	})

	t.Run("TransactionNotFound", func(t *testing.T) {
		err := tm.CommitTransaction("nonexistent_tx")
		if err == nil {
			t.Error("Expected error for nonexistent transaction")
		}

		err = tm.RollbackTransaction("nonexistent_tx")
		if err == nil {
			t.Error("Expected error for nonexistent transaction")
		}
	})

	t.Run("InvalidTransactionState", func(t *testing.T) {
		tx, err := tm.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Commit transaction
		err = tm.CommitTransaction(tx.ID)
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Try to commit again
		err = tm.CommitTransaction(tx.ID)
		if err == nil {
			t.Error("Expected error when committing already committed transaction")
		}
	})
}

func TestTransactionACID(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_acid_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)

	t.Run("Atomicity", func(t *testing.T) {
		// Create table
		db.CreateTable("accounts", []string{"id", "balance"})
		db.Insert("accounts", []string{"1", "100"})
		db.Insert("accounts", []string{"2", "50"})

		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Transfer money (debit account 1, credit account 2)
		db.UpdateTx("accounts", 0, []string{"1", "90"}) // Debit 10
		db.UpdateTx("accounts", 1, []string{"2", "60"}) // Credit 10

		// Rollback transaction
		err = db.RollbackTransaction()
		if err != nil {
			t.Fatalf("Failed to rollback transaction: %v", err)
		}

		// Verify no changes were applied (atomicity)
		table := db.Tables["accounts"]
		if table.Rows[0][1] != "100" {
			t.Errorf("Expected account 1 balance to be 100, got %s", table.Rows[0][1])
		}
		if table.Rows[1][1] != "50" {
			t.Errorf("Expected account 2 balance to be 50, got %s", table.Rows[1][1])
		}
	})

	t.Run("Consistency", func(t *testing.T) {
		// Create table with constraints
		db.CreateTable("products", []string{"id", "name", "price"})

		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert valid data
		db.InsertTx("products", []string{"1", "Laptop", "999.99"})
		db.InsertTx("products", []string{"2", "Mouse", "29.99"})

		// Commit transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify data is consistent
		table := db.Tables["products"]
		if len(table.Rows) != 2 {
			t.Errorf("Expected 2 products, got %d", len(table.Rows))
		}

		// Verify data integrity
		for _, row := range table.Rows {
			if len(row) != 3 {
				t.Errorf("Expected 3 columns per row, got %d", len(row))
			}
		}
	})

	t.Run("Isolation", func(t *testing.T) {
		// Create table
		db.CreateTable("inventory", []string{"id", "quantity"})
		db.Insert("inventory", []string{"1", "10"})

		// Begin first transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction 1: %v", err)
		}

		// Update in first transaction
		db.UpdateTx("inventory", 0, []string{"1", "8"})

		// Begin second transaction
		_, err = db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction 2: %v", err)
		}

		// Read in second transaction (should see original value due to isolation)
		table := db.Tables["inventory"]
		if table.Rows[0][1] != "10" {
			t.Errorf("Expected to see original value 10 in second transaction, got %s", table.Rows[0][1])
		}

		// Commit first transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit transaction 1: %v", err)
		}

		// Commit second transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit transaction 2: %v", err)
		}

		// Verify final state
		table = db.Tables["inventory"]
		if table.Rows[0][1] != "8" {
			t.Errorf("Expected final value to be 8, got %s", table.Rows[0][1])
		}
	})

	t.Run("Durability", func(t *testing.T) {
		// Create table
		db.CreateTable("logs", []string{"id", "message"})

		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert data
		db.InsertTx("logs", []string{"1", "Transaction log entry"})

		// Commit transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify data is persisted
		table := db.Tables["logs"]
		if len(table.Rows) != 1 {
			t.Errorf("Expected 1 log entry, got %d", len(table.Rows))
		}

		// Verify WAL was written
		walPath := filepath.Join(tempDir, "wal.log")
		if _, err := os.Stat(walPath); os.IsNotExist(err) {
			t.Error("WAL file should exist after transaction")
		}
	})
}

func TestTransactionEdgeCases(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_edge_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)

	t.Run("NestedTransactions", func(t *testing.T) {
		// Create table
		db.CreateTable("test", []string{"id", "value"})

		// Begin first transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction 1: %v", err)
		}

		// Try to begin second transaction (should fail or be handled gracefully)
		_, err = db.BeginTransaction(ReadCommitted)
		if err != nil {
			// This is expected behavior - only one transaction at a time
			t.Logf("Expected error for nested transaction: %v", err)
		} else {
			// If it succeeds, clean up
			db.RollbackTransaction()
		}

		// Clean up first transaction
		db.RollbackTransaction()
	})

	t.Run("CommitWithoutTransaction", func(t *testing.T) {
		err := db.CommitTransaction()
		if err == nil {
			t.Error("Expected error when committing without active transaction")
		}
	})

	t.Run("RollbackWithoutTransaction", func(t *testing.T) {
		err := db.RollbackTransaction()
		if err == nil {
			t.Error("Expected error when rolling back without active transaction")
		}
	})

	t.Run("SavepointWithoutTransaction", func(t *testing.T) {
		err := db.CreateSavepoint("sp1")
		if err == nil {
			t.Error("Expected error when creating savepoint without active transaction")
		}
	})

	t.Run("RollbackToNonexistentSavepoint", func(t *testing.T) {
		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Try to rollback to nonexistent savepoint
		err = db.RollbackToSavepoint("nonexistent")
		if err == nil {
			t.Error("Expected error when rolling back to nonexistent savepoint")
		}

		// Clean up
		db.RollbackTransaction()
	})

	t.Run("LargeTransaction", func(t *testing.T) {
		// Create table
		db.CreateTable("large_test", []string{"id", "data"})

		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Add many operations
		for i := 0; i < 1000; i++ {
			db.InsertTx("large_test", []string{string(rune(i)), "data"})
		}

		// Commit transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit large transaction: %v", err)
		}

		// Verify all operations were applied
		table := db.Tables["large_test"]
		if len(table.Rows) != 1000 {
			t.Errorf("Expected 1000 rows, got %d", len(table.Rows))
		}
	})

	t.Run("TransactionTimeout", func(t *testing.T) {
		// Create table
		db.CreateTable("timeout_test", []string{"id", "value"})

		// Begin transaction
		_, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Simulate long-running transaction
		time.Sleep(100 * time.Millisecond)

		// Add operation
		db.InsertTx("timeout_test", []string{"1", "test"})

		// Commit transaction
		err = db.CommitTransaction()
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}

		// Verify operation was applied
		table := db.Tables["timeout_test"]
		if len(table.Rows) != 1 {
			t.Errorf("Expected 1 row, got %d", len(table.Rows))
		}
	})
}

func TestTransactionIsolationLevels(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "harudb_isolation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)

	t.Run("ReadUncommitted", func(t *testing.T) {
		// Create table
		db.CreateTable("isolation_test", []string{"id", "value"})
		db.Insert("isolation_test", []string{"1", "original"})

		// Begin transaction with ReadUncommitted
		tx, err := db.BeginTransaction(ReadUncommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if tx.IsolationLevel != ReadUncommitted {
			t.Errorf("Expected ReadUncommitted isolation level, got %d", tx.IsolationLevel)
		}

		// Clean up
		db.RollbackTransaction()
	})

	t.Run("ReadCommitted", func(t *testing.T) {
		// Begin transaction with ReadCommitted
		tx, err := db.BeginTransaction(ReadCommitted)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if tx.IsolationLevel != ReadCommitted {
			t.Errorf("Expected ReadCommitted isolation level, got %d", tx.IsolationLevel)
		}

		// Clean up
		db.RollbackTransaction()
	})

	t.Run("RepeatableRead", func(t *testing.T) {
		// Begin transaction with RepeatableRead
		tx, err := db.BeginTransaction(RepeatableRead)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if tx.IsolationLevel != RepeatableRead {
			t.Errorf("Expected RepeatableRead isolation level, got %d", tx.IsolationLevel)
		}

		// Clean up
		db.RollbackTransaction()
	})

	t.Run("Serializable", func(t *testing.T) {
		// Begin transaction with Serializable
		tx, err := db.BeginTransaction(Serializable)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if tx.IsolationLevel != Serializable {
			t.Errorf("Expected Serializable isolation level, got %d", tx.IsolationLevel)
		}

		// Clean up
		db.RollbackTransaction()
	})
}
