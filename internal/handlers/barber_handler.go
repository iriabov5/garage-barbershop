package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"
)

// BarberHandler обрабатывает HTTP запросы, связанные с барберами
type BarberHandler struct {
	barberService services.BarberService
}

// NewBarberHandler создает новый экземпляр BarberHandler
func NewBarberHandler(barberService services.BarberService) *BarberHandler {
	return &BarberHandler{barberService: barberService}
}

// AdminGetAllBarbers получает всех барберов (только админ)
func (h *BarberHandler) AdminGetAllBarbers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	barbers, err := h.barberService.GetAllBarbers()
	if err != nil {
		http.Error(w, "Ошибка получения барберов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"barbers": barbers,
		"count":   len(barbers),
	})
}

// AdminGetBarber получает барбера по ID (только админ)
func (h *BarberHandler) AdminGetBarber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	barberID, err := h.extractIDFromURL(r.URL.Path, "/api/admin/barbers/")
	if err != nil {
		http.Error(w, "Неверный ID барбера: "+err.Error(), http.StatusBadRequest)
		return
	}

	barber, err := h.barberService.GetBarberByID(barberID)
	if err != nil {
		http.Error(w, "Ошибка получения барбера: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barber)
}

// AdminUpdateBarber обновляет барбера (только админ)
func (h *BarberHandler) AdminUpdateBarber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	barberID, err := h.extractIDFromURL(r.URL.Path, "/api/admin/barbers/")
	if err != nil {
		http.Error(w, "Неверный ID барбера: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req models.BarberUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	barber, err := h.barberService.UpdateBarber(barberID, req)
	if err != nil {
		http.Error(w, "Ошибка обновления барбера: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barber)
}

// AdminDeleteBarber удаляет барбера (только админ)
func (h *BarberHandler) AdminDeleteBarber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL
	barberID, err := h.extractIDFromURL(r.URL.Path, "/api/admin/barbers/")
	if err != nil {
		http.Error(w, "Неверный ID барбера: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.barberService.DeleteBarber(barberID); err != nil {
		http.Error(w, "Ошибка удаления барбера: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Барбер успешно удален",
	})
}

// BarberGetSelf получает собственный профиль барбера
func (h *BarberHandler) BarberGetSelf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID барбера из контекста
	barberID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	barber, err := h.barberService.GetBarberSelf(barberID)
	if err != nil {
		http.Error(w, "Ошибка получения профиля: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barber)
}

// BarberUpdateSelf обновляет собственный профиль барбера
func (h *BarberHandler) BarberUpdateSelf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID барбера из контекста
	barberID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	var req models.BarberSelfUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные: "+err.Error(), http.StatusBadRequest)
		return
	}

	barber, err := h.barberService.UpdateBarberSelf(barberID, req)
	if err != nil {
		http.Error(w, "Ошибка обновления профиля: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barber)
}

// extractIDFromURL извлекает ID из URL
func (h *BarberHandler) extractIDFromURL(path, prefix string) (uint, error) {
	idStr := strings.TrimPrefix(path, prefix)
	if idStr == "" {
		return 0, fmt.Errorf("ID не указан")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("неверный формат ID")
	}

	return uint(id), nil
}
