package unit

import (
	"testing"

	"garage-barbershop/internal/models"
	"garage-barbershop/internal/services"

	"github.com/stretchr/testify/assert"
)



// Тесты RoleService
func TestRoleService_CreateRole(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	role := &models.Role{
		Name:        "test_role",
		DisplayName: "Test Role",
		Description: "Test role for testing",
		IsActive:    true,
	}

	mockRepo.On("CreateRole", role).Return(nil)

	// Act
	err := roleService.CreateRole(role)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRoleByName(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	expectedRole := &models.Role{
		ID:          1,
		Name:        "admin",
		DisplayName: "Administrator",
		IsActive:    true,
	}

	mockRepo.On("GetRoleByName", "admin").Return(expectedRole, nil)

	// Act
	role, err := roleService.GetRoleByName("admin")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "admin", role.Name)
	assert.Equal(t, "Administrator", role.DisplayName)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_AssignRoleToUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)
	roleID := uint(2)
	assignedBy := uint(3)

	role := &models.Role{
		ID:   roleID,
		Name: "barber",
	}

	mockRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockRepo.On("HasUserRole", userID, "barber").Return(false)
	mockRepo.On("AssignRoleToUser", userID, roleID, assignedBy).Return(nil)

	// Act
	err := roleService.AssignRoleToUser(userID, roleID, assignedBy)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_AssignRoleToUser_AlreadyAssigned(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)
	roleID := uint(2)
	assignedBy := uint(3)

	role := &models.Role{
		ID:   roleID,
		Name: "barber",
	}

	mockRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockRepo.On("HasUserRole", userID, "barber").Return(true)

	// Act
	err := roleService.AssignRoleToUser(userID, roleID, assignedBy)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "роль уже назначена")
	mockRepo.AssertExpectations(t)
}

func TestRoleService_HasUserRole(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)
	roleName := "admin"

	mockRepo.On("HasUserRole", userID, roleName).Return(true)

	// Act
	hasRole := roleService.HasUserRole(userID, roleName)

	// Assert
	assert.True(t, hasRole)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_HasAnyRole(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)

	mockRepo.On("HasUserRole", userID, "admin").Return(false)
	mockRepo.On("HasUserRole", userID, "barber").Return(true)

	// Act
	hasAnyRole := roleService.HasAnyRole(userID, "admin", "barber")

	// Assert
	assert.True(t, hasAnyRole)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_HasAllRoles(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)

	mockRepo.On("HasUserRole", userID, "admin").Return(true)
	mockRepo.On("HasUserRole", userID, "barber").Return(true)

	// Act
	hasAllRoles := roleService.HasAllRoles(userID, "admin", "barber")

	// Assert
	assert.True(t, hasAllRoles)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_IsAdmin(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)

	mockRepo.On("HasUserRole", userID, "admin").Return(true)

	// Act
	isAdmin := roleService.IsAdmin(userID)

	// Assert
	assert.True(t, isAdmin)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_IsBarber(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)

	mockRepo.On("HasUserRole", userID, "barber").Return(true)

	// Act
	isBarber := roleService.IsBarber(userID)

	// Assert
	assert.True(t, isBarber)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_IsClient(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	userID := uint(1)

	mockRepo.On("HasUserRole", userID, "client").Return(true)

	// Act
	isClient := roleService.IsClient(userID)

	// Assert
	assert.True(t, isClient)
	mockRepo.AssertExpectations(t)
}
