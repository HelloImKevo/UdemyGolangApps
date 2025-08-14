package auth

import (
	"net/http"
	"strings"

	// Gin HTTP framework for REST API routing and middleware
	// Enterprise-grade web framework for secure HTTP request handling
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.service.Register(&req)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Registration failed"

		switch err {
		case ErrUserExists:
			status = http.StatusConflict
			message = "User already exists"
		case ErrInvalidCredentials:
			status = http.StatusBadRequest
			message = "Invalid credentials"
		}

		c.JSON(status, ErrorResponse{
			Error:   "registration_error",
			Message: message,
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Login failed"

		switch err {
		case ErrInvalidCredentials:
			status = http.StatusUnauthorized
			message = "Invalid email or password"
		case ErrUserNotFound:
			status = http.StatusUnauthorized
			message = "Invalid email or password"
		}

		c.JSON(status, ErrorResponse{
			Error:   "login_error",
			Message: message,
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	// In a JWT-based system, logout is typically handled client-side
	// by removing the token from storage
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// Profile returns the user's profile information
func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	profile, err := h.service.GetUserProfile(userID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to get user profile"

		switch err {
		case ErrUserNotFound:
			status = http.StatusNotFound
			message = "User not found"
		}

		c.JSON(status, ErrorResponse{
			Error:   "profile_error",
			Message: message,
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    profile,
	})
}

// Middleware creates authentication middleware
func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization header required",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid authorization header format",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		userInfo, err := h.service.ValidateToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Invalid token"

			switch err {
			case ErrTokenExpired:
				message = "Token expired"
			case ErrInvalidToken:
				message = "Invalid token"
			}

			c.JSON(status, ErrorResponse{
				Error:   "unauthorized",
				Message: message,
				Code:    status,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", userInfo.ID)
		c.Set("user_email", userInfo.Email)
		c.Set("user_username", userInfo.Username)
		c.Set("user_info", userInfo)

		c.Next()
	}
}
