package services

import "context"

type EmbeddingProvider interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type EmbeddingService struct {
	Provider EmbeddingProvider
}

func (s EmbeddingService) Embed(ctx context.Context, text string) ([]float32, error) {
	if s.Provider == nil {
		return DeterministicEmbedding(text, 1536), nil
	}
	return s.Provider.Embed(ctx, text)
}

func DeterministicEmbedding(text string, dimensions int) []float32 {
	if dimensions <= 0 {
		dimensions = 1536
	}

	vector := make([]float32, dimensions)
	if text == "" {
		return vector
	}

	bytes := []byte(text)
	for index := range vector {
		value := bytes[index%len(bytes)]
		vector[index] = float32((int(value)%31)-15) / 15.0
	}
	return vector
}
