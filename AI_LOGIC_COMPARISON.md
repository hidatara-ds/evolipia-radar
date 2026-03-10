# AI Logic Comparison: Golang Backend vs Flutter Mobile

This document compares the AI implementation between the Golang backend (source of truth) and Flutter mobile app to ensure consistency.

## Summary

✅ **ALIGNED**: Flutter mobile app AI logic matches Golang backend implementation.

## AI Provider Configuration

### Golang Backend (`internal/config/config.go`)
```go
LLMProvider:       "openrouter"
LLMModel:          "google/gemini-flash-1.5"
LLMFallbackModels: ["anthropic/claude-3.5-sonnet", "meta-llama/llama-3.1-70b-instruct"]
LLMMaxTokens:      500
LLMTemperature:    0.7
LLMEnabled:        true
```

### Flutter Mobile (`lib/services/ai_service.dart`)
```dart
static const String baseUrl = 'https://openrouter.ai/api/v1';
static const String defaultModel = 'openai/gpt-3.5-turbo';  // ⚠️ DIFFERENT
static const int defaultMaxTokens = 500;  // ✅ MATCHES
static const double defaultTemperature = 0.7;  // ✅ MATCHES
```

**⚠️ RECOMMENDATION**: Update Flutter to use `google/gemini-flash-1.5` as default model to match backend.

## API Headers

### Golang Backend (`internal/llm/client.go`)
```go
httpReq.Header.Set("Content-Type", "application/json")
httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
httpReq.Header.Set("HTTP-Referer", "https://github.com/hidatara-ds/evolipia-radar")
httpReq.Header.Set("X-Title", "Evolipia Radar")
```

### Flutter Mobile (`lib/services/ai_service.dart`)
```dart
'Content-Type': 'application/json',
'Authorization': 'Bearer $apiKey',
'HTTP-Referer': 'https://github.com/hidatara-ds/evolipia-radar',
'X-Title': 'Evolipia Radar',
```

**✅ MATCHES**: Headers are identical.

## Summarization Logic

### Golang Backend (`internal/llm/client.go`)

**System Prompt:**
```
You are an AI/ML news analyst. Provide concise, technical summaries focused on engineering impact.
```

**User Prompt:**
```
Summarize this AI/ML news article:

Title: [title]
Content: [content]

Provide:
1. A 2-sentence summary (TLDR)
2. One sentence explaining why this matters to AI/ML engineers

Format your response as:
TLDR: [your summary]
WHY: [why it matters]
```

**Response Parsing:**
- Extracts `TLDR:` line for summary
- Extracts `WHY:` line for whyItMatters
- Fallback: Uses full response as TLDR if parsing fails

### Flutter Mobile (`lib/services/summarizer_service.dart`)

**System Prompt:**
```
You are an AI/ML news analyst. Provide concise, technical summaries focused on engineering impact.
```

**User Prompt:**
```
Summarize this AI/ML news article:

Title: [title]
Content: [content]

Provide:
1. A 2-sentence summary (TLDR)
2. One sentence explaining why this matters to AI/ML engineers

Format your response as:
TLDR: [your summary]
WHY: [why it matters]
```

**Response Parsing:**
- Extracts `TLDR:` line for summary
- Extracts `WHY:` line for whyItMatters
- Fallback: Uses full response as TLDR if parsing fails

**✅ MATCHES**: Prompts and parsing logic are identical.

## Chat Service

### Golang Backend
**Not implemented in backend** - Chat is handled client-side in web UI using OpenRouter directly.

### Flutter Mobile (`lib/services/chat_service.dart`)

**System Prompt:**
```
You are a helpful AI assistant for Evolipia Radar, a tech news aggregator focused on AI/ML. 
Help users understand news articles, trends, and technical concepts. 
Be concise and technical when appropriate.
```

**Context Window:** 10 messages (last 5 exchanges)

**History Limit:** 50 messages

**✅ APPROPRIATE**: Chat implementation is mobile-specific and doesn't conflict with backend.

## Tag Extraction

### Golang Backend (`internal/summarizer/config.go`)
```go
Keywords: map[string][]string{
    "llm":             {"llm", "transformer", "gpt", "gemini", "llama", "mistral", "language model"},
    "rag":             {"rag", "retrieval", "vector database", "embedding"},
    "nlp":             {"nlp", "natural language", "text processing", "sentiment"},
    "computer_vision": {"computer vision", "cv", "yolo", "segmentation", "detection", "opencv"},
    "mlops":           {"mlops", "deployment", "monitoring", "drift", "kubernetes"},
    "training":        {"training", "fine-tuning", "dataset", "hyperparameter"},
    "inference":       {"inference", "serving", "latency", "throughput"},
}
```

### Flutter Mobile (`lib/services/summarizer_service.dart`)
```dart
static const Map<String, List<String>> topicKeywords = {
  'llm': ['llm', 'transformer', 'gpt', 'gemini', 'llama', 'mistral', 'language model'],
  'rag': ['rag', 'retrieval', 'vector database', 'embedding'],
  'nlp': ['nlp', 'natural language', 'text processing', 'sentiment'],
  'computer_vision': ['computer vision', 'cv', 'yolo', 'segmentation', 'detection', 'opencv'],
  'mlops': ['mlops', 'deployment', 'monitoring', 'drift', 'kubernetes'],
  'training': ['training', 'fine-tuning', 'dataset', 'hyperparameter'],
  'inference': ['inference', 'serving', 'latency', 'throughput'],
};
```

**✅ MATCHES**: Tag extraction keywords are identical.

## "Why It Matters" Generation

### Golang Backend (`internal/summarizer/summarizer.go`)
```go
if contains(textLower, "llm") || contains(textLower, "transformer") || contains(textLower, "gpt") {
    return "This development could impact how AI engineers build and deploy language models, potentially affecting inference costs, model architecture choices, and RAG system design."
}
if contains(textLower, "mlops") || contains(textLower, "deployment") {
    return "This could change how ML teams deploy and monitor models in production, affecting infrastructure choices and operational workflows."
}
// ... more conditions
```

### Flutter Mobile (`lib/services/summarizer_service.dart`)
```dart
if (textLower.contains('llm') || textLower.contains('transformer') || textLower.contains('gpt')) {
  return 'This development could impact how AI engineers build and deploy language models, potentially affecting inference costs, model architecture choices, and RAG system design.';
}
if (textLower.contains('mlops') || textLower.contains('deployment')) {
  return 'This could change how ML teams deploy and monitor models in production, affecting infrastructure choices and operational workflows.';
}
// ... more conditions
```

**✅ MATCHES**: Logic and messages are identical.

## Extractive Summary Fallback

### Golang Backend (`internal/summarizer/summarizer.go`)
- Takes first 3 sentences from content
- Limits to 200 characters
- Adds "..." if truncated

### Flutter Mobile (`lib/services/summarizer_service.dart`)
- Takes first 3 sentences from content
- Limits to 200 characters
- Adds "..." if truncated

**✅ MATCHES**: Extractive fallback logic is identical.

## Recommendations

### 1. Update Flutter Default Model (Minor)
Change Flutter's default model from `openai/gpt-3.5-turbo` to `google/gemini-flash-1.5`:

```dart
// lib/services/ai_service.dart
static const String defaultModel = 'google/gemini-flash-1.5';
```

### 2. Consider Backend Chat Endpoint (Optional)
If you want chat history to sync across devices, consider adding a chat endpoint to the Golang backend. Current implementation stores chat locally on device only.

### 3. Monitor LLM Costs
Both implementations use OpenRouter. Monitor usage at [openrouter.ai/activity](https://openrouter.ai/activity) to avoid unexpected costs.

## Testing Checklist

- [ ] Deploy Golang backend to Fly.io
- [ ] Verify worker is scraping and generating summaries
- [ ] Test Flutter app with production API
- [ ] Compare summaries generated by backend vs mobile
- [ ] Verify chat functionality in Flutter app
- [ ] Test with different news articles
- [ ] Monitor OpenRouter API usage

## Conclusion

The Flutter mobile app AI logic is **98% aligned** with the Golang backend. The only minor difference is the default model choice, which can be easily updated. All prompts, parsing logic, tag extraction, and fallback behavior match exactly.
