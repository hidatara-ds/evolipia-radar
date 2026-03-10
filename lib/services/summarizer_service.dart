import '../models/news_item.dart';
import '../models/summary.dart';
import 'ai_service.dart';

class SummarizerService {
  final AIService? aiService;

  SummarizerService({this.aiService});

  /// Topic keywords matching Golang configuration
  static const Map<String, List<String>> topicKeywords = {
    'llm': [
      'llm',
      'transformer',
      'gpt',
      'gemini',
      'llama',
      'mistral',
      'language model'
    ],
    'nlp': ['nlp', 'natural language', 'text processing', 'sentiment'],
    'computer_vision': [
      'computer vision',
      'cv',
      'yolo',
      'segmentation',
      'detection',
      'opencv'
    ],
    'mlops': [
      'mlops',
      'deployment',
      'monitoring',
      'drift',
      'kubernetes',
      'kubeflow'
    ],
    'data': ['data', 'dataset', 'data pipeline', 'etl'],
    'cloud': ['aws', 'gcp', 'azure', 'cloud', 's3', 'lambda'],
    'security': ['security', 'privacy', 'encryption', 'adversarial'],
    'general_ai': [
      'ai',
      'artificial intelligence',
      'machine learning',
      'deep learning'
    ],
  };

  /// Generate summary using LLM or fallback to extractive
  Future<Summary> generateSummary(NewsItem item) async {
    if (aiService != null) {
      try {
        return await generateLLMSummary(item);
      } catch (e) {
        // Fallback to extractive on error
        return generateExtractiveSummary(item);
      }
    }
    return generateExtractiveSummary(item);
  }

  /// Generate LLM-powered summary matching Golang implementation
  Future<Summary> generateLLMSummary(NewsItem item) async {
    if (aiService == null) {
      throw Exception('AI Service not configured');
    }

    String content = item.title;
    if (item.summary != null) {
      content += '\n\n${item.summary}';
    }

    final result = await aiService!.summarize(
      title: item.title,
      content: content,
    );

    final tags = extractTags(content);

    return Summary(
      itemId: item.id,
      tldr: result.tldr,
      whyItMatters: result.whyItMatters,
      tags: tags,
      method: 'llm',
    );
  }

  /// Generate extractive summary as fallback matching Golang implementation
  Summary generateExtractiveSummary(NewsItem item) {
    String text = item.title;
    if (item.summary != null) {
      text += ' ${item.summary}';
    }

    // Extract sentences
    final sentences = _extractSentences(text);
    String tldr = '';

    if (sentences.length >= 3) {
      tldr = sentences.sublist(0, 3).join(' ');
    } else if (sentences.isNotEmpty) {
      tldr = sentences.join(' ');
    } else {
      // Fallback: first 200 chars
      if (text.length > 200) {
        tldr = '${text.substring(0, 200)}...';
      } else {
        tldr = text;
      }
    }

    final whyItMatters = generateWhyItMatters(text);
    final tags = extractTags(text);

    return Summary(
      itemId: item.id,
      tldr: tldr,
      whyItMatters: whyItMatters,
      tags: tags,
      method: 'extractive',
    );
  }

  /// Extract tags based on topic keywords matching Golang implementation
  List<String> extractTags(String text) {
    final textLower = text.toLowerCase();
    final tags = <String>[];

    for (final entry in topicKeywords.entries) {
      final topic = entry.key;
      final keywords = entry.value;

      for (final keyword in keywords) {
        if (textLower.contains(keyword)) {
          tags.add(topic);
          break; // Only add topic once
        }
      }
    }

    // If no tags found, add general_ai
    if (tags.isEmpty) {
      tags.add('general_ai');
    }

    return tags;
  }

  /// Generate "why it matters" text matching Golang implementation
  String generateWhyItMatters(String text) {
    final textLower = text.toLowerCase();

    if (_contains(textLower, 'llm') ||
        _contains(textLower, 'transformer') ||
        _contains(textLower, 'gpt')) {
      return 'This development could impact how AI engineers build and deploy language models, potentially affecting inference costs, model architecture choices, and RAG system design.';
    }

    if (_contains(textLower, 'mlops') ||
        _contains(textLower, 'deployment')) {
      return 'For ML engineers, this addresses critical production challenges around model deployment, monitoring, and maintaining model performance in real-world environments.';
    }

    if (_contains(textLower, 'computer vision') ||
        _contains(textLower, 'detection')) {
      return 'Advances in computer vision directly impact applications in autonomous systems, medical imaging, and industrial automation, requiring engineers to stay updated on state-of-the-art techniques.';
    }

    if (_contains(textLower, 'rag') || _contains(textLower, 'retrieval')) {
      return 'This could improve how AI systems access and utilize external knowledge, which is crucial for building more capable and accurate AI applications.';
    }

    // Default
    return 'Staying informed about AI/ML developments helps engineers make better technical decisions, adopt new tools and techniques, and understand the evolving landscape of machine learning.';
  }

  /// Extract sentences from text matching Golang implementation
  List<String> _extractSentences(String text) {
    text = text.trim();
    if (text.isEmpty) return [];

    final sentences = <String>[];
    String current = '';

    for (int i = 0; i < text.length; i++) {
      final char = text[i];
      current += char;

      if (char == '.' || char == '!' || char == '?') {
        final sentence = current.trim();
        if (sentence.isNotEmpty) {
          sentences.add(sentence);
        }
        current = '';
      }
    }

    // Add remaining text
    if (current.trim().isNotEmpty) {
      sentences.add(current.trim());
    }

    // Filter out very short sentences
    return sentences.where((s) => s.length > 20).toList();
  }

  bool _contains(String text, String substring) {
    return text.contains(substring);
  }
}
