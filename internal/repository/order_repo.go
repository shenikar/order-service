package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shenikar/order-service/internal/models"
)

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
	order.Delivery.OrderUID = order.OrederUID
	_, err = tx.NamedExec(`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
        VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)
        ON CONFLICT (order_uid) DO NOTHING`, order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	// Сохраняем платеж
	order.Payment.OrderUID = order.OrederUID
	_, err = tx.NamedExec(`INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES (:order_uid, :transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
        ON CONFLICT (order_uid) DO NOTHING`, order.Payment)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		item.OrderUID = order.OrederUID
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
