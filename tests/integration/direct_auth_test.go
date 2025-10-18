package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"garage-barbershop/internal/database"
	"garage-barbershop/internal/handlers"
	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"
	"garage-barbershop/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DirectAuthTestSuite набор тестов для прямой авторизации
type DirectAuthTestSuite struct {
	suite.Suite
	db          *database.Database
	authService services.AuthService
	authHandler *handlers.AuthHTTPHandler
	router      *gin.Engine
}

// SetupSuite инициализирует тестовую среду
func (suite *DirectAuthTestSuite) SetupSuite() {
	// Инициализируем in-memory SQLite базу данных
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	suite.Require().NoError(err)

	// Создаем тестовую БД
	testDB := &database.Database{DB: db}

	// Выполняем миграции
	err = testDB.Migrate(&models.User{}, &models.Role{}, &models.UserRole{})
	suite.Require().NoError(err)

	suite.db = testDB

	// Создаем сервисы (Redis = nil для упрощения)
	userRepo := repositories.NewUserRepository(suite.db.DB)
	roleRepo := repositories.NewRoleRepository(suite.db.DB)
	suite.authService = services.NewAuthService(userRepo, roleRepo, nil, "test_secret", "test_bot_token")
	suite.authHandler = handlers.NewAuthHTTPHandler(suite.authService)

	// Настраиваем Gin роутер
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Добавляем маршруты
	suite.router.POST("/api/auth/register", func(c *gin.Context) {
		suite.authHandler.RegisterDirect(c.Writer, c.Request)
	})
	suite.router.POST("/api/auth/login", func(c *gin.Context) {
		suite.authHandler.LoginDirect(c.Writer, c.Request)
	})
}

// TearDownSuite очищает тестовую среду
func (suite *DirectAuthTestSuite) TearDownSuite() {
	sqlDB, err := suite.db.DB.DB()
	suite.Require().NoError(err)
	sqlDB.Close()
}

// SetupTest очищает данные перед каждым тестом
func (suite *DirectAuthTestSuite) SetupTest() {
	// Очищаем таблицу пользователей
	suite.db.DB.Exec("DELETE FROM users")
}

// TestDirectRegister_Success тестирует успешную регистрацию
func (suite *DirectAuthTestSuite) TestDirectRegister_Success() {
	// Arrange
	registerData := models.DirectRegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
	}

	jsonData, err := json.Marshal(registerData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert
	suite.Equal(http.StatusOK, w.Code)

	var authResponse models.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	suite.NoError(err)

	// Проверяем, что токены сгенерированы
	suite.NotEmpty(authResponse.AccessToken)
	suite.NotEmpty(authResponse.RefreshToken)
	suite.Equal(int64(900), authResponse.ExpiresIn)

	// Проверяем данные пользователя
	suite.Equal("test@example.com", authResponse.User.Email)
	suite.Equal("John", authResponse.User.FirstName)
	suite.Equal("Doe", authResponse.User.LastName)
	// Роли теперь проверяются через RoleService
	suite.Equal("direct", authResponse.User.AuthMethod)
	suite.True(authResponse.User.IsActive)

	// Проверяем, что пользователь создан в БД
	var user models.User
	err = suite.db.DB.Where("email = ?", "test@example.com").First(&user).Error
	suite.NoError(err)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("direct", user.AuthMethod)
}

// TestDirectRegister_DuplicateEmail тестирует регистрацию с дублирующимся email
func (suite *DirectAuthTestSuite) TestDirectRegister_DuplicateEmail() {
	// Arrange - создаем пользователя заранее
	existingUser := &models.User{
		Email:      "existing@example.com",
		FirstName:  "Existing",
		LastName:   "User",
		AuthMethod: "direct",
		IsActive:   true,
	}
	err := suite.db.DB.Create(existingUser).Error
	suite.Require().NoError(err)

	// Пытаемся зарегистрировать с тем же email
	registerData := models.DirectRegisterRequest{
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
	}

	jsonData, err := json.Marshal(registerData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert
	suite.Equal(http.StatusBadRequest, w.Code)
}

// TestDirectLogin_Success тестирует успешную авторизацию
func (suite *DirectAuthTestSuite) TestDirectLogin_Success() {
	// Arrange - создаем пользователя заранее
	user := &models.User{
		Email:      "login@example.com",
		FirstName:  "Login",
		LastName:   "User",
		AuthMethod: "direct",
		IsActive:   true,
	}

	// Хешируем пароль
	passwordHash, err := suite.authService.HashPassword("password123")
	suite.Require().NoError(err)
	user.PasswordHash = passwordHash

	err = suite.db.DB.Create(user).Error
	suite.Require().NoError(err)

	// Данные для авторизации
	loginData := models.DirectLoginRequest{
		Email:    "login@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(loginData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert
	suite.Equal(http.StatusOK, w.Code)

	var authResponse models.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	suite.NoError(err)

	// Проверяем, что токены сгенерированы
	suite.NotEmpty(authResponse.AccessToken)
	suite.NotEmpty(authResponse.RefreshToken)

	// Проверяем данные пользователя
	suite.Equal("login@example.com", authResponse.User.Email)
	suite.Equal("Login", authResponse.User.FirstName)
	suite.Equal("User", authResponse.User.LastName)
	// Роли теперь проверяются через RoleService
	suite.Equal("direct", authResponse.User.AuthMethod)
}

// TestDirectLogin_WrongPassword тестирует авторизацию с неверным паролем
func (suite *DirectAuthTestSuite) TestDirectLogin_WrongPassword() {
	// Arrange - создаем пользователя заранее
	user := &models.User{
		Email:      "wrongpass@example.com",
		FirstName:  "Wrong",
		LastName:   "Password",
		AuthMethod: "direct",
		IsActive:   true,
	}

	// Хешируем правильный пароль
	passwordHash, err := suite.authService.HashPassword("correctpassword")
	suite.Require().NoError(err)
	user.PasswordHash = passwordHash

	err = suite.db.DB.Create(user).Error
	suite.Require().NoError(err)

	// Пытаемся авторизоваться с неверным паролем
	loginData := models.DirectLoginRequest{
		Email:    "wrongpass@example.com",
		Password: "wrongpassword",
	}

	jsonData, err := json.Marshal(loginData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert
	suite.Equal(http.StatusUnauthorized, w.Code)
}

// TestDirectAuthTestSuite запускает все тесты
func TestDirectAuthTestSuite(t *testing.T) {
	suite.Run(t, new(DirectAuthTestSuite))
}
