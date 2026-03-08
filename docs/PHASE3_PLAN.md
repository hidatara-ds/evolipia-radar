# Phase 3 Plan - Admin Dashboard, Personalization & Mobile Apps

## Overview
Phase 3 focuses on polish, advanced features, and production readiness with admin tools, personalization, and mobile app wrappers.

## 1. Admin Dashboard

### Features
- **Source Management**
  - View all sources with health metrics
  - Enable/disable sources
  - Test connections
  - View fetch history and errors
  - Add/edit/delete sources via UI

- **Content Moderation**
  - Review flagged items
  - Manual scoring adjustments
  - Blacklist domains/keywords
  - Whitelist trusted sources

- **Analytics Dashboard**
  - Total items ingested (daily/weekly/monthly)
  - Top sources by volume
  - Average scores by category
  - User engagement metrics (if tracking enabled)
  - API usage statistics

- **Scoring Algorithm Tuning**
  - Adjust weights (w1, w2, w3, w4) via UI
  - Preview score changes
  - A/B test different algorithms
  - Export/import configurations

- **System Health**
  - Worker status and last run time
  - Database connection pool stats
  - API response times
  - Error rates and logs
  - Memory/CPU usage

### Tech Stack
- **Backend**: New admin API endpoints
- **Frontend**: React or Vue.js SPA
- **Auth**: JWT-based authentication
- **Charts**: Chart.js or Recharts

### Implementation Plan

#### Step 1: Admin API Endpoints
```go
// internal/http/handlers/admin.go
type AdminHandlers struct {
    db *db.DB
}

// GET /admin/sources - List all sources with metrics
func (h *AdminHandlers) ListSourcesWithMetrics(c *gin.Context) {
    // Return sources with:
    // - Last fetch time
    // - Success rate
    // - Items fetched (last 24h)
    // - Average latency
}

// GET /admin/analytics/overview
func (h *AdminHandlers) GetAnalyticsOverview(c *gin.Context) {
    // Return:
    // - Total items
    // - Items by category
    // - Top domains
    // - Score distribution
}

// POST /admin/scoring/weights
func (h *AdminHandlers) UpdateScoringWeights(c *gin.Context) {
    // Update w1, w2, w3, w4
    // Trigger re-scoring
}

// GET /admin/logs?level=error&limit=100
func (h *AdminHandlers) GetLogs(c *gin.Context) {
    // Return recent logs
}
```

#### Step 2: Admin UI
```bash
# Create admin UI directory
mkdir -p web/admin

# Use Vite + React
cd web/admin
npm create vite@latest . -- --template react
npm install recharts axios react-router-dom

# Build admin UI
npm run build

# Serve from Go
router.Static("/admin", "./web/admin/dist")
```

#### Step 3: Authentication
```go
// internal/auth/jwt.go
type AuthService struct {
    secretKey []byte
}

func (a *AuthService) GenerateToken(userID string) (string, error) {
    // Generate JWT token
}

func (a *AuthService) ValidateToken(token string) (*Claims, error) {
    // Validate JWT token
}

// Middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // Validate token
        // Set user in context
    }
}
```

### Database Schema
```sql
-- Admin users table
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP
);

-- Audit log
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES admin_users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Scoring configurations
CREATE TABLE scoring_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    w1 FLOAT NOT NULL,
    w2 FLOAT NOT NULL,
    w3 FLOAT NOT NULL,
    w4 FLOAT NOT NULL,
    active BOOLEAN DEFAULT false,
    created_by UUID REFERENCES admin_users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## 2. Personalization Engine

### Features
- **User Profiles**
  - Track reading history (optional, privacy-first)
  - Preferred topics
  - Preferred sources
  - Reading time preferences

- **Personalized Feed**
  - Boost items matching user interests
  - Learn from clicks and bookmarks
  - Collaborative filtering (users with similar interests)
  - Diversity injection (avoid filter bubbles)

- **Smart Recommendations**
  - "You might also like" based on reading history
  - Topic-based recommendations
  - Author-based recommendations
  - Time-based recommendations (morning vs evening)

### Implementation Plan

#### Step 1: User Tracking (Privacy-First)
```go
// internal/models/user.go
type User struct {
    ID              uuid.UUID
    AnonymousID     string  // No PII, just random ID
    Preferences     Preferences
    CreatedAt       time.Time
}

type Preferences struct {
    Topics          []string  // ["llm", "cv", "mlops"]
    Sources         []string  // Preferred source IDs
    ReadingTime     string    // "morning", "evening", "anytime"
    ContentLength   string    // "short", "medium", "long"
}

// internal/models/interaction.go
type Interaction struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    ItemID      uuid.UUID
    Type        string    // "view", "click", "bookmark", "share"
    Duration    int       // Seconds spent reading
    CreatedAt   time.Time
}
```

#### Step 2: Personalization Algorithm
```go
// internal/personalization/engine.go
type Engine struct {
    db *db.DB
}

func (e *Engine) PersonalizeScore(userID uuid.UUID, item *models.Item, baseScore float64) float64 {
    // Get user preferences
    prefs := e.getUserPreferences(userID)
    
    // Boost if matches preferred topics
    topicBoost := e.calculateTopicBoost(item, prefs.Topics)
    
    // Boost if from preferred sources
    sourceBoost := e.calculateSourceBoost(item, prefs.Sources)
    
    // Boost if similar to previously liked items
    similarityBoost := e.calculateSimilarityBoost(userID, item)
    
    // Combine boosts
    personalizedScore := baseScore * (1 + topicBoost + sourceBoost + similarityBoost)
    
    return personalizedScore
}

func (e *Engine) GetRecommendations(userID uuid.UUID, limit int) ([]models.Item, error) {
    // Collaborative filtering
    // Find users with similar interests
    // Recommend items they liked
}
```

#### Step 3: API Endpoints
```go
// GET /v1/feed/personalized?user_id=xxx
// GET /v1/recommendations?user_id=xxx&limit=10
// POST /v1/interactions - Track user interactions
// GET /v1/users/:id/preferences
// PUT /v1/users/:id/preferences
```

### Privacy Considerations
- **Anonymous by Default**: Use random UUIDs, no email/name required
- **Opt-in Tracking**: Users must explicitly enable personalization
- **Data Retention**: Delete interaction data after 90 days
- **Export/Delete**: Users can export or delete their data anytime
- **No Cross-site Tracking**: Only track within the app

## 3. Mobile Apps

### Approach: Progressive Web App + Native Wrappers

#### Option A: PWA Only (Recommended for Phase 3)
- Already implemented in Phase 1
- "Add to Home Screen" on iOS/Android
- Offline support via service worker
- Push notifications via Web Push API
- No app store approval needed

#### Option B: Capacitor Wrapper
```bash
# Install Capacitor
npm install @capacitor/core @capacitor/cli
npx cap init

# Add platforms
npx cap add ios
npx cap add android

# Build web app
npm run build

# Copy to native projects
npx cap copy

# Open in Xcode/Android Studio
npx cap open ios
npx cap open android
```

**Pros:**
- Native app store presence
- Better performance
- Native features (camera, biometrics, etc.)
- Push notifications more reliable

**Cons:**
- App store approval process
- Maintenance overhead
- Separate builds for iOS/Android

#### Option C: Tauri (Desktop Apps)
```bash
# Install Tauri
cargo install tauri-cli

# Create Tauri project
cargo tauri init

# Build
cargo tauri build
```

**Pros:**
- Native desktop apps (Windows, macOS, Linux)
- Small bundle size (~3MB)
- Rust backend for security

**Cons:**
- Desktop-only (no mobile)
- Requires Rust toolchain

### Recommended: Start with PWA, Add Capacitor Later

**Phase 3A: Enhance PWA**
1. Add Web Push notifications
2. Improve offline experience
3. Add share target (share to app)
4. Add shortcuts (quick actions)

**Phase 3B: Capacitor Wrapper (Optional)**
1. Wrap PWA with Capacitor
2. Add native splash screen
3. Submit to app stores
4. Add native features as needed

## 4. Additional Features

### A. Email Digests
```go
// internal/email/digest.go
type DigestService struct {
    smtpConfig SMTPConfig
}

func (d *DigestService) SendDailyDigest(userID uuid.UUID) error {
    // Get top items from last 24h
    // Personalize for user
    // Generate HTML email
    // Send via SMTP
}

// Cron job in worker
c.AddFunc("0 8 * * *", func() {
    // Send daily digest at 8 AM
    digestService.SendDailyDigests()
})
```

### B. Slack/Discord Integration
```go
// internal/integrations/slack.go
type SlackIntegration struct {
    webhookURL string
}

func (s *SlackIntegration) PostNewItem(item *models.Item, score float64) error {
    // Format message
    // Post to Slack webhook
}

// internal/integrations/discord.go
type DiscordIntegration struct {
    webhookURL string
}

func (d *DiscordIntegration) PostNewItem(item *models.Item, score float64) error {
    // Format embed
    // Post to Discord webhook
}
```

### C. Browser Extension
```javascript
// chrome-extension/background.js
chrome.runtime.onInstalled.addListener(() => {
    // Set up context menu
    chrome.contextMenus.create({
        id: "save-to-radar",
        title: "Save to Evolipia Radar",
        contexts: ["link", "page"]
    });
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
    // Send URL to API
    fetch('http://localhost:8080/v1/items/submit', {
        method: 'POST',
        body: JSON.stringify({ url: info.linkUrl || info.pageUrl })
    });
});
```

### D. API SDKs
Generate SDKs for popular languages:

```bash
# Python SDK
pip install evolipia-radar

from evolipia_radar import RadarClient

client = RadarClient(api_key="xxx")
feed = client.get_feed(date="today", topic="llm")

# TypeScript SDK
npm install @evolipia/radar-sdk

import { RadarClient } from '@evolipia/radar-sdk';

const client = new RadarClient({ apiKey: 'xxx' });
const feed = await client.getFeed({ date: 'today', topic: 'llm' });
```

## Implementation Timeline

### Phase 3A: Admin Dashboard (4-6 weeks)
- Week 1-2: Admin API endpoints
- Week 3-4: Admin UI (React)
- Week 5: Authentication & authorization
- Week 6: Testing & deployment

### Phase 3B: Personalization (3-4 weeks)
- Week 1: User tracking & preferences
- Week 2: Personalization algorithm
- Week 3: Recommendations engine
- Week 4: Testing & tuning

### Phase 3C: Mobile Apps (2-3 weeks)
- Week 1: Enhance PWA features
- Week 2: Capacitor wrapper (optional)
- Week 3: App store submission (optional)

### Phase 3D: Integrations (2-3 weeks)
- Week 1: Email digests
- Week 2: Slack/Discord webhooks
- Week 3: Browser extension

## Testing Strategy

### Admin Dashboard
- Unit tests for admin API endpoints
- Integration tests for auth flow
- E2E tests with Playwright
- Load testing for analytics queries

### Personalization
- A/B testing framework
- Metrics: CTR, time on site, return rate
- Privacy audit
- Performance testing (recommendation latency)

### Mobile Apps
- Test on real devices (iOS/Android)
- Offline functionality testing
- Push notification testing
- App store review guidelines compliance

## Deployment Considerations

### Scaling
- **Horizontal Scaling**: Multiple API servers behind load balancer
- **Database**: Read replicas for analytics queries
- **Caching**: Redis for hot data (feed, scores)
- **CDN**: CloudFlare for static assets

### Monitoring
- **APM**: New Relic or Datadog
- **Logs**: ELK stack or Loki
- **Metrics**: Prometheus + Grafana
- **Alerts**: PagerDuty for critical issues

### Security
- **Rate Limiting**: Per-user and per-IP
- **DDoS Protection**: CloudFlare
- **SQL Injection**: Parameterized queries (already done)
- **XSS**: Content Security Policy headers
- **HTTPS**: Let's Encrypt certificates

## Cost Estimates (Monthly)

### Infrastructure
- **Compute**: $50-200 (2-4 VMs)
- **Database**: $50-100 (Managed PostgreSQL)
- **CDN**: $10-50 (CloudFlare)
- **Monitoring**: $50-100 (Datadog/New Relic)

### APIs
- **LLM Summarization**: $10-50 (depends on volume)
- **Embeddings**: $5-20 (for vector search)
- **Email**: $10-30 (SendGrid/Mailgun)

### Total: $185-550/month for production deployment

## Success Metrics

### Phase 3 Goals
- **Admin Dashboard**: 100% source health visibility
- **Personalization**: 20% increase in user engagement
- **Mobile Apps**: 1000+ installs in first month
- **Integrations**: 50+ Slack/Discord installations

### KPIs to Track
- Daily Active Users (DAU)
- Average session duration
- Items clicked per session
- Return rate (7-day, 30-day)
- API response time (p95 < 200ms)
- Error rate (< 0.1%)

## Next Steps After Phase 3

### Future Enhancements
- **Multi-language Support**: Translate summaries
- **Audio Summaries**: Text-to-speech for commuting
- **Video Content**: YouTube/podcast integration
- **Social Features**: Share, comment, discuss
- **Premium Tier**: Advanced features for power users
- **API Marketplace**: Let others build on top

---

## Summary

Phase 3 transforms Evolipia Radar from a functional aggregator into a polished, production-ready platform with:
- Professional admin tools for operators
- Personalized experiences for users
- Mobile apps for on-the-go access
- Integrations for workflow automation

The modular approach allows implementing features incrementally based on user feedback and priorities.
