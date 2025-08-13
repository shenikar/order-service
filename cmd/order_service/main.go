package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/cache"
	"github.com/shenikar/order-service/internal/db"
	"github.com/shenikar/order-service/internal/kafka"
	"github.com/shenikar/order-service/internal/repository"
	"github.com/shenikar/order-service/internal/server"
	"github.com/shenikar/order-service/internal/service"
)

func main() {
	// Загружаем конфигурацию
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Выполняем миграции базы данных
	if err := runMigrations(config); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Подключаемся к БД
	dbConn, err := db.Connect(config)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	// Создаем компоненты приложения
	repo := repository.NewOrderRepository(dbConn)
	cache := cache.NewCache()
	orderService := service.NewOrderService(repo, cache)

	// Восстанавливаем кэш из БД
	if err := orderService.RestoreCacheFromDB(); err != nil {
		log.Println("Warning: Failed to restore cache from DB:", err)
	} else {
		log.Println("Cache restored from database successfully")
	}

	// Запускаем Kafka consumer в отдельной горутине
	go kafka.StartConsumer(config, orderService)

	// Запускаем HTTP сервер
	server.StartServer(config, orderService)
}

func runMigrations(config *config.Config) error {
	dbURL := config.GetDatabaseUrl()

	m, err := migrate.New(
		"file://migrations", dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations applied successfully")
	return nil
}
