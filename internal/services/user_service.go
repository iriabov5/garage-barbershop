package services

import (
	"fmt"
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
	GetAllUsers() ([]models.User, error)
	GetUsersByRole(role string) ([]models.User, error)
	RegisterBarber(telegramID int64, username, firstName, lastName, email string) (*models.User, error)
	RegisterClient(telegramID int64, username, firstName, lastName, email string) (*models.User, error)
}

// userService реализация сервиса пользователей
type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
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
func (s *userService) RegisterBarber(telegramID int64, username, firstName, lastName, email string) (*models.User, error) {
	barber := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		IsActive:   true,
		Rating:     5.0, // Начальный рейтинг
	}

	err := s.userRepo.Create(barber)
	if err != nil {
		return nil, err
	}

	// Назначаем роль "barber"
	barberRole, err := s.roleRepo.GetRoleByName("barber")
	if err != nil {
		return nil, fmt.Errorf("роль barber не найдена: %v", err)
	}
	if err := s.roleRepo.AssignRoleToUser(barber.ID, barberRole.ID, barber.ID); err != nil {
		return nil, fmt.Errorf("ошибка назначения роли: %v", err)
	}

	return barber, nil
}

// RegisterClient регистрирует нового клиента
func (s *userService) RegisterClient(telegramID int64, username, firstName, lastName, email string) (*models.User, error) {
	client := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
	}

	err := s.userRepo.Create(client)
	if err != nil {
		return nil, err
	}

	// Назначаем роль "client"
	clientRole, err := s.roleRepo.GetRoleByName("client")
	if err != nil {
		return nil, fmt.Errorf("роль client не найдена: %v", err)
	}
	if err := s.roleRepo.AssignRoleToUser(client.ID, clientRole.ID, client.ID); err != nil {
		return nil, fmt.Errorf("ошибка назначения роли: %v", err)
	}

	return client, nil
}

// GetAllUsers возвращает всех пользователей
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetUsersByRole возвращает пользователей по роли
func (s *userService) GetUsersByRole(role string) ([]models.User, error) {
	// Используем RoleService для получения пользователей по роли
	roleObj, err := s.roleRepo.GetRoleByName(role)
	if err != nil {
		return nil, fmt.Errorf("роль %s не найдена: %v", role, err)
	}

	return s.roleRepo.GetUsersWithRole(roleObj.ID)
}
