package storage

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never include in JSON
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
}

// UserStore defines the interface for user storage operations
type UserStore interface {
	// CreateUser creates a new user
	CreateUser(user *User) error

	// GetUserByID retrieves a user by ID
	GetUserByID(id string) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(email string) (*User, error)

	// GetUserByUsername retrieves a user by username
	GetUserByUsername(username string) (*User, error)

	// UpdateUser updates an existing user
	UpdateUser(user *User) error

	// DeleteUser deletes a user by ID
	DeleteUser(id string) error

	// ListUsers returns all users (for admin purposes)
	ListUsers() ([]*User, error)
}

// MemoryUserStore implements UserStore using in-memory storage
type MemoryUserStore struct {
	mu          sync.RWMutex
	users       map[string]*User
	emailIdx    map[string]string // email -> user_id mapping
	usernameIdx map[string]string // username -> user_id mapping
}

// NewMemoryUserStore creates a new in-memory user store
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users:       make(map[string]*User),
		emailIdx:    make(map[string]string),
		usernameIdx: make(map[string]string),
	}
}

// CreateUser creates a new user
func (s *MemoryUserStore) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user ID already exists
	if _, exists := s.users[user.ID]; exists {
		return ErrUserExists
	}

	// Check if email already exists
	if _, exists := s.emailIdx[user.Email]; exists {
		return ErrUserExists
	}

	// Check if username already exists
	if _, exists := s.usernameIdx[user.Username]; exists {
		return ErrUserExists
	}

	// Create user
	userCopy := *user
	userCopy.CreatedAt = time.Now()
	userCopy.UpdatedAt = time.Now()
	userCopy.IsActive = true

	s.users[user.ID] = &userCopy
	s.emailIdx[user.Email] = user.ID
	s.usernameIdx[user.Username] = user.ID

	return nil
}

// GetUserByID retrieves a user by ID
func (s *MemoryUserStore) GetUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Return a copy to prevent external modification
	userCopy := *user
	return &userCopy, nil
}

// GetUserByEmail retrieves a user by email
func (s *MemoryUserStore) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.emailIdx[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	user := s.users[userID]
	userCopy := *user
	return &userCopy, nil
}

// GetUserByUsername retrieves a user by username
func (s *MemoryUserStore) GetUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.usernameIdx[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	user := s.users[userID]
	userCopy := *user
	return &userCopy, nil
}

// UpdateUser updates an existing user
func (s *MemoryUserStore) UpdateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingUser, exists := s.users[user.ID]
	if !exists {
		return ErrUserNotFound
	}

	// Check if email changed and if new email already exists
	if user.Email != existingUser.Email {
		if _, exists := s.emailIdx[user.Email]; exists {
			return ErrUserExists
		}
		// Update email index
		delete(s.emailIdx, existingUser.Email)
		s.emailIdx[user.Email] = user.ID
	}

	// Check if username changed and if new username already exists
	if user.Username != existingUser.Username {
		if _, exists := s.usernameIdx[user.Username]; exists {
			return ErrUserExists
		}
		// Update username index
		delete(s.usernameIdx, existingUser.Username)
		s.usernameIdx[user.Username] = user.ID
	}

	// Update user
	userCopy := *user
	userCopy.UpdatedAt = time.Now()
	s.users[user.ID] = &userCopy

	return nil
}

// DeleteUser deletes a user by ID
func (s *MemoryUserStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return ErrUserNotFound
	}

	// Remove from all indexes
	delete(s.users, id)
	delete(s.emailIdx, user.Email)
	delete(s.usernameIdx, user.Username)

	return nil
}

// ListUsers returns all users
func (s *MemoryUserStore) ListUsers() ([]*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	return users, nil
}
