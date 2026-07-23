// Package utils provides common utility functions for hashing, text processing, and validation.
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

// HashString generates a SHA-256 hex hash of the given string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(h.Sum(nil))
}

// NormalizeTitle trims whitespace, lowers case, and removes excess punctuation for deduplication.
func NormalizeTitle(title string) string {
	lowered := strings.ToLower(strings.TrimSpace(title))
	var sb strings.Builder
	for _, r := range lowered {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			sb.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(sb.String()), " ")
}

// IsValidURL verifies if the string is a valid HTTP or HTTPS URL.
func IsValidURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
