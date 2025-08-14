package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig `json:"server"`
	Auth   AuthConfig   `json:"auth"`
	Log    LogConfig    `json:"log"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// AuthConfig contains authentication-related configuration
type AuthConfig struct {
	JWTSecret      string        `json:"jwt_secret"`
	TokenDuration  time.Duration `json:"token_duration"`
	BCryptCost     int           `json:"bcrypt_cost"`
	SessionTimeout time.Duration `json:"session_timeout"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Load loads configuration based on the environment
func Load(environment string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Auth: AuthConfig{
			JWTSecret:      getEnv("JWT_SECRET", "your-256-bit-secret-key-here-make-sure-its-long-enough"),
			TokenDuration:  24 * time.Hour,
			BCryptCost:     getBcryptCost(),
			SessionTimeout: 24 * time.Hour,
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}

	// Environment-specific overrides
	switch environment {
	case "production":
		cfg.Auth.BCryptCost = 12 // Higher cost for production
		cfg.Log.Level = "warn"
		if cfg.Auth.JWTSecret == "your-256-bit-secret-key-here-make-sure-its-long-enough" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production environment")
		}
	case "development":
		cfg.Auth.BCryptCost = 8 // Lower cost for development
		cfg.Log.Level = "debug"
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBcryptCost gets the bcrypt cost from environment or returns default
func getBcryptCost() int {
	if cost := os.Getenv("BCRYPT_COST"); cost != "" {
		if c, err := strconv.Atoi(cost); err == nil && c >= 4 && c <= 31 {
			return c
		}
	}
	return 10 // Default bcrypt cost
}
