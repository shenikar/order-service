package cache

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/shenikar/order-service/internal/models"
)

type Cache struct {
	lru *expirable.LRU[string, models.Order]
}

// NewCache создает новый экземпляр кэша
func NewCache(capacity int, ttl time.Duration) (*Cache, error) {
	lru := expirable.NewLRU[string, models.Order](capacity, nil, ttl)
	return &Cache{lru: lru}, nil
}

// Set добавляет или обновляет заказ в кэше
func (c *Cache) Set(order models.Order) {
	c.lru.Add(order.OrderUID, order)
}

// Get извлекает заказ из кэша по OrderUID
func (c *Cache) Get(orderUID string) (models.Order, bool) {
	return c.lru.Get(orderUID)
}
