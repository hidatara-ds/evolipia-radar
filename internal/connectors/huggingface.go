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

const hfAPIBase = "https://huggingface.co/api"

// FetchHuggingFaceTrending fetches trending models from HuggingFace
func FetchHuggingFaceTrending(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	apiURL := hfAPIBase + "/models?sort=trending&limit=50"

	body, err := fetchWithLimits(ctx, apiURL, cfg)
	if err != nil {
		return nil, err
	}

	var models []struct {
		ID           string    `json:"id"`
		ModelID      string    `json:"modelId"`
		Author       string    `json:"author"`
		Downloads    int       `json:"downloads"`
		Likes        int       `json:"likes"`
		Tags         []string  `json:"tags"`
		LastModified time.Time `json:"lastModified"`
	}

	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("failed to parse HuggingFace response: %w", err)
	}

	var items []dto.ContentItem
	for _, model := range models {
		modelID := model.ModelID
		if modelID == "" {
			modelID = model.ID
		}

		modelURL := fmt.Sprintf("https://huggingface.co/%s", modelID)
		title := fmt.Sprintf("🤗 Trending Model: %s", modelID)

		points := model.Likes
		comments := model.Downloads / 100 // Approximate engagement

		item := dto.ContentItem{
			Title:       title,
			URL:         modelURL,
			PublishedAt: model.LastModified,
			Domain:      "huggingface.co",
			Category:    "models",
			Points:      &points,
			Comments:    &comments,
			Tags:        model.Tags,
		}

		items = append(items, item)
	}

	return items, nil
}

// FetchPapersWithCode fetches trending papers from Papers with Code
func FetchPapersWithCode(ctx context.Context, cfg *config.Config) ([]dto.ContentItem, error) {
	apiURL := "https://paperswithcode.com/api/v1/papers/"

	body, err := fetchWithLimits(ctx, apiURL, cfg)
	if err != nil {
		return nil, err
	}

	var response struct {
		Results []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Abstract    string `json:"abstract"`
			URL         string `json:"url_abs"`
			PaperURL    string `json:"paper_url"`
			PublishedAt string `json:"published"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Papers with Code response: %w", err)
	}

	var items []dto.ContentItem
	for _, paper := range response.Results {
		paperURL := paper.PaperURL
		if paperURL == "" {
			paperURL = paper.URL
		}

		publishedAt := time.Now()
		if paper.PublishedAt != "" {
			if t, err := time.Parse(time.RFC3339, paper.PublishedAt); err == nil {
				publishedAt = t
			}
		}

		parsedURL, _ := url.Parse(paperURL)
		domain := "paperswithcode.com"
		if parsedURL != nil && parsedURL.Hostname() != "" {
			domain = normalizer.NormalizeDomain(parsedURL.Hostname())
		}

		item := dto.ContentItem{
			Title:       paper.Title,
			URL:         paperURL,
			PublishedAt: publishedAt,
			Excerpt:     paper.Abstract,
			Domain:      domain,
			Category:    "research",
			Tags:        []string{"papers", "research"},
		}

		items = append(items, item)
	}

	return items, nil
}
