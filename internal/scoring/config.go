package scoring

// CredibilityConfig holds credibility scoring configuration
type CredibilityConfig struct {
	Whitelist map[string]bool
	Blacklist map[string]bool
}

// DefaultCredibilityConfig returns the default credibility configuration
func DefaultCredibilityConfig() CredibilityConfig {
	return CredibilityConfig{
		Whitelist: map[string]bool{
			"openai.com":           true,
			"ai.googleblog.com":    true,
			"deepmind.google":      true,
			"arxiv.org":            true,
			"acm.org":              true,
			"ieee.org":             true,
			"github.com":           true,
			"docs.github.com":      true,
		},
		Blacklist: map[string]bool{
			"medium.com": true,
		},
	}
}

// RelevanceKeywords holds keyword configuration for relevance scoring
type RelevanceKeywords struct {
	LLM    []string
	MLOps  []string
	CV     []string
	Weights map[string]float64
}

// DefaultRelevanceKeywords returns the default relevance keywords configuration
func DefaultRelevanceKeywords() RelevanceKeywords {
	return RelevanceKeywords{
		LLM: []string{"llm", "transformer", "rag", "prompt", "inference", "fine-tune", "gpt", "gemini", "llama", "mistral"},
		MLOps: []string{"mlops", "deployment", "monitoring", "drift", "kubernetes", "kubeflow", "airflow", "feature store"},
		CV: []string{"computer vision", "yolo", "segmentation", "detection", "opencv"},
		Weights: map[string]float64{
			"llm":   0.3,
			"mlops": 0.25,
			"cv":    0.2,
			"tag":   0.2,
		},
	}
}

