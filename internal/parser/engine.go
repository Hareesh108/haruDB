// internal/parser/engine.go
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Hareesh108/haruDB/internal/auth"
	"github.com/Hareesh108/haruDB/internal/storage"
)

const (
	ErrSyntaxError = "Syntax error"
)

type Engine struct {
	DB             *storage.Database
	UserManager    *auth.UserManager
	BackupManager  *storage.BackupManager
	CurrentSession *auth.Session
}

func NewEngine(dataDir string) *Engine {
	return &Engine{
		DB:            storage.NewDatabase(dataDir),
		UserManager:   auth.NewUserManager(dataDir),
		BackupManager: storage.NewBackupManager(dataDir),
	}
}

func (e *Engine) Execute(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimSuffix(input, ";") // remove trailing semicolon

	upper := strings.ToUpper(input)

	switch {
	case strings.HasPrefix(upper, "BEGIN"):
		// BEGIN TRANSACTION [ISOLATION LEVEL level]
		return e.handleBeginTransaction(input)

	case strings.HasPrefix(upper, "COMMIT"):
		// COMMIT [TRANSACTION]
		return e.handleCommitTransaction()

	case strings.HasPrefix(upper, "ROLLBACK"):
		// ROLLBACK [TRANSACTION] [TO SAVEPOINT name]
		return e.handleRollbackTransaction(input)

	case strings.HasPrefix(upper, "SAVEPOINT"):
		// SAVEPOINT name
		return e.handleSavepoint(input)

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
		return e.DB.CreateTableTx(tableName, columns)

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
		return e.DB.InsertTx(tableName, values)

	case strings.HasPrefix(upper, "SELECT * FROM"):
		// SELECT * FROM users [WHERE conditions]
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

		// Extract WHERE clause
		whereClause := strings.Join(parts[whereIdx+1:], " ")

		// Parse advanced WHERE clause
		whereExpr, err := ParseWhereClause(whereClause)
		if err != nil {
			return fmt.Sprintf("WHERE clause error: %v", err)
		}

		// Use advanced WHERE evaluation
		return e.DB.SelectWhereAdvanced(tableName, whereExpr)

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

		return e.DB.UpdateTx(tableName, rowIndex, newRow)

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

		return e.DB.DeleteTx(tableName, rowIndex)

	case strings.HasPrefix(upper, "DROP TABLE"):
		// DROP TABLE users
		parts := strings.Fields(input)
		if len(parts) < 3 {
			return "Syntax error: DROP TABLE table_name"
		}
		tableName := strings.ToLower(parts[2])
		return e.DB.DropTableTx(tableName)

	case strings.HasPrefix(upper, "LOGIN"):
		// LOGIN username password
		return e.handleLogin(input)

	case strings.HasPrefix(upper, "LOGOUT"):
		// LOGOUT
		return e.handleLogout()

	case strings.HasPrefix(upper, "CREATE USER"):
		// CREATE USER username password [role]
		return e.handleCreateUser(input)

	case strings.HasPrefix(upper, "DROP USER"):
		// DROP USER username
		return e.handleDropUser(input)

	case strings.HasPrefix(upper, "LIST USERS"):
		// LIST USERS
		return e.handleListUsers()

	case strings.HasPrefix(upper, "BACKUP"):
		// BACKUP [TO path] [DESCRIPTION description]
		return e.handleBackup(input)

	case strings.HasPrefix(upper, "RESTORE"):
		// RESTORE FROM path
		return e.handleRestore(input)

	case strings.HasPrefix(upper, "BACKUP INFO"):
		// BACKUP INFO path
		return e.handleBackupInfo(input)

	case strings.HasPrefix(upper, "LIST BACKUPS"):
		// LIST BACKUPS [directory]
		return e.handleListBackups(input)

	default:
		return "Unknown command"
	}
}

// Transaction handler methods

// handleBeginTransaction handles BEGIN TRANSACTION commands
func (e *Engine) handleBeginTransaction(input string) string {
	parts := strings.Fields(input)

	// Default isolation level
	isolationLevel := storage.ReadCommitted

	// Parse isolation level if specified
	if len(parts) >= 3 && strings.ToUpper(parts[1]) == "TRANSACTION" {
		if len(parts) >= 6 && strings.ToUpper(parts[2]) == "ISOLATION" &&
			strings.ToUpper(parts[3]) == "LEVEL" {
			switch strings.ToUpper(parts[4]) {
			case "READ":
				if len(parts) >= 6 && strings.ToUpper(parts[5]) == "UNCOMMITTED" {
					isolationLevel = storage.ReadUncommitted
				} else {
					isolationLevel = storage.ReadCommitted
				}
			case "REPEATABLE":
				if len(parts) >= 6 && strings.ToUpper(parts[5]) == "READ" {
					isolationLevel = storage.RepeatableRead
				}
			case "SERIALIZABLE":
				isolationLevel = storage.Serializable
			default:
				return "Invalid isolation level"
			}
		}
	}

	tx, err := e.DB.BeginTransaction(isolationLevel)
	if err != nil {
		return fmt.Sprintf("Failed to begin transaction: %v", err)
	}

	return fmt.Sprintf("Transaction %s started with isolation level %d", tx.ID, isolationLevel)
}

// handleCommitTransaction handles COMMIT commands
func (e *Engine) handleCommitTransaction() string {

	fmt.Printf("Hello")

	err := e.DB.CommitTransaction()

	fmt.Printf("commit err = %#v", err)

	if err != nil {
		return fmt.Sprintf("Failed to commit transaction: %v", err)
	}
	return "Transaction committed successfully"
}

// handleRollbackTransaction handles ROLLBACK commands
func (e *Engine) handleRollbackTransaction(input string) string {
	parts := strings.Fields(input)

	// Check for ROLLBACK TO SAVEPOINT
	if len(parts) >= 4 && strings.ToUpper(parts[1]) == "TO" &&
		strings.ToUpper(parts[2]) == "SAVEPOINT" {
		savepointName := parts[3]
		err := e.DB.RollbackToSavepoint(savepointName)
		if err != nil {
			return fmt.Sprintf("Failed to rollback to savepoint %s: %v", savepointName, err)
		}
		return fmt.Sprintf("Rolled back to savepoint %s", savepointName)
	}

	// Regular rollback
	err := e.DB.RollbackTransaction()
	if err != nil {
		return fmt.Sprintf("Failed to rollback transaction: %v", err)
	}
	return "Transaction rolled back successfully"
}

// handleSavepoint handles SAVEPOINT commands
func (e *Engine) handleSavepoint(input string) string {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "Syntax error: SAVEPOINT name"
	}

	savepointName := parts[1]
	err := e.DB.CreateSavepoint(savepointName)
	if err != nil {
		return fmt.Sprintf("Failed to create savepoint %s: %v", savepointName, err)
	}
	return fmt.Sprintf("Savepoint %s created", savepointName)
}

// Authentication handler methods

// handleLogin handles LOGIN commands
func (e *Engine) handleLogin(input string) string {
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return "Syntax error: LOGIN username password"
	}

	username := parts[1]
	password := parts[2]

	user, err := e.UserManager.AuthenticateUser(username, password)
	if err != nil {
		return fmt.Sprintf("Login failed: %v", err)
	}

	session, err := e.UserManager.CreateSession(user)
	if err != nil {
		return fmt.Sprintf("Failed to create session: %v", err)
	}

	e.CurrentSession = session
	return fmt.Sprintf("Login successful. Welcome, %s!", username)
}

// handleLogout handles LOGOUT commands
func (e *Engine) handleLogout() string {
	if e.CurrentSession == nil {
		return "No active session"
	}

	err := e.UserManager.LogoutSession(e.CurrentSession.SessionID)
	if err != nil {
		return fmt.Sprintf("Logout failed: %v", err)
	}

	e.CurrentSession = nil
	return "Logout successful"
}

// handleCreateUser handles CREATE USER commands
func (e *Engine) handleCreateUser(input string) string {
	if e.CurrentSession == nil || e.CurrentSession.Role != auth.RoleAdmin {
		return "Access denied: Admin privileges required"
	}

	parts := strings.Fields(input)
	if len(parts) < 4 {
		return "Syntax error: CREATE USER username password [role]"
	}

	username := parts[2]
	password := parts[3]
	role := auth.RoleUser

	// Parse role if specified
	if len(parts) >= 5 {
		switch strings.ToUpper(parts[4]) {
		case "ADMIN":
			role = auth.RoleAdmin
		case "USER":
			role = auth.RoleUser
		case "READONLY":
			role = auth.RoleReadOnly
		default:
			return "Invalid role. Use: ADMIN, USER, or READONLY"
		}
	}

	err := e.UserManager.CreateUser(username, password, role)
	if err != nil {
		return fmt.Sprintf("Failed to create user: %v", err)
	}

	return fmt.Sprintf("User %s created successfully", username)
}

// handleDropUser handles DROP USER commands
func (e *Engine) handleDropUser(input string) string {
	if e.CurrentSession == nil || e.CurrentSession.Role != auth.RoleAdmin {
		return "Access denied: Admin privileges required"
	}

	parts := strings.Fields(input)
	if len(parts) < 3 {
		return "Syntax error: DROP USER username"
	}

	username := parts[2]
	err := e.UserManager.DeleteUser(username)
	if err != nil {
		return fmt.Sprintf("Failed to delete user: %v", err)
	}

	return fmt.Sprintf("User %s deleted successfully", username)
}

// handleListUsers handles LIST USERS commands
func (e *Engine) handleListUsers() string {
	if e.CurrentSession == nil || e.CurrentSession.Role != auth.RoleAdmin {
		return "Access denied: Admin privileges required"
	}

	users := e.UserManager.ListUsers()
	if len(users) == 0 {
		return "No users found"
	}

	result := "Users:\n"
	for _, user := range users {
		roleStr := "USER"
		switch user.Role {
		case auth.RoleAdmin:
			roleStr = "ADMIN"
		case auth.RoleReadOnly:
			roleStr = "READONLY"
		}
		result += fmt.Sprintf("- %s (%s) - Created: %s, Last Login: %s\n",
			user.Username, roleStr, user.CreatedAt.Format("2006-01-02 15:04:05"),
			user.LastLogin.Format("2006-01-02 15:04:05"))
	}

	return result
}

// Backup handler methods

// handleBackup handles BACKUP commands
func (e *Engine) handleBackup(input string) string {
	if e.CurrentSession == nil || e.CurrentSession.Role == auth.RoleReadOnly {
		return "Access denied: Write privileges required"
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "Syntax error: BACKUP [TO path] [DESCRIPTION description]"
	}

	// Default backup path
	backupPath := fmt.Sprintf("./backups/harudb_backup_%s.backup", time.Now().Format("20060102_150405"))
	description := "Manual backup"

	// Parse optional parameters
	for i := 1; i < len(parts); i++ {
		if strings.ToUpper(parts[i]) == "TO" && i+1 < len(parts) {
			backupPath = parts[i+1]
			i++
		} else if strings.ToUpper(parts[i]) == "DESCRIPTION" && i+1 < len(parts) {
			description = parts[i+1]
			i++
		}
	}

	err := e.BackupManager.CreateBackup(backupPath, description)
	if err != nil {
		return fmt.Sprintf("Backup failed: %v", err)
	}

	return fmt.Sprintf("Backup created successfully: %s", backupPath)
}

// handleRestore handles RESTORE commands
func (e *Engine) handleRestore(input string) string {
	if e.CurrentSession == nil || e.CurrentSession.Role != auth.RoleAdmin {
		return "Access denied: Admin privileges required"
	}

	parts := strings.Fields(input)
	if len(parts) < 3 || strings.ToUpper(parts[1]) != "FROM" {
		return "Syntax error: RESTORE FROM path"
	}

	backupPath := parts[2]
	err := e.BackupManager.RestoreBackup(backupPath)
	if err != nil {
		return fmt.Sprintf("Restore failed: %v", err)
	}

	return fmt.Sprintf("Database restored successfully from: %s", backupPath)
}

// handleBackupInfo handles BACKUP INFO commands
func (e *Engine) handleBackupInfo(input string) string {
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return "Syntax error: BACKUP INFO path"
	}

	backupPath := parts[2]
	info, err := e.BackupManager.GetBackupInfo(backupPath)
	if err != nil {
		return fmt.Sprintf("Failed to get backup info: %v", err)
	}

	return fmt.Sprintf("Backup Info:\n"+
		"Timestamp: %s\n"+
		"Version: %s\n"+
		"Table Count: %d\n"+
		"Backup Size: %d bytes\n"+
		"Description: %s",
		info.Timestamp.Format("2006-01-02 15:04:05"),
		info.Version,
		info.TableCount,
		info.BackupSize,
		info.Description)
}

// handleListBackups handles LIST BACKUPS commands
func (e *Engine) handleListBackups(input string) string {
	parts := strings.Fields(input)
	backupDir := "./backups"

	if len(parts) >= 3 {
		backupDir = parts[2]
	}

	backups, err := e.BackupManager.ListBackups(backupDir)
	if err != nil {
		return fmt.Sprintf("Failed to list backups: %v", err)
	}

	if len(backups) == 0 {
		return "No backups found"
	}

	result := "Available backups:\n"
	for _, backup := range backups {
		result += fmt.Sprintf("- %s\n", backup)
	}

	return result
}
