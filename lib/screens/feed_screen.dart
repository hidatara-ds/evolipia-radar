import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:url_launcher/url_launcher.dart';
import '../theme.dart';
import '../models/news_item.dart';
import '../widgets/news_card.dart';
import '../services/api_service.dart';

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
      final newsItems = await ApiService.getNews();
      
      setState(() {
        items = newsItems;
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
                                  color: AppTheme.muted.withValues(alpha: 0.5),
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
                                  color: AppTheme.muted.withValues(alpha: 0.7),
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
