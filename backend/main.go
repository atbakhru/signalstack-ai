package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"signalstack-ai/backend/adapters"
	"signalstack-ai/backend/config"
	"signalstack-ai/backend/controllers"
	"signalstack-ai/backend/middleware"
	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/routes"
	"signalstack-ai/backend/services"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	pool, err := config.ConnectDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	deps, err := buildDependencies(ctx, cfg, pool)
	if err != nil {
		log.Fatal(err)
	}

	router := NewRouter(cfg, deps)

	log.Printf("SignalStack AI listening on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal(err)
	}
}

func NewRouter(cfg config.Config, deps routes.Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.CORSMiddleware())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "signalstack-ai",
		})
	})

	api := router.Group("/api")
	routes.RegisterAllRoutes(api, cfg, deps)

	return router
}

func buildDependencies(ctx context.Context, cfg config.Config, pool *pgxpool.Pool) (routes.Dependencies, error) {
	awsConfig, err := config.LoadAWSConfig(ctx, cfg.AWSRegion)
	if err != nil {
		return routes.Dependencies{}, err
	}

	openAIClient := services.NewOpenAIClient(config.NewOpenAIConfig(cfg))
	embeddingService := services.EmbeddingService{}
	var chatCompletion services.ChatCompletionProvider
	if cfg.OpenAIAPIKey != "" {
		embeddingService.Provider = openAIClient
		chatCompletion = openAIClient
	}

	chatService := services.NewChatService(services.ChatDependencies{
		ChatRepository:     repositories.ChatRepository{Pool: pool},
		CitationRepository: repositories.CitationRepository{Pool: pool},
		DocumentRepository: repositories.DocumentRepository{Pool: pool},
		ChunkRepository:    repositories.ChunkRepository{Pool: pool},
		Embedding:          embeddingService,
		RAG:                services.RAGService{},
		ChatCompletion:     chatCompletion,
	})

	authController := controllers.AuthController{
		UserRepository: repositories.UserRepository{Pool: pool},
		AuthService:    services.AuthService{},
		JWTService: services.JWTService{
			Secret:   cfg.JWTSecret,
			Issuer:   "signalstack-ai",
			Duration: 24 * time.Hour,
		},
	}

	documentController := controllers.DocumentController{
		DocumentRepository: repositories.DocumentRepository{Pool: pool},
		ChunkRepository:    repositories.ChunkRepository{Pool: pool},
		SourceRepository:   repositories.SourceRepository{Pool: pool},
	}

	evaluationController := controllers.EvaluationController{
		Repository:     repositories.EvaluationRepository{Pool: pool},
		UserRepository: repositories.UserRepository{Pool: pool},
		AuthService:    services.AuthService{},
		ChatService:    chatService,
		JWTSecret:      cfg.JWTSecret,
		DemoEmail:      "demo@signalstack.local",
		DemoName:       "Demo User",
		DemoPassword:   "demo-password",
	}

	metricsController := controllers.MetricsController{
		Pool:                 pool,
		SourceRepository:     repositories.SourceRepository{Pool: pool},
		DocumentRepository:   repositories.DocumentRepository{Pool: pool},
		ChunkRepository:      repositories.ChunkRepository{Pool: pool},
		IngestionRepository:  repositories.IngestionRepository{Pool: pool},
		EvaluationRepository: repositories.EvaluationRepository{Pool: pool},
		ChatRepository:       repositories.ChatRepository{Pool: pool},
	}

	ingestionService := services.NewIngestionService(cfg, services.IngestionDependencies{
		SourceRepository:    repositories.SourceRepository{Pool: pool},
		DocumentRepository:  repositories.DocumentRepository{Pool: pool},
		ChunkRepository:     repositories.ChunkRepository{Pool: pool},
		IngestionRepository: repositories.IngestionRepository{Pool: pool},
		S3Service: services.S3Service{
			Client: s3.NewFromConfig(awsConfig),
			Bucket: cfg.S3Bucket,
		},
		Adapters: map[string]any{
			"gdelt":       adapters.GDELTAdapter{},
			"guardian":    adapters.GuardianAdapter{},
			"hackernews":  adapters.HackerNewsAdapter{},
			"arxiv":       adapters.ArxivAdapter{},
			"spaceflight": adapters.SpaceflightAdapter{},
		},
	})

	return routes.Dependencies{
		AuthController:       authController,
		DocumentController:   documentController,
		IngestionController:  controllers.IngestionController{Service: ingestionService},
		EvaluationController: evaluationController,
		MetricsController:    metricsController,
		ChatController: controllers.ChatController{
			Service:        chatService,
			UserRepository: repositories.UserRepository{Pool: pool},
			AuthService:    services.AuthService{},
			JWTSecret:      cfg.JWTSecret,
			DemoEmail:      "demo@signalstack.local",
			DemoName:       "Demo User",
			DemoPassword:   "demo-password",
		},
	}, nil
}
