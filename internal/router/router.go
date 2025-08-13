package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shenikar/order-service/internal/handler"
)

func SetupRoutes(engine *gin.Engine, orderHandler *handler.OrderHandler) {
	engine.LoadHTMLFiles("web/index.html")

	// Основные маршруты
	engine.GET("/", orderHandler.Index)
	engine.GET("/orders/:order_uid", orderHandler.GetOrderByUID)
}
