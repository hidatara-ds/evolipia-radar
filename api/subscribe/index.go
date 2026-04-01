// Package handler is the serverless function for the /api/subscribe endpoint.
// It validates the email, applies per-IP rate limiting (3 req/hr),
// and forwards the subscription to your email provider via a secret API key
// stored in environment variables — never exposed to the client bundle.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

// ── Rate limiter ────────────────────────────────────────────────────────────

type rateBucket struct {
	count     int
	windowEnd time.Time
}

var (
	mu      sync.Mutex
	buckets = make(map[string]*rateBucket)
)

const (
	rateLimit  = 3         // max requests
	rateWindow = time.Hour // per window
)

func isRateLimited(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	b, ok := buckets[ip]
	if !ok || now.After(b.windowEnd) {
		buckets[ip] = &rateBucket{count: 1, windowEnd: now.Add(rateWindow)}
		return false
	}
	if b.count >= rateLimit {
		return true
	}
	b.count++
	return false
}

// ── Email validation ────────────────────────────────────────────────────────

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidEmail(email string) bool {
	return len(email) <= 320 && emailRegex.MatchString(email)
}

// ── Response helpers ────────────────────────────────────────────────────────

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

// ── Handler ──────────────────────────────────────────────────────────────────

// Handler is the Vercel serverless entry point.
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS preflight
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Rate limiting (use X-Forwarded-For behind Vercel's proxy)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if isRateLimited(ip) {
		jsonError(w, http.StatusTooManyRequests, "Too many requests. Please wait before trying again.")
		return
	}

	// Parse body
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !isValidEmail(body.Email) {
		jsonError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	// ── Connect your email provider here ────────────────────────────────────
	// The API key is read from an environment variable — never bundled in client JS.
	//
	// Example using Resend (https://resend.com):
	//   apiKey := os.Getenv("RESEND_API_KEY")
	//   payload := fmt.Sprintf(`{"from":"Evolipia Radar <noreply@evolipia.com>","to":[%q],"subject":"Welcome!","html":"<p>Thanks for subscribing!</p>"}`, body.Email)
	//   req, _ := http.NewRequest("POST", "https://api.resend.com/emails", strings.NewReader(payload))
	//   req.Header.Set("Authorization", "Bearer "+apiKey)
	//   req.Header.Set("Content-Type", "application/json")
	//   client := &http.Client{Timeout: 10 * time.Second}
	//   resp, err := client.Do(req)
	//   if err != nil || resp.StatusCode >= 400 { ... handle error ... }
	//
	// Until you wire up a provider, log the subscription:
	log.Printf("[subscribe] new subscriber: %s", body.Email)
	// ────────────────────────────────────────────────────────────────────────

	jsonOK(w, "Subscribed successfully!")
}
