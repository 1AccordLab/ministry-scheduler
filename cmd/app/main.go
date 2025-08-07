package main

import (
	"context"
	"fmt"
	"log"
	"ministry-scheduler/internal/features/users"
	"ministry-scheduler/internal/shared/database"
	"ministry-scheduler/internal/shared/env"
	"ministry-scheduler/internal/shared/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 30 * time.Second
)

func main() {
	dbPath := env.GetEnvOrDefault("DB_PATH", "users.db")
	port := env.GetEnvOrDefault("PORT", "8080")

	log.Println("Starting Ministry Scheduler API...")
	log.Printf("Database path: %s", dbPath)
	log.Printf("Port: %s", port)

	// Initialize database
	db, err := database.InitializeDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	// Initialize user feature
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	// Setup HTTP routes
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	// Health and info endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Ministry Scheduler API", "version": "1.0.0"}`)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "healthy"}`)
	})

	// Setup server with middleware
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           middleware.LoggingMiddleware(mux),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if serverErr := server.ListenAndServe(); serverErr != nil && serverErr != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", serverErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
		log.Printf("Server forced to shutdown: %v", shutdownErr)
		return
	}

	log.Println("Server exited")
}
