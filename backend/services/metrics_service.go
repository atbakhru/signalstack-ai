package services

type MetricsOverview struct {
	DocumentsIngested   int     `json:"documents_ingested"`
	ChunksGenerated     int     `json:"chunks_generated"`
	EmbeddingsGenerated int     `json:"embeddings_generated"`
	ActiveSources       int     `json:"active_sources"`
	AverageRetrievalMS  float64 `json:"average_retrieval_ms"`
	AverageAnswerMS     float64 `json:"average_answer_ms"`
	TokenUsage          int     `json:"token_usage"`
}

type MetricsService struct{}
