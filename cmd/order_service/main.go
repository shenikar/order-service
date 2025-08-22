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
	kf "github.com/segmentio/kafka-go"
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

	// Создаем context для Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka consumer
	reader := kafka.StartConsumer(ctx, config, orderService)

	// Запускаем HTTP сервер
	server.StartServer(config, orderService)

	// Корректное завершение работы приложения
	gracefulShutdown(orderService, dbConn, reader, cancel)
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
func gracefulShutdown(orderService *service.OrderService, dbConn *sqlx.DB, reader *kf.Reader, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Отменяем context для Kafka
	cancel()

	// Timeout для завершения HTTP сервера
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Завершаем Kafka consumer
	kafka.StopConsumer(reader, nil)
	log.Println("Kafka consumer stopped")

	// Завершаем HTTP сервер
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
