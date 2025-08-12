package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/handler"
	"github.com/shenikar/order-service/internal/router"
	"github.com/shenikar/order-service/internal/service"
)

func StartServer(config *config.Config, orderService *service.OrderService) {
	r := gin.Default()

	// создаем обработчик заказов
	orderHandler := handler.NewOrderHandler(orderService)

	// настраиваем маршруты
	router.SetupRoutes(r, orderHandler)

	// запускаем сервер
	addr := config.GetServerAddress()
	log.Printf("Starting server on %s\n", addr)
	r.Run(addr)
}
