# Dockerfile для Railway с автоматическими тестами
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты для тестов
RUN apk add --no-cache \
    git \
    make \
    postgresql-client \
    redis

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Запускаем тесты (если тесты не пройдут, сборка упадет)
RUN echo "🧪 Запуск тестов в Railway..." && \
    make test-all && \
    echo "✅ Все тесты пройдены!"

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
