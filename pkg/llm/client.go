package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client handles LLM API interactions via OpenRouter
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Config for LLM client
type Config struct {
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	FallbackModels []string `json:"fallback_models"`
	APIKey         string   `json:"api_key"`
	MaxTokens      int      `json:"max_tokens"`
	Temperature    float64  `json:"temperature"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type   string      `json:"type"`
	Schema interface{} `json:"schema,omitempty"`
}

// CompletionRequest for OpenRouter API
type CompletionRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Temperature    float64         `json:"temperature,omitempty"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// CompletionResponse from OpenRouter API
type CompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewClient creates a new LLM client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://openrouter.ai/api/v1",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Complete sends a completion request to the LLM
func (c *Client) Complete(ctx context.Context, model string, messages []Message, maxTokens int, temperature float64, format *ResponseFormat) (string, error) {
	req := CompletionRequest{
		Model:          model,
		Messages:       messages,
		MaxTokens:      maxTokens,
		Temperature:    temperature,
		ResponseFormat: format,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/hidatara-ds/evolipia-radar")
	httpReq.Header.Set("X-Title", "Evolipia Radar")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var completionResp CompletionResponse
	if err := json.Unmarshal(respBody, &completionResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if completionResp.Error != nil {
		return "", fmt.Errorf("API error: %s", completionResp.Error.Message)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return completionResp.Choices[0].Message.Content, nil
}

type ArticleAnalysis struct {
	TLDR             string  `json:"tldr"`
	WhyItMatters     string  `json:"why_it_matters"`
	Novelty          float64 `json:"novelty"`
	Impact           float64 `json:"impact"`
	EngineeringValue float64 `json:"engineering_value"`
	Reasoning        string  `json:"reasoning"`
}

// AnalyzeArticle generates a structured analysis using LLM
func (c *Client) AnalyzeArticle(ctx context.Context, config Config, title, content string) (*ArticleAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze this AI/ML news article:

Title: %s
Content: %s

You must output a strict JSON object with the following fields:
- tldr: A 2-sentence summary of the article.
- why_it_matters: One sentence explaining why this matters to AI/ML engineers.
- reasoning: A 1-2 sentence explanation of your scores below.
- novelty: Score from 1-10. (10 = open weights SOTA paper, 1 = generic API wrapper).
- impact: Score from 1-10. (10 = industry-changing, 1 = unnoticeable).
- engineering_value: Score from 1-10. (10 = highly useful for builders/engineers, 1 = irrelevant).

Return ONLY the raw JSON object, no markdown blocks.`, title, content)

	messages := []Message{
		{Role: "system", Content: "You are an expert AI/ML news analyst and scoring engine. You strictly output valid JSON."},
		{Role: "user", Content: prompt},
	}

	response, err := c.Complete(ctx, config.Model, messages, config.MaxTokens, config.Temperature, &ResponseFormat{Type: "json_object"})
	if err != nil {
		return nil, err
	}

	var analysis ArticleAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis JSON: %w", err)
	}

	return &analysis, nil
}
