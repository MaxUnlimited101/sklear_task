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
	userService         services.IUserService
	subscriptionService services.ISubscriptionService
	config              *config.Config
}

// NewHandler Creates a new Handler instance
func NewHandler(userServicer services.IUserService, subscriptionService services.ISubscriptionService, config *config.Config) *Handler {
	return &Handler{userService: userServicer, subscriptionService: subscriptionService, config: config}
}

// --- Endpoints ---

const openweathermapAPIBaseUrl = "https://api.openweathermap.org/data/2.5/weather"

func (h *Handler) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	city := query.Get("city")

	// Make a request to the OpenWeatherMap API
	url := openweathermapAPIBaseUrl + "?q=" + city + "&appid=" + h.config.OpenWeatherMapAPIKey + "&units=metric"
	log.Print("GET ", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		log.Fatal("Error fetching weather data: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error fetching weather data", resp.StatusCode)
		log.Fatal("Error fetching weather data: ", err)
		return
	}

	// Read the response body
	var weatherResponse models.WeatherResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&weatherResponse)
	if err != nil {
		http.Error(w, "Error decoding weather data", http.StatusInternalServerError)
		log.Fatal("Error decoding weather data: ", err)
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
