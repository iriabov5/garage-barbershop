# Garage Barbershop - Makefile для тестирования

.PHONY: test test-unit test-integration test-e2e test-all build run clean

# Запуск всех тестов
test-all: test-unit test-integration test-e2e

# Юнит тесты (быстрые, изолированные)
test-unit:
	@echo "🧪 Запуск юнит тестов..."
	go test -v ./tests/unit/... -short

# Интеграционные тесты (API, база данных)
test-integration:
	@echo "🔗 Запуск интеграционных тестов..."
	go test -v ./tests/integration/... -timeout 30s

# E2E тесты (полные сценарии)
test-e2e:
	@echo "🌐 Запуск E2E тестов..."
	go test -v ./tests/e2e/... -timeout 60s

# Все тесты
test:
	@echo "🚀 Запуск всех тестов..."
	go test -v ./tests/... -timeout 120s

# Сборка приложения
build:
	@echo "🔨 Сборка приложения..."
	go build -o main .

# Запуск приложения
run:
	@echo "🚀 Запуск приложения..."
	go run main.go

# Очистка
clean:
	@echo "🧹 Очистка..."
	rm -f main
	go clean

# Покрытие тестами
coverage:
	@echo "📊 Анализ покрытия тестами..."
	go test -v ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📄 Отчет сохранен в coverage.html"

# Линтинг
lint:
	@echo "🔍 Проверка кода..."
	golangci-lint run

# Форматирование
fmt:
	@echo "✨ Форматирование кода..."
	go fmt ./...

# Установка зависимостей
deps:
	@echo "📦 Установка зависимостей..."
	go mod tidy
	go mod download

# Помощь
help:
	@echo "📋 Доступные команды:"
	@echo "  test-unit        - Юнит тесты (быстрые)"
	@echo "  test-integration - Интеграционные тесты (API)"
	@echo "  test-e2e         - E2E тесты (полные сценарии)"
	@echo "  test-all         - Все тесты"
	@echo "  test             - Все тесты (краткая версия)"
	@echo "  build            - Сборка приложения"
	@echo "  run              - Запуск приложения"
	@echo "  clean            - Очистка"
	@echo "  coverage         - Анализ покрытия"
	@echo "  lint             - Линтинг кода"
	@echo "  fmt              - Форматирование"
	@echo "  deps             - Установка зависимостей"
	@echo "  help             - Эта справка"
