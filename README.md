# Evolipia Radar

> AI Research Intelligence Platform - Real-time tracking and analysis of AI/ML developments

[![Production](https://img.shields.io/badge/status-production-success)](https://evolipia-radar.vercel.app)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)
[![Go Version](https://img.shields.io/badge/go-1.24.1-00ADD8.svg)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/next.js-14.2.0-000000.svg)](https://nextjs.org/)

## 🎯 Overview

Evolipia Radar is an intelligent news aggregation and analysis platform designed for AI researchers and engineers. It automatically discovers, scores, and categorizes the latest developments in artificial intelligence and machine learning.

**Live Demo:** [https://evolipia-radar.vercel.app](https://evolipia-radar.vercel.app)

### Key Features

- 🔍 **Intelligent Crawling** - Automated discovery from 40+ curated AI/ML sources
- 🏷️ **Smart Tagging** - Auto-categorization into LLM, Vision, RL, Robotics, IDE, and more
- 📊 **Relevance Scoring** - Multi-factor scoring algorithm for content quality
- 🎨 **Modern Dashboard** - Real-time updates with topic filtering
- 🚀 **Production Ready** - Deployed on Vercel with PostgreSQL backend

## 🏗️ Architecture

```
┌─────────────────┐
│  GitHub Actions │ ──► Scrapes news every 30min
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Neon.tech DB   │ ──► PostgreSQL database
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Vercel API     │ ──► Serverless Go functions
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Next.js UI     │ ──► React dashboard
└─────────────────┘
```

**Tech Stack:**
- **Backend:** Go 1.24.1, PostgreSQL (Neon.tech)
- **Frontend:** Next.js 14, React, TypeScript, Tailwind CSS
- **Deployment:** Vercel (Frontend + API), GitHub Actions (Scraper)
- **Database:** Neon.tech PostgreSQL with connection pooling

## 🚀 Quick Start

### Prerequisites

- Go 1.24.1+
- Node.js 18+
- PostgreSQL database (or Neon.tech account)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/hidatara-ds/evolipia-radar.git
   cd evolipia-radar
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env.local
   ```
   
   Edit `.env.local`:
   ```env
   DATABASE_URL=postgresql://user:pass@host/dbname
   LLM_API_KEY=your_openrouter_key
   LLM_PROVIDER=openrouter
   LLM_MODEL=google/gemini-flash-1.5
   ```

3. **Install dependencies**
   ```bash
   # Backend
   go mod download
   
   # Frontend
   npm install
   ```

4. **Run database migrations**
   ```bash
   # Apply schema
   psql $DATABASE_URL < migrations/001_initial_schema.sql
   ```

5. **Start development servers**
   ```bash
   # Terminal 1: Frontend
   npm run dev
   
   # Terminal 2: API (optional for local testing)
   go run cmd/api/main.go
   
   # Terminal 3: Worker (optional for local scraping)
   go run cmd/worker/main.go
   ```

6. **Access the application**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080

## 📚 Documentation

Comprehensive documentation is available in the [`docs/`](docs/) directory:

- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and components
- **[API Documentation](docs/API.md)** - REST API endpoints and schemas
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment instructions
- **[Development Guide](docs/DEVELOPMENT.md)** - Local development setup
- **[Database Schema](docs/DATABASE.md)** - PostgreSQL schema and migrations
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `LLM_API_KEY` | OpenRouter API key for AI features | No | - |
| `LLM_PROVIDER` | LLM provider (openrouter) | No | openrouter |
| `LLM_MODEL` | Model to use | No | google/gemini-flash-1.5 |
| `LLM_ENABLED` | Enable AI summarization | No | false |

### Topic Filters

The platform supports the following topic categories:

- **LLM** - Large Language Models
- **Vision** - Computer Vision & Image Generation
- **Data** - Data Science & Analytics
- **Security** - AI Security & Privacy
- **RL** - Reinforcement Learning
- **Robotics** - Robotics & Automation
- **IDE** - Developer Tools & IDEs
- **Free Credits** - Student Programs & Free Resources

## 🛠️ Development

### Project Structure

```
evolipia-radar/
├── api/              # Vercel serverless functions (Go)
├── app/              # Next.js app directory
├── cmd/              # Go command-line tools
│   ├── api/          # API server
│   ├── worker/       # Background scraper
│   └── worker-json/  # JSON export worker
├── pkg/              # Go packages
│   ├── db/           # Database layer
│   ├── models/       # Data models
│   ├── services/     # Business logic
│   └── tagging/      # Auto-tagging system
├── scripts/          # Utility scripts
├── migrations/       # Database migrations
├── docs/             # Documentation
└── public/           # Static assets
```

### Running Tests

```bash
# Go tests
go test ./...

# Frontend tests
npm test

# Integration tests
npm run test:e2e
```

### Code Quality

```bash
# Linting
golangci-lint run
npm run lint

# Formatting
go fmt ./...
npm run format
```

## 🚢 Deployment

### Vercel (Recommended)

1. **Connect to Vercel**
   ```bash
   vercel
   ```

2. **Set environment variables**
   ```bash
   vercel env add DATABASE_URL
   vercel env add LLM_API_KEY
   ```

3. **Deploy**
   ```bash
   vercel --prod
   ```

See [Deployment Guide](docs/DEPLOYMENT.md) for detailed instructions.

## 📊 Monitoring

- **Application Logs:** Vercel Dashboard → Logs
- **Database Metrics:** Neon.tech Dashboard
- **GitHub Actions:** Repository → Actions tab

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 🙏 Acknowledgments

- Built with [Next.js](https://nextjs.org/)
- Powered by [Neon.tech](https://neon.tech/)
- Deployed on [Vercel](https://vercel.com/)
- AI features by [OpenRouter](https://openrouter.ai/)

## 📧 Contact

- **Website:** [evolipia-radar.vercel.app](https://evolipia-radar.vercel.app)
- **GitHub:** [@hidatara-ds](https://github.com/hidatara-ds)
- **Issues:** [GitHub Issues](https://github.com/hidatara-ds/evolipia-radar/issues)

---

Made with ❤️ by the Evolipia team
