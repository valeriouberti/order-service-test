// cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/valeriouberti/order-service-test/internal/api"
	"github.com/valeriouberti/order-service-test/internal/api/handlers"
	"github.com/valeriouberti/order-service-test/internal/config"
	"github.com/valeriouberti/order-service-test/internal/repository"
	"github.com/valeriouberti/order-service-test/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to the database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to the database successfully")

	// Initialize repositories
	productRepo := repository.NewProductRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// Initialize services
	orderService := services.NewOrderService(orderRepo, productRepo)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize router
	router := api.NewRouter(orderHandler)

	// Configure HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine so it doesn't block graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
