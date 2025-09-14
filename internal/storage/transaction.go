// internal/storage/transaction.go
package storage

import (
	"fmt"
	"sync"
	"time"
)

// TransactionState represents the state of a transaction
type TransactionState int

const (
	TransactionActive TransactionState = iota
	TransactionCommitted
	TransactionRolledBack
	TransactionAborted
)

// IsolationLevel represents the transaction isolation level
type IsolationLevel int

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
	RepeatableRead
	Serializable
)

// Transaction represents a database transaction
type Transaction struct {
	ID             string
	State          TransactionState
	IsolationLevel IsolationLevel
	StartTime      time.Time
	EndTime        time.Time
	Operations     []TransactionOperation
	Savepoints     map[string]int // savepoint name -> operation index
	mu             sync.RWMutex
}

// TransactionOperation represents a single operation within a transaction
type TransactionOperation struct {
	Type      WALEntryType
	TableName string
	Data      interface{}
	Timestamp time.Time
}

// TransactionManager manages all active transactions
type TransactionManager struct {
	transactions map[string]*Transaction
	nextID       int64
	mu           sync.RWMutex
	db           *Database
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *Database) *TransactionManager {
	return &TransactionManager{
		transactions: make(map[string]*Transaction),
		nextID:       1,
		db:           db,
	}
}

// BeginTransaction starts a new transaction
func (tm *TransactionManager) BeginTransaction(isolationLevel IsolationLevel) (*Transaction, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	txID := fmt.Sprintf("tx_%d_%d", time.Now().UnixNano(), tm.nextID)
	tm.nextID++

	tx := &Transaction{
		ID:             txID,
		State:          TransactionActive,
		IsolationLevel: isolationLevel,
		StartTime:      time.Now(),
		Operations:     make([]TransactionOperation, 0),
		Savepoints:     make(map[string]int),
	}

	tm.transactions[txID] = tx

	// Log transaction begin to WAL
	if tm.db.WAL != nil {
		data := map[string]interface{}{
			"isolation_level": int(isolationLevel),
		}
		if err := tm.db.WAL.WriteEntry(WAL_BEGIN_TRANSACTION, "", data); err != nil {
			return nil, fmt.Errorf("failed to write transaction begin to WAL: %w", err)
		}
	}

	return tx, nil
}

// GetTransaction retrieves a transaction by ID
func (tm *TransactionManager) GetTransaction(txID string) (*Transaction, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tx, exists := tm.transactions[txID]
	return tx, exists
}

// CommitTransaction commits a transaction
func (tm *TransactionManager) CommitTransaction(txID string) error {
	fmt.Printf("[COMMIT] start txID=%s", txID)

	// 1️⃣ Grab the transaction safely
	fmt.Printf("[COMMIT] locking tm.mu to fetch tx")
	tm.mu.Lock()
	tx, exists := tm.transactions[txID]
	if !exists {
		tm.mu.Unlock()
		fmt.Printf("[COMMIT] tx %s not found", txID)
		return fmt.Errorf("transaction %s not found", txID)
	}
	fmt.Printf("[COMMIT] tx %s fetched, releasing tm.mu", txID)
	tm.mu.Unlock()

	// 2️⃣ Lock the transaction itself
	fmt.Printf("[COMMIT] locking tx.mu")
	tx.mu.Lock()
	defer func() {
		fmt.Printf("[COMMIT] unlocked tx.mu for %s", txID)
		tx.mu.Unlock()
	}()

	if tx.State != TransactionActive {
		fmt.Printf("[COMMIT] tx %s not active (state=%d)", txID, tx.State)
		return fmt.Errorf("transaction %s is not active (state: %d)", txID, tx.State)
	}
	fmt.Printf("[COMMIT] tx %s is active with %d ops", txID, len(tx.Operations))

	// 3️⃣ Apply operations atomically
	for i, op := range tx.Operations {
		fmt.Printf("[COMMIT] applying op %d: %+v", i, op)
		if err := tm.applyOperation(op); err != nil {
			fmt.Printf("[COMMIT] FAILED op %d: %v — rolling back", i, err)
			tm.rollbackTransactionUnsafe(tx)
			return fmt.Errorf("failed to apply operation %d: %w", i, err)
		}
		fmt.Printf("[COMMIT] op %d applied successfully", i)
	}

	// 4️⃣ Mark committed
	tx.State = TransactionCommitted
	tx.EndTime = time.Now()
	fmt.Printf("[COMMIT] tx %s marked committed", txID)

	// 5️⃣ Clean up safely
	fmt.Printf("[COMMIT] locking tm.mu for cleanup")
	tm.mu.Lock()
	delete(tm.transactions, txID)
	tm.mu.Unlock()
	fmt.Printf("[COMMIT] tx %s removed from manager", txID)

	fmt.Printf("[COMMIT] completed successfully for tx %s", txID)
	return nil
}

// RollbackTransaction rolls back a transaction
func (tm *TransactionManager) RollbackTransaction(txID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tx, exists := tm.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	return tm.rollbackTransactionUnsafe(tx)
}

// rollbackTransactionUnsafe performs rollback without acquiring the manager lock
func (tm *TransactionManager) rollbackTransactionUnsafe(tx *Transaction) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != TransactionActive {
		return fmt.Errorf("transaction %s is not active (state: %d)", tx.ID, tx.State)
	}

	// Mark transaction as rolled back
	tx.State = TransactionRolledBack
	tx.EndTime = time.Now()

	// Log transaction rollback to WAL
	if tm.db.WAL != nil {
		data := map[string]interface{}{
			"transaction_id": tx.ID,
		}
		if err := tm.db.WAL.WriteEntry(WAL_ROLLBACK_TRANSACTION, "", data); err != nil {
			return fmt.Errorf("failed to write transaction rollback to WAL: %w", err)
		}
	}

	// Clean up transaction
	delete(tm.transactions, tx.ID)

	return nil
}

// CreateSavepoint creates a savepoint within a transaction
func (tm *TransactionManager) CreateSavepoint(txID, savepointName string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tx, exists := tm.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != TransactionActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	// Record savepoint at current operation count
	tx.Savepoints[savepointName] = len(tx.Operations)

	// Log savepoint creation to WAL
	if tm.db.WAL != nil {
		data := map[string]interface{}{
			"transaction_id":  txID,
			"savepoint_name":  savepointName,
			"operation_index": len(tx.Operations),
		}
		if err := tm.db.WAL.WriteEntry(WAL_SAVEPOINT, "", data); err != nil {
			return fmt.Errorf("failed to write savepoint to WAL: %w", err)
		}
	}

	return nil
}

// RollbackToSavepoint rolls back to a specific savepoint
func (tm *TransactionManager) RollbackToSavepoint(txID, savepointName string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tx, exists := tm.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != TransactionActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	operationIndex, exists := tx.Savepoints[savepointName]
	if !exists {
		return fmt.Errorf("savepoint %s not found", savepointName)
	}

	// Truncate operations to the savepoint
	tx.Operations = tx.Operations[:operationIndex]

	// Log rollback to savepoint to WAL
	if tm.db.WAL != nil {
		data := map[string]interface{}{
			"transaction_id":  txID,
			"savepoint_name":  savepointName,
			"operation_index": operationIndex,
		}
		if err := tm.db.WAL.WriteEntry(WAL_ROLLBACK_TO_SAVEPOINT, "", data); err != nil {
			return fmt.Errorf("failed to write rollback to savepoint to WAL: %w", err)
		}
	}

	return nil
}

// AddOperation adds an operation to a transaction
func (tm *TransactionManager) AddOperation(txID string, opType WALEntryType, tableName string, data interface{}) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tx, exists := tm.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != TransactionActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	// ✅ Normalize UPDATE payload so applyOperation gets the types it expects.
	if opType == WAL_UPDATE {
		if m, ok := data.(map[string]interface{}); ok {
			// ensure row_index is float64
			if ri, ok := m["row_index"].(int); ok {
				m["row_index"] = float64(ri)
			}
			// ensure values is []interface{}
			if vals, ok := m["values"].([]string); ok {
				intfVals := make([]interface{}, len(vals))
				for i, v := range vals {
					intfVals[i] = v
				}
				m["values"] = intfVals
			}
			data = m
		}
	}

	op := TransactionOperation{
		Type:      opType,
		TableName: tableName,
		Data:      data,
		Timestamp: time.Now(),
	}
	tx.Operations = append(tx.Operations, op)

	return nil
}

// applyOperation applies a single transaction operation to the database
func (tm *TransactionManager) applyOperation(op TransactionOperation) error {
	switch op.Type {
	case WAL_CREATE_TABLE:
		if data, ok := op.Data.(map[string]interface{}); ok {
			if columns, ok := data["columns"].([]interface{}); ok {
				colStrs := make([]string, len(columns))
				for i, col := range columns {
					colStrs[i] = col.(string)
				}
				return tm.applyCreateTable(op.TableName, colStrs)
			}
		}
		return fmt.Errorf("invalid CREATE TABLE operation data")

	case WAL_INSERT:
		if data, ok := op.Data.(map[string]interface{}); ok {
			if values, ok := data["values"].([]interface{}); ok {
				valStrs := make([]string, len(values))
				for i, val := range values {
					valStrs[i] = val.(string)
				}
				return tm.applyInsert(op.TableName, valStrs)
			}
		}
		return fmt.Errorf("invalid INSERT operation data")

	case WAL_UPDATE:
		if data, ok := op.Data.(map[string]interface{}); ok {
			if rowIndex, ok := data["row_index"].(float64); ok {
				if values, ok := data["values"].([]interface{}); ok {
					valStrs := make([]string, len(values))
					for i, val := range values {
						valStrs[i] = val.(string)
					}
					return tm.applyUpdate(op.TableName, int(rowIndex), valStrs)
				}
			}
		}
		return fmt.Errorf("invalid UPDATE operation data")

	case WAL_DELETE:
		if data, ok := op.Data.(map[string]interface{}); ok {
			if rowIndex, ok := data["row_index"].(float64); ok {
				return tm.applyDelete(op.TableName, int(rowIndex))
			}
		}
		return fmt.Errorf("invalid DELETE operation data")

	case WAL_DROP_TABLE:
		return tm.applyDropTable(op.TableName)

	default:
		return fmt.Errorf("unsupported operation type: %d", op.Type)
	}
}

// applyCreateTable applies CREATE TABLE operation
func (tm *TransactionManager) applyCreateTable(tableName string, columns []string) error {
	if _, exists := tm.db.Tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}

	tm.db.Tables[tableName] = &Table{
		Name:           tableName,
		Columns:        columns,
		Rows:           [][]string{},
		IndexedColumns: []string{},
		Indexes:        make(map[string]map[string][]int),
	}

	return tm.db.saveTable(tm.db.Tables[tableName])
}

// applyInsert applies INSERT operation
func (tm *TransactionManager) applyInsert(tableName string, values []string) error {
	table, exists := tm.db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	if len(values) != len(table.Columns) {
		return fmt.Errorf("column count mismatch")
	}

	table.Rows = append(table.Rows, values)
	tm.db.applyIndexesOnInsert(table, len(table.Rows)-1)

	return tm.db.saveTable(table)
}

// applyUpdate applies UPDATE operation
func (tm *TransactionManager) applyUpdate(tableName string, rowIndex int, values []string) error {
	table, exists := tm.db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return fmt.Errorf("row index %d out of bounds (table has %d rows)", rowIndex, len(table.Rows))
	}

	if len(values) != len(table.Columns) {
		return fmt.Errorf("column count mismatch: expected %d, got %d", len(table.Columns), len(values))
	}

	table.Rows[rowIndex] = values
	tm.db.rebuildAllIndexes(table)

	return tm.db.saveTable(table)
}

// applyDelete applies DELETE operation
func (tm *TransactionManager) applyDelete(tableName string, rowIndex int) error {
	table, exists := tm.db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return fmt.Errorf("row index out of bounds")
	}

	table.Rows = append(table.Rows[:rowIndex], table.Rows[rowIndex+1:]...)
	tm.db.rebuildAllIndexes(table)

	return tm.db.saveTable(table)
}

// applyDropTable applies DROP TABLE operation
func (tm *TransactionManager) applyDropTable(tableName string) error {
	if _, exists := tm.db.Tables[tableName]; !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	delete(tm.db.Tables, tableName)
	return nil
}

// GetActiveTransactions returns all active transactions
func (tm *TransactionManager) GetActiveTransactions() []*Transaction {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var active []*Transaction
	for _, tx := range tm.transactions {
		if tx.State == TransactionActive {
			active = append(active, tx)
		}
	}
	return active
}

// CleanupCompletedTransactions removes completed transactions
func (tm *TransactionManager) CleanupCompletedTransactions() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for txID, tx := range tm.transactions {
		if tx.State != TransactionActive {
			delete(tm.transactions, txID)
		}
	}
}
