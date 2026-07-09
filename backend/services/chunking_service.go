package services

import (
	"strings"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/utils"
)

type Chunk struct {
	ChunkIndex int
	Content    string
	TokenCount int
}

type ChunkingService struct {
	ChunkSizeTokens int
	OverlapTokens   int
}

func NewChunkingService(chunkSizeTokens, overlapTokens int) ChunkingService {
	if chunkSizeTokens <= 0 {
		chunkSizeTokens = 700
	}
	if overlapTokens < 0 {
		overlapTokens = 0
	}
	if overlapTokens >= chunkSizeTokens {
		overlapTokens = chunkSizeTokens / 5
	}

	return ChunkingService{
		ChunkSizeTokens: chunkSizeTokens,
		OverlapTokens:   overlapTokens,
	}
}

func (s ChunkingService) ChunkText(text string) []Chunk {
	cleaned := utils.CleanText(text)
	if cleaned == "" {
		return nil
	}

	words := strings.Fields(cleaned)
	if len(words) == 0 {
		return nil
	}

	chunkSize := s.ChunkSizeTokens
	if chunkSize <= 0 {
		chunkSize = 700
	}
	overlap := s.OverlapTokens
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 5
	}

	step := chunkSize - overlap
	if step <= 0 {
		step = chunkSize
	}

	chunks := make([]Chunk, 0)
	for start := 0; start < len(words); start += step {
		end := start + chunkSize
		if end > len(words) {
			end = len(words)
		}

		chunkWords := words[start:end]
		if len(chunkWords) == 0 {
			break
		}

		chunkContent := strings.Join(chunkWords, " ")
		chunks = append(chunks, Chunk{
			ChunkIndex: len(chunks),
			Content:    chunkContent,
			TokenCount: len(chunkWords),
		})

		if end == len(words) {
			break
		}
	}

	return chunks
}

func (s ChunkingService) ChunkDocument(document models.NormalizedDocument) []Chunk {
	return s.ChunkText(document.Content)
}
