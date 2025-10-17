package e2e

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

// UserJourneyTestSuite - E2E тесты пользовательских сценариев
type UserJourneyTestSuite struct {
	suite.Suite
	db          *database.Database
	userRepo    repositories.UserRepository
	userService services.UserService
	userHandler *handlers.UserHandler
	server      *httptest.Server
}

// SetupSuite - настройка перед всеми тестами
func (suite *UserJourneyTestSuite) SetupSuite() {
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
	suite.userService = services.NewUserService(suite.userRepo)
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
func (suite *UserJourneyTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
}

// SetupTest - настройка перед каждым тестом
func (suite *UserJourneyTestSuite) SetupTest() {
	// Очищаем базу данных перед каждым тестом
	suite.db.DB.Exec("DELETE FROM users")
}

// TestCompleteUserJourney - полный пользовательский сценарий
func (suite *UserJourneyTestSuite) TestCompleteUserJourney() {
	// Сценарий: Регистрация барбера → Регистрация клиента → Создание услуги → Запись на услугу

	// 1. Регистрация барбера
	barberData := map[string]interface{}{
		"telegram_id": 11111,
		"username":    "ivan_barber",
		"first_name":  "Ivan",
		"last_name":   "Barber",
		"role":        "barber",
	}

	jsonData, _ := json.Marshal(barberData)
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var barber models.User
	err = json.NewDecoder(resp.Body).Decode(&barber)
	suite.NoError(err)
	suite.Equal("barber", barber.Role)
	suite.True(barber.IsActive)

	// 2. Регистрация клиента
	clientData := map[string]interface{}{
		"telegram_id": 22222,
		"username":    "jane_client",
		"first_name":  "Jane",
		"last_name":   "Client",
		"role":        "client",
	}

	jsonData, _ = json.Marshal(clientData)
	resp, err = http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var client models.User
	err = json.NewDecoder(resp.Body).Decode(&client)
	suite.NoError(err)
	suite.Equal("client", client.Role)

	// 3. Проверяем, что оба пользователя созданы
	resp, err = http.Get(suite.server.URL + "/api/users")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	users := response["users"].([]interface{})
	suite.Len(users, 2)

	// 4. Проверяем фильтрацию по ролям
	resp, err = http.Get(suite.server.URL + "/api/users?role=barber")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	barbers := response["users"].([]interface{})
	suite.Len(barbers, 1)
	suite.Equal("barber", barbers[0].(map[string]interface{})["role"])

	// 5. Проверяем фильтрацию клиентов
	resp, err = http.Get(suite.server.URL + "/api/users?role=client")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	clients := response["users"].([]interface{})
	suite.Len(clients, 1)
	suite.Equal("client", clients[0].(map[string]interface{})["role"])
}

// TestBarberRegistrationFlow - сценарий регистрации барбера
func (suite *UserJourneyTestSuite) TestBarberRegistrationFlow() {
	// 1. Регистрация барбера через API
	barberData := map[string]interface{}{
		"telegram_id": 33333,
		"username":    "master_barber",
		"first_name":  "Master",
		"last_name":   "Barber",
		"role":        "barber",
	}

	jsonData, _ := json.Marshal(barberData)
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var barber models.User
	err = json.NewDecoder(resp.Body).Decode(&barber)
	suite.NoError(err)

	// 2. Проверяем, что барбер создан с правильными параметрами
	suite.Equal(int64(33333), barber.TelegramID)
	suite.Equal("master_barber", barber.Username)
	suite.Equal("Master", barber.FirstName)
	suite.Equal("Barber", barber.LastName)
	suite.Equal("barber", barber.Role)
	suite.True(barber.IsActive)
	suite.Equal(5.0, barber.Rating)

	// 3. Проверяем, что барбер сохранился в базе
	savedBarber, err := suite.userRepo.GetByID(barber.ID)
	suite.NoError(err)
	suite.Equal(barber.ID, savedBarber.ID)
	suite.Equal(barber.TelegramID, savedBarber.TelegramID)

	// 4. Проверяем, что барбер появляется в списке барберов
	barbers, err := suite.userRepo.GetBarbers()
	suite.NoError(err)
	suite.Len(barbers, 1)
	suite.Equal(barber.ID, barbers[0].ID)
}

// TestClientRegistrationFlow - сценарий регистрации клиента
func (suite *UserJourneyTestSuite) TestClientRegistrationFlow() {
	// 1. Регистрация клиента через API
	clientData := map[string]interface{}{
		"telegram_id": 44444,
		"username":    "regular_client",
		"first_name":  "Regular",
		"last_name":   "Client",
		"role":        "client",
	}

	jsonData, _ := json.Marshal(clientData)
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var client models.User
	err = json.NewDecoder(resp.Body).Decode(&client)
	suite.NoError(err)

	// 2. Проверяем, что клиент создан с правильными параметрами
	suite.Equal(int64(44444), client.TelegramID)
	suite.Equal("regular_client", client.Username)
	suite.Equal("Regular", client.FirstName)
	suite.Equal("Client", client.LastName)
	suite.Equal("client", client.Role)

	// 3. Проверяем, что клиент сохранился в базе
	savedClient, err := suite.userRepo.GetByID(client.ID)
	suite.NoError(err)
	suite.Equal(client.ID, savedClient.ID)
	suite.Equal(client.TelegramID, savedClient.TelegramID)

	// 4. Проверяем, что клиент появляется в списке клиентов
	clients, err := suite.userRepo.GetClients()
	suite.NoError(err)
	suite.Len(clients, 1)
	suite.Equal(client.ID, clients[0].ID)
}

// TestErrorHandling - тест обработки ошибок
func (suite *UserJourneyTestSuite) TestErrorHandling() {
	// 1. Тест создания пользователя с невалидными данными
	invalidData := map[string]interface{}{
		"telegram_id": "invalid", // Неправильный тип
		"username":    "test",
		"first_name":  "Test",
		"last_name":   "User",
		"role":        "invalid_role", // Неправильная роль
	}

	jsonData, _ := json.Marshal(invalidData)
	resp, err := http.Post(
		suite.server.URL+"/api/users/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	// 2. Тест получения несуществующего пользователя
	resp, err = http.Get(suite.server.URL + "/api/users/99999")
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestAPIStatus - тест статуса API
func (suite *UserJourneyTestSuite) TestAPIStatus() {
	resp, err := http.Get(suite.server.URL + "/api/status")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Equal("ok", response["status"])
}

// Запуск E2E тестов
func TestUserJourneyTestSuite(t *testing.T) {
	suite.Run(t, new(UserJourneyTestSuite))
}
