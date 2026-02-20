package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/hidatara-ds/evolipia-radar/internal/services"
)

type Handlers struct {
	sourceService *services.SourceService
	feedService   *services.FeedService
}

func New(database *db.DB) *Handlers {
	return &Handlers{
		sourceService: services.NewSourceService(database),
		feedService:   services.NewFeedService(database),
	}
}

func (h *Handlers) GetFeed(c *gin.Context) {
	dateStr := c.DefaultQuery("date", "today")
	topic := c.Query("topic")

	var date time.Time
	if dateStr == "today" {
		// Use Asia/Jakarta timezone
		loc, _ := time.LoadLocation("Asia/Jakarta")
		date = time.Now().In(loc)
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
	}

	var topicPtr *string
	if topic != "" {
		topicPtr = &topic
	}

	items, err := h.feedService.GetTopDaily(c.Request.Context(), date, topicPtr, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.feedService.BuildFeedResponse(c.Request.Context(), items, date, topicPtr)
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetRising(c *gin.Context) {
	windowStr := c.DefaultQuery("window", "2h")
	window, err := time.ParseDuration(windowStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid window format"})
		return
	}

	items, err := h.feedService.GetRising(c.Request.Context(), window, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.feedService.BuildRisingResponse(c.Request.Context(), items, window)
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	item, signal, score, summary, err := h.feedService.GetItemWithDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	source, err := h.sourceService.GetSourceByID(c.Request.Context(), item.SourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"id":           item.ID,
		"title":        item.Title,
		"url":          item.URL,
		"domain":       item.Domain,
		"published_at": item.PublishedAt.Format(time.RFC3339),
		"category":     item.Category,
		"source": gin.H{
			"id":   source.ID,
			"name": source.Name,
			"type": source.Type,
		},
	}

	if signal != nil {
		response["signals_latest"] = gin.H{
			"points":     signal.Points,
			"comments":   signal.Comments,
			"rank_pos":   signal.RankPos,
			"fetched_at": signal.FetchedAt.Format(time.RFC3339),
		}
	}

	if score != nil {
		response["scores"] = gin.H{
			"final":       score.Final,
			"hot":         score.Hot,
			"relevance":   score.Relevance,
			"credibility": score.Credibility,
			"novelty":     score.Novelty,
			"computed_at": score.ComputedAt.Format(time.RFC3339),
		}
	}

	if summary != nil {
		response["summary"] = gin.H{
			"tldr":           summary.TLDR,
			"why_it_matters": summary.WhyItMatters,
			"tags":           summary.Tags,
			"method":         summary.Method,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	topic := c.Query("topic")
	var topicPtr *string
	if topic != "" {
		topicPtr = &topic
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	items, total, err := h.feedService.SearchItems(c.Request.Context(), query, topicPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseItems := make([]gin.H, 0, len(items))
	for _, item := range items {
		_, _, score, summary, _ := h.feedService.GetItemWithDetails(c.Request.Context(), item.ID)

		itemResp := gin.H{
			"id":           item.ID,
			"title":        item.Title,
			"url":          item.URL,
			"published_at": item.PublishedAt.Format(time.RFC3339),
			"domain":       item.Domain,
			"final_score":  0.0,
			"tags":         []string{},
		}

		if score != nil {
			itemResp["final_score"] = score.Final
		}
		if summary != nil {
			itemResp["tags"] = summary.Tags
		}

		responseItems = append(responseItems, itemResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"q":              query,
		"topic":          topic,
		"total_estimate": total,
		"items":          responseItems,
	})
}

func (h *Handlers) ListSources(c *gin.Context) {
	sources, err := h.sourceService.ListSources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseItems := make([]gin.H, 0, len(sources))
	for _, s := range sources {
		responseItems = append(responseItems, gin.H{
			"id":                s.ID,
			"name":              s.Name,
			"type":              s.Type,
			"category":          s.Category,
			"url":               s.URL,
			"enabled":           s.Enabled,
			"status":            s.Status,
			"last_test_status":  s.LastTestStatus,
			"last_test_message": s.LastTestMessage,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sources": responseItems})
}

func (h *Handlers) CreateSource(c *gin.Context) {
	var req struct {
		Name        string          `json:"name" binding:"required"`
		Type        string          `json:"type" binding:"required"`
		Category    string          `json:"category" binding:"required"`
		URL         string          `json:"url" binding:"required"`
		MappingJSON json.RawMessage `json:"mapping_json,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source := &models.Source{
		Name:     req.Name,
		Type:     req.Type,
		Category: req.Category,
		URL:      req.URL,
		Enabled:  false,
		Status:   "pending",
	}

	if req.MappingJSON != nil {
		source.MappingJSON = req.MappingJSON
	}

	if err := h.sourceService.CreateSource(c.Request.Context(), source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      source.ID,
		"status":  source.Status,
		"enabled": source.Enabled,
	})
}

func (h *Handlers) TestSource(c *gin.Context) {
	var req struct {
		Type        string          `json:"type" binding:"required"`
		Category    string          `json:"category" binding:"required"`
		URL         string          `json:"url" binding:"required"`
		MappingJSON json.RawMessage `json:"mapping_json,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.sourceService.TestConnection(c.Request.Context(), req.Type, req.Category, req.URL, req.MappingJSON)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":     "failed",
			"error_code": "TEST_ERROR",
			"message":    err.Error(),
		})
		return
	}

	if result.Status == "ok" {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handlers) EnableSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.sourceService.SetEnabled(c.Request.Context(), id, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	source, err := h.sourceService.GetSourceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      source.ID,
		"enabled": source.Enabled,
		"status":  source.Status,
	})
}
