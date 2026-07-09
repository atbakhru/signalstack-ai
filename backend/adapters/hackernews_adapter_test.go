package adapters

import (
	"testing"
	"time"
)

func TestHackerNewsAdapterNormalize(t *testing.T) {
	adapter := HackerNewsAdapter{}
	postedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := adapter.Normalize(HackerNewsItem{
		ID:       123,
		Title:    "HN Post",
		URL:      "https://example.com/hn",
		Author:   "hn-user",
		Text:     "Text",
		Type:     "story",
		Score:    100,
		Comments: 25,
		PostedAt: &postedAt,
	})

	if document.Source != "hackernews" || document.ExternalID != "123" {
		t.Fatalf("unexpected normalized document: %+v", document)
	}
}
