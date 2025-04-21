// internal/services/SubscriptionService.go
package services

import (
	"fmt"

	"maxcool.com/weatherapp/internal/database"
	"maxcool.com/weatherapp/internal/models"
)

type ISubscriptionService interface {
	GetSubscriptionsByUserID(userID int) ([]*models.Subscription, error)
	CreateSubscription(subscription *models.Subscription) error
	UpdateSubscription(subscription *models.Subscription) error
	DeleteSubscription(id int) error
	GetSubscriptionByID(id int) (*models.Subscription, error)
}

type SubscriptionService struct {
	DB *database.DB
}

// NewSubscriptionService creates a new SubscriptionService instance
func NewSubscriptionService(db *database.DB) *SubscriptionService {
	return &SubscriptionService{DB: db}
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
	id, err := s.DB.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	subscription.Id = id
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
