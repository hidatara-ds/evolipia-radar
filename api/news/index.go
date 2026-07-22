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

const allConst = "all"

var TopicAliases = map[string][]string{
	"llm":         {"llm", "large language models", "gpt", "claude", "llama", "general_ai"},
	"agents":      {"agents", "autonomous", "autogpt", "babyagi"},
	"vision":      {"computer vision", "cv", "image generation", "midjourney", "dalle"},
	"open-source": {"open source", "oss", "huggingface", "open_source"},
	"infra":       {"infrastructure", "mlops", "deployment", "kubernetes"},
	"robotics":    {"robotics", "embodied ai", "control"},
	"security":    {"security", "safety", "alignment", "jailbreak"},
}

func getAliases(topic string) []string {
	topic = strings.ToLower(topic)
	aliases := []string{topic}
	if mapped, ok := TopicAliases[topic]; ok {
		aliases = append(aliases, mapped...)
	}
	return aliases
}

type NewsItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	PublishedAt  time.Time `json:"published_at"`
	Category     string    `json:"category"`
	Score        float64   `json:"score"`
	RawScore     float64   `json:"raw_score"`
	HeatLevel    string    `json:"heat_level"`
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

func parseListQuery(val []string, rawVal string) []string {
	if len(val) == 0 && rawVal != "" {
		return strings.Split(rawVal, ",")
	}
	return val
}

func getTimeThreshold(timeRange string) time.Time {
	switch timeRange {
	case "24h":
		return time.Now().Add(-24 * time.Hour)
	case "30d":
		return time.Now().Add(-30 * 24 * time.Hour)
	case allConst:
		return time.Time{}
	default:
		return time.Now().Add(-7 * 24 * time.Hour)
	}
}

func connectToDB() (*sql.DB, bool) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("⚠️ DATABASE_URL not set, using JSON fallback")
		return nil, true
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("⚠️ Failed to connect to database: %v. Using JSON fallback", err)
		return nil, true
	}

	db.SetMaxOpenConns(3)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("⚠️ Failed to ping database: %v. Using JSON fallback", err)
		_ = db.Close()
		return nil, true
	}

	return db, false
}

func buildSQLQuery(topics, domains []string, sortMode string, timeThreshold time.Time, searchQuery string) (string, []interface{}) {
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
			if t != "" && !strings.EqualFold(t, allConst) {
				aliases := getAliases(t)
				var placeholders []string
				for _, alias := range aliases {
					placeholders = append(placeholders, "$"+strconv.Itoa(argIdx))
					args = append(args, strings.ToLower(alias))
					argIdx++
				}
				topicChecks = append(topicChecks, `EXISTS (SELECT 1 FROM jsonb_array_elements_text(sm.tags) as tag WHERE LOWER(tag) IN (`+strings.Join(placeholders, ",")+`))`)
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
	return sqlQuery, args
}

func scanNewsItems(rows *sql.Rows) []NewsItem {
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
			if err := json.Unmarshal(tagsJSON, &item.Tags); err != nil {
				log.Printf("⚠️ Failed to unmarshal tags for item %s: %v", item.ID, err)
			}
		}
		items = append(items, item)
	}
	return items
}

func Handler(w http.ResponseWriter, r *http.Request) {
	api.EnableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query()
	topics := parseListQuery(query["topic"], query.Get("topic"))
	domains := parseListQuery(query["domain"], query.Get("domain"))
	sortMode := query.Get("sort")
	searchQuery := strings.ToLower(query.Get("q"))
	timeThreshold := getTimeThreshold(query.Get("time"))

	db, useFallback := connectToDB()
	if db != nil {
		defer db.Close()
	}

	if useFallback {
		handleJSONFallback(w, topics, domains, sortMode, timeThreshold, searchQuery)
		return
	}

	sqlQuery, args := buildSQLQuery(topics, domains, sortMode, timeThreshold, searchQuery)
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("❌ Failed to query database: %v. Using JSON fallback", err)
		handleJSONFallback(w, topics, domains, sortMode, timeThreshold, searchQuery)
		return
	}
	defer rows.Close()

	items := scanNewsItems(rows)
	if err := rows.Err(); err != nil {
		log.Printf("⚠️ Error iterating rows: %v", err)
	}

	sendSuccessResponse(w, items, time.Now())
}

func isMatchTopic(rawTags []string, reqTopics []string) bool {
	if len(reqTopics) == 0 {
		return true
	}
	for _, reqTopic := range reqTopics {
		if strings.EqualFold(reqTopic, allConst) || reqTopic == "" {
			return true
		}
		aliases := getAliases(reqTopic)
		for _, tag := range rawTags {
			tagLower := strings.ToLower(tag)
			for _, alias := range aliases {
				if tagLower == strings.ToLower(alias) {
					return true
				}
			}
		}
	}
	return false
}

func isMatchDomain(itemDomain string, reqDomains []string) bool {
	if len(reqDomains) == 0 {
		return true
	}
	for _, d := range reqDomains {
		if strings.EqualFold(itemDomain, d) {
			return true
		}
	}
	return false
}

func filterJSONItems(dataItems []api.NewsItem, topics, domains []string, timeThreshold time.Time, searchQuery string) []NewsItem {
	filtered := make([]NewsItem, 0, len(dataItems))
	for _, rawItem := range dataItems {
		if rawItem.PublishedAt.Before(timeThreshold) {
			continue
		}

		if !isMatchDomain(rawItem.Domain, domains) {
			continue
		}

		if searchQuery != "" {
			if !strings.Contains(strings.ToLower(rawItem.Title), searchQuery) &&
				!strings.Contains(strings.ToLower(rawItem.TLDR), searchQuery) {
				continue
			}
		}

		if !isMatchTopic(rawItem.Tags, topics) {
			continue
		}

		item := NewsItem{
			ID:           rawItem.ID,
			Title:        rawItem.Title,
			URL:          rawItem.URL,
			Domain:       rawItem.Domain,
			PublishedAt:  rawItem.PublishedAt,
			Category:     rawItem.Category,
			Score:        convertToScale10(rawItem.Score),
			RawScore:     rawItem.Score,
			HeatLevel:    getHeatLevel(convertToScale10(rawItem.Score)),
			TLDR:         rawItem.TLDR,
			WhyItMatters: rawItem.WhyItMatters,
			Tags:         rawItem.Tags,
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func sortJSONItems(filtered []NewsItem, sortMode string) {
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
}

func handleJSONFallback(w http.ResponseWriter, topics, domains []string, sortMode string, timeThreshold time.Time, searchQuery string) {
	data, err := api.LoadNewsData()
	if err != nil {
		if encErr := json.NewEncoder(w).Encode(Response{Success: false, Error: "Failed to load fallback JSON data"}); encErr != nil {
			log.Printf("Failed to encode error response: %v", encErr)
		}
		return
	}

	filtered := filterJSONItems(data.Items, topics, domains, timeThreshold, searchQuery)
	sortJSONItems(filtered, sortMode)

	if len(filtered) > 30 {
		filtered = filtered[:30]
	}

	sendSuccessResponse(w, filtered, data.LastUpdated)
}

func sendSuccessResponse(w http.ResponseWriter, items []NewsItem, lastUpdated time.Time) {
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"items":        items,
			"total_count":  len(items),
			"last_updated": lastUpdated,
		},
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
