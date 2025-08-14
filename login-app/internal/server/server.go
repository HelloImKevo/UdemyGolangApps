package server

import (
	"net/http"

	// Gin HTTP web framework for REST API and web page serving
	// Provides routing, middleware, input validation, and security features
	"github.com/gin-gonic/gin"

	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/auth"
	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/config"
	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	router      *gin.Engine
	authService *auth.Service
	config      *config.Config
}

// New creates a new server instance
func New(cfg *config.Config, userStore storage.UserStore) (*Server, error) {
	// Set Gin mode based on environment
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Create auth service
	authService := auth.NewService(userStore, cfg)

	server := &Server{
		router:      router,
		authService: authService,
		config:      cfg,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// Handler returns the HTTP handler
func (s *Server) Handler() http.Handler {
	return s.router
}

// setupMiddleware configures global middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware (conditional)
	if s.config.Log.Level == "debug" {
		s.router.Use(gin.Logger())
	}

	// CORS middleware (basic implementation)
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Security headers
	s.router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Load HTML templates
	s.router.LoadHTMLGlob("web/templates/*")

	// Health check
	s.router.GET("/health", s.healthCheck)

	// Static files
	s.router.Static("/static", "./web/static")

	// API routes
	api := s.router.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", s.handleRegister)
			authGroup.POST("/login", s.handleLogin)
			authGroup.POST("/logout", s.handleLogout)
			authGroup.GET("/profile", s.authMiddleware(), s.handleProfile)
		}
	}

	// Web routes (will serve HTML pages)
	s.router.GET("/", s.handleHome)
	s.router.GET("/login", s.handleLoginPage)
	s.router.GET("/register", s.handleRegisterPage)
	s.router.GET("/dashboard", s.authMiddleware(), s.handleDashboard)
} // healthCheck returns the service health status
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "login-app",
	})
}
