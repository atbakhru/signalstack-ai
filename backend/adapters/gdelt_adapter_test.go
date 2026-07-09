package adapters

import (
	"testing"
	"time"
)

func TestGDELTAdapterNormalize(t *testing.T) {
	adapter := GDELTAdapter{}
	publishedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := adapter.Normalize(GDELTArticle{
		Title:       "News Title",
		URL:         "https://example.com/news",
		SourceName:  "Example News",
		ArticleID:   "gdelt-1",
		Summary:     "Summary",
		Author:      "Reporter",
		PublishedAt: &publishedAt,
		Text:        "Body",
	})

	if document.Source != "gdelt" || document.ExternalID != "gdelt-1" {
		t.Fatalf("unexpected normalized document: %+v", document)
	}
}
