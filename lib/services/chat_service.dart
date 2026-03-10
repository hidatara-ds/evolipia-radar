import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/ai_message.dart';
import 'ai_service.dart';

class ChatService {
  final AIService aiService;
  static const String _storageKey = 'chat_history';
  static const int _historyLimit = 50;
  static const int _contextWindow = 10;

  /// System prompt matching Golang implementation
  static const String systemPrompt = '''Kamu adalah Evolipia Radar AI — asisten pintar untuk aplikasi radar tren.

KARAKTER:
- Bicara dalam Bahasa Indonesia natural, friendly, tapi profesional
- Respons singkat & padat (maksimal 3-4 kalimat kecuali diminta detail)
- Gunakan bullet points untuk list
- Fokus pada tren, teknologi, dan insight data

KEMAMPUAN:
- Jelaskan tren/sinyal yang sedang naik
- Berikan rekomendasi aksi berdasarkan data
- Ringkas artikel panjang jadi poin-poin kunci
- Bantu analisis kompetitor atau market

JAWAB LANGSUNG tanpa basa-basi.''';

  ChatService({required this.aiService});

  /// Send message and get AI response
  Future<String> sendMessage(
    String message,
    List<ChatMessage> history,
  ) async {
    // Build messages array with system prompt + context window
    final messages = <AIMessage>[
      AIMessage(role: 'system', content: systemPrompt),
    ];

    // Add last N messages for context
    final contextMessages = history.length > _contextWindow
        ? history.sublist(history.length - _contextWindow)
        : history;

    for (final msg in contextMessages) {
      messages.add(AIMessage(
        role: msg.role,
        content: msg.content,
      ));
    }

    // Add current user message
    messages.add(AIMessage(role: 'user', content: message));

    // Call AI service
    final response = await aiService.complete(
      model: 'openai/gpt-3.5-turbo',
      messages: messages,
      maxTokens: 500,
      temperature: 0.7,
    );

    return response;
  }

  /// Load chat history from storage
  Future<List<ChatMessage>> loadHistory() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final jsonString = prefs.getString(_storageKey);

      if (jsonString == null) return [];

      final List<dynamic> jsonList = jsonDecode(jsonString);
      return jsonList
          .map((json) => ChatMessage.fromJson(json as Map<String, dynamic>))
          .toList();
    } catch (e) {
      return [];
    }
  }

  /// Save chat history to storage
  Future<void> saveHistory(List<ChatMessage> history) async {
    try {
      // Limit history size
      final limitedHistory = history.length > _historyLimit
          ? history.sublist(history.length - _historyLimit)
          : history;

      final prefs = await SharedPreferences.getInstance();
      final jsonString = jsonEncode(
        limitedHistory.map((msg) => msg.toJson()).toList(),
      );
      await prefs.setString(_storageKey, jsonString);
    } catch (e) {
      // Ignore save errors
    }
  }

  /// Clear chat history
  Future<void> clearHistory() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.remove(_storageKey);
    } catch (e) {
      // Ignore clear errors
    }
  }
}
