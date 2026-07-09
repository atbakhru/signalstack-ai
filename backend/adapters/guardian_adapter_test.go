package adapters

import (
	"testing"
	"time"
)

func TestGuardianAdapterNormalize(t *testing.T) {
	adapter := GuardianAdapter{}
	publishedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := adapter.Normalize(GuardianArticle{
		ID:          "guardian-1",
		Title:       "Guardian Title",
		URL:         "https://example.com/guardian",
		Section:     "technology",
		Summary:     "Summary",
		Body:        "Body",
		Author:      "Writer",
		PublishedAt: &publishedAt,
	})

	if document.Source != "guardian" || document.Metadata["section"] != "technology" {
		t.Fatalf("unexpected normalized document: %+v", document)
	}
}
