package connectors

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/dto"
)

// FetchLMSYSArena scrapes LMSYS Chatbot Arena leaderboard
// Note: This is a simple HTML scraper. For production, consider using their API if available
func FetchLMSYSArena(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	leaderboardURL := "https://chat.lmsys.org/"
	
	body, err := fetchWithLimits(ctx, leaderboardURL, cfg)
	if err != nil {
		return nil, err
	}

	html := string(body)
	
	// Simple pattern matching for model rankings
	// This is a basic implementation - in production, use proper HTML parsing
	modelPattern := regexp.MustCompile(`(?i)(gpt-4|claude|gemini|llama|mistral|palm)[-\w]*`)
	matches := modelPattern.FindAllString(html, -1)

	if len(matches) == 0 {
		// Return a placeholder item indicating the leaderboard exists
		return []dto.ContentItem{
			{
				Title:       "LMSYS Chatbot Arena Leaderboard Updated",
				URL:         leaderboardURL,
				PublishedAt: time.Now(),
				Domain:      "chat.lmsys.org",
				Category:    "benchmarks",
				Excerpt:     "Check the latest LLM rankings on the LMSYS Chatbot Arena leaderboard",
				Tags:        []string{"benchmarks", "llm", "leaderboard"},
			},
		}, nil
	}

	// Create items for top models found
	seen := make(map[string]bool)
	var items []dto.ContentItem
	
	for i, model := range matches {
		if i >= 10 { // Limit to top 10
			break
		}
		
		modelName := strings.ToLower(model)
		if seen[modelName] {
			continue
		}
		seen[modelName] = true

		rank := i + 1
		title := fmt.Sprintf("LMSYS Arena: %s (Rank ~%d)", model, rank)
		
		// Simulate engagement based on rank
		points := 100 - (rank * 5)
		
		item := dto.ContentItem{
			Title:       title,
			URL:         leaderboardURL,
			PublishedAt: time.Now(),
			Domain:      "chat.lmsys.org",
			Category:    "benchmarks",
			Points:      &points,
			Tags:        []string{"benchmarks", "llm", "arena"},
		}
		
		items = append(items, item)
	}

	return items, nil
}

// FetchOpenAIStatus fetches OpenAI API status and updates
func FetchOpenAIStatus(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	// OpenAI status RSS feed
	statusURL := "https://status.openai.com/history.rss"
	
	items, err := FetchRSSAtom(ctx, statusURL, cfg)
	if err != nil {
		return nil, err
	}

	// Tag all items as OpenAI status updates
	for i := range items {
		items[i].Category = "status"
		items[i].Tags = append(items[i].Tags, "openai", "status", "api")
		items[i].Domain = "status.openai.com"
	}

	return items, nil
}

// FetchAnthropicDocs fetches Anthropic API release notes
// Note: Anthropic doesn't have a public RSS feed, so we create a placeholder
func FetchAnthropicDocs(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	docsURL := "https://docs.anthropic.com/en/release-notes"
	
	// Try to fetch the page
	body, err := fetchWithLimits(ctx, docsURL, cfg)
	if err != nil {
		return nil, err
	}

	html := string(body)
	
	// Look for version numbers or dates in the content
	versionPattern := regexp.MustCompile(`(?i)(claude|version|v?\d+\.\d+)`)
	datePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	
	hasUpdates := versionPattern.MatchString(html) || datePattern.MatchString(html)
	
	if !hasUpdates {
		return nil, fmt.Errorf("no updates found")
	}

	// Extract dates if available
	dates := datePattern.FindAllString(html, -1)
	publishedAt := time.Now()
	if len(dates) > 0 {
		if t, err := time.Parse("2006-01-02", dates[0]); err == nil {
			publishedAt = t
		}
	}

	return []dto.ContentItem{
		{
			Title:       "Anthropic API Release Notes Updated",
			URL:         docsURL,
			PublishedAt: publishedAt,
			Domain:      "docs.anthropic.com",
			Category:    "docs",
			Excerpt:     "Check the latest Claude API updates and release notes",
			Tags:        []string{"anthropic", "claude", "api", "docs"},
		},
	}, nil
}

// FetchGitHubTrending fetches trending AI/ML repositories from GitHub
func FetchGitHubTrending(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	// GitHub trending page for AI/ML topics
	trendingURL := "https://github.com/trending?spoken_language_code=en"
	
	body, err := fetchWithLimits(ctx, trendingURL, cfg)
	if err != nil {
		return nil, err
	}

	html := string(body)
	
	// Simple pattern matching for repository names
	// Format: /owner/repo
	repoPattern := regexp.MustCompile(`href="/([\w-]+)/([\w-]+)"`)
	starsPattern := regexp.MustCompile(`(\d+(?:,\d+)*)\s+stars`)
	
	repoMatches := repoPattern.FindAllStringSubmatch(html, -1)
	starsMatches := starsPattern.FindAllStringSubmatch(html, -1)
	
	var items []dto.ContentItem
	seen := make(map[string]bool)
	
	for i, match := range repoMatches {
		if i >= 20 || len(match) < 3 { // Limit to top 20
			break
		}
		
		owner := match[1]
		repo := match[2]
		repoKey := owner + "/" + repo
		
		if seen[repoKey] {
			continue
		}
		seen[repoKey] = true
		
		repoURL := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
		title := fmt.Sprintf("⭐ Trending: %s", repoKey)
		
		// Extract stars if available
		var points *int
		if i < len(starsMatches) && len(starsMatches[i]) > 1 {
			starsStr := strings.ReplaceAll(starsMatches[i][1], ",", "")
			if stars, err := strconv.Atoi(starsStr); err == nil {
				points = &stars
			}
		}
		
		item := dto.ContentItem{
			Title:       title,
			URL:         repoURL,
			PublishedAt: time.Now(),
			Domain:      "github.com",
			Category:    "tools",
			Points:      points,
			Tags:        []string{"github", "trending", "opensource"},
		}
		
		items = append(items, item)
	}
	
	if len(items) == 0 {
		// Return placeholder
		return []dto.ContentItem{
			{
				Title:       "GitHub Trending Repositories",
				URL:         trendingURL,
				PublishedAt: time.Now(),
				Domain:      "github.com",
				Category:    "tools",
				Excerpt:     "Check trending AI/ML repositories on GitHub",
				Tags:        []string{"github", "trending"},
			},
		}, nil
	}

	return items, nil
}
