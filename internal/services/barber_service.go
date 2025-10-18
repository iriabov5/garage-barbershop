package services

import (
	"fmt"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"
)

// BarberService интерфейс для управления барберами
type BarberService interface {
	// Управление барберами (только админ)
	UpdateBarber(barberID uint, req models.BarberUpdateRequest) (*models.User, error)
	DeleteBarber(barberID uint) error
	GetBarberByID(barberID uint) (*models.User, error)
	GetAllBarbers() ([]models.User, error)

	// Управление собственным профилем барбера
	UpdateBarberSelf(barberID uint, req models.BarberSelfUpdateRequest) (*models.User, error)
	GetBarberSelf(barberID uint) (*models.User, error)
}

// barberService реализация BarberService
type barberService struct {
	userRepo  repositories.UserRepository
	roleRepo  repositories.RoleRepository
}

// NewBarberService создает новый экземпляр BarberService
func NewBarberService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) BarberService {
	return &barberService{userRepo: userRepo, roleRepo: roleRepo}
}

// UpdateBarber обновляет барбера (только админ)
func (s *barberService) UpdateBarber(barberID uint, req models.BarberUpdateRequest) (*models.User, error) {
	// Получаем барбера
	barber, err := s.userRepo.GetByID(barberID)
	if err != nil {
		return nil, fmt.Errorf("барбер не найден: %v", err)
	}

	// Проверяем, что это барбер
	if !s.roleRepo.HasUserRole(barberID, "barber") {
		return nil, fmt.Errorf("пользователь не является барбером")
	}

	// Обновляем поля, если они переданы
	if req.Email != "" {
		// Проверяем, что email не занят другим пользователем
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != barberID {
			return nil, fmt.Errorf("email уже занят")
		}
		barber.Email = req.Email
	}

	if req.FirstName != "" {
		barber.FirstName = req.FirstName
	}

	if req.LastName != "" {
		barber.LastName = req.LastName
	}

	if req.Specialties != "" {
		barber.Specialties = req.Specialties
	}

	if req.Experience > 0 {
		barber.Experience = req.Experience
	}

	// Обновляем поля, которые могут быть nil (используем указатели)
	if req.IsActive != nil {
		barber.IsActive = *req.IsActive
	}

	if req.Rating != nil {
		barber.Rating = *req.Rating
	}

	// Сохраняем изменения
	if err := s.userRepo.Update(barber); err != nil {
		return nil, fmt.Errorf("ошибка обновления барбера: %v", err)
	}

	return barber, nil
}

// DeleteBarber удаляет барбера (только админ)
func (s *barberService) DeleteBarber(barberID uint) error {
	// Проверяем, что это барбер
	if !s.roleRepo.HasUserRole(barberID, "barber") {
		return fmt.Errorf("пользователь не является барбером")
	}

	// Удаляем барбера
	if err := s.userRepo.Delete(barberID); err != nil {
		return fmt.Errorf("ошибка удаления барбера: %v", err)
	}

	return nil
}

// GetBarberByID получает барбера по ID (только админ)
func (s *barberService) GetBarberByID(barberID uint) (*models.User, error) {
	barber, err := s.userRepo.GetByID(barberID)
	if err != nil {
		return nil, fmt.Errorf("барбер не найден: %v", err)
	}

	// Проверяем, что это барбер
	if !s.roleRepo.HasUserRole(barberID, "barber") {
		return nil, fmt.Errorf("пользователь не является барбером")
	}

	return barber, nil
}

// GetAllBarbers получает всех барберов (только админ)
func (s *barberService) GetAllBarbers() ([]models.User, error) {
	return s.userRepo.GetByRole("barber")
}

// UpdateBarberSelf обновляет собственный профиль барбера
func (s *barberService) UpdateBarberSelf(barberID uint, req models.BarberSelfUpdateRequest) (*models.User, error) {
	// Получаем барбера
	barber, err := s.userRepo.GetByID(barberID)
	if err != nil {
		return nil, fmt.Errorf("барбер не найден: %v", err)
	}

	// Проверяем, что это барбер
	if !s.roleRepo.HasUserRole(barberID, "barber") {
		return nil, fmt.Errorf("пользователь не является барбером")
	}

	// Обновляем только разрешенные поля
	if req.FirstName != "" {
		barber.FirstName = req.FirstName
	}

	if req.LastName != "" {
		barber.LastName = req.LastName
	}

	if req.Specialties != "" {
		barber.Specialties = req.Specialties
	}

	if req.Experience > 0 {
		barber.Experience = req.Experience
	}

	// Сохраняем изменения
	if err := s.userRepo.Update(barber); err != nil {
		return nil, fmt.Errorf("ошибка обновления профиля: %v", err)
	}

	return barber, nil
}

// GetBarberSelf получает собственный профиль барбера
func (s *barberService) GetBarberSelf(barberID uint) (*models.User, error) {
	barber, err := s.userRepo.GetByID(barberID)
	if err != nil {
		return nil, fmt.Errorf("барбер не найден: %v", err)
	}

	// Проверяем, что это барбер
	if !s.roleRepo.HasUserRole(barberID, "barber") {
		return nil, fmt.Errorf("пользователь не является барбером")
	}

	return barber, nil
}
