package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/repositories"
	"signalstack-ai/backend/services"
)

type AuthController struct {
	UserRepository repositories.UserRepository
	AuthService    services.AuthService
	JWTService     services.JWTService
}

func (c AuthController) Register(ctx *gin.Context) {
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := c.AuthService.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created, err := c.UserRepository.Create(ctx.Request.Context(), models.User{
		Name:         strings.TrimSpace(request.Name),
		Email:        strings.TrimSpace(strings.ToLower(request.Email)),
		PasswordHash: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.JWTService.GenerateToken(created.ID, created.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": created, "token": token})
}

func (c AuthController) Login(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserRepository.GetByEmail(ctx.Request.Context(), strings.TrimSpace(strings.ToLower(request.Email)))
	if err != nil {
		if repositories.IsNotFound(err) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.AuthService.ComparePassword(user.PasswordHash, request.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := c.JWTService.GenerateToken(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (c AuthController) Me(ctx *gin.Context) {
	head := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(head, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		return
	}

	claims, err := c.JWTService.ParseToken(strings.TrimSpace(strings.TrimPrefix(head, "Bearer ")))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserRepository.GetByID(ctx.Request.Context(), claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
