package repositories

import (
	"time"

	"garage-barbershop/internal/models"

	"gorm.io/gorm"
)

// RoleRepository интерфейс для работы с ролями
type RoleRepository interface {
	// Управление ролями
	CreateRole(role *models.Role) error
	GetRoleByID(id uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error

	// Управление связями пользователь-роль
	AssignRoleToUser(userID, roleID uint, assignedBy uint) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	GetUsersWithRole(roleID uint) ([]models.User, error)
	GetUserRole(userID, roleID uint) (*models.UserRole, error)
	HasUserRole(userID uint, roleName string) bool

	// Получение пользователей с ролями
	GetUserWithRoles(userID uint) (*models.UserWithRoles, error)
	GetAllUsersWithRoles() ([]models.UserWithRoles, error)
}

// roleRepository реализация репозитория ролей
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository создает новый репозиторий ролей
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// CreateRole создает новую роль
func (r *roleRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

// GetRoleByID получает роль по ID
func (r *roleRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName получает роль по имени
func (r *roleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles получает все роли
func (r *roleRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// UpdateRole обновляет роль
func (r *roleRepository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

// DeleteRole удаляет роль
func (r *roleRepository) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// AssignRoleToUser назначает роль пользователю
func (r *roleRepository) AssignRoleToUser(userID, roleID uint, assignedBy uint) error {
	userRole := &models.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
		IsActive:   1, // 1 = true
	}
	return r.db.Create(userRole).Error
}

// RemoveRoleFromUser снимает роль с пользователя
func (r *roleRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// GetUserRoles получает роли пользователя
func (r *roleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var userRoles []models.UserRole
	err := r.db.Where("user_id = ? AND is_active = ?", userID, 1).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []models.Role{}, nil
	}

	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	var roles []models.Role
	err = r.db.Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}

// GetUsersWithRole получает пользователей с определенной ролью
func (r *roleRepository) GetUsersWithRole(roleID uint) ([]models.User, error) {
	var userRoles []models.UserRole
	err := r.db.Where("role_id = ? AND is_active = ?", roleID, 1).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []models.User{}, nil
	}

	var userIDs []uint
	for _, userRole := range userRoles {
		userIDs = append(userIDs, userRole.UserID)
	}

	var users []models.User
	err = r.db.Where("id IN ?", userIDs).Find(&users).Error
	return users, err
}

// GetUserRole получает связь пользователь-роль
func (r *roleRepository) GetUserRole(userID, roleID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

// HasUserRole проверяет, есть ли у пользователя указанная роль
func (r *roleRepository) HasUserRole(userID uint, roleName string) bool {
	// Сначала получаем роль по имени
	var role models.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		return false
	}

	// Проверяем, есть ли связь пользователь-роль
	var count int64
	err = r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ? AND is_active = ?", userID, role.ID, 1).
		Count(&count).Error
	return err == nil && count > 0
}

// GetUserWithRoles получает пользователя с его ролями
func (r *roleRepository) GetUserWithRoles(userID uint) (*models.UserWithRoles, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &models.UserWithRoles{
		User:  user,
		Roles: user.Roles,
	}, nil
}

// GetAllUsersWithRoles получает всех пользователей с их ролями
func (r *roleRepository) GetAllUsersWithRoles() ([]models.UserWithRoles, error) {
	var users []models.User
	err := r.db.Preload("Roles").Find(&users).Error
	if err != nil {
		return nil, err
	}

	var usersWithRoles []models.UserWithRoles
	for _, user := range users {
		usersWithRoles = append(usersWithRoles, models.UserWithRoles{
			User:  user,
			Roles: user.Roles,
		})
	}

	return usersWithRoles, nil
}
