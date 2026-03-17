package ai

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// TokenTracker monitors daily and monthly token usage in memory.
// For a multi-instance deployment, this should use Redis/PostgreSQL, 
// but for our constraint of "no complex infrastructure", in-memory with a mutex suffices.
type TokenTracker struct {
	mu sync.Mutex

	dailyTokens   int
	monthlyTokens int

	lastResetDaily   time.Time
	lastResetMonthly time.Time

	maxDailyTokens   int
	maxMonthlyTokens int
}

// NewTokenTracker initializes the tracker with strict budgets.
func NewTokenTracker(maxDaily, maxMonthly int) *TokenTracker {
	now := time.Now()
	return &TokenTracker{
		lastResetDaily:   now,
		lastResetMonthly: now,
		maxDailyTokens:   maxDaily,
		maxMonthlyTokens: maxMonthly,
	}
}

// checkResets safely handles daily and monthly rollover.
func (t *TokenTracker) checkResets(now time.Time) {
	// Reset daily if day changed
	if now.YearDay() != t.lastResetDaily.YearDay() || now.Year() != t.lastResetDaily.Year() {
		t.dailyTokens = 0
		t.lastResetDaily = now
	}

	// Reset monthly if month changed
	if now.Month() != t.lastResetMonthly.Month() || now.Year() != t.lastResetMonthly.Year() {
		t.monthlyTokens = 0
		t.lastResetMonthly = now
	}
}

// RecordUsage logs tokens consumed by a request and checks for warnings.
func (t *TokenTracker) RecordUsage(usage Usage) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.checkResets(time.Now())

	t.dailyTokens += usage.TotalTokens
	t.monthlyTokens += usage.TotalTokens

	// 80% Warning Triggers
	dailyWarningThresh := float64(t.maxDailyTokens) * 0.8
	monthlyWarningThresh := float64(t.maxMonthlyTokens) * 0.8

	if float64(t.dailyTokens) >= dailyWarningThresh && float64(t.dailyTokens-usage.TotalTokens) < dailyWarningThresh {
		log.Printf("[BUDGET WARNING] Daily token usage has reached 80%% (%d / %d tokens)", t.dailyTokens, t.maxDailyTokens)
	}

	if float64(t.monthlyTokens) >= monthlyWarningThresh && float64(t.monthlyTokens-usage.TotalTokens) < monthlyWarningThresh {
		log.Printf("[BUDGET WARNING] Monthly token usage has reached 80%% (%d / %d tokens)", t.monthlyTokens, t.maxMonthlyTokens)
	}
}

// IsBudgetExhausted returns true if we are over our strict free-tier limits.
func (t *TokenTracker) IsBudgetExhausted() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.checkResets(time.Now())

	if t.dailyTokens >= t.maxDailyTokens {
		return fmt.Errorf("daily token limit exceeded (%d / %d)", t.dailyTokens, t.maxDailyTokens)
	}
	if t.monthlyTokens >= t.maxMonthlyTokens {
		return fmt.Errorf("monthly token limit exceeded (%d / %d)", t.monthlyTokens, t.maxMonthlyTokens)
	}

	return nil
}

// Ensure interface compliance
var _ LLMProvider = (*TrackerMiddleware)(nil)

// TrackerMiddleware wrapping the LLM provider to intercept limits
type TrackerMiddleware struct {
	next    LLMProvider
	tracker *TokenTracker
}

func NewTrackerMiddleware(next LLMProvider, maxDaily, maxMonthly int) *TrackerMiddleware {
	t := NewTokenTracker(maxDaily, maxMonthly)
	return &TrackerMiddleware{next: next, tracker: t}
}

func (m *TrackerMiddleware) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if err := m.tracker.IsBudgetExhausted(); err != nil {
		return nil, fmt.Errorf("BUDGET EXHAUSTED: %w", err)
	}
	resp, err := m.next.ChatCompletion(ctx, req)
	if resp != nil {
		m.tracker.RecordUsage(resp.Usage)
	}
	return resp, err
}

func (m *TrackerMiddleware) Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error) {
	if err := m.tracker.IsBudgetExhausted(); err != nil {
		return nil, fmt.Errorf("BUDGET EXHAUSTED: %w", err)
	}
	resp, err := m.next.Summarize(ctx, req)
	if resp != nil {
		m.tracker.RecordUsage(resp.Usage)
	}
	return resp, err
}

func (m *TrackerMiddleware) Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if err := m.tracker.IsBudgetExhausted(); err != nil {
		return nil, fmt.Errorf("BUDGET EXHAUSTED: %w", err)
	}
	resp, err := m.next.Embed(ctx, req)
	if resp != nil {
		m.tracker.RecordUsage(resp.Usage)
	}
	return resp, err
}
