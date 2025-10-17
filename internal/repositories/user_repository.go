package repositories

import (
	"garage-barbershop/internal/models"

	"gorm.io/gorm"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByTelegramID(telegramID int64) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetBarbers() ([]models.User, error)
	GetClients() ([]models.User, error)
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

// Update обновляет пользователя
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete удаляет пользователя
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetBarbers получает всех барберов
func (r *userRepository) GetBarbers() ([]models.User, error) {
	var barbers []models.User
	err := r.db.Where("role = ?", "barber").Find(&barbers).Error
	return barbers, err
}

// GetClients получает всех клиентов
func (r *userRepository) GetClients() ([]models.User, error) {
	var clients []models.User
	err := r.db.Where("role = ?", "client").Find(&clients).Error
	return clients, err
}

