package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/models"
)

func randomOrder() *models.Order {
	if err := gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Fatalf("failed to seed gofakeit: %v", err)
	}
	now := time.Now()
	return &models.Order{
		OrderUID:        gofakeit.UUID(),
		TrackNumber:     gofakeit.Regex("WBILM[0-9A-Z]{8}"),
		Entry:           "WBIL",
		Locale:          gofakeit.Language(),
		CustomerID:      gofakeit.Username(),
		DeliveryService: gofakeit.Company(),
		ShardKey:        gofakeit.DigitN(1),
		SmID:            gofakeit.Number(0, 100),
		DateCreated:     now.Format(time.RFC3339),
		Delivery: models.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: models.Payment{
			Transaction:  gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       int(gofakeit.Price(100, 2000)),
			PaymentDT:    now.Unix(),
			Bank:         gofakeit.Company(),
			DeliveryCost: int(gofakeit.Price(0, 500)),
			GoodsTotal:   int(gofakeit.Price(0, 1500)),
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      gofakeit.Number(1000, 999999),
				TrackNumber: gofakeit.Regex("WBILM[0-9A-Z]{8}"),
				Price:       int(gofakeit.Price(100, 1000)),
				Name:        gofakeit.ProductName(),
				TotalPrice:  int(gofakeit.Price(100, 1000)),
				Brand:       gofakeit.Company(),
				Status:      202,
				NmID:        gofakeit.Number(1, 999999),
			},
		},
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Настройка Kafka writer
	writer := kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Println("Failed to close Kafka writer:", err)
		}
	}()

	for i := 0; i < 5; i++ {
		order := randomOrder()
		data, err := json.Marshal(order)
		if err != nil {
			log.Println("Failed to marshal order:", err)
			continue
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(order.OrderUID),
				Value: data,
			},
		)
		if err != nil {
			log.Println("Failed to send order:", err)
		} else {
			log.Println("Order sent:", order.OrderUID)
		}

		time.Sleep(1 * time.Second)
	}
}
