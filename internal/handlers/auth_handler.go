package handlers

import (
	"net/http"
	"time"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler обработчик для аутентификации
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// TelegramAuth обрабатывает аутентификацию через Telegram
func (h *AuthHandler) TelegramAuth(c *gin.Context) {
	var authData models.TelegramAuthData
	if err := c.ShouldBindJSON(&authData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Валидируем Telegram данные
	// Получаем bot token из конфигурации (упрощенная версия)
	botToken := "your_bot_token_here" // В реальном приложении получать из конфигурации
	if !h.authService.ValidateTelegramAuth(authData, botToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидная аутентификация Telegram"})
		return
	}

	// Находим или создаем пользователя
	user, err := h.authService.AuthenticateUser(authData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка аутентификации пользователя"})
		return
	}

	// Генерируем access token
	accessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Генерируем refresh token
	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации refresh токена"})
		return
	}

	// Сохраняем refresh token в Redis
	if err := h.authService.StoreRefreshToken(user.ID, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения токена"})
		return
	}

	// Возвращаем ответ
	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		User:         *user,
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken обновляет токены
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token обязателен"})
		return
	}

	// Парсим refresh token
	claims, err := h.authService.ParseJWT(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный refresh token"})
		return
	}

	// Проверяем тип токена
	if !claims.IsRefreshToken() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный тип токена"})
		return
	}

	// Проверяем срок действия
	if claims.IsExpired() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token истек"})
		return
	}

	// Проверяем, что токен существует в Redis
	if !h.authService.IsRefreshTokenValid(claims.UserID, req.RefreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token не найден"})
		return
	}

	// Получаем пользователя (в реальном приложении нужно получить из БД)
	// Для упрощения создаем временного пользователя
	user := &models.User{
		ID:         claims.UserID,
		TelegramID: claims.TelegramID,
		Role:       claims.Role,
	}

	// Генерируем новую пару токенов
	newAccessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	newRefreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации refresh токена"})
		return
	}

	// Обновляем refresh token в Redis
	if err := h.authService.UpdateRefreshToken(claims.UserID, req.RefreshToken, newRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления токена"})
		return
	}

	// Возвращаем новые токены
	response := models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		User:         *user,
	}

	c.JSON(http.StatusOK, response)
}

// Logout выходит из системы
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	// Отзываем refresh token
	if err := h.authService.RevokeRefreshToken(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка выхода из системы"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешный выход из системы"})
}

// GetProfile возвращает профиль текущего пользователя
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}

	// В реальном приложении нужно получить пользователя из БД
	// Для упрощения возвращаем базовую информацию
	telegramID, _ := c.Get("telegram_id")
	userRole, _ := c.Get("user_role")

	profile := gin.H{
		"user_id":          userID,
		"telegram_id":      telegramID,
		"role":             userRole,
		"authenticated_at": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, profile)
}
