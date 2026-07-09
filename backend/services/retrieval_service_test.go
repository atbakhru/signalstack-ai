package services

import (
	"testing"

	"signalstack-ai/backend/models"
)

func TestRetrievalServiceRanksBySimilarity(t *testing.T) {
	service := NewRetrievalService()

	query := []float32{1, 0, 0}
	chunks := []RetrievedChunk{
		{Chunk: models.DocumentChunk{Embedding: []float32{0.9, 0.1, 0}}, DocumentTitle: "A"},
		{Chunk: models.DocumentChunk{Embedding: []float32{0, 1, 0}}, DocumentTitle: "B"},
	}

	ranked := service.RankChunks(query, chunks, 1, 0.1)
	if len(ranked) != 1 {
		t.Fatalf("expected 1 ranked chunk, got %d", len(ranked))
	}

	if ranked[0].DocumentTitle != "A" {
		t.Fatalf("expected chunk A to rank first, got %s", ranked[0].DocumentTitle)
	}
}
