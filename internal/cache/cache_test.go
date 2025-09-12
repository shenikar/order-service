package cache

import (
	"testing"
	"time"

	"github.com/shenikar/order-service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	c, err := NewCache(10, time.Minute)
	assert.NoError(t, err)

	order := models.Order{OrderUID: "order1"}
	c.Set(order)

	got, found := c.Get("order1")
	assert.True(t, found)
	assert.Equal(t, order, got)
}

func TestCache_TTLExpiration(t *testing.T) {
	c, err := NewCache(10, time.Millisecond*100)
	assert.NoError(t, err)

	order := models.Order{OrderUID: "order_ttl"}
	c.Set(order)

	_, found := c.Get("order_ttl")
	assert.True(t, found)

	time.Sleep(time.Millisecond * 150)
	_, found = c.Get("order_ttl")
	assert.False(t, found)
}

func TestCache_LRUCapacity(t *testing.T) {
	c, err := NewCache(2, time.Minute)
	assert.NoError(t, err)

	order1 := models.Order{OrderUID: "order1"}
	order2 := models.Order{OrderUID: "order2"}
	order3 := models.Order{OrderUID: "order3"}

	c.Set(order1)
	c.Set(order2)
	c.Set(order3)

	_, found1 := c.Get("order1")
	_, found2 := c.Get("order2")
	_, found3 := c.Get("order3")

	assert.False(t, found1) // order1 должен быть вытеснен
	assert.True(t, found2)
	assert.True(t, found3)
}
