// internal/services/UserService.go
package services

import (
	"fmt"

	"maxcool.com/weatherapp/internal/database"
	"maxcool.com/weatherapp/internal/models"
)

type IUserService interface {
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type UserService struct {
	DB database.IDB
}

// NewUserService creates a new UserService instance
func NewUserService(db database.IDB) *UserService {
	return &UserService{DB: db}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.DB.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.DB.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	id, err := s.DB.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.Id = id
	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *models.User) error {
	if err := s.DB.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(id int) error {
	if err := s.DB.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
