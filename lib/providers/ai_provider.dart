import 'package:flutter/foundation.dart';
import '../services/ai_service.dart';
import '../services/summarizer_service.dart';

class AIProvider extends ChangeNotifier {
  AIService? _aiService;
  SummarizerService? _summarizerService;
  String? _apiKey;
  bool _isLoading = false;

  AIService? get aiService => _aiService;
  SummarizerService? get summarizerService => _summarizerService;
  String? get apiKey => _apiKey;
  bool get isConfigured => _apiKey != null && _apiKey!.isNotEmpty;
  bool get isLoading => _isLoading;

  /// Initialize provider and load saved API key
  Future<void> initialize() async {
    _isLoading = true;
    notifyListeners();

    // For demo purposes, no persistent storage
    // In production, you'd load from secure storage
    
    _isLoading = false;
    notifyListeners();
  }

  /// Set API key and initialize services
  Future<void> setApiKey(String key, {bool saveToStorage = true}) async {
    _apiKey = key;

    if (key.isEmpty) {
      _aiService = null;
      _summarizerService = null;
    } else {
      _aiService = AIService(apiKey: key);
      _summarizerService = SummarizerService(aiService: _aiService);
    }

    // For demo purposes, no persistent storage
    // In production, you'd save to secure storage

    notifyListeners();
  }

  /// Clear API key
  Future<void> clearApiKey() async {
    await setApiKey('');
  }

  @override
  void dispose() {
    _aiService?.dispose();
    super.dispose();
  }
}
