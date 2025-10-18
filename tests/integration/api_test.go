package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"garage-barbershop/internal/database"
	"garage-barbershop/internal/handlers"
	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"
	"garage-barbershop/internal/services"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// APITestSuite - набор интеграционных тестов
type APITestSuite struct {
	suite.Suite
	db          *database.Database
	userRepo    repositories.UserRepository
	roleRepo    repositories.RoleRepository
	userService services.UserService
	userHandler *handlers.UserHandler
	server      *httptest.Server
}

// SetupSuite - настройка перед всеми тестами
func (suite *APITestSuite) SetupSuite() {
	// Создаем in-memory SQLite базу для тестов
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		suite.T().Fatal("Failed to connect to test database:", err)
	}

	// Создаем database wrapper
	suite.db = &database.Database{DB: db}

	// Выполняем миграции
	err = suite.db.Migrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Service{},
		&models.Appointment{},
		&models.WorkingHours{},
		&models.Payment{},
		&models.Review{},
	)
	if err != nil {
		suite.T().Fatal("Failed to migrate test database:", err)
	}

	// Создаем зависимости
	suite.userRepo = repositories.NewUserRepository(db)
	suite.roleRepo = repositories.NewRoleRepository(db)
	suite.userService = services.NewUserService(suite.userRepo, suite.roleRepo)
	suite.userHandler = handlers.NewUserHandler(suite.userService)

	// Создаем тестовый HTTP сервер
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", suite.userHandler.GetUsers)
	mux.HandleFunc("/api/users/", suite.userHandler.GetUser)
	mux.HandleFunc("/api/users/create", suite.userHandler.CreateUser)
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	suite.server = httptest.NewServer(mux)
}

// TearDownSuite - очистка после всех тестов
func (suite *APITestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
}

// SetupTest - настройка перед каждым тестом
func (suite *APITestSuite) SetupTest() {
	// Очищаем базу данных перед каждым тестом
	suite.db.DB.Exec("DELETE FROM users")
}

// TestGetUsers_Empty - тест получения пустого списка пользователей
func (suite *APITestSuite) TestGetUsers_Empty() {
	// Act
	resp, err := http.Get(suite.server.URL + "/api/users")

	// Assert
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	users := response["users"].([]interface{})
	suite.Empty(users)
}

// TestCreateUser_Success - тест успешного создания пользователя
func (suite *APITestSuite) TestCreateUser_Success() {
	// Arrange
	userData := map[string]interface{}{
		"telegram_id": 12345,
		"username":    "testuser",
		"first_name":  "John",
		"last_name":   "Doe",
		"role":        "client",
	}

	jsonData, _ := json.Marshal(userData)

	// Act
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	// Assert
	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var response models.User
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	suite.Equal(int64(12345), response.TelegramID)
	suite.Equal("testuser", response.Username)
	suite.Equal("John", response.FirstName)
	suite.Equal("Doe", response.LastName)
	// Роли теперь проверяются через RoleService
	suite.NotZero(response.ID)
}

// TestCreateUser_InvalidData - тест создания пользователя с невалидными данными
func (suite *APITestSuite) TestCreateUser_InvalidData() {
	// Arrange
	userData := map[string]interface{}{
		"telegram_id": "invalid", // Неправильный тип
		"username":    "testuser",
		"first_name":  "John",
		"last_name":   "Doe",
		"role":        "invalid_role", // Неправильная роль
	}

	jsonData, _ := json.Marshal(userData)

	// Act
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	// Assert
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

// TestGetUsers_WithData - тест получения пользователей с данными
func (suite *APITestSuite) TestGetUsers_WithData() {
	// Arrange - создаем тестовых пользователей
	barber := &models.User{
		TelegramID: 11111,
		Username:   "barber1",
		FirstName:  "Ivan",
		LastName:   "Barber",
		Email:      "barber1@example.com",
		IsActive:   true,
	}

	client := &models.User{
		TelegramID: 22222,
		Username:   "client1",
		FirstName:  "Jane",
		LastName:   "Client",
		Email:      "client1@example.com",
	}

	suite.userRepo.Create(barber)
	suite.userRepo.Create(client)

	// Act
	resp, err := http.Get(suite.server.URL + "/api/users")

	// Assert
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	users := response["users"].([]interface{})
	suite.Len(users, 2)
}

// TestGetUsers_ByRole - тест получения пользователей по роли
func (suite *APITestSuite) TestGetUsers_ByRole() {
	// Arrange - создаем тестовых пользователей
	barber := &models.User{
		TelegramID: 11111,
		Username:   "barber1",
		FirstName:  "Ivan",
		LastName:   "Barber",
		Email:      "barber1@example.com",
		IsActive:   true,
	}

	client := &models.User{
		TelegramID: 22222,
		Username:   "client1",
		FirstName:  "Jane",
		LastName:   "Client",
		Email:      "client1@example.com",
		IsActive:   true,
	}

	// Создаем пользователей
	err := suite.db.DB.Create(barber).Error
	suite.Require().NoError(err)
	err = suite.db.DB.Create(client).Error
	suite.Require().NoError(err)

	// Назначаем роли (используем существующие роли из миграции)
	barberRole, err := suite.roleRepo.GetRoleByName("barber")
	suite.Require().NoError(err)
	clientRole, err := suite.roleRepo.GetRoleByName("client")
	suite.Require().NoError(err)

	err = suite.roleRepo.AssignRoleToUser(barber.ID, barberRole.ID, barber.ID)
	suite.Require().NoError(err)
	err = suite.roleRepo.AssignRoleToUser(client.ID, clientRole.ID, client.ID)
	suite.Require().NoError(err)

	// Act - запрашиваем пользователей с ролью "barber"
	resp, err := http.Get(suite.server.URL + "/api/users?role=barber")
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	users := response["users"].([]interface{})
	suite.Len(users, 1)

	// Проверяем, что вернулся только барбер
	userData := users[0].(map[string]interface{})
	suite.Equal("Ivan", userData["first_name"])
	suite.Equal("Barber", userData["last_name"])
}

// TestAPIStatus - тест статуса API
func (suite *APITestSuite) TestAPIStatus() {
	// Act
	resp, err := http.Get(suite.server.URL + "/api/status")

	// Assert
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	suite.Equal("ok", response["status"])
}

// TestUserService_RegisterBarber_Integration - интеграционный тест регистрации барбера
func (suite *APITestSuite) TestUserService_RegisterBarber_Integration() {
	// Act
	barber, err := suite.userService.RegisterBarber(12345, "barber_user", "Ivan", "Barber", "barber@example.com")

	// Assert
	suite.NoError(err)
	suite.NotNil(barber)
	suite.Equal(int64(12345), barber.TelegramID)
	suite.Equal("barber_user", barber.Username)
	suite.Equal("Ivan", barber.FirstName)
	suite.Equal("Barber", barber.LastName)
	// Роли теперь проверяются через RoleService
	suite.True(barber.IsActive)
	suite.Equal(5.0, barber.Rating)

	// Проверяем, что барбер сохранился в базе
	savedBarber, err := suite.userRepo.GetByID(barber.ID)
	suite.NoError(err)
	suite.Equal(barber.ID, savedBarber.ID)
	suite.Equal(barber.TelegramID, savedBarber.TelegramID)
}

// Запуск тестов
func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
