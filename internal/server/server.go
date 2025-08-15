package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/handler"
	"github.com/shenikar/order-service/internal/router"
	"github.com/shenikar/order-service/internal/service"
)

var httpServer *http.Server

func StartServer(config *config.Config, orderService *service.OrderService) {
	r := gin.Default()

	// Создаем обработчик
	orderHandler := handler.NewOrderHandler(orderService)

	// настраиваем маршруты
	router.SetupRoutes(r, orderHandler)

	// запускаем сервер
	addr := config.GetServerAddress()
	log.Printf("Starting server on %s\n", addr)

	httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

// ShutdownServer корректно завершает работу HTTP сервера
func ShutdownServer(ctx context.Context) error {
	if httpServer == nil {
		return nil
	}
	log.Println("Shutting down server")
	return httpServer.Shutdown(ctx)

}
