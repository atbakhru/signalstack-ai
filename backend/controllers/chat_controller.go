package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/services"
)

type ChatController struct {
	Service        services.ChatService
	UserRepository repositories.UserRepository
	AuthService    services.AuthService
	JWTSecret      string
	DemoEmail      string
	DemoName       string
	DemoPassword   string
}

func (c ChatController) Ask(ctx *gin.Context) {
	var request services.ChatAskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if c.Service.ChatRepository.Pool == nil || c.UserRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "chat ask stub"})
		return
	}

	userID, err := c.resolveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := c.Service.Ask(ctx.Request.Context(), userID, request)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "openai") {
			status = http.StatusBadGateway
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c ChatController) ListSessions(ctx *gin.Context) {
	if c.Service.ChatRepository.Pool == nil || c.UserRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "chat sessions stub"})
		return
	}
	userID, err := c.resolveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessions, err := c.Service.ChatRepository.ListSessions(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (c ChatController) GetSession(ctx *gin.Context) {
	if c.Service.ChatRepository.Pool == nil || c.UserRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "chat session detail stub"})
		return
	}
	userID, err := c.resolveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	session, err := c.Service.ChatRepository.GetSessionByID(ctx.Request.Context(), sessionID)
	if err != nil {
		if repositories.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session.UserID != userID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	messages, err := c.Service.ChatRepository.ListMessages(ctx.Request.Context(), sessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session": session, "messages": messages})
}

func (c ChatController) DeleteSession(ctx *gin.Context) {
	if c.Service.ChatRepository.Pool == nil || c.UserRepository.Pool == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "chat session delete stub"})
		return
	}
	userID, err := c.resolveUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	if err := c.Service.ChatRepository.DeleteSession(ctx.Request.Context(), sessionID, userID); err != nil {
		if repositories.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c ChatController) resolveUserID(ctx *gin.Context) (int64, error) {
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
	if !repositories.IsNotFound(err) {
		return 0, err
	}

	hashedPassword, err := c.AuthService.HashPassword(password)
	if err != nil {
		return 0, err
	}

	created, createErr := c.UserRepository.Create(ctx.Request.Context(), models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	})
	if createErr == nil {
		return created.ID, nil
	}

	user, err = c.UserRepository.GetByEmail(ctx.Request.Context(), email)
	if err != nil {
		return 0, createErr
	}
	return user.ID, nil
}
