// internal/tests/MockDB.go
package tests

import (
	"github.com/stretchr/testify/mock"
	"maxcool.com/weatherapp/internal/database"
	"maxcool.com/weatherapp/internal/models"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreateUser(user *models.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) CreateSubscription(sub *models.Subscription) (int, error) {
	args := m.Called(sub)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetSubscriptionByID(subID int) (*models.Subscription, error) {
	args := m.Called(subID)
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockDB) UpdateSubscription(sub *models.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func (m *MockDB) DeleteSubscription(subID int) error {
	args := m.Called(subID)
	return args.Error(0)
}

func (m *MockDB) GetSubscriptions() ([]models.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (m *MockDB) Close() {
}

func (m *MockDB) CreateNotification(notification *models.Notification) (int, error) {
	args := m.Called(notification)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetSubscriptionsByUserID(userID int) ([]models.Subscription, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Subscription), args.Error(1)
}

var _ database.IDB = &MockDB{}
