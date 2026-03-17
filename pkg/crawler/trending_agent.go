package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TrendingAgent hits free APIs (like Algolia HN Search) for viral buzz.
type TrendingAgent struct {
	client *http.Client
}

func NewTrendingAgent() *TrendingAgent {
	return &TrendingAgent{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *TrendingAgent) Name() string {
	return "TrendingAgent"
}

// Algolia HN Search response
type hnResponse struct {
	Hits []struct {
		Title     string `json:"title"`
		URL       string `json:"url"`
		Points    int    `json:"points"`
		CreatedAt string `json:"created_at"`
		StoryText string `json:"story_text"` // Sometimes empty, sometimes body
	} `json:"hits"`
}

func (a *TrendingAgent) Crawl(ctx context.Context, maxItems int) ([]Article, error) {
	// Free API: fetch front page HN hits specifically tagged "story"
	url := "https://hn.algolia.com/api/v1/search?tags=front_page&hitsPerPage=" + fmt.Sprint(maxItems)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hn api returned status %d", resp.StatusCode)
	}

	var data hnResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var discovered []Article
	for _, hit := range data.Hits {
		if hit.Title == "" || hit.URL == "" {
			continue // skip empty or non-link stories (e.g. Ask HN mostly, though they are fine sometimes)
		}

		pubDate, _ := time.Parse(time.RFC3339, hit.CreatedAt)

		discovered = append(discovered, Article{
			Title:       hit.Title,
			Link:        hit.URL,
			Content:     hit.StoryText, // might be empty, we rely on title mostly
			PublishedAt: pubDate,
			Source:      a.Name(),
		})
	}

	return discovered, nil
}
