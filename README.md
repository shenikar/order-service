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
