package cluster

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Article represents a crawled document inside the semantic engine.
type Article struct {
	ID        string
	Title     string
	Content   string
	Link      string
	Embedding []float64
	CreatedAt time.Time
}

// Cluster represents a semantic grouping of related Articles.
type Cluster struct {
	ID        string
	Label     string
	Articles  []Article
	Score     float64
	CreatedAt time.Time
}

// ClusterService handles in-memory vector intelligence and deduplication.
type ClusterService struct {
	mu           sync.RWMutex
	clusters     []*Cluster
	embedder     EmbeddingProvider
	threshold    float64
}

func NewClusterService(embedder EmbeddingProvider) *ClusterService {
	return &ClusterService{
		clusters:  make([]*Cluster, 0),
		embedder:  embedder,
		threshold: 0.80, // cosine similarity threshold for grouping
	}
}

// ProcessArticle analyzes, embeds, and clusters an incoming article.
func (s *ClusterService) ProcessArticle(ctx context.Context, title, content, link string) error {
	// 1. Generate semantic embedding
	emb, err := s.embedder.Embed(title + " " + content)
	if err != nil {
		return err // Or fallback gracefully up the chain
	}

	art := Article{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		Link:      link,
		Embedding: emb,
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var bestMatch *Cluster
	var maxSim float64 = -1.0

	// 2. Greedy cosine similarity scan against existing cluster centroids (first article representation)
	for _, c := range s.clusters {
		if len(c.Articles) > 0 {
			sim := cosineSimilarity(emb, c.Articles[0].Embedding)
			if sim > maxSim {
				maxSim = sim
				bestMatch = c
			}
		}
	}

	// 3. Cluster Assignment or Spawning
	if maxSim >= s.threshold && bestMatch != nil {
		bestMatch.Articles = append(bestMatch.Articles, art)
		bestMatch.Score += 1.0 // Simple engagement score boost
	} else {
		newC := &Cluster{
			ID:        uuid.New().String(),
			Label:     title, // Use title as generative fallback
			Articles:  []Article{art},
			Score:     1.0,
			CreatedAt: time.Now(),
		}
		s.clusters = append(s.clusters, newC)
	}

	return nil
}

// GetTopClusters returns the most active semantic groups.
func (s *ClusterService) GetTopClusters(limit int) []*Cluster {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make([]*Cluster, len(s.clusters))
	copy(cp, s.clusters)

	sort.Slice(cp, func(i, j int) bool {
		return cp[i].Score > cp[j].Score
	})

	if len(cp) > limit {
		return cp[:limit]
	}
	return cp
}

func (s *ClusterService) GetTotalClusters() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clusters)
}

// cosineSimilarity helper function without external heavy ML libraries.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
