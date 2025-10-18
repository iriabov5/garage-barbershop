package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"garage-barbershop/internal/database"
	"garage-barbershop/internal/handlers"
	"garage-barbershop/internal/models"
	"garage-barbershop/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TelegramAuthTestSuite набор тестов для Telegram аутентификации
type TelegramAuthTestSuite struct {
	suite.Suite
	db          *database.Database
	authService *TestAuthService
	authHandler *handlers.AuthHandler
	router      *gin.Engine
}

// SetupSuite инициализирует тестовую среду
func (suite *TelegramAuthTestSuite) SetupSuite() {
	// Инициализируем in-memory SQLite базу данных
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	suite.Require().NoError(err)

	// Создаем тестовую БД
	testDB := &database.Database{DB: db}

	// Выполняем миграции
	err = testDB.Migrate(&models.User{})
	suite.Require().NoError(err)

	suite.db = testDB

	// Создаем тестовые сервисы (Redis = nil для упрощения)
	userRepo := repositories.NewUserRepository(suite.db.DB)
	testAuthService := NewTestAuthService(userRepo, nil, "test_secret", "test_bot_token")
	suite.authService = testAuthService
	suite.authHandler = handlers.NewAuthHandler(testAuthService)

	// Настраиваем Gin роутер
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Добавляем маршруты
	suite.router.POST("/api/auth/telegram", suite.authHandler.TelegramAuth)
}

// TearDownSuite очищает тестовую среду
func (suite *TelegramAuthTestSuite) TearDownSuite() {
	sqlDB, err := suite.db.DB.DB()
	suite.Require().NoError(err)
	sqlDB.Close()
}

// SetupTest очищает данные перед каждым тестом
func (suite *TelegramAuthTestSuite) SetupTest() {
	// Очищаем таблицу пользователей
	suite.db.DB.Exec("DELETE FROM users")
}

// TestTelegramAuth_Success тестирует успешную аутентификацию через Telegram
func (suite *TelegramAuthTestSuite) TestTelegramAuth_Success() {
	// Arrange - подготавливаем данные для Telegram аутентификации
	authData := models.TelegramAuthData{
		ID:        12345,
		Username:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		AuthDate:  time.Now().Unix(),
		Hash:      "test_hash", // В реальном тесте нужна правильная подпись
	}

	jsonData, err := json.Marshal(authData)
	suite.Require().NoError(err)

	// Act - отправляем POST запрос на аутентификацию
	req := httptest.NewRequest("POST", "/api/auth/telegram", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert - проверяем результат
	suite.Equal(http.StatusOK, w.Code)

	var authResponse models.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	suite.NoError(err)

	// Проверяем, что токены сгенерированы
	suite.NotEmpty(authResponse.AccessToken)
	suite.NotEmpty(authResponse.RefreshToken)
	suite.Equal(int64(900), authResponse.ExpiresIn) // 15 минут

	// Проверяем данные пользователя
	suite.Equal("testuser", authResponse.User.Username)
	suite.Equal("John", authResponse.User.FirstName)
	suite.Equal("Doe", authResponse.User.LastName)
	suite.Equal("client", authResponse.User.Role) // По умолчанию клиент
	suite.True(authResponse.User.IsActive)

	// Проверяем, что пользователь создан в БД
	var user models.User
	err = suite.db.DB.Where("telegram_id = ?", 12345).First(&user).Error
	suite.NoError(err)
	suite.Equal("testuser", user.Username)
}

// TestTelegramAuth_InvalidData тестирует аутентификацию с неверными данными
func (suite *TelegramAuthTestSuite) TestTelegramAuth_InvalidData() {
	// Arrange - неверные данные (отсутствует FirstName)
	authData := models.TelegramAuthData{
		ID:        12345,
		Username:  "testuser",
		FirstName: "", // Отсутствует имя - невалидные данные
		LastName:  "Doe",
		AuthDate:  time.Now().Unix(),
		Hash:      "test_hash",
	}

	jsonData, err := json.Marshal(authData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/telegram", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert - должен вернуть ошибку валидации
	suite.Equal(http.StatusUnauthorized, w.Code)
}

// TestTelegramAuth_ExistingUser тестирует аутентификацию существующего пользователя
func (suite *TelegramAuthTestSuite) TestTelegramAuth_ExistingUser() {
	// Arrange - создаем пользователя заранее
	existingUser := &models.User{
		TelegramID: 54321,
		Username:   "existing_user",
		FirstName:  "Existing",
		LastName:   "User",
		Role:       "barber", // Уже барбер
		IsActive:   true,
	}
	err := suite.db.DB.Create(existingUser).Error
	suite.Require().NoError(err)

	// Данные для аутентификации существующего пользователя
	authData := models.TelegramAuthData{
		ID:        54321,
		Username:  "updated_username",
		FirstName: "Updated",
		LastName:  "Name",
		AuthDate:  time.Now().Unix(),
		Hash:      "test_hash",
	}

	jsonData, err := json.Marshal(authData)
	suite.Require().NoError(err)

	// Act
	req := httptest.NewRequest("POST", "/api/auth/telegram", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// Assert
	suite.Equal(http.StatusOK, w.Code)

	var authResponse models.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	suite.NoError(err)

	// Проверяем, что данные пользователя обновились
	suite.Equal("updated_username", authResponse.User.Username)
	suite.Equal("Updated", authResponse.User.FirstName)
	suite.Equal("Name", authResponse.User.LastName)
	suite.Equal("barber", authResponse.User.Role) // Роль сохранилась

	// Проверяем, что в БД данные обновились
	var user models.User
	err = suite.db.DB.Where("telegram_id = ?", 54321).First(&user).Error
	suite.NoError(err)
	suite.Equal("updated_username", user.Username)
	suite.Equal("Updated", user.FirstName)
	suite.Equal("Name", user.LastName)
}

// TestTelegramAuthTestSuite запускает все тесты
func TestTelegramAuthTestSuite(t *testing.T) {
	suite.Run(t, new(TelegramAuthTestSuite))
}

// TestAuthService упрощенная версия AuthService для тестов
type TestAuthService struct {
	userRepo  repositories.UserRepository
	rdb       *redis.Client
	jwtSecret string
	botToken  string
}

// NewTestAuthService создает тестовый сервис аутентификации
func NewTestAuthService(userRepo repositories.UserRepository, rdb *redis.Client, jwtSecret, botToken string) *TestAuthService {
	return &TestAuthService{
		userRepo:  userRepo,
		rdb:       rdb,
		jwtSecret: jwtSecret,
		botToken:  botToken,
	}
}

// ValidateTelegramAuth упрощенная валидация для тестов
func (s *TestAuthService) ValidateTelegramAuth(authData models.TelegramAuthData, botToken string) bool {
	// Для тестов проверяем, что ID не равен 0 и есть имя
	return authData.ID != 0 && authData.FirstName != ""
}

// AuthenticateUser находит или создает пользователя
func (s *TestAuthService) AuthenticateUser(authData models.TelegramAuthData) (*models.User, error) {
	// Ищем пользователя по TelegramID
	user, err := s.userRepo.GetByTelegramID(authData.ID)
	if err == nil {
		// Пользователь найден, обновляем данные
		user.Username = authData.Username
		user.FirstName = authData.FirstName
		user.LastName = authData.LastName
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// Пользователь не найден, создаем нового
	newUser := &models.User{
		TelegramID: authData.ID,
		Username:   authData.Username,
		FirstName:  authData.FirstName,
		LastName:   authData.LastName,
		Role:       "client", // По умолчанию клиент
		IsActive:   true,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// GenerateAccessToken генерирует access token
func (s *TestAuthService) GenerateAccessToken(user *models.User) (string, error) {
	claims := models.TokenClaims{
		UserID:     user.ID,
		TelegramID: user.TelegramID,
		Role:       user.Role,
		Type:       "access",
		Exp:        time.Now().Add(15 * time.Minute).Unix(),
		Iat:        time.Now().Unix(),
		Jti:        "test_jti",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// GenerateRefreshToken генерирует refresh token
func (s *TestAuthService) GenerateRefreshToken(user *models.User) (string, error) {
	claims := models.TokenClaims{
		UserID:     user.ID,
		TelegramID: user.TelegramID,
		Role:       user.Role,
		Type:       "refresh",
		Exp:        time.Now().Add(7 * 24 * time.Hour).Unix(),
		Iat:        time.Now().Unix(),
		Jti:        "test_jti",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ParseJWT парсит JWT токен
func (s *TestAuthService) ParseJWT(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// StoreRefreshToken сохраняет refresh token (для тестов не реализовано)
func (s *TestAuthService) StoreRefreshToken(userID uint, refreshToken string) error {
	return nil
}

// IsRefreshTokenValid проверяет refresh token (для тестов всегда true)
func (s *TestAuthService) IsRefreshTokenValid(userID uint, refreshToken string) bool {
	return true
}

// UpdateRefreshToken обновляет refresh token (для тестов не реализовано)
func (s *TestAuthService) UpdateRefreshToken(userID uint, oldToken, newToken string) error {
	return nil
}

// RevokeRefreshToken отзывает refresh token (для тестов не реализовано)
func (s *TestAuthService) RevokeRefreshToken(userID uint) error {
	return nil
}
