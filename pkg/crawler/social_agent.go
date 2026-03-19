package crawler

import (
	"context"
	"log"

	"github.com/hidatara-ds/evolipia-radar/pkg/db"
)

type SocialAgent struct {
	Platform string
	repo     *db.SettingRepository
}

func NewSocialAgent(platform string, pool *db.DB) *SocialAgent {
	return &SocialAgent{
		Platform: platform,
		repo:     db.NewSettingRepository(pool),
	}
}

func (a *SocialAgent) Name() string {
	return "SocialAgent-" + a.Platform
}

func (a *SocialAgent) Crawl(ctx context.Context, maxItems int) ([]Article, error) {
	// Check if API key exists in settings
	keyName := ""
	switch a.Platform {
	case "X":
		keyName = "x_api_key"
	case "Threads":
		keyName = "threads_api_key"
	}

	if keyName == "" {
		return nil, nil
	}

	apiKey, err := a.repo.Get(ctx, keyName)
	if err != nil || apiKey == "" {
		log.Printf("[SOCIAL] Skipping %s agent: No API key found in settings.", a.Platform)
		return nil, nil
	}

	log.Printf("[SOCIAL] %s agent active with provided API key. (Placeholder fetching logic)", a.Platform)
	
	// Real implementation would use the apiKey to fetch from X/Threads API
	// For now, we'll return empty list but log that it's "authorized"
	
	return []Article{}, nil
}
