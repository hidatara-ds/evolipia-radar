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
		// Vision articles
		{
			ID:           uuid.New().String(),
			Title:        "Stable Diffusion 3.0 Released with Improved Image Quality",
			URL:          "https://stability.ai/news/stable-diffusion-3",
			Domain:       "stability.ai",
			PublishedAt:  now.Add(-2 * time.Hour),
			Category:     "tech",
			Score:        0.85,
			TLDR:         "Stability AI releases Stable Diffusion 3.0 with better text-to-image generation",
			WhyItMatters: "Major update to popular open-source image generation model",
			Tags:         []string{"vision", "tools"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "DALL-E 3 Now Available in ChatGPT Plus",
			URL:          "https://openai.com/blog/dall-e-3-chatgpt",
			Domain:       "openai.com",
			PublishedAt:  now.Add(-5 * time.Hour),
			Category:     "tech",
			Score:        0.82,
			TLDR:         "OpenAI integrates DALL-E 3 image generation into ChatGPT",
			WhyItMatters: "Makes advanced image generation accessible to millions of users",
			Tags:         []string{"vision", "llm"},
		},
		// RL articles
		{
			ID:           uuid.New().String(),
			Title:        "DeepMind's New RL Algorithm Achieves Human-Level Performance",
			URL:          "https://deepmind.google/research/rl-breakthrough",
			Domain:       "deepmind.google",
			PublishedAt:  now.Add(-3 * time.Hour),
			Category:     "tech",
			Score:        0.88,
			TLDR:         "New reinforcement learning approach matches human experts in complex tasks",
			WhyItMatters: "Breakthrough in RL could enable more capable autonomous systems",
			Tags:         []string{"rl", "research"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "PPO Algorithm Improvements for Robotics Control",
			URL:          "https://arxiv.org/abs/2024.12345",
			Domain:       "arxiv.org",
			PublishedAt:  now.Add(-6 * time.Hour),
			Category:     "tech",
			Score:        0.75,
			TLDR:         "Researchers improve Proximal Policy Optimization for robot manipulation",
			WhyItMatters: "Better RL algorithms enable more efficient robot training",
			Tags:         []string{"rl", "robotics", "research"},
		},
		// Robotics articles
		{
			ID:           uuid.New().String(),
			Title:        "Boston Dynamics Unveils New Humanoid Robot Atlas",
			URL:          "https://bostondynamics.com/blog/atlas-next-gen",
			Domain:       "bostondynamics.com",
			PublishedAt:  now.Add(-4 * time.Hour),
			Category:     "tech",
			Score:        0.90,
			TLDR:         "Next-generation Atlas robot demonstrates advanced manipulation capabilities",
			WhyItMatters: "Humanoid robots getting closer to practical real-world deployment",
			Tags:         []string{"robotics"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "Tesla Optimus Robot Shows Improved Dexterity",
			URL:          "https://tesla.com/blog/optimus-update",
			Domain:       "tesla.com",
			PublishedAt:  now.Add(-7 * time.Hour),
			Category:     "tech",
			Score:        0.87,
			TLDR:         "Tesla's humanoid robot demonstrates fine motor control improvements",
			WhyItMatters: "Progress toward general-purpose household robots",
			Tags:         []string{"robotics"},
		},
		// Free credits articles
		{
			ID:           uuid.New().String(),
			Title:        "Anthropic Offers $10 Free Credits for Students",
			URL:          "https://anthropic.com/student-program",
			Domain:       "anthropic.com",
			PublishedAt:  now.Add(-1 * time.Hour),
			Category:     "tech",
			Score:        0.78,
			TLDR:         "Students can now access Claude API with $10 in free credits",
			WhyItMatters: "Makes advanced AI accessible to students and educators",
			Tags:         []string{"free-credits", "llm"},
		},
		{
			ID:           uuid.New().String(),
			Title:        "GitHub Student Developer Pack Adds New AI Tools",
			URL:          "https://education.github.com/pack",
			Domain:       "education.github.com",
			PublishedAt:  now.Add(-8 * time.Hour),
			Category:     "tech",
			Score:        0.80,
			TLDR:         "GitHub expands student benefits with free access to AI coding tools",
			WhyItMatters: "Students get free access to premium developer tools and AI assistants",
			Tags:         []string{"free-credits", "ide", "tools"},
		},
	}

	// Add sample articles to the beginning (highest scores)
	newsData.Items = append(sampleArticles, newsData.Items...)
	newsData.TotalCount = len(newsData.Items)
	newsData.LastUpdated = now

	// Keep only top 100
	if len(newsData.Items) > 100 {
		newsData.Items = newsData.Items[:100]
		newsData.TotalCount = 100
	}

	// Write back to file
	file, err := os.Create(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newsData); err != nil {
		log.Fatalf("❌ Failed to write JSON: %v", err)
	}

	log.Printf("✅ Added %d sample articles", len(sampleArticles))
	log.Println("📊 Sample articles by tag:")
	log.Println("   - vision: 2 articles")
	log.Println("   - rl: 2 articles")
	log.Println("   - robotics: 3 articles")
	log.Println("   - free-credits: 2 articles")
	log.Printf("\n💾 Updated file: %s", inputPath)
}
