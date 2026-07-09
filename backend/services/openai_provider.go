package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"signalstack-ai/backend/config"
)

type ChatCompletionProvider interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type OpenAIClient struct {
	APIKey         string
	EmbeddingModel string
	ChatModel      string
	BaseURL        string
	HTTPClient     *http.Client
}

func NewOpenAIClient(cfg config.OpenAIConfig) *OpenAIClient {
	return &OpenAIClient{
		APIKey:         cfg.APIKey,
		EmbeddingModel: cfg.EmbeddingModel,
		ChatModel:      cfg.ChatModel,
		BaseURL:        "https://api.openai.com/v1",
		HTTPClient:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *OpenAIClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if c == nil || c.APIKey == "" {
		return nil, fmt.Errorf("openai api key not configured")
	}
	model := c.EmbeddingModel
	if model == "" {
		model = "text-embedding-3-small"
	}
	reqBody := map[string]any{
		"model": model,
		"input": text,
	}
	var response embeddingResponse
	if err := c.doJSON(ctx, http.MethodPost, "/embeddings", reqBody, &response); err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("openai embedding response was empty")
	}
	return float64SliceToFloat32(response.Data[0].Embedding), nil
}

func (c *OpenAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	if c == nil || c.APIKey == "" {
		return "", fmt.Errorf("openai api key not configured")
	}
	model := c.ChatModel
	if model == "" {
		model = "gpt-4o-mini"
	}
	reqBody := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{{
			Role:    "user",
			Content: prompt,
		}},
		Temperature: 0.2,
	}
	var response chatCompletionResponse
	if err := c.doJSON(ctx, http.MethodPost, "/chat/completions", reqBody, &response); err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("openai chat completion response was empty")
	}
	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

func (c *OpenAIClient) doJSON(ctx context.Context, method, endpoint string, requestBody any, responseBody any) error {
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+c.APIKey)
	request.Header.Set("Content-Type", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("openai api request failed: %s", strings.TrimSpace(string(body)))
	}
	return json.Unmarshal(body, responseBody)
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func float64SliceToFloat32(values []float64) []float32 {
	result := make([]float32, len(values))
	for index, value := range values {
		result[index] = float32(value)
	}
	return result
}
