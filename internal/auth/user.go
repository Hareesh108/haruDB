// internal/auth/user.go
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UserRole represents the role of a user
type UserRole int

const (
	RoleAdmin UserRole = iota
	RoleUser
	RoleReadOnly
)

// User represents a database user
type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
	IsActive     bool      `json:"is_active"`
}

// Session represents an active user session
type Session struct {
	SessionID  string
	Username   string
	Role       UserRole
	CreatedAt  time.Time
	LastAccess time.Time
	IsActive   bool
}

// UserManager handles user authentication and management
type UserManager struct {
	users     map[string]*User
	sessions  map[string]*Session
	usersFile string
	mu        sync.RWMutex
}

// NewUserManager creates a new user manager
func NewUserManager(dataDir string) *UserManager {
	usersFile := filepath.Join(dataDir, "users.json")

	um := &UserManager{
		users:     make(map[string]*User),
		sessions:  make(map[string]*Session),
		usersFile: usersFile,
	}

	// Load existing users
	um.loadUsers()

	// Create default admin user if no users exist
	if len(um.users) == 0 {
		um.createDefaultAdmin()
	}

	return um
}

// createDefaultAdmin creates a default admin user
func (um *UserManager) createDefaultAdmin() {
	adminUser := &User{
		Username:     "admin",
		PasswordHash: um.hashPassword("admin123"),
		Role:         RoleAdmin,
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	um.users["admin"] = adminUser
	um.saveUsers()
}

// hashPassword hashes a password using SHA-256
func (um *UserManager) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// AuthenticateUser authenticates a user with username and password
func (um *UserManager) AuthenticateUser(username, password string) (*User, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()

	user, exists := um.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is disabled")
	}

	if user.PasswordHash != um.hashPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	// Update last login
	user.LastLogin = time.Now()
	um.saveUsers()

	return user, nil
}

// CreateSession creates a new session for a user
func (um *UserManager) CreateSession(user *User) (*Session, error) {
	um.mu.Lock()
	defer um.mu.Unlock()

	// Generate random session ID
	sessionID := um.generateSessionID()

	session := &Session{
		SessionID:  sessionID,
		Username:   user.Username,
		Role:       user.Role,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
		IsActive:   true,
	}

	um.sessions[sessionID] = session
	return session, nil
}

// ValidateSession validates a session ID
func (um *UserManager) ValidateSession(sessionID string) (*Session, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()

	session, exists := um.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("invalid session")
	}

	if !session.IsActive {
		return nil, fmt.Errorf("session expired")
	}

	// Update last access
	session.LastAccess = time.Now()

	return session, nil
}

// CreateUser creates a new user
func (um *UserManager) CreateUser(username, password string, role UserRole) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	if _, exists := um.users[username]; exists {
		return fmt.Errorf("user already exists")
	}

	user := &User{
		Username:     username,
		PasswordHash: um.hashPassword(password),
		Role:         role,
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	um.users[username] = user
	return um.saveUsers()
}

// DeleteUser deletes a user
func (um *UserManager) DeleteUser(username string) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	if _, exists := um.users[username]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(um.users, username)

	// Remove all sessions for this user
	for sessionID, session := range um.sessions {
		if session.Username == username {
			delete(um.sessions, sessionID)
		}
	}

	return um.saveUsers()
}

// ListUsers returns a list of all users
func (um *UserManager) ListUsers() []*User {
	um.mu.RLock()
	defer um.mu.RUnlock()

	users := make([]*User, 0, len(um.users))
	for _, user := range um.users {
		// Don't expose password hashes
		userCopy := *user
		userCopy.PasswordHash = "***"
		users = append(users, &userCopy)
	}

	return users
}

// generateSessionID generates a random session ID
func (um *UserManager) generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// loadUsers loads users from file
func (um *UserManager) loadUsers() error {
	if _, err := os.Stat(um.usersFile); os.IsNotExist(err) {
		return nil // File doesn't exist, start with empty users
	}

	data, err := os.ReadFile(um.usersFile)
	if err != nil {
		return fmt.Errorf("failed to read users file: %w", err)
	}

	var users map[string]*User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to unmarshal users: %w", err)
	}

	um.users = users
	return nil
}

// saveUsers saves users to file
func (um *UserManager) saveUsers() error {
	data, err := json.MarshalIndent(um.users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(um.usersFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

// LogoutSession logs out a session
func (um *UserManager) LogoutSession(sessionID string) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	if session, exists := um.sessions[sessionID]; exists {
		session.IsActive = false
		delete(um.sessions, sessionID)
	}

	return nil
}

// CleanupExpiredSessions removes expired sessions
func (um *UserManager) CleanupExpiredSessions() {
	um.mu.Lock()
	defer um.mu.Unlock()

	now := time.Now()
	for sessionID, session := range um.sessions {
		// Remove sessions older than 24 hours
		if now.Sub(session.LastAccess) > 24*time.Hour {
			delete(um.sessions, sessionID)
		}
	}
}

// UpdateUserPassword updates a user's password
func (um *UserManager) UpdateUserPassword(username, newPassword string) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	user, exists := um.users[username]
	if !exists {
		return fmt.Errorf("user not found")
	}

	user.PasswordHash = um.hashPassword(newPassword)
	return um.saveUsers()
}
