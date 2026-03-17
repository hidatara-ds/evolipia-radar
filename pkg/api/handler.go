package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	"github.com/hidatara-ds/evolipia-radar/pkg/config"
)

// AIHandler exposes HTTP endpoints for the AI Service.
type AIHandler struct {
	service *ai.Service
}

// NewAIHandler creates a new instance of AIHandler.
func NewAIHandler(service *ai.Service) *AIHandler {
	return &AIHandler{
		service: service,
	}
}

// RegisterRoutes registers the AI endpoints to the router group.
func (h *AIHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Apply Hardening Middlewares specifically to the AI /v2 group
	v2 := router.Group("/v2")
	v2.Use(AILoggerMiddleware())
	v2.Use(AIRecoveryMiddleware())
	cfg := config.LoadAIConfig()
	v2.Use(TimeoutMiddleware(cfg.Timeout))

	{
		v2.POST("/chat", h.HandleChat)
		v2.POST("/summarize", h.HandleSummarize)
	}
}

// HandleChat godoc
// @Summary Chat with the LLM
// @Description Send a sequence of messages to the central AI service and get a completion.
// @Tags AI
// @Accept json
// @Produce json
// @Param request body ai.ChatRequest true "Chat Request"
// @Success 200 {object} ai.ChatResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/chat [post]
func (h *AIHandler) HandleChat(c *gin.Context) {
	var req ai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, ErrCodeValidation, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	if err := req.Validate(); err != nil {
		RespondWithError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	resp, err := h.service.Chat(c.Request.Context(), req)
	if err != nil {
		log.Printf("[AI GATEWAY ERROR] Chat provider error: %v\n", err)
		RespondWithError(c, http.StatusInternalServerError, ErrCodeInternal, "failed to process chat completion")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleSummarize godoc
// @Summary Summarize text
// @Description Summarize provided text with optional custom instructions via the central AI service.
// @Tags AI
// @Accept json
// @Produce json
// @Param request body ai.SummarizeRequest true "Summarize Request"
// @Success 200 {object} ai.SummarizeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/summarize [post]
func (h *AIHandler) HandleSummarize(c *gin.Context) {
	var req ai.SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, ErrCodeValidation, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	if err := req.Validate(); err != nil {
		RespondWithError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	resp, err := h.service.Summarize(c.Request.Context(), req)
	if err != nil {
		log.Printf("[AI GATEWAY ERROR] Summarize provider error: %v\n", err)
		RespondWithError(c, http.StatusInternalServerError, ErrCodeInternal, "failed to summarize text")
		return
	}

	c.JSON(http.StatusOK, resp)
}
