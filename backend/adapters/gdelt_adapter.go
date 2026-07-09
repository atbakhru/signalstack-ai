package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"signalstack-ai/backend/models"
)

type GDELTArticle struct {
	Title       string
	URL         string
	SourceName  string
	ArticleID   string
	Summary     string
	Author      string
	PublishedAt *time.Time
	Text        string
}

type GDELTAdapter struct{}

func (a GDELTAdapter) Name() string { return "gdelt" }

func (a GDELTAdapter) Fetch(ctx context.Context, baseURL string, limit int) ([]models.NormalizedDocument, any, error) {
	if limit <= 0 {
		limit = 10
	}
	if baseURL == "" {
		baseURL = "https://api.gdeltproject.org/api/v2/doc/doc"
	}
	queryURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	q := queryURL.Query()
	q.Set("query", "technology OR AI OR cloud OR cybersecurity")
	q.Set("mode", "ArtList")
	q.Set("format", "json")
	q.Set("maxrecords", fmt.Sprintf("%d", limit))
	queryURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		Articles []struct {
			Title    string `json:"title"`
			URL      string `json:"url"`
			Source   string `json:"sourceCountry"`
			Domain   string `json:"domain"`
			Seendate string `json:"seendate"`
			Snippet  string `json:"snippet"`
		} `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, nil, err
	}

	documents := make([]models.NormalizedDocument, 0, len(payload.Articles))
	for index, article := range payload.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, strings.ReplaceAll(article.Seendate, " ", "T"))
		articleTime := publishedAt
		documents = append(documents, models.NormalizedDocument{
			Source:      a.Name(),
			ExternalID:  fmt.Sprintf("gdelt-%d", index),
			Title:       article.Title,
			Summary:     article.Snippet,
			URL:         article.URL,
			Author:      article.Domain,
			PublishedAt: &articleTime,
			Content:     article.Snippet,
			Metadata: map[string]string{
				"source_country": article.Source,
				"domain":         article.Domain,
			},
		})
	}
	return documents, payload, nil
}

func (a GDELTAdapter) Normalize(article GDELTArticle) models.NormalizedDocument {
	return models.NormalizedDocument{
		Source:      "gdelt",
		ExternalID:  article.ArticleID,
		Title:       article.Title,
		Summary:     article.Summary,
		URL:         article.URL,
		Author:      article.Author,
		PublishedAt: article.PublishedAt,
		Content:     article.Text,
		Metadata: map[string]string{
			"source_name": article.SourceName,
		},
	}
}
