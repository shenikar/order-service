package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/shenikar/order-service/docs" // Import Swagger docs for router initialization
	"github.com/shenikar/order-service/internal/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(engine *gin.Engine, orderHandler *handler.OrderHandler) {
	engine.LoadHTMLFiles("web/index.html")

	// Основные маршруты
	engine.GET("/", orderHandler.Index)
	engine.GET("/orders/:order_uid", orderHandler.GetOrderByUID)
	engine.GET("/health", orderHandler.HealthCheck)

	// Swagger UI
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
