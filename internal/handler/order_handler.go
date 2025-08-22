package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/order-service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler создает новый экземпляр OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetOrderBy получает заказ по OrderUID
func (h *OrderHandler) GetOrderByUID(c *gin.Context) {
	orderUID := c.Param("order_uid")

	// Проверка на наличие OrderUID
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order UID is required"})
		return
	}

	order, err := h.orderService.GetOrderByUID(orderUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Index отображает главную страницу
func (h *OrderHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// HealthCheck проверяет состояние сервиса
func (h *OrderHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
