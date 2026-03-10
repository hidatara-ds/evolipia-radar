class NewsItem {
  final String id;
  final String title;
  final String url;
  final String domain;
  final DateTime publishedAt;
  final String category;
  final double? score;
  final String? summary;

  NewsItem({
    required this.id,
    required this.title,
    required this.url,
    required this.domain,
    required this.publishedAt,
    required this.category,
    this.score,
    this.summary,
  });
}
