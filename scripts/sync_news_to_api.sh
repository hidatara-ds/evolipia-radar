#!/bin/bash
# Script to sync news.json from data/ to api/ folder
# This ensures Vercel can bundle the file with Go functions

echo "🔄 Syncing news.json to api folder..."

if [ -f "data/news.json" ]; then
    cp data/news.json api/news.json
    echo "✅ Copied data/news.json to api/news.json"
    
    # Show file size
    size=$(wc -c < api/news.json)
    echo "📦 File size: $size bytes"
else
    echo "❌ data/news.json not found!"
    exit 1
fi
