package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/models"
	"github.com/shenikar/order-service/internal/service"
)

var DLQWriter *kafka.Writer

// StartConsumer запускает Kafka consumer для обработки сообщений
func StartConsumer(ctx context.Context, cfg *config.Config, orderService *service.OrderService) *kafka.Reader {
	InitDLQWriter(cfg)
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Проверяем и создаём топик при необходимости
	if err := ensureTopic(cfg, dialer); err != nil {
		log.Fatalf("failed to ensure topic exists: %v", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          cfg.Kafka.Topic,
		GroupID:        cfg.Kafka.GroupID,
		StartOffset:    kafka.FirstOffset,
		Dialer:         dialer,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	go func() {
		for {
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Kafka consumer context canceled, stopping")
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					log.Println("Kafka consumer context deadline exceeded, stopping")
					return
				}
				log.Printf("Failed to read message: %v", err)
				continue
			}

			order := &models.Order{}
			if err := json.Unmarshal(msg.Value, order); err != nil {
				log.Printf("Invalid JSON, ignoring: %v", err)
				_ = reader.CommitMessages(ctx, msg)
				continue
			}

			// Валидация всех полей через validator
			if !orderService.ValidateOrder(order) {
				log.Printf("Invalid order data, ignoring: %+v", order)
				sendToDLQ(msg.Value)
				_ = reader.CommitMessages(ctx, msg)
				continue
			}

			// Сохранение в БД
			if err := orderService.SaveOrder(order); err != nil {
				log.Printf("Failed to save order: %v", err)
				// Не коммитим сообщение — оно будет прочитано снова, можно будет восстановить обработку
				continue
			}

			// Коммитим сообщение после успешной обработки
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit message %s: %v", order.OrderUID, err)
			} else {
				log.Printf("Order processed: %s", order.OrderUID)
			}

		}
	}()

	return reader
}

// StopConsumer корректно завершает работу Kafka consumer
func StopConsumer(reader *kafka.Reader, cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
	if reader != nil {
		log.Println("Shutting down Kafka consumer")
		if err := reader.Close(); err != nil {
			log.Printf("Failed to close Kafka reader: %v", err)
		}
	}
}

// ensureTopic проверяет, что топик существует, и создаёт его при необходимости
func ensureTopic(cfg *config.Config, dialer *kafka.Dialer) error {
	conn, err := dialer.Dial("tcp", cfg.Kafka.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	ctrlConn, err := dialer.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		return err
	}
	defer ctrlConn.Close()

	return ctrlConn.CreateTopics(kafka.TopicConfig{
		Topic:             cfg.Kafka.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

// InitDLQWriter инициализирует Kafka writer для DLQ
func InitDLQWriter(cfg *config.Config) {
	DLQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.DLQTopic,
	})
}

// sendToDLQ отправляет сообщение в DLQ
func sendToDLQ(data []byte) {
	if DLQWriter == nil {
		log.Println("DLQ writer not initialized")
		return
	}
	err := DLQWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: data,
		},
	)
	if err != nil {
		log.Printf("Failed to send message to DLQ: %v", err)
	}
}
