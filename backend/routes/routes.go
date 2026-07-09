package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

type Dependencies struct {
	AuthController       controllers.AuthController
	DocumentController   controllers.DocumentController
	IngestionController  controllers.IngestionController
	EvaluationController controllers.EvaluationController
	MetricsController    controllers.MetricsController
	ChatController       controllers.ChatController
}

func RegisterAllRoutes(api *gin.RouterGroup, cfg config.Config, deps Dependencies) {
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	RegisterAuthRoutes(api.Group("/auth"), cfg, deps.AuthController)
	RegisterDocumentRoutes(api.Group("/documents"), cfg, deps.DocumentController)
	RegisterIngestionRoutes(api.Group("/ingest"), cfg, deps.IngestionController)
	RegisterChatRoutes(api.Group("/chat"), cfg, deps.ChatController)
	RegisterEvaluationRoutes(api.Group("/evaluate"), cfg, deps.EvaluationController)
	RegisterMetricsRoutes(api.Group("/metrics"), cfg, deps.MetricsController)
}
