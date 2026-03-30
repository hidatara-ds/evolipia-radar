package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RedditAgent struct {
	Subreddits []string
}

func NewRedditAgent() *RedditAgent {
	return &RedditAgent{
		Subreddits: []string{"technology", "artificialintelligence", "programming"},
	}
}

func (a *RedditAgent) Name() string {
	return "RedditAgent"
}

func (a *RedditAgent) Crawl(ctx context.Context, maxItems int) ([]Article, error) {
	var results []Article

	for _, sub := range a.Subreddits {
		url := fmt.Sprintf("https://www.reddit.com/r/%s/new.json?limit=25", sub)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "EvolipiaRadar/1.0")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var data struct {
			Data struct {
				Children []struct {
					Data struct {
						Title     string  `json:"title"`
						URL       string  `json:"url"`
						Selftext  string  `json:"selftext"`
						Created   float64 `json:"created_utc"`
						ID        string  `json:"id"`
						Permalink string  `json:"permalink"`
					} `json:"data"`
				} `json:"children"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			continue
		}

		for _, child := range data.Data.Children {
			d := child.Data
			results = append(results, Article{
				Title:       d.Title,
				Link:        d.URL,
				Content:     d.Selftext,
				PublishedAt: time.Unix(int64(d.Created), 0),
				Source:      fmt.Sprintf("reddit/r/%s", sub),
			})

			if len(results) >= maxItems {
				break
			}
		}

		if len(results) >= maxItems {
			break
		}
	}

	return results, nil
}
