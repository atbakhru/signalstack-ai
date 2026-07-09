package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"signalstack-ai/backend/models"
)

type SpaceflightArticle struct {
	ID          string
	Title       string
	URL         string
	Summary     string
	Body        string
	Author      string
	PublishedAt *time.Time
	NewsSite    string
}

type SpaceflightAdapter struct{}

func (a SpaceflightAdapter) Name() string { return "spaceflight" }

func (a SpaceflightAdapter) Fetch(ctx context.Context, baseURL string, limit int) ([]models.NormalizedDocument, any, error) {
	if limit <= 0 {
		limit = 10
	}
	if baseURL == "" {
		baseURL = "https://api.spaceflightnewsapi.net/v4"
	}
	queryURL, err := url.Parse(baseURL + "/articles/")
	if err != nil {
		return nil, nil, err
	}
	q := queryURL.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("ordering", "-published_at")
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
		Results []struct {
			ID          int64    `json:"id"`
			Title       string   `json:"title"`
			URL         string   `json:"url"`
			Summary     string   `json:"summary"`
			NewsSite    string   `json:"news_site"`
			PublishedAt string   `json:"published_at"`
			Authors     []string `json:"authors"`
			ImageURL    string   `json:"image_url"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, nil, err
	}

	documents := make([]models.NormalizedDocument, 0, len(payload.Results))
	for _, article := range payload.Results {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)
		articleTime := publishedAt
		author := ""
		if len(article.Authors) > 0 {
			author = article.Authors[0]
		}
		documents = append(documents, models.NormalizedDocument{
			Source:      a.Name(),
			ExternalID:  fmt.Sprintf("%d", article.ID),
			Title:       article.Title,
			Summary:     article.Summary,
			URL:         article.URL,
			Author:      author,
			PublishedAt: &articleTime,
			Content:     article.Summary,
			Metadata: map[string]string{
				"news_site": article.NewsSite,
				"image_url": article.ImageURL,
			},
		})
	}
	return documents, payload, nil
}

func (a SpaceflightAdapter) Normalize(article SpaceflightArticle) models.NormalizedDocument {
	return models.NormalizedDocument{
		Source:      "spaceflight",
		ExternalID:  article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		URL:         article.URL,
		Author:      article.Author,
		PublishedAt: article.PublishedAt,
		Content:     article.Body,
		Metadata: map[string]string{
			"news_site": article.NewsSite,
		},
	}
}
