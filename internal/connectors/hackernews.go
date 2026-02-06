package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/dto"
	"github.com/hidatara-ds/evolipia-radar/internal/normalizer"
)

const hnAPIBase = "https://hacker-news.firebaseio.com/v0"

func FetchHackerNews(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	// Fetch top stories
	topStoriesURL := hnAPIBase + "/topstories.json"
	body, err := fetchWithLimits(ctx, topStoriesURL, cfg)
	if err != nil {
		return nil, err
	}

	var storyIDs []int
	if err := json.Unmarshal(body, &storyIDs); err != nil {
		return nil, fmt.Errorf("failed to parse HN top stories: %w", err)
	}

	// Fetch top 100 stories
	limit := 100
	if limit > len(storyIDs) {
		limit = len(storyIDs)
	}

	var items []dto.ContentItem
	for i := 0; i < limit; i++ {
		storyID := storyIDs[i]
		item, err := fetchHNItem(ctx, storyID, cfg)
		if err != nil {
			continue // Skip failed items
		}
		if item != nil {
			items = append(items, *item)
		}
	}

	return items, nil
}

func fetchHNItem(ctx context.Context, id int, cfg *config.Config) (*dto.ContentItem, error) {
	itemURL := fmt.Sprintf("%s/item/%d.json", hnAPIBase, id)
	body, err := fetchWithLimits(ctx, itemURL, cfg)
	if err != nil {
		return nil, err
	}

	var hnItem struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		Score       int    `json:"score"`
		Descendants int    `json:"descendants"`
		Time        int64  `json:"time"`
		Type        string `json:"type"`
	}

	if err := json.Unmarshal(body, &hnItem); err != nil {
		return nil, err
	}

	if hnItem.Type != "story" || hnItem.Title == "" {
		return nil, nil
	}

	// HN items without URL are "Ask HN" posts - use discussion URL
	if hnItem.URL == "" {
		hnItem.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", hnItem.ID)
	}

	parsedURL, err := url.Parse(hnItem.URL)
	if err != nil {
		return nil, err
	}

	item := &dto.ContentItem{
		Title:       hnItem.Title,
		URL:         hnItem.URL,
		PublishedAt: time.Unix(hnItem.Time, 0),
		Domain:      normalizer.NormalizeDomain(parsedURL.Hostname()),
		Category:    "news",
		Points:      &hnItem.Score,
		Comments:    &hnItem.Descendants,
		RankPos:     &id, // Use story ID as rank position indicator
		Tags:        []string{},
	}

	return item, nil
}
