package adapters

import (
	"testing"
	"time"
)

func TestArxivAdapterNormalize(t *testing.T) {
	adapter := ArxivAdapter{}
	publishedAt := time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC)

	document := adapter.Normalize(ArxivPaper{
		ID:          "arxiv-1",
		Title:       "Paper Title",
		URL:         "https://arxiv.org/abs/1",
		Abstract:    "Abstract",
		Authors:     []string{"Researcher"},
		Categories:  []string{"cs.AI"},
		PublishedAt: &publishedAt,
	})

	if document.Source != "arxiv" || document.Metadata["categories"] != "cs.AI" {
		t.Fatalf("unexpected normalized document: %+v", document)
	}
}
