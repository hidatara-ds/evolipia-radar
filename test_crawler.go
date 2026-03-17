package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hidatara-ds/evolipia-radar/internal/crawler"
)

func main() {
	metrics := &crawler.Metrics{}
	// Nil ClusterService is fine since DryRun=true bypasses it
	orc := crawler.NewOrchestrator(nil, nil, metrics, true)

	ctx := context.Background()

	log.Println("--- TRIGGER 1 ---")
	stats1 := orc.RunCycle(ctx)
	b1, _ := json.MarshalIndent(stats1, "", "  ")
	fmt.Println(string(b1))

	log.Println("--- TRIGGER 2 (Should Dedupe) ---")
	stats2 := orc.RunCycle(ctx)
	b2, _ := json.MarshalIndent(stats2, "", "  ")
	fmt.Println(string(b2))
	
	orc.UpdateClusterMetrics(ctx)

	log.Println("--- FINAL METRICS ---")
	bm, _ := json.MarshalIndent(metrics, "", "  ")
	fmt.Println(string(bm))
}
