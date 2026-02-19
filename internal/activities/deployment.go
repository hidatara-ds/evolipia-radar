package activities

import (
	"context"

	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// DeploymentActivity promotes a candidate model to the desired stage
// (e.g., Staging/Production) in the model registry and updates serving config.
func DeploymentActivity(ctx context.Context, trainingResult mlpipeline.TrainingResult, evalResult mlpipeline.EvaluationResult, approval mlpipeline.HumanApprovalSignal) (mlpipeline.DeploymentResult, error) {
	// TODO: integrate with MLflow model registry and update serving configuration.
	return mlpipeline.DeploymentResult{
		RunID:        "deployment-stub",
		ModelVersion: trainingResult.ModelVersion,
		Stage:        "Production",
	}, nil
}

