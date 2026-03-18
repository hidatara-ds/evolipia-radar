package crawler

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/pkg/ai"
	"github.com/hidatara-ds/evolipia-radar/pkg/cluster"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Orchestrator manages the crawling lifecycle and agents.
type Orchestrator struct {
	agents          []DiscoveryAgent
	budget          *CrawlBudget
	clusterService  *ai.ClusterService
	inMemClusterSvc *cluster.Service
	DryRun          bool
	metrics         *Metrics
}

// NewOrchestrator wires together all agents and binds them to the AI clustering brain.
func NewOrchestrator(clusterSvc *ai.ClusterService, inMemSvc *cluster.Service, metrics *Metrics, pool *pgxpool.Pool, dryRun bool) *Orchestrator {
	// Initialize with strict zero-cost budget: 50 requests per hour max
	budget := NewCrawlBudget(50, metrics, pool)

	return &Orchestrator{
		agents: []DiscoveryAgent{
			NewRSSAgent(),
			NewTrendingAgent(),
		},
		budget:          budget,
		clusterService:  clusterSvc,
		inMemClusterSvc: inMemSvc,
		DryRun:          dryRun,
		metrics:         metrics,
	}
}

// Start begins a blocking loop that triggers Discovery on an interval.
func (o *Orchestrator) Start(ctx context.Context, interval time.Duration) {
	log.Printf("[ORCHESTRATOR] Starting Multi-Agent Discovery System. Interval: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately once
	o.RunCycle(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[ORCHESTRATOR] Shutting down...")
			return
		case <-ticker.C:
			o.RunCycle(ctx)
			o.UpdateClusterMetrics(ctx)
		}
	}
}

// RunCycle executes one pass of all Discovery agents.
func (o *Orchestrator) RunCycle(ctx context.Context) map[string]int {
	log.Printf("[ORCHESTRATOR] Beginning Discovery Cycle (DryRun: %v)", o.DryRun)
	o.budget.LogStatus()

	stats := map[string]int{
		"discovered": 0,
		"accepted":   0,
		"rejected":   0,
	}

	for _, agent := range o.agents {
		// Phase 3.5: Agent Jitter (0-10 seconds)
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		jitter := time.Duration(n.Int64()) * time.Second
		log.Printf("[ORCHESTRATOR] Applying jitter %v before dispatching %s...", jitter, agent.Name())
		time.Sleep(jitter)

		// Limit each agent to fetching up to 10 candidates to inspect
		articles, err := agent.Crawl(ctx, 10)
		if err != nil {
			log.Printf("[ORCHESTRATOR] Agent %s failed: %v", agent.Name(), err)
			continue
		}

		stats["discovered"] += len(articles)

		for _, art := range articles {
			// 1. Budget & Deduplication Check (Fast rejection)
			if !o.budget.Consume(ctx, art.Link) {
				stats["rejected"]++
				continue // Skip if already seen or over hourly limit
			}

			stats["accepted"]++

			// Phase 3.5: DRY RUN Mode
			if o.DryRun {
				log.Printf("[DRY-RUN] Discovered: %s | Source: %s", art.Title, art.Source)
				continue // Bypass cluster ingestion
			}

			// Phase 5: Fast In-Memory Clustering Routing
			if o.inMemClusterSvc != nil {
				err := o.inMemClusterSvc.ProcessArticle(ctx, art.Title, art.Content, art.Link)
				if err != nil {
					log.Printf("[ORCHESTRATOR] In-Memory Clustering failed for %s: %v", art.Link, err)
				}
			}

			// Generate a fake article ID for ingestion since URLs are our real primary keys here
			artID := uuid.New()

			// 2. Feed into the Persistence AI Cluster Engine (If not dry run)
			if o.clusterService != nil && !o.DryRun {
				err := o.clusterService.ProcessArticle(ctx, artID, art.Title, art.Content, art.Link)
				if err != nil {
					log.Printf("[ORCHESTRATOR] Cluster pipeline failed for article %s: %v", art.Link, err)
				}
			}
		}
	}

	return stats
}

// UpdateClusterMetrics fetches DB stats for the /metrics endpoint
func (o *Orchestrator) UpdateClusterMetrics(ctx context.Context) {
	o.metrics.mu.Lock()
	defer o.metrics.mu.Unlock()

	var totalScore float64
	var titles []string
	var clustersCount int

	// Prefer Phase 5 In-Memory metrics if active
	if o.inMemClusterSvc != nil {
		clusters := o.inMemClusterSvc.GetTopClusters(10)
		clustersCount = o.inMemClusterSvc.GetTotalClusters()
		if len(clusters) > 0 {
			for _, c := range clusters {
				totalScore += c.Score
				titles = append(titles, c.Label)
			}
		}
	} else if o.clusterService != nil && !o.DryRun {
		// Legacy DB metrics
		clusters, err := o.clusterService.GetTopClusters(ctx, 10)
		if err == nil && len(clusters) > 0 {
			clustersCount = len(clusters) // Simplified count
			for _, c := range clusters {
				totalScore += c.Score
				titles = append(titles, c.Title)
			}
		}
	}

	o.metrics.ClustersCount = clustersCount
	if clustersCount > 0 {
		o.metrics.AvgClusterScore = totalScore / float64(len(titles))
	} else {
		o.metrics.AvgClusterScore = 0
	}
	o.metrics.TopClusterTitles = titles
}
