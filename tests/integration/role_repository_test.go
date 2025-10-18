package integration

import (
	"testing"

	"garage-barbershop/internal/database"
	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB создает тестовую базу данных
func setupTestDB(t *testing.T) *database.Database {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	testDB := &database.Database{DB: db}

	// Выполняем миграции
	err = testDB.Migrate(&models.User{}, &models.Role{}, &models.UserRole{})
	require.NoError(t, err)

	return testDB
}

// cleanupTestDB очищает тестовую базу данных
func cleanupTestDB(t *testing.T, db *database.Database) {
	if db != nil && db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func TestRoleRepository_CreateAndGetRole(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)

	role := &models.Role{
		Name:        "test_role_1",
		DisplayName: "Test Role 1",
		Description: "Test role for testing",
		IsActive:    true,
	}

	// Act
	err := roleRepo.CreateRole(role)
	require.NoError(t, err)

	// Assert
	retrievedRole, err := roleRepo.GetRoleByID(role.ID)
	require.NoError(t, err)
	assert.Equal(t, role.Name, retrievedRole.Name)
	assert.Equal(t, role.DisplayName, retrievedRole.DisplayName)
	assert.Equal(t, role.Description, retrievedRole.Description)
	assert.True(t, retrievedRole.IsActive)
}

func TestRoleRepository_GetRoleByName(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)

	// Используем существующую роль "admin" (созданную миграцией)
	// Act
	retrievedRole, err := roleRepo.GetRoleByName("admin")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "admin", retrievedRole.Name)
	assert.Equal(t, "Администратор", retrievedRole.DisplayName)
}

func TestRoleRepository_AssignRoleToUser(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)
	userRepo := repositories.NewUserRepository(db.DB)

	// Создаем пользователя
	user := &models.User{
		TelegramID: 12345,
		FirstName:  "Test",
		LastName:   "User",
		Username:   "testuser",
		AuthMethod: "telegram",
	}

	err := userRepo.Create(user)
	require.NoError(t, err)

	// Создаем роль
	role := &models.Role{
		Name:        "barber_test",
		DisplayName: "Barber",
		IsActive:    true,
	}

	err = roleRepo.CreateRole(role)
	require.NoError(t, err)

	// Act
	err = roleRepo.AssignRoleToUser(user.ID, role.ID, user.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что роль назначена
	hasRole := roleRepo.HasUserRole(user.ID, "barber_test")
	assert.True(t, hasRole)

	// Проверяем роли пользователя
	userRoles, err := roleRepo.GetUserRoles(user.ID)
	require.NoError(t, err)
	assert.Len(t, userRoles, 1)
	assert.Equal(t, "barber_test", userRoles[0].Name)
}

func TestRoleRepository_RemoveRoleFromUser(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)
	userRepo := repositories.NewUserRepository(db.DB)

	// Создаем пользователя
	user := &models.User{
		TelegramID: 12345,
		FirstName:  "Test",
		LastName:   "User",
		Username:   "testuser",
		AuthMethod: "telegram",
	}

	err := userRepo.Create(user)
	require.NoError(t, err)

	// Создаем роль
	role := &models.Role{
		Name:        "barber_test",
		DisplayName: "Barber",
		IsActive:    true,
	}

	err = roleRepo.CreateRole(role)
	require.NoError(t, err)

	// Назначаем роль
	err = roleRepo.AssignRoleToUser(user.ID, role.ID, user.ID)
	require.NoError(t, err)

	// Act
	err = roleRepo.RemoveRoleFromUser(user.ID, role.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что роль снята
	hasRole := roleRepo.HasUserRole(user.ID, "barber")
	assert.False(t, hasRole)
}

func TestRoleRepository_GetUsersWithRole(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)
	userRepo := repositories.NewUserRepository(db.DB)

	// Создаем пользователей
	user1 := &models.User{
		TelegramID: 12345,
		FirstName:  "User1",
		LastName:   "Test",
		Username:   "user1",
		Email:      "user1@test.com",
		AuthMethod: "telegram",
	}

	user2 := &models.User{
		TelegramID: 67890,
		FirstName:  "User2",
		LastName:   "Test",
		Username:   "user2",
		Email:      "user2@test.com",
		AuthMethod: "telegram",
	}

	err := userRepo.Create(user1)
	require.NoError(t, err)
	err = userRepo.Create(user2)
	require.NoError(t, err)

	// Создаем роль
	role := &models.Role{
		Name:        "barber_test",
		DisplayName: "Barber",
		IsActive:    true,
	}

	err = roleRepo.CreateRole(role)
	require.NoError(t, err)

	// Назначаем роль обоим пользователям
	err = roleRepo.AssignRoleToUser(user1.ID, role.ID, user1.ID)
	require.NoError(t, err)
	err = roleRepo.AssignRoleToUser(user2.ID, role.ID, user2.ID)
	require.NoError(t, err)

	// Act
	users, err := roleRepo.GetUsersWithRole(role.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestRoleRepository_GetUserWithRoles(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	roleRepo := repositories.NewRoleRepository(db.DB)
	userRepo := repositories.NewUserRepository(db.DB)

	// Создаем пользователя
	user := &models.User{
		TelegramID: 12345,
		FirstName:  "Test",
		LastName:   "User",
		Username:   "testuser",
		AuthMethod: "telegram",
	}

	err := userRepo.Create(user)
	require.NoError(t, err)

	// Создаем роли
	barberRole := &models.Role{
		Name:        "barber_test",
		DisplayName: "Barber",
		IsActive:    true,
	}

	clientRole := &models.Role{
		Name:        "client_test",
		DisplayName: "Client",
		IsActive:    true,
	}

	err = roleRepo.CreateRole(barberRole)
	require.NoError(t, err)
	err = roleRepo.CreateRole(clientRole)
	require.NoError(t, err)

	// Назначаем обе роли
	err = roleRepo.AssignRoleToUser(user.ID, barberRole.ID, user.ID)
	require.NoError(t, err)
	err = roleRepo.AssignRoleToUser(user.ID, clientRole.ID, user.ID)
	require.NoError(t, err)

	// Act
	userWithRoles, err := roleRepo.GetUserWithRoles(user.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, userWithRoles.User.ID)
	assert.Len(t, userWithRoles.Roles, 2)

	// Проверяем, что роли содержат ожидаемые имена
	roleNames := make([]string, len(userWithRoles.Roles))
	for i, role := range userWithRoles.Roles {
		roleNames[i] = role.Name
	}
	assert.Contains(t, roleNames, "barber_test")
	assert.Contains(t, roleNames, "client_test")
}
