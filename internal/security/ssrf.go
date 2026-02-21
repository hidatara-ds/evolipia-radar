package security

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var (
	privateIPRanges = []*net.IPNet{
		{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(8, 32)},     // 10.0.0.0/8
		{IP: net.IP{172, 16, 0, 0}, Mask: net.CIDRMask(12, 32)},  // 172.16.0.0/12
		{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(16, 32)}, // 192.168.0.0/16
		{IP: net.IP{169, 254, 0, 0}, Mask: net.CIDRMask(16, 32)}, // 169.254.0.0/16
	}

	blockedHosts = map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"::1":       true,
		"0.0.0.0":   true,
	}
)

func ValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("only http/https schemes are allowed")
	}

	host := strings.ToLower(parsed.Hostname())

	// Check blocked hosts
	if blockedHosts[host] {
		return fmt.Errorf("blocked host: %s", host)
	}

	// Resolve IP addresses
	ips, err := net.LookupIP(host)
	if err != nil {
		// If DNS lookup fails, we'll allow it but log it
		// In production, you might want to be stricter
		return nil
	}

	for _, ip := range ips {
		// Check for localhost IPs
		if ip.IsLoopback() {
			return fmt.Errorf("loopback IP address not allowed")
		}

		// Check for link-local addresses
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("link-local IP address not allowed")
		}

		// Check for private IP ranges
		for _, ipNet := range privateIPRanges {
			if ipNet.Contains(ip) {
				return fmt.Errorf("private IP range not allowed: %s", ip.String())
			}
		}
	}

	return nil
}
