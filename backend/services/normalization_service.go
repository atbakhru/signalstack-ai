package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"signalstack-ai/backend/models"
	"signalstack-ai/backend/utils"
)

type NormalizationService struct{}

func (s NormalizationService) NormalizeDocument(sourceName, externalID, title, summary, url, author string, publishedAt *time.Time, content string, metadata map[string]string) models.NormalizedDocument {
	return models.NormalizedDocument{
		Source:      strings.TrimSpace(sourceName),
		ExternalID:  strings.TrimSpace(externalID),
		Title:       utils.CleanText(title),
		Summary:     utils.CleanText(summary),
		URL:         strings.TrimSpace(url),
		Author:      utils.CleanText(author),
		PublishedAt: publishedAt,
		Content:     utils.CleanText(content),
		Metadata:    metadata,
	}
}

func (s NormalizationService) BuildContent(title, summary, body string, extras ...string) string {
	parts := []string{title, summary, body}
	parts = append(parts, extras...)
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		cleaned := utils.CleanText(part)
		if cleaned != "" {
			normalized = append(normalized, cleaned)
		}
	}
	return strings.Join(normalized, "\n\n")
}

func (s NormalizationService) BuildContentHash(document models.NormalizedDocument) string {
	var keys []string
	for key := range document.Metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	builder := strings.Builder{}
	builder.WriteString(document.Source)
	builder.WriteString("|")
	builder.WriteString(document.ExternalID)
	builder.WriteString("|")
	builder.WriteString(document.Title)
	builder.WriteString("|")
	builder.WriteString(document.Summary)
	builder.WriteString("|")
	builder.WriteString(document.URL)
	builder.WriteString("|")
	builder.WriteString(document.Author)
	builder.WriteString("|")
	if document.PublishedAt != nil {
		builder.WriteString(document.PublishedAt.UTC().Format(time.RFC3339Nano))
	}
	builder.WriteString("|")
	builder.WriteString(document.Content)

	for _, key := range keys {
		builder.WriteString("|")
		builder.WriteString(fmt.Sprintf("%s=%s", key, document.Metadata[key]))
	}

	return utils.HashString(builder.String())
}
