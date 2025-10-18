package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"garage-barbershop/internal/services"
)

// UserHandler обрабатывает HTTP запросы для пользователей
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers обрабатывает GET /api/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")

	var users interface{}
	var err error

	if role != "" {
		users, err = h.userService.GetUsersByRole(role)
	} else {
		users, err = h.userService.GetAllUsers()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}

// GetUser обрабатывает GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL (упрощенная версия)
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser обрабатывает POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		TelegramID int64  `json:"telegram_id"`
		Username   string `json:"username"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Role       string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createdUser interface{}
	var err error

	switch user.Role {
	case "barber":
		createdUser, err = h.userService.RegisterBarber(
			user.TelegramID, user.Username, user.FirstName, user.LastName, user.Email,
		)
	case "client":
		createdUser, err = h.userService.RegisterClient(
			user.TelegramID, user.Username, user.FirstName, user.LastName, user.Email,
		)
	default:
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
