package services

import (
	"testing"
	"time"
)

func TestNormalizationServiceBuildContentAndHash(t *testing.T) {
	service := NormalizationService{}
	publishedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := service.NormalizeDocument(
		"guardian",
		"abc-123",
		"  Example Title  ",
		" Example Summary ",
		"https://example.com/article",
		" Example Author ",
		&publishedAt,
		service.BuildContent("Title", "Summary", "Body text"),
		map[string]string{"section": "technology"},
	)

	if document.Title != "Example Title" {
		t.Fatalf("unexpected title: %q", document.Title)
	}

	if document.Content != "Title Summary Body text" {
		t.Fatalf("unexpected content: %q", document.Content)
	}

	if service.BuildContentHash(document) == "" {
		t.Fatal("expected non-empty content hash")
	}
}
