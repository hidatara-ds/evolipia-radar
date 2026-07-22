package scoring

type CredibilityConfig struct {
	Tier1 map[string]float64
	Tier2 map[string]float64
	Tier3 map[string]float64
}

// DefaultCredibilityConfig returns the default credibility configuration
func DefaultCredibilityConfig() CredibilityConfig {
	return CredibilityConfig{
		Tier1: map[string]float64{
			"arxiv.org":             1.2,
			"deepmind.google":       1.2,
			"openai.com":            1.2,
			"ai.googleblog.com":     1.2,
			"research.facebook.com": 1.2,
		},
		Tier2: map[string]float64{
			"anthropic.com": 1.0,
			"acm.org":       1.0,
			"ieee.org":      1.0,
		},
		Tier3: map[string]float64{
			"github.com":      0.7,
			"techcrunch.com":  0.7,
			"venturebeat.com": 0.7,
		},
	}
}

// RelevanceKeywords holds keyword configuration for relevance scoring
type RelevanceKeywords struct {
	LLM     []string
	MLOps   []string
	CV      []string
	Weights map[string]float64
}

// DefaultRelevanceKeywords returns the default relevance keywords configuration
func DefaultRelevanceKeywords() RelevanceKeywords {
	return RelevanceKeywords{
		LLM:   []string{"llm", "transformer", "rag", "prompt", "inference", "fine-tune", "gpt", "gemini", "llama", "mistral"},
		MLOps: []string{"mlops", "deployment", "monitoring", "drift", "kubernetes", "kubeflow", "airflow", "feature store"},
		CV:    []string{"computer vision", "yolo", "segmentation", "detection", "opencv"},
		Weights: map[string]float64{
			"llm":   0.3,
			"mlops": 0.25,
			"cv":    0.2,
			"tag":   0.2,
		},
	}
}
