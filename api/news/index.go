package news

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/pkg/api"
	_ "github.com/lib/pq"
)

type NewsItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	PublishedAt  time.Time `json:"published_at"`
	Category     string    `json:"category"`
	Score        float64   `json:"score"`      // Always 1-10 scale for frontend
	RawScore     float64   `json:"raw_score"`  // Internal 0.0-1.0 for debugging
	HeatLevel    string    `json:"heat_level"` // "hot", "rising", "signal", "low"
	TLDR         string    `json:"tldr,omitempty"`
	WhyItMatters string    `json:"why_it_matters,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func convertToScale10(score float64) float64 {
	if score <= 0 {
		return 1.0
	}
	if score >= 1.0 {
		return 10.0
	}
	scaled := (score * 9.0) + 1.0
	return math.Round(scaled*10) / 10
}

func getHeatLevel(score10 float64) string {
	switch {
	case score10 >= 7.0:
		return "hot"
	case score10 >= 5.0:
		return "rising"
	case score10 >= 3.0:
		return "signal"
	default:
		return "low"
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query()
	topics := query["topic"]
	if len(topics) == 0 && query.Get("topic") != "" {
		topics = strings.Split(query.Get("topic"), ",")
	}
	domains := query["domain"]
	if len(domains) == 0 && query.Get("domain") != "" {
		domains = strings.Split(query.Get("domain"), ",")
	}
	sortMode := query.Get("sort")
	timeRange := query.Get("time") // "24h", "7d", "30d"
	searchQuery := strings.ToLower(query.Get("q"))

	timeThreshold := time.Now().Add(-7 * 24 * time.Hour) // default 7d
	if timeRange == "24h" {
		timeThreshold = time.Now().Add(-24 * time.Hour)
	} else if timeRange == "30d" {
		timeThreshold = time.Now().Add(-30 * 24 * time.Hour)
	} else if timeRange == "all" {
		timeThreshold = time.Time{}
	}

	useFallback := false
	dbURL := os.Getenv("DATABASE_URL")
	var db *sql.DB
	var err error

	if dbURL == "" {
		log.Println("⚠️ DATABASE_URL not set, using JSON fallback")
		useFallback = true
	} else {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("⚠️ Failed to connect to database: %v. Using JSON fallback", err)
			useFallback = true
		} else {
			defer db.Close()
			db.SetMaxOpenConns(3)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := db.PingContext(ctx); err != nil {
				log.Printf("⚠️ Failed to ping database: %v. Using JSON fallback", err)
				useFallback = true
			}
		}
	}

	if useFallback {
		handleJSONFallback(w, topics, domains, sortMode, timeThreshold, searchQuery)
		return
	}

	sqlQuery := `
		SELECT 
			i.id, i.title, i.url, i.domain, i.published_at, i.category,
			COALESCE(s.final, 0) as raw_score,
			COALESCE(sm.tldr, '') as tldr,
			COALESCE(sm.why_it_matters, '') as why_it_matters,
			COALESCE(sm.tags, '[]'::jsonb) as tags
		FROM items i
		LEFT JOIN scores s ON i.id = s.item_id
		LEFT JOIN summaries sm ON i.id = sm.item_id
		WHERE i.published_at >= $1
	`
	args := []interface{}{timeThreshold}
	argIdx := 2

	if len(domains) > 0 {
		var placeholders []string
		for _, d := range domains {
			placeholders = append(placeholders, "$"+strconv.Itoa(argIdx))
			args = append(args, d)
			argIdx++
		}
		sqlQuery += ` AND i.domain IN (` + strings.Join(placeholders, ",") + `)`
	}

	if searchQuery != "" {
		sqlQuery += ` AND (LOWER(i.title) LIKE $` + strconv.Itoa(argIdx) + ` OR LOWER(COALESCE(sm.tldr, '')) LIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+searchQuery+"%")
		argIdx++
	}

	if len(topics) > 0 {
		var topicChecks []string
		for _, t := range topics {
			if t != "" && strings.ToLower(t) != "all" {
				topicChecks = append(topicChecks, `sm.tags @> $`+strconv.Itoa(argIdx)+`::jsonb`)
				args = append(args, `["`+t+`"]`)
				argIdx++
			}
		}
		if len(topicChecks) > 0 {
			sqlQuery += ` AND (` + strings.Join(topicChecks, " OR ") + `)`
		}
	}

	switch sortMode {
	case "newest":
		sqlQuery += ` ORDER BY i.published_at DESC`
	case "oldest":
		sqlQuery += ` ORDER BY i.published_at ASC`
	default:
		sqlQuery += ` ORDER BY (
			COALESCE(s.final, 0) * 0.7 + 
			CASE WHEN i.published_at > NOW() - INTERVAL '24 hours' THEN 0.3 
			     WHEN i.published_at > NOW() - INTERVAL '48 hours' THEN 0.2 
			     ELSE 0.1 END
		) DESC, i.published_at DESC`
	}
	sqlQuery += ` LIMIT 30`

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("❌ Failed to query database: %v. Using JSON fallback", err)
		handleJSONFallback(w, topics, domains, sortMode, timeThreshold, searchQuery)
		return
	}
	defer rows.Close()

	var items []NewsItem
	for rows.Next() {
		var item NewsItem
		var rawScore float64
		var tagsJSON []byte

		if err := rows.Scan(&item.ID, &item.Title, &item.URL, &item.Domain, &item.PublishedAt, &item.Category, &rawScore, &item.TLDR, &item.WhyItMatters, &tagsJSON); err != nil {
			continue
		}

		item.RawScore = rawScore
		item.Score = convertToScale10(rawScore)
		item.HeatLevel = getHeatLevel(item.Score)

		if len(tagsJSON) > 0 && string(tagsJSON) != "[]" {
			_ = json.Unmarshal(tagsJSON, &item.Tags)
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        items,
			"total_count":  len(items),
			"last_updated": time.Now(),
		},
	})
}

func handleJSONFallback(w http.ResponseWriter, topics, domains []string, sortMode string, timeThreshold time.Time, searchQuery string) {
	data, err := api.LoadNewsData()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Error: "Failed to load fallback JSON data"})
		return
	}

	var filtered []NewsItem
	for _, rawItem := range data.Items {
		if rawItem.PublishedAt.Before(timeThreshold) {
			continue
		}

		if len(domains) > 0 {
			matchDomain := false
			for _, d := range domains {
				if strings.EqualFold(rawItem.Domain, d) {
					matchDomain = true
					break
				}
			}
			if !matchDomain {
				continue
			}
		}

		if searchQuery != "" {
			if !strings.Contains(strings.ToLower(rawItem.Title), searchQuery) &&
				!strings.Contains(strings.ToLower(rawItem.TLDR), searchQuery) {
				continue
			}
		}

		if len(topics) > 0 {
			matchTopic := false
			for _, reqTopic := range topics {
				if strings.ToLower(reqTopic) == "all" || reqTopic == "" {
					matchTopic = true
					break
				}
				for _, tag := range rawItem.Tags {
					if strings.EqualFold(tag, reqTopic) {
						matchTopic = true
						break
					}
				}
				if matchTopic {
					break
				}
			}
			if !matchTopic {
				continue
			}
		}

		item := NewsItem{
			ID:           rawItem.ID,
			Title:        rawItem.Title,
			URL:          rawItem.URL,
			Domain:       rawItem.Domain,
			PublishedAt:  rawItem.PublishedAt,
			Category:     rawItem.Category,
			Tags:         rawItem.Tags,
			TLDR:         rawItem.TLDR,
			WhyItMatters: rawItem.WhyItMatters,
			RawScore:     rawItem.Score, 
		}
		item.Score = convertToScale10(rawItem.Score)
		item.HeatLevel = getHeatLevel(item.Score)
		filtered = append(filtered, item)
	}

	if sortMode == "newest" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].PublishedAt.After(filtered[j].PublishedAt) })
	} else if sortMode == "oldest" {
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].PublishedAt.Before(filtered[j].PublishedAt) })
	} else {
		sort.Slice(filtered, func(i, j int) bool {
			scoreI := filtered[i].RawScore * 0.7
			if filtered[i].PublishedAt.After(time.Now().Add(-24 * time.Hour)) {
				scoreI += 0.3
			}
			scoreJ := filtered[j].RawScore * 0.7
			if filtered[j].PublishedAt.After(time.Now().Add(-24 * time.Hour)) {
				scoreJ += 0.3
			}
			return scoreI > scoreJ
		})
	}

	if len(filtered) > 30 {
		filtered = filtered[:30]
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        filtered,
			"total_count":  len(filtered),
			"last_updated": data.LastUpdated,
		},
	})
}
