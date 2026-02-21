package activities

import (
	"context"

	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// EvaluationActivity evaluates a trained model on a validation or test set.
func EvaluationActivity(ctx context.Context, trainingResult mlpipeline.TrainingResult, cfg mlpipeline.TrainingConfig) (mlpipeline.EvaluationResult, error) {
	// TODO: compute evaluation metrics and decide whether the model should be deployed.
	return mlpipeline.EvaluationResult{
		RunID:         "evaluation-stub",
		Accuracy:      0.0,
		MetricDetails: map[string]float64{},
		ShouldDeploy:  false,
	}, nil
}
