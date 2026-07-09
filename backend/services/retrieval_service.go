package services

import (
	"math"
	"sort"

	"signalstack-ai/backend/models"
)

type RetrievedChunk struct {
	Chunk           models.DocumentChunk `json:"chunk"`
	DocumentTitle   string               `json:"document_title"`
	SourceName      string               `json:"source_name"`
	URL             string               `json:"url"`
	SimilarityScore float64              `json:"similarity_score"`
}

type RetrievalService struct {
	DefaultK              int
	MinimumRelevanceScore float64
}

func NewRetrievalService() RetrievalService {
	return RetrievalService{
		DefaultK:              5,
		MinimumRelevanceScore: 0.20,
	}
}

func (s RetrievalService) CosineSimilarity(left, right []float32) float64 {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}

	length := len(left)
	if len(right) < length {
		length = len(right)
	}

	var dot, leftNorm, rightNorm float64
	for index := 0; index < length; index++ {
		leftValue := float64(left[index])
		rightValue := float64(right[index])
		dot += leftValue * rightValue
		leftNorm += leftValue * leftValue
		rightNorm += rightValue * rightValue
	}

	if leftNorm == 0 || rightNorm == 0 {
		return 0
	}

	return dot / (math.Sqrt(leftNorm) * math.Sqrt(rightNorm))
}

func (s RetrievalService) RankChunks(queryEmbedding []float32, chunks []RetrievedChunk, topK int, minimumScore float64) []RetrievedChunk {
	if topK <= 0 {
		topK = s.DefaultK
	}
	if minimumScore <= 0 {
		minimumScore = s.MinimumRelevanceScore
	}

	for index := range chunks {
		chunks[index].SimilarityScore = s.CosineSimilarity(queryEmbedding, chunks[index].Chunk.Embedding)
	}

	filtered := make([]RetrievedChunk, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk.SimilarityScore >= minimumScore {
			filtered = append(filtered, chunk)
		}
	}

	sort.Slice(filtered, func(leftIndex, rightIndex int) bool {
		return filtered[leftIndex].SimilarityScore > filtered[rightIndex].SimilarityScore
	})

	if len(filtered) > topK {
		filtered = filtered[:topK]
	}
	return filtered
}
