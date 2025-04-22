// internal/services/SubscriptionService.go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"
	"maxcool.com/weatherapp/internal/config"
	"maxcool.com/weatherapp/internal/database"
	"maxcool.com/weatherapp/internal/models"
)

type ISubscriptionService interface {
	GetSubscriptionsByUserID(userID int) ([]*models.Subscription, error)
	CreateSubscription(subscription *models.Subscription) error
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id int) error
	GetSubscriptionByID(id int) (*models.Subscription, error)
	GetWeather(city string) (models.WeatherResponse, error)
	CheckCondition(condition, city string) (bool, error)
	SendNotificationToUsers() error
	CheckWhetherCityExists(city string) (bool, error)
}

type SubscriptionService struct {
	DB     database.IDB
	Config *config.Config
}

// NewSubscriptionService creates a new SubscriptionService instance
func NewSubscriptionService(db database.IDB, cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{DB: db, Config: cfg}
}

// GetSubscriptionsByUserID retrieves all subscriptions for a given user ID
func (s *SubscriptionService) GetSubscriptionsByUserID(userID int) ([]*models.Subscription, error) {
	subscriptions, err := s.DB.GetSubscriptionsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by user ID: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*models.Subscription, len(subscriptions))
	for i := range subscriptions {
		result[i] = &subscriptions[i]
	}
	return result, nil
}

// CreateSubscription creates a new subscription
func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) error {
	user, err := s.DB.GetUserByEmail(subscription.UserEmail)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found: %s", subscription.UserEmail)
	}
	subscription.UserId = user.Id
	subscription.UserEmail = user.Email

	_, err = s.DB.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// UpdateSubscription updates an existing subscription
func (s *SubscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	if err := s.DB.UpdateSubscription(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

// DeleteSubscription deletes a subscription by its ID
func (s *SubscriptionService) DeleteSubscription(id int) error {
	if err := s.DB.DeleteSubscription(id); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

// GetSubscriptionByID retrieves a subscription by its ID
func (s *SubscriptionService) GetSubscriptionByID(id int) (*models.Subscription, error) {
	subscription, err := s.DB.GetSubscriptionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription by ID: %w", err)
	}
	return subscription, nil
}

// SendEmail sends an email using the provided mailer
// It takes the recipient's email address, subject, and body as parameters
// It returns an error if sending the email fails
func SendEmail(to string, subject string, body string, config *config.Config) error {
	client := resend.NewClient(config.ResendApiKey)

	log.Printf("Sending email to %s.", to)

	params := &resend.SendEmailRequest{
		From:    "weatherapp@resend.dev",
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Printf("Email sent successfully: %s", sent.Id)
	return nil
}

const openweathermapAPIBaseUrl = "https://api.openweathermap.org/data/2.5/weather"

func (s *SubscriptionService) GetWeather(city string) (models.WeatherResponse, error) {
	// Make a request to the OpenWeatherMap API
	url := openweathermapAPIBaseUrl + "?q=" + city + "&appid=" + s.Config.OpenWeatherMapAPIKey + "&units=metric"
	log.Print("GET ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching weather data: ", err)
		return models.WeatherResponse{}, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching weather data: %v", err)
		return models.WeatherResponse{}, fmt.Errorf("error fetching weather data: %s", resp.Status)
	}

	// Read the response body
	var weatherResponse models.WeatherResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&weatherResponse)
	if err != nil {
		log.Print("Error decoding weather data: ", err)
		return models.WeatherResponse{}, err
	}
	return weatherResponse, nil
}

// GetWeather retrieves the weather data for a given city
// It takes the city name as a parameter
// It returns the weather data and an error if the request fails
// It returns an error if the weather data cannot be fetched
func (s *SubscriptionService) CheckWhetherCityExists(city string) (bool, error) {
	// Make a request to the OpenWeatherMap API
	log.Print("checking whether city exists")
	url := openweathermapAPIBaseUrl + "?q=" + city + "&appid=" + s.Config.OpenWeatherMapAPIKey + "&units=metric"
	log.Print("GET ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching weather data: ", err)
		return false, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("City not found: %s", city)
			return false, nil
		}
		log.Printf("Error fetching weather data: %v", err)
		return false, fmt.Errorf("error fetching weather data: %s", resp.Status)
	}

	return true, nil
}

// CheckCondition checks if the weather condition is met for a given city
// It takes the condition string and city name as parameters
// It returns true if the condition is met, false otherwise
// It returns an error if the weather data cannot be fetched
func (s *SubscriptionService) CheckCondition(condition, city string) (bool, error) {
	weatherResponse, err := s.GetWeather(city)
	if err != nil {
		return false, fmt.Errorf("failed to get weather data: %w", err)
	}

	splitCondition := strings.Split(condition, ":")

	if len(splitCondition) < 2 {
		return false, fmt.Errorf("invalid condition format: %s", condition)
	}
	if splitCondition[0] != "main" && len(splitCondition) != 3 {
		return false, fmt.Errorf("invalid condition format: %s", condition)
	}

	var prop float32
	var operator func(a, b float32) bool
	var condVal float32

	switch splitCondition[0] {
	case "feels_like":
		prop = float32(weatherResponse.Main.Feels_like)
	case "humidity":
		prop = float32(weatherResponse.Main.Humidity)
	case "main": // special case
		return strings.EqualFold(weatherResponse.Weather[0].Main, splitCondition[1]), nil
	case "temperature":
		prop = float32(weatherResponse.Main.Temp)
	default:
		return false, fmt.Errorf("invalid property: %s", splitCondition[0])
	}

	switch splitCondition[1] {
	case "==":
		operator = func(a, b float32) bool { return a == b }
	case "!=":
		operator = func(a, b float32) bool { return a != b }
	case "<":
		operator = func(a, b float32) bool { return a < b }
	case "<=":
		operator = func(a, b float32) bool { return a <= b }
	case ">":
		operator = func(a, b float32) bool { return a > b }
	case ">=":
		operator = func(a, b float32) bool { return a >= b }
	default:
		return false, fmt.Errorf("invalid operator: %s", splitCondition[1])
	}

	condVal64, err := strconv.ParseFloat(splitCondition[2], 32)
	condVal = float32(condVal64)
	if err != nil {
		return false, fmt.Errorf("invalid condition value: %s", splitCondition[2])
	}

	if operator(prop, condVal) {
		log.Printf("Condition met: %s %s %f for city %s", splitCondition[0], splitCondition[1], condVal, city)
		return true, nil
	}

	return false, nil
}

// SendNotificationToUsers sends notifications to users based on their subscriptions
// It retrieves all subscriptions from the database
// It checks if the weather condition is met for each subscription
// If the condition is met, it sends an email notification to the user
func (s *SubscriptionService) SendNotificationToUsers() error {
	subscriptions, err := s.DB.GetSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	for _, subscription := range subscriptions {
		// Check the weather condition
		met, err := s.CheckCondition(subscription.Condition, subscription.City)
		if err != nil {
			log.Printf("Error for subscription %d: %v", subscription.Id, err)
			continue
		}

		if met {
			// Send email notification
			subject := "Weather Update"
			body := fmt.Sprintf("The weather condition `%s` is met for city `%s`.", subscription.Condition, subscription.City)
			err = SendEmail(subscription.UserEmail, subject, body, s.Config)
			if err != nil {
				log.Printf("Failed to send email to user %s: %v", subscription.UserEmail, err)
				continue
			}
			log.Printf("Notification sent to user %s for subscription %d", subscription.UserEmail, subscription.Id)

			// Add new notification to the database
			notif := models.Notification{
				UserId:         subscription.UserId,
				SubscriptionId: subscription.Id,
				SentAt:         time.Now(),
			}

			log.Printf("Creating notification in DB for subscription with userId %d and city %s",
				subscription.UserId, subscription.City)
			_, err = s.DB.CreateNotification(&notif)
			if err != nil {
				log.Printf("Failed to create notification in DB for subscription with userId %d and city %s: %v",
					subscription.UserId, subscription.City, err)
				continue
			}
		}
	}

	return nil
}
