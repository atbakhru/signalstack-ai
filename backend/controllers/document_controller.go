package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
)

type DocumentController struct {
	DocumentRepository repositories.DocumentRepository
	ChunkRepository    repositories.ChunkRepository
	SourceRepository   repositories.SourceRepository
}

func (c DocumentController) List(ctx *gin.Context) {
	var sourceID *int64
	if raw := ctx.Query("source_id"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err == nil {
			sourceID = &value
		}
	}
	filter := repositories.DocumentListFilter{SourceID: sourceID, Search: ctx.Query("search"), Limit: 100, Offset: 0}
	documents, err := c.DocumentRepository.List(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type item struct {
		models.Document
		SourceName string `json:"source_name"`
	}
	results := make([]item, 0, len(documents))
	for _, document := range documents {
		source, sourceErr := c.SourceRepository.GetByID(ctx.Request.Context(), document.SourceID)
		sourceName := ""
		if sourceErr == nil {
			sourceName = source.Name
		}
		results = append(results, item{Document: document, SourceName: sourceName})
	}
	ctx.JSON(http.StatusOK, gin.H{"documents": results})
}

func (c DocumentController) Get(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	document, err := c.DocumentRepository.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if repositories.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	source, sourceErr := c.SourceRepository.GetByID(ctx.Request.Context(), document.SourceID)
	ctx.JSON(http.StatusOK, gin.H{"document": gin.H{"data": document, "source_name": func() string {
		if sourceErr == nil {
			return source.Name
		}
		return ""
	}()}})
}

func (c DocumentController) Chunks(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	chunks, err := c.ChunkRepository.ListByDocumentID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"chunks": chunks})
}

func (c DocumentController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	if err := c.DocumentRepository.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
