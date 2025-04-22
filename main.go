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

	"github.com/go-co-op/gocron/v2"
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

	if err = database.EnsureDatabaseExists("weatherapp", cfg.PostgresConnectionString); err != nil {
		log.Fatalf("Failed to ensure database existence: %v", err)
	}
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
	appHandler := handlers.NewHandler(services.NewUserService(db), services.NewSubscriptionService(db, cfg), cfg)

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

	// Create a new scheduler instance
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	defer func() {
		_ = s.Shutdown()
	}()

	// Define the task function
	taskFunc := func() {
		log.Println("Running user notification task...")
		if err := appHandler.SubscriptionService.SendNotificationToUsers(); err != nil {
			log.Printf("Error sending notifications: %v", err)
		}
	}

	// Schedule the task
	_, _ = s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(12, 0, 0),
			),
		),
		gocron.NewTask(
			taskFunc,
		),
	)
	// debug task
	// _, _ = s.NewJob(
	// 	gocron.DurationJob(time.Second*20),
	// 	gocron.NewTask(
	// 		func() {
	// 			log.Print("Running scheduled task every 20 seconds...")
	// 			taskFunc()
	// 		},
	// 	),
	// )
	// s.Start()

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
