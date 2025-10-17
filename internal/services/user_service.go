package services

import (
	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"
)

// UserService интерфейс для бизнес-логики пользователей
type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByTelegramID(telegramID int64) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetBarbers() ([]models.User, error)
	GetClients() ([]models.User, error)
	RegisterBarber(telegramID int64, username, firstName, lastName string) (*models.User, error)
	RegisterClient(telegramID int64, username, firstName, lastName string) (*models.User, error)
}

// userService реализация сервиса пользователей
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser создает нового пользователя
func (s *userService) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}

// GetUserByID получает пользователя по ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *userService) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	return s.userRepo.GetByTelegramID(telegramID)
}

// UpdateUser обновляет пользователя
func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUser удаляет пользователя
func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// GetBarbers получает всех барберов
func (s *userService) GetBarbers() ([]models.User, error) {
	return s.userRepo.GetBarbers()
}

// GetClients получает всех клиентов
func (s *userService) GetClients() ([]models.User, error) {
	return s.userRepo.GetClients()
}

// RegisterBarber регистрирует нового барбера
func (s *userService) RegisterBarber(telegramID int64, username, firstName, lastName string) (*models.User, error) {
	barber := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Role:       "barber",
		IsActive:   true,
		Rating:     5.0, // Начальный рейтинг
	}

	err := s.userRepo.Create(barber)
	if err != nil {
		return nil, err
	}

	return barber, nil
}

// RegisterClient регистрирует нового клиента
func (s *userService) RegisterClient(telegramID int64, username, firstName, lastName string) (*models.User, error) {
	client := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Role:       "client",
	}

	err := s.userRepo.Create(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

