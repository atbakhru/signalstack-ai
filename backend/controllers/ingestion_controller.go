package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/services"
)

type IngestionController struct {
	Service services.IngestionService
}

func (c IngestionController) IngestAll(ctx *gin.Context) {
	if c.Service.SourceRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"results": []any{}})
		return
	}
	result, err := c.Service.IngestAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"results": result})
}

func (c IngestionController) IngestSource(ctx *gin.Context) {
	if c.Service.SourceRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "ingestion stub"})
		return
	}
	source := ctx.Param("source")
	result, err := c.Service.IngestSourceByName(ctx.Request.Context(), source)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
