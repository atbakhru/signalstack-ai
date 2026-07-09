package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
)

func RegisterIngestionRoutes(group *gin.RouterGroup, _ config.Config, controller controllers.IngestionController) {
	group.POST("/all", controller.IngestAll)
	group.POST("/:source", controller.IngestSource)
	group.GET("/runs", func(c *gin.Context) {
		if controller.Service.IngestionRepository.Pool == nil {
			c.JSON(http.StatusOK, gin.H{"runs": []any{}})
			return
		}
		runs, err := controller.Service.IngestionRepository.ListRuns(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"runs": runs})
	})
	group.GET("/runs/:id", func(c *gin.Context) {
		if controller.Service.IngestionRepository.Pool == nil {
			c.JSON(http.StatusOK, gin.H{"message": "ingestion run detail stub"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run id"})
			return
		}
		run, err := controller.Service.IngestionRepository.GetRunByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"run": run})
	})
	group.GET("/sources", func(c *gin.Context) {
		if controller.Service.SourceRepository.Pool == nil {
			c.JSON(http.StatusOK, gin.H{"sources": []any{}})
			return
		}
		sources, err := controller.Service.SourceRepository.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"sources": sources})
	})
	group.PATCH("/sources/:id/toggle", func(c *gin.Context) {
		if controller.Service.SourceRepository.Pool == nil {
			c.JSON(http.StatusOK, gin.H{"message": "toggle source stub"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
			return
		}
		source, err := controller.Service.SourceRepository.ToggleEnabled(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"source": source})
	})
}
