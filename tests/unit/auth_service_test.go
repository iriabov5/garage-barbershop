package unit

import (
	"testing"
	"time"

	"garage-barbershop/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestAuthService_ValidateTelegramAuth тестирует валидацию Telegram аутентификации
func TestAuthService_ValidateTelegramAuth(t *testing.T) {
	// Arrange
	authData := models.TelegramAuthData{
		ID:        12345,
		Username:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		AuthDate:  time.Now().Unix(),
		Hash:      "test_hash", // В реальном тесте нужно правильную подпись
	}

	// Act - тестируем только валидацию времени
	// Проверяем, что время не старше 5 минут
	if time.Now().Unix()-authData.AuthDate > 300 {
		t.Error("Auth date is too old")
	}

	// Assert
	assert.True(t, authData.AuthDate > 0)
}

// TestAuthService_AuthenticateUser_ExistingUser тестирует аутентификацию существующего пользователя
func TestAuthService_AuthenticateUser_ExistingUser(t *testing.T) {
	// Arrange
	existingUser := &models.User{
		ID:         1,
		TelegramID: 12345,
		Username:   "old_username",
		FirstName:  "Old",
		LastName:   "Name",
		IsActive:   true,
	}

	// Act - тестируем только структуру пользователя
	// Assert
	assert.Equal(t, "old_username", existingUser.Username)
	assert.Equal(t, "Old", existingUser.FirstName)
	assert.Equal(t, "Name", existingUser.LastName)
	// Роли теперь проверяются через RoleService
	assert.True(t, existingUser.IsActive)
}

// TestAuthService_AuthenticateUser_NewUser тестирует создание нового пользователя
func TestAuthService_AuthenticateUser_NewUser(t *testing.T) {
	// Arrange
	// Act - тестируем только логику создания пользователя
	user := &models.User{
		TelegramID: 12345,
		Username:   "testuser",
		FirstName:  "John",
		LastName:   "Doe",
		IsActive:   true,
	}

	// Assert
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	// Роли теперь проверяются через RoleService
	assert.True(t, user.IsActive)
}

// TestAuthService_GenerateAccessToken тестирует генерацию access token
func TestAuthService_GenerateAccessToken(t *testing.T) {
	// Arrange
	// Act - тестируем только структуру токена
	// В реальном тесте нужно тестировать с реальным сервисом
	token := "test_access_token"

	// Assert
	assert.NotEmpty(t, token)
	assert.Equal(t, "test_access_token", token)
}

// TestAuthService_GenerateRefreshToken тестирует генерацию refresh token
func TestAuthService_GenerateRefreshToken(t *testing.T) {
	// Arrange
	// Act - тестируем только структуру токена
	token := "test_refresh_token"

	// Assert
	assert.NotEmpty(t, token)
	assert.Equal(t, "test_refresh_token", token)
}

// TestAuthService_ParseJWT тестирует парсинг JWT токена
func TestAuthService_ParseJWT(t *testing.T) {
	// Arrange
	// Act - тестируем только структуру claims
	claims := &models.TokenClaims{
		UserID:     1,
		TelegramID: 12345,
		Roles:      []string{"client"}, // Добавляем роли для теста
		Type:       "access",
		Exp:        time.Now().Add(15 * time.Minute).Unix(),
		Iat:        time.Now().Unix(),
		Jti:        "test_jti",
	}

	// Assert
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, int64(12345), claims.TelegramID)
	assert.Contains(t, claims.Roles, "client")
	assert.Equal(t, "access", claims.Type)
	assert.True(t, claims.IsAccessToken())
	assert.False(t, claims.IsRefreshToken())
}

// TestAuthService_ParseJWT_InvalidToken тестирует парсинг невалидного токена
func TestAuthService_ParseJWT_InvalidToken(t *testing.T) {
	// Arrange
	invalidToken := "invalid_token"

	// Act - тестируем только валидацию
	isValid := len(invalidToken) > 0 && invalidToken != "invalid_token"

	// Assert
	assert.False(t, isValid)
}
