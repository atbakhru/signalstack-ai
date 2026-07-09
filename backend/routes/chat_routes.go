package routes

import (
	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterChatRoutes(group *gin.RouterGroup, cfg config.Config, controller controllers.ChatController) {
	group.POST("/ask", controller.Ask)
	group.GET("/sessions", controller.ListSessions)
	group.GET("/sessions/:id", controller.GetSession)
	group.DELETE("/sessions/:id", controller.DeleteSession)
	_ = cfg
}
