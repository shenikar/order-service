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

// GetOrderByUID получает заказ по OrderUID
// @Summary Получить заказ по UID
// @Description Получает заказ с товарами по уникальному идентификатору
// @Tags orders
// @Param order_uid path string true "Order UID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{order_uid} [get]
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

// Index godoc
// @Summary Главная страница сервиса
// @Description Отображает главную страницу
// @Tags general
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router / [get]
func (h *OrderHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// HealthCheck godoc
// @Summary Проверка состояния сервиса
// @Description Возвращает статус сервиса (ok)
// @Tags general
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *OrderHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
