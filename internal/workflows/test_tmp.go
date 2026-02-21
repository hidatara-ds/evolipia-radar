package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/hidatara-ds/evolipia-radar/internal/activities"
	"github.com/hidatara-ds/evolipia-radar/internal/mlpipeline"
)

// MLPipelineWorkflow orchestrates the end-to-end ML lifecycle:
// Ingestion → Feature Engineering → Training → Evaluation → (conditional) Deployment.
func MLPipelineWorkflow(ctx workflow.Context, config mlpipeline.PipelineConfig) error {
	// Default training config if not provided
	if config.TrainingConfig.MaxTrainingDuration == 0 {
		config.TrainingConfig.MaxTrainingDuration = 2 * time.Hour
	}
	if config.TrainingConfig.MinAcceptableScore == 0 {
		config.TrainingConfig.MinAcceptableScore = 0.85
	}
	if config.TrainingConfig.TargetMetric == "" {
		config.TrainingConfig.TargetMetric = "accuracy"
	}

	// Common retry policy for all activities.
	retryPolicy := &temporal.RetryPolicy{
		MaximumAttempts: 3,
	}

	// Ingestion
	ingestionOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ingestionOpts)

	var ingestionResult mlpipeline.IngestionResult
	if err := workflow.ExecuteActivity(ctx, activities.IngestionActivity, config.Sources).Get(ctx, &ingestionResult); err != nil {
		return err
	}

	// Feature engineering
	featureOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, featureOpts)

	var featureResult mlpipeline.FeatureEngResult
	if err := workflow.ExecuteActivity(ctx, activities.FeatureEngineeringActivity, ingestionResult).Get(ctx, &featureResult); err != nil {
		return err
	}

	// Training
	trainingOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    config.TrainingConfig.MaxTrainingDuration,
		ScheduleToCloseTimeout: config.TrainingConfig.MaxTrainingDuration + 30*time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, trainingOpts)

	var trainingResult mlpipeline.TrainingResult
	if err := workflow.ExecuteActivity(ctx, activities.TrainingActivity, featureResult, config.TrainingConfig).Get(ctx, &trainingResult); err != nil {
		return err
	}

	// Evaluation
	evalOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, evalOpts)

	var evalResult mlpipeline.EvaluationResult
	if err := workflow.ExecuteActivity(ctx, activities.EvaluationActivity, trainingResult, config.TrainingConfig).Get(ctx, &evalResult); err != nil {
		return err
	}

	// Automatic gate based on accuracy threshold.
	if !evalResult.ShouldDeploy || evalResult.Accuracy < config.TrainingConfig.MinAcceptableScore {
		return nil
	}

	// Human-in-the-loop approval signal gate.
	signalChan := workflow.GetSignalChannel(ctx, "human_approval")
	selector := workflow.NewSelector(ctx)

	var approval mlpipeline.HumanApprovalSignal
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &approval)
	})

	// Optionally add a timeout to wait for approval.
	approvalTimeout := workflow.NewTimer(ctx, 24*time.Hour)
	selector.AddFuture(approvalTimeout, func(f workflow.Future) {
		// If timeout fires first, leave approval as default (false).
	})

	selector.Select(ctx)
	if !approval.Approved {
		return nil
	}

	// Deployment
	deployOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, deployOpts)

	var _ mlpipeline.DeploymentResult
	if err := workflow.ExecuteActivity(ctx, activities.DeploymentActivity, trainingResult, evalResult, approval).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
