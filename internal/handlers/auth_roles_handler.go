package handlers

import (
	"encoding/json"
	"net/http"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"
)

// AuthRolesHandler обрабатывает HTTP запросы с ролевой авторизацией
type AuthRolesHandler struct {
	authService services.AuthService
}

// NewAuthRolesHandler создает новый экземпляр AuthRolesHandler
func NewAuthRolesHandler(authService services.AuthService) *AuthRolesHandler {
	return &AuthRolesHandler{authService: authService}
}

// RegisterClient обрабатывает регистрацию клиента (публичный endpoint)
func (h *AuthRolesHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req models.ClientRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Регистрируем клиента
	user, err := h.authService.RegisterClient(req)
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

// RegisterBarber обрабатывает регистрацию барбера (только админ)
func (h *AuthRolesHandler) RegisterBarber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req models.BarberRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Регистрируем барбера
	user, err := h.authService.RegisterBarber(req)
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
