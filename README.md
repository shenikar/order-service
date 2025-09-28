# Order Service

Order Service — это микросервис на Go для приёма, обработки и хранения заказов.  
Сообщения с заказами поступают в **Kafka** от генератора заказов (`order_generator`) и сохраняются в **PostgreSQL**.

---

## Содержание

- [Технологии](#технологии)
- [Архитектура](#архитектура)
- [Запуск](#запуск)
- [Переменные окружения](#переменные-окружения)
- [Примеры использования](#примеры-использования)
- [Swagger документация](#swagger-документация)
- [Мониторинг](#мониторинг)
- [Тестирование](#тестирование)
- [Лицензия](#лицензия)

---

## Технологии

- Go 1.24+
- PostgreSQL 16
- Kafka (KRaft mode)
- Gin (HTTP сервер)
- sqlx (работа с PostgreSQL)
- segmentio/kafka-go (Kafka client)
- Golang-migrate (миграции базы данных)
- Testify (юнит-тестирование)
- Swagger (API документация)

---

## Архитектура

- **order_service** — основной сервис:
  - Консьюмит сообщения из Kafka.
  - Валидирует заказы.
  - Сохраняет их в PostgreSQL.
  - Предоставляет REST API для работы с заказами.

- **order_generator** — сервис для генерации и отправки тестовых заказов в Kafka.

- **PostgreSQL** — база данных для хранения заказов.

- **Kafka (KRaft mode)** — брокер сообщений для асинхронной обработки заказов.

---

## Запуск

### Требования

- Docker
- Docker Compose

### Шаги

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/shenikar/order-service.git
   cd order-service
    ```

2. Создайте .env на основе .env_example:

    ```bash
    cp .env_example .env
    ```

3. Поднимите сервисы:

    ```bash
    docker compose up --build
    ```

4. Убедитесь, что все контейнеры запущены:

    ```bash
    docker ps
    ```

---

## Переменные окружения

Все переменные окружения описаны в файле `.env_example`.

---

## Примеры использования

### Отправка тестового заказа

`order_generator` автоматически начинает слать тестовые заказы после запуска `order_service`.

### Проверка сохранённых заказов

Сервис хранит заказы в PostgreSQL. Для проверки можно подключиться к БД:

```bash
docker exec -it <postgres_container_id> psql -U postgres -d orders
```

И выполнить SQL-запрос:

```sql
SELECT * FROM orders;
```

---

## Swagger документация

REST API сервиса описан с помощью Swagger (OpenAPI).

После запуска сервиса документация доступна по адресу:
```bash
    http://localhost:8081/swagger/index.html
```
Для обновления Swagger-спецификации используется swaggo/swag
### Генерация документации
```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g cmd/order_service/main.go -o ./docs
```
После этого в папке ./docs будут сгенерированы файлы OpenAPI.

---

## Мониторинг

Проект включает в себя систему мониторинга на базе `Prometheus` и `Grafana`.

- **Prometheus** — сбор метрик.
- **Grafana** — визуализация метрик.

### Доступ к сервисам

- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Grafana**: [http://localhost:3000](http://localhost:3000) (логин: `admin`, пароль: `admin`)

### Настройка Grafana

1.  Откройте Grafana и войдите в систему.
2.  Перейдите в `Configuration` > `Data Sources`.
3.  Нажмите `Add data source` и выберите `Prometheus`.
4.  В поле `URL` укажите `http://prometheus:9090`.
5.  Нажмите `Save & Test`.

После этого вы сможете создавать дашборды, используя метрики из `Prometheus`.

### Основные метрики

- `http_requests_total` — общее количество HTTP-запросов.
  - `method` — HTTP-метод запроса (GET, POST и т.д.).
  - `path` — путь запроса (например, `/orders/:order_uid`).
  - `status` — HTTP-статус ответа (200, 404, 500 и т.д.).

---

## Тестирование

Для проекта написаны unit-тесты с использованием библиотеки Testify
### Запуск всех тестов
```bash
    go test -v ./...
```
### Запуск тестов с покрытием
```bash
    go test -cover ./...
```
### Запуск тестов конкретного пакета
```bash
    go test -v ./internal/service
```
Отчёт о покрытии можно сгенерировать в HTML:
```bash
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
```
---
## Лицензия

Этот проект распространяется под лицензией MIT.  
Подробнее см. файл [LICENSE](./LICENSE).
