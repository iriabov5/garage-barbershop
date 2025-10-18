package middleware

import (
	"net/http"
	"strings"

	"garage-barbershop/internal/services"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware middleware для проверки JWT токенов
func JWTMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		token := parts[1]

		// Парсим JWT токен
		claims, err := authService.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			c.Abort()
			return
		}

		// Проверяем, что это access token
		if !claims.IsAccessToken() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный тип токена"})
			c.Abort()
			return
		}

		// Проверяем срок действия
		if claims.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен истек"})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контекст
		c.Set("user_id", claims.UserID)
		c.Set("telegram_id", claims.TelegramID)
		c.Set("user_roles", claims.Roles)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RequireRole middleware для проверки роли пользователя
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Роль пользователя не найдена"})
			c.Abort()
			return
		}

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware для проверки любой из ролей
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Роль пользователя не найдена"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав доступа"})
		c.Abort()
	}
}

// OptionalAuth middleware для опциональной аутентификации
func OptionalAuth(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := authService.ParseJWT(token)
		if err != nil {
			c.Next()
			return
		}

		if claims.IsAccessToken() && !claims.IsExpired() {
			c.Set("user_id", claims.UserID)
			c.Set("telegram_id", claims.TelegramID)
			c.Set("user_roles", claims.Roles)
			c.Set("jwt_claims", claims)
		}

		c.Next()
	}
}
