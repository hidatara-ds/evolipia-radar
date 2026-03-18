# Contributing to EVOLIPIA-RADAR

First off, thank you for considering contributing to EVOLIPIA-RADAR! 🎉

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce**
- **Provide specific examples**
- **Describe the behavior you observed**
- **Explain which behavior you expected**
- **Include code samples and screenshots**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. Create an issue and provide:

- **Use a clear and descriptive title**
- **Provide a step-by-step description**
- **Provide specific examples**
- **Explain why this enhancement would be useful**

### Pull Requests

1. Fork the repository
2. Create a branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linters (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make
- golangci-lint

### Setup

```bash
# Clone your fork
git clone https://github.com/your-username/evolipia-radar.git
cd evolipia-radar

# Install dependencies
make setup

# Setup pre-commit hooks
make setup-hooks

# Start services
docker-compose up -d postgres temporal

# Run migrations
make migrate-up
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./internal/scoring/... -run TestCalculateScore
```

### Code Style

We use `golangci-lint` for linting. Run it before committing:

```bash
make lint
```

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Test changes
- `chore:` Build process or auxiliary tool changes

Example:
```
feat(scoring): add ML-based relevance scoring

- Implement XGBoost model for content relevance
- Add feature engineering pipeline
- Update tests
```

## Project Structure

```
evolipia-radar/
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── config/       # Configuration
│   ├── connectors/   # Data source connectors
│   ├── db/           # Database repositories
│   ├── dto/          # Data transfer objects
│   ├── http/         # HTTP handlers
│   ├── models/       # Domain models
│   ├── scoring/      # Scoring algorithms
│   ├── services/     # Business logic
│   └── ...
├── tests/            # Test files
├── k8s/              # Kubernetes manifests
├── terraform/        # Infrastructure code
└── docs/             # Documentation
```

## Testing Guidelines

### Unit Tests

- Test one thing per test
- Use table-driven tests
- Mock external dependencies
- Aim for >80% coverage

```go
func TestCalculateScore(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected float64
    }{
        {"zero", 0, 0.0},
        {"positive", 10, 0.5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateScore(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Tests

- Use build tag `//go:build integration`
- Test with real database
- Clean up test data

### Benchmarks

- Include benchmarks for performance-critical code
- Use `Benchmark` prefix

## Documentation

- Update README.md if adding features
- Add Go doc comments for exported functions
- Update architecture diagrams

## Questions?

Feel free to open an issue with your question or join our [Discord community](https://discord.gg/evolipia-radar).

Thank you for contributing! 🚀
