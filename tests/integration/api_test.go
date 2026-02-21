//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/http/handlers"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	db          *sql.DB
	appDB       *db.DB
	router      *gin.Engine
	testSrcID   uuid.UUID
	testItemIDs []uuid.UUID
}

func (s *IntegrationTestSuite) SetupSuite() {
	databaseURL := os.Getenv("DATABASE_URL")
	require.NotEmpty(s.T(), databaseURL, "DATABASE_URL must be set")

	var err error
	s.db, err = sql.Open("pgx", databaseURL)
	require.NoError(s.T(), err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.db.PingContext(ctx)
	require.NoError(s.T(), err)

	s.setupTestData()

	gin.SetMode(gin.TestMode)
	s.router = gin.New()

	cfg := config.Load()
	appDB, err := db.New(cfg)
	require.NoError(s.T(), err)
	s.appDB = appDB

	h := handlers.New(appDB)

	s.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.router.GET("/v1/feed", h.GetFeed)
	s.router.GET("/v1/items/:id", h.GetItem)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		s.cleanupTestData()
		s.db.Close()
	}
	if s.appDB != nil {
		s.appDB.Close()
	}
}

func (s *IntegrationTestSuite) setupTestData() {
	ctx := context.Background()

	s.testSrcID = uuid.New()
	s.testItemIDs = make([]uuid.UUID, 5)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sources (id, name, type, category, url, enabled, status, created_at, updated_at)
		VALUES ($1, 'Test Source', 'rss_atom', 'news', 'https://test.com/feed', true, 'active', NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, s.testSrcID)
	require.NoError(s.T(), err)

	for i := 0; i < 5; i++ {
		itemID := uuid.New()
		s.testItemIDs[i] = itemID
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO items (id, source_id, title, url, published_at, content_hash, domain, category, raw_excerpt, created_at)
			VALUES ($1, $2, $3, $4, NOW() - ($5::integer * INTERVAL '1 hour'), $6, 'test.com', 'news', NULL, NOW())
			ON CONFLICT (id) DO NOTHING
		`,
			itemID,
			s.testSrcID,
			fmt.Sprintf("Test Article %d", i),
			fmt.Sprintf("https://test.com/article-%d", i),
			i,
			"hash-"+itemID.String(),
		)
		require.NoError(s.T(), err)

		_, err = s.db.ExecContext(ctx, `
			INSERT INTO scores (item_id, hot, relevance, credibility, novelty, final, computed_at)
			VALUES ($1, 0.5, 0.5, 0.5, 0.5, 0.5 + ($2 * 0.1), NOW())
			ON CONFLICT (item_id) DO UPDATE SET final = 0.5 + ($2 * 0.1)
		`, itemID, i)
		require.NoError(s.T(), err)
	}
}

func (s *IntegrationTestSuite) cleanupTestData() {
	ctx := context.Background()
	for _, id := range s.testItemIDs {
		_, _ = s.db.ExecContext(ctx, `DELETE FROM scores WHERE item_id = $1`, id)
		_, _ = s.db.ExecContext(ctx, `DELETE FROM items WHERE id = $1`, id)
	}
	_, _ = s.db.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, s.testSrcID)
}

func (s *IntegrationTestSuite) TestHealthCheck() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", http.NoBody)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "ok")
}

func (s *IntegrationTestSuite) TestGetFeed() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/feed?date=today", http.NoBody)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "items")
}

func (s *IntegrationTestSuite) TestGetItem() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/items/"+s.testItemIDs[0].String(), http.NoBody)
	s.router.ServeHTTP(w, req)
	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "Test Article 0")
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(IntegrationTestSuite))
}
