import 'package:flutter/material.dart';
import 'package:postgres/postgres.dart';
import 'package:intl/intl.dart';
import 'package:url_launcher/url_launcher.dart';
import 'config.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Evolipia',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const NewsListPage(),
    );
  }
}

class NewsListPage extends StatefulWidget {
  const NewsListPage({super.key});

  @override
  State<NewsListPage> createState() => _NewsListPageState();
}

class _NewsListPageState extends State<NewsListPage> {
  List<NewsItem> items = [];
  bool isLoading = true;
  String? error;

  @override
  void initState() {
    super.initState();
    loadNews();
  }

  Future<void> loadNews() async {
    setState(() {
      isLoading = true;
      error = null;
    });

    try {
      final connection = await Connection.open(
        Endpoint(
          host: DatabaseConfig.host,
          database: DatabaseConfig.database,
          username: DatabaseConfig.username,
          password: DatabaseConfig.password,
          port: 5432,
        ),
        settings: const ConnectionSettings(sslMode: SslMode.require),
      );

      final results = await connection.execute('''
        SELECT 
          i.id,
          i.title,
          i.url,
          i.domain,
          i.published_at,
          i.category,
          s.final as score,
          sm.tldr as summary
        FROM items i
        LEFT JOIN scores s ON i.id = s.item_id
        LEFT JOIN summaries sm ON i.id = sm.item_id
        WHERE i.published_at >= NOW() - INTERVAL '7 days'
        ORDER BY s.final DESC NULLS LAST, i.published_at DESC
        LIMIT 50
      ''');

      await connection.close();

      setState(() {
        items = results.map((row) {
          return NewsItem(
            id: row[0] as String,
            title: row[1] as String,
            url: row[2] as String,
            domain: row[3] as String,
            publishedAt: row[4] as DateTime,
            category: row[5] as String,
            score: row[6] as double?,
            summary: row[7] as String?,
          );
        }).toList();
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        error = e.toString();
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Evolipia'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: loadNews,
          ),
        ],
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.error, size: 48, color: Colors.red),
                      const SizedBox(height: 16),
                      Text('Error: $error'),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: loadNews,
                        child: const Text('Retry'),
                      ),
                    ],
                  ),
                )
              : items.isEmpty
                  ? const Center(child: Text('No news yet'))
                  : RefreshIndicator(
                      onRefresh: loadNews,
                      child: ListView.builder(
                        itemCount: items.length,
                        itemBuilder: (context, index) {
                          final item = items[index];
                          return NewsCard(item: item);
                        },
                      ),
                    ),
    );
  }
}

class NewsCard extends StatelessWidget {
  final NewsItem item;

  const NewsCard({super.key, required this.item});

  Future<void> _launchUrl() async {
    final uri = Uri.parse(item.url);
    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
      throw Exception('Could not launch ${item.url}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final dateFormat = DateFormat('MMM dd, HH:mm');

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: InkWell(
        onTap: _launchUrl,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: _getCategoryColor(item.category),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      item.category.toUpperCase(),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  if (item.score != null)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.orange,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        '${item.score!.toStringAsFixed(0)}',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  const Spacer(),
                  Text(
                    item.domain,
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey[600],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                item.title,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              if (item.summary != null) ...[
                const SizedBox(height: 8),
                Text(
                  item.summary!,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.grey[700],
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
              const SizedBox(height: 8),
              Text(
                dateFormat.format(item.publishedAt),
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey[600],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getCategoryColor(String category) {
    switch (category.toLowerCase()) {
      case 'tech':
        return Colors.blue;
      case 'ai':
        return Colors.purple;
      default:
        return Colors.grey;
    }
  }
}

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
