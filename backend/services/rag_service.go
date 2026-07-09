package services

import (
	"fmt"
	"strings"

	"signalstack-ai/backend/models"
)

type RAGResponse struct {
	Answer       string            `json:"answer"`
	Citations    []models.Citation `json:"citations"`
	UsedContext  bool              `json:"used_context"`
	Insufficient bool              `json:"insufficient"`
}

type RAGService struct{}

func (s RAGService) SystemPrompt() string {
	return "You are SignalStack AI. Answer the user’s question using only the retrieved source context. Include citations for every factual claim. If the context is insufficient, say you could not find enough relevant information in the indexed sources. Do not invent facts."
}

func (s RAGService) BuildPrompt(question string, contexts []RetrievedChunk) string {
	var builder strings.Builder
	builder.WriteString(s.SystemPrompt())
	builder.WriteString("\n\nQuestion: ")
	builder.WriteString(question)
	builder.WriteString("\n\nRetrieved context:\n")

	for index, context := range contexts {
		builder.WriteString(fmt.Sprintf("[%d] %s\nSource: %s\nURL: %s\nScore: %.4f\nChunk: %s\n\n",
			index+1,
			context.DocumentTitle,
			context.SourceName,
			context.URL,
			context.SimilarityScore,
			context.Chunk.Content,
		))
	}

	builder.WriteString("Only answer from the retrieved context and cite each factual claim.")
	return builder.String()
}

func (s RAGService) BuildInsufficientResponse() RAGResponse {
	return RAGResponse{
		Answer:       "I could not find enough relevant information in the indexed sources.",
		UsedContext:  false,
		Insufficient: true,
	}
}
