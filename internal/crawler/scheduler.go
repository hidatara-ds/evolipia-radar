package crawler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/robfig/cron/v3"
)

// CrawlTaskFunc is the signature for triggering a crawl run with progress reporting.
type CrawlTaskFunc func(ctx context.Context, onProgress func(models.CrawlProgressEvent)) (int, error)

// Scheduler manages automated background crawl jobs using cron.
type Scheduler struct {
	cron             *cron.Cron
	crawlInterval    string
	crawlTask        CrawlTaskFunc
	isRunning        atomic.Bool
	wg               sync.WaitGroup
	lastRunTime      time.Time
	lastRunStatus    string
	lastItemsCount   int
	lastRunError     string
	mu               sync.RWMutex
	progressReporter func(models.CrawlProgressEvent)
}

// NewScheduler creates a Scheduler instance with interval schedule.
func NewScheduler(interval string, task CrawlTaskFunc, progressReporter func(models.CrawlProgressEvent)) (*Scheduler, error) {
	if interval == "" {
		interval = "@every 6h"
	}

	c := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))

	s := &Scheduler{
		cron:             c,
		crawlInterval:    interval,
		crawlTask:        task,
		lastRunStatus:    "never_run",
		progressReporter: progressReporter,
	}

	_, err := c.AddFunc(interval, func() {
		s.RunCrawl(context.Background(), "auto")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to schedule crawl job with spec '%s': %w", interval, err)
	}

	return s, nil
}

// Start launches the cron scheduler loop in the background.
func (s *Scheduler) Start() {
	slog.Info("Starting Auto-Scheduler...", "interval", s.crawlInterval)
	s.cron.Start()
}

// Stop gracefully shuts down the cron scheduler and waits for active crawls to complete.
func (s *Scheduler) Stop() {
	slog.Info("Shutting down Auto-Scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.wg.Wait()
	slog.Info("Auto-Scheduler stopped gracefully")
}

// RunCrawl executes one crawl cycle if not already running.
func (s *Scheduler) RunCrawl(ctx context.Context, triggerType string) (int, error) {
	if !s.isRunning.CompareAndSwap(false, true) {
		slog.Warn("Crawl cycle skipped: previous crawl is still in progress")
		return 0, fmt.Errorf("crawl already in progress")
	}

	s.wg.Add(1)
	defer func() {
		s.isRunning.Store(false)
		s.wg.Done()
	}()

	startTime := time.Now()
	slog.Info("Starting crawl cycle", "trigger", triggerType, "time", startTime.Format(time.RFC3339))

	reporter := func(ev models.CrawlProgressEvent) {
		if s.progressReporter != nil {
			s.progressReporter(ev)
		}
	}

	reporter(models.CrawlProgressEvent{
		Step:       1,
		Message:    "Starting crawl cycle...",
		Progress:   10,
		Timestamp:  time.Now(),
		IsComplete: false,
	})

	itemsProcessed, err := s.crawlTask(ctx, reporter)

	s.mu.Lock()
	s.lastRunTime = time.Now()
	s.lastItemsCount = itemsProcessed

	if err != nil {
		s.lastRunStatus = "failed"
		s.lastRunError = err.Error()
		slog.Error("Crawl cycle finished with error", "duration", time.Since(startTime), "err", err)
		reporter(models.CrawlProgressEvent{
			Step:       6,
			Message:    fmt.Sprintf("Crawl failed: %v", err),
			Progress:   100,
			HasError:   true,
			Error:      err.Error(),
			IsComplete: true,
			Timestamp:  time.Now(),
		})
	} else {
		s.lastRunStatus = "success"
		s.lastRunError = ""
		slog.Info("Finished crawl cycle successfully", "total_items", itemsProcessed, "duration", time.Since(startTime))
		reporter(models.CrawlProgressEvent{
			Step:           6,
			Message:        fmt.Sprintf("Done! %d items processed", itemsProcessed),
			Progress:       100,
			ProcessedItems: itemsProcessed,
			IsComplete:     true,
			Timestamp:      time.Now(),
		})
	}
	s.mu.Unlock()

	return itemsProcessed, err
}

// GetStatus returns the last crawl execution metrics.
func (s *Scheduler) GetStatus() (time.Time, string, int, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastRunTime, s.lastRunStatus, s.lastItemsCount, s.lastRunError, s.isRunning.Load()
}
