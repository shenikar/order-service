package models

// вспомогательная структура для чтения из БД
type OrderDB struct {
	// Orders
	OrderUID          string `db:"order_uid"`
	TrackNumber       string `db:"track_number"`
	Entry             string `db:"entry"`
	Locale            string `db:"locale"`
	InternalSignature string `db:"internal_signature"`
	CustomerID        string `db:"customer_id"`
	DeliveryService   string `db:"delivery_service"`
	ShardKey          string `db:"shardkey"`
	SmID              int    `db:"sm_id"`
	DateCreated       string `db:"date_created"`
	OofShard          string `db:"oof_shard"`

	// Delivery
	DeliveryName    string `db:"name"`
	DeliveryPhone   string `db:"phone"`
	DeliveryZip     string `db:"zip"`
	DeliveryCity    string `db:"city"`
	DeliveryAddress string `db:"address"`
	DeliveryRegion  string `db:"region"`
	DeliveryEmail   string `db:"email"`

	// Payment
	PaymentTransaction  string `db:"transaction"`
	PaymentRequestID    string `db:"request_id"`
	PaymentCurrency     string `db:"currency"`
	PaymentProvider     string `db:"provider"`
	PaymentAmount       int    `db:"amount"`
	PaymentDT           int64  `db:"payment_dt"`
	PaymentBank         string `db:"bank"`
	PaymentDeliveryCost int    `db:"delivery_cost"`
	PaymentGoodsTotal   int    `db:"goods_total"`
	PaymentCustomFee    int    `db:"custom_fee"`
}
