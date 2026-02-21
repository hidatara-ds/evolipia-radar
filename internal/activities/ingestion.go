package activities

import (
	"context"

	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// IngestionActivity is a Temporal activity responsible for ingesting raw content
// from configured sources and materializing it into the offline store.
// For now this is a thin stub that will be wired to existing connectors/worker logic.
func IngestionActivity(ctx context.Context, sources []string) (mlpipeline.IngestionResult, error) {
	// TODO: integrate with existing connectors and ingestion flow.
	return mlpipeline.IngestionResult{
		RunID:          "ingestion-stub",
		ItemsIngested:  0,
		FeatureSetName: "",
	}, nil
}
