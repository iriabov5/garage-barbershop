package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TelegramAuthData представляет данные аутентификации от Telegram
type TelegramAuthData struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

// AuthResponse представляет ответ при аутентификации
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

// RefreshTokenRequest представляет запрос на обновление токена
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// DirectLoginRequest представляет запрос на прямую авторизацию
type DirectLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// DirectRegisterRequest представляет запрос на прямую регистрацию
type DirectRegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=client barber"`
}

// ClientRegisterRequest представляет запрос на регистрацию клиента
type ClientRegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// BarberRegisterRequest представляет запрос на регистрацию барбера (только админ)
type BarberRegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Specialties string `json:"specialties"` // специализации
	Experience  int    `json:"experience"`  // опыт в годах
}

// BarberUpdateRequest представляет запрос на обновление барбера (админ)
type BarberUpdateRequest struct {
	Email       string   `json:"email" binding:"omitempty,email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Specialties string   `json:"specialties"`
	Experience  int      `json:"experience"`
	IsActive    *bool    `json:"is_active"` // указатель для различения false и отсутствия поля
	Rating      *float64 `json:"rating"`    // указатель для различения 0 и отсутствия поля
}

// BarberSelfUpdateRequest представляет запрос на обновление собственного профиля барбера
type BarberSelfUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Specialties string `json:"specialties"`
	Experience  int    `json:"experience"`
}

// TokenClaims представляет claims JWT токена
type TokenClaims struct {
	UserID     uint   `json:"user_id"`
	TelegramID int64  `json:"telegram_id"`
	Role       string `json:"role"`
	Type       string `json:"type"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
	Jti        string `json:"jti"`
	jwt.RegisteredClaims
}

// IsExpired проверяет, истек ли токен
func (tc *TokenClaims) IsExpired() bool {
	return time.Now().Unix() > tc.Exp
}

// IsAccessToken проверяет, является ли токен access token
func (tc *TokenClaims) IsAccessToken() bool {
	return tc.Type == "access"
}

// IsRefreshToken проверяет, является ли токен refresh token
func (tc *TokenClaims) IsRefreshToken() bool {
	return tc.Type == "refresh"
}
