package routes

import (
	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterAuthRoutes(group *gin.RouterGroup, _ config.Config, controller controllers.AuthController) {
	group.POST("/register", controller.Register)
	group.POST("/login", controller.Login)
	group.GET("/me", controller.Me)
}
