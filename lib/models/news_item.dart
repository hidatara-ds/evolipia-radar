class NewsItem {
  final String id;
  final String title;
  final String url;
  final String domain;
  final DateTime publishedAt;
  final String category;
  final double? score;
  final String? summary;
  final String? tldr;
  final String? whyItMatters;
  final List<String>? tags;

  NewsItem({
    required this.id,
    required this.title,
    required this.url,
    required this.domain,
    required this.publishedAt,
    required this.category,
    this.score,
    this.summary,
    this.tldr,
    this.whyItMatters,
    this.tags,
  });

  factory NewsItem.fromJson(Map<String, dynamic> json) {
    return NewsItem(
      id: json['id'] as String,
      title: json['title'] as String,
      url: json['url'] as String,
      domain: json['domain'] as String,
      publishedAt: DateTime.parse(json['published_at'] as String),
      category: json['category'] as String,
      score: json['score']?.toDouble(),
      summary: json['summary'] as String?,
      tldr: json['tldr'] as String?,
      whyItMatters: json['why_it_matters'] as String?,
      tags: json['tags'] != null ? List<String>.from(json['tags']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'url': url,
      'domain': domain,
      'published_at': publishedAt.toIso8601String(),
      'category': category,
      'score': score,
      'summary': summary,
      'tldr': tldr,
      'why_it_matters': whyItMatters,
      'tags': tags,
    };
  }
}
