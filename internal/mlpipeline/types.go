package mlpipeline

import "time"

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

