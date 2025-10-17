#!/bin/bash

# Скрипт для разработки с автоматическими тестами

set -e

echo "🚀 Запуск Garage Barbershop в режиме разработки..."

# Функция для очистки
cleanup() {
    echo "🧹 Остановка контейнеров..."
    docker-compose down
}

# Обработчик сигналов для корректного завершения
trap cleanup EXIT INT TERM

# Проверяем, что Docker запущен
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker не запущен. Пожалуйста, запустите Docker Desktop."
    exit 1
fi

# Собираем и запускаем контейнеры
echo "🔨 Сборка контейнеров..."
docker-compose build

echo "🚀 Запуск сервисов..."
docker-compose up

# Если нужно запустить только тесты
if [ "$1" = "test" ]; then
    echo "🧪 Запуск только тестов..."
    docker-compose --profile test-only up tests
fi

# Если нужно запустить без тестов
if [ "$1" = "no-test" ]; then
    echo "🚀 Запуск без тестов..."
    docker-compose up postgres redis app
fi
