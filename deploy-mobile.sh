#!/bin/bash
# Quick deployment script untuk mobile PWA

set -e

echo "🚀 Evolipia Radar - Mobile Deployment"
echo "======================================"
echo ""

# Check if git is clean
if [[ -n $(git status -s) ]]; then
    echo "⚠️  Ada perubahan yang belum di-commit"
    read -p "Commit sekarang? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git add .
        read -p "Commit message: " commit_msg
        git commit -m "$commit_msg"
    fi
fi

echo ""
echo "Pilih platform deployment:"
echo "1) Fly.io (Recommended - Fast & Free)"
echo "2) Render.com (Easy - Auto deploy from GitHub)"
echo "3) Railway.app (Simple - Good free tier)"
echo ""
read -p "Pilihan (1-3): " choice

case $choice in
    1)
        echo ""
        echo "📦 Deploying to Fly.io..."
        echo ""
        
        # Check if flyctl installed
        if ! command -v flyctl &> /dev/null; then
            echo "❌ flyctl belum terinstall"
            echo "Install: https://fly.io/docs/hands-on/install-flyctl/"
            exit 1
        fi
        
        # Check if logged in
        if ! flyctl auth whoami &> /dev/null; then
            echo "🔐 Login ke Fly.io..."
            flyctl auth login
        fi
        
        # Check if app exists
        if ! flyctl status &> /dev/null; then
            echo "🆕 Membuat app baru..."
            flyctl launch --no-deploy
        fi
        
        # Create/attach Postgres
        echo "🗄️  Setup database..."
        if ! flyctl postgres list | grep -q "evolipia-radar-db"; then
            flyctl postgres create --name evolipia-radar-db --region sin
            flyctl postgres attach evolipia-radar-db
        fi
        
        # Deploy
        echo "🚀 Deploying..."
        flyctl deploy
        
        # Run migrations
        echo "📊 Running migrations..."
        flyctl ssh console -C "migrate -path /app/migrations -database \$DATABASE_URL up"
        
        echo ""
        echo "✅ Deployment selesai!"
        flyctl open
        ;;
        
    2)
        echo ""
        echo "📦 Setup Render.com deployment..."
        echo ""
        echo "Langkah manual:"
        echo "1. Push code ke GitHub: git push"
        echo "2. Buka https://render.com"
        echo "3. New > Blueprint"
        echo "4. Connect repository ini"
        echo "5. Render akan auto-detect render.yaml"
        echo "6. Klik 'Apply'"
        echo ""
        echo "Atau gunakan Render CLI:"
        echo "  npm install -g @render/cli"
        echo "  render login"
        echo "  render blueprint launch"
        echo ""
        read -p "Push ke GitHub sekarang? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git push
            echo "✅ Pushed! Sekarang setup di render.com"
        fi
        ;;
        
    3)
        echo ""
        echo "📦 Deploying to Railway..."
        echo ""
        
        # Check if railway installed
        if ! command -v railway &> /dev/null; then
            echo "Installing Railway CLI..."
            npm install -g @railway/cli
        fi
        
        # Login
        if ! railway whoami &> /dev/null; then
            echo "🔐 Login ke Railway..."
            railway login
        fi
        
        # Init if needed
        if [ ! -f "railway.json" ]; then
            echo "🆕 Membuat project baru..."
            railway init
        fi
        
        # Add Postgres
        echo "🗄️  Adding PostgreSQL..."
        railway add --plugin postgresql
        
        # Deploy
        echo "🚀 Deploying..."
        railway up
        
        # Get database URL and run migrations
        echo "📊 Running migrations..."
        DB_URL=$(railway variables get DATABASE_URL)
        migrate -path migrations -database "$DB_URL" up
        
        echo ""
        echo "✅ Deployment selesai!"
        railway open
        ;;
        
    *)
        echo "❌ Pilihan tidak valid"
        exit 1
        ;;
esac

echo ""
echo "🎉 Deployment selesai!"
echo ""
echo "📱 Test di HP:"
echo "1. Buka URL aplikasi di browser HP"
echo "2. Chrome: Menu (⋮) > Add to Home Screen"
echo "3. Safari: Share > Add to Home Screen"
echo ""
echo "🔧 Jangan lupa:"
echo "- Set environment variables di dashboard platform"
echo "- Update API_BASE di web/index.html kalau API beda domain"
echo "- Test PWA dengan Lighthouse di Chrome DevTools"
echo ""
