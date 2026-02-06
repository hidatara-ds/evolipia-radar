# Dependencies

This document lists all external dependencies used in this project.

## Core Dependencies

### HTTP Framework
- **github.com/gin-gonic/gin** - High-performance HTTP web framework
  - Used for: REST API endpoints, request routing, middleware
  - License: MIT

### Database
- **github.com/jackc/pgx/v5** - PostgreSQL driver and toolkit
  - Used for: Database connection pooling, query execution
  - License: MIT

### Utilities
- **github.com/google/uuid** - UUID package for Go
  - Used for: Generating unique identifiers for entities
  - License: BSD-3-Clause

### Scheduling
- **github.com/robfig/cron/v3** - Cron library for Go
  - Used for: Worker scheduling and periodic tasks
  - License: MIT

## Development Dependencies

### Database Migrations
- **github.com/golang-migrate/migrate** - Database migration tool
  - Used for: Database schema versioning and migrations
  - License: MIT

### Linting
- **github.com/golangci/golangci-lint** - Fast linter for Go
  - Used for: Code quality checks and static analysis
  - License: GPL-3.0

## Dependency Management

This project uses Go modules for dependency management.

### Commands

```bash
# Download all dependencies
go mod download

# Add a new dependency
go get <package>

# Update dependencies
go get -u ./...

# Remove unused dependencies
go mod tidy

# View dependency graph
go mod graph

# List all dependencies
go list -m all

# Verify dependencies
go mod verify
```

### Updating Dependencies

To update all dependencies to their latest versions:

```bash
go get -u ./...
go mod tidy
```

To update a specific dependency:

```bash
go get -u github.com/gin-gonic/gin@latest
```

## Security

Regularly check for security vulnerabilities:

```bash
go list -json -m all | nancy sleuth
# or use go-audit
```

## License Compatibility

All dependencies use permissive licenses (MIT, BSD) that are compatible with this project's MIT license.

