class AIConfig {
  final String provider;
  final String model;
  final String apiKey;
  final int maxTokens;
  final double temperature;

  AIConfig({
    this.provider = 'openrouter',
    this.model = 'openai/gpt-3.5-turbo',
    required this.apiKey,
    this.maxTokens = 500,
    this.temperature = 0.7,
  });

  AIConfig copyWith({
    String? provider,
    String? model,
    String? apiKey,
    int? maxTokens,
    double? temperature,
  }) {
    return AIConfig(
      provider: provider ?? this.provider,
      model: model ?? this.model,
      apiKey: apiKey ?? this.apiKey,
      maxTokens: maxTokens ?? this.maxTokens,
      temperature: temperature ?? this.temperature,
    );
  }
}
