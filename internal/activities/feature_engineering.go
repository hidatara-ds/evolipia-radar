package activities

import (
	"context"

	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// FeatureEngineeringActivity is a Temporal activity that builds feature tables
// in both the offline and online feature stores.
func FeatureEngineeringActivity(ctx context.Context, ingestionResult mlpipeline.IngestionResult) (mlpipeline.FeatureEngResult, error) {
	// TODO: implement feature computation and write to feature store backends.
	return mlpipeline.FeatureEngResult{}, nil
}
