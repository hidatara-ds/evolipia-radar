package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Sample articles with different tags
	samples := []struct {
		title        string
		url          string
		domain       string
		category     string
		tldr         string
		whyItMatters string
		tags         []string
		score        float64
	}{
		{
			title:        "Stable Diffusion 3.5 Released with Better Image Quality",
			url:          "https://stability.ai/news/stable-diffusion-3-5",
			domain:       "stability.ai",
			category:     "tech",
			tldr:         "Stability AI releases major update to image generation model",
			whyItMatters: "Advances in vision models enable better creative tools",
			tags:         []string{"vision", "tools"},
			score:        0.85,
		},
		{
			title:        "DALL-E 3 Integration in ChatGPT",
			url:          "https://openai.com/blog/dall-e-3-integration",
			domain:       "openai.com",
			category:     "tech",
			tldr:         "OpenAI integrates advanced image generation into ChatGPT",
			whyItMatters: "Makes vision AI accessible to millions of users",
			tags:         []string{"vision", "llm"},
			score:        0.82,
		},
		{
			title:        "DeepMind's New RL Algorithm Achieves Human-Level Performance",
			url:          "https://deepmind.google/research/rl-breakthrough-2026",
			domain:       "deepmind.google",
			category:     "tech",
			tldr:         "New reinforcement learning approach matches human experts",
			whyItMatters: "Breakthrough in RL enables more capable autonomous systems",
			tags:         []string{"rl", "research"},
			score:        0.88,
		},
		{
			title:        "PPO Algorithm Improvements for Robot Control",
			url:          "https://arxiv.org/abs/2026.12345",
			domain:       "arxiv.org",
			category:     "tech",
			tldr:         "Researchers improve PPO for better robot manipulation",
			whyItMatters: "Better RL algorithms enable efficient robot training",
			tags:         []string{"rl", "robotics", "research"},
			score:        0.75,
		},
		{
			title:        "Boston Dynamics Atlas Robot Shows New Capabilities",
			url:          "https://bostondynamics.com/blog/atlas-2026",
			domain:       "bostondynamics.com",
			category:     "tech",
			tldr:         "Next-gen Atlas demonstrates advanced manipulation",
			whyItMatters: "Humanoid robots getting closer to practical deployment",
			tags:         []string{"robotics"},
			score:        0.90,
		},
		{
			title:        "Tesla Optimus Robot Update March 2026",
			url:          "https://tesla.com/blog/optimus-march-2026",
			domain:       "tesla.com",
			category:     "tech",
			tldr:         "Tesla's humanoid robot shows improved dexterity",
			whyItMatters: "Progress toward general-purpose household robots",
			tags:         []string{"robotics"},
			score:        0.87,
		},
		{
			title:        "Cursor IDE Adds New AI Features",
			url:          "https://cursor.sh/blog/new-features-2026",
			domain:       "cursor.sh",
			category:     "tech",
			tldr:         "Cursor releases major update with enhanced AI coding",
			whyItMatters: "AI-powered IDEs are transforming developer workflows",
			tags:         []string{"ide", "tools"},
			score:        0.78,
		},
		{
			title:        "GitHub Copilot Workspace Now Available",
			url:          "https://github.blog/copilot-workspace",
			domain:       "github.com",
			category:     "tech",
			tldr:         "GitHub launches AI-powered development environment",
			whyItMatters: "New IDE features make AI coding more accessible",
			tags:         []string{"ide", "tools"},
			score:        0.83,
		},
		{
			title:        "Anthropic Offers $15 Free Credits for Students",
			url:          "https://anthropic.com/student-program-2026",
			domain:       "anthropic.com",
			category:     "tech",
			tldr:         "Students can access Claude API with free credits",
			whyItMatters: "Makes advanced AI accessible to students and educators",
			tags:         []string{"free-credits", "llm"},
			score:        0.80,
		},
		{
			title:        "GitHub Student Pack Adds More AI Tools",
			url:          "https://education.github.com/pack/2026",
			domain:       "education.github.com",
			category:     "tech",
			tldr:         "GitHub expands student benefits with AI coding tools",
			whyItMatters: "Students get free access to premium developer tools",
			tags:         []string{"free-credits", "ide", "tools"},
			score:        0.81,
		},
	}

	// Get or create a default source
	var sourceID uuid.UUID
	err = db.QueryRowContext(ctx, `
		INSERT INTO sources (name, type, category, url, enabled, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (url) DO UPDATE SET updated_at = now()
		RETURNING id
	`, "Sample Data", "manual", "tech", "https://example.com/sample", true, "active").Scan(&sourceID)

	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}

	log.Printf("Using source ID: %s", sourceID)

	// Insert sample articles
	for i, sample := range samples {
		// Create item
		var itemID uuid.UUID
		publishedAt := time.Now().Add(-time.Duration(i) * time.Hour)

		err = db.QueryRowContext(ctx, `
			INSERT INTO items (source_id, title, url, published_at, content_hash, domain, category, raw_excerpt)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (content_hash) DO NOTHING
			RETURNING id
		`, sourceID, sample.title, sample.url, publishedAt,
			fmt.Sprintf("sample-%d", i), sample.domain, sample.category, sample.tldr).Scan(&itemID)

		if err != nil {
			log.Printf("Failed to insert item %s: %v", sample.title, err)
			continue
		}

		// Insert score
		_, err = db.ExecContext(ctx, `
			INSERT INTO scores (item_id, hot, relevance, credibility, novelty, final)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (item_id) DO UPDATE SET
				hot = EXCLUDED.hot,
				relevance = EXCLUDED.relevance,
				credibility = EXCLUDED.credibility,
				novelty = EXCLUDED.novelty,
				final = EXCLUDED.final,
				computed_at = now()
		`, itemID, sample.score, sample.score, sample.score, sample.score, sample.score)

		if err != nil {
			log.Printf("Failed to insert score for %s: %v", sample.title, err)
		}

		// Insert summary with tags
		tagsJSON, _ := json.Marshal(sample.tags)
		_, err = db.ExecContext(ctx, `
			INSERT INTO summaries (item_id, tldr, why_it_matters, tags, method)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (item_id) DO UPDATE SET
				tldr = EXCLUDED.tldr,
				why_it_matters = EXCLUDED.why_it_matters,
				tags = EXCLUDED.tags,
				method = EXCLUDED.method
		`, itemID, sample.tldr, sample.whyItMatters, tagsJSON, "manual")

		if err != nil {
			log.Printf("Failed to insert summary for %s: %v", sample.title, err)
		}

		log.Printf("✅ Inserted: %s (tags: %v)", sample.title, sample.tags)
	}

	log.Println("✅ Sample data population completed!")
}
