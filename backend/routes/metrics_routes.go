package routes

import (
	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterMetricsRoutes(group *gin.RouterGroup, _ config.Config, controller controllers.MetricsController) {
	group.GET("/overview", controller.Overview)
	group.GET("/sources", controller.Sources)
	group.GET("/latency", controller.Latency)
	group.GET("/embeddings", controller.Embeddings)
}
