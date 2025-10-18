package repositories

import (
	"fmt"
	"garage-barbershop/internal/models"

	"gorm.io/gorm"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByTelegramID(telegramID int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetBarbers() ([]models.User, error)
	GetClients() ([]models.User, error)
	GetAll() ([]models.User, error)
	GetByRole(role string) ([]models.User, error)
}

// userRepository реализация репозитория пользователей
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByTelegramID получает пользователя по Telegram ID
func (r *userRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update обновляет пользователя
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete удаляет пользователя
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetBarbers получает всех барберов (DEPRECATED - используйте RoleService.GetUsersWithRole)
func (r *userRepository) GetBarbers() ([]models.User, error) {
	// Этот метод больше не работает с новой системой ролей
	// Используйте RoleService.GetUsersWithRole(barberRoleID) вместо этого
	return []models.User{}, fmt.Errorf("GetBarbers deprecated - используйте RoleService.GetUsersWithRole")
}

// GetClients получает всех клиентов (DEPRECATED - используйте RoleService.GetUsersWithRole)
func (r *userRepository) GetClients() ([]models.User, error) {
	// Этот метод больше не работает с новой системой ролей
	// Используйте RoleService.GetUsersWithRole(clientRoleID) вместо этого
	return []models.User{}, fmt.Errorf("GetClients deprecated - используйте RoleService.GetUsersWithRole")
}

// GetAll получает всех пользователей
func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetByRole получает пользователей по роли (использует RoleRepository)
func (r *userRepository) GetByRole(role string) ([]models.User, error) {
	// Этот метод теперь должен работать через RoleRepository
	// Пока возвращаем пустой массив, так как нужен RoleRepository
	return []models.User{}, fmt.Errorf("GetByRole требует RoleRepository - используйте RoleService.GetUsersWithRole")
}
