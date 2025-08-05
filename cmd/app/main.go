package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ministry-scheduler/internal/handler"
	"ministry-scheduler/internal/infra"
	"ministry-scheduler/internal/usecase"
)

const (
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 30 * time.Second
)

func main() {
	dbPath := getEnvOrDefault("DB_PATH", "users.db")
	port := getEnvOrDefault("PORT", "8080")

	log.Println("Starting Ministry Scheduler API...")
	log.Printf("Database path: %s", dbPath)
	log.Printf("Port: %s", port)

	db, err := infra.InitializeDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	userRepo := infra.NewSQLUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Ministry Scheduler API", "version": "1.0.0"}`)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "healthy"}`)
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           loggingMiddleware(mux),
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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
