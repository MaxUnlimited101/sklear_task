package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"maxcool.com/weatherapp/internal/config"
	"maxcool.com/weatherapp/internal/models"
	"maxcool.com/weatherapp/internal/services"
)

var Validator *validator.Validate = validator.New()

// Validate validates the given struct using the validator
// It returns an error if validation fails
func Validate(i any) error {
	err := Validator.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Validation error: %s", err)
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// SendJsonResponse sends a JSON response with the given status code and data
// It sets the Content-Type header to application/json and encodes the data into JSON format
// If encoding fails, it logs the error and sends an internal server error response
func SendJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handler struct holds dependencies for handlers
type Handler struct {
	UserService         services.IUserService
	SubscriptionService services.ISubscriptionService
	Config              *config.Config
}

// NewHandler Creates a new Handler instance
func NewHandler(userServicer services.IUserService, subscriptionService services.ISubscriptionService, config *config.Config) *Handler {
	return &Handler{UserService: userServicer, SubscriptionService: subscriptionService, Config: config}
}

// --- Endpoints ---

func (h *Handler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	city := query.Get("city")

	// Validate the city parameter
	if city == "" || len(city) == 0 {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		log.Println("City parameter is required")
		return
	}

	if f, err := h.SubscriptionService.CheckWhetherCityExists(city); err == nil && !f {
		http.Error(w, "City not found", http.StatusNotFound)
		log.Println("City not found: ", city)
		return
	} else {
		if err != nil {
			log.Println("Failed to check city existence: ", err)
		}
	}

	weatherResponse, err := h.SubscriptionService.GetWeather(city)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		log.Println("Failed to fetch weather data: ", err)
		return
	}

	// Create a response object
	response := map[string]any{
		"city":        city,
		"temperature": weatherResponse.Main.Temp,
		"humidity":    weatherResponse.Main.Humidity,
		"feels_like":  weatherResponse.Main.Feels_like,
		"main":        weatherResponse.Weather[0].Main,
	}

	SendJsonResponse(w, http.StatusOK, response)
}

func (h *Handler) PostSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Invalid request payload: ", err)
		return
	}

	if err := Validate(subscription); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		log.Println("Validation failed: ", err)
		return
	}

	if f, err := h.SubscriptionService.CheckWhetherCityExists(subscription.City); err == nil && !f {
		http.Error(w, "City not found", http.StatusNotFound)
		log.Println("City not found: ", subscription.City)
		return
	} else {
		if err != nil {
			log.Println("Failed to check city existence: ", err)
		}
	}

	if err := h.SubscriptionService.CreateSubscription(&subscription); err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		log.Println("Failed to create subscription: ", err)
		return
	}

	SendJsonResponse(w, http.StatusCreated, "Subscription created successfully")
}

func (h *Handler) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Invalid request payload: ", err)
		return
	}

	if err := Validate(user); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		log.Println("Validation failed: ", err)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("Failed to create user: ", err)
		return
	}

	SendJsonResponse(w, http.StatusCreated, "User created successfully")
}
