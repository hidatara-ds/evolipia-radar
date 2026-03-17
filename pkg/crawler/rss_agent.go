package crawler

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// RSSAgent crawls high-signal standard RSS feeds. (Zero Cost).
type RSSAgent struct {
	client *http.Client
	feeds  []string
}

func NewRSSAgent() *RSSAgent {
	return &RSSAgent{
		client: &http.Client{Timeout: 10 * time.Second},
		feeds: []string{
			"https://news.ycombinator.com/rss",          // HackerNews
			"https://techcrunch.com/category/artificial-intelligence/feed/", // TechCrunch AI
			// Add more high-signal feeds here
		},
	}
}

func (a *RSSAgent) Name() string {
	return "RSSAgent"
}

// Simple RSS struct parser
type rss struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (a *RSSAgent) Crawl(ctx context.Context, maxItems int) ([]Article, error) {
	var discovered []Article

	for _, feedURL := range a.feeds {
		if len(discovered) >= maxItems {
			break
		}

		req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
		if err != nil {
			continue
		}

		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var feed rss
		if err := xml.Unmarshal(body, &feed); err != nil {
			continue // Skip malformed feeds
		}

		for _, item := range feed.Channel.Items {
			if len(discovered) >= maxItems {
				break
			}
			
			// Try parsing pubDate or fall back to now
			pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				pubDate = time.Now()
			}

			discovered = append(discovered, Article{
				Title:       item.Title,
				Link:        item.Link,
				Content:     item.Description, // description might be summary or full HTML
				PublishedAt: pubDate,
				Source:      a.Name(),
			})
		}
	}

	return discovered, nil
}
