package models

type Order struct {
	OrderUID          string   `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" db:"track_number" validate:"required"`
	Entry             string   `json:"entry" db:"entry" validate:"required"`
	Locale            string   `json:"locale" db:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature" db:"internal_signature"`
	CustomerID        string   `json:"customer_id" db:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" db:"delivery_service"`
	ShardKey          string   `json:"shardkey" db:"shardkey"`
	SmID              int      `json:"sm_id" db:"sm_id"`
	DateCreated       string   `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" db:"oof_shard"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required,dive,required"`
}

type Delivery struct {
	OrderUID string `json:"-" db:"order_uid"`
	Name     string `json:"name" db:"name" validate:"required"`
	Phone    string `json:"phone" db:"phone" validate:"required"`
	Zip      string `json:"zip" db:"zip" validate:"required"`
	City     string `json:"city" db:"city" validate:"required"`
	Address  string `json:"address" db:"address" validate:"required"`
	Region   string `json:"region" db:"region"`
	Email    string `json:"email" db:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string `json:"-" db:"order_uid"`
	Transaction  string `json:"transaction" db:"transaction" validate:"required"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency" validate:"required"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount" validate:"gte=0"`
	PaymentDT    int64  `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"gte=0"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"gte=0"`
}

type Item struct {
	OrderUID    string `json:"-" db:"order_uid"`
	ChrtID      int    `json:"chrt_id" db:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price       int    `json:"price" db:"price" validate:"gte=0"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name" validate:"required"`
	Sale        int    `json:"sale" db:"sale" validate:"gte=0"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price" validate:"gte=0"`
	NmID        int    `json:"nm_id" db:"nm_id" validate:"required"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}
