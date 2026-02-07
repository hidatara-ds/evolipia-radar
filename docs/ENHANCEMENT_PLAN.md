# Enhancement Plan: evolipia-radar → User-Facing App Backend

## A) Gap Analysis

### Current State
- ✅ News ingestion from multiple sources (HN, RSS, arXiv, JSON API)
- ✅ Ranking/scoring system (popularity, relevance, credibility, novelty)
- ✅ Deduplication and normalization
- ✅ REST API for feeds and search
- ✅ SSRF protection
- ✅ Worker-based scheduled ingestion
- ✅ PostgreSQL with migrations

### Missing Components (Mapped to Goals)

#### 1. Authentication + User Profiles
**Missing:**
- User registration/login endpoints
- JWT token management (access + refresh)
- Password hashing (bcrypt/argon2)
- User profile storage (email, timezone, preferences)
- Session management
- Email verification (optional)

**Impact:** No user identity, no personalization possible

#### 2. Personalization (Preferences + Feedback Loop)
**Missing:**
- User preferences table (topics, blocked sources, language)
- Feedback mechanism (upvote/downvote/hide/duplicate flags)
- Personalization scoring adjustment
- User-item interaction tracking

**Impact:** All users see same ranking, no learning from behavior

#### 3. Bookmark & Read History
**Missing:**
- Bookmarks table (user_id, item_id, created_at)
- Read history table (user_id, item_id, read_at)
- API endpoints for CRUD operations
- Privacy controls (delete history)

**Impact:** No way to save items or track reading

#### 4. Daily/Weekly Digest + Keyword Alerts
**Missing:**
- Digest generation logic
- User digest preferences (frequency, time, format)
- Keyword alert subscriptions
- Email/notification delivery (email service integration)
- Timezone-aware scheduling
- Unsubscribe mechanism

**Impact:** No proactive engagement, users must check manually

#### 5. Score Explainability
**Missing:**
- Score breakdown storage (already computed but not exposed)
- Explanation text generation
- API response enhancement to include breakdown
- Visual explanation format

**Impact:** Users can't understand why items rank high/low

#### 6. Observability + Production CI/CD
**Missing:**
- Structured logging (request_id, user_id, source_id)
- Metrics collection (Prometheus/StatsD)
- Distributed tracing (OpenTelemetry)
- Health check endpoints with dependencies
- CI/CD pipeline (GitHub Actions)
- Migration safety checks
- Docker build automation
- Error tracking (Sentry/rollbar)

**Impact:** Hard to debug, monitor, or deploy safely

---

## B) Target Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                             │
│  (Web App / Mobile App / API Consumers)                         │
└────────────────────┬────────────────────────────────────────────┘
                     │ HTTPS
                     │
┌────────────────────▼────────────────────────────────────────────┐
│                      API SERVER (cmd/api)                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  HTTP Handlers (internal/http/handlers)                   │  │
│  │  - Auth handlers (login, register, refresh)              │  │
│  │  - Feed handlers (personalized)                           │  │
│  │  - Bookmark handlers                                      │  │
│  │  - Digest handlers                                        │  │
│  │  - Middleware: Auth, Logging, Rate Limiting               │  │
│  └──────────────┬───────────────────────────────────────────┘  │
│                 │                                                │
│  ┌──────────────▼───────────────────────────────────────────┐  │
│  │  Services Layer (internal/services)                       │  │
│  │  - AuthService (JWT, password hashing)                    │  │
│  │  - UserService (profiles, preferences)                   │  │
│  │  - FeedService (personalized ranking)                    │  │
│  │  - BookmarkService                                       │  │
│  │  - DigestService                                         │  │
│  │  - AlertService                                          │  │
│  └──────────────┬───────────────────────────────────────────┘  │
└─────────────────┼───────────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────────┐
│                    DATABASE (PostgreSQL)                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Existing Tables:                                        │  │
│  │  - sources, items, signals, scores, summaries, fetch_runs │  │
│  │                                                            │  │
│  │  New Tables:                                              │  │
│  │  - users (id, email, password_hash, timezone, ...)       │  │
│  │  - user_preferences (user_id, topics[], blocked_sources[])│  │
│  │  - bookmarks (user_id, item_id, created_at)              │  │
│  │  - read_history (user_id, item_id, read_at)              │  │
│  │  - feedback (user_id, item_id, type, created_at)        │  │
│  │  - digest_subscriptions (user_id, frequency, time, ...) │  │
│  │  - keyword_alerts (user_id, keyword, enabled)           │  │
│  │  - digest_jobs (id, user_id, scheduled_at, status)      │  │
│  │  - refresh_tokens (token_hash, user_id, expires_at)      │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    WORKER (cmd/worker)                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Scheduled Jobs:                                         │  │
│  │  1. IngestionWorker (existing)                           │  │
│  │     - Fetches from sources                               │  │
│  │     - Normalizes & deduplicates                         │  │
│  │     - Computes scores                                    │  │
│  │                                                            │  │
│  │  2. DigestWorker (NEW)                                    │  │
│  │     - Timezone-aware scheduling                          │  │
│  │     - Generates daily/weekly digests                     │  │
│  │     - Sends via email (SMTP/SendGrid)                    │  │
│  │                                                            │  │
│  │  3. AlertWorker (NEW)                                     │  │
│  │     - Checks keyword matches                             │  │
│  │     - Sends real-time alerts                             │  │
│  │                                                            │  │
│  │  4. ScoreRecalculationWorker (NEW)                        │  │
│  │     - Recomputes scores with feedback                    │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    EXTERNAL SERVICES                            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Email Provider (SMTP/SendGrid/Mailgun)                  │  │
│  │  - Digest delivery                                       │  │
│  │  - Alert notifications                                   │  │
│  │  - Verification emails                                   │  │
│  │                                                            │  │
│  │  Optional: Telegram Bot (for alerts)                      │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    OBSERVABILITY STACK                          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Logging: Structured logs → stdout (JSON)                │  │
│  │  Metrics: Prometheus → Grafana                           │  │
│  │  Tracing: OpenTelemetry → Jaeger/Tempo                   │  │
│  │  Errors: Sentry                                          │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

**Data Flow Examples:**

1. **User Login Flow:**
   ```
   Client → API Handler → AuthService → DB (users)
   AuthService → Generate JWT → Return access + refresh tokens
   ```

2. **Personalized Feed Flow:**
   ```
   Client (with JWT) → API Handler → Auth Middleware → FeedService
   FeedService → UserService (get preferences) → FeedService (apply personalization)
   FeedService → ItemRepository (query with filters) → Return personalized feed
   ```

3. **Digest Generation Flow:**
   ```
   Cron Trigger → DigestWorker → Get users with digest subscriptions
   For each user: Generate digest → EmailService → Send email
   ```

---

## C) Database Design

### New Tables + Indexes

**File: `migrations/000002_add_user_features.up.sql`**

```sql
-- ============================================
-- USER AUTHENTICATION & PROFILES
-- ============================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL, -- bcrypt/argon2 hash
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    timezone TEXT NOT NULL DEFAULT 'UTC', -- IANA timezone (e.g., 'America/New_York')
    language TEXT NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL -- Soft delete
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- ============================================
-- REFRESH TOKENS (for JWT rotation)
-- ============================================

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE, -- SHA256 hash of refresh token
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id, expires_at DESC);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash) WHERE revoked_at IS NULL;

-- ============================================
-- USER PREFERENCES
-- ============================================

CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    preferred_topics JSONB NOT NULL DEFAULT '[]'::jsonb, -- ["llm", "mlops", "computer_vision"]
    blocked_sources JSONB NOT NULL DEFAULT '[]'::jsonb, -- [source_id UUIDs]
    blocked_domains JSONB NOT NULL DEFAULT '[]'::jsonb, -- ["medium.com", "example.com"]
    personalization_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    personalization_weight DOUBLE PRECISION NOT NULL DEFAULT 0.1, -- 0.0-0.3 range
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_preferences_topics_gin ON user_preferences USING GIN (preferred_topics);
CREATE INDEX idx_user_preferences_blocked_sources_gin ON user_preferences USING GIN (blocked_sources);

-- ============================================
-- BOOKMARKS
-- ============================================

CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    notes TEXT NULL, -- User's personal notes
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, item_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id, created_at DESC);
CREATE INDEX idx_bookmarks_item_id ON bookmarks(item_id);

-- ============================================
-- READ HISTORY
-- ============================================

CREATE TABLE read_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, item_id)
);

CREATE INDEX idx_read_history_user_id ON read_history(user_id, read_at DESC);
CREATE INDEX idx_read_history_item_id ON read_history(item_id);
CREATE INDEX idx_read_history_user_read_at ON read_history(user_id, read_at DESC) WHERE read_at > now() - interval '90 days'; -- Partial index for recent reads

-- ============================================
-- USER FEEDBACK
-- ============================================

CREATE TABLE feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- 'upvote', 'downvote', 'hide', 'duplicate', 'not_relevant'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, item_id, type) -- One feedback type per user-item pair
);

CREATE INDEX idx_feedback_user_id ON feedback(user_id, created_at DESC);
CREATE INDEX idx_feedback_item_id ON feedback(item_id, type);
CREATE INDEX idx_feedback_type ON feedback(type, created_at DESC);

-- ============================================
-- DIGEST SUBSCRIPTIONS
-- ============================================

CREATE TABLE digest_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    frequency TEXT NOT NULL, -- 'daily', 'weekly'
    send_at TIME NOT NULL, -- Local time (e.g., '09:00:00')
    timezone TEXT NOT NULL, -- User's timezone for send_at
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_sent_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, frequency)
);

CREATE INDEX idx_digest_subscriptions_enabled ON digest_subscriptions(enabled, frequency) WHERE enabled = TRUE;
CREATE INDEX idx_digest_subscriptions_user_id ON digest_subscriptions(user_id);

-- ============================================
-- KEYWORD ALERTS
-- ============================================

CREATE TABLE keyword_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    keyword TEXT NOT NULL, -- Simple substring match (sanitized, no regex)
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, keyword)
);

CREATE INDEX idx_keyword_alerts_user_id ON keyword_alerts(user_id, enabled) WHERE enabled = TRUE;
CREATE INDEX idx_keyword_alerts_keyword ON keyword_alerts(keyword) WHERE enabled = TRUE;

-- ============================================
-- DIGEST JOBS (for tracking digest generation)
-- ============================================

CREATE TABLE digest_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES digest_subscriptions(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    error_message TEXT NULL,
    items_count INT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_digest_jobs_status ON digest_jobs(status, scheduled_at) WHERE status IN ('pending', 'processing');
CREATE INDEX idx_digest_jobs_user_id ON digest_jobs(user_id, scheduled_at DESC);

-- ============================================
-- ENHANCE EXISTING SCORES TABLE
-- ============================================

-- Add explanation fields to scores table (for explainability)
ALTER TABLE scores ADD COLUMN IF NOT EXISTS explanation JSONB NULL;
-- Format: {"popularity": {"value": 0.8, "reason": "High HN points (500) with recent decay"}, ...}

-- Add index for explanation queries (if needed)
CREATE INDEX IF NOT EXISTS idx_scores_explanation_gin ON scores USING GIN (explanation);

-- ============================================
-- RATE LIMITING (per-domain tracking)
-- ============================================

CREATE TABLE rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain TEXT NOT NULL,
    source_id UUID NULL REFERENCES sources(id) ON DELETE SET NULL,
    request_count INT NOT NULL DEFAULT 0,
    window_start TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(domain, source_id, window_start)
);

CREATE INDEX idx_rate_limits_domain ON rate_limits(domain, window_start DESC);
CREATE INDEX idx_rate_limits_source_id ON rate_limits(source_id, window_start DESC);
```

**File: `migrations/000002_add_user_features.down.sql`**

```sql
DROP TABLE IF EXISTS digest_jobs CASCADE;
DROP TABLE IF EXISTS keyword_alerts CASCADE;
DROP TABLE IF EXISTS digest_subscriptions CASCADE;
DROP TABLE IF EXISTS feedback CASCADE;
DROP TABLE IF EXISTS read_history CASCADE;
DROP TABLE IF EXISTS bookmarks CASCADE;
DROP TABLE IF EXISTS user_preferences CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS rate_limits CASCADE;

ALTER TABLE scores DROP COLUMN IF EXISTS explanation;
DROP INDEX IF EXISTS idx_scores_explanation_gin;
```

---

## D) API Design

### Authentication Endpoints

#### POST /v1/auth/register
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "timezone": "America/New_York",
  "language": "en"
}
```

**Response (201):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "timezone": "America/New_York",
    "email_verified": false
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
  "expires_in": 3600
}
```

#### POST /v1/auth/login
**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200):** Same as register

#### POST /v1/auth/refresh
**Request:**
```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

#### POST /v1/auth/logout
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

### User Profile Endpoints

#### GET /v1/users/me
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "timezone": "America/New_York",
  "language": "en",
  "email_verified": true,
  "created_at": "2024-01-15T10:00:00Z",
  "last_login_at": "2024-01-20T14:30:00Z"
}
```

#### PATCH /v1/users/me
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "timezone": "Europe/London",
  "language": "en"
}
```

**Response (200):** Updated user object

### Preferences Endpoints

#### GET /v1/users/me/preferences
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "preferred_topics": ["llm", "mlops", "computer_vision"],
  "blocked_sources": ["source-uuid-1", "source-uuid-2"],
  "blocked_domains": ["medium.com"],
  "personalization_enabled": true,
  "personalization_weight": 0.15
}
```

#### PATCH /v1/users/me/preferences
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "preferred_topics": ["llm", "nlp"],
  "blocked_domains": ["medium.com", "example.com"],
  "personalization_weight": 0.2
}
```

**Response (200):** Updated preferences

### Personalized Feed Endpoints

#### GET /v1/feed?date=today&topic=llm&limit=20&offset=0
**Headers:** `Authorization: Bearer <access_token>` (optional, for personalization)

**Response (200):**
```json
{
  "date": "2024-01-20",
  "topic": "llm",
  "items": [
    {
      "id": "item-uuid",
      "rank": 1,
      "title": "New LLM Architecture",
      "url": "https://example.com/article",
      "domain": "example.com",
      "published_at": "2024-01-20T10:00:00Z",
      "scores": {
        "final": 0.85,
        "hot": 0.75,
        "relevance": 0.90,
        "credibility": 0.80,
        "novelty": 0.70
      },
      "score_explanation": {
        "popularity": {
          "value": 0.75,
          "reason": "High HN engagement (450 points, 120 comments) with moderate recency decay"
        },
        "relevance": {
          "value": 0.90,
          "reason": "Strong match: 'llm', 'transformer' keywords found; tagged as 'llm'"
        },
        "credibility": {
          "value": 0.80,
          "reason": "Domain 'example.com' is in whitelist"
        },
        "novelty": {
          "value": 0.70,
          "reason": "Published 2 hours ago"
        },
        "personalization_boost": {
          "value": 0.05,
          "reason": "User preference match: topic 'llm'"
        }
      },
      "summary": {
        "tldr": "Researchers introduce new transformer architecture...",
        "why_it_matters": "This could impact how AI engineers build language models...",
        "tags": ["llm", "nlp"],
        "method": "extractive"
      },
      "user_context": {
        "bookmarked": false,
        "read": false,
        "feedback": null
      }
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

### Bookmark Endpoints

#### GET /v1/bookmarks?limit=20&offset=0
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "items": [
    {
      "id": "bookmark-uuid",
      "item": {
        "id": "item-uuid",
        "title": "Article Title",
        "url": "https://example.com/article",
        "domain": "example.com",
        "published_at": "2024-01-15T10:00:00Z"
      },
      "notes": "Remember to read this",
      "created_at": "2024-01-16T08:00:00Z"
    }
  ],
  "pagination": {
    "total": 45,
    "limit": 20,
    "offset": 0
  }
}
```

#### POST /v1/bookmarks
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "item_id": "item-uuid",
  "notes": "Optional notes"
}
```

**Response (201):**
```json
{
  "id": "bookmark-uuid",
  "item_id": "item-uuid",
  "notes": "Optional notes",
  "created_at": "2024-01-20T15:00:00Z"
}
```

#### DELETE /v1/bookmarks/{bookmark_id}
**Headers:** `Authorization: Bearer <access_token>`

**Response (204):** No content

### Read History Endpoints

#### GET /v1/history?limit=20&offset=0
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "items": [
    {
      "item_id": "item-uuid",
      "item": {
        "id": "item-uuid",
        "title": "Article Title",
        "url": "https://example.com/article",
        "domain": "example.com"
      },
      "read_at": "2024-01-19T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 120,
    "limit": 20,
    "offset": 0
  }
}
```

#### POST /v1/history
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "item_id": "item-uuid"
}
```

**Response (201):**
```json
{
  "item_id": "item-uuid",
  "read_at": "2024-01-20T15:00:00Z"
}
```

#### DELETE /v1/history
**Headers:** `Authorization: Bearer <access_token>`

**Query params:** `?item_id=<uuid>` or `?all=true`

**Response (204):** No content

### Feedback Endpoints

#### POST /v1/feedback
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "item_id": "item-uuid",
  "type": "upvote" // or "downvote", "hide", "duplicate", "not_relevant"
}
```

**Response (201):**
```json
{
  "id": "feedback-uuid",
  "item_id": "item-uuid",
  "type": "upvote",
  "created_at": "2024-01-20T15:00:00Z"
}
```

#### DELETE /v1/feedback/{feedback_id}
**Headers:** `Authorization: Bearer <access_token>`

**Response (204):** No content

### Digest Endpoints

#### GET /v1/digests/subscriptions
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "subscriptions": [
    {
      "id": "subscription-uuid",
      "frequency": "daily",
      "send_at": "09:00:00",
      "timezone": "America/New_York",
      "enabled": true,
      "last_sent_at": "2024-01-20T09:00:00Z"
    }
  ]
}
```

#### POST /v1/digests/subscriptions
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "frequency": "daily", // or "weekly"
  "send_at": "09:00:00", // HH:MM:SS in user's timezone
  "timezone": "America/New_York"
}
```

**Response (201):** Subscription object

#### PATCH /v1/digests/subscriptions/{subscription_id}
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "enabled": false
}
```

**Response (200):** Updated subscription

#### DELETE /v1/digests/subscriptions/{subscription_id}
**Headers:** `Authorization: Bearer <access_token>`

**Response (204):** No content

### Keyword Alerts Endpoints

#### GET /v1/alerts/keywords
**Headers:** `Authorization: Bearer <access_token>`

**Response (200):**
```json
{
  "alerts": [
    {
      "id": "alert-uuid",
      "keyword": "GPT-5",
      "enabled": true,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### POST /v1/alerts/keywords
**Headers:** `Authorization: Bearer <access_token>`

**Request:**
```json
{
  "keyword": "GPT-5" // Simple substring, sanitized
}
```

**Response (201):** Alert object

#### DELETE /v1/alerts/keywords/{alert_id}
**Headers:** `Authorization: Bearer <access_token>`

**Response (204):** No content

### Score Explainability (Enhanced Existing Endpoint)

#### GET /v1/items/{id}/explanation
**Headers:** `Authorization: Bearer <access_token>` (optional)

**Response (200):**
```json
{
  "item_id": "item-uuid",
  "scores": {
    "final": 0.85,
    "hot": 0.75,
    "relevance": 0.90,
    "credibility": 0.80,
    "novelty": 0.70
  },
  "explanation": {
    "popularity": {
      "value": 0.75,
      "weight": 0.55,
      "contribution": 0.4125,
      "reason": "High HN engagement (450 points, 120 comments) with moderate recency decay (published 3 hours ago)"
    },
    "relevance": {
      "value": 0.90,
      "weight": 0.25,
      "contribution": 0.225,
      "reason": "Strong match: keywords 'llm', 'transformer' found in title; item tagged as 'llm'"
    },
    "credibility": {
      "value": 0.80,
      "weight": 0.15,
      "contribution": 0.12,
      "reason": "Domain 'example.com' is in credibility whitelist"
    },
    "novelty": {
      "value": 0.70,
      "weight": 0.05,
      "contribution": 0.035,
      "reason": "Published 3 hours ago (within 7-day novelty window)"
    },
    "personalization": {
      "value": 0.05,
      "weight": 0.15,
      "contribution": 0.0075,
      "reason": "User preference match: topic 'llm' in preferred topics",
      "enabled": true
    }
  },
  "breakdown": {
    "base_score": 0.7925,
    "personalization_boost": 0.0075,
    "final_score": 0.85
  }
}
```

---

## E) Implementation Plan (10-20 PR-sized Steps)

### PR 1: Database Schema - User Authentication
**Files:**
- `migrations/000002_add_user_features.up.sql`
- `migrations/000002_add_user_features.down.sql`
- `internal/models/users.go` (new)
- `internal/db/repositories.go` (add UserRepository)

**Acceptance Criteria:**
- Users table created with email, password_hash, timezone
- Refresh tokens table created
- UserRepository with Create, GetByEmail, GetByID methods
- Migration runs successfully up and down

**Tests:**
- `internal/db/repositories_test.go` - UserRepository tests
- Migration test (up/down)

---

### PR 2: Authentication Service & JWT
**Files:**
- `internal/auth/jwt.go` (new)
- `internal/auth/password.go` (new)
- `internal/services/auth_service.go` (new)
- `internal/config/config.go` (add JWT_SECRET, JWT_EXPIRY)

**Acceptance Criteria:**
- JWT generation/validation with access + refresh tokens
- Password hashing with bcrypt (cost 10)
- AuthService: Register, Login, Refresh, ValidateToken
- Config via env vars (JWT_SECRET, JWT_ACCESS_EXPIRY, JWT_REFRESH_EXPIRY)

**Tests:**
- `internal/auth/jwt_test.go`
- `internal/auth/password_test.go`
- `internal/services/auth_service_test.go`

---

### PR 3: Auth HTTP Handlers & Middleware
**Files:**
- `internal/http/handlers/auth_handlers.go` (new)
- `internal/http/middleware/auth.go` (new)
- `cmd/api/main.go` (register auth routes)

**Acceptance Criteria:**
- POST /v1/auth/register, /login, /refresh, /logout endpoints
- Auth middleware extracts user from JWT
- Error handling for invalid credentials, expired tokens
- Rate limiting on auth endpoints (5 req/min)

**Tests:**
- `internal/http/handlers/auth_handlers_test.go`
- Integration test: register → login → access protected endpoint

---

### PR 4: User Preferences & Profile
**Files:**
- `internal/models/user_preferences.go` (new)
- `internal/db/repositories.go` (add UserPreferencesRepository)
- `internal/services/user_service.go` (new)
- `internal/http/handlers/user_handlers.go` (new)

**Acceptance Criteria:**
- UserPreferencesRepository CRUD
- UserService: GetProfile, UpdateProfile, GetPreferences, UpdatePreferences
- GET/PATCH /v1/users/me, /v1/users/me/preferences endpoints
- Validation: timezone (IANA), topics (allowed list)

**Tests:**
- `internal/services/user_service_test.go`
- `internal/http/handlers/user_handlers_test.go`

---

### PR 5: Bookmarks Feature
**Files:**
- `internal/models/bookmarks.go` (new)
- `internal/db/repositories.go` (add BookmarkRepository)
- `internal/services/bookmark_service.go` (new)
- `internal/http/handlers/bookmark_handlers.go` (new)

**Acceptance Criteria:**
- BookmarkRepository: Create, List, Delete, GetByUserAndItem
- BookmarkService: AddBookmark, ListBookmarks, RemoveBookmark
- GET/POST/DELETE /v1/bookmarks endpoints
- Pagination support (limit/offset)

**Tests:**
- `internal/services/bookmark_service_test.go`
- `internal/http/handlers/bookmark_handlers_test.go`

---

### PR 6: Read History Feature
**Files:**
- `internal/models/read_history.go` (new)
- `internal/db/repositories.go` (add ReadHistoryRepository)
- `internal/services/history_service.go` (new)
- `internal/http/handlers/history_handlers.go` (new)

**Acceptance Criteria:**
- ReadHistoryRepository: Create, List, Delete (by item or all)
- HistoryService: MarkAsRead, GetHistory, ClearHistory
- GET/POST/DELETE /v1/history endpoints
- Privacy: users can delete their own history

**Tests:**
- `internal/services/history_service_test.go`
- `internal/http/handlers/history_handlers_test.go`

---

### PR 7: Feedback System
**Files:**
- `internal/models/feedback.go` (new)
- `internal/db/repositories.go` (add FeedbackRepository)
- `internal/services/feedback_service.go` (new)
- `internal/http/handlers/feedback_handlers.go` (new)

**Acceptance Criteria:**
- FeedbackRepository: Create, GetByUserAndItem, Delete
- FeedbackService: SubmitFeedback, RemoveFeedback
- POST/DELETE /v1/feedback endpoints
- Feedback types: upvote, downvote, hide, duplicate, not_relevant

**Tests:**
- `internal/services/feedback_service_test.go`
- `internal/http/handlers/feedback_handlers_test.go`

---

### PR 8: Personalized Feed Ranking
**Files:**
- `internal/services/feed_service.go` (modify)
- `internal/scoring/personalization.go` (new)

**Acceptance Criteria:**
- FeedService applies personalization when user authenticated
- Personalization boost based on:
  - Preferred topics match (+0.05-0.15)
  - Blocked sources/domains filter
  - Feedback signals (hide items user downvoted)
- Keep global ranking as baseline, add small boost (0.1-0.3 weight)

**Tests:**
- `internal/scoring/personalization_test.go`
- `internal/services/feed_service_test.go` (personalized vs non-personalized)

---

### PR 9: Score Explainability
**Files:**
- `internal/scoring/explanation.go` (new)
- `internal/db/repositories.go` (modify ScoreRepository - add explanation)
- `internal/services/feed_service.go` (add explanation to responses)

**Acceptance Criteria:**
- Explanation generation for each score component
- Store explanation JSONB in scores table
- Return explanation in GET /v1/items/{id} and feed endpoints
- Format: value, weight, contribution, reason text

**Tests:**
- `internal/scoring/explanation_test.go`
- Verify explanation in API responses

---

### PR 10: Digest Subscriptions (Database + API)
**Files:**
- `internal/models/digest_subscriptions.go` (new)
- `internal/db/repositories.go` (add DigestSubscriptionRepository)
- `internal/services/digest_service.go` (new, partial)
- `internal/http/handlers/digest_handlers.go` (new)

**Acceptance Criteria:**
- DigestSubscriptionRepository CRUD
- GET/POST/PATCH/DELETE /v1/digests/subscriptions endpoints
- Validation: frequency (daily/weekly), send_at (HH:MM:SS), timezone

**Tests:**
- `internal/services/digest_service_test.go`
- `internal/http/handlers/digest_handlers_test.go`

---

### PR 11: Keyword Alerts (Database + API)
**Files:**
- `internal/models/keyword_alerts.go` (new)
- `internal/db/repositories.go` (add KeywordAlertRepository)
- `internal/services/alert_service.go` (new, partial)
- `internal/http/handlers/alert_handlers.go` (new)

**Acceptance Criteria:**
- KeywordAlertRepository CRUD
- GET/POST/DELETE /v1/alerts/keywords endpoints
- Keyword sanitization (no regex, max length 100)

**Tests:**
- `internal/services/alert_service_test.go`
- `internal/http/handlers/alert_handlers_test.go`

---

### PR 12: Digest Generation Worker
**Files:**
- `internal/services/digest_service.go` (complete)
- `internal/services/worker.go` (add DigestWorker)
- `cmd/worker/main.go` (register digest cron job)

**Acceptance Criteria:**
- DigestWorker runs on schedule (every hour)
- Timezone-aware scheduling (convert user's send_at to UTC)
- Generate digest content (top N items with explanations)
- Store digest_jobs for tracking
- Support daily and weekly frequencies

**Tests:**
- `internal/services/digest_service_test.go`
- `internal/services/worker_test.go` (digest generation)

---

### PR 13: Email Service Integration
**Files:**
- `internal/email/service.go` (new)
- `internal/email/templates.go` (new)
- `internal/config/config.go` (add SMTP config)
- `internal/services/digest_service.go` (integrate email)

**Acceptance Criteria:**
- EmailService with SMTP support
- HTML email templates for digests
- Config via env vars (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS)
- Send digest emails
- Handle email errors gracefully

**Tests:**
- `internal/email/service_test.go` (mock SMTP)
- Integration test with test SMTP server

---

### PR 14: Keyword Alert Worker
**Files:**
- `internal/services/alert_service.go` (complete)
- `internal/services/worker.go` (add AlertWorker)
- `cmd/worker/main.go` (register alert cron job)

**Acceptance Criteria:**
- AlertWorker runs every 15 minutes
- Check new items for keyword matches
- Send email alerts for matches
- Rate limit: max 5 alerts per user per hour

**Tests:**
- `internal/services/alert_service_test.go`
- Keyword matching logic tests

---

### PR 15: Structured Logging
**Files:**
- `internal/logging/logger.go` (new)
- `internal/http/middleware/logging.go` (new)
- `cmd/api/main.go` (setup logger)
- `cmd/worker/main.go` (setup logger)

**Acceptance Criteria:**
- Structured JSON logging (logrus/zap)
- Fields: request_id, user_id, source_id, method, path, status, latency
- Request ID generation and propagation
- Log levels: DEBUG, INFO, WARN, ERROR

**Tests:**
- Verify log output format
- Request ID propagation test

---

### PR 16: Metrics Collection
**Files:**
- `internal/metrics/metrics.go` (new)
- `internal/http/middleware/metrics.go` (new)
- `internal/services/worker.go` (add metrics)

**Acceptance Criteria:**
- Prometheus metrics:
  - `http_requests_total` (counter, labels: method, path, status)
  - `http_request_duration_seconds` (histogram)
  - `ingestion_items_total` (counter, labels: source_id, status)
  - `ingestion_duration_seconds` (histogram)
  - `dedup_rate` (gauge)
  - `db_query_duration_seconds` (histogram)
- `/metrics` endpoint for Prometheus

**Tests:**
- `internal/metrics/metrics_test.go`
- Verify metrics exported correctly

---

### PR 17: Health Check Enhancement
**Files:**
- `internal/http/handlers/health_handlers.go` (new)
- `cmd/api/main.go` (register /healthz)

**Acceptance Criteria:**
- `/healthz` returns 200 if healthy
- `/healthz/ready` checks database connectivity
- `/healthz/live` simple liveness check
- Response includes version, timestamp

**Tests:**
- Health check endpoint tests
- Database connectivity test

---

### PR 18: Rate Limiting (Per-Domain)
**Files:**
- `internal/db/repositories.go` (add RateLimitRepository)
- `internal/services/rate_limiter.go` (new)
- `internal/http/middleware/rate_limit.go` (new)

**Acceptance Criteria:**
- Per-domain rate limiting (100 req/hour per domain)
- Per-source rate limiting (configurable)
- Rate limit headers in responses
- Store rate limits in database

**Tests:**
- `internal/services/rate_limiter_test.go`
- Rate limit middleware tests

---

### PR 19: CI/CD Pipeline
**Files:**
- `.github/workflows/ci.yml` (new)
- `.github/workflows/release.yml` (new)
- `Dockerfile.api` (update)
- `Dockerfile.worker` (update)

**Acceptance Criteria:**
- GitHub Actions: lint, test, build on PR
- Migration safety check (no destructive migrations in PR)
- Docker build and push on tag
- Release workflow with versioning

**Tests:**
- CI pipeline runs successfully
- Docker images build correctly

---

### PR 20: Documentation & OpenAPI
**Files:**
- `docs/openapi.yaml` (new)
- `README.md` (update with new endpoints)
- `DEPLOYMENT.md` (new)

**Acceptance Criteria:**
- OpenAPI 3.0 spec for all endpoints
- Updated README with auth, personalization, bookmarks sections
- Deployment guide with env vars, migration steps
- Architecture diagram in docs

**Tests:**
- OpenAPI spec validates
- All endpoints documented

---

## F) Explainability Spec

### Score Breakdown Storage

**Database Schema:**
```sql
ALTER TABLE scores ADD COLUMN explanation JSONB NULL;
```

**Explanation JSON Structure:**
```json
{
  "popularity": {
    "value": 0.75,
    "weight": 0.55,
    "contribution": 0.4125,
    "reason": "High HN engagement (450 points, 120 comments) with moderate recency decay (published 3 hours ago)"
  },
  "relevance": {
    "value": 0.90,
    "weight": 0.25,
    "contribution": 0.225,
    "reason": "Strong match: keywords 'llm', 'transformer' found in title; item tagged as 'llm'"
  },
  "credibility": {
    "value": 0.80,
    "weight": 0.15,
    "contribution": 0.12,
    "reason": "Domain 'example.com' is in credibility whitelist"
  },
  "novelty": {
    "value": 0.70,
    "weight": 0.05,
    "contribution": 0.035,
    "reason": "Published 3 hours ago (within 7-day novelty window)"
  },
  "personalization": {
    "value": 0.05,
    "weight": 0.15,
    "contribution": 0.0075,
    "reason": "User preference match: topic 'llm' in preferred topics",
    "enabled": true
  }
}
```

### Computation Logic

**File: `internal/scoring/explanation.go`**

```go
type ScoreExplanation struct {
    Popularity      ComponentExplanation `json:"popularity"`
    Relevance       ComponentExplanation `json:"relevance"`
    Credibility     ComponentExplanation `json:"credibility"`
    Novelty         ComponentExplanation `json:"novelty"`
    Personalization *ComponentExplanation `json:"personalization,omitempty"`
}

type ComponentExplanation struct {
    Value        float64 `json:"value"`
    Weight       float64 `json:"weight"`
    Contribution float64 `json:"contribution"`
    Reason       string  `json:"reason"`
    Enabled      bool    `json:"enabled,omitempty"`
}

func GenerateExplanation(
    item *models.Item,
    signal *models.Signal,
    summary *models.Summary,
    weights scoring.Weights,
    personalizationBoost float64,
    personalizationReason string,
) *ScoreExplanation {
    hot := computeHotScore(signal, item.PublishedAt)
    relevance := computeRelevanceScore(item, summary)
    credibility := computeCredibilityScore(item.Domain)
    novelty := computeNoveltyScore(item.PublishedAt)
    
    return &ScoreExplanation{
        Popularity: ComponentExplanation{
            Value:        hot,
            Weight:       weights.W1,
            Contribution: hot * weights.W1,
            Reason:       generatePopularityReason(signal, item.PublishedAt),
        },
        Relevance: ComponentExplanation{
            Value:        relevance,
            Weight:       weights.W2,
            Contribution: relevance * weights.W2,
            Reason:       generateRelevanceReason(item, summary),
        },
        Credibility: ComponentExplanation{
            Value:        credibility,
            Weight:       weights.W3,
            Contribution: credibility * weights.W3,
            Reason:       generateCredibilityReason(item.Domain),
        },
        Novelty: ComponentExplanation{
            Value:        novelty,
            Weight:       weights.W4,
            Contribution: novelty * weights.W4,
            Reason:       generateNoveltyReason(item.PublishedAt),
        },
        Personalization: func() *ComponentExplanation {
            if personalizationBoost > 0 {
                return &ComponentExplanation{
                    Value:        personalizationBoost,
                    Weight:       0.15, // User's personalization_weight
                    Contribution: personalizationBoost * 0.15,
                    Reason:       personalizationReason,
                    Enabled:      true,
                }
            }
            return nil
        }(),
    }
}

func generatePopularityReason(signal *models.Signal, publishedAt time.Time) string {
    if signal == nil {
        return "No engagement signals available"
    }
    points := 0
    comments := 0
    if signal.Points != nil {
        points = *signal.Points
    }
    if signal.Comments != nil {
        comments = *signal.Comments
    }
    ageHours := time.Since(publishedAt).Hours()
    return fmt.Sprintf("HN engagement: %d points, %d comments. Published %.1f hours ago (recency decay applied)", 
        points, comments, ageHours)
}

func generateRelevanceReason(item *models.Item, summary *models.Summary) string {
    reasons := []string{}
    text := strings.ToLower(item.Title)
    if item.RawExcerpt != nil {
        text += " " + strings.ToLower(*item.RawExcerpt)
    }
    
    if strings.Contains(text, "llm") || strings.Contains(text, "transformer") {
        reasons = append(reasons, "keywords 'llm', 'transformer' found")
    }
    if summary != nil && len(summary.Tags) > 0 {
        reasons = append(reasons, fmt.Sprintf("tagged as: %s", strings.Join(summary.Tags, ", ")))
    }
    
    if len(reasons) == 0 {
        return "Baseline relevance score (no strong keyword matches)"
    }
    return strings.Join(reasons, "; ")
}

func generateCredibilityReason(domain string) string {
    config := scoring.DefaultCredibilityConfig()
    if config.Whitelist[domain] {
        return fmt.Sprintf("Domain '%s' is in credibility whitelist", domain)
    }
    if config.Blacklist[domain] {
        return fmt.Sprintf("Domain '%s' is in credibility blacklist", domain)
    }
    return fmt.Sprintf("Domain '%s' has baseline credibility score", domain)
}

func generateNoveltyReason(publishedAt time.Time) string {
    ageHours := time.Since(publishedAt).Hours()
    if ageHours > 168 {
        return fmt.Sprintf("Published %.1f days ago (beyond 7-day novelty window)", ageHours/24)
    }
    return fmt.Sprintf("Published %.1f hours ago (within 7-day novelty window)", ageHours)
}
```

### API Response Integration

**Modify `internal/services/feed_service.go`:**
```go
func (s *FeedService) BuildFeedResponse(ctx context.Context, items []models.Item, date time.Time, topic *string, userID *uuid.UUID) map[string]interface{} {
    // ... existing code ...
    
    for rank, item := range items {
        score, _ := s.scoreRepo.GetByItemID(ctx, item.ID)
        summary, _ := s.summaryRepo.GetByItemID(ctx, item.ID)
        
        itemResp := map[string]interface{}{
            // ... existing fields ...
            "scores": map[string]float64{...},
        }
        
        if score != nil {
            itemResp["scores"] = map[string]float64{
                "final": score.Final,
                "hot": score.Hot,
                "relevance": score.Relevance,
                "credibility": score.Credibility,
                "novelty": score.Novelty,
            }
            
            // Add explanation if available
            if score.Explanation != nil {
                var explanation scoring.ScoreExplanation
                if err := json.Unmarshal(score.Explanation, &explanation); err == nil {
                    itemResp["score_explanation"] = explanation
                }
            }
        }
        
        // ... rest of code ...
    }
}
```

---

## G) Digest/Alerts Spec

### Scheduling Strategy

**Worker Architecture:**
```
DigestWorker runs every hour (cron: "0 * * * *")
1. Query all enabled digest_subscriptions
2. For each subscription:
   a. Convert user's send_at (local time) to UTC based on timezone
   b. Check if current UTC time matches scheduled time (±5 min window)
   c. If match and last_sent_at < scheduled_at:
      - Generate digest
      - Send email
      - Update last_sent_at
```

**File: `internal/services/digest_service.go`**

```go
func (s *DigestService) ProcessDigests(ctx context.Context) error {
    // Get all enabled subscriptions
    subscriptions, err := s.digestRepo.GetEnabledSubscriptions(ctx)
    if err != nil {
        return err
    }
    
    now := time.Now().UTC()
    
    for _, sub := range subscriptions {
        // Convert user's send_at (local) to UTC
        userTZ, err := time.LoadLocation(sub.Timezone)
        if err != nil {
            log.Printf("Invalid timezone %s for user %s", sub.Timezone, sub.UserID)
            continue
        }
        
        // Parse send_at as local time
        localTime := time.Date(now.Year(), now.Month(), now.Day(), 
            sub.SendAt.Hour(), sub.SendAt.Minute(), 0, 0, userTZ)
        scheduledUTC := localTime.UTC()
        
        // Check if we're in the send window (±5 minutes)
        if now.Before(scheduledUTC.Add(-5*time.Minute)) || 
           now.After(scheduledUTC.Add(5*time.Minute)) {
            continue
        }
        
        // Check if already sent today
        if sub.LastSentAt != nil {
            lastSent := sub.LastSentAt.In(userTZ)
            if lastSent.Year() == now.Year() && 
               lastSent.Month() == now.Month() && 
               lastSent.Day() == now.Day() {
                continue // Already sent today
            }
        }
        
        // Generate and send digest
        if err := s.generateAndSendDigest(ctx, sub); err != nil {
            log.Printf("Error sending digest for user %s: %v", sub.UserID, err)
            continue
        }
    }
    
    return nil
}
```

### Timezone Handling

**Strategy:**
1. Store user timezone (IANA format: "America/New_York")
2. Store `send_at` as TIME (HH:MM:SS) in user's local time
3. Convert to UTC at send time using Go's `time.LoadLocation()`
4. Handle DST automatically (Go time package handles this)

**Example:**
```go
// User in New York wants digest at 9:00 AM local
// send_at = "09:00:00", timezone = "America/New_York"

// At 9:00 AM EST (UTC-5), scheduledUTC = 14:00 UTC
// At 9:00 AM EDT (UTC-4), scheduledUTC = 13:00 UTC
// Go's time package handles DST automatically
```

### Rate Limits

**Per User:**
- Max 1 digest per frequency per day (daily = 1/day, weekly = 1/week)
- Max 5 keyword alerts per hour
- Max 10 keyword alerts per day

**Implementation:**
```go
// In digest_service.go
func (s *DigestService) canSendDigest(ctx context.Context, sub *models.DigestSubscription) (bool, error) {
    if sub.LastSentAt == nil {
        return true, nil
    }
    
    now := time.Now()
    lastSent := sub.LastSentAt
    
    if sub.Frequency == "daily" {
        // Check if sent in last 23 hours
        return now.Sub(lastSent) >= 23*time.Hour, nil
    } else if sub.Frequency == "weekly" {
        // Check if sent in last 6 days
        return now.Sub(lastSent) >= 6*24*time.Hour, nil
    }
    
    return false, nil
}
```

### Unsubscribe Mechanism

**Email Unsubscribe Link:**
```
https://app.example.com/unsubscribe?token=<jwt_token>
```

**Token Generation:**
```go
func (s *DigestService) GenerateUnsubscribeToken(userID uuid.UUID, subscriptionID uuid.UUID) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "subscription_id": subscriptionID.String(),
        "type": "unsubscribe",
        "exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.cfg.JWTSecret))
}
```

**Unsubscribe Endpoint:**
```
GET /v1/digests/unsubscribe?token=<token>
```

**Response:** Disables subscription and returns confirmation

---

## H) Observability

### Structured Logging

**Implementation: `internal/logging/logger.go`**

```go
package logging

import (
    "github.com/sirupsen/logrus"
    "github.com/google/uuid"
)

type Logger struct {
    *logrus.Logger
}

func New() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    return &Logger{Logger: logger}
}

func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
    return l.WithField("request_id", requestID)
}

func (l *Logger) WithUserID(userID uuid.UUID) *logrus.Entry {
    return l.WithField("user_id", userID.String())
}

func (l *Logger) WithSourceID(sourceID uuid.UUID) *logrus.Entry {
    return l.WithField("source_id", sourceID.String())
}

// Usage in handlers:
logger.WithRequestID(requestID).WithUserID(userID).Info("User accessed feed")
```

**Middleware: `internal/http/middleware/logging.go`**

```go
func RequestLoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method
        
        c.Next()
        
        latency := time.Since(start)
        status := c.Writer.Status()
        
        entry := logger.WithFields(logrus.Fields{
            "request_id": requestID,
            "method": method,
            "path": path,
            "status": status,
            "latency_ms": latency.Milliseconds(),
            "ip": c.ClientIP(),
        })
        
        if userID, exists := c.Get("user_id"); exists {
            entry = entry.WithField("user_id", userID)
        }
        
        if status >= 500 {
            entry.Error("Request failed")
        } else if status >= 400 {
            entry.Warn("Request error")
        } else {
            entry.Info("Request completed")
        }
    }
}
```

### Metrics List

**Prometheus Metrics:**

```go
// internal/metrics/metrics.go
var (
    HTTPRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
        },
        []string{"method", "path"},
    )
    
    IngestionItemsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ingestion_items_total",
            Help: "Total items ingested",
        },
        []string{"source_id", "status"}, // status: "inserted", "duplicate", "error"
    )
    
    IngestionDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "ingestion_duration_seconds",
            Help: "Ingestion duration per source",
            Buckets: []float64{1, 5, 10, 30, 60},
        },
        []string{"source_id"},
    )
    
    DedupRate = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "dedup_rate",
            Help: "Deduplication rate (0-1)",
        },
        []string{"source_id"},
    )
    
    DBQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Database query duration",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
        },
        []string{"operation"}, // "select", "insert", "update"
    )
    
    ActiveUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_total",
            Help: "Number of active users (logged in last 24h)",
        },
    )
)
```

### Tracing Approach

**OpenTelemetry Integration:**

```go
// internal/tracing/tracing.go
package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracing(serviceName string, endpoint string) (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return tp, nil
}

// Usage in handlers:
func (h *Handlers) GetFeed(c *gin.Context) {
    ctx, span := otel.Tracer("api").Start(c.Request.Context(), "GetFeed")
    defer span.End()
    
    // ... handler logic ...
}
```

---

## I) CI/CD

### GitHub Actions Workflow

**File: `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  pull_request:
    branches: [main, develop]
  push:
    branches: [main, develop]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: radar_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run migrations
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/radar_test?sslmode=disable
        run: |
          go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
          migrate -path migrations -database $DATABASE_URL up
      - name: Run tests
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/radar_test?sslmode=disable
        run: go test ./... -v -coverprofile=coverage.out
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build API
        run: go build -o bin/api ./cmd/api
      - name: Build Worker
        run: go build -o bin/worker ./cmd/worker

  migration-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check for destructive migrations
        run: |
          # Check if any migration contains DROP TABLE, DROP COLUMN, etc.
          if grep -r "DROP TABLE\|DROP COLUMN\|TRUNCATE" migrations/*.up.sql; then
            echo "ERROR: Destructive migration detected in PR"
            exit 1
          fi
```

**File: `.github/workflows/release.yml`**

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      - name: Build and push API image
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.api
          push: true
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/evolipia-radar-api:${{ github.ref_name }}
            ${{ secrets.DOCKER_USERNAME }}/evolipia-radar-api:latest
      - name: Build and push Worker image
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.worker
          push: true
          tags: |
            ${{ secrets.DOCKER_USERNAME }}/evolipia-radar-worker:${{ github.ref_name }}
            ${{ secrets.DOCKER_USERNAME }}/evolipia-radar-worker:latest
```

### Migration Safety Checks

**Pre-deployment Checklist:**
1. No DROP TABLE/COLUMN in PR migrations (checked in CI)
2. All migrations are reversible (down migration exists)
3. Migration tested on staging database
4. Backup created before production migration
5. Migration runs in transaction (PostgreSQL supports this)

**Migration Script: `scripts/check_migrations.sh`**

```bash
#!/bin/bash
# Check migration safety

echo "Checking migration safety..."

# Check for destructive operations
if grep -r "DROP TABLE\|DROP COLUMN\|TRUNCATE" migrations/*.up.sql; then
    echo "ERROR: Destructive migration detected"
    exit 1
fi

# Check all up migrations have down migrations
for up_file in migrations/*.up.sql; do
    down_file="${up_file%.up.sql}.down.sql"
    if [ ! -f "$down_file" ]; then
        echo "ERROR: Missing down migration for $up_file"
        exit 1
    fi
done

echo "Migration safety checks passed"
```

---

## J) Security Checklist

### SSRF Hardening

**Current:** Basic SSRF protection in `internal/security/ssrf.go`

**Enhancements:**
```go
// internal/security/ssrf.go (enhance)

func ValidateURL(url string) error {
    parsed, err := url.Parse(url)
    if err != nil {
        return fmt.Errorf("invalid URL format")
    }
    
    // Block private IPs
    if isPrivateIP(parsed.Hostname()) {
        return fmt.Errorf("private IP addresses not allowed")
    }
    
    // Block localhost variants
    hostname := strings.ToLower(parsed.Hostname())
    blockedHosts := []string{"localhost", "127.0.0.1", "0.0.0.0", "::1"}
    for _, blocked := range blockedHosts {
        if hostname == blocked {
            return fmt.Errorf("localhost not allowed")
        }
    }
    
    // Block link-local addresses
    if strings.HasPrefix(parsed.Hostname(), "169.254.") {
        return fmt.Errorf("link-local addresses not allowed")
    }
    
    // Only allow HTTP/HTTPS
    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return fmt.Errorf("only HTTP/HTTPS allowed")
    }
    
    return nil
}
```

### Per-Domain Rate Limits

**Implementation:** See PR 18 in Implementation Plan

**Configuration:**
- Default: 100 requests/hour per domain
- Configurable per source via `sources` table
- Stored in `rate_limits` table

### Auth Security

**Password Hashing:**
```go
// Use bcrypt with cost 10 (or argon2id)
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(hash), err
}

func VerifyPassword(hashedPassword, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
```

**JWT Rotation:**
- Access token: 1 hour expiry
- Refresh token: 30 days expiry
- Refresh tokens stored hashed in database
- On refresh, old token revoked, new token issued
- Max 5 active refresh tokens per user

**JWT Security:**
```go
// Use HS256 with strong secret (min 32 bytes)
// Include: user_id, exp, iat, jti (token ID)
// Validate: exp, iat, signature
```

### Input Validation

**Email Validation:**
```go
import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email) && len(email) <= 254
}
```

**Keyword Sanitization:**
```go
func SanitizeKeyword(keyword string) (string, error) {
    // Remove regex special characters
    sanitized := regexp.MustCompile(`[^a-zA-Z0-9\s\-_]`).ReplaceAllString(keyword, "")
    sanitized = strings.TrimSpace(sanitized)
    
    if len(sanitized) < 2 || len(sanitized) > 100 {
        return "", fmt.Errorf("keyword must be 2-100 characters")
    }
    
    return sanitized, nil
}
```

**Pagination Abuse Prevention:**
```go
// Max limit: 100
// Default limit: 20
func ValidatePagination(limit, offset int) (int, int, error) {
    if limit < 1 {
        limit = 20
    }
    if limit > 100 {
        return 0, 0, fmt.Errorf("limit cannot exceed 100")
    }
    if offset < 0 {
        offset = 0
    }
    if offset > 10000 {
        return 0, 0, fmt.Errorf("offset too large")
    }
    return limit, offset, nil
}
```

### Privacy Considerations

**User Data:**
1. **GDPR Compliance:**
   - Users can delete their account (soft delete)
   - Users can export their data (bookmarks, history, preferences)
   - Users can delete read history
   - Data retention: 90 days for read history (partial index)

2. **Data Minimization:**
   - Don't store full article content (only excerpt)
   - Don't store user IP addresses (only for rate limiting, 24h retention)
   - Don't log passwords (even hashed in logs)

3. **Access Control:**
   - Users can only access their own data
   - Admin endpoints separated (future: admin role)

4. **Encryption:**
   - Passwords hashed (bcrypt)
   - Sensitive data encrypted at rest (database encryption)
   - HTTPS only in production

**Privacy Endpoints:**
```
GET /v1/users/me/data-export - Export user data (JSON)
DELETE /v1/users/me - Delete account (soft delete)
DELETE /v1/history?all=true - Delete all read history
```

---

## Summary

This enhancement plan transforms evolipia-radar from a news aggregation service into a full user-facing app backend with:

✅ **Authentication & Profiles** - JWT-based auth with refresh tokens
✅ **Personalization** - User preferences and feedback-driven ranking
✅ **Bookmarks & History** - Save and track reading
✅ **Digests & Alerts** - Proactive engagement via email
✅ **Explainability** - Transparent score breakdowns
✅ **Observability** - Structured logs, metrics, tracing
✅ **CI/CD** - Automated testing and deployment
✅ **Security** - Hardened SSRF, rate limiting, input validation, privacy

All changes are incremental, backward-compatible, and follow existing architecture patterns.

