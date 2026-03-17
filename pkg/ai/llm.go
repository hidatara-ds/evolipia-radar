package ai

import (
	"context"
	"fmt"
)

// MessageRole defines the role of a chat participant.
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// ChatRequest contains the parameters for a chat completion request.
type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	Model       string        `json:"model,omitempty"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

// Validate checks the chat request for critical errors.
func (r *ChatRequest) Validate() error {
	if len(r.Messages) == 0 {
		return fmt.Errorf("messages array cannot be empty")
	}
	for i, msg := range r.Messages {
		if msg.Role == "" {
			return fmt.Errorf("message at index %d is missing a role", i)
		}
		if msg.Content == "" {
			return fmt.Errorf("message at index %d context cannot be empty", i)
		}
	}
	return nil
}

// ChatResponse contains the result of a chat completion.
type ChatResponse struct {
	Content string `json:"content"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
}

// Usage tracks token usage for billing and monitoring.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// SummarizeRequest contains the parameters for summarizing text.
type SummarizeRequest struct {
	Text        string   `json:"text"`
	Instruction string   `json:"instruction,omitempty"` // Custom instruction, e.g., "Provide a 3 bullet-point summary"
	Model       string   `json:"model,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
}

// Validate checks the summarize request for critical errors.
func (r *SummarizeRequest) Validate() error {
	if r.Text == "" {
		return fmt.Errorf("text to summarize cannot be empty")
	}
	if len(r.Text) > 100000 {
		return fmt.Errorf("text size exceeds maximum allowed limit")
	}
	return nil
}

// SummarizeResponse contains the resulting summary.
type SummarizeResponse struct {
	Summary string `json:"summary"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
}

// EmbeddingRequest contains parameters for generating vector embeddings.
type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model,omitempty"` // typically text-embedding-3-small
}

// Validate checks the embedding request for critical errors.
func (r *EmbeddingRequest) Validate() error {
	if r.Input == "" {
		return fmt.Errorf("input text to embed cannot be empty")
	}
	return nil
}

// EmbeddingResponse contains the generated embedding vector.
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
	Model     string    `json:"model"`
	Usage     Usage     `json:"usage"`
}

// LLMProvider is the core interface any AI provider (OpenRouter, OpenAI, Anthropic) must implement.
type LLMProvider interface {
	ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error)
	Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error)
}
