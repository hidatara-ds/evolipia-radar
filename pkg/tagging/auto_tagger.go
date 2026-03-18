package tagging

import (
	"regexp"
	"strings"
)

// AutoTagger automatically assigns tags based on content analysis
type AutoTagger struct {
	rules []TagRule
}

// TagRule defines a pattern and corresponding tag
type TagRule struct {
	Tag      string
	Patterns []string // Regex patterns
	Keywords []string // Simple keyword matching (case-insensitive)
}

// NewAutoTagger creates a new auto-tagger with predefined rules
func NewAutoTagger() *AutoTagger {
	return &AutoTagger{
		rules: []TagRule{
			// LLM & Language Models
			{
				Tag: "llm",
				Keywords: []string{
					"gpt", "llama", "claude", "gemini", "mistral", "falcon",
					"language model", "large language model", "transformer",
					"chatgpt", "bard", "palm", "llm", "generative ai",
				},
			},
			// Computer Vision
			{
				Tag: "vision",
				Keywords: []string{
					"diffusion", "stable diffusion", "dall-e", "dalle", "midjourney",
					"computer vision", "image generation", "text-to-image",
					"vision transformer", "vit", "clip", "image recognition",
					"object detection", "segmentation",
				},
			},
			// AI Safety & Alignment
			{
				Tag: "safety",
				Keywords: []string{
					"alignment", "rlhf", "safety", "constitutional",
					"jailbreak", "adversarial", "robustness", "interpretability",
					"explainability", "ai safety", "responsible ai",
				},
			},
			// Reinforcement Learning
			{
				Tag: "rl",
				Keywords: []string{
					"reinforcement learning", "rl", "reward", "policy gradient",
					"q-learning", "dqn", "ppo", "actor-critic", "markov",
				},
			},
			// Robotics & Embodied AI
			{
				Tag: "robotics",
				Keywords: []string{
					"robot", "manipulation", "locomotion", "embodied",
					"robotics", "autonomous", "drone", "self-driving",
				},
			},
			// Data & Datasets
			{
				Tag: "data",
				Keywords: []string{
					"dataset", "benchmark", "evaluation", "leaderboard",
					"data", "corpus", "annotation", "labeling",
				},
			},
			// Security & Privacy
			{
				Tag: "security",
				Keywords: []string{
					"security", "privacy", "encryption", "vulnerability",
					"attack", "defense", "backdoor", "poisoning",
				},
			},
			// IDE & Developer Tools (NEW!)
			{
				Tag: "ide",
				Keywords: []string{
					"kiro", "cursor", "windsurf", "codeium", "copilot",
					"github copilot", "tabnine", "replit", "ghostwriter",
					"warp", "fig", "zed", "fleet", "nova",
					"ai ide", "code editor", "code assistant",
					"intellij", "vscode", "visual studio code",
					"jetbrains", "sublime", "atom",
				},
			},
			// Free Credits & Student Programs (NEW!)
			{
				Tag: "free-credits",
				Keywords: []string{
					"free credit", "free token", "free api", "student program",
					"education program", "academic program", "free tier",
					"free access", "student discount", "github student",
					"anthropic student", "openai credit", "azure credit",
					"gcp credit", "aws educate", "free trial",
				},
			},
			// Research & Papers
			{
				Tag: "research",
				Keywords: []string{
					"arxiv", "paper", "research", "study", "conference",
					"neurips", "icml", "iclr", "cvpr", "acl", "emnlp",
				},
			},
			// Tools & Frameworks
			{
				Tag: "tools",
				Keywords: []string{
					"framework", "library", "tool", "sdk", "api",
					"pytorch", "tensorflow", "jax", "huggingface",
					"langchain", "llamaindex",
				},
			},
		},
	}
}

// AssignTags analyzes title and content to assign relevant tags
func (at *AutoTagger) AssignTags(title, content string) []string {
	// Combine title and content for analysis (title weighted more)
	text := strings.ToLower(title + " " + title + " " + content)
	
	tagSet := make(map[string]bool)
	
	for _, rule := range at.rules {
		// Check keywords
		for _, keyword := range rule.Keywords {
			if strings.Contains(text, strings.ToLower(keyword)) {
				tagSet[rule.Tag] = true
				break
			}
		}
		
		// Check regex patterns
		for _, pattern := range rule.Patterns {
			matched, err := regexp.MatchString(pattern, text)
			if err == nil && matched {
				tagSet[rule.Tag] = true
				break
			}
		}
	}
	
	// Convert set to slice
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	
	// If no tags matched, assign "general_ai"
	if len(tags) == 0 {
		tags = append(tags, "general_ai")
	}
	
	return tags
}

// MergeTags combines existing tags with auto-generated tags (deduplicates)
func MergeTags(existingTags, newTags []string) []string {
	tagSet := make(map[string]bool)
	
	// Add existing tags
	for _, tag := range existingTags {
		if tag != "" {
			tagSet[tag] = true
		}
	}
	
	// Add new tags
	for _, tag := range newTags {
		if tag != "" {
			tagSet[tag] = true
		}
	}
	
	// Convert to slice
	result := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		result = append(result, tag)
	}
	
	return result
}
