package activities

import (
	"context"

	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// TrainingActivity trains a ranking model using the prepared offline feature set.
func TrainingActivity(ctx context.Context, featureResult mlpipeline.FeatureEngResult, cfg mlpipeline.TrainingConfig) (mlpipeline.TrainingResult, error) {
	// TODO: call out to Python ML service / MLflow-backed training pipeline.
	return mlpipeline.TrainingResult{
		RunID:        "training-stub",
		ModelURI:     "",
		ModelVersion: "",
		Metrics:      map[string]float64{},
	}, nil
}

