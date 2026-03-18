package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EmbeddingProvider defines the specific interface for generating vector maps in Phase 5.
type EmbeddingProvider interface {
	Embed(text string) ([]float64, error)
}

// OpenRouterEmbedder implements EmbeddingProvider using the OpenRouter embeddings API.
type OpenRouterEmbedder struct {
	apiKey string
	client *http.Client
}

func NewOpenRouterEmbedder(apiKey string) *OpenRouterEmbedder {
	return &OpenRouterEmbedder{
		apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (e *OpenRouterEmbedder) Embed(text string) ([]float64, error) {
	if e.apiKey == "" {
		return nil, fmt.Errorf("openrouter API key is not set")
	}

	payload := map[string]interface{}{
		"input": text,
		"model": "text-embedding-3-small", // fallback or default model
	}

	reqBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://openrouter.ai/api/v1/embeddings", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/hidatara-ds/evolipia-radar")
	req.Header.Set("X-Title", "Evolipia Radar")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var orResp struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(orResp.Data) == 0 {
		return nil, fmt.Errorf("returned empty embedding data")
	}

	return orResp.Data[0].Embedding, nil
}
