package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/models"
	"github.com/shenikar/order-service/internal/service"
)

var (
	reader     *kafka.Reader
	readerLock sync.Mutex
)

// StartConsumer запускает Kafka consumer для обработки сообщений
func StartConsumer(config *config.Config, orderService *service.OrderService) {
	readerLock.Lock()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.GroupID,
	})
	readerLock.Unlock()

	log.Printf("Kafka consumer started. Brokers: %v, Topic: %s, GroupID: %s",
		config.Kafka.Brokers, config.Kafka.Topic, config.Kafka.GroupID)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err == context.Canceled {
				break
			}
			log.Printf("Failed to read message: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Invalid JSON format: %v", err)
			continue
		}

		if err := orderService.SaveOrder(&order); err != nil {
			log.Printf("Failed to save order: %v", err)
			continue
		}

		log.Printf("Order processed successfully: %s", order.OrderUID)
	}
}

// StopConsumer корректно завершает работу Kafka consumer
func StopConsumer() {
	readerLock.Lock()
	defer readerLock.Unlock()

	if reader != nil {
		log.Println("Shutting down Kafka consumer")
		if err := reader.Close(); err != nil {
			log.Printf("Failed to close Kafka reader: %v", err)
		}
		reader = nil
	}
}
