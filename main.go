package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Middleware для логирования HTTP запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("🌐 %s %s %s - %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		
		next(w, r)
		
		duration := time.Since(start)
		log.Printf("⏱️  %s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

func main() {
	log.Println("🚀 Запуск Garage Barbershop сервера...")
	
	// Обработчик для главной страницы
	http.HandleFunc("/", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		log.Println("📄 Обслуживание главной страницы")
		
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
		log.Println("📊 Запрос статуса API")
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
		log.Println("💚 Health check запрос")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))

	// Получаем порт из переменной окружения или используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Логируем информацию о запуске
	log.Printf("🚀 Garage Barbershop сервер запускается на порту %s", port)
	log.Printf("📱 Откройте http://localhost:%s в браузере", port)
	log.Printf("🌍 Environment: %s", os.Getenv("ENVIRONMENT"))
	log.Printf("⏰ Время запуска: %s", time.Now().Format(time.RFC3339))
	
	// Запускаем сервер с логированием
	log.Println("✅ Сервер готов к работе!")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
