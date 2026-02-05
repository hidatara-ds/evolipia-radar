package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/internal/normalizer"
)

var (
	ErrTimeout   = errors.New("request timeout")
	ErrSizeLimit = errors.New("response size limit exceeded")
)

func fetchWithLimits(ctx context.Context, url string, cfg *config.Config) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "evolipia-radar/1.0")

	client := &http.Client{
		Timeout: cfg.FetchTimeout(),
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Limit response size
	limitedReader := io.LimitReader(resp.Body, cfg.MaxFetchBytes)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(body) >= int(cfg.MaxFetchBytes) {
		return nil, ErrSizeLimit
	}

	return body, nil
}

func FetchRSSAtom(ctx context.Context, feedURL string, cfg *config.Config) ([]models.ContentItem, error) {
	body, err := fetchWithLimits(ctx, feedURL, cfg)
	if err != nil {
		return nil, err
	}

	// Simple RSS/Atom parser (MVP - can be enhanced)
	items := parseRSSAtom(body)
	return items, nil
}

func parseRSSAtom(body []byte) []models.ContentItem {
	// MVP: Simple XML parsing
	// In production, use a proper RSS/Atom library
	content := string(body)
	var items []models.ContentItem

	// Very basic RSS parsing (item tags)
	itemStart := "<item>"
	itemEnd := "</item>"
	
	// Also handle Atom entries
	entryStart := "<entry>"
	entryEnd := "</entry>"

	parseItem := func(itemContent string) models.ContentItem {
		item := models.ContentItem{
			Category: "news",
			Tags:     []string{},
		}

		// Extract title
		if titleStart := strings.Index(itemContent, "<title>"); titleStart != -1 {
			titleStart += len("<title>")
			if titleEnd := strings.Index(itemContent[titleStart:], "</title>"); titleEnd != -1 {
				item.Title = strings.TrimSpace(itemContent[titleStart : titleStart+titleEnd])
			}
		}

		// Extract link
		if linkStart := strings.Index(itemContent, "<link>"); linkStart != -1 {
			linkStart += len("<link>")
			if linkEnd := strings.Index(itemContent[linkStart:], "</link>"); linkEnd != -1 {
				item.URL = strings.TrimSpace(itemContent[linkStart : linkStart+linkEnd])
			}
		} else if linkStart := strings.Index(itemContent, `href="`); linkStart != -1 {
			linkStart += len(`href="`)
			if linkEnd := strings.Index(itemContent[linkStart:], `"`); linkEnd != -1 {
				item.URL = strings.TrimSpace(itemContent[linkStart : linkStart+linkEnd])
			}
		}

		// Extract pubDate or published
		if pubStart := strings.Index(itemContent, "<pubDate>"); pubStart != -1 {
			pubStart += len("<pubDate>")
			if pubEnd := strings.Index(itemContent[pubStart:], "</pubDate>"); pubEnd != -1 {
				dateStr := strings.TrimSpace(itemContent[pubStart : pubStart+pubEnd])
				if t, err := time.Parse(time.RFC1123Z, dateStr); err == nil {
					item.PublishedAt = t
				} else if t, err := time.Parse(time.RFC1123, dateStr); err == nil {
					item.PublishedAt = t
				}
			}
		} else if pubStart := strings.Index(itemContent, "<published>"); pubStart != -1 {
			pubStart += len("<published>")
			if pubEnd := strings.Index(itemContent[pubStart:], "</published>"); pubEnd != -1 {
				dateStr := strings.TrimSpace(itemContent[pubStart : pubStart+pubEnd])
				if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
					item.PublishedAt = t
				}
			}
		}

		// Extract description/summary
		if descStart := strings.Index(itemContent, "<description>"); descStart != -1 {
			descStart += len("<description>")
			if descEnd := strings.Index(itemContent[descStart:], "</description>"); descEnd != -1 {
				item.Excerpt = strings.TrimSpace(itemContent[descStart : descStart+descEnd])
			}
		} else if sumStart := strings.Index(itemContent, "<summary>"); sumStart != -1 {
			sumStart += len("<summary>")
			if sumEnd := strings.Index(itemContent[sumStart:], "</summary>"); sumEnd != -1 {
				item.Excerpt = strings.TrimSpace(itemContent[sumStart : sumStart+sumEnd])
			}
		}

		if item.PublishedAt.IsZero() {
			item.PublishedAt = time.Now()
		}

		if item.URL != "" {
			parsedURL, err := url.Parse(item.URL)
			if err == nil {
				item.Domain = normalizer.NormalizeDomain(parsedURL.Hostname())
			}
		}

		return item
	}

	// Parse RSS items
	start := 0
	for {
		idx := strings.Index(content[start:], itemStart)
		if idx == -1 {
			break
		}
		itemStartIdx := start + idx
		endIdx := strings.Index(content[itemStartIdx:], itemEnd)
		if endIdx == -1 {
			break
		}
		itemContent := content[itemStartIdx : itemStartIdx+endIdx+len(itemEnd)]
		item := parseItem(itemContent)
		if item.Title != "" && item.URL != "" {
			items = append(items, item)
		}
		start = itemStartIdx + endIdx + len(itemEnd)
	}

	// Parse Atom entries if no RSS items found
	if len(items) == 0 {
		start = 0
		for {
			idx := strings.Index(content[start:], entryStart)
			if idx == -1 {
				break
			}
			entryStartIdx := start + idx
			endIdx := strings.Index(content[entryStartIdx:], entryEnd)
			if endIdx == -1 {
				break
			}
			entryContent := content[entryStartIdx : entryStartIdx+endIdx+len(entryEnd)]
			item := parseItem(entryContent)
			if item.Title != "" && item.URL != "" {
				items = append(items, item)
			}
			start = entryStartIdx + endIdx + len(entryEnd)
		}
	}

	return items
}

func FetchJSONAPI(ctx context.Context, apiURL string, mapping map[string]interface{}, cfg *config.Config) ([]models.ContentItem, error) {
	body, err := fetchWithLimits(ctx, apiURL, cfg)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	itemsPath := getString(mapping, "items_path", "items")
	itemsArray := getNestedValue(data, itemsPath)
	if itemsArray == nil {
		return nil, fmt.Errorf("items array not found at path: %s", itemsPath)
	}

	itemsSlice, ok := itemsArray.([]interface{})
	if !ok {
		return nil, fmt.Errorf("items_path does not point to an array")
	}

	titlePath := getString(mapping, "title_path", "title")
	urlPath := getString(mapping, "url_path", "url")
	publishedAtPath := getString(mapping, "published_at_path", "published_at")
	summaryPath := getString(mapping, "summary_path", "")

	var items []models.ContentItem
	for _, itemRaw := range itemsSlice {
		itemMap, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}

		item := models.ContentItem{
			Category: "news",
			Tags:     []string{},
		}

		item.Title = getStringValue(itemMap, titlePath)
		item.URL = getStringValue(itemMap, urlPath)
		
		if item.Title == "" || item.URL == "" {
			continue
		}

		// Parse published_at
		dateStr := getStringValue(itemMap, publishedAtPath)
		if dateStr != "" {
			if t, err := parseDate(dateStr); err == nil {
				item.PublishedAt = t
			}
		}
		if item.PublishedAt.IsZero() {
			item.PublishedAt = time.Now()
		}

		if summaryPath != "" {
			item.Excerpt = getStringValue(itemMap, summaryPath)
		}

		parsedURL, err := url.Parse(item.URL)
		if err == nil {
			item.Domain = normalizer.NormalizeDomain(parsedURL.Hostname())
		}

		items = append(items, item)
	}

	return items, nil
}

func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := interface{}(data)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			return nil
		}
	}
	return current
}

func getStringValue(m map[string]interface{}, path string) string {
	val := getNestedValue(m, path)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
