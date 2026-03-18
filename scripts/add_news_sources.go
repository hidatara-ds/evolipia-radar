// Script to add multiple AI/ML news sources to the database
// Run with: go run scripts/add_news_sources.go
package main

import (
	"context"
	"log"

	"github.com/hidatara-ds/evolipia-radar/pkg/config"
	"github.com/hidatara-ds/evolipia-radar/pkg/db"
	"github.com/hidatara-ds/evolipia-radar/pkg/models"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	sourceRepo := db.NewSourceRepository(database)
	ctx := context.Background()

	// Define all news sources
	newsSources := []models.Source{
		// Existing sources
		{
			Name:     "Hacker News",
			Type:     "hacker_news",
			Category: "news",
			URL:      "https://news.ycombinator.com",
			Enabled:  true,
			Status:   "active",
		},
		
		// arXiv Research Papers
		{
			Name:     "arXiv AI Papers",
			Type:     "arxiv",
			Category: "research",
			URL:      "http://export.arxiv.org/api/query?search_query=cat:cs.AI&sortBy=submittedDate&sortOrder=descending&max_results=50",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "arXiv Machine Learning",
			Type:     "arxiv",
			Category: "research",
			URL:      "http://export.arxiv.org/api/query?search_query=cat:cs.LG&sortBy=submittedDate&sortOrder=descending&max_results=50",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "arXiv Computer Vision",
			Type:     "arxiv",
			Category: "research",
			URL:      "http://export.arxiv.org/api/query?search_query=cat:cs.CV&sortBy=submittedDate&sortOrder=descending&max_results=50",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "arXiv NLP",
			Type:     "arxiv",
			Category: "research",
			URL:      "http://export.arxiv.org/api/query?search_query=cat:cs.CL&sortBy=submittedDate&sortOrder=descending&max_results=50",
			Enabled:  true,
			Status:   "active",
		},
		
		// Company Blogs
		{
			Name:     "HuggingFace Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://huggingface.co/blog/feed.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Google AI Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://blog.research.google/feeds/posts/default",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "OpenAI News",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://openai.com/news/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Anthropic News",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://www.anthropic.com/news/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "DeepMind Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://deepmind.google/blog/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		
		// Tech News Sites
		{
			Name:     "MIT Technology Review AI",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://www.technologyreview.com/topic/artificial-intelligence/feed",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "VentureBeat AI",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://venturebeat.com/category/ai/feed/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "The Verge AI",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://www.theverge.com/ai-artificial-intelligence/rss/index.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "TechCrunch AI",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://techcrunch.com/category/artificial-intelligence/feed/",
			Enabled:  true,
			Status:   "active",
		},
		
		// Research & Analysis
		{
			Name:     "The Gradient",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://thegradient.pub/rss/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Distill.pub",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://distill.pub/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "AI Alignment Forum",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://www.alignmentforum.org/feed.xml",
			Enabled:  true,
			Status:   "active",
		},
		
		// Community & Aggregators
		{
			Name:     "Papers with Code",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://paperswithcode.com/latest/rss",
			Enabled:  true,
			Status:   "active",
		},
		
		// IDE & Developer Tools (NEW!)
		{
			Name:     "GitHub Blog - Copilot",
			Type:     "rss_atom",
			Category: "tools",
			URL:      "https://github.blog/tag/github-copilot/feed/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Cursor Changelog",
			Type:     "rss_atom",
			Category: "tools",
			URL:      "https://changelog.cursor.sh/rss",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "JetBrains Blog - AI",
			Type:     "rss_atom",
			Category: "tools",
			URL:      "https://blog.jetbrains.com/feed/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Replit Blog",
			Type:     "rss_atom",
			Category: "tools",
			URL:      "https://blog.replit.com/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Codeium Blog",
			Type:     "rss_atom",
			Category: "tools",
			URL:      "https://codeium.com/blog/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		
		// Free Credits & Student Programs (NEW!)
		{
			Name:     "GitHub Education Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://github.blog/category/education/feed/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Dev.to - Free Resources",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://dev.to/feed/tag/free",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Indie Hackers",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://www.indiehackers.com/feed",
			Enabled:  true,
			Status:   "active",
		},
		
		// Computer Vision & Image Generation (NEW!)
		{
			Name:     "Stability AI Blog",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://stability.ai/news/rss",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "Midjourney News",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://www.midjourney.com/feed",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "RunwayML Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://runwayml.com/blog/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		
		// Robotics & RL (NEW!)
		{
			Name:     "Boston Dynamics Blog",
			Type:     "rss_atom",
			Category: "news",
			URL:      "https://bostondynamics.com/blog/feed/",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "OpenAI Robotics",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://openai.com/research/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
		{
			Name:     "DeepMind Research",
			Type:     "rss_atom",
			Category: "research",
			URL:      "https://deepmind.google/discover/blog/rss.xml",
			Enabled:  true,
			Status:   "active",
		},
	}

	// Insert sources with upsert logic (skip if already exists)
	added := 0
	skipped := 0
	failed := 0

	for _, source := range newsSources {
		// Check if source already exists by URL
		sources, err := sourceRepo.List(ctx)
		if err != nil {
			log.Printf("❌ Error listing sources: %v", err)
			failed++
			continue
		}

		exists := false
		for _, s := range sources {
			if s.URL == source.URL {
				exists = true
				log.Printf("⏭️  Source already exists: %s (%s)", source.Name, source.URL)
				skipped++
				break
			}
		}

		if !exists {
			if err := sourceRepo.Create(ctx, &source); err != nil {
				log.Printf("❌ Error creating source %s: %v", source.Name, err)
				failed++
			} else {
				log.Printf("✅ Added source: %s (%s)", source.Name, source.Type)
				added++
			}
		}
	}

	log.Println("\n" + strings.Repeat("=", 60))
	log.Printf("📊 Summary:")
	log.Printf("   ✅ Added: %d sources", added)
	log.Printf("   ⏭️  Skipped (already exists): %d sources", skipped)
	log.Printf("   ❌ Failed: %d sources", failed)
	log.Printf("   📦 Total sources configured: %d", len(newsSources))
	log.Println(strings.Repeat("=", 60))
}

// Helper function to repeat string
func (s string) repeat(count int) string {
	return strings.Repeat(s, count)
}
