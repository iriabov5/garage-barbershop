package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Обработчик для главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
<!DOCTYPE html>
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
</html>
		`)
	})

	// Обработчик для API статуса
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "ok",
			"service": "Garage Barbershop",
			"version": "1.0.0",
			"message": "Сервер работает корректно"
		}`)
	})

	// Обработчик для health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Получаем порт из переменной окружения или используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	fmt.Printf("🚀 Garage Barbershop сервер запускается на порту %s\n", port)
	fmt.Printf("📱 Откройте http://localhost:%s в браузере\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
