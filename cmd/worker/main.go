package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evolipia/radar/internal/config"
	"github.com/evolipia/radar/internal/db"
	"github.com/evolipia/radar/internal/services"
	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	w := services.NewWorker(database, cfg)

	c := cron.New()
	_, err = c.AddFunc(cfg.WorkerCron, func() {
		log.Println("Starting scheduled ingestion...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := w.RunIngestion(ctx); err != nil {
			log.Printf("Ingestion error: %v", err)
		} else {
			log.Println("Ingestion completed successfully")
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	log.Printf("Worker started with cron schedule: %s", cfg.WorkerCron)

	// Run once immediately
	log.Println("Running initial ingestion...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := w.RunIngestion(ctx); err != nil {
		log.Printf("Initial ingestion error: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	c.Stop()
	log.Println("Worker exited")
}
