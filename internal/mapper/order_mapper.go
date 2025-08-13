package mapper

import (
	"github.com/shenikar/order-service/internal/models"
)

// MapOrderDBToModel преобразует OrderDB в Order
func MapOrderDBToModel(dbo models.OrderDB) models.Order {
	return models.Order{
		OrderUID:          dbo.OrderUID,
		TrackNumber:       dbo.TrackNumber,
		Entry:             dbo.Entry,
		Locale:            dbo.Locale,
		InternalSignature: dbo.InternalSignature,
		CustomerID:        dbo.CustomerID,
		DeliveryService:   dbo.DeliveryService,
		ShardKey:          dbo.ShardKey,
		SmID:              dbo.SmID,
		DateCreated:       dbo.DateCreated,
		OofShard:          dbo.OofShard,
		Delivery: models.Delivery{
			Name:    dbo.DeliveryName,
			Phone:   dbo.DeliveryPhone,
			Zip:     dbo.DeliveryZip,
			City:    dbo.DeliveryCity,
			Address: dbo.DeliveryAddress,
			Region:  dbo.DeliveryRegion,
			Email:   dbo.DeliveryEmail,
		},
		Payment: models.Payment{
			Transaction:  dbo.PaymentTransaction,
			RequestID:    dbo.PaymentRequestID,
			Currency:     dbo.PaymentCurrency,
			Provider:     dbo.PaymentProvider,
			Amount:       dbo.PaymentAmount,
			PaymentDT:    dbo.PaymentDT,
			Bank:         dbo.PaymentBank,
			DeliveryCost: dbo.PaymentDeliveryCost,
			GoodsTotal:   dbo.PaymentGoodsTotal,
			CustomFee:    dbo.PaymentCustomFee,
		},
	}
}

// MapOrdersDBToModels преобразует срез OrderDB в срез Order
func MapOrdersDBToModels(dbOrders []models.OrderDB) []models.Order {
	orders := make([]models.Order, 0, len(dbOrders))
	for _, dbo := range dbOrders {
		orders = append(orders, MapOrderDBToModel(dbo))
	}
	return orders
}
