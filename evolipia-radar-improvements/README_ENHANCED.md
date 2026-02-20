<div align="center">

# 🎯 EVOLIPIA-RADAR

**E**ngineering **V**erified **O**verview of **L**atest **I**nsights, **P**riorities, **I**mpact & **A**nalytics

[![CI](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ci.yml/badge.svg)](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ci.yml)
[![CD](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/cd.yml/badge.svg)](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/cd.yml)
[![Security](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/security.yml/badge.svg)](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/security.yml)
[![ML Pipeline](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ml-pipeline.yml/badge.svg)](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ml-pipeline.yml)
[![codecov](https://codecov.io/gh/hidatara-ds/evolipia-radar/branch/main/graph/badge.svg)](https://codecov.io/gh/hidatara-ds/evolipia-radar)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidatara-ds/evolipia-radar)](https://goreportcard.com/report/github.com/hidatara-ds/evolipia-radar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<p align="center">
  <img src="docs/assets/architecture.png" alt="Architecture" width="800"/>
</p>

**A production-ready AI/ML tech news aggregator with MLOps best practices**

[📖 Documentation](https://docs.evolipia-radar.dev) • [🚀 Live Demo](https://evolipia-radar.dev) • [📊 MLflow](https://mlflow.evolipia-radar.dev)

</div>

---

## 🌟 Features

### Core Functionality
- **🔍 Multi-source Aggregation**: Hacker News, RSS/Atom feeds, arXiv, and custom JSON APIs
- **🧠 Intelligent Ranking**: ML-powered scoring combining popularity, relevance, credibility, and novelty
- **📝 Automatic Summarization**: AI-generated summaries with ML engineer-focused insights
- **🔄 Deduplication**: Prevents duplicate items across sources using URL normalization
- **⚡ Real-time Processing**: Temporal.io workflows for reliable background processing

### MLOps Features
- **📊 ML Experiment Tracking**: MLflow integration for model versioning and experiments
- **✅ Data Quality**: Great Expectations for data validation
- **🧪 A/B Testing**: Built-in experimentation framework
- **📈 Model Monitoring**: Prometheus metrics for model performance
- **🔄 Feature Store**: Centralized feature management

### Infrastructure
- **☸️ Kubernetes Native**: Helm charts with HPA, PDB, and network policies
- **🏗️ Infrastructure as Code**: Terraform modules for AWS/GCP/Azure
- **📡 Observability**: Prometheus, Grafana, Loki, Tempo, and Jaeger
- **🔒 Security First**: SSRF protection, rate limiting, secret scanning

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Web UI     │  │   Mobile     │  │    CLI       │  │   API Clients│   │
│  │  (React)     │  │  (PWA)       │  │   (Go)       │  │              │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │
└─────────┼─────────────────┼─────────────────┼─────────────────┼───────────┘
          │                 │                 │                 │
          └─────────────────┴────────┬────────┴─────────────────┘
                                     │
┌────────────────────────────────────┼────────────────────────────────────────┐
│                         API GATEWAY │ LAYER                                 │
│  ┌──────────────────────────────────┴──────────────────────────────────┐   │
│  │                         Kong / Nginx Ingress                         │   │
│  │  (Rate Limiting, Auth, SSL Termination, Request Routing)            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌────────────────────────────────────┼────────────────────────────────────────┐
│                         SERVICE LAYER (Kubernetes)                          │
│                                                                              │
│  ┌─────────────────────────────────┐    ┌─────────────────────────────────┐ │
│  │         API Server              │    │           Worker                │ │
│  │    ┌─────────────────────┐      │    │    ┌─────────────────────┐      │ │
│  │    │   REST API (Gin)    │      │    │    │  Temporal Worker    │      │ │
│  │    │  - /v1/feed         │      │    │    │  - Ingestion        │      │ │
│  │    │  - /v1/search       │      │    │    │  - Scoring          │      │ │
│  │    │  - /v1/sources      │      │    │    │  - Summarization    │      │ │
│  │    │  - /healthz         │      │    │    │  - Deduplication    │      │ │
│  │    └─────────────────────┘      │    │    └─────────────────────┘      │ │
│  │                                 │    │                                 │ │
│  │  ┌─────────────────────────┐   │    │  ┌─────────────────────────┐   │ │
│  │  │   OpenTelemetry SDK     │   │    │  │   OpenTelemetry SDK     │   │ │
│  │  │  (Metrics, Logs, Traces)│   │    │  │  (Metrics, Logs, Traces)│   │ │
│  │  └─────────────────────────┘   │    │  └─────────────────────────┘   │ │
│  └─────────────────────────────────┘    └─────────────────────────────────┘ │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Temporal Server (Workflow Engine)                 │   │
│  │         (Reliable background processing, retries, scheduling)        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌────────────────────────────────────┼────────────────────────────────────────┐
│                         DATA LAYER                                           │
│                                                                              │
│  ┌─────────────────────────┐    ┌─────────────────────────────────────────┐ │
│  │    PostgreSQL           │    │         Redis (Cache)                   │ │
│  │  - content_items        │    │    - Feed cache                         │ │
│  │  - sources              │    │    - Session store                      │ │
│  │  - scores               │    │    - Rate limiting                      │ │
│  │  - embeddings           │    │                                         │ │
│  └─────────────────────────┘    └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
┌────────────────────────────────────┼────────────────────────────────────────┐
│                      ML & OBSERVABILITY LAYER                               │
│                                                                              │
│  ┌─────────────────────────┐    ┌─────────────────────────────────────────┐ │
│  │       MLflow            │    │      Observability Stack                │ │
│  │  - Experiment Tracking  │    │  ┌─────────┐ ┌─────────┐ ┌─────────┐  │ │
│  │  - Model Registry       │    │  │Prometheus│ │  Loki   │ │  Tempo  │  │ │
│  │  - Artifact Store       │    │  │Metrics  │ │  Logs   │ │ Traces  │  │ │
│  │  - Model Comparison     │    │  └────┬────┘ └────┬────┘ └────┬────┘  │ │
│  └─────────────────────────┘    │       └───────────┴───────────┘       │ │
│                                 │                  │                      │ │
│  ┌─────────────────────────┐    │            ┌─────────┐                 │ │
│  │  Great Expectations     │    │            │ Grafana │                 │ │
│  │  - Data Validation      │    │            │Dashboard│                 │ │
│  │  - Quality Reports      │    │            └─────────┘                 │ │
│  └─────────────────────────┘    └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- kubectl (for Kubernetes deployment)
- Terraform 1.5+ (for infrastructure)

### Local Development

```bash
# Clone the repository
git clone https://github.com/hidatara-ds/evolipia-radar.git
cd evolipia-radar

# Setup development environment
make setup

# Start infrastructure services
make docker-compose-up

# Run database migrations
make migrate-up

# Start API server (terminal 1)
make run-api

# Start worker (terminal 2)
make run-worker

# Access the web UI
open http://localhost:8080
```

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run with coverage
make test-coverage

# Run benchmarks
make test-benchmark
```

---

## 📦 Deployment

### Docker Compose (Development)

```bash
docker-compose up -d
```

### Kubernetes (Production)

```bash
# Install with Helm
helm repo add evolipia-radar https://charts.evolipia-radar.dev
helm install evolipia-radar evolipia-radar/evolipia-radar   --namespace evolipia-radar   --create-namespace   --values values-production.yaml
```

### Terraform (Infrastructure)

```bash
cd terraform/environments/production
terraform init
terraform plan
terraform apply
```

---

## 📊 Monitoring & Observability

Access the observability stack:

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| Jaeger | http://localhost:16686 | - |
| MLflow | http://localhost:5000 | - |

### Key Metrics

- **API Latency**: p50, p95, p99 request latency
- **Error Rate**: 4xx and 5xx error percentages
- **Throughput**: Requests per second
- **ML Pipeline**: Model inference time, accuracy
- **Data Quality**: Validation success rate

---

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://localhost/radar` |
| `TEMPORAL_HOST` | Temporal server address | `localhost:7233` |
| `MLFLOW_TRACKING_URI` | MLflow tracking URI | `http://localhost:5000` |
| `WORKER_CRON` | Ingestion schedule | `*/10 * * * *` |
| `LOG_LEVEL` | Logging level | `info` |
| `METRICS_ENABLED` | Enable Prometheus metrics | `true` |

### Scoring Configuration

```yaml
# internal/scoring/config.go
weights:
  popularity: 0.55   # HN points, comments
  relevance: 0.25    # AI/ML keyword matching
  credibility: 0.15  # Domain reputation
  novelty: 0.05      # Recency score
```

---

## 🧪 ML Pipeline

### Training a New Model

```bash
# Run data quality checks
make ml-data-quality

# Train the model
make ml-train

# Evaluate the model
make ml-evaluate
```

### Experiment Tracking

All experiments are automatically logged to MLflow:

```python
import mlflow

with mlflow.start_run():
    mlflow.log_param("model_type", "xgboost")
    mlflow.log_metric("accuracy", 0.95)
    mlflow.sklearn.log_model(model, "model")
```

---

## 🔒 Security

- **SSRF Protection**: Blocks private IP ranges
- **Rate Limiting**: Per-source and per-IP limits
- **Input Validation**: Strict URL and JSON validation
- **Secret Scanning**: Automated detection in CI/CD
- **Container Scanning**: Trivy vulnerability scans
- **Dependency Updates**: Automated via Dependabot

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Temporal.io](https://temporal.io/)
- [MLflow](https://mlflow.org/)
- [Great Expectations](https://greatexpectations.io/)

---

<div align="center">

**Made with ❤️ for the MLOps Community**

[⭐ Star this repo](https://github.com/hidatara-ds/evolipia-radar) • [🐛 Report Bug](https://github.com/hidatara-ds/evolipia-radar/issues) • [💡 Request Feature](https://github.com/hidatara-ds/evolipia-radar/issues)

</div>
