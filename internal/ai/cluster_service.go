package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Cluster represents a synthesized group of related articles.
type Cluster struct {
	ID        uuid.UUID
	Title     string
	Summary   string
	Embedding []float32
	Score     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ClusterService handles semantic embedding generation, vector search, and clustering logic.
type ClusterService struct {
	aiService *Service
	db        *pgxpool.Pool
	threshold float64
}

// NewClusterService creates a new ClusterService instance.
func NewClusterService(aiService *Service, db *pgxpool.Pool) *ClusterService {
	return &ClusterService{
		aiService: aiService,
		db:        db,
		threshold: 0.85, // Cosine similarity threshold for clustering
	}
}

// ProcessArticle is the core pipeline: it takes a raw article, generates its embedding,
// searches the vector DB for the nearest cluster, and assigns/creates a cluster.
func (s *ClusterService) ProcessArticle(ctx context.Context, articleID uuid.UUID, title, content, link string) error {
	// Phase 2.8: PreProcessFilter - Stop noise early
	if err := PreProcessFilter(link, title, content); err != nil {
		log.Printf("[CLUSTER] Article %s rejected: %v", articleID, err)
		return nil // Drop gracefully; not an internal error
	}

	// 1. Generate text embedding for the article.
	// We merge title and content for a richer semantic representation.
	textToEmbed := fmt.Sprintf("%s\n\n%s", title, content)

	embedResp, err := s.aiService.Embed(ctx, EmbeddingRequest{
		Input: textToEmbed,
	})

	// Phase 2.9 Fallback Behavior
	if err != nil {
		if strings.Contains(err.Error(), "BUDGET EXHAUSTED") || strings.Contains(err.Error(), "429") {
			log.Printf("[BUDGET-FALLBACK] Skipping vector embedding for article %s due to budget limits. Creating zero-vector fallback.", articleID)
			return s.createFallbackCluster(ctx, articleID, title, content)
		}
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
	embedding := embedResp.Embedding

	// 2. Search for the nearest existing cluster using pgvector.
	nearestClusterID, nearestText, distance, err := s.findNearestCluster(ctx, embedding)
	if err != nil {
		return fmt.Errorf("failed during nearest cluster search: %w", err)
	}

	// pgvector `vector_cosine_ops` uses cosine distance. Cosine Similarity = 1 - Cosine Distance
	embedSimilarity := 1.0 - distance

	// Phase 2.5: Hybrid Similarity
	// Combine Embedding Similarity (0.8 weight) + Keyword Overlap (0.2 weight).
	keywordSimilarity := KeywordOverlap(textToEmbed, nearestText)
	hybridSimilarity := (embedSimilarity * 0.8) + (keywordSimilarity * 0.2)

	if nearestClusterID != uuid.Nil && hybridSimilarity >= s.threshold {
		// 3a. Match found! Attach the article to this existing cluster.
		log.Printf("[CLUSTER] Matched article %s to cluster %s (embed: %.3f, keyword: %.3f, hybrid: %.3f)",
			articleID, nearestClusterID, embedSimilarity, keywordSimilarity, hybridSimilarity)

		err = s.assignToCluster(ctx, nearestClusterID, articleID)
		if err == nil {
			// Update the cluster ranking score dynamically
			_ = s.incrementClusterScore(ctx, nearestClusterID, textToEmbed)
		}
		return err
	}

	// 3b. No match found or similarity < threshold. Create a new cluster.
	log.Printf("[CLUSTER] Creating new cluster for article %s (nearest hybrid was %.3f)", articleID, hybridSimilarity)

	// Base score on the initial text signals
	initialScore := 1.0 + (float64(len(extractKeywords(textToEmbed))) * 0.5)
	return s.createNewCluster(ctx, articleID, title, content, embedding, initialScore)
}

// findNearestCluster queries pgvector for the closest cluster by cosine distance
// and returns the ID, the string contents (title+summary) of the cluster for hybrid scoring,
// and the raw embedding distance.
func (s *ClusterService) findNearestCluster(ctx context.Context, embedding []float32) (uuid.UUID, string, float64, error) {
	// Convert Go float32 slice to pgvector format literal (e.g., '[0.1, 0.2, ...]')
	embeddingStr := formatVector(embedding)

	query := `
		SELECT id, title, summary, embedding <=> $1::vector AS distance
		FROM clusters
		ORDER BY embedding <=> $1::vector
		LIMIT 1;
	`

	var clusterID uuid.UUID
	var clusterTitle, clusterSummary string
	var distance float64

	err := s.db.QueryRow(ctx, query, embeddingStr).Scan(&clusterID, &clusterTitle, &clusterSummary, &distance)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return uuid.Nil, "", 1.0, nil
		}
		return uuid.Nil, "", 0, err
	}

	clusterText := fmt.Sprintf("%s\n\n%s", clusterTitle, clusterSummary)
	return clusterID, clusterText, distance, nil
}

// assignToCluster links an article to an existing cluster.
func (s *ClusterService) assignToCluster(ctx context.Context, clusterID, articleID uuid.UUID) error {
	query := `
		INSERT INTO cluster_sources (cluster_id, article_id)
		VALUES ($1, $2)
		ON CONFLICT (cluster_id, article_id) DO NOTHING;
	`
	_, err := s.db.Exec(ctx, query, clusterID, articleID)
	return err
}

// incrementClusterScore increases a cluster's score when a new source is attached.
// If the score boost pushes the cluster past major thresholds, it will trigger a recompute of the insight.
func (s *ClusterService) incrementClusterScore(ctx context.Context, clusterID uuid.UUID, newArticleText string) error {
	// Base score +1 for a new source. Add fractional points based on keyword density.
	boost := 1.0 + (float64(len(extractKeywords(newArticleText))) * 0.2)

	// Phase 2.8: Fetch the current score to determine if a recompute is needed
	var currentScore float64
	var clusterTitle, clusterSummary string

	// (Transaction omitted for simplicity of demonstration, assuming sync flow)
	err := s.db.QueryRow(ctx, `SELECT score, title, summary FROM clusters WHERE id = $1`, clusterID).Scan(&currentScore, &clusterTitle, &clusterSummary)
	if err != nil {
		return err
	}

	query := `UPDATE clusters SET score = score + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = s.db.Exec(ctx, query, boost, clusterID)
	if err != nil {
		return err
	}

	newScore := currentScore + boost

	// Phase 2.8: Summary Recompute Logic
	// Phase 2.9 (Smart Usage): Only recompute if going from 1 signal to massive virality (cost saving)
	// We skip intermediate updates to save budget.
	if currentScore < 5.0 && newScore >= 5.0 {
		log.Printf("[CLUSTER] Cluster %s crossed High-Importance Threshold (Score: %.1f). Recomputing Summary.", clusterID, newScore)

		// In a real system you'd pull down *all* existing text from cluster_sources, but
		// for this sync flow we'll blend the old summary + new breaking article.
		blendedText := fmt.Sprintf("Previous Insight: %s\n\nNew Breaking Detail: %s", clusterSummary, newArticleText)

		recomputed, rErr := s.generateGroundedInsight(ctx, blendedText)
		if rErr == nil {
			log.Printf("[CLUSTER] Summary recomputed for %s", clusterID)
			s.db.Exec(ctx, `UPDATE clusters SET summary = $1 WHERE id = $2`, recomputed, clusterID)
		} else if strings.Contains(rErr.Error(), "BUDGET EXHAUSTED") {
			log.Printf("[BUDGET-FALLBACK] Skipping expensive recompute for cluster %s", clusterID)
		}
	}

	return nil
}

// generateGroundedInsight applies Phase 2.8 Guardrails against hallucination and speculation
func (s *ClusterService) generateGroundedInsight(ctx context.Context, content string) (string, error) {
	guardrailPrompt := `You are an elite, highly strict AI tech analyst. 
Summarize the event in exactly two sentences.

RULES:
1. Sentence one must state the objective, verifiable technical event.
2. Sentence two must state the direct engineering or market implication.
3. DO NOT speculate. DO NOT use words like "might", "could", "perhaps", or "promises to".
4. If a claim is unverified, state "claims to".
5. Base the insight strictly on the provided text.`

	summaryResp, err := s.aiService.Summarize(ctx, SummarizeRequest{
		Text:        content,
		Instruction: guardrailPrompt,
	})
	if err != nil {
		return "", err
	}
	return summaryResp.Summary, nil
}

// createNewCluster creates a brand new insight cluster from the source article.
func (s *ClusterService) createNewCluster(ctx context.Context, articleID uuid.UUID, title, content string, embedding []float32, initialScore float64) error {

	// 1. Generate an Intelligent Summary (Insight Format) with Guardrails
	summary, err := s.generateGroundedInsight(ctx, content)
	if err != nil {
		if strings.Contains(err.Error(), "BUDGET EXHAUSTED") || strings.Contains(err.Error(), "429") {
			log.Printf("[BUDGET-FALLBACK] LLM summary failed. Falling back to raw text for %s", articleID)
			summary = content // Fallback: Just use the raw article text
		} else {
			return fmt.Errorf("failed to generate cluster insight: %w", err)
		}
	}

	// 2. Generate a Generative Semantic Title
	// Instead of copying the raw article headline (e.g., "Meta just dropped Llama 3"),
	// generate a standardized cluster name (e.g., "Meta Releases Llama 3 Open-Source Framework").
	titlePrompt := "Return ONLY a highly professional, 4-7 word title for this tech cluster."
	titleResp, err := s.aiService.Summarize(ctx, SummarizeRequest{
		Text:        content,
		Instruction: titlePrompt,
	})
	if err != nil {
		// Graceful fallback to raw title if the LLM call fails
		titleResp = &SummarizeResponse{Summary: title}
	}

	embeddingStr := formatVector(embedding)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert new cluster
	var newClusterID uuid.UUID
	insertClusterQuery := `
		INSERT INTO clusters (title, summary, embedding, score)
		VALUES ($1, $2, $3::vector, $4)
		RETURNING id;
	`
	err = tx.QueryRow(ctx, insertClusterQuery, titleResp.Summary, summary, embeddingStr, initialScore).Scan(&newClusterID)
	if err != nil {
		return fmt.Errorf("failed to insert cluster: %w", err)
	}

	// Link the source article
	insertSourceQuery := `
		INSERT INTO cluster_sources (cluster_id, article_id)
		VALUES ($1, $2);
	`
	_, err = tx.Exec(ctx, insertSourceQuery, newClusterID, articleID)
	if err != nil {
		return fmt.Errorf("failed to insert cluster source: %w", err)
	}

	return tx.Commit(ctx)
}

// createFallbackCluster completely bypasses all LLM generation when the free tier budget is exhausted.
// It creates a standalone cluster using entirely raw text and a zero-vector so the system stays online.
func (s *ClusterService) createFallbackCluster(ctx context.Context, articleID uuid.UUID, title, content string) error {
	zeroVector := make([]float32, 1536) // Zero array so pgvector doesn't crash
	
	// Truncate fallback content to prevent massive UI rendering issues
	fallbackSummary := content
	if len(content) > 300 {
		fallbackSummary = content[:297] + "..."
	}

	// Assign basic low score
	initialScore := 1.0 + (float64(len(extractKeywords(content))) * 0.5)

	embeddingStr := formatVector(zeroVector)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var newClusterID uuid.UUID
	insertClusterQuery := `
		INSERT INTO clusters (title, summary, embedding, score)
		VALUES ($1, $2, $3::vector, $4)
		RETURNING id;
	`
	err = tx.QueryRow(ctx, insertClusterQuery, title, fallbackSummary, embeddingStr, initialScore).Scan(&newClusterID)
	if err != nil {
		return fmt.Errorf("failed to insert fallback cluster: %w", err)
	}

	insertSourceQuery := `
		INSERT INTO cluster_sources (cluster_id, article_id)
		VALUES ($1, $2);
	`
	_, err = tx.Exec(ctx, insertSourceQuery, newClusterID, articleID)
	return tx.Commit(ctx)
}

// formatVector converts a []float32 into a string formatted like "[0.1,-0.2,0.3]" for pgvector
func formatVector(v []float32) string {
	res := "["
	for i, val := range v {
		res += fmt.Sprintf("%f", val)
		if i < len(v)-1 {
			res += ","
		}
	}
	res += "]"
	return res
}

// MergeFragmentedClusters is a Phase 2.8 stability function meant to be run on a cron job.
// It searches for clusters that belong to the exact same event but were slightly under
// the hybrid threshold during the initial race condition of ingestion.
func (s *ClusterService) MergeFragmentedClusters(ctx context.Context) (int, error) {
	// 1. Find pairs of highly identical clusters
	query := `
		WITH pairs AS (
			SELECT a.id as id_a, b.id as id_b,
			       a.title as title_a, b.title as title_b,
				   a.score as score_a, b.score as score_b,
				   1 - (a.embedding <=> b.embedding) as similarity
			FROM clusters a
			JOIN clusters b ON a.id > b.id
			WHERE (1 - (a.embedding <=> b.embedding)) > 0.94 -- extremely strict merge threshold
		)
		SELECT id_a, id_b, score_a, score_b FROM pairs
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed searching for fragmented pairs: %w", err)
	}
	defer rows.Close()

	merges := 0
	for rows.Next() {
		var idA, idB uuid.UUID
		var scoreA, scoreB float64
		if err := rows.Scan(&idA, &idB, &scoreA, &scoreB); err != nil {
			continue // Skip problematic rows
		}

		// Keep the one with the higher score, merge the weaker into the stronger
		winner, loser := idA, idB
		if scoreB > scoreA {
			winner, loser = idB, idA
		}

		// Move all sources from loser -> winner
		updateSources := `UPDATE cluster_sources SET cluster_id = $1 WHERE cluster_id = $2`
		if _, err := s.db.Exec(ctx, updateSources, winner, loser); err == nil {
			// Add loser's score to the winner
			s.db.Exec(ctx, `UPDATE clusters SET score = score + $1 WHERE id = $2`, scoreB*0.5, winner) // give half points
			// Delete the fragmented loser
			s.db.Exec(ctx, `DELETE FROM clusters WHERE id = $1`, loser)
			merges++
		}
	}

	if merges > 0 {
		log.Printf("[CLUSTER-MERGE] Successfully identified and merged %d fragmented cluster pairs", merges)
	}
	return merges, nil
}

// GetTopClusters fetches clusters with a Time-Decay applied to the ranking score.
func (s *ClusterService) GetTopClusters(ctx context.Context, limit int) ([]Cluster, error) {
	// Formula: Effective Score = Base_Score / (1 + age_in_hours)^decay_gravity
	// Using PostgreSQL mathematically directly on read.
	query := `
		SELECT id, title, summary, score, created_at, updated_at,
		       (score / POWER((1 + EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - created_at))/3600.0), 1.5)) AS decay_score
		FROM clusters
		ORDER BY decay_score DESC
		LIMIT $1;
	`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Cluster
	for rows.Next() {
		var c Cluster
		var decayScore float64 // Read but not stored

		if err := rows.Scan(&c.ID, &c.Title, &c.Summary, &c.Score, &c.CreatedAt, &c.UpdatedAt, &decayScore); err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	return results, nil
}
