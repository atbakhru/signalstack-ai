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

type GuardianArticle struct {
	ID          string
	Title       string
	URL         string
	Section     string
	Summary     string
	Body        string
	Author      string
	PublishedAt *time.Time
}

type GuardianAdapter struct{}

func (a GuardianAdapter) Name() string { return "guardian" }

func (a GuardianAdapter) Fetch(ctx context.Context, baseURL, apiKey string, limit int) ([]models.NormalizedDocument, any, error) {
	if limit <= 0 {
		limit = 10
	}
	if baseURL == "" {
		baseURL = "https://content.guardianapis.com"
	}
	queryURL, err := url.Parse(baseURL + "/search")
	if err != nil {
		return nil, nil, err
	}
	q := queryURL.Query()
	q.Set("page-size", fmt.Sprintf("%d", limit))
	q.Set("show-fields", "bodyText")
	q.Set("order-by", "newest")
	q.Set("api-key", apiKey)
	q.Set("q", "technology OR AI OR cloud OR cybersecurity")
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
		Response struct {
			Results []struct {
				ID                 string `json:"id"`
				WebTitle           string `json:"webTitle"`
				WebURL             string `json:"webUrl"`
				SectionName        string `json:"sectionName"`
				WebPublicationDate string `json:"webPublicationDate"`
				Fields             struct {
					BodyText string `json:"bodyText"`
				} `json:"fields"`
			} `json:"results"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, nil, err
	}

	documents := make([]models.NormalizedDocument, 0, len(payload.Response.Results))
	for _, article := range payload.Response.Results {
		publishedAt, _ := time.Parse(time.RFC3339, article.WebPublicationDate)
		articleTime := publishedAt
		documents = append(documents, models.NormalizedDocument{
			Source:      a.Name(),
			ExternalID:  article.ID,
			Title:       article.WebTitle,
			Summary:     article.Fields.BodyText,
			URL:         article.WebURL,
			Author:      "The Guardian",
			PublishedAt: &articleTime,
			Content:     article.Fields.BodyText,
			Metadata: map[string]string{
				"section": article.SectionName,
			},
		})
	}
	return documents, payload, nil
}

func (a GuardianAdapter) Normalize(article GuardianArticle) models.NormalizedDocument {
	return models.NormalizedDocument{
		Source:      "guardian",
		ExternalID:  article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		URL:         article.URL,
		Author:      article.Author,
		PublishedAt: article.PublishedAt,
		Content:     article.Body,
		Metadata: map[string]string{
			"section": article.Section,
		},
	}
}
