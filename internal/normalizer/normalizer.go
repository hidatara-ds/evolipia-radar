package normalizer

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
)

var trackingParams = map[string]bool{
	"utm_source":   true,
	"utm_medium":   true,
	"utm_campaign": true,
	"utm_term":     true,
	"utm_content":  true,
	"ref":          true,
	"source":       true,
	"fbclid":       true,
	"gclid":        true,
	"utm_id":       true,
}

func NormalizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Lowercase host
	parsed.Host = strings.ToLower(parsed.Host)

	// Remove fragment
	parsed.Fragment = ""

	// Remove tracking query params
	query := parsed.Query()
	for key := range query {
		if trackingParams[strings.ToLower(key)] {
			query.Del(key)
		}
	}
	parsed.RawQuery = query.Encode()

	// Normalize trailing slash (remove for consistency)
	path := parsed.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		parsed.Path = strings.TrimSuffix(path, "/")
	}

	return parsed.String(), nil
}

func NormalizeDomain(hostname string) string {
	hostname = strings.ToLower(hostname)
	// Remove port if present
	if idx := strings.Index(hostname, ":"); idx != -1 {
		hostname = hostname[:idx]
	}
	// Remove www. prefix for consistency
	hostname = strings.TrimPrefix(hostname, "www.")
	return hostname
}

func ContentHash(title, url string) string {
	normalizedURL, err := NormalizeURL(url)
	if err != nil {
		normalizedURL = url
	}

	normalizedTitle := strings.ToLower(strings.TrimSpace(title))
	input := normalizedTitle + "|" + normalizedURL

	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
