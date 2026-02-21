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
	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/http/handlers"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	db     *sql.DB  // raw DB handle for setting up test data
	appDB  *db.DB   // application DB pool used by handlers
	router *gin.Engine
}

func (s *IntegrationTestSuite) SetupSuite() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	require.NotEmpty(s.T(), databaseURL, "DATABASE_URL must be set")

	// Connect to database
	var err error
	s.db, err = sql.Open("pgx", databaseURL)
	require.NoError(s.T(), err)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.db.PingContext(ctx)
	require.NoError(s.T(), err)

	// Setup test data
	s.setupTestData()

	// Setup router
	gin.SetMode(gin.TestMode)
	s.router = gin.New()

	cfg := config.Load()
	appDB, err := db.New(cfg)
	require.NoError(s.T(), err)
	s.appDB = appDB

	h := handlers.New(appDB)

	// Mirror the routing in cmd/api/main.go
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

	// Insert test source
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sources (id, name, type, url, category, enabled, created_at, updated_at)
		VALUES ('test-source-1', 'Test Source', 'rss_atom', 'https://test.com/feed', 'news', true, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(s.T(), err)

	// Insert test items
	for i := 0; i < 5; i++ {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO content_items (id, source_id, title, url, summary, score, published_at, created_at)
			VALUES (
				$1, 
				'test-source-1', 
				$2, 
				$3, 
				'Test summary', 
				$4, 
				NOW() - INTERVAL '$5 hours',
				NOW()
			)
			ON CONFLICT (id) DO NOTHING
		`,
			fmt.Sprintf("test-item-%d", i),
			fmt.Sprintf("Test Article %d", i),
			fmt.Sprintf("https://test.com/article-%d", i),
			0.5+float64(i)*0.1,
			i,
		)
		require.NoError(s.T(), err)
	}
}

func (s *IntegrationTestSuite) cleanupTestData() {
	ctx := context.Background()
	_, err := s.db.ExecContext(ctx, `DELETE FROM content_items WHERE id LIKE 'test-%'`)
	require.NoError(s.T(), err)
	_, err = s.db.ExecContext(ctx, `DELETE FROM sources WHERE id LIKE 'test-%'`)
	require.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestHealthCheck() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "ok")
}

func (s *IntegrationTestSuite) TestGetFeed() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/feed?date=today", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "items")
}

func (s *IntegrationTestSuite) TestGetItem() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/items/test-item-0", nil)
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
