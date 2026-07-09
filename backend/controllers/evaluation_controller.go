package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/services"
)

type EvaluationController struct {
	Repository     repositories.EvaluationRepository
	UserRepository repositories.UserRepository
	AuthService    services.AuthService
	ChatService    services.ChatService
	JWTSecret      string
	DemoEmail      string
	DemoName       string
	DemoPassword   string
}

func (c EvaluationController) Run(ctx *gin.Context) {
	if c.Repository.Pool == nil || c.ChatService.ChatRepository.Pool == nil || c.UserRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"summary": services.EvaluationSummary{}, "results": []any{}})
		return
	}
	questions, err := c.Repository.ListQuestions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userID, err := c.resolveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]models.EvaluationResult, 0, len(questions))
	for _, question := range questions {
		response, askErr := c.ChatService.Ask(ctx.Request.Context(), userID, services.ChatAskRequest{Question: question.Question, TopK: 5})
		if askErr != nil {
			continue
		}
		retrievedIDs := make([]int64, 0, len(response.Citations))
		topKAccuracy := 0.0
		for _, citation := range response.Citations {
			retrievedIDs = append(retrievedIDs, citation.DocumentID)
			if strings.EqualFold(citation.SourceName, question.ExpectedSource) {
				topKAccuracy = 1
			}
		}
		result := models.EvaluationResult{
			QuestionID:           question.ID,
			RetrievedDocumentIDs: retrievedIDs,
			TopKAccuracy:         topKAccuracy,
			AnswerLatencyMS:      response.AnswerLatencyMS,
			RetrievalLatencyMS:   response.RetrievalLatencyMS,
			CitationCount:        len(response.Citations),
		}
		created, createErr := c.Repository.CreateResult(ctx.Request.Context(), result)
		if createErr == nil {
			results = append(results, created)
		}
	}
	summary := services.EvaluationService{}.BuildSummary(results)
	ctx.JSON(http.StatusOK, gin.H{"summary": summary, "results": results})
}

func (c EvaluationController) Results(ctx *gin.Context) {
	if c.Repository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"summary": services.EvaluationSummary{}, "results": []any{}})
		return
	}
	results, err := c.Repository.ListResults(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	summary := services.EvaluationService{}.BuildSummary(results)
	ctx.JSON(http.StatusOK, gin.H{"summary": summary, "results": results})
}

func (c EvaluationController) Questions(ctx *gin.Context) {
	if c.Repository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"questions": []any{}})
		return
	}
	questions, err := c.Repository.ListQuestions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (c EvaluationController) resolveUserID(ctx *gin.Context) (int64, error) {
	head := ctx.GetHeader("Authorization")
	if strings.HasPrefix(head, "Bearer ") && c.JWTSecret != "" {
		token := strings.TrimSpace(strings.TrimPrefix(head, "Bearer "))
		claims, err := services.JWTService{Secret: c.JWTSecret}.ParseToken(token)
		if err == nil && claims != nil && claims.UserID > 0 {
			return claims.UserID, nil
		}
	}

	email := c.DemoEmail
	if email == "" {
		email = "demo@signalstack.local"
	}
	name := c.DemoName
	if name == "" {
		name = "Demo User"
	}
	password := c.DemoPassword
	if password == "" {
		password = "demo-password"
	}

	user, err := c.UserRepository.GetByEmail(ctx.Request.Context(), email)
	if err == nil {
		return user.ID, nil
	}
	hashedPassword, err := c.AuthService.HashPassword(password)
	if err != nil {
		return 0, err
	}
	created, createErr := c.UserRepository.Create(ctx.Request.Context(), models.User{Name: name, Email: email, PasswordHash: hashedPassword})
	if createErr == nil {
		return created.ID, nil
	}
	return 0, createErr
}
