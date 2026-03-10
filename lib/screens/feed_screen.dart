import 'package:flutter/material.dart';
import 'package:postgres/postgres.dart';
import 'package:intl/intl.dart';
import 'package:url_launcher/url_launcher.dart';
import '../config.dart';
import '../theme.dart';
import '../models/news_item.dart';
import '../widgets/news_card.dart';

class FeedScreen extends StatefulWidget {
  const FeedScreen({super.key});

  @override
  State<FeedScreen> createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
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
        ORDER BY i.published_at DESC
        LIMIT 50
      ''');

      await connection.close();

      print('DEBUG: Query returned ${results.length} rows'); // Debug log

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
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Top Hari Ini',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: AppTheme.textColor,
                ),
              ),
              SizedBox(height: 4),
              Text(
                'Berita & tren terpopuler',
                style: TextStyle(
                  fontSize: 14,
                  color: AppTheme.muted,
                ),
              ),
            ],
          ),
        ),
        Expanded(
          child: isLoading
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(
                        width: 40,
                        height: 40,
                        child: CircularProgressIndicator(
                          strokeWidth: 3,
                          valueColor: const AlwaysStoppedAnimation<Color>(
                            AppTheme.accent,
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      const Text(
                        'Memuat berita...',
                        style: TextStyle(color: AppTheme.muted),
                      ),
                    ],
                  ),
                )
              : error != null
                  ? Center(
                      child: Padding(
                        padding: const EdgeInsets.all(24),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Text('⚠️', style: TextStyle(fontSize: 48)),
                            const SizedBox(height: 16),
                            const Text(
                              'Gagal memuat data',
                              style: TextStyle(
                                color: AppTheme.textColor,
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              error!,
                              style: const TextStyle(
                                color: AppTheme.muted,
                                fontSize: 14,
                              ),
                              textAlign: TextAlign.center,
                            ),
                            const SizedBox(height: 24),
                            ElevatedButton(
                              onPressed: loadNews,
                              style: ElevatedButton.styleFrom(
                                backgroundColor: AppTheme.accent,
                                foregroundColor: AppTheme.bg,
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 24,
                                  vertical: 14,
                                ),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(12),
                                ),
                              ),
                              child: const Text(
                                'Retry',
                                style: TextStyle(fontWeight: FontWeight.w600),
                              ),
                            ),
                          ],
                        ),
                      ),
                    )
                  : items.isEmpty
                      ? Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                '📭',
                                style: TextStyle(
                                  fontSize: 48,
                                  color: AppTheme.muted.withOpacity(0.5),
                                ),
                              ),
                              const SizedBox(height: 16),
                              const Text(
                                'Belum ada berita hari ini',
                                style: TextStyle(color: AppTheme.muted),
                              ),
                              const SizedBox(height: 8),
                              Text(
                                'Coba lagi nanti',
                                style: TextStyle(
                                  color: AppTheme.muted.withOpacity(0.7),
                                  fontSize: 14,
                                ),
                              ),
                            ],
                          ),
                        )
                      : RefreshIndicator(
                          onRefresh: loadNews,
                          color: AppTheme.accent,
                          backgroundColor: AppTheme.surface,
                          child: ListView.builder(
                            padding: const EdgeInsets.all(16),
                            itemCount: items.length,
                            itemBuilder: (context, index) {
                              return NewsCard(item: items[index]);
                            },
                          ),
                        ),
        ),
      ],
    );
  }
}
