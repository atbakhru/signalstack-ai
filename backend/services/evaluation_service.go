package services

import "signalstack-ai/backend/models"

type EvaluationSummary struct {
	TopKAccuracy       float64 `json:"top_k_accuracy"`
	AverageRetrievalMS int     `json:"average_retrieval_ms"`
	AverageAnswerMS    int     `json:"average_answer_ms"`
	AverageCitations   float64 `json:"average_citations"`
	FailedRetrievals   int     `json:"failed_retrievals"`
}

type EvaluationService struct{}

func (s EvaluationService) BuildSummary(results []models.EvaluationResult) EvaluationSummary {
	if len(results) == 0 {
		return EvaluationSummary{}
	}

	var topKAccuracy float64
	var retrievalMS int
	var answerMS int
	var citations int

	for _, result := range results {
		topKAccuracy += result.TopKAccuracy
		retrievalMS += result.RetrievalLatencyMS
		answerMS += result.AnswerLatencyMS
		citations += result.CitationCount
	}

	count := float64(len(results))
	return EvaluationSummary{
		TopKAccuracy:       topKAccuracy / count,
		AverageRetrievalMS: retrievalMS / len(results),
		AverageAnswerMS:    answerMS / len(results),
		AverageCitations:   float64(citations) / count,
	}
}
