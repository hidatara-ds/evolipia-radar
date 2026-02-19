package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// PipelineConfig holds high-level configuration for a single ML pipeline run.
// It is intentionally generic so it can be evolved without changing the workflow shape.
type PipelineConfig struct {
	Sources        []string
	ModelType      string
	TrainingConfig TrainingConfig
}

// TrainingConfig captures knobs relevant for the training activity.
type TrainingConfig struct {
	MaxTrainingDuration time.Duration
	TargetMetric        string
	MinAcceptableScore  float64
}

// IngestionResult represents the outcome of the ingestion activity.
type IngestionResult struct {
	RunID          string
	ItemsIngested  int
	FeatureSetName string
}

// FeatureEngResult represents the outcome of the feature engineering activity.
type FeatureEngResult struct {
	FeatureTable       string
	FeatureVersion     string
	OfflineStoreRef    string
	OnlineStoreRef     string
	PointInTimeCutover time.Time
}

// TrainingResult represents the outcome of the training activity.
type TrainingResult struct {
	RunID        string
	ModelURI     string
	ModelVersion string
	Metrics      map[string]float64
}

// EvaluationResult represents the outcome of the evaluation activity.
type EvaluationResult struct {
	RunID         string
	Accuracy      float64
	AUC           *float64
	Loss          *float64
	MetricDetails map[string]float64
	ShouldDeploy  bool
}

// DeploymentResult represents the outcome of the deployment activity.
type DeploymentResult struct {
	RunID        string
	ModelVersion string
	Stage        string
}

// HumanApprovalSignal is used to gate deployment on an external human approval.
type HumanApprovalSignal struct {
	Approved bool
	Reason   string
}

// MLPipelineWorkflow orchestrates the end-to-end ML lifecycle:
// Ingestion → Feature Engineering → Training → Evaluation → (conditional) Deployment.
func MLPipelineWorkflow(ctx workflow.Context, config PipelineConfig) error {
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
	retryPolicy := &workflow.RetryPolicy{
		MaximumAttempts: 3,
	}

	// Ingestion
	ingestionOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ingestionOpts)

	var ingestionResult IngestionResult
	if err := workflow.ExecuteActivity(ctx, IngestionActivity, config.Sources).Get(ctx, &ingestionResult); err != nil {
		return err
	}

	// Feature engineering
	featureOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, featureOpts)

	var featureResult FeatureEngResult
	if err := workflow.ExecuteActivity(ctx, FeatureEngineeringActivity, ingestionResult).Get(ctx, &featureResult); err != nil {
		return err
	}

	// Training
	trainingOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    config.TrainingConfig.MaxTrainingDuration,
		ScheduleToCloseTimeout: config.TrainingConfig.MaxTrainingDuration + 30*time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, trainingOpts)

	var trainingResult TrainingResult
	if err := workflow.ExecuteActivity(ctx, TrainingActivity, featureResult, config.TrainingConfig).Get(ctx, &trainingResult); err != nil {
		return err
	}

	// Evaluation
	evalOpts := workflow.ActivityOptions{
		StartToCloseTimeout:    30 * time.Minute,
		ScheduleToCloseTimeout: 60 * time.Minute,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, evalOpts)

	var evalResult EvaluationResult
	if err := workflow.ExecuteActivity(ctx, EvaluationActivity, trainingResult, config.TrainingConfig).Get(ctx, &evalResult); err != nil {
		return err
	}

	// Automatic gate based on accuracy threshold.
	if !evalResult.ShouldDeploy || evalResult.Accuracy < config.TrainingConfig.MinAcceptableScore {
		return nil
	}

	// Human-in-the-loop approval signal gate.
	signalChan := workflow.GetSignalChannel(ctx, "human_approval")
	selector := workflow.NewSelector(ctx)

	var approval HumanApprovalSignal
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

	var _ DeploymentResult
	if err := workflow.ExecuteActivity(ctx, DeploymentActivity, trainingResult, evalResult, approval).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

