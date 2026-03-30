package ai

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/models"
)

// ============================================================================
// Keyword Helpers (kept from original scaffold)
// ============================================================================

// extractKeywords is a super lightweight local keyword extractor to aid in hybrid scoring.
// For production, this could be an NER agent, but we are keeping it simple per constraints.
func extractKeywords(text string) map[string]bool {
	// A small set of high-signal tech keywords
	signals := []string{
		"llm", "gpt", "rag", "gemini", "llama", "mistral", "transformer",
		"openai", "anthropic", "google", "meta", "microsoft", "apple",
		"mlops", "kubernetes", "GPU", "nvidia", "funding", "startup", "release",
	}

	lowerText := strings.ToLower(text)
	found := make(map[string]bool)

	for _, s := range signals {
		if strings.Contains(lowerText, strings.ToLower(s)) {
			found[strings.ToLower(s)] = true
		}
	}

	return found
}

// KeywordOverlap calculates a fast Jaccard-like similarity (0.0 to 1.0)
// based on the overlapping presence of high-signal tech keywords.
func KeywordOverlap(textA, textB string) float64 {
	aTerms := extractKeywords(textA)
	bTerms := extractKeywords(textB)

	if len(aTerms) == 0 && len(bTerms) == 0 {
		return 0.0 // No overlap and no keywords
	}

	intersection := 0
	for k := range aTerms {
		if bTerms[k] {
			intersection++
		}
	}

	union := len(aTerms) + len(bTerms) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// ============================================================================
// Hybrid Search Engine
// ============================================================================

// HybridResult represents a single search result with scoring breakdown.
type HybridResult struct {
	Item          models.Item `json:"item"`
	FinalScore    float64     `json:"final_score"`
	SemanticScore float64     `json:"semantic_score"`
	TextScore     float64     `json:"text_score"`
	Mode          string      `json:"mode"` // "text", "semantic", or "hybrid"
}

// HybridSearcher combines text-based (ILIKE) and semantic (vector cosine) search.
type HybridSearcher struct {
	aiService *Service
	database  *db.DB
}

// NewHybridSearcher creates a new HybridSearcher instance.
func NewHybridSearcher(aiService *Service, database *db.DB) *HybridSearcher {
	return &HybridSearcher{
		aiService: aiService,
		database:  database,
	}
}

// Search executes a search with the specified mode: "text", "semantic", or "hybrid".
// For "hybrid", it runs text and semantic searches concurrently and merges with RRF.
func (h *HybridSearcher) Search(ctx context.Context, query string, mode string, limit int) ([]HybridResult, error) {
	switch mode {
	case "semantic":
		return h.semanticSearch(ctx, query, limit)
	case "text":
		return h.textSearch(ctx, query, limit)
	case "hybrid":
		return h.hybridSearch(ctx, query, limit)
	default:
		return h.hybridSearch(ctx, query, limit)
	}
}

// textSearch uses the existing ILIKE-based search in ItemRepository.
func (h *HybridSearcher) textSearch(ctx context.Context, query string, limit int) ([]HybridResult, error) {
	itemRepo := db.NewItemRepository(h.database)
	items, _, err := itemRepo.Search(ctx, query, nil, limit, 0)
	if err != nil {
		return nil, err
	}

	results := make([]HybridResult, len(items))
	for i, item := range items {
		// Text results get a position-based score: top result = 1.0, decaying
		textScore := 1.0 - (float64(i) / float64(max(len(items), 1)))
		results[i] = HybridResult{
			Item:          item,
			FinalScore:    textScore,
			SemanticScore: 0,
			TextScore:     textScore,
			Mode:          "text",
		}
	}
	return results, nil
}

// semanticSearch embeds the query and searches by cosine similarity.
func (h *HybridSearcher) semanticSearch(ctx context.Context, query string, limit int) ([]HybridResult, error) {
	// Step 1: Generate query embedding via the AI service (budget-tracked)
	embedResp, err := h.aiService.Embed(ctx, EmbeddingRequest{
		Input: query,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: Search by embedding in DB
	itemRepo := db.NewItemRepository(h.database)
	scoredItems, err := itemRepo.SearchByEmbedding(ctx, embedResp.Embedding, limit)
	if err != nil {
		return nil, err
	}

	results := make([]HybridResult, len(scoredItems))
	for i, si := range scoredItems {
		results[i] = HybridResult{
			Item:          si.Item,
			FinalScore:    si.Similarity,
			SemanticScore: si.Similarity,
			TextScore:     0,
			Mode:          "semantic",
		}
	}
	return results, nil
}

// hybridSearch runs text and semantic searches concurrently, then merges with RRF.
func (h *HybridSearcher) hybridSearch(ctx context.Context, query string, limit int) ([]HybridResult, error) {
	var (
		textResults     []HybridResult
		semanticResults []HybridResult
		textErr         error
		semanticErr     error
		wg              sync.WaitGroup
	)

	// Fetch more than `limit` from each source for better RRF fusion quality
	fetchLimit := limit * 2
	if fetchLimit < 20 {
		fetchLimit = 20
	}

	// Run both searches concurrently
	wg.Add(2)

	go func() {
		defer wg.Done()
		textResults, textErr = h.textSearch(ctx, query, fetchLimit)
	}()

	go func() {
		defer wg.Done()
		semanticResults, semanticErr = h.semanticSearch(ctx, query, fetchLimit)
	}()

	wg.Wait()

	// If semantic fails (budget exhausted, no API key, etc.), fallback to text-only
	if semanticErr != nil {
		log.Printf("[HYBRID] Semantic search failed, falling back to text: %v", semanticErr)
		if textErr != nil {
			return nil, textErr
		}
		// Cap to requested limit
		if len(textResults) > limit {
			textResults = textResults[:limit]
		}
		return textResults, nil
	}

	// If text fails but semantic succeeded, return semantic only
	if textErr != nil {
		log.Printf("[HYBRID] Text search failed, using semantic only: %v", textErr)
		if len(semanticResults) > limit {
			semanticResults = semanticResults[:limit]
		}
		return semanticResults, nil
	}

	// Merge using Reciprocal Rank Fusion
	merged := reciprocalRankFusion(textResults, semanticResults, limit)
	return merged, nil
}

// reciprocalRankFusion merges two ranked result lists using the RRF algorithm.
//
// RRF score = Σ 1/(k + rank_i) across all lists where the item appears.
// k=60 is the standard constant from the original RRF paper (Cormack et al., 2009).
//
// Weights: semantic gets 60% weight, text gets 40% weight.
func reciprocalRankFusion(textResults, semanticResults []HybridResult, limit int) []HybridResult {
	const k = 60
	const semanticWeight = 0.6
	const textWeight = 0.4

	type fusedEntry struct {
		result        HybridResult
		rrfScore      float64
		textRank      int // 0 = not present in text results
		semanticRank  int // 0 = not present in semantic results
		semanticScore float64
		textScore     float64
	}

	// Map item ID → fusedEntry
	fusionMap := make(map[string]*fusedEntry)

	// Process text results
	for rank, r := range textResults {
		idStr := r.Item.ID.String()
		entry, exists := fusionMap[idStr]
		if !exists {
			entry = &fusedEntry{result: r}
			fusionMap[idStr] = entry
		}
		entry.textRank = rank + 1
		entry.textScore = r.TextScore
		entry.rrfScore += textWeight * (1.0 / float64(k+rank+1))
	}

	// Process semantic results
	for rank, r := range semanticResults {
		idStr := r.Item.ID.String()
		entry, exists := fusionMap[idStr]
		if !exists {
			entry = &fusedEntry{result: r}
			fusionMap[idStr] = entry
		}
		entry.semanticRank = rank + 1
		entry.semanticScore = r.SemanticScore
		entry.rrfScore += semanticWeight * (1.0 / float64(k+rank+1))
	}

	// Collect all entries
	entries := make([]*fusedEntry, 0, len(fusionMap))
	for _, entry := range fusionMap {
		entries = append(entries, entry)
	}

	// Sort by RRF score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].rrfScore > entries[j].rrfScore
	})

	// Cap to limit
	if len(entries) > limit {
		entries = entries[:limit]
	}

	// Convert to HybridResult
	results := make([]HybridResult, len(entries))
	for i, entry := range entries {
		results[i] = HybridResult{
			Item:          entry.result.Item,
			FinalScore:    entry.rrfScore,
			SemanticScore: entry.semanticScore,
			TextScore:     entry.textScore,
			Mode:          "hybrid",
		}
	}

	return results
}
