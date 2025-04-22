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

	r.HandleFunc("/subscribe", handler.PostSubscriptionHandler).Methods("POST")

	r.HandleFunc("/weather", handler.GetWeatherHandler).Methods("GET")

	r.HandleFunc("/user", handler.PostUserHandler).Methods("POST")

	// Add a simple health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

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
