package server

import (
	"net/http"

	// Gin HTTP framework for request handling and routing
	// Provides secure HTTP context, parameter binding, and response formatting
	"github.com/gin-gonic/gin"

	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/auth"
)

// Auth API handlers

func (s *Server) handleRegister(c *gin.Context) {
	handler := auth.NewHandler(s.authService)
	handler.Register(c)
}

func (s *Server) handleLogin(c *gin.Context) {
	handler := auth.NewHandler(s.authService)
	handler.Login(c)
}

func (s *Server) handleLogout(c *gin.Context) {
	handler := auth.NewHandler(s.authService)
	handler.Logout(c)
}

func (s *Server) handleProfile(c *gin.Context) {
	handler := auth.NewHandler(s.authService)
	handler.Profile(c)
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	handler := auth.NewHandler(s.authService)
	return handler.Middleware()
}

// Web page handlers

func (s *Server) handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcome to Login App",
	})
}

func (s *Server) handleLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func (s *Server) handleRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (s *Server) handleDashboard(c *gin.Context) {
	userInfo, exists := c.Get("user_info")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard",
		"user":  userInfo,
	})
}
