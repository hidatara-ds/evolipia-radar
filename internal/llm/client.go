package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// CompletionRequest for OpenRouter API
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
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
func (c *Client) Complete(ctx context.Context, model string, messages []Message, maxTokens int, temperature float64) (string, error) {
	req := CompletionRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
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
	defer resp.Body.Close()

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

// Summarize generates an abstractive summary using LLM
func (c *Client) Summarize(ctx context.Context, config Config, title, content string) (string, string, error) {
	prompt := fmt.Sprintf(`Summarize this AI/ML news article:

Title: %s
Content: %s

Provide:
1. A 2-sentence summary (TLDR)
2. One sentence explaining why this matters to AI/ML engineers

Format your response as:
TLDR: [your summary]
WHY: [why it matters]`, title, content)

	messages := []Message{
		{Role: "system", Content: "You are an AI/ML news analyst. Provide concise, technical summaries focused on engineering impact."},
		{Role: "user", Content: prompt},
	}

	response, err := c.Complete(ctx, config.Model, messages, config.MaxTokens, config.Temperature)
	if err != nil {
		return "", "", err
	}

	// Parse response
	tldr, why := parseSummaryResponse(response)
	return tldr, why, nil
}

func parseSummaryResponse(response string) (string, string) {
	lines := bytes.Split([]byte(response), []byte("\n"))
	tldr := ""
	why := ""

	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("TLDR:")) {
			tldr = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("TLDR:"))))
		} else if bytes.HasPrefix(line, []byte("WHY:")) {
			why = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("WHY:"))))
		}
	}

	// Fallback if parsing fails
	if tldr == "" {
		tldr = response
	}

	return tldr, why
}
