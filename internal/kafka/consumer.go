package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/models"
	"github.com/shenikar/order-service/internal/service"
)

// StartConsumer запускает Kafka consumer для обработки сообщений
func StartConsumer(config *config.Config, orderService *service.OrderService) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: config.Kafka.GroupID,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Failed to close Kafka reader: %v", err)
		}
	}()

	log.Printf("Kafka consumer started. Brokers: %v, Topic: %s, GroupID: %s",
		config.Kafka.Brokers, config.Kafka.Topic, config.Kafka.GroupID)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
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
