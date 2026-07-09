package routes

import (
	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterDocumentRoutes(group *gin.RouterGroup, _ config.Config, controller controllers.DocumentController) {
	group.GET("", controller.List)
	group.GET("/:id", controller.Get)
	group.GET("/:id/chunks", controller.Chunks)
	group.DELETE("/:id", controller.Delete)
}
