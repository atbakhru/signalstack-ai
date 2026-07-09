package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/utils"
)

type ChatAskRequest struct {
	Question  string  `json:"question"`
	SessionID *int64  `json:"session_id,omitempty"`
	TopK      int     `json:"top_k,omitempty"`
	SourceIDs []int64 `json:"source_ids,omitempty"`
}

type ChatAskResult struct {
	Session            models.ChatSession `json:"session"`
	UserMessage        models.ChatMessage `json:"user_message"`
	AssistantMessage   models.ChatMessage `json:"assistant_message"`
	Citations          []models.Citation  `json:"citations"`
	RetrievedContexts  []RetrievedChunk   `json:"retrieved_contexts"`
	Answer             string             `json:"answer"`
	UsedContext        bool               `json:"used_context"`
	Insufficient       bool               `json:"insufficient"`
	RetrievalLatencyMS int                `json:"retrieval_latency_ms"`
	AnswerLatencyMS    int                `json:"answer_latency_ms"`
}

type ChatDependencies struct {
	ChatRepository     repositories.ChatRepository
	CitationRepository repositories.CitationRepository
	DocumentRepository repositories.DocumentRepository
	ChunkRepository    repositories.ChunkRepository
	Embedding          EmbeddingService
	RAG                RAGService
	ChatCompletion     ChatCompletionProvider
}

type ChatService struct {
	ChatRepository     repositories.ChatRepository
	CitationRepository repositories.CitationRepository
	DocumentRepository repositories.DocumentRepository
	ChunkRepository    repositories.ChunkRepository
	Embedding          EmbeddingService
	RAG                RAGService
	ChatCompletion     ChatCompletionProvider
}

func NewChatService(deps ChatDependencies) ChatService {
	return ChatService{
		ChatRepository:     deps.ChatRepository,
		CitationRepository: deps.CitationRepository,
		DocumentRepository: deps.DocumentRepository,
		ChunkRepository:    deps.ChunkRepository,
		Embedding:          deps.Embedding,
		RAG:                deps.RAG,
		ChatCompletion:     deps.ChatCompletion,
	}
}

func (s ChatService) Ask(ctx context.Context, userID int64, request ChatAskRequest) (ChatAskResult, error) {
	question := utils.CleanText(request.Question)
	if question == "" {
		question = "What information is available in the indexed sources?"
	}

	session, err := s.resolveSession(ctx, userID, request, question)
	if err != nil {
		return ChatAskResult{}, err
	}

	userMessage, err := s.ChatRepository.CreateMessage(ctx, models.ChatMessage{
		SessionID: session.ID,
		UserID:    userID,
		Role:      "user",
		Content:   question,
	})
	if err != nil {
		return ChatAskResult{}, err
	}

	retrievalStarted := time.Now()
	queryEmbedding, err := s.Embedding.Embed(ctx, question)
	if err != nil {
		return ChatAskResult{}, err
	}

	matches, err := s.ChunkRepository.SearchByEmbeddingWithMetadata(ctx, queryEmbedding, request.TopK, request.SourceIDs)
	if err != nil {
		return ChatAskResult{}, err
	}
	retrievalLatencyMS := int(time.Since(retrievalStarted).Milliseconds())

	contexts := make([]RetrievedChunk, 0, len(matches))
	for _, match := range matches {
		contexts = append(contexts, RetrievedChunk{
			Chunk:           match.Chunk,
			DocumentTitle:   match.DocumentTitle,
			SourceName:      match.SourceName,
			URL:             match.URL,
			SimilarityScore: match.SimilarityScore,
		})
	}

	if len(contexts) == 0 {
		answer := s.RAG.BuildInsufficientResponse().Answer
		assistantMessage, err := s.ChatRepository.CreateMessage(ctx, models.ChatMessage{
			SessionID: session.ID,
			UserID:    userID,
			Role:      "assistant",
			Content:   answer,
		})
		if err != nil {
			return ChatAskResult{}, err
		}
		return ChatAskResult{
			Session:            session,
			UserMessage:        userMessage,
			AssistantMessage:   assistantMessage,
			Answer:             answer,
			UsedContext:        false,
			Insufficient:       true,
			RetrievalLatencyMS: retrievalLatencyMS,
		}, nil
	}

	prompt := s.RAG.BuildPrompt(question, contexts)
	answerStarted := time.Now()
	answer := s.buildAnswer(ctx, prompt, contexts)
	answerLatencyMS := int(time.Since(answerStarted).Milliseconds())

	assistantMessage, err := s.ChatRepository.CreateMessage(ctx, models.ChatMessage{
		SessionID: session.ID,
		UserID:    userID,
		Role:      "assistant",
		Content:   answer,
	})
	if err != nil {
		return ChatAskResult{}, err
	}

	citations := make([]models.Citation, 0, len(contexts))
	for _, contextChunk := range contexts {
		citations = append(citations, models.Citation{
			MessageID:      assistantMessage.ID,
			DocumentID:     contextChunk.Chunk.DocumentID,
			ChunkID:        contextChunk.Chunk.ID,
			SourceName:     contextChunk.SourceName,
			Title:          contextChunk.DocumentTitle,
			URL:            contextChunk.URL,
			RelevanceScore: contextChunk.SimilarityScore,
		})
	}

	savedCitations, err := s.CitationRepository.CreateMany(ctx, citations)
	if err != nil {
		return ChatAskResult{}, err
	}

	return ChatAskResult{
		Session:            session,
		UserMessage:        userMessage,
		AssistantMessage:   assistantMessage,
		Citations:          savedCitations,
		RetrievedContexts:  contexts,
		Answer:             assistantMessage.Content,
		UsedContext:        true,
		Insufficient:       false,
		RetrievalLatencyMS: retrievalLatencyMS,
		AnswerLatencyMS:    answerLatencyMS,
	}, nil
}

func (s ChatService) resolveSession(ctx context.Context, userID int64, request ChatAskRequest, question string) (models.ChatSession, error) {
	if request.SessionID != nil {
		session, err := s.ChatRepository.GetSessionByID(ctx, *request.SessionID)
		if err != nil {
			return models.ChatSession{}, err
		}
		if session.UserID != userID {
			return models.ChatSession{}, fmt.Errorf("session does not belong to the current user")
		}
		return session, nil
	}

	return s.ChatRepository.CreateSession(ctx, models.ChatSession{
		UserID: userID,
		Title:  buildSessionTitle(question),
	})
}

func (s ChatService) buildAnswer(ctx context.Context, prompt string, contexts []RetrievedChunk) string {
	if s.ChatCompletion != nil {
		if answer, err := s.ChatCompletion.Complete(ctx, prompt); err == nil && strings.TrimSpace(answer) != "" {
			return strings.TrimSpace(answer)
		}
	}

	var builder strings.Builder
	builder.WriteString("I found these relevant sources:\n")
	limit := len(contexts)
	if limit > 3 {
		limit = 3
	}
	for index := 0; index < limit; index++ {
		contextChunk := contexts[index]
		builder.WriteString(fmt.Sprintf("%d. %s (%s): %s\n", index+1, contextChunk.DocumentTitle, contextChunk.SourceName, summarizeText(contextChunk.Chunk.Content, 220)))
	}
	return strings.TrimSpace(builder.String())
}

func buildSessionTitle(question string) string {
	question = utils.CleanText(question)
	if question == "" {
		return "New chat"
	}
	words := strings.Fields(question)
	if len(words) > 8 {
		words = words[:8]
	}
	title := strings.Join(words, " ")
	if len(title) > 72 {
		title = title[:72]
	}
	return title
}

func summarizeText(text string, limit int) string {
	cleaned := utils.CleanText(text)
	if limit <= 0 || len(cleaned) <= limit {
		return cleaned
	}
	return strings.TrimSpace(cleaned[:limit]) + "..."
}
