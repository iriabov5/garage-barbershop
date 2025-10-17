package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Глобальные переменные для подключений
var (
	db  *gorm.DB
	rdb *redis.Client
)

// Подключение к PostgreSQL
func connectDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("⚠️  DATABASE_URL не установлен, пропускаем подключение к БД")
		return nil
	}

	var err error
	db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")
	return nil
}

// Подключение к Redis
func connectRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("⚠️  REDIS_URL не установлен, пропускаем подключение к Redis")
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга Redis URL: %v", err)
	}

	rdb = redis.NewClient(opt)
	
	// Проверяем подключение
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	log.Println("✅ Подключение к Redis успешно")
	return nil
}

// Middleware для логирования HTTP запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем только в development режиме
		if os.Getenv("ENVIRONMENT") != "production" {
			log.Printf("🌐 %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		}

		next(w, r)

		// Логируем только медленные запросы в production
		duration := time.Since(start)
		if os.Getenv("ENVIRONMENT") == "production" && duration > 100*time.Millisecond {
			log.Printf("SLOW: %s %s took %v", r.Method, r.URL.Path, duration)
		} else if os.Getenv("ENVIRONMENT") != "production" {
			log.Printf("⏱️  %s %s completed in %v", r.Method, r.URL.Path, duration)
		}
	}
}

func main() {
	// Логируем запуск только в development
	if os.Getenv("ENVIRONMENT") != "production" {
		log.Println("🚀 Запуск Garage Barbershop сервера...")
	}

	// Подключаемся к базам данных
	if err := connectDB(); err != nil {
		log.Printf("❌ Ошибка подключения к PostgreSQL: %v", err)
	}

	if err := connectRedis(); err != nil {
		log.Printf("❌ Ошибка подключения к Redis: %v", err)
	}

	// Обработчик для главной страницы
	http.HandleFunc("/", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Логируем только в development
		if os.Getenv("ENVIRONMENT") != "production" {
			log.Println("📄 Обслуживание главной страницы")
		}

		html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Garage Barbershop</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            text-align: center;
            max-width: 500px;
        }
        h1 {
            color: #333;
            margin-bottom: 1rem;
        }
        .status {
            color: #28a745;
            font-weight: bold;
            margin: 1rem 0;
        }
        .info {
            color: #666;
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🪒 Garage Barbershop</h1>
        <div class="status">✅ Сервер работает!</div>
        <p>Добро пожаловать в систему управления барбершопом</p>
        <div class="info">
            <p>Версия: 1.0.0</p>
            <p>Статус: Готов к разработке</p>
        </div>
    </div>
</body>
</html>`

		fmt.Fprint(w, html)
	}))

	// Обработчик для API статуса
	http.HandleFunc("/api/status", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Логируем только в development
		if os.Getenv("ENVIRONMENT") != "production" {
			log.Println("📊 Запрос статуса API")
		}
		fmt.Fprintf(w, `{
			"status": "ok",
			"service": "Garage Barbershop",
			"version": "1.0.0",
			"message": "Сервер работает корректно",
			"timestamp": "%s"
		}`, time.Now().Format(time.RFC3339))
	}))

	// Обработчик для health check
	http.HandleFunc("/health", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Health check не логируем - он вызывается часто
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))

	// Обработчик для проверки статуса баз данных
	http.HandleFunc("/api/db-status", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		status := map[string]interface{}{
			"postgresql": "disconnected",
			"redis":      "disconnected",
		}

		// Проверяем PostgreSQL
		if db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					status["postgresql"] = "connected"
				}
			}
		}

		// Проверяем Redis
		if rdb != nil {
			ctx := context.Background()
			if _, err := rdb.Ping(ctx).Result(); err == nil {
				status["redis"] = "connected"
			}
		}

		fmt.Fprintf(w, `{
			"databases": %+v,
			"timestamp": "%s"
		}`, status, time.Now().Format(time.RFC3339))
	}))

	// Получаем порт из переменной окружения или используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Логируем информацию о запуске только в development
	if os.Getenv("ENVIRONMENT") != "production" {
		log.Printf("🚀 Garage Barbershop сервер запускается на порту %s", port)
		log.Printf("📱 Откройте http://localhost:%s в браузере", port)
		log.Printf("🌍 Environment: %s", os.Getenv("ENVIRONMENT"))
		log.Printf("⏰ Время запуска: %s", time.Now().Format(time.RFC3339))
		log.Println("✅ Сервер готов к работе!")
	} else {
		// В production только минимальная информация
		log.Printf("Server starting on port %s", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
