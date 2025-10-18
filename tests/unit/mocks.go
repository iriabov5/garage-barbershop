package unit

import (
	"garage-barbershop/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	args := m.Called(telegramID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetBarbers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetClients() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(role string) ([]models.User, error) {
	args := m.Called(role)
	return args.Get(0).([]models.User), args.Error(1)
}

// MockRoleRepository для тестирования
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) CreateRole(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByName(name string) (*models.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetAllRoles() ([]models.Role, error) {
	args := m.Called()
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateRole(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRole(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignRoleToUser(userID, roleID, assignedBy uint) error {
	args := m.Called(userID, roleID, assignedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveRoleFromUser(userID, roleID uint) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetUsersWithRole(roleID uint) ([]models.User, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockRoleRepository) HasUserRole(userID uint, roleName string) bool {
	args := m.Called(userID, roleName)
	return args.Bool(0)
}

func (m *MockRoleRepository) GetUserRole(userID, roleID uint) (*models.UserRole, error) {
	args := m.Called(userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRole), args.Error(1)
}

func (m *MockRoleRepository) GetUserWithRoles(userID uint) (*models.UserWithRoles, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserWithRoles), args.Error(1)
}

func (m *MockRoleRepository) GetAllUsersWithRoles() ([]models.UserWithRoles, error) {
	args := m.Called()
	return args.Get(0).([]models.UserWithRoles), args.Error(1)
}
