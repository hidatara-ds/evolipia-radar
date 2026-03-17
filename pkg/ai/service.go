package ai

import (
	"context"
	"fmt"
)

// Service is the primary business logic layer for AI operations.
// It wraps an underlying LLMProvider, allowing for centralized business rules
// (e.g., logging, default fallback, structured formatting) to be applied
// independently of the provider implementation.
type Service struct {
	provider LLMProvider
}

// NewService creates a new AI Service.
func NewService(provider LLMProvider) *Service {
	return &Service{
		provider: provider,
	}
}

// Chat handles a conversational exchange with the LLM.
func (s *Service) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("chat request must contain at least one message")
	}
	// Future: We can inject centralized behavior here, such as:
	// - System prompt injection for persona enforcing
	// - Pre-computation text safety scanning
	// - Caching mechanisms

	return s.provider.ChatCompletion(ctx, req)
}

// Summarize orchestrates a summarization task.
func (s *Service) Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("summarization text cannot be empty")
	}
	// Future: Centralized business rules:
	// - Text chunking if req.Text > max context window
	// - Fallback models if primary summary fails

	return s.provider.Summarize(ctx, req)
}

// Embed generates an embedding vector from text.
func (s *Service) Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if req.Input == "" {
		return nil, fmt.Errorf("embedding text cannot be empty")
	}
	return s.provider.Embed(ctx, req)
}
