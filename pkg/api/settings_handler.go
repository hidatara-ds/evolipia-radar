package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

type SettingsHandler struct {
	repo *db.SettingRepository
}

func NewSettingsHandler(database *db.DB) *SettingsHandler {
	return &SettingsHandler{
		repo: db.NewSettingRepository(database),
	}
}

func (h *SettingsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/settings", h.List)
	rg.POST("/settings", h.Update)
}

func (h *SettingsHandler) List(c *gin.Context) {
	settings, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	for k, v := range req {
		if err := h.repo.Set(ctx, k, v); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
