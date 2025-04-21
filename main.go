// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"maxcool.com/weatherapp/internal/config"
	"maxcool.com/weatherapp/internal/database"
	"maxcool.com/weatherapp/internal/handlers"
	"maxcool.com/weatherapp/internal/server"
	"maxcool.com/weatherapp/internal/services"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	database.EnsureDatabaseExists("weatherapp", cfg.PostgresConnectionString)
	log.Println("Correct database existance ensured.")

	// Migrate Database
	database.MigrateUpAll(cfg)

	// Initialize Database Connection
	db, err := database.NewDB(cfg.PostgresConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established.")

	// Create Handlers with Dependencies
	appHandler := handlers.NewHandler(services.NewUserService(db), services.NewSubscriptionService(db), cfg)

	// Setup Router and Server
	router := server.NewRouter(appHandler)
	srv := server.NewServer(":"+cfg.ServerPort, router)

	// Start Server
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful Shutdown
	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	// Register to receive SIGINT (Ctrl+C) and SIGTERM (termination signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	log.Println("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 seconds to wait for active requests
	defer cancel()                                                          // Release resources associated with this context

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
