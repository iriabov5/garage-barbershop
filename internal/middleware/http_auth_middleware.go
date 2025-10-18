package middleware

import (
	"context"
	"net/http"
	"strings"

	"garage-barbershop/internal/services"
)

// HTTPAuthMiddleware проверяет JWT токен и добавляет данные пользователя в контекст запроса
func HTTPAuthMiddleware(authService services.AuthService) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Требуется токен аутентификации", http.StatusUnauthorized)
				return
			}

			// Токен должен быть в формате "Bearer <token>"
			if !strings.HasPrefix(tokenString, "Bearer ") {
				http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
				return
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := authService.ParseJWT(tokenString)
			if err != nil {
				http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Проверяем, что это access token
			if !claims.IsAccessToken() {
				http.Error(w, "Неверный тип токена: требуется access token", http.StatusUnauthorized)
				return
			}

			// Добавляем данные пользователя в контекст запроса
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "telegramID", claims.TelegramID)
			ctx = context.WithValue(ctx, "userRoles", claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// HTTPRequireRoleMiddleware проверяет роль пользователя
func HTTPRequireRoleMiddleware(requiredRole string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("userRoles").(string)
			if !ok || userRoles != requiredRole {
				http.Error(w, "Недостаточно прав", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

// HTTPRequireAnyRoleMiddleware проверяет любую из ролей
func HTTPRequireAnyRoleMiddleware(roles ...string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("userRoles").(string)
			if !ok {
				http.Error(w, "Роль пользователя не найдена", http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, role := range roles {
				if userRoles == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Недостаточно прав", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
