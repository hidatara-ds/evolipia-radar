import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/ai_message.dart';

class AIService {
  final String apiKey;
  final String baseUrl = 'https://openrouter.ai/api/v1';
  final http.Client _client;

  AIService({
    required this.apiKey,
    http.Client? client,
  }) : _client = client ?? http.Client();

  /// Generic completion method matching Golang implementation
  Future<String> complete({
    required String model,
    required List<AIMessage> messages,
    int maxTokens = 500,
    double temperature = 0.7,
  }) async {
    final url = Uri.parse('$baseUrl/chat/completions');

    final body = jsonEncode({
      'model': model,
      'messages': messages.map((m) => m.toJson()).toList(),
      'max_tokens': maxTokens,
      'temperature': temperature,
    });

    try {
      final response = await _client
          .post(
            url,
            headers: {
              'Content-Type': 'application/json',
              'Authorization': 'Bearer $apiKey',
              'HTTP-Referer': 'https://github.com/hidatara-ds/evolipia-radar',
              'X-Title': 'Evolipia Radar',
            },
            body: body,
          )
          .timeout(const Duration(seconds: 60));

      if (response.statusCode != 200) {
        final error = jsonDecode(response.body);
        throw AIException(
          error['error']?['message'] ?? 'API Error (${response.statusCode})',
        );
      }

      final data = jsonDecode(response.body);

      if (data['error'] != null) {
        throw AIException(data['error']['message'] ?? 'API Error');
      }

      if (data['choices'] == null || (data['choices'] as List).isEmpty) {
        throw AIException('No completion choices returned');
      }

      return data['choices'][0]['message']['content'] as String;
    } on TimeoutException {
      throw AIException('Request timeout');
    } on http.ClientException catch (e) {
      throw AIException('Network error: ${e.message}');
    } catch (e) {
      if (e is AIException) rethrow;
      throw AIException('Failed to complete: $e');
    }
  }

  /// Summarize news article matching Golang implementation
  Future<SummaryResult> summarize({
    required String title,
    required String content,
    String model = 'openai/gpt-3.5-turbo',
    int maxTokens = 500,
    double temperature = 0.7,
  }) async {
    final prompt = '''Summarize this AI/ML news article:

Title: $title
Content: $content

Provide:
1. A 2-sentence summary (TLDR)
2. One sentence explaining why this matters to AI/ML engineers

Format your response as:
TLDR: [your summary]
WHY: [why it matters]''';

    final messages = [
      AIMessage(
        role: 'system',
        content:
            'You are an AI/ML news analyst. Provide concise, technical summaries focused on engineering impact.',
      ),
      AIMessage(
        role: 'user',
        content: prompt,
      ),
    ];

    final response = await complete(
      model: model,
      messages: messages,
      maxTokens: maxTokens,
      temperature: temperature,
    );

    return _parseSummaryResponse(response);
  }

  /// Parse summary response matching Golang logic
  SummaryResult _parseSummaryResponse(String response) {
    String tldr = '';
    String why = '';

    final lines = response.split('\n');

    for (final line in lines) {
      final trimmed = line.trim();
      if (trimmed.startsWith('TLDR:')) {
        tldr = trimmed.substring(5).trim();
      } else if (trimmed.startsWith('WHY:')) {
        why = trimmed.substring(4).trim();
      }
    }

    // Fallback if parsing fails
    if (tldr.isEmpty) {
      tldr = response;
    }

    return SummaryResult(tldr: tldr, whyItMatters: why);
  }

  void dispose() {
    _client.close();
  }
}

class SummaryResult {
  final String tldr;
  final String whyItMatters;

  SummaryResult({
    required this.tldr,
    required this.whyItMatters,
  });
}

class AIException implements Exception {
  final String message;

  AIException(this.message);

  @override
  String toString() => message;
}
