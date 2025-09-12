package service

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/shenikar/order-service/internal/cache"
	"github.com/shenikar/order-service/internal/models"
	"github.com/stretchr/testify/assert"
)

// mockRepo реализует интерфейс OrderRepository
type mockRepo struct {
	saveOrder func(order *models.Order) error
	getByUID  func(uid string) (*models.Order, error)
	getItems  func(uid string) ([]models.Item, error)
	getAll    func() ([]models.Order, error)
}

func (m *mockRepo) SaveOrder(order *models.Order) error {
	if m.saveOrder != nil {
		return m.saveOrder(order)
	}
	return nil
}
func (m *mockRepo) GetOrderByUID(uid string) (*models.Order, error) {
	if m.getByUID != nil {
		return m.getByUID(uid)
	}
	return nil, nil
}
func (m *mockRepo) GetItemByOrderUID(uid string) ([]models.Item, error) {
	if m.getItems != nil {
		return m.getItems(uid)
	}
	return nil, nil
}
func (m *mockRepo) GetAllOrders() ([]models.Order, error) {
	if m.getAll != nil {
		return m.getAll()
	}
	return nil, nil
}

func TestSaveOrder_Success(t *testing.T) {
	repo := &mockRepo{
		saveOrder: func(order *models.Order) error { return nil },
		getItems: func(uid string) ([]models.Item, error) {
			return []models.Item{{ChrtID: 1, TrackNumber: "TN123"}}, nil
		},
	}
	c, err := cache.NewCache(100, time.Minute*5)
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
	}
	svc := NewOrderService(repo, c)

	order := &models.Order{OrderUID: "uid123"}
	err = svc.SaveOrder(order)

	assert.NoError(t, err)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, "TN123", order.Items[0].TrackNumber)
}

func TestSaveOrder_RepoError(t *testing.T) {
	repo := &mockRepo{
		saveOrder: func(order *models.Order) error {
			return errors.New("db error")
		},
	}
	c, err := cache.NewCache(100, time.Minute*5)
	assert.NoError(t, err)
	svc := NewOrderService(repo, c)

	order := &models.Order{OrderUID: "uid123"}
	err = svc.SaveOrder(order)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestGetOrderByUID_FromCache(t *testing.T) {
	repo := &mockRepo{}
	c, err := cache.NewCache(100, time.Minute*5)
	assert.NoError(t, err)
	svc := NewOrderService(repo, c)

	// кладём заказ в кэш
	expected := models.Order{OrderUID: "uid123"}
	c.Set(expected)

	order, err := svc.GetOrderByUID("uid123")

	assert.NoError(t, err)
	assert.Equal(t, "uid123", order.OrderUID)
}

func TestGetOrderByUID_FromDB(t *testing.T) {
	repo := &mockRepo{
		getByUID: func(uid string) (*models.Order, error) {
			return &models.Order{OrderUID: uid}, nil
		},
		getItems: func(uid string) ([]models.Item, error) {
			return []models.Item{{ChrtID: 2, TrackNumber: "TN999"}}, nil
		},
	}
	c, err := cache.NewCache(100, time.Minute*5)
	assert.NoError(t, err)
	svc := NewOrderService(repo, c)

	order, err := svc.GetOrderByUID("uid456")

	assert.NoError(t, err)
	assert.Equal(t, "uid456", order.OrderUID)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, "TN999", order.Items[0].TrackNumber)

	// проверяем, что в кэш положилось
	cached, found := c.Get("uid456")
	assert.True(t, found)
	assert.Equal(t, "uid456", cached.OrderUID)
}

func TestGetOrderByUID_DBError(t *testing.T) {
	repo := &mockRepo{
		getByUID: func(uid string) (*models.Order, error) {
			return nil, errors.New("not found")
		},
	}
	c, err := cache.NewCache(100, time.Minute*5)
	assert.NoError(t, err)
	svc := NewOrderService(repo, c)

	order, err := svc.GetOrderByUID("bad_uid")

	assert.Nil(t, order)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}

func TestRestoreCacheFromDB(t *testing.T) {
	repo := &mockRepo{
		getAll: func() ([]models.Order, error) {
			return []models.Order{
				{OrderUID: "uid111"},
				{OrderUID: "uid222"},
			}, nil
		},
		getItems: func(uid string) ([]models.Item, error) {
			return []models.Item{{ChrtID: 1, TrackNumber: "TN_" + uid}}, nil
		},
	}
	c, err := cache.NewCache(100, time.Minute*5)
	assert.NoError(t, err)
	svc := NewOrderService(repo, c)

	err = svc.RestoreCacheFromDB()

	assert.NoError(t, err)

	// проверяем, что оба заказа оказались в кэше
	o1, found1 := c.Get("uid111")
	o2, found2 := c.Get("uid222")

	assert.True(t, found1)
	assert.True(t, found2)
	assert.Equal(t, "uid111", o1.OrderUID)
	assert.Equal(t, "uid222", o2.OrderUID)
}
