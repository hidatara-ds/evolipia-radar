package ai

import (
	"fmt"
	"net/url"
	"strings"
)

// PreProcessFilter checks an incoming article for noise, low quality, and blacklists.
// Returns an error if the article should be dropped.
func PreProcessFilter(urlStr, title, content string) error {
	// 1. Minimum Content Length
	// If an article is too short (e.g., just a tweet or link-farm), we don't cluster it.
	wordCount := len(strings.Fields(content))
	if wordCount < 50 {
		return fmt.Errorf("dropped: content too short (%d words)", wordCount)
	}

	// 2. Simple Keyword Blacklist (Noise words that often pollute AI news)
	blacklist := []string{"crypto", "nft", "token", "bitcoin", "doge", "memecoin", "gamble", "casino"}
	contentLower := strings.ToLower(content)
	titleLower := strings.ToLower(title)

	for _, badWord := range blacklist {
		if strings.Contains(contentLower, badWord) || strings.Contains(titleLower, badWord) {
			return fmt.Errorf("dropped: hits blacklist keyword '%s'", badWord)
		}
	}

	// 3. Basic Domain Quality Scoring
	// Check if the URL belongs to a known spam or low-quality farm.
	parsedURL, err := url.Parse(urlStr)
	if err == nil && parsedURL.Host != "" {
		host := strings.ToLower(parsedURL.Host)
		domainBlacklist := []string{
			"medium.com", // often too noisy, unless heavily curated
			"buzzfeed.com",
			"slashdot.org", // too broad
			"prweb.com",    // press releases
			"businesswire", // pure PR
		}

		for _, badDomain := range domainBlacklist {
			if strings.Contains(host, badDomain) {
				return fmt.Errorf("dropped: domain %s is on the quality blacklist", host)
			}
		}
	}

	return nil
}
