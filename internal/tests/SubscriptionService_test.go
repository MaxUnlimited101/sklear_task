// internal/tests/SubscriptionService_test.go
package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"maxcool.com/weatherapp/internal/models"
	"maxcool.com/weatherapp/internal/services"
)

func TestGetSubscriptionsByUserID(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	expectedSubscriptions := []models.Subscription{
		{Id: 1, UserId: 1, City: "New York", Condition: "temperature:>:30"},
		{Id: 2, UserId: 1, City: "Los Angeles", Condition: "humidity:<:50"},
	}
	mockDB.On("GetSubscriptionsByUserID", 1).Return(expectedSubscriptions, nil)

	subscriptions, err := subscriptionService.GetSubscriptionsByUserID(1)

	assert.NoError(t, err)
	assert.Len(t, subscriptions, 2)
	assert.Equal(t, expectedSubscriptions[0].City, subscriptions[0].City)
	mockDB.AssertExpectations(t)
}

func TestCreateSubscription(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	user := &models.User{Id: 1, Email: "test@example.com"}
	subscription := &models.Subscription{City: "New York", Condition: "temperature:>:30", UserEmail: "test@example.com"}

	mockDB.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockDB.On("CreateSubscription", subscription).Return(1, nil)

	err := subscriptionService.CreateSubscription(subscription)

	assert.NoError(t, err)
	assert.Equal(t, 1, subscription.UserId)
	mockDB.AssertExpectations(t)
}

func TestCreateSubscription_UserNotFound(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	subscription := &models.Subscription{City: "New York", Condition: "temperature:>:30", UserEmail: "test@example.com"}

	mockDB.On("GetUserByEmail", "test@example.com").Return((*models.User)(nil), nil)

	err := subscriptionService.CreateSubscription(subscription)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockDB.AssertExpectations(t)
}

func TestUpdateSubscription(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	subscription := &models.Subscription{Id: 1, City: "New York", Condition: "temperature:>:30"}

	mockDB.On("UpdateSubscription", subscription).Return(nil)

	err := subscriptionService.UpdateSubscription(subscription)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestDeleteSubscription(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	mockDB.On("DeleteSubscription", 1).Return(nil)

	err := subscriptionService.DeleteSubscription(1)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetSubscriptionByID(t *testing.T) {
	mockDB := new(MockDB)
	subscriptionService := services.NewSubscriptionService(mockDB, nil)

	expectedSubscription := &models.Subscription{Id: 1, City: "New York", Condition: "temperature:>:30"}
	mockDB.On("GetSubscriptionByID", 1).Return(expectedSubscription, nil)

	subscription, err := subscriptionService.GetSubscriptionByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubscription, subscription)
	mockDB.AssertExpectations(t)
}
