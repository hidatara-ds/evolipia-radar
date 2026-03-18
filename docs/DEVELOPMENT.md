# Development Guide

Complete guide for local development of Evolipia Radar.

## Prerequisites

### Required

- **Go:** 1.24.1 or higher
- **Node.js:** 18.x or higher
- **PostgreSQL:** 15+ (or Neon.tech account)
- **Git:** Latest version

### Optional

- **Docker:** For local PostgreSQL
- **Make:** For build automation
- **golangci-lint:** For Go linting

## Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/hidatara-ds/evolipia-radar.git
cd evolipia-radar
```

### 2. Install Dependencies

**Backend (Go):**
```bash
go mod download
go mod verify
```

**Frontend (Node.js):**
```bash
npm install
```

### 3. Database Setup

**Option A: Use Neon.tech (Recommended)**

1. Create account at [neon.tech](https://neon.tech)
2. Create new project
3. Copy connection string
4. Run migrations:
   ```bash
   psql "your-connection-string" < migrations/001_initial_schema.sql
   ```

**Option B: Local PostgreSQL**

```bash
# Start PostgreSQL with Docker
docker run --name evolipia-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=radar \
  -p 5432:5432 \
  -d postgres:15

# Run migrations
psql postgresql://postgres:postgres@localhost:5432/radar < migrations/001_initial_schema.sql
```

### 4. Environment Variables

Create `.env.local`:

```env
# Database
DATABASE_URL=postgresql://user:pass@host/dbname

# LLM (Optional)
LLM_API_KEY=your_openrouter_key
LLM_PROVIDER=openrouter
LLM_MODEL=google/gemini-flash-1.5
LLM_ENABLED=true
LLM_MAX_TOKENS=500
LLM_TEMPERATURE=0.7

# Worker
WORKER_CRON=0 */30 * * * *
OUTPUT_JSON=true
JSON_OUTPUT_PATH=data/news.json
```

### 5. Populate Sample Data

```bash
export DATABASE_URL="your-connection-string"
go run scripts/populate_sample_data.go
```

## Development Workflow

### Running Services

**Terminal 1: Frontend**
```bash
npm run dev
# Runs on http://localhost:3000
```

**Terminal 2: API Server (Optional)**
```bash
go run cmd/api/main.go
# Runs on http://localhost:8080
```

**Terminal 3: Worker (Optional)**
```bash
go run cmd/worker/main.go
# Runs scraper with configured cron
```

### Hot Reload

- **Frontend:** Auto-reloads on file changes
- **Backend:** Restart manually or use `air`:
  ```bash
  go install github.com/cosmtrek/air@latest
  air
  ```

## Project Structure

```
evolipia-radar/
├── api/                    # Vercel serverless functions
│   ├── health/
│   ├── metrics/
│   ├── news/              # Main news API
│   ├── search/
│   ├── trending/
│   └── trigger/
├── app/                    # Next.js app directory
│   ├── globals.css
│   ├── layout.tsx
│   └── page.tsx           # Main dashboard
├── cmd/                    # Go command-line tools
│   ├── api/               # API server
│   ├── worker/            # Background worker
│   └── worker-json/       # JSON export worker
├── pkg/                    # Go packages
│   ├── config/            # Configuration
│   ├── crawler/           # Web crawling
│   ├── db/                # Database layer
│   ├── models/            # Data models
│   ├── scoring/           # Scoring algorithms
│   ├── services/          # Business logic
│   └── tagging/           # Auto-tagging
├── scripts/                # Utility scripts
│   ├── add_news_sources.go
│   ├── populate_sample_data.go
│   └── retag_news.go
├── migrations/             # Database migrations
├── docs/                   # Documentation
├── public/                 # Static assets
├── .github/workflows/      # GitHub Actions
└── tests/                  # Test files
```

## Code Style

### Go

Follow [Effective Go](https://golang.org/doc/effective_go.html):

```go
// Good
func LoadNewsData() (*NewsData, error) {
    // Implementation
}

// Bad
func load_news_data() (*NewsData, error) {
    // Implementation
}
```

**Formatting:**
```bash
go fmt ./...
goimports -w .
```

**Linting:**
```bash
golangci-lint run
```

### TypeScript/React

Follow [Airbnb Style Guide](https://github.com/airbnb/javascript):

```typescript
// Good
const fetchNews = async (topic?: string): Promise<NewsItem[]> => {
  // Implementation
};

// Bad
function fetchNews(topic) {
  // Implementation
}
```

**Formatting:**
```bash
npm run format
```

**Linting:**
```bash
npm run lint
```

## Testing

### Go Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/tagging/...

# Verbose output
go test -v ./...
```

**Writing Tests:**

```go
// pkg/tagging/auto_tagger_test.go
package tagging

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAutoTagger_TagContent(t *testing.T) {
    tagger := NewAutoTagger()
    
    tags := tagger.TagContent("GPT-4 released", "OpenAI announces GPT-4")
    
    assert.Contains(t, tags, "llm")
    assert.NotContains(t, tags, "vision")
}
```

### Frontend Tests

```bash
# Run tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage
```

**Writing Tests:**

```typescript
// app/__tests__/page.test.tsx
import { render, screen } from '@testing/library/react';
import Dashboard from '../page';

describe('Dashboard', () => {
  it('renders dashboard title', () => {
    render(<Dashboard />);
    expect(screen.getByText('Evolipia Radar')).toBeInTheDocument();
  });
});
```

### Integration Tests

```bash
# Run integration tests
npm run test:e2e
```

## Debugging

### Go Debugging

**Using Delve:**

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug API
dlv debug cmd/api/main.go

# Set breakpoint
(dlv) break main.main
(dlv) continue
```

**VS Code:**

```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug API",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/api/main.go"
    }
  ]
}
```

### Frontend Debugging

**Browser DevTools:**
- Open Chrome DevTools (F12)
- Sources tab → Set breakpoints
- Console tab → View logs

**VS Code:**

```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Next.js: debug",
      "type": "node-terminal",
      "request": "launch",
      "command": "npm run dev"
    }
  ]
}
```

## Database Development

### Migrations

**Create Migration:**

```bash
# Create new migration file
touch migrations/002_add_bookmarks.sql
```

**Apply Migration:**

```bash
psql $DATABASE_URL < migrations/002_add_bookmarks.sql
```

**Rollback (Manual):**

```sql
-- migrations/002_add_bookmarks_down.sql
DROP TABLE IF EXISTS bookmarks;
```

### Database Tools

**psql:**

```bash
# Connect
psql $DATABASE_URL

# List tables
\dt

# Describe table
\d items

# Run query
SELECT COUNT(*) FROM items;
```

**pgAdmin:**

1. Download from [pgadmin.org](https://www.pgadmin.org/)
2. Add server with DATABASE_URL
3. Browse tables, run queries

## Common Tasks

### Add New News Source

1. **Add to database:**
   ```sql
   INSERT INTO sources (name, type, category, url, enabled)
   VALUES ('New Source', 'rss', 'tech', 'https://example.com/feed', true);
   ```

2. **Test scraping:**
   ```bash
   go run cmd/worker/main.go
   ```

### Add New Topic Tag

1. **Update auto_tagger.go:**
   ```go
   func (t *AutoTagger) TagContent(title, content string) []string {
       // Add new tag logic
       if containsKeywords(text, []string{"quantum", "qbit"}) {
           tags = append(tags, "quantum")
       }
   }
   ```

2. **Update frontend:**
   ```typescript
   // app/page.tsx
   const TOPICS = [
       // ...
       { id: "quantum", label: "Quantum", color: "indigo" },
   ];
   ```

### Add New API Endpoint

1. **Create handler:**
   ```go
   // api/bookmarks/index.go
   package bookmarks
   
   func Handler(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```

2. **Update vercel.json:**
   ```json
   {
     "routes": [
       {
         "src": "/api/bookmarks",
         "dest": "/api/bookmarks/index.go"
       }
     ]
   }
   ```

3. **Test locally:**
   ```bash
   vercel dev
   curl http://localhost:3000/api/bookmarks
   ```

## Performance Profiling

### Go Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Frontend Profiling

**React DevTools:**
1. Install React DevTools extension
2. Open Profiler tab
3. Record interaction
4. Analyze render times

**Lighthouse:**
```bash
npm run build
npm start
# Open Chrome DevTools → Lighthouse → Run audit
```

## Troubleshooting

### Common Issues

**"go: module not found"**
```bash
go mod tidy
go mod download
```

**"npm ERR! peer dependencies"**
```bash
npm install --legacy-peer-deps
```

**"Database connection refused"**
- Check DATABASE_URL
- Verify database is running
- Check network/firewall

**"Port already in use"**
```bash
# Find process
lsof -i :3000
# Kill process
kill -9 <PID>
```

### Getting Help

1. Check [Documentation](https://github.com/hidatara-ds/evolipia-radar/tree/main/docs)
2. Search [GitHub Issues](https://github.com/hidatara-ds/evolipia-radar/issues)
3. Create new issue with:
   - Environment details
   - Steps to reproduce
   - Error messages
   - Expected vs actual behavior

## Best Practices

### Code Quality

- Write tests for new features
- Keep functions small and focused
- Use meaningful variable names
- Add comments for complex logic
- Handle errors properly

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes
git add .
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/new-feature
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new topic filter
fix: resolve database connection issue
docs: update API documentation
refactor: simplify scoring algorithm
test: add tests for auto-tagger
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
