// internal/server/server.go
package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"maxcool.com/weatherapp/internal/handlers"
)

// NewRouter creates and configures the mux router
func NewRouter(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// // Define your routes and map them to handler methods
	// // POST endpoint to create a subscription (and user if not exists)
	// r.HandleFunc("/subscriptions", handler.HandlePostSubscription).Methods("POST")

	// // GET endpoint to get subscriptions for a specific user
	// // The {userID} is a path variable captured by mux
	// r.HandleFunc("/users/{userID}/subscriptions", handler.HandleGetUserSubscriptions).Methods("GET")

	r.HandleFunc("/weather", handler.GetWeatherHandler).Methods("GET")

	// Add a simple health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Optional: Add middleware here
	// r.Use(middleware.LoggingMiddleware)
	// r.Use(handlers.JSONContentTypeMiddleware)

	return r
}

// NewServer creates a custom http.Server
func NewServer(addr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv
}
