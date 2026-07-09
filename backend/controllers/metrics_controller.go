package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/services"
)

type MetricsController struct {
	Pool                 *pgxpool.Pool
	SourceRepository     repositories.SourceRepository
	DocumentRepository   repositories.DocumentRepository
	ChunkRepository      repositories.ChunkRepository
	IngestionRepository  repositories.IngestionRepository
	EvaluationRepository repositories.EvaluationRepository
	ChatRepository       repositories.ChatRepository
}

func (c MetricsController) Overview(ctx *gin.Context) {
	if c.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"documents_ingested": 0, "chunks_generated": 0, "embeddings_generated": 0, "active_sources": 0, "average_retrieval_ms": 0, "average_answer_ms": 0, "token_usage": 0})
		return
	}
	var documents, chunks, embeddings, activeSources, tokenUsage int
	var avgRetrieval, avgAnswer float64
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COUNT(*) FROM documents`).Scan(&documents)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COUNT(*) FROM document_chunks`).Scan(&chunks)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COUNT(*) FROM document_chunks WHERE embedding IS NOT NULL`).Scan(&embeddings)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COUNT(*) FROM sources WHERE enabled = TRUE`).Scan(&activeSources)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COALESCE(AVG(retrieval_latency_ms), 0) FROM evaluation_results`).Scan(&avgRetrieval)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COALESCE(AVG(answer_latency_ms), 0) FROM evaluation_results`).Scan(&avgAnswer)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COALESCE(SUM(token_count), 0) FROM document_chunks`).Scan(&tokenUsage)
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COALESCE(SUM(cardinality(regexp_split_to_array(content, '\\s+'))), 0) FROM chat_messages`).Scan(&tokenUsage)
	ctx.JSON(http.StatusOK, gin.H{"documents_ingested": documents, "chunks_generated": chunks, "embeddings_generated": embeddings, "active_sources": activeSources, "average_retrieval_ms": avgRetrieval, "average_answer_ms": avgAnswer, "token_usage": tokenUsage})
}

func (c MetricsController) Sources(ctx *gin.Context) {
	sources, err := c.SourceRepository.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sources": sources})
}

func (c MetricsController) Latency(ctx *gin.Context) {
	results, err := c.EvaluationRepository.ListResults(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"summary": gin.H{"average_retrieval_ms": services.EvaluationService{}.BuildSummary(results).AverageRetrievalMS, "average_answer_ms": services.EvaluationService{}.BuildSummary(results).AverageAnswerMS}, "results": results})
}

func (c MetricsController) Embeddings(ctx *gin.Context) {
	if c.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"embeddings_generated": 0})
		return
	}
	var embeddings int
	_ = c.Pool.QueryRow(ctx.Request.Context(), `SELECT COUNT(*) FROM document_chunks WHERE embedding IS NOT NULL`).Scan(&embeddings)
	ctx.JSON(http.StatusOK, gin.H{"embeddings_generated": embeddings})
}
