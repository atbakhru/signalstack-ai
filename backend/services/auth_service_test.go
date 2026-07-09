package services

import "testing"

func TestAuthServiceHashAndCompare(t *testing.T) {
	service := AuthService{}

	hashed, err := service.HashPassword("secret-password")
	if err != nil {
		t.Fatalf("unexpected hash error: %v", err)
	}

	if err := service.ComparePassword(hashed, "secret-password"); err != nil {
		t.Fatalf("expected password comparison to succeed: %v", err)
	}
}
