package services

import "testing"

func TestChunkingServiceChunkText(t *testing.T) {
	service := NewChunkingService(5, 2)
	chunks := service.ChunkText("one two three four five six seven eight")

	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}

	if chunks[0].Content != "one two three four five" {
		t.Fatalf("unexpected first chunk content: %q", chunks[0].Content)
	}

	if chunks[1].Content != "four five six seven eight" {
		t.Fatalf("unexpected second chunk content: %q", chunks[1].Content)
	}
}
