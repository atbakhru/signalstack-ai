package adapters

import (
	"testing"
	"time"
)

func TestSpaceflightAdapterNormalize(t *testing.T) {
	adapter := SpaceflightAdapter{}
	publishedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := adapter.Normalize(SpaceflightArticle{
		ID:          "space-1",
		Title:       "Space News",
		URL:         "https://example.com/space",
		Summary:     "Summary",
		Body:        "Body",
		Author:      "Editor",
		PublishedAt: &publishedAt,
		NewsSite:    "Spaceflight",
	})

	if document.Source != "spaceflight" || document.Metadata["news_site"] != "Spaceflight" {
		t.Fatalf("unexpected normalized document: %+v", document)
	}
}
