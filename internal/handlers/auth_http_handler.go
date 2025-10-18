package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"
)

// AuthHTTPHandler HTTP обработчик для аутентификации (без Gin)
type AuthHTTPHandler struct {
	authService services.AuthService
}

// NewAuthHTTPHandler создает новый HTTP обработчик аутентификации
func NewAuthHTTPHandler(authService services.AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

// TelegramAuth обрабатывает аутентификацию через Telegram
func (h *AuthHTTPHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var authData models.TelegramAuthData
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидируем Telegram данные
	botToken := "your_bot_token_here" // В реальном приложении получать из конфигурации
	if !h.authService.ValidateTelegramAuth(authData, botToken) {
		http.Error(w, "Invalid Telegram authentication", http.StatusUnauthorized)
		return
	}

	// Находим или создаем пользователя
	user, err := h.authService.AuthenticateUser(authData)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Генерируем access token
	accessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Генерируем refresh token
	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "Refresh token generation failed", http.StatusInternalServerError)
		return
	}

	// Сохраняем refresh token в Redis
	if err := h.authService.StoreRefreshToken(user.ID, refreshToken); err != nil {
		http.Error(w, "Token storage failed", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		User:         *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken обновляет токены
func (h *AuthHTTPHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Парсим refresh token
	claims, err := h.authService.ParseJWT(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Проверяем тип токена
	if !claims.IsRefreshToken() {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	// Проверяем срок действия
	if claims.IsExpired() {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	// Проверяем, что токен существует в Redis
	if !h.authService.IsRefreshTokenValid(claims.UserID, req.RefreshToken) {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	// Получаем пользователя из БД
	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Генерируем новую пару токенов
	newAccessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "Refresh token generation failed", http.StatusInternalServerError)
		return
	}

	// Обновляем refresh token в Redis
	if err := h.authService.UpdateRefreshToken(claims.UserID, req.RefreshToken, newRefreshToken); err != nil {
		http.Error(w, "Token update failed", http.StatusInternalServerError)
		return
	}

	// Возвращаем новые токены
	response := models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		User:         *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout выходит из системы
func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// В реальном приложении нужно извлечь user_id из JWT токена
	// Для упрощения возвращаем успех
	response := map[string]string{
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile возвращает профиль текущего пользователя
func (h *AuthHTTPHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// В реальном приложении нужно извлечь данные из JWT токена
	// Для упрощения возвращаем базовую информацию
	profile := map[string]interface{}{
		"user_id":          1,
		"telegram_id":      12345,
		"role":             "client",
		"authenticated_at": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// RegisterDirect обрабатывает прямую регистрацию пользователя
func (h *AuthHTTPHandler) RegisterDirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req models.DirectRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Регистрируем пользователя
	user, err := h.authService.RegisterUserDirect(req)
	if err != nil {
		http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Генерируем токены
	accessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "Ошибка генерации access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "Ошибка генерации refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем refresh token
	if err := h.authService.StoreRefreshToken(user.ID, refreshToken); err != nil {
		http.Error(w, "Ошибка сохранения refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 минут
		User:         *user,
	})
}

// LoginDirect обрабатывает прямую авторизацию пользователя
func (h *AuthHTTPHandler) LoginDirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req models.DirectLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Авторизуем пользователя
	user, err := h.authService.LoginDirect(req)
	if err != nil {
		http.Error(w, "Ошибка авторизации: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Генерируем токены
	accessToken, err := h.authService.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "Ошибка генерации access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "Ошибка генерации refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем refresh token
	if err := h.authService.StoreRefreshToken(user.ID, refreshToken); err != nil {
		http.Error(w, "Ошибка сохранения refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 минут
		User:         *user,
	})
}
