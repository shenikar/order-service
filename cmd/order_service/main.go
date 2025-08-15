package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
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

	//Корректное завершение работы приложения
	gracefulShutdown(orderService, dbConn)
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

// Graceful shutdown
func gracefulShutdown(orderService *service.OrderService, dbConn *sqlx.DB) {
	// канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Timeout для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Завершение Kafka consumer
	kafka.StopConsumer()
	log.Println("Kafka consumer stopped")

	// Завершение HTTP сервера
	if err := server.ShutdownServer(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	// Закрываем БД
	if err := dbConn.Close(); err != nil {
		log.Printf("Error closing DB connection: %v", err)
	}
	log.Println("Database connection closed")

	log.Println("Shutdown completed successfully")
}
