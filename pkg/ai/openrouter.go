package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultOpenRouterAPIURL       = "https://openrouter.ai/api/v1/chat/completions"
	defaultOpenRouterEmbeddingURL = "https://openrouter.ai/api/v1/embeddings"
	openRouterAppName             = "Evolipia Radar"
	openRouterAppURL              = "https://github.com/hidatara-ds/evolipia-radar"
)

// OpenRouterProviderConfig holds configuration for the OpenRouter client.
type OpenRouterProviderConfig struct {
	APIKey       string
	DefaultModel string
	Client       *http.Client
}

type openRouterProvider struct {
	apiKey       string
	defaultModel string
	client       *http.Client
}

// NewOpenRouterProvider creates a new LLMProvider implementation using OpenRouter.
func NewOpenRouterProvider(cfg OpenRouterProviderConfig) LLMProvider {
	client := cfg.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	defaultModel := cfg.DefaultModel
	if defaultModel == "" {
		defaultModel = "google/gemini-flash-1.5"
	}

	return &openRouterProvider{
		apiKey:       cfg.APIKey,
		defaultModel: defaultModel,
		client:       client,
	}
}

// openRouterRequest models the OpenRouter API request payload.
type openRouterRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

// openRouterResponse models the OpenRouter API response payload.
type openRouterResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// openRouterEmbeddingRequest models the API request for embeddings.
type openRouterEmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// openRouterEmbeddingResponse models the API response for embeddings.
type openRouterEmbeddingResponse struct {
	Model string `json:"model"`
	Data  []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// doRequest performs the actual HTTP call to OpenRouter.
func (p *openRouterProvider) doRequest(ctx context.Context, apiURL string, payload interface{}) (*http.Response, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("openrouter API key is not set")
	}

	reqBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Required Headers for OpenRouter
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", openRouterAppURL)
	req.Header.Set("X-Title", openRouterAppName)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// doChatRequest handles mapping and executing chat requests.
func (p *openRouterProvider) doChatRequest(ctx context.Context, payload openRouterRequest) (*openRouterResponse, error) {
	resp, err := p.doRequest(ctx, defaultOpenRouterAPIURL, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var orResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if orResp.Error != nil {
		return nil, fmt.Errorf("openrouter returned error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) == 0 {
		return nil, fmt.Errorf("openrouter returned empty choices")
	}

	return &orResp, nil
}

// ChatCompletion implements LLMProvider.
func (p *openRouterProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	model := req.Model
	if model == "" {
		model = p.defaultModel
	}

	payload := openRouterRequest{
		Model:       model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	orResp, err := p.doChatRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		Content: orResp.Choices[0].Message.Content,
		Model:   orResp.Model,
		Usage: Usage{
			PromptTokens:     orResp.Usage.PromptTokens,
			CompletionTokens: orResp.Usage.CompletionTokens,
			TotalTokens:      orResp.Usage.TotalTokens,
		},
	}, nil
}

// Summarize implements LLMProvider.
func (p *openRouterProvider) Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error) {
	instruction := req.Instruction
	if instruction == "" {
		instruction = "You are an AI/ML news analyst. Provide a concise, highly technical summary focused on engineering impact. Format cleanly."
	}

	systemMsg := ChatMessage{
		Role:    RoleSystem,
		Content: instruction,
	}
	userMsg := ChatMessage{
		Role:    RoleUser,
		Content: fmt.Sprintf("Please summarize the following text:\n\n%s", req.Text),
	}

	chatReq := ChatRequest{
		Messages:    []ChatMessage{systemMsg, userMsg},
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	chatResp, err := p.ChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, err
	}

	return &SummarizeResponse{
		Summary: chatResp.Content,
		Model:   chatResp.Model,
		Usage:   chatResp.Usage,
	}, nil
}

// Embed implements LLMProvider to generate dense vector embeddings.
func (p *openRouterProvider) Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	model := req.Model
	if model == "" {
		model = "text-embedding-3-small" // Fallback reasonable default
	}

	payload := openRouterEmbeddingRequest{
		Input: req.Input,
		Model: model,
	}

	resp, err := p.doRequest(ctx, defaultOpenRouterEmbeddingURL, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter embedding API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var orResp openRouterEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	if orResp.Error != nil {
		return nil, fmt.Errorf("openrouter returned error: %s", orResp.Error.Message)
	}

	if len(orResp.Data) == 0 {
		return nil, fmt.Errorf("openrouter returned empty embedding data")
	}

	return &EmbeddingResponse{
		Embedding: orResp.Data[0].Embedding,
		Model:     orResp.Model,
		Usage: Usage{
			PromptTokens: orResp.Usage.PromptTokens,
			TotalTokens:  orResp.Usage.TotalTokens,
		},
	}, nil
}
