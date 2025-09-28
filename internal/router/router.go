package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/shenikar/order-service/docs" // Import Swagger docs for router initialization
	"github.com/shenikar/order-service/internal/handler"
	"github.com/shenikar/order-service/internal/metrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(engine *gin.Engine, orderHandler *handler.OrderHandler, metricsMiddleware gin.HandlerFunc) {
	engine.LoadHTMLFiles("web/index.html")

	// Группа для метрик - без нашего middleware
	engine.GET("/metrics", metrics.PrometheusHandler())

	// Группа для API - с middleware для метрик
	apiGroup := engine.Group("/")
	apiGroup.Use(metricsMiddleware)
	{
		apiGroup.GET("/", orderHandler.Index)
		apiGroup.GET("/orders/:order_uid", orderHandler.GetOrderByUID)
		apiGroup.GET("/health", orderHandler.HealthCheck)
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
