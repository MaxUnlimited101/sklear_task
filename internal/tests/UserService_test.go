// internal/tests/UserService_test.go
package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"maxcool.com/weatherapp/internal/models"
	"maxcool.com/weatherapp/internal/services"
)

func TestGetUserByID(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	expectedUser := &models.User{Id: 1, Name: "John Doe", Email: "john.doe@example.com"}
	mockDB.On("GetUserByID", 1).Return(expectedUser, nil)

	user, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockDB.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	mockDB.On("GetUserByID", 1).Return((*models.User)(nil), nil)

	user, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Nil(t, user)
	mockDB.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	newUser := &models.User{Name: "Jane Doe", Email: "jane.doe@example.com"}
	mockDB.On("CreateUser", newUser).Return(1, nil)

	err := userService.CreateUser(newUser)

	assert.NoError(t, err)
	assert.Equal(t, 1, newUser.Id)
	mockDB.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	newUser := &models.User{Name: "Jane Doe", Email: "jane.doe@example.com"}
	mockDB.On("CreateUser", newUser).Return(0, errors.New("database error"))

	err := userService.CreateUser(newUser)

	assert.Error(t, err)
	assert.Equal(t, 0, newUser.Id)
	mockDB.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	updatedUser := &models.User{Id: 1, Name: "John Smith", Email: "john.smith@example.com"}
	mockDB.On("UpdateUser", updatedUser).Return(nil)

	err := userService.UpdateUser(updatedUser)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateUser_Error(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	updatedUser := &models.User{Id: 1, Name: "John Smith", Email: "john.smith@example.com"}
	mockDB.On("UpdateUser", updatedUser).Return(errors.New("database error"))

	err := userService.UpdateUser(updatedUser)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	mockDB.On("DeleteUser", 1).Return(nil)

	err := userService.DeleteUser(1)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestDeleteUser_Error(t *testing.T) {
	mockDB := new(MockDB)
	userService := services.NewUserService(mockDB)

	mockDB.On("DeleteUser", 1).Return(errors.New("database error"))

	err := userService.DeleteUser(1)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}
