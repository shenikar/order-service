package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/shenikar/order-service/internal/cache"
	"github.com/shenikar/order-service/internal/models"
	"github.com/shenikar/order-service/internal/repository"
)

var validate = validator.New()

type OrderService struct {
	repo  repository.OrderRepositoryInterface
	cache *cache.Cache // Добавляем кэш для оптимизации
}

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(repo repository.OrderRepositoryInterface, cache *cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

// SaveOrder сохраняет заказ
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

	return nil
}

// GetOrderByUID извлекает заказ из кэша или БД
func (s *OrderService) GetOrderByUID(orderUID string) (*models.Order, error) {
	// проверяем кэш
	if order, found := s.cache.Get(orderUID); found {
		log.Printf("Order %s retrieved from cache", orderUID)
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
	log.Printf("Order %s retrieved from database", orderUID)

	return order, nil
}

// RestoreCacheFromDB восстанавливает кэш из БД
func (s *OrderService) RestoreCacheFromDB() error {
	orders, err := s.repo.GetAllOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		items, err := s.repo.GetItemByOrderUID(order.OrderUID)
		if err != nil {
			continue // Пропускаем заказы с ошибками
		}
		order.Items = items
		s.cache.Set(order)
	}

	return nil
}

func (s *OrderService) ValidateOrder(order *models.Order) bool {
	err := validate.Struct(order)
	if err != nil {
		log.Printf("Invalid order %s: %v", order.OrderUID, err)
		return false
	}
	return true
}
