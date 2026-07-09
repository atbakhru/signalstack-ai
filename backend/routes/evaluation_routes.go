package routes

import (
	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterEvaluationRoutes(group *gin.RouterGroup, _ config.Config, controller controllers.EvaluationController) {
	group.POST("/run", controller.Run)
	group.GET("/results", controller.Results)
	group.GET("/questions", controller.Questions)
}
