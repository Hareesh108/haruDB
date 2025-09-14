// internal/storage/memory.go
package storage

import (
	"fmt"
	"os"
	"strings"
)

const (
	ErrTableNotFound = "Table %s not found"
	ErrWALCheckpoint = "Warning: Failed to write WAL checkpoint: %v\n"
)

type Table struct {
	Name    string
	Columns []string
	Rows    [][]string
	// IndexedColumns lists column names that are indexed
	IndexedColumns []string
	// Indexes maps column name -> value -> list of row indexes
	Indexes map[string]map[string][]int
}

type Database struct {
	DataDir           string
	Tables            map[string]*Table
	WAL               *WALManager
	TransactionManager *TransactionManager
	currentTransaction *Transaction
}

func NewDatabase(dataDir string) *Database {
	db := &Database{
		DataDir: dataDir,
		Tables:  make(map[string]*Table),
	}

	// Initialize WAL manager
	var err error
	db.WAL, err = NewWALManager(dataDir)
	if err != nil {
		// If WAL initialization fails, continue without WAL (degraded mode)
		fmt.Printf("Warning: Failed to initialize WAL: %v\n", err)
	}

	// Initialize Transaction Manager
	db.TransactionManager = NewTransactionManager(db)

	// Load any existing .harudb files first
	_ = db.loadTables()

	// Replay WAL entries if WAL is available (this will apply any changes not yet persisted)
	if db.WAL != nil {
		if err := db.WAL.ReplayWAL(db); err != nil {
			fmt.Printf("Warning: Failed to replay WAL: %v\n", err)
		}
	}

	return db
}

func (db *Database) CreateTable(name string, columns []string) string {
	name = strings.ToLower(name)
	if _, exists := db.Tables[name]; exists {
		return fmt.Sprintf("Table %s already exists", name)
	}

	// Write to WAL (Write Ahead Logs) first
	if db.WAL != nil {
		data := map[string]interface{}{
			"columns": columns,
		}
		if err := db.WAL.WriteEntry(WAL_CREATE_TABLE, name, data); err != nil {
			return fmt.Sprintf("Table %s created (warning: failed to write to WAL: %v)", name, err)
		}
	}

	// Apply changes to memory
	db.Tables[name] = &Table{Name: name, Columns: columns, Rows: [][]string{}, IndexedColumns: []string{}, Indexes: make(map[string]map[string][]int)}

	// Persist to disk
	if err := db.saveTable(db.Tables[name]); err != nil {
		return fmt.Sprintf("Table %s created (warning: failed to persist: %v)", name, err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return fmt.Sprintf("Table %s created", name)
}

func (db *Database) Insert(tableName string, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}
	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"values": values,
		}
		if err := db.WAL.WriteEntry(WAL_INSERT, tableName, data); err != nil {
			return fmt.Sprintf("1 row inserted (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows = append(table.Rows, values)
	// Maintain indexes for this row
	db.applyIndexesOnInsert(table, len(table.Rows)-1)

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("1 row inserted (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row inserted"
}

func (db *Database) SelectAll(tableName string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	result := strings.Join(table.Columns, " | ") + "\n"
	for _, row := range table.Rows {
		result += strings.Join(row, " | ") + "\n"
	}
	if len(table.Rows) == 0 {
		result += "(no rows)\n"
	}
	return result
}

// Update updates a row in the specified table
func (db *Database) Update(tableName string, rowIndex int, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
			"values":    values,
		}
		if err := db.WAL.WriteEntry(WAL_UPDATE, tableName, data); err != nil {
			return fmt.Sprintf("Row updated (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows[rowIndex] = values
	// Rebuild indexes as row positions and values may have changed
	db.rebuildAllIndexes(table)

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("Row updated (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row updated"
}

// Delete deletes a row from the specified table
func (db *Database) Delete(tableName string, rowIndex int) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	// Write to WAL first
	if db.WAL != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
		}
		if err := db.WAL.WriteEntry(WAL_DELETE, tableName, data); err != nil {
			return fmt.Sprintf("Row deleted (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	table.Rows = append(table.Rows[:rowIndex], table.Rows[rowIndex+1:]...)
	// Rebuild indexes as row positions shifted
	db.rebuildAllIndexes(table)

	// Persist to disk
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("Row deleted (warning: failed to persist: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return "1 row deleted"
}

// DropTable drops the specified table
func (db *Database) DropTable(tableName string) string {
	tableName = strings.ToLower(tableName)
	_, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// Write to WAL first
	if db.WAL != nil {
		if err := db.WAL.WriteEntry(WAL_DROP_TABLE, tableName, nil); err != nil {
			return fmt.Sprintf("Table dropped (warning: failed to write to WAL: %v)", err)
		}
	}

	// Apply changes to memory
	delete(db.Tables, tableName)

	// Remove table file from disk
	tablePath := db.tablePath(tableName)
	if err := os.Remove(tablePath); err != nil && !os.IsNotExist(err) {
		return fmt.Sprintf("Table dropped (warning: failed to remove table file: %v)", err)
	}

	// Write checkpoint to WAL
	if db.WAL != nil {
		if err := db.WAL.WriteCheckpoint(); err != nil {
			fmt.Printf(ErrWALCheckpoint, err)
		}
	}

	return fmt.Sprintf("Table %s dropped", tableName)
}

// CreateIndex creates an in-memory hash index on a given column and
// persists the indexed column metadata so indexes can be rebuilt on load.
func (db *Database) CreateIndex(tableName string, columnName string) string {
	tableName = strings.ToLower(tableName)
	columnName = strings.TrimSpace(columnName)

	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// Validate column exists
	colIdx := -1
	for i, c := range table.Columns {
		if c == columnName {
			colIdx = i
			break
		}
	}
	if colIdx == -1 {
		return fmt.Sprintf("Column %s not found", columnName)
	}

	// Initialize maps if needed
	if table.Indexes == nil {
		table.Indexes = make(map[string]map[string][]int)
	}

	if _, ok := table.Indexes[columnName]; !ok {
		table.Indexes[columnName] = make(map[string][]int)
	}

	// Add to IndexedColumns if not present
	found := false
	for _, ic := range table.IndexedColumns {
		if ic == columnName {
			found = true
			break
		}
	}
	if !found {
		table.IndexedColumns = append(table.IndexedColumns, columnName)
	}

	// Build index for this column
	db.buildIndexForColumn(table, columnName)

	// Persist table metadata so indexes can be rebuilt on restart
	if err := db.saveTable(table); err != nil {
		return fmt.Sprintf("Index created with warnings: failed to persist: %v", err)
	}

	return fmt.Sprintf("Index created on %s(%s)", tableName, columnName)
}

// SelectWhere returns rows where columnName == value. Uses index if available.
func (db *Database) SelectWhere(tableName, columnName, value string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// Header
	result := strings.Join(table.Columns, " | ") + "\n"

	// If index exists, use it
	if table.Indexes != nil {
		if idxMap, ok := table.Indexes[columnName]; ok {
			if rowIdxs, ok2 := idxMap[value]; ok2 {
				for _, ri := range rowIdxs {
					if ri >= 0 && ri < len(table.Rows) {
						result += strings.Join(table.Rows[ri], " | ") + "\n"
					}
				}
				if len(rowIdxs) == 0 {
					result += "(no rows)\n"
				}
				return result
			}
			// no matching value
			result += "(no rows)\n"
			return result
		}
	}

	// Fallback: full scan
	colIdx := -1
	for i, c := range table.Columns {
		if c == columnName {
			colIdx = i
			break
		}
	}
	if colIdx == -1 {
		return fmt.Sprintf("Column %s not found", columnName)
	}
	matched := 0
	for _, row := range table.Rows {
		if row[colIdx] == value {
			result += strings.Join(row, " | ") + "\n"
			matched++
		}
	}
	if matched == 0 {
		result += "(no rows)\n"
	}
	return result
}

// SelectWhereAdvanced returns rows matching complex WHERE conditions
func (db *Database) SelectWhereAdvanced(tableName string, whereExpr interface{}) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// Build column index map
	columnIndexes := make(map[string]int)
	for i, col := range table.Columns {
		columnIndexes[col] = i
	}

	// Header
	result := strings.Join(table.Columns, " | ") + "\n"

	// Evaluate each row against the WHERE expression
	matched := 0
	for _, row := range table.Rows {
		// Use reflection to call EvaluateExpression method
		if expr, ok := whereExpr.(interface {
			EvaluateExpression([]string, map[string]int) (bool, error)
		}); ok {
			match, err := expr.EvaluateExpression(row, columnIndexes)
			if err != nil {
				return fmt.Sprintf("Error evaluating WHERE condition: %v", err)
			}
			if match {
				result += strings.Join(row, " | ") + "\n"
				matched++
			}
		} else {
			return "Invalid WHERE expression type"
		}
	}

	if matched == 0 {
		result += "(no rows)\n"
	}
	return result
}

// buildIndexForColumn builds index for a specific column from scratch
func (db *Database) buildIndexForColumn(table *Table, columnName string) {
	if table.Indexes == nil {
		table.Indexes = make(map[string]map[string][]int)
	}
	idx, ok := table.Indexes[columnName]
	if !ok {
		idx = make(map[string][]int)
		table.Indexes[columnName] = idx
	} else {
		// reset existing
		for k := range idx {
			delete(idx, k)
		}
	}
	// find column index
	colIdx := -1
	for i, c := range table.Columns {
		if c == columnName {
			colIdx = i
			break
		}
	}
	if colIdx == -1 {
		return
	}
	for ri, row := range table.Rows {
		if colIdx < len(row) {
			val := row[colIdx]
			idx[val] = append(idx[val], ri)
		}
	}
}

// rebuildAllIndexes rebuilds all configured indexes for a table
func (db *Database) rebuildAllIndexes(table *Table) {
	if table == nil || len(table.IndexedColumns) == 0 {
		return
	}
	for _, col := range table.IndexedColumns {
		db.buildIndexForColumn(table, col)
	}
}

// applyIndexesOnInsert updates indexes for a newly inserted row at rowIndex
func (db *Database) applyIndexesOnInsert(table *Table, rowIndex int) {
	if table == nil || len(table.IndexedColumns) == 0 {
		return
	}
	row := table.Rows[rowIndex]
	for _, col := range table.IndexedColumns {
		// find column index
		colIdx := -1
		for i, c := range table.Columns {
			if c == col {
				colIdx = i
				break
			}
		}
		if colIdx == -1 || colIdx >= len(row) {
			continue
		}
		val := row[colIdx]
		if table.Indexes == nil {
			table.Indexes = make(map[string]map[string][]int)
		}
		if _, ok := table.Indexes[col]; !ok {
			table.Indexes[col] = make(map[string][]int)
		}
		table.Indexes[col][val] = append(table.Indexes[col][val], rowIndex)
	}
}

// Transaction-aware methods

// BeginTransaction starts a new transaction
func (db *Database) BeginTransaction(isolationLevel IsolationLevel) (*Transaction, error) {
	tx, err := db.TransactionManager.BeginTransaction(isolationLevel)
	if err != nil {
		return nil, err
	}
	db.currentTransaction = tx
	return tx, nil
}

// CommitTransaction commits the current transaction
func (db *Database) CommitTransaction() error {
	if db.currentTransaction == nil {
		return fmt.Errorf("no active transaction")
	}
	
	txID := db.currentTransaction.ID
	err := db.TransactionManager.CommitTransaction(txID)
	if err == nil {
		db.currentTransaction = nil
	}
	return err
}

// RollbackTransaction rolls back the current transaction
func (db *Database) RollbackTransaction() error {
	if db.currentTransaction == nil {
		return fmt.Errorf("no active transaction")
	}
	
	txID := db.currentTransaction.ID
	err := db.TransactionManager.RollbackTransaction(txID)
	if err == nil {
		db.currentTransaction = nil
	}
	return err
}

// CreateSavepoint creates a savepoint in the current transaction
func (db *Database) CreateSavepoint(name string) error {
	if db.currentTransaction == nil {
		return fmt.Errorf("no active transaction")
	}
	return db.TransactionManager.CreateSavepoint(db.currentTransaction.ID, name)
}

// RollbackToSavepoint rolls back to a savepoint in the current transaction
func (db *Database) RollbackToSavepoint(name string) error {
	if db.currentTransaction == nil {
		return fmt.Errorf("no active transaction")
	}
	return db.TransactionManager.RollbackToSavepoint(db.currentTransaction.ID, name)
}

// GetCurrentTransaction returns the current active transaction
func (db *Database) GetCurrentTransaction() *Transaction {
	return db.currentTransaction
}

// Transaction-aware versions of existing methods

// CreateTableTx creates a table within a transaction
func (db *Database) CreateTableTx(name string, columns []string) string {
	name = strings.ToLower(name)
	if _, exists := db.Tables[name]; exists {
		return fmt.Sprintf("Table %s already exists", name)
	}

	// If we're in a transaction, add operation to transaction
	if db.currentTransaction != nil {
		data := map[string]interface{}{
			"columns": columns,
		}
		if err := db.TransactionManager.AddOperation(db.currentTransaction.ID, WAL_CREATE_TABLE, name, data); err != nil {
			return fmt.Sprintf("Failed to add operation to transaction: %v", err)
		}
		return fmt.Sprintf("Table %s creation queued in transaction", name)
	}

	// Original non-transactional behavior
	return db.CreateTable(name, columns)
}

// InsertTx inserts a row within a transaction
func (db *Database) InsertTx(tableName string, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}
	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// If we're in a transaction, add operation to transaction
	if db.currentTransaction != nil {
		data := map[string]interface{}{
			"values": values,
		}
		if err := db.TransactionManager.AddOperation(db.currentTransaction.ID, WAL_INSERT, tableName, data); err != nil {
			return fmt.Sprintf("Failed to add operation to transaction: %v", err)
		}
		return "1 row insert queued in transaction"
	}

	// Original non-transactional behavior
	return db.Insert(tableName, values)
}

// UpdateTx updates a row within a transaction
func (db *Database) UpdateTx(tableName string, rowIndex int, values []string) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	if len(values) != len(table.Columns) {
		return "Column count does not match"
	}

	// If we're in a transaction, add operation to transaction
	if db.currentTransaction != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
			"values":    values,
		}
		if err := db.TransactionManager.AddOperation(db.currentTransaction.ID, WAL_UPDATE, tableName, data); err != nil {
			return fmt.Sprintf("Failed to add operation to transaction: %v", err)
		}
		return "1 row update queued in transaction"
	}

	// Original non-transactional behavior
	return db.Update(tableName, rowIndex, values)
}

// DeleteTx deletes a row within a transaction
func (db *Database) DeleteTx(tableName string, rowIndex int) string {
	tableName = strings.ToLower(tableName)
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	if rowIndex < 0 || rowIndex >= len(table.Rows) {
		return "Row index out of bounds"
	}

	// If we're in a transaction, add operation to transaction
	if db.currentTransaction != nil {
		data := map[string]interface{}{
			"row_index": rowIndex,
		}
		if err := db.TransactionManager.AddOperation(db.currentTransaction.ID, WAL_DELETE, tableName, data); err != nil {
			return fmt.Sprintf("Failed to add operation to transaction: %v", err)
		}
		return "1 row delete queued in transaction"
	}

	// Original non-transactional behavior
	return db.Delete(tableName, rowIndex)
}

// DropTableTx drops a table within a transaction
func (db *Database) DropTableTx(tableName string) string {
	tableName = strings.ToLower(tableName)
	_, exists := db.Tables[tableName]
	if !exists {
		return fmt.Sprintf(ErrTableNotFound, tableName)
	}

	// If we're in a transaction, add operation to transaction
	if db.currentTransaction != nil {
		if err := db.TransactionManager.AddOperation(db.currentTransaction.ID, WAL_DROP_TABLE, tableName, nil); err != nil {
			return fmt.Sprintf("Failed to add operation to transaction: %v", err)
		}
		return fmt.Sprintf("Table %s drop queued in transaction", tableName)
	}

	// Original non-transactional behavior
	return db.DropTable(tableName)
}
