package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	// JWT library for secure token-based authentication
	// Enterprise-grade implementation of RFC 7519 JSON Web Token standard
	"github.com/golang-jwt/jwt/v5"

	// Official Go cryptography library for secure password hashing
	// Uses bcrypt algorithm for enterprise-grade password security
	"golang.org/x/crypto/bcrypt"

	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/config"
	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/storage"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
)

// JWTClaims extends the basic claims with JWT standard claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Service handles authentication business logic
type Service struct {
	userStore storage.UserStore
	config    *config.Config
}

// NewService creates a new authentication service
func NewService(userStore storage.UserStore, cfg *config.Config) *Service {
	return &Service{
		userStore: userStore,
		config:    cfg,
	}
}

// Register creates a new user account
func (s *Service) Register(req *RegisterRequest) (*LoginResponse, error) {
	// Check if user already exists
	if _, err := s.userStore.GetUserByEmail(req.Email); err == nil {
		return nil, ErrUserExists
	}

	if _, err := s.userStore.GetUserByUsername(req.Username); err == nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate user ID
	userID, err := s.generateID()
	if err != nil {
		return nil, err
	}

	// Create user
	user := &storage.User{
		ID:           userID,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := s.userStore.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate token
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		User:      s.userToUserInfo(user),
		ExpiresAt: expiresAt,
	}, nil
}

// Login authenticates a user and returns a token
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.userStore.GetUserByEmail(req.Email)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := s.verifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		User:      s.userToUserInfo(user),
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the user information
func (s *Service) ValidateToken(tokenString string) (*UserInfo, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	// Validate token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Get user from store to ensure it still exists and is active
	user, err := s.userStore.GetUserByID(claims.UserID)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInvalidToken
	}

	userInfo := s.userToUserInfo(user)
	return &userInfo, nil
}

// GetUserProfile returns user profile information
func (s *Service) GetUserProfile(userID string) (*UserInfo, error) {
	user, err := s.userStore.GetUserByID(userID)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userInfo := s.userToUserInfo(user)
	return &userInfo, nil
}

// hashPassword hashes a password using bcrypt
func (s *Service) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Auth.BCryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func (s *Service) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateToken generates a JWT token for a user
func (s *Service) generateToken(user *storage.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.Auth.TokenDuration)

	claims := &JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "login-app",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateID generates a random ID
func (s *Service) generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// userToUserInfo converts a User to UserInfo (removes sensitive data)
func (s *Service) userToUserInfo(user *storage.User) UserInfo {
	return UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
