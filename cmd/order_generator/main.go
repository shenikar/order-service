package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shenikar/order-service/internal/models"
)

var cities = []string{"Kiryat Mozkin", "Tel Aviv", "Haifa", "Jerusalem"}
var names = []string{"Test Testov", "Alice Smith", "Bob Johnson"}
var emails = []string{"test@gmail.com", "alice@example.com", "bob@example.com"}
var brands = []string{"Vivienne Sabo", "Maybelline", "L'Oreal"}
var products = []string{"Mascaras", "Lipstick", "Foundation"}

func randomStringFromSlice(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func randomOrder() *models.Order {
	now := time.Now()
	uid := fmt.Sprintf("%x", rand.Int63())

	return &models.Order{
		OrderUID:        uid,
		TrackNumber:     fmt.Sprintf("WBILM%s", uid[:8]),
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      fmt.Sprintf("customer%d", rand.Intn(1000)),
		DeliveryService: "meest",
		ShardKey:        fmt.Sprintf("%d", rand.Intn(10)),
		SmID:            rand.Intn(100),
		DateCreated:     now.Format(time.RFC3339),
		Delivery: models.Delivery{
			Name:    randomStringFromSlice(names),
			Phone:   fmt.Sprintf("+972%07d", rand.Intn(10000000)),
			Zip:     fmt.Sprintf("%06d", rand.Intn(1000000)),
			City:    randomStringFromSlice(cities),
			Address: fmt.Sprintf("Street %d", rand.Intn(100)),
			Region:  "Kraiot",
			Email:   randomStringFromSlice(emails),
		},
		Payment: models.Payment{
			Transaction:  uid,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(2000),
			PaymentDT:    now.Unix(),
			Bank:         "alpha",
			DeliveryCost: rand.Intn(500),
			GoodsTotal:   rand.Intn(1500),
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      rand.Intn(1000000),
				TrackNumber: fmt.Sprintf("WBILM%s", uid[:8]),
				Price:       rand.Intn(1000),
				Name:        randomStringFromSlice(products),
				TotalPrice:  rand.Intn(1000),
				Brand:       randomStringFromSlice(brands),
				Status:      202,
			},
		},
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})
	defer writer.Close()

	for i := 0; i < 5; i++ {
		order := randomOrder()
		data, _ := json.Marshal(order)
		err := writer.WriteMessages(context.TODO(),
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
