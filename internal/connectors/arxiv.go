package connectors

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/internal/normalizer"
)

const arxivAPIBase = "http://export.arxiv.org/api/query"

func FetchArxiv(ctx context.Context, query string, cfg *config.Config) ([]models.ContentItem, error) {
	// Build query URL
	params := url.Values{}
	params.Set("search_query", query)
	params.Set("start", "0")
	params.Set("max_results", "100")
	params.Set("sortBy", "submittedDate")
	params.Set("sortOrder", "descending")

	feedURL := arxivAPIBase + "?" + params.Encode()
	body, err := fetchWithLimits(ctx, feedURL, cfg)
	if err != nil {
		return nil, err
	}

	var feed ArxivFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse arXiv feed: %w", err)
	}

	var items []models.ContentItem
	for _, entry := range feed.Entries {
		item := models.ContentItem{
			Title:       entry.Title,
			URL:         entry.ID,
			PublishedAt: entry.Published,
			Domain:      normalizer.NormalizeDomain("arxiv.org"),
			Category:    "news",
			Excerpt:     entry.Summary,
			Tags:        []string{},
		}

		// Extract categories as tags
		for _, cat := range entry.Categories {
			item.Tags = append(item.Tags, cat.Term)
		}

		items = append(items, item)
	}

	return items, nil
}

type ArxivFeed struct {
	XMLName xml.Name      `xml:"feed"`
	Entries []ArxivEntry `xml:"entry"`
}

type ArxivEntry struct {
	ID        string         `xml:"id"`
	Title     string         `xml:"title"`
	Published time.Time      `xml:"published"`
	Summary   string         `xml:"summary"`
	Categories []ArxivCategory `xml:"category"`
}

type ArxivCategory struct {
	Term string `xml:"term,attr"`
}
