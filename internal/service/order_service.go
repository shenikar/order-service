package service

import (
	"github.com/shenikar/order-service/internal/cache"
	"github.com/shenikar/order-service/internal/models"
	"github.com/shenikar/order-service/internal/repository"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache *cache.Cache // Добавляем кэш для оптимизации
}

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(repo *repository.OrderRepository, cache *cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

// SaveOrder сохраняет заказ и обновляет кэш
func (s *OrderService) SaveOrder(order *models.Order) error {
	// Сохраняем заказ в БД
	if err := s.repo.SaveOrder(order); err != nil {
		return err
	}

	// Получаем items для кэширования
	items, err := s.repo.GetItemByOrderUID(order.OrderUID)
	if err != nil {
		return err
	}
	order.Items = items

	// Обновляем кэш
	s.cache.Set(*order)

	return nil
}

// GetOrderByUID извлекает заказ из кэша или БД
func (s *OrderService) GetOrderByUID(orderUID string) (*models.Order, error) {
	// проверяем кэш
	if order, found := s.cache.Get(orderUID); found {
		return &order, nil
	}

	// если нет в кэше, извлекаем из БД
	order, err := s.repo.GetOrderByUID(orderUID)
	if err != nil {
		return nil, err
	}

	// Загружаем items для заказа
	items, err := s.repo.GetItemByOrderUID(orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	// Обновляем кэш
	s.cache.Set(*order)

	return order, nil
}
