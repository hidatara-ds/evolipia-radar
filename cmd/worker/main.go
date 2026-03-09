package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hidatara-ds/evolipia-radar/internal/config"
	"github.com/hidatara-ds/evolipia-radar/internal/db"
	"github.com/hidatara-ds/evolipia-radar/internal/services"
)

func main() {
	startTime := time.Now()
	log.Printf("========================================")
	log.Printf("Worker started at: %s", startTime.Format(time.RFC3339))
	log.Printf("========================================")

	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create scrape log entry
	scrapeLogID := uuid.New()
	if err := createScrapeLog(database, scrapeLogID); err != nil {
		log.Printf("Warning: Failed to create scrape log: %v", err)
	}

	// Run one-shot ingestion
	w := services.NewWorker(database, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	log.Println("Starting ingestion...")
	itemsProcessed := 0
	itemsNew := 0
	status := "success"
	var errorMsg string

	if err := w.RunIngestion(ctx); err != nil {
		log.Printf("Ingestion error: %v", err)
		status = "failed"
		errorMsg = err.Error()
		
		// Update scrape log with failure
		if updateErr := updateScrapeLog(database, scrapeLogID, itemsProcessed, itemsNew, status, errorMsg); updateErr != nil {
			log.Printf("Warning: Failed to update scrape log: %v", updateErr)
		}
		
		log.Printf("========================================")
		log.Printf("Worker finished at: %s", time.Now().Format(time.RFC3339))
		log.Printf("Status: FAILED")
		log.Printf("========================================")
		os.Exit(1)
	}

	// Get stats (you may need to modify services.Worker to return these)
	itemsProcessed = getItemsProcessed(database)
	itemsNew = getItemsNew(database, startTime)

	log.Printf("Ingestion completed successfully")
	log.Printf("Items processed: %d", itemsProcessed)
	log.Printf("New items: %d", itemsNew)

	// Update scrape log with success
	if err := updateScrapeLog(database, scrapeLogID, itemsProcessed, itemsNew, status, errorMsg); err != nil {
		log.Printf("Warning: Failed to update scrape log: %v", err)
	}

	log.Printf("========================================")
	log.Printf("Worker finished at: %s", time.Now().Format(time.RFC3339))
	log.Printf("Duration: %s", time.Since(startTime))
	log.Printf("Status: SUCCESS")
	log.Printf("========================================")
}

func createScrapeLog(database *db.DB, id uuid.UUID) error {
	query := `
		INSERT INTO scrape_logs (id, started_at, status, trigger_source)
		VALUES ($1, $2, $3, $4)
	`
	_, err := database.Pool.Exec(context.Background(), query, id, time.Now(), "running", getenv("TRIGGER_SOURCE", "github_actions"))
	return err
}

func updateScrapeLog(database *db.DB, id uuid.UUID, itemsProcessed, itemsNew int, status, errorMsg string) error {
	query := `
		UPDATE scrape_logs
		SET completed_at = $1, items_processed = $2, items_new = $3, status = $4, error_message = $5
		WHERE id = $6
	`
	_, err := database.Pool.Exec(context.Background(), query, time.Now(), itemsProcessed, itemsNew, status, errorMsg, id)
	return err
}

func getItemsProcessed(database *db.DB) int {
	var count int
	query := `SELECT COUNT(*) FROM items`
	if err := database.Pool.QueryRow(context.Background(), query).Scan(&count); err != nil {
		log.Printf("Warning: Failed to get items count: %v", err)
		return 0
	}
	return count
}

func getItemsNew(database *db.DB, since time.Time) int {
	var count int
	query := `SELECT COUNT(*) FROM items WHERE created_at >= $1`
	if err := database.Pool.QueryRow(context.Background(), query, since).Scan(&count); err != nil {
		log.Printf("Warning: Failed to get new items count: %v", err)
		return 0
	}
	return count
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
