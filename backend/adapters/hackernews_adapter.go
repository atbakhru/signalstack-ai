package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"signalstack-ai/backend/models"
)

type HackerNewsItem struct {
	ID       int64
	Title    string
	URL      string
	Author   string
	Text     string
	Type     string
	Score    int
	Comments int
	PostedAt *time.Time
}

type HackerNewsAdapter struct{}

func (a HackerNewsAdapter) Name() string { return "hackernews" }

func (a HackerNewsAdapter) Fetch(ctx context.Context, baseURL string, limit int) ([]models.NormalizedDocument, any, error) {
	if limit <= 0 {
		limit = 10
	}
	if baseURL == "" {
		baseURL = "https://hacker-news.firebaseio.com/v0"
	}
	storyLists := []string{"topstories", "askstories", "showstories"}
	documents := make([]models.NormalizedDocument, 0, limit*len(storyLists))
	rawPayload := make(map[string]any)

	for _, listName := range storyLists {
		idsURL := fmt.Sprintf("%s/%s.json", strings.TrimRight(baseURL, "/"), listName)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, idsURL, nil)
		if err != nil {
			return nil, nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		var ids []int64
		if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
			resp.Body.Close()
			return nil, nil, err
		}
		resp.Body.Close()
		rawPayload[listName] = ids

		for index, id := range ids {
			if index >= limit {
				break
			}
			itemURL := fmt.Sprintf("%s/item/%d.json", strings.TrimRight(baseURL, "/"), id)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, itemURL, nil)
			if err != nil {
				return nil, nil, err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, nil, err
			}
			var item struct {
				ID    int64   `json:"id"`
				Title string  `json:"title"`
				URL   string  `json:"url"`
				By    string  `json:"by"`
				Text  string  `json:"text"`
				Type  string  `json:"type"`
				Score int     `json:"score"`
				Kids  []int64 `json:"kids"`
				Time  int64   `json:"time"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				resp.Body.Close()
				return nil, nil, err
			}
			resp.Body.Close()
			postedAt := time.Unix(item.Time, 0).UTC()
			itemTime := postedAt
			documents = append(documents, models.NormalizedDocument{
				Source:      a.Name(),
				ExternalID:  strconv.FormatInt(item.ID, 10),
				Title:       item.Title,
				URL:         item.URL,
				Author:      item.By,
				PublishedAt: &itemTime,
				Content:     item.Text,
				Metadata: map[string]string{
					"type":     item.Type,
					"score":    strconv.Itoa(item.Score),
					"comments": strconv.Itoa(len(item.Kids)),
					"list":     listName,
				},
			})
			rawPayload[fmt.Sprintf("%s_%d", listName, item.ID)] = item
		}
	}

	return documents, rawPayload, nil
}

func (a HackerNewsAdapter) Normalize(item HackerNewsItem) models.NormalizedDocument {
	return models.NormalizedDocument{
		Source:      "hackernews",
		ExternalID:  strconv.FormatInt(item.ID, 10),
		Title:       item.Title,
		URL:         item.URL,
		Author:      item.Author,
		PublishedAt: item.PostedAt,
		Content:     item.Text,
		Metadata: map[string]string{
			"type":     item.Type,
			"score":    strconv.Itoa(item.Score),
			"comments": strconv.Itoa(item.Comments),
		},
	}
}
