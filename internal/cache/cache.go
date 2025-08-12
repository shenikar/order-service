package cache

import (
	"sync"

	"github.com/shenikar/order-service/internal/models"
)

type Cache struct {
	data map[string]models.Order
	mu   sync.RWMutex
}

// NewCache создает новый экземпляр кэша
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]models.Order),
	}
}

// Set добавляет или обновляет заказ в кэше
func (c *Cache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

// Get извлекает заказ из кэша по OrderUID
func (c *Cache) Get(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.data[orderUID]
	return order, ok
}
