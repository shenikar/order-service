package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shenikar/order-service/internal/mapper"
	"github.com/shenikar/order-service/internal/models"
)

type OrderRepositoryInterface interface {
	SaveOrder(order *models.Order) error
	GetOrderByUID(orderUID string) (*models.Order, error)
	GetItemByOrderUID(orderUID string) ([]models.Item, error)
	GetAllOrders() ([]models.Order, error)
}

type OrderRepository struct {
	db *sqlx.DB
}

// NewOrderRepository создает новый экземпляр OrderRepository
func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// SaveOrder сохраняет заказ в базе данных
func (r *OrderRepository) SaveOrder(order *models.Order) error {
	// Используем транзакции
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сохраняем заказ
	_, err = tx.NamedExec(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)
        ON CONFLICT (order_uid) DO NOTHING`, order)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	// Сохраняем доставку
	order.Delivery.OrderUID = order.OrderUID
	_, err = tx.NamedExec(`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
        VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)
        ON CONFLICT (order_uid) DO NOTHING`, order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	// Сохраняем платеж
	order.Payment.OrderUID = order.OrderUID
	_, err = tx.NamedExec(`INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES (:order_uid, :transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
        ON CONFLICT (order_uid) DO NOTHING`, order.Payment)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		item.OrderUID = order.OrderUID
		_, err = tx.NamedExec(`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES (:order_uid, :chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)
            ON CONFLICT (chrt_id) DO NOTHING`, item)
		if err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
	}
	// Фиксируем транзакцию
	return tx.Commit()
}

// GetAllOrders возвращает все заказы из базы данных
func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	query := `SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
               o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank,
               p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        JOIN deliveries d ON o.order_uid = d.order_uid
        JOIN payments p ON o.order_uid = p.order_uid`

	var dbOrders []models.OrderDB
	if err := r.db.Select(&dbOrders, query); err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return mapper.MapOrdersDBToModels(dbOrders), nil
}

// GetItemByOrderUID возвращает товары по для конкретного заказа
func (r *OrderRepository) GetItemByOrderUID(orderUID string) ([]models.Item, error) {
	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1`

	var items []models.Item
	err := r.db.Select(&items, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items for order %s: %w", orderUID, err)
	}

	return items, nil
}

// GetOrderByUID возвращает заказ по его уникальному идентификатору
func (r *OrderRepository) GetOrderByUID(orderUID string) (*models.Order, error) {
	query := `SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
               o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank,
               p.delivery_cost, p.goods_total, p.custom_fee
        	FROM orders o
        	JOIN deliveries d ON o.order_uid = d.order_uid
        	JOIN payments p ON o.order_uid = p.order_uid
        	WHERE o.order_uid = $1`

	var dbo models.OrderDB
	if err := r.db.Get(&dbo, query, orderUID); err != nil {
		return nil, fmt.Errorf("failed to get order by UID %s: %w", orderUID, err)
	}

	order := mapper.MapOrderDBToModel(dbo)

	items, err := r.GetItemByOrderUID(orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}
