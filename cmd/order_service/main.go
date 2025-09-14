package main

import (
	"context"
	"errors"
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

// @title Order Service API
// @version 1.0
// @description API для управления заказами
// @host localhost:8080
// @BasePath /
func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Инициализируем Kafka DLQ writer
	kafka.InitDLQWriter(cfg)

	// Выполняем миграции базы данных
	if err := runMigrations(cfg); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Подключаемся к БД
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Создаем компоненты приложения
	repo := repository.NewOrderRepository(dbConn)

	cacheOrder, err := cache.NewCache(cfg.Cache.Capacity, time.Duration(cfg.Cache.TTL)*time.Minute)
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
	}

	orderService := service.NewOrderService(repo, cacheOrder)

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
	consumer := kafka.StartConsumer(ctx, cfg, orderService)

	// Запускаем HTTP сервер
	server.StartServer(cfg, orderService)

	// Корректное завершение работы приложения
	gracefulShutdown(dbConn, consumer, cancel)
}

func runMigrations(cfg *config.Config) error {
	dbURL := cfg.GetDatabaseURL()

	m, err := migrate.New(
		"file://migrations", dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Database migrations applied successfully")
	return nil
}

// Graceful shutdown
func gracefulShutdown(dbConn *sqlx.DB, consumer *kf.Reader, cancel context.CancelFunc) {
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
	kafka.StopConsumer(consumer, nil)
	log.Println("Kafka consumer stopped")

	// Завершаем DLQ writer
	if kafka.DLQWriter != nil {
		if err := kafka.DLQWriter.Close(); err != nil {
			log.Printf("Error closing DLQ writer: %v", err)
		}
		log.Println("DLQ writer closed")
	}

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
