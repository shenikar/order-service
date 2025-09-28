package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/order-service/config"
	"github.com/shenikar/order-service/internal/handler"
	"github.com/shenikar/order-service/internal/metrics"
	"github.com/shenikar/order-service/internal/router"
	"github.com/shenikar/order-service/internal/service"
)

var httpServer *http.Server

func StartServer(cfg *config.Config, orderService *service.OrderService) {
	r := gin.Default()

	// Prometheus middleware
	metricsMiddleware := func(c *gin.Context) {
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "not_found"
		}

		metrics.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			strconv.Itoa(c.Writer.Status()),
		).Inc()
	}

	// Создаем обработчик
	orderHandler := handler.NewOrderHandler(orderService)

	// настраиваем маршруты
	router.SetupRoutes(r, orderHandler, metricsMiddleware)

	// запускаем сервер
	addr := cfg.GetServerAddress()
	log.Printf("Starting server on %s\n", addr)

	httpServer = &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: time.Duration(cfg.Server.ReadHeaderTimeout) * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
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
