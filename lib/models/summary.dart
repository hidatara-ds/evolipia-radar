class Summary {
  final String itemId;
  final String tldr;
  final String whyItMatters;
  final List<String> tags;
  final String method; // 'extractive' or 'llm'

  Summary({
    required this.itemId,
    required this.tldr,
    required this.whyItMatters,
    required this.tags,
    required this.method,
  });

  factory Summary.fromJson(Map<String, dynamic> json) {
    return Summary(
      itemId: json['item_id'] as String,
      tldr: json['tldr'] as String,
      whyItMatters: json['why_it_matters'] as String,
      tags: (json['tags'] as List<dynamic>).cast<String>(),
      method: json['method'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'item_id': itemId,
      'tldr': tldr,
      'why_it_matters': whyItMatters,
      'tags': tags,
      'method': method,
    };
  }
}
