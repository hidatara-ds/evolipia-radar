import 'package:flutter/foundation.dart';
import '../models/ai_message.dart';
import '../services/chat_service.dart';

class ChatProvider extends ChangeNotifier {
  final ChatService chatService;
  List<ChatMessage> _messages = [];
  bool _isLoading = false;
  String? _error;

  ChatProvider({required this.chatService});

  List<ChatMessage> get messages => _messages;
  bool get isLoading => _isLoading;
  String? get error => _error;
  bool get hasMessages => _messages.isNotEmpty;

  /// Initialize and load chat history
  Future<void> initialize() async {
    _messages = await chatService.loadHistory();
    notifyListeners();
  }

  /// Send message to AI
  Future<void> sendMessage(String text) async {
    if (text.trim().isEmpty || _isLoading) return;

    _isLoading = true;
    _error = null;

    // Add user message
    final userMessage = ChatMessage(
      role: 'user',
      content: text,
      timestamp: DateTime.now(),
    );
    _messages.add(userMessage);
    notifyListeners();

    try {
      // Get AI response
      final response = await chatService.sendMessage(text, _messages);

      // Add assistant message
      final assistantMessage = ChatMessage(
        role: 'assistant',
        content: response,
        timestamp: DateTime.now(),
      );
      _messages.add(assistantMessage);

      // Save history
      await chatService.saveHistory(_messages);
    } catch (e) {
      _error = e.toString();

      // Add error message
      final errorMessage = ChatMessage(
        role: 'assistant',
        content: 'Error: ${e.toString()}. Coba periksa API key atau coba lagi.',
        timestamp: DateTime.now(),
        isError: true,
      );
      _messages.add(errorMessage);
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Clear chat history
  Future<void> clearHistory() async {
    _messages.clear();
    await chatService.clearHistory();
    notifyListeners();
  }

  /// Reload history from storage
  Future<void> reloadHistory() async {
    _messages = await chatService.loadHistory();
    notifyListeners();
  }
}
