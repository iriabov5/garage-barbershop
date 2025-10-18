package unit

import (
	"errors"
	"testing"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserService_CreateUser - тест создания пользователя
func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	user := &models.User{
		TelegramID: 12345,
		Username:   "testuser",
		FirstName:  "John",
		LastName:   "Doe",
	}

	// Настраиваем мок
	mockRepo.On("Create", user).Return(nil)

	// Act
	err := userService.CreateUser(user)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_CreateUser_Error - тест ошибки при создании
func TestUserService_CreateUser_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	user := &models.User{
		TelegramID: 12345,
		Username:   "testuser",
		FirstName:  "John",
		LastName:   "Doe",
	}

	// Настраиваем мок для возврата ошибки
	mockRepo.On("Create", user).Return(errors.New("database error"))

	// Act
	err := userService.CreateUser(user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

// TestUserService_RegisterBarber - тест регистрации барбера
func TestUserService_RegisterBarber(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	telegramID := int64(12345)
	username := "barber_user"
	firstName := "Ivan"
	lastName := "Barber"
	email := "barber@example.com"

	// Настраиваем моки
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Моки для ролей
	barberRole := &models.Role{ID: 1, Name: "barber"}
	mockRoleRepo.On("GetRoleByName", "barber").Return(barberRole, nil)
	mockRoleRepo.On("AssignRoleToUser", mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	// Act
	barber, err := userService.RegisterBarber(telegramID, username, firstName, lastName, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, barber)
	assert.Equal(t, telegramID, barber.TelegramID)
	assert.Equal(t, username, barber.Username)
	assert.Equal(t, firstName, barber.FirstName)
	assert.Equal(t, lastName, barber.LastName)
	// Роли теперь проверяются через RoleService
	assert.True(t, barber.IsActive)
	assert.Equal(t, 5.0, barber.Rating)
	mockRepo.AssertExpectations(t)
}

// TestUserService_RegisterClient - тест регистрации клиента
func TestUserService_RegisterClient(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	telegramID := int64(67890)
	username := "client_user"
	firstName := "Jane"
	lastName := "Client"
	email := "client@example.com"

	// Настраиваем моки
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Моки для ролей
	clientRole := &models.Role{ID: 2, Name: "client"}
	mockRoleRepo.On("GetRoleByName", "client").Return(clientRole, nil)
	mockRoleRepo.On("AssignRoleToUser", mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)

	// Act
	client, err := userService.RegisterClient(telegramID, username, firstName, lastName, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, telegramID, client.TelegramID)
	assert.Equal(t, username, client.Username)
	assert.Equal(t, firstName, client.FirstName)
	assert.Equal(t, lastName, client.LastName)
	// Роли теперь проверяются через RoleService
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByID - тест получения пользователя по ID
func TestUserService_GetUserByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	userID := uint(1)
	expectedUser := &models.User{
		ID:         userID,
		TelegramID: 12345,
		Username:   "testuser",
		FirstName:  "John",
		LastName:   "Doe",
	}

	// Настраиваем мок
	mockRepo.On("GetByID", userID).Return(expectedUser, nil)

	// Act
	user, err := userService.GetUserByID(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.TelegramID, user.TelegramID)
	assert.Equal(t, expectedUser.Username, user.Username)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetUserByID_NotFound - тест пользователя не найден
func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockRepo, mockRoleRepo)

	userID := uint(999)

	// Настраиваем мок для возврата ошибки
	mockRepo.On("GetByID", userID).Return((*models.User)(nil), errors.New("user not found"))

	// Act
	user, err := userService.GetUserByID(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
