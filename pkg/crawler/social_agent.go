package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

type xSearchResponse struct {
	Data []struct {
		ID        string    `json:"id"`
		Text      string    `json:"text"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type SocialAgent struct {
	Platform string
	repo     *db.SettingRepository
}

func NewSocialAgent(platform string, pool *db.DB) *SocialAgent {
	return &SocialAgent{
		Platform: platform,
		repo:     db.NewSettingRepository(pool),
	}
}

func (a *SocialAgent) Name() string {
	return "SocialAgent-" + a.Platform
}

func (a *SocialAgent) Crawl(ctx context.Context, maxItems int) ([]Article, error) {
	// 1. Resolve API Key from Settings
	keyName := ""
	query := ""
	switch a.Platform {
	case "X":
		keyName = "x_api_key"
		query = "(AI OR ML OR Robotics OR LLM) -is:retweet lang:en"
	case "Threads":
		keyName = "threads_api_key"
		query = "AI Robotics Innovation"
	}

	if keyName == "" {
		return nil, nil
	}

	apiKey, err := a.repo.Get(ctx, keyName)
	if err != nil || apiKey == "" {
		log.Printf("[SOCIAL] Skipping %s agent: No API key found in settings.", a.Platform)
		return nil, nil
	}

	// 2. Fetch from X (Twitter) API V2
	if a.Platform == "X" {
		encodedQuery := url.QueryEscape(query)
		urlStr := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s&max_results=%d&tweet.fields=created_at",
			encodedQuery,
			10,
		)

		req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("X API request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("[SOCIAL] X API error: status %d. Check if Bearer Token is valid for V2 access.", resp.StatusCode)
			return nil, nil
		}

		var xResp xSearchResponse
		if err := json.NewDecoder(resp.Body).Decode(&xResp); err != nil {
			return nil, fmt.Errorf("failed to decode X response: %w", err)
		}

		var articles []Article
		for _, tweet := range xResp.Data {
			articles = append(articles, Article{
				Title:       "X Post by ID " + tweet.ID,
				Content:     tweet.Text,
				Link:        "https://x.com/x/status/" + tweet.ID,
				PublishedAt: tweet.CreatedAt,
				Source:      "X (Twitter)",
			})
		}

		log.Printf("[SOCIAL] Successfully fetched %d posts from X.", len(articles))
		return articles, nil
	}

	// Placeholder for Threads/Other platforms
	log.Printf("[SOCIAL] %s agent active but fetching logic is placeholder.", a.Platform)
	return []Article{}, nil
}
