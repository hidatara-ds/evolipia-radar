package summarizer

// TopicKeywordsConfig holds topic keyword configuration for tag extraction
type TopicKeywordsConfig struct {
	Keywords map[string][]string
}

// DefaultTopicKeywordsConfig returns the default topic keywords configuration
func DefaultTopicKeywordsConfig() TopicKeywordsConfig {
	return TopicKeywordsConfig{
		Keywords: map[string][]string{
			"llm":             {"llm", "transformer", "gpt", "gemini", "llama", "mistral", "language model"},
			"nlp":             {"nlp", "natural language", "text processing", "sentiment"},
			"computer_vision": {"computer vision", "cv", "yolo", "segmentation", "detection", "opencv"},
			"mlops":           {"mlops", "deployment", "monitoring", "drift", "kubernetes", "kubeflow"},
			"data":            {"data", "dataset", "data pipeline", "etl"},
			"cloud":           {"aws", "gcp", "azure", "cloud", "s3", "lambda"},
			"security":        {"security", "privacy", "encryption", "adversarial"},
			"general_ai":      {"ai", "artificial intelligence", "machine learning", "deep learning"},
		},
	}
}
