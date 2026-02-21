package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/dto"
	"github.com/hidatara-ds/evolipia-radar/internal/normalizer"
)

var (
	ErrTimeout        = errors.New("request timeout")
	ErrSizeLimit      = errors.New("response size limit exceeded")
	ErrInvalidURL     = errors.New("invalid outbound url")
	ErrDisallowedHost = errors.New("host not allowed")
)

// Optional allowlist (recommended): comma-separated hosts/domains.
// Examples:
//   EVOLIPIA_ALLOWED_FETCH_HOSTS="kompas.com,tempo.co,cnnindonesia.com,antaranews.com,.googleapis.com"
// Rules:
// - "example.com" allows example.com + subdomains (*.example.com)
// - ".example.com" also allows subdomains (suffix match)
func allowedFetchHostsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("EVOLIPIA_ALLOWED_FETCH_HOSTS"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(p))
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func fetchWithLimits(ctx context.Context, rawURL string, cfg *config.Config) ([]byte, error) {
	u, err := validateOutboundURL(ctx, rawURL, allowedFetchHostsFromEnv())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "evolipia-radar/1.0")

	client := newSafeHTTPClient(cfg)

	resp, err := client.Do(req)
	if err != nil {
		// keep your old timeout behavior
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Limit response size
	limitedReader := io.LimitReader(resp.Body, cfg.MaxFetchBytes)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(body) >= int(cfg.MaxFetchBytes) {
		return nil, ErrSizeLimit
	}

	return body, nil
}

// Disable redirects so attacker can't redirect from public URL -> internal URL.
func newSafeHTTPClient(cfg *config.Config) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// Optional extra hardening
	transport.ResponseHeaderTimeout = cfg.FetchTimeout()
	transport.TLSHandshakeTimeout = 10 * time.Second

	return &http.Client{
		Timeout:   cfg.FetchTimeout(),
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func validateOutboundURL(ctx context.Context, raw string, allowlist []string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, ErrInvalidURL
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: parse failed", ErrInvalidURL)
	}

	// Require absolute URL (scheme + host)
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("%w: url must be absolute", ErrInvalidURL)
	}

	// Recommend https only
	if u.Scheme != "https" {
		return nil, fmt.Errorf("%w: scheme not allowed", ErrInvalidURL)
	}

	// No credentials in URL
	if u.User != nil {
		return nil, fmt.Errorf("%w: userinfo not allowed", ErrInvalidURL)
	}

	host := strings.ToLower(u.Hostname())
	if host == "" {
		return nil, fmt.Errorf("%w: missing host", ErrInvalidURL)
	}

	// Block local hostnames
	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return nil, fmt.Errorf("%w: local hostnames blocked", ErrInvalidURL)
	}

	// Allowlist (if configured)
	if len(allowlist) > 0 && !hostAllowed(host, allowlist) {
		return nil, fmt.Errorf("%w: %s", ErrDisallowedHost, host)
	}

	// DNS resolve and block private/local/link-local/etc
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("%w: dns lookup failed", ErrInvalidURL)
	}
	for _, ip := range ips {
		if isPrivateOrLocalIP(ip.IP) {
			return nil, fmt.Errorf("%w: private/local ip blocked (%s)", ErrInvalidURL, ip.IP.String())
		}
	}

	return u, nil
}

func hostAllowed(host string, allowed []string) bool {
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" {
			continue
		}
		// "example.com" => allow example.com and subdomains
		if host == a || strings.HasSuffix(host, "."+a) {
			return true
		}
		// ".example.com" => suffix match (subdomains)
		if strings.HasPrefix(a, ".") && strings.HasSuffix(host, a) {
			return true
		}
	}
	return false
}

func isPrivateOrLocalIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	// Normalize to 16 bytes
	ip16 := ip.To16()
	if ip16 == nil {
		return true
	}

	// Built-in checks: loopback/link-local/multicast/unspecified/private
	if ip16.IsLoopback() ||
		ip16.IsLinkLocalUnicast() ||
		ip16.IsLinkLocalMulticast() ||
		ip16.IsMulticast() ||
		ip16.IsUnspecified() ||
		ip16.IsPrivate() {
		return true
	}

	// Extra IPv4 block: CGNAT 100.64.0.0/10
	if v4 := ip.To4(); v4 != nil {
		_, cgnat, _ := net.ParseCIDR("100.64.0.0/10")
		if cgnat.Contains(v4) {
			return true
		}
	}

	return false
}

func FetchRSSAtom(ctx context.Context, feedURL string, cfg *config.Config) ([]dto.ContentItem, error) {
	body, err := fetchWithLimits(ctx, feedURL, cfg)
	if err != nil {
		return nil, err
	}

	// Simple RSS/Atom parser (MVP - can be enhanced)
	items := parseRSSAtom(body)
	return items, nil
}

func parseRSSAtom(body []byte) []dto.ContentItem {
	// MVP: Simple XML parsing
	// In production, use a proper RSS/Atom library
	content := string(body)
	var items []dto.ContentItem

	// Very basic RSS parsing (item tags)
	itemStart := "<item>"
	itemEnd := "</item>"

	// Also handle Atom entries
	entryStart := "<entry>"
	entryEnd := "</entry>"

	parseItem := func(itemContent string) dto.ContentItem {
		item := dto.ContentItem{
			Category: "news",
			Tags:     []string{},
		}

		// Extract title
		if titleStart := strings.Index(itemContent, "<title>"); titleStart != -1 {
			titleStart += len("<title>")
			if titleEnd := strings.Index(itemContent[titleStart:], "</title>"); titleEnd != -1 {
				item.Title = strings.TrimSpace(itemContent[titleStart : titleStart+titleEnd])
			}
		}

		// Extract link
		if linkStart := strings.Index(itemContent, "<link>"); linkStart != -1 {
			linkStart += len("<link>")
			if linkEnd := strings.Index(itemContent[linkStart:], "</link>"); linkEnd != -1 {
				item.URL = strings.TrimSpace(itemContent[linkStart : linkStart+linkEnd])
			}
		} else if linkStart := strings.Index(itemContent, `href="`); linkStart != -1 {
			linkStart += len(`href="`)
			if linkEnd := strings.Index(itemContent[linkStart:], `"`); linkEnd != -1 {
				item.URL = strings.TrimSpace(itemContent[linkStart : linkStart+linkEnd])
			}
		}

		// Extract pubDate or published
		if pubStart := strings.Index(itemContent, "<pubDate>"); pubStart != -1 {
			pubStart += len("<pubDate>")
			if pubEnd := strings.Index(itemContent[pubStart:], "</pubDate>"); pubEnd != -1 {
				dateStr := strings.TrimSpace(itemContent[pubStart : pubStart+pubEnd])
				if t, err := time.Parse(time.RFC1123Z, dateStr); err == nil {
					item.PublishedAt = t
				} else if t, err := time.Parse(time.RFC1123, dateStr); err == nil {
					item.PublishedAt = t
				}
			}
		} else if pubStart := strings.Index(itemContent, "<published>"); pubStart != -1 {
			pubStart += len("<published>")
			if pubEnd := strings.Index(itemContent[pubStart:], "</published>"); pubEnd != -1 {
				dateStr := strings.TrimSpace(itemContent[pubStart : pubStart+pubEnd])
				if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
					item.PublishedAt = t
				}
			}
		}

		// Extract description/summary
		if descStart := strings.Index(itemContent, "<description>"); descStart != -1 {
			descStart += len("<description>")
			if descEnd := strings.Index(itemContent[descStart:], "</description>"); descEnd != -1 {
				item.Excerpt = strings.TrimSpace(itemContent[descStart : descStart+descEnd])
			}
		} else if sumStart := strings.Index(itemContent, "<summary>"); sumStart != -1 {
			sumStart += len("<summary>")
			if sumEnd := strings.Index(itemContent[sumStart:], "</summary>"); sumEnd != -1 {
				item.Excerpt = strings.TrimSpace(itemContent[sumStart : sumStart+sumEnd])
			}
		}

		if item.PublishedAt.IsZero() {
			item.PublishedAt = time.Now()
		}

		if item.URL != "" {
			parsedURL, err := url.Parse(item.URL)
			if err == nil {
				item.Domain = normalizer.NormalizeDomain(parsedURL.Hostname())
			}
		}

		return item
	}

	// Parse RSS items
	start := 0
	for {
		idx := strings.Index(content[start:], itemStart)
		if idx == -1 {
			break
		}
		itemStartIdx := start + idx
		endIdx := strings.Index(content[itemStartIdx:], itemEnd)
		if endIdx == -1 {
			break
		}
		itemContent := content[itemStartIdx : itemStartIdx+endIdx+len(itemEnd)]
		item := parseItem(itemContent)
		if item.Title != "" && item.URL != "" {
			items = append(items, item)
		}
		start = itemStartIdx + endIdx + len(itemEnd)
	}

	// Parse Atom entries if no RSS items found
	if len(items) == 0 {
		start = 0
		for {
			idx := strings.Index(content[start:], entryStart)
			if idx == -1 {
				break
			}
			entryStartIdx := start + idx
			endIdx := strings.Index(content[entryStartIdx:], entryEnd)
			if endIdx == -1 {
				break
			}
			entryContent := content[entryStartIdx : entryStartIdx+endIdx+len(entryEnd)]
			item := parseItem(entryContent)
			if item.Title != "" && item.URL != "" {
				items = append(items, item)
			}
			start = entryStartIdx + endIdx + len(entryEnd)
		}
	}

	return items
}

func FetchJSONAPI(ctx context.Context, apiURL string, mapping map[string]interface{}, cfg *config.Config) ([]dto.ContentItem, error) {
	body, err := fetchWithLimits(ctx, apiURL, cfg)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	itemsPath := getString(mapping, "items_path", "items")
	itemsArray := getNestedValue(data, itemsPath)
	if itemsArray == nil {
		return nil, fmt.Errorf("items array not found at path: %s", itemsPath)
	}

	itemsSlice, ok := itemsArray.([]interface{})
	if !ok {
		return nil, fmt.Errorf("items_path does not point to an array")
	}

	titlePath := getString(mapping, "title_path", "title")
	urlPath := getString(mapping, "url_path", "url")
	publishedAtPath := getString(mapping, "published_at_path", "published_at")
	summaryPath := getString(mapping, "summary_path", "")

	var items []dto.ContentItem
	for _, itemRaw := range itemsSlice {
		itemMap, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}

		item := dto.ContentItem{
			Category: "news",
			Tags:     []string{},
		}

		item.Title = getStringValue(itemMap, titlePath)
		item.URL = getStringValue(itemMap, urlPath)

		if item.Title == "" || item.URL == "" {
			continue
		}

		// Parse published_at
		dateStr := getStringValue(itemMap, publishedAtPath)
		if dateStr != "" {
			if t, err := parseDate(dateStr); err == nil {
				item.PublishedAt = t
			}
		}
		if item.PublishedAt.IsZero() {
			item.PublishedAt = time.Now()
		}

		if summaryPath != "" {
			item.Excerpt = getStringValue(itemMap, summaryPath)
		}

		parsedURL, err := url.Parse(item.URL)
		if err == nil {
			item.Domain = normalizer.NormalizeDomain(parsedURL.Hostname())
		}

		items = append(items, item)
	}

	return items, nil
}

func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		if part == "" {
			continue
		}
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			return nil
		}
	}
	return current
}

func getStringValue(m map[string]interface{}, path string) string {
	val := getNestedValue(m, path)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
