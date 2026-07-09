package adapters

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"signalstack-ai/backend/models"
)

type ArxivPaper struct {
	ID          string
	Title       string
	URL         string
	Abstract    string
	Authors     []string
	Categories  []string
	PublishedAt *time.Time
}

type ArxivAdapter struct{}

func (a ArxivAdapter) Name() string { return "arxiv" }

func (a ArxivAdapter) Fetch(ctx context.Context, baseURL string, limit int) ([]models.NormalizedDocument, any, error) {
	if limit <= 0 {
		limit = 10
	}
	if baseURL == "" {
		baseURL = "http://export.arxiv.org/api/query"
	}
	queryURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	q := queryURL.Query()
	q.Set("search_query", "all:cat OR all:machine learning OR all:large language models")
	q.Set("start", "0")
	q.Set("max_results", fmt.Sprintf("%d", limit))
	q.Set("sortBy", "submittedDate")
	q.Set("sortOrder", "descending")
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

	type feed struct {
		Entries []struct {
			ID        string `xml:"id"`
			Title     string `xml:"title"`
			Summary   string `xml:"summary"`
			Published string `xml:"published"`
			Authors   []struct {
				Name string `xml:"name"`
			} `xml:"author"`
			Links []struct {
				Href string `xml:"href,attr"`
				Rel  string `xml:"rel,attr"`
			} `xml:"link"`
			Categories []struct {
				Term string `xml:"term,attr"`
			} `xml:"category"`
		} `xml:"entry"`
	}
	var payload feed
	if err := xml.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, nil, err
	}

	documents := make([]models.NormalizedDocument, 0, len(payload.Entries))
	for _, entry := range payload.Entries {
		publishedAt, _ := time.Parse(time.RFC3339, entry.Published)
		articleTime := publishedAt
		author := ""
		if len(entry.Authors) > 0 {
			author = entry.Authors[0].Name
		}
		paperURL := ""
		for _, link := range entry.Links {
			if strings.EqualFold(link.Rel, "alternate") {
				paperURL = link.Href
				break
			}
		}
		category := ""
		if len(entry.Categories) > 0 {
			category = entry.Categories[0].Term
		}
		documents = append(documents, models.NormalizedDocument{
			Source:      a.Name(),
			ExternalID:  entry.ID,
			Title:       strings.TrimSpace(entry.Title),
			Summary:     strings.TrimSpace(entry.Summary),
			URL:         paperURL,
			Author:      author,
			PublishedAt: &articleTime,
			Content:     strings.TrimSpace(entry.Summary),
			Metadata: map[string]string{
				"category": category,
			},
		})
	}
	return documents, payload, nil
}

func (a ArxivAdapter) Normalize(paper ArxivPaper) models.NormalizedDocument {
	metadata := map[string]string{}
	if len(paper.Categories) > 0 {
		metadata["categories"] = paper.Categories[0]
	}

	author := ""
	if len(paper.Authors) > 0 {
		author = paper.Authors[0]
	}

	return models.NormalizedDocument{
		Source:      "arxiv",
		ExternalID:  paper.ID,
		Title:       paper.Title,
		Summary:     paper.Abstract,
		URL:         paper.URL,
		Author:      author,
		PublishedAt: paper.PublishedAt,
		Content:     paper.Abstract,
		Metadata:    metadata,
	}
}
