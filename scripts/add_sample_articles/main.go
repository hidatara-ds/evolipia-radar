// Script to add sample articles for testing missing tags
// Run with: go run scripts/add_sample_articles.go
package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type NewsItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	PublishedAt  time.Time `json:"published_at"`
	Category     string    `json:"category"`
	Score        float64   `json:"score"`
	TLDR         string    `json:"tldr,omitempty"`
	WhyItMatters string    `json:"why_it_matters,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

type NewsData struct {
	Items       []NewsItem `json:"items"`
	LastUpdated time.Time  `json:"last_updated"`
	TotalCount  int        `json:"total_count"`
}

func main() {
	inputPath := "data/news.json"

	log.Println("📝 Adding sample articles for missing tags...")

	// Read existing news.json
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to read %s: %v", inputPath, err)
	}

	var newsData NewsData
	if err := json.Unmarshal(data, &newsData); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	now := time.Now()

	// Sample articles for missing tags
	sampleArticles := []NewsItem{
		// Hacker News
		{
			ID:           uuid.New().String(),
			Title:        "Show HN: Open Source Agentic Workflow Framework in Go",
			URL:          "https://news.ycombinator.com/item?id=39123456",
			Domain:       "news.ycombinator.com",
			PublishedAt:  now.Add(-30 * time.Minute),
			Category:     "agents",
			Score:        0.92,
			TLDR:         "A lightweight zero-dependency Go library for orchestrating autonomous LLM agent networks",
			WhyItMatters: "Enables production deployment of complex multi-agent workflows with low latency",
			Tags:         []string{"agents", "open-source"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "PostgreSQL 17 Vector Extensions Benchmark for RAG Workloads",
			URL:          "https://news.ycombinator.com/item?id=39124000",
			Domain:       "news.ycombinator.com",
			PublishedAt:  now.Add(-90 * time.Minute),
			Category:     "infra",
			Score:        0.89,
			TLDR:         "Benchmark comparing pgvector HNSW indexing performance against dedicated vector DBs",
			WhyItMatters: "Proves traditional RDBMS can handle production scale vector similarity search",
			Tags:         []string{"infra", "llm"},
		},

		// TechCrunch AI
		{
			ID:           uuid.New().String(),
			Title:        "TechCrunch: Startup Raises $50M for Autonomous Code Refactoring Agents",
			URL:          "https://techcrunch.com/2026/07/30/autonomous-code-agents-series-a",
			Domain:       "techcrunch.com",
			PublishedAt:  now.Add(-1 * time.Hour),
			Category:     "agents",
			Score:        0.88,
			TLDR:         "AI coding startup secures Series A to automate legacy codebase migrations and test generation",
			WhyItMatters: "Accelerates enterprise software maintenance and reduces technical debt dramatically",
			Tags:         []string{"agents", "infra"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "TechCrunch: Next-Gen Open Models Challenge Closed Frontier APIs",
			URL:          "https://techcrunch.com/2026/07/29/open-source-ai-frontier-models",
			Domain:       "techcrunch.com",
			PublishedAt:  now.Add(-4 * time.Hour),
			Category:     "open-source",
			Score:        0.86,
			TLDR:         "New open-weight models achieve parity with proprietary APIs on coding and reasoning benchmarks",
			WhyItMatters: "Shifts competitive advantage toward customized on-premise AI deployments",
			Tags:         []string{"open-source", "llm"},
		},

		// ArXiv AI
		{
			ID:           uuid.New().String(),
			Title:        "ArXiv: Efficient Speculative Decoding for Long-Context Transformer Models",
			URL:          "https://arxiv.org/abs/2026.09876",
			Domain:       "arxiv.org",
			PublishedAt:  now.Add(-2 * time.Hour),
			Category:     "research",
			Score:        0.94,
			TLDR:         "A novel speculative decoding method achieving 3.5x inference speedup without precision loss",
			WhyItMatters: "Crucial optimization for real-time LLM inference serving",
			Tags:         []string{"research", "llm", "infra"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "ArXiv: Multi-Modal Mixture-of-Experts for Real-Time Robot Vision",
			URL:          "https://arxiv.org/abs/2026.05432",
			Domain:       "arxiv.org",
			PublishedAt:  now.Add(-5 * time.Hour),
			Category:     "vision",
			Score:        0.91,
			TLDR:         "Sparse MoE vision architecture designed for low-power edge computing on mobile robots",
			WhyItMatters: "Unlocks spatial understanding for autonomous physical agents",
			Tags:         []string{"vision", "robotics", "research"},
		},

		// GitHub Trending
		{
			ID:           uuid.New().String(),
			Title:        "GitHub Trending: LocalAI - Self-Hosted OpenAI Equivalent REST API",
			URL:          "https://github.com/mudler/LocalAI",
			Domain:       "github.com",
			PublishedAt:  now.Add(-3 * time.Hour),
			Category:     "open-source",
			Score:        0.90,
			TLDR:         "Drop-in replacement REST API for OpenAI specs supporting local model inference",
			WhyItMatters: "Essential tool for local privacy-conscious AI application development",
			Tags:         []string{"open-source", "infra", "llm"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "GitHub Trending: AutoRAG - Automated RAG Pipeline Evaluation Framework",
			URL:          "https://github.com/Marker-IncDB/AutoRAG",
			Domain:       "github.com",
			PublishedAt:  now.Add(-6 * time.Hour),
			Category:     "infra",
			Score:        0.87,
			TLDR:         "Automated evaluation and parameter optimization tool for retrieval-augmented generation",
			WhyItMatters: "Simplifies RAG pipeline tuning and chunking strategy benchmarking",
			Tags:         []string{"infra", "llm"},
		},

		// Reddit MachineLearning
		{
			ID:           uuid.New().String(),
			Title:        "Reddit r/MachineLearning: Discussion on Small Language Models vs Large Model Distillation",
			URL:          "https://reddit.com/r/MachineLearning/comments/slm_vs_distillation",
			Domain:       "reddit.com",
			PublishedAt:  now.Add(-4 * time.Hour),
			Category:     "llm",
			Score:        0.84,
			TLDR:         "Community discussion analyzing trade-offs of 3B-7B parameter distilled models vs large APIs",
			WhyItMatters: "Provides practical deployment insights from engineers scaling local LLM pipelines",
			Tags:         []string{"llm", "open-source"},
		},

		// Vision & AI Labs
		{
			ID:           uuid.New().String(),
			Title:        "Stable Diffusion 3.5 Released with Improved Image Quality",
			URL:          "https://stability.ai/news/stable-diffusion-3",
			Domain:       "stability.ai",
			PublishedAt:  now.Add(-2 * time.Hour),
			Category:     "vision",
			Score:        0.85,
			TLDR:         "Stability AI releases Stable Diffusion 3.5 with better text-to-image generation and prompt adherence",
			WhyItMatters: "Major update to popular open-source image generation model family",
			Tags:         []string{"vision", "open-source"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "DALL-E 3 Now Available in ChatGPT Plus",
			URL:          "https://openai.com/blog/dall-e-3-chatgpt",
			Domain:       "openai.com",
			PublishedAt:  now.Add(-5 * time.Hour),
			Category:     "vision",
			Score:        0.82,
			TLDR:         "OpenAI integrates DALL-E 3 image generation directly into ChatGPT conversational interface",
			WhyItMatters: "Makes advanced vision generation accessible to non-technical users",
			Tags:         []string{"vision", "llm"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "DeepMind's New RL Algorithm Achieves Human-Level Control Performance",
			URL:          "https://deepmind.google/research/rl-breakthrough",
			Domain:       "deepmind.google",
			PublishedAt:  now.Add(-3 * time.Hour),
			Category:     "robotics",
			Score:        0.88,
			TLDR:         "New reinforcement learning approach matches human expert controllers in complex physics tasks",
			WhyItMatters: "Breakthrough in RL could enable more capable autonomous robotics systems",
			Tags:         []string{"robotics", "security"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "Anthropic Offers Free Credits and API Tiers for AI Researchers",
			URL:          "https://anthropic.com/student-program",
			Domain:       "anthropic.com",
			PublishedAt:  now.Add(-1 * time.Hour),
			Category:     "llm",
			Score:        0.80,
			TLDR:         "Researchers can now access Claude 3.5 Sonnet API with upgraded grant allocation",
			WhyItMatters: "Lowers barriers for academic research on model alignment and safety",
			Tags:         []string{"llm", "security"},
		},
	}

	// Read existing if present to keep dataset healthy
	if len(newsData.Items) > 0 {
		// Merge existing items avoiding duplicates by URL
		existingURLs := make(map[string]bool)
		for _, art := range sampleArticles {
			existingURLs[art.URL] = true
		}
		for _, old := range newsData.Items {
			if !existingURLs[old.URL] {
				sampleArticles = append(sampleArticles, old)
			}
		}
	}

	newsData.Items = sampleArticles
	newsData.TotalCount = len(newsData.Items)
	newsData.LastUpdated = now

	// Write to data/news.json
	file, err := os.Create(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to create output file: %v", err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newsData); err != nil {
		file.Close()
		log.Fatalf("❌ Failed to write JSON: %v", err)
	}
	file.Close()

	// Also write to api/news.json
	apiPath := "api/news.json"
	apiFile, err := os.Create(apiPath)
	if err == nil {
		apiEncoder := json.NewEncoder(apiFile)
		apiEncoder.SetIndent("", "  ")
		_ = apiEncoder.Encode(newsData)
		apiFile.Close()
		log.Printf("💾 Synced to %s", apiPath)
	}

	log.Printf("✅ Added %d sample articles across diverse sources", len(newsData.Items))
	log.Printf("\n💾 Updated file: %s", inputPath)
}
