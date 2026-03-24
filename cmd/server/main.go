// Package main is the entrypoint for the calculator API server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calculator-api/internal/database"
	"calculator-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read configuration from environment variables with defaults.
	dbPath := os.Getenv("CALC_DB_PATH")
	if dbPath == "" {
		dbPath = "calculator.db"
	}

	listenAddr := os.Getenv("CALC_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:8000"
	}

	// Initialize the database.
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set up Gin router and register routes.
	router := gin.Default()
	handler.RegisterRoutes(router, db)

	// Set up HTTP server with graceful shutdown.
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	// Start server in a goroutine.
	go func() {
		log.Printf("Starting server on %s", listenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
