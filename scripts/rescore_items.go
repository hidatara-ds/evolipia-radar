package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/scoring"
)

// Script to re-score all items in the database
// This is useful after changing scoring algorithm or weights
func main() {
	fmt.Println("=== Re-scoring All Items ===")
	
	cfg := config.Load()
	
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	
	ctx := context.Background()
	
	itemRepo := db.NewItemRepository(database)
	signalRepo := db.NewSignalRepository(database)
	scoreRepo := db.NewScoreRepository(database)
	summaryRepo := db.NewSummaryRepository(database)
	
	// Get all items from last 30 days
	items, err := itemRepo.GetItemsNeedingScoring(ctx, 30, 10000)
	if err != nil {
		log.Fatalf("Failed to get items: %v", err)
	}
	
	fmt.Printf("Found %d items to re-score\n", len(items))
	
	updated := 0
	for i, item := range items {
		if i%100 == 0 {
			fmt.Printf("Progress: %d/%d (%.1f%%)\n", i, len(items), float64(i)/float64(len(items))*100)
		}
		
		// Get latest signal
		signal, _ := signalRepo.GetLatestByItemID(ctx, item.ID)
		
		// Get summary
		summary, _ := summaryRepo.GetByItemID(ctx, item.ID)
		
		// Compute new score
		score := scoring.ComputeScore(&item, signal, summary, scoring.DefaultWeights)
		
		// Update score in database
		if err := scoreRepo.Upsert(ctx, score); err != nil {
			log.Printf("Error updating score for item %s: %v", item.ID, err)
			continue
		}
		
		updated++
	}
	
	fmt.Printf("\n=== Re-scoring Complete ===\n")
	fmt.Printf("Total items processed: %d\n", len(items))
	fmt.Printf("Successfully updated: %d\n", updated)
	fmt.Printf("Failed: %d\n", len(items)-updated)
}
