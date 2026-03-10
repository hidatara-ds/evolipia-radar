class AIMessage {
  final String role; // 'system', 'user', or 'assistant'
  final String content;

  AIMessage({
    required this.role,
    required this.content,
  });

  Map<String, dynamic> toJson() {
    return {
      'role': role,
      'content': content,
    };
  }

  factory AIMessage.fromJson(Map<String, dynamic> json) {
    return AIMessage(
      role: json['role'] as String,
      content: json['content'] as String,
    );
  }
}

class ChatMessage {
  final String role; // 'user' or 'assistant'
  final String content;
  final DateTime timestamp;
  final bool isError;

  ChatMessage({
    required this.role,
    required this.content,
    required this.timestamp,
    this.isError = false,
  });

  Map<String, dynamic> toJson() {
    return {
      'role': role,
      'content': content,
      'timestamp': timestamp.toIso8601String(),
      'error': isError,
    };
  }

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      role: json['role'] as String,
      content: json['content'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
      isError: json['error'] as bool? ?? false,
    );
  }
}
