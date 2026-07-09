package services

import (
	"testing"
	"time"
)

func TestJWTServiceGenerateAndParse(t *testing.T) {
	service := JWTService{
		Secret:   "test-secret",
		Issuer:   "signalstack-ai",
		Duration: time.Hour,
	}

	token, err := service.GenerateToken(42, "user@example.com")
	if err != nil {
		t.Fatalf("unexpected token generation error: %v", err)
	}

	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
}
