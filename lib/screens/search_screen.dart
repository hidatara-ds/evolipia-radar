import 'package:flutter/material.dart';
import '../theme.dart';
import '../models/news_item.dart';
import '../widgets/news_card.dart';
import '../services/api_service.dart';

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  final TextEditingController _searchController = TextEditingController();
  List<NewsItem> searchResults = [];
  bool isLoading = false;
  String? error;
  bool hasSearched = false;

  Future<void> performSearch(String query) async {
    if (query.trim().isEmpty) return;

    setState(() {
      isLoading = true;
      error = null;
      hasSearched = true;
    });

    try {
      final results = await ApiService.search(query.trim());
      
      setState(() {
        searchResults = results;
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
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          TextField(
            controller: _searchController,
            onSubmitted: performSearch,
            decoration: InputDecoration(
              hintText: 'Cari tren, topik, atau kata kunci...',
              hintStyle: const TextStyle(color: AppTheme.muted),
              filled: true,
              fillColor: AppTheme.surface,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: AppTheme.border),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: AppTheme.border),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: AppTheme.accent, width: 2),
              ),
              prefixIcon: const Icon(Icons.search, color: AppTheme.muted),
              suffixIcon: _searchController.text.isNotEmpty
                  ? IconButton(
                      icon: const Icon(Icons.clear, color: AppTheme.muted),
                      onPressed: () {
                        _searchController.clear();
                        setState(() {
                          searchResults.clear();
                          hasSearched = false;
                          error = null;
                        });
                      },
                    )
                  : null,
            ),
            style: const TextStyle(color: AppTheme.textColor),
          ),
          const SizedBox(height: 24),
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
                          'Mencari...',
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
                                'Gagal mencari',
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
                                onPressed: () => performSearch(_searchController.text),
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
                    : hasSearched && searchResults.isEmpty
                        ? Center(
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  '🔍',
                                  style: TextStyle(
                                    fontSize: 48,
                                    color: AppTheme.muted.withValues(alpha: 0.5),
                                  ),
                                ),
                                const SizedBox(height: 16),
                                const Text(
                                  'Tidak ada hasil',
                                  style: TextStyle(color: AppTheme.muted),
                                ),
                                const SizedBox(height: 8),
                                Text(
                                  'Coba kata kunci lain',
                                  style: TextStyle(
                                    color: AppTheme.muted.withValues(alpha: 0.7),
                                    fontSize: 14,
                                  ),
                                ),
                              ],
                            ),
                          )
                        : hasSearched
                            ? ListView.builder(
                                itemCount: searchResults.length,
                                itemBuilder: (context, index) {
                                  return NewsCard(item: searchResults[index]);
                                },
                              )
                            : Center(
                                child: Column(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Text(
                                      '🔍',
                                      style: TextStyle(
                                        fontSize: 48,
                                        color: AppTheme.muted.withValues(alpha: 0.5),
                                      ),
                                    ),
                                    const SizedBox(height: 16),
                                    const Text(
                                      'Cari berita atau topik',
                                      style: TextStyle(color: AppTheme.muted),
                                    ),
                                  ],
                                ),
                              ),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }
}
