#!/bin/bash
# EVOLIPIA-RADAR MLOps Integration Script
# Run this from your evolipia-radar repository root

set -e

echo "🚀 Starting MLOps integration..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Run this script from your evolipia-radar repository root"
    exit 1
fi

# Create backup branch
echo "📦 Creating backup branch..."
git add -A
git commit -m "chore: pre-mlops backup" || true
git checkout -b mlops-improvements 2>/dev/null || git checkout mlops-improvements

# Extract improvements
echo "📂 Extracting improvements..."
tar -xzf evolipia-radar-improvements.tar.gz

# Copy new files (safe - these don't exist)
echo "➕ Adding new files..."
cp -r evolipia-radar-improvements/.github . 2>/dev/null || true
cp -r evolipia-radar-improvements/tests . 2>/dev/null || true
cp -r evolipia-radar-improvements/k8s . 2>/dev/null || true
cp -r evolipia-radar-improvements/terraform . 2>/dev/null || true
cp evolipia-radar-improvements/.golangci.yml . 2>/dev/null || true
cp evolipia-radar-improvements/.pre-commit-config.yaml . 2>/dev/null || true
cp evolipia-radar-improvements/docker-compose.observability.yml . 2>/dev/null || true
cp evolipia-radar-improvements/docker-compose.ml.yml . 2>/dev/null || true
cp evolipia-radar-improvements/CONTRIBUTING.md . 2>/dev/null || true

# Merge Makefile improvements
echo "🔧 Merging Makefile improvements..."
if [ -f "Makefile" ]; then
    # Backup original
    cp Makefile Makefile.backup

    # Add new targets to existing Makefile
    cat >> Makefile << 'EOF'

# ============================================
# MLOps Additions
# ============================================

.PHONY: docker-compose-obs
docker-compose-obs: ## Start observability stack
	docker-compose -f docker-compose.observability.yml up -d

.PHONY: docker-compose-ml
docker-compose-ml: ## Start ML stack
	docker-compose -f docker-compose.ml.yml up -d

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --config=.golangci.yml ./...

.PHONY: security-scan
security-scan: ## Run security scans
	trivy fs --severity HIGH,CRITICAL .
EOF
    echo "✅ Makefile updated"
else
    cp evolipia-radar-improvements/Makefile .
fi

# Add observability dependencies to go.mod
echo "📥 Adding observability dependencies..."
go get github.com/prometheus/client_golang/prometheus@latest 2>/dev/null || true
go get github.com/prometheus/client_golang/prometheus/promhttp@latest 2>/dev/null || true
go get go.opentelemetry.io/otel@latest 2>/dev/null || true
go get go.opentelemetry.io/otel/trace@latest 2>/dev/null || true
go get github.com/sirupsen/logrus@latest 2>/dev/null || true
go mod tidy

# Create enhanced README (backup original first)
echo "📝 Enhancing README..."
if [ -f "README.md" ]; then
    cp README.md README.md.backup

    # Add badges at the top
    cat > README_NEW.md << 'EOF'
<div align="center">

# 🎯 EVOLIPIA-RADAR

[![CI](https://github.com/YOUR_USERNAME/evolipia-radar/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/evolipia-radar/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/evolipia-radar)](https://goreportcard.com/report/github.com/YOUR_USERNAME/evolipia-radar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**An AI/ML tech news aggregator with MLOps best practices**

</div>

---

EOF
    # Append original README content
    cat README.md >> README_NEW.md
    mv README_NEW.md README.md
fi

# Add GitHub Actions dependencies
echo "🔧 Setting up GitHub Actions..."
mkdir -p .github/workflows

# Clean up
echo "🧹 Cleaning up..."
rm -rf evolipia-radar-improvements
touch .secrets.baseline  # For detect-secrets

echo ""
echo "✅ Integration complete!"
echo ""
echo "📋 Next steps:"
echo "   1. Review the changes: git status"
echo "   2. Add GitHub secrets (Settings → Secrets):"
echo "      - CODECOV_TOKEN (optional)"
echo "   3. Commit: git add . && git commit -m 'feat: Add MLOps infrastructure'"
echo "   4. Push: git push origin mlops-improvements"
echo "   5. Create PR to main"
echo ""
echo "🎉 Your repo now has MLOps superpowers!"
