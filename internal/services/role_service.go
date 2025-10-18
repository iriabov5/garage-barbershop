package services

import (
	"fmt"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"
)

// RoleService интерфейс для управления ролями
type RoleService interface {
	// Управление ролями
	CreateRole(role *models.Role) error
	GetRoleByID(id uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id uint) error

	// Управление ролями пользователей
	AssignRoleToUser(userID, roleID uint, assignedBy uint) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	GetUsersWithRole(roleID uint) ([]models.User, error)
	HasUserRole(userID uint, roleName string) bool
	GetUserWithRoles(userID uint) (*models.UserWithRoles, error)
	GetAllUsersWithRoles() ([]models.UserWithRoles, error)

	// Проверка разрешений
	HasAnyRole(userID uint, roleNames ...string) bool
	HasAllRoles(userID uint, roleNames ...string) bool
	IsAdmin(userID uint) bool
	IsBarber(userID uint) bool
	IsClient(userID uint) bool
}

// roleService реализация RoleService
type roleService struct {
	roleRepo repositories.RoleRepository
}

// NewRoleService создает новый экземпляр RoleService
func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

// CreateRole создает новую роль
func (s *roleService) CreateRole(role *models.Role) error {
	return s.roleRepo.CreateRole(role)
}

// GetRoleByID получает роль по ID
func (s *roleService) GetRoleByID(id uint) (*models.Role, error) {
	return s.roleRepo.GetRoleByID(id)
}

// GetRoleByName получает роль по имени
func (s *roleService) GetRoleByName(name string) (*models.Role, error) {
	return s.roleRepo.GetRoleByName(name)
}

// GetAllRoles получает все роли
func (s *roleService) GetAllRoles() ([]models.Role, error) {
	return s.roleRepo.GetAllRoles()
}

// UpdateRole обновляет роль
func (s *roleService) UpdateRole(role *models.Role) error {
	return s.roleRepo.UpdateRole(role)
}

// DeleteRole удаляет роль
func (s *roleService) DeleteRole(id uint) error {
	return s.roleRepo.DeleteRole(id)
}

// AssignRoleToUser назначает роль пользователю
func (s *roleService) AssignRoleToUser(userID, roleID uint, assignedBy uint) error {
	// Проверяем, что роль не назначена уже
	role, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("роль не найдена: %v", err)
	}
	if s.roleRepo.HasUserRole(userID, role.Name) {
		return fmt.Errorf("роль уже назначена пользователю")
	}

	return s.roleRepo.AssignRoleToUser(userID, roleID, assignedBy)
}

// RemoveRoleFromUser снимает роль с пользователя
func (s *roleService) RemoveRoleFromUser(userID, roleID uint) error {
	return s.roleRepo.RemoveRoleFromUser(userID, roleID)
}

// GetUserRoles получает роли пользователя
func (s *roleService) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(userID)
}

// GetUsersWithRole получает пользователей с определенной ролью
func (s *roleService) GetUsersWithRole(roleID uint) ([]models.User, error) {
	return s.roleRepo.GetUsersWithRole(roleID)
}

// HasUserRole проверяет, есть ли у пользователя указанная роль
func (s *roleService) HasUserRole(userID uint, roleName string) bool {
	return s.roleRepo.HasUserRole(userID, roleName)
}

// GetUserWithRoles получает пользователя с его ролями
func (s *roleService) GetUserWithRoles(userID uint) (*models.UserWithRoles, error) {
	return s.roleRepo.GetUserWithRoles(userID)
}

// GetAllUsersWithRoles получает всех пользователей с их ролями
func (s *roleService) GetAllUsersWithRoles() ([]models.UserWithRoles, error) {
	return s.roleRepo.GetAllUsersWithRoles()
}

// HasAnyRole проверяет, есть ли у пользователя хотя бы одна из указанных ролей
func (s *roleService) HasAnyRole(userID uint, roleNames ...string) bool {
	for _, roleName := range roleNames {
		if s.roleRepo.HasUserRole(userID, roleName) {
			return true
		}
	}
	return false
}

// HasAllRoles проверяет, есть ли у пользователя все указанные роли
func (s *roleService) HasAllRoles(userID uint, roleNames ...string) bool {
	for _, roleName := range roleNames {
		if !s.roleRepo.HasUserRole(userID, roleName) {
			return false
		}
	}
	return true
}

// IsAdmin проверяет, является ли пользователь админом
func (s *roleService) IsAdmin(userID uint) bool {
	return s.roleRepo.HasUserRole(userID, "admin")
}

// IsBarber проверяет, является ли пользователь барбером
func (s *roleService) IsBarber(userID uint) bool {
	return s.roleRepo.HasUserRole(userID, "barber")
}

// IsClient проверяет, является ли пользователь клиентом
func (s *roleService) IsClient(userID uint) bool {
	return s.roleRepo.HasUserRole(userID, "client")
}
