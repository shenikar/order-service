package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/models"
)

var cities = []string{"Kiryat Mozkin", "Tel Aviv", "Haifa", "Jerusalem"}
var names = []string{"Test Testov", "Alice Smith", "Bob Johnson"}
var emails = []string{"test@gmail.com", "alice@example.com", "bob@example.com"}
var brands = []string{"Vivienne Sabo", "Maybelline", "L'Oreal"}
var products = []string{"Mascaras", "Lipstick", "Foundation"}

func randomStringFromSlice(r *rand.Rand, slice []string) string {
	return slice[r.Intn(len(slice))]
}

func randomOrder(r *rand.Rand) *models.Order {
	now := time.Now()
	uid := fmt.Sprintf("%x", r.Int63())

	return &models.Order{
		OrderUID:        uid,
		TrackNumber:     fmt.Sprintf("WBILM%s", uid[:8]),
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      fmt.Sprintf("customer%d", r.Intn(1000)),
		DeliveryService: "meest",
		ShardKey:        fmt.Sprintf("%d", r.Intn(10)),
		SmID:            r.Intn(100),
		DateCreated:     now.Format(time.RFC3339),
		Delivery: models.Delivery{
			Name:    randomStringFromSlice(r, names),
			Phone:   fmt.Sprintf("+972%07d", r.Intn(10000000)),
			Zip:     fmt.Sprintf("%06d", r.Intn(1000000)),
			City:    randomStringFromSlice(r, cities),
			Address: fmt.Sprintf("Street %d", r.Intn(100)),
			Region:  "Kraiot",
			Email:   randomStringFromSlice(r, emails),
		},
		Payment: models.Payment{
			Transaction:  uid,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       r.Intn(2000),
			PaymentDT:    now.Unix(),
			Bank:         "alpha",
			DeliveryCost: r.Intn(500),
			GoodsTotal:   r.Intn(1500),
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      r.Intn(1000000),
				TrackNumber: fmt.Sprintf("WBILM%s", uid[:8]),
				Price:       r.Intn(1000),
				Name:        randomStringFromSlice(r, products),
				TotalPrice:  r.Intn(1000),
				Brand:       randomStringFromSlice(r, brands),
				Status:      202,
			},
		},
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Локальный генератор случайных чисел
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Настройка Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})
	defer func() {
		if err := writer.Close(); err != nil {
			log.Println("Failed to close Kafka writer:", err)
		}
	}()

	for i := 0; i < 5; i++ {
		order := randomOrder(r)
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
