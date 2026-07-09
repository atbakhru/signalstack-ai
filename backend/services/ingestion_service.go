package services

import (
	"context"
	"fmt"
	"strings"

	"signalstack-ai/backend/adapters"
	"signalstack-ai/backend/config"
	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
)

type IngestionService struct {
	Config              config.Config
	SourceRepository    repositories.SourceRepository
	DocumentRepository  repositories.DocumentRepository
	ChunkRepository     repositories.ChunkRepository
	IngestionRepository repositories.IngestionRepository
	S3Service           S3Service
	Normalization       NormalizationService
	Chunking            ChunkingService
	Embedding           EmbeddingService
	Adapters            map[string]any
}

type IngestionOutcome struct {
	Run               models.IngestionRun `json:"run"`
	Source            models.Source       `json:"source"`
	DocumentsFetched  int                 `json:"documents_fetched"`
	DocumentsInserted int                 `json:"documents_inserted"`
	ChunksCreated     int                 `json:"chunks_created"`
	EmbeddingsCreated int                 `json:"embeddings_created"`
}

type IngestionDependencies struct {
	SourceRepository    repositories.SourceRepository
	DocumentRepository  repositories.DocumentRepository
	ChunkRepository     repositories.ChunkRepository
	IngestionRepository repositories.IngestionRepository
	S3Service           S3Service
	Adapters            map[string]any
}

func NewIngestionService(cfg config.Config, deps IngestionDependencies) IngestionService {
	return IngestionService{
		Config:              cfg,
		SourceRepository:    deps.SourceRepository,
		DocumentRepository:  deps.DocumentRepository,
		ChunkRepository:     deps.ChunkRepository,
		IngestionRepository: deps.IngestionRepository,
		S3Service:           deps.S3Service,
		Normalization:       NormalizationService{},
		Chunking:            NewChunkingService(700, 120),
		Embedding:           EmbeddingService{},
		Adapters:            deps.Adapters,
	}
}

func (s IngestionService) IngestAll(ctx context.Context) ([]IngestionOutcome, error) {
	sources, err := s.SourceRepository.ListEnabled(ctx)
	if err != nil {
		return nil, err
	}

	outcomes := make([]IngestionOutcome, 0, len(sources))
	for _, source := range sources {
		outcome, err := s.IngestSourceByName(ctx, source.Name)
		if err != nil {
			return outcomes, err
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}

func (s IngestionService) IngestSourceByName(ctx context.Context, sourceName string) (IngestionOutcome, error) {
	source, err := s.SourceRepository.GetByName(ctx, sourceName)
	if err != nil {
		return IngestionOutcome{}, err
	}

	run, err := s.IngestionRepository.CreateRun(ctx, source.ID)
	if err != nil {
		return IngestionOutcome{}, err
	}

	outcome := IngestionOutcome{Run: run, Source: source}
	fail := func(runErr error) (IngestionOutcome, error) {
		_, _ = s.IngestionRepository.FinishRun(ctx, run.ID, "failed", outcome.DocumentsFetched, outcome.DocumentsInserted, outcome.ChunksCreated, outcome.EmbeddingsCreated, runErr.Error())
		return outcome, runErr
	}

	documents, rawPayload, fetchErr := s.fetchSource(ctx, source)
	if fetchErr != nil {
		return fail(fetchErr)
	}
	outcome.DocumentsFetched = len(documents)

	for _, normalized := range documents {
		rawS3Key, uploadErr := s.S3Service.UploadJSON(ctx, source.Name, normalized.ExternalID, rawPayload)
		if uploadErr != nil {
			return fail(uploadErr)
		}

		contentHash := s.Normalization.BuildContentHash(normalized)
		document, inserted, upsertErr := s.DocumentRepository.Upsert(ctx, source.ID, normalized, rawS3Key, contentHash)
		if upsertErr != nil {
			return fail(upsertErr)
		}
		if !inserted {
			continue
		}
		outcome.DocumentsInserted++

		chunks := s.Chunking.ChunkDocument(normalized)
		documentChunks := make([]models.DocumentChunk, 0, len(chunks))
		for _, chunk := range chunks {
			embedding, embeddingErr := s.Embedding.Embed(ctx, chunk.Content)
			if embeddingErr != nil {
				return fail(embeddingErr)
			}
			outcome.EmbeddingsCreated++
			documentChunks = append(documentChunks, models.DocumentChunk{
				DocumentID: document.ID,
				SourceID:   source.ID,
				ChunkIndex: chunk.ChunkIndex,
				Content:    chunk.Content,
				TokenCount: chunk.TokenCount,
				Embedding:  embedding,
			})
		}

		savedChunks, chunkErr := s.ChunkRepository.InsertMany(ctx, documentChunks)
		if chunkErr != nil {
			return fail(chunkErr)
		}
		outcome.ChunksCreated += len(savedChunks)
	}

	finishedRun, finishErr := s.IngestionRepository.FinishRun(ctx, run.ID, "completed", outcome.DocumentsFetched, outcome.DocumentsInserted, outcome.ChunksCreated, outcome.EmbeddingsCreated, "")
	if finishErr != nil {
		return fail(finishErr)
	}
	outcome.Run = finishedRun
	return outcome, nil
}

func (s IngestionService) fetchSource(ctx context.Context, source models.Source) ([]models.NormalizedDocument, any, error) {
	adapter, ok := s.Adapters[strings.ToLower(source.Name)]
	if !ok {
		return nil, nil, fmt.Errorf("no adapter registered for source %s", source.Name)
	}

	switch concrete := adapter.(type) {
	case adapters.GDELTAdapter:
		return concrete.Fetch(ctx, source.BaseURL, 10)
	case adapters.GuardianAdapter:
		return concrete.Fetch(ctx, source.BaseURL, s.Config.GuardianAPIKey, 10)
	case adapters.HackerNewsAdapter:
		return concrete.Fetch(ctx, source.BaseURL, 10)
	case adapters.ArxivAdapter:
		return concrete.Fetch(ctx, source.BaseURL, 10)
	case adapters.SpaceflightAdapter:
		return concrete.Fetch(ctx, source.BaseURL, 10)
	default:
		return nil, nil, fmt.Errorf("unsupported adapter type %T", adapter)
	}
}
