# Сборка
FROM golang:1.24-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git curl

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum и качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o order_service ./cmd/order_service/main.go

# Сборка финального образа
FROM alpine:latest

# Установка зависимостей для работы приложения
RUN apk add --no-cache ca-certificates \
    curl

WORKDIR /app

# Миграции и .env
COPY migrations ./migrations
COPY .env .env
COPY web /app/web


# Копирование собранного приложения из builder
COPY --from=builder /app/order_service .

# Порт
EXPOSE 8081

# Запуск приложения
CMD ["./order_service"]