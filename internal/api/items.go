package api

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

const colPublishedAt = "i.published_at"

// ItemsHandler handles item listing, filtering, search, and pagination.
type ItemsHandler struct {
	database *db.DB
}

// NewItemsHandler constructs an ItemsHandler.
func NewItemsHandler(database *db.DB) *ItemsHandler {
	return &ItemsHandler{
		database: database,
	}
}

// HandleGetItems godoc
// @Summary List items with advanced filtering and pagination
// @Tags Items
// @Produce json
// @Param search query string false "Search query"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param sources query string false "Comma separated source IDs or names"
// @Param min_relevance query int false "Minimum relevance score 0-100"
// @Param status query string false "Crawl status filter (verified, pending, done, failed, all)"
// @Param sort_by query string false "Sort field (date, relevance, credibility, impact)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 20)"
// @Success 200 {object} models.PaginatedItemsResponse
// @Router /api/items [get]
func (h *ItemsHandler) HandleGetItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	search := strings.TrimSpace(c.Query("search"))
	dateFrom := strings.TrimSpace(c.Query("date_from"))
	dateTo := strings.TrimSpace(c.Query("date_to"))

	rawSources := c.QueryArray("sources[]")
	if len(rawSources) == 0 {
		if s := c.Query("sources"); s != "" {
			rawSources = strings.Split(s, ",")
		}
	}

	rawCategories := c.QueryArray("categories[]")
	if len(rawCategories) == 0 {
		if cat := c.Query("categories"); cat != "" {
			rawCategories = strings.Split(cat, ",")
		}
	}

	minRelevance, _ := strconv.Atoi(c.DefaultQuery("min_relevance", "0"))
	status := strings.TrimSpace(c.Query("status"))
	sortBy := strings.TrimSpace(c.Query("sort_by"))
	sortOrder := strings.TrimSpace(c.Query("sort_order"))

	items, totalCount, filteredCount, err := h.queryItems(c.Request.Context(), search, dateFrom, dateTo, rawSources, rawCategories, minRelevance, status, sortBy, sortOrder, page, limit)
	if err != nil {
		slog.Error("Failed to query items", "err", err)
		c.JSON(http.StatusInternalServerError, models.PaginatedItemsResponse{
			Success: false,
			Error:   stringPtr(fmt.Sprintf("failed to fetch items: %v", err)),
		})
		return
	}

	totalPages := 1
	if filteredCount > 0 {
		totalPages = int(math.Ceil(float64(filteredCount) / float64(limit)))
	}

	c.JSON(http.StatusOK, models.PaginatedItemsResponse{
		Success:       true,
		Data:          items,
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Page:          page,
		TotalPages:    totalPages,
		LastUpdated:   time.Now().Format(time.RFC3339),
	})
}

func (h *ItemsHandler) queryItems(
	ctx context.Context,
	search, dateFrom, dateTo string,
	sources, categories []string,
	minRelevance int,
	status, sortBy, sortOrder string,
	page, limit int,
) ([]models.Item, int64, int64, error) {
	if h.database == nil || h.database.Pool == nil {
		// Mock fallback if running without live database connection
		return mockItems(), 5, 5, nil
	}

	whereStmt, args := buildWhereClauses(search, dateFrom, dateTo, sources, minRelevance, status)

	// Count Total
	var totalCount int64
	_ = h.database.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM items").Scan(&totalCount)

	// Count Filtered
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items i LEFT JOIN sources s ON i.source_id = s.id %s", whereStmt)
	var filteredCount int64
	err := h.database.Pool.QueryRow(ctx, countQuery, args...).Scan(&filteredCount)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error counting filtered items: %w", err)
	}

	orderCol, dir := getSortOrdering(sortBy, sortOrder)
	offset := (page - 1) * limit
	argIdx := len(args) + 1

	query := fmt.Sprintf(`
		SELECT 
			i.id, i.source_id, COALESCE(s.name, 'Unknown') as source_name, i.title, i.url, 
			i.published_at, i.content_hash, i.domain, i.category, i.raw_excerpt, 
			COALESCE(i.crawl_status, 'done') as crawl_status, i.crawl_error, COALESCE(i.relevance_score, 0) as relevance_score, 
			i.validated_at, i.created_at,
			sc.hot, sc.relevance, sc.credibility, sc.novelty, sc.impact, sc.engineering_value, sc.final
		FROM items i
		LEFT JOIN sources s ON i.source_id = s.id
		LEFT JOIN scores sc ON i.id = sc.item_id
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereStmt, orderCol, dir, argIdx, argIdx+1)

	argsWithPagination := append(append([]interface{}{}, args...), limit, offset)
	rows, err := h.database.Pool.Query(ctx, query, argsWithPagination...)
	if err != nil {
		return nil, totalCount, filteredCount, fmt.Errorf("error querying items page: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		it, scanErr := scanItemRow(rows)
		if scanErr != nil {
			slog.Error("Error scanning item row", "err", scanErr)
			continue
		}
		items = append(items, it)
	}

	return items, totalCount, filteredCount, nil
}

func buildWhereClauses(search, dateFrom, dateTo string, sources []string, minRelevance int, status string) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(i.title ILIKE $%d OR i.domain ILIKE $%d OR i.raw_excerpt ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= $%d", colPublishedAt, argIdx))
			args = append(args, t)
			argIdx++
		}
	}

	if dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= $%d", colPublishedAt, argIdx))
			args = append(args, t.Add(24*time.Hour))
			argIdx++
		}
	}

	if minRelevance > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("i.relevance_score >= $%d", argIdx))
		args = append(args, minRelevance)
		argIdx++
	}

	if status != "" && status != "all" {
		targetStatus := status
		if status == "verified" {
			targetStatus = "done"
		}
		whereClauses = append(whereClauses, fmt.Sprintf("i.crawl_status = $%d", argIdx))
		args = append(args, targetStatus)
		argIdx++
	}

	if len(sources) > 0 && sources[0] != "" {
		var srcOr []string
		for _, s := range sources {
			if s != "" {
				srcOr = append(srcOr, fmt.Sprintf("s.name ILIKE $%d OR i.source_id::text = $%d", argIdx, argIdx))
				args = append(args, "%"+s+"%")
				argIdx++
			}
		}
		if len(srcOr) > 0 {
			whereClauses = append(whereClauses, "("+strings.Join(srcOr, " OR ")+")")
		}
	}

	if len(whereClauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(whereClauses, " AND "), args
}

func getSortOrdering(sortBy, sortOrder string) (string, string) {
	dir := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		dir = "ASC"
	}

	switch strings.ToLower(sortBy) {
	case "relevance":
		return "i.relevance_score", dir
	case "credibility":
		return "sc.credibility", dir
	case "impact":
		return "sc.impact", dir
	case "oldest":
		return colPublishedAt, "ASC"
	default:
		return colPublishedAt, dir
	}
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanItemRow(scanner rowScanner) (models.Item, error) {
	var it models.Item
	var score models.Score
	var hot, rel, cred, nov, imp, eng, fin *float64

	err := scanner.Scan(
		&it.ID, &it.SourceID, &it.SourceName, &it.Title, &it.URL,
		&it.PublishedAt, &it.ContentHash, &it.Domain, &it.Category, &it.RawExcerpt,
		&it.CrawlStatus, &it.CrawlError, &it.RelevanceScore,
		&it.ValidatedAt, &it.CreatedAt,
		&hot, &rel, &cred, &nov, &imp, &eng, &fin,
	)
	if err != nil {
		return it, err
	}

	if fin != nil {
		score.Final = *fin
		if hot != nil {
			score.Hot = *hot
		}
		if rel != nil {
			score.Relevance = *rel
		}
		if cred != nil {
			score.Credibility = *cred
		}
		if nov != nil {
			score.Novelty = *nov
		}
		if imp != nil {
			score.Impact = *imp
		}
		if eng != nil {
			score.EngineeringValue = *eng
		}
		it.Score = &score
		it.ScaledScore = score.Final * 10
	} else {
		it.ScaledScore = float64(it.RelevanceScore) / 10.0
	}

	return it, nil
}

func mockItems() []models.Item {
	now := time.Now()
	excerpt := "Comprehensive deep-dive analysis on autonomous AI agents and vector retrieval pipelines."
	return []models.Item{
		{
			ID:             uuid.New(),
			SourceID:       uuid.New(),
			SourceName:     "Hacker News",
			Title:          "DeepSeek-R1: Open Reasoning Models with Reinforcement Learning",
			URL:            "https://news.ycombinator.com/item?id=123456",
			PublishedAt:    now.Add(-2 * time.Hour),
			ContentHash:    "hash1",
			Domain:         "ycombinator.com",
			Category:       "llm",
			RawExcerpt:     &excerpt,
			CrawlStatus:    "done",
			RelevanceScore: 95,
			CreatedAt:      now,
			ScaledScore:    9.5,
		},
		{
			ID:             uuid.New(),
			SourceID:       uuid.New(),
			SourceName:     "ArXiv AI",
			Title:          "Scaling Laws for Autonomous Code Generation Agents in Production",
			URL:            "https://arxiv.org/abs/2501.00001",
			PublishedAt:    now.Add(-8 * time.Hour),
			ContentHash:    "hash2",
			Domain:         "arxiv.org",
			Category:       "agents",
			RawExcerpt:     &excerpt,
			CrawlStatus:    "done",
			RelevanceScore: 88,
			CreatedAt:      now,
			ScaledScore:    8.8,
		},
	}
}

func stringPtr(s string) *string {
	return &s
}
