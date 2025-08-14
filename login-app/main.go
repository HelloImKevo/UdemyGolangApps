package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/config"
	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/server"
	"github.com/HelloImKevo/UdemyGolangApps/login-app/internal/storage"
)

var (
	flagVersion = flag.Bool("version", false, "show version")
	flagPort    = flag.String("port", "8080", "port to listen on")
	flagEnv     = flag.String("env", "development", "environment (development, production)")
)

// buildVersion is set at compile time
var buildVersion = "dev"

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Fprintf(os.Stderr, "login-app version: %s\nGo version: %s (%s/%s)\n",
			buildVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	log.Printf("Starting login-app version %s; Go %s (%s/%s)",
		buildVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Load configuration
	cfg, err := config.Load(*flagEnv)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override port from command line if provided
	if *flagPort != "8080" {
		cfg.Server.Port = *flagPort
	}

	// Initialize storage (in-memory for this demo)
	userStore := storage.NewMemoryUserStore()

	// Create server
	srv, err := server.New(cfg, userStore)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
