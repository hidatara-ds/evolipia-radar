import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:intl/intl.dart';
import '../theme.dart';
import '../models/news_item.dart';

class DetailScreen extends StatelessWidget {
  final NewsItem item;

  const DetailScreen({super.key, required this.item});

  Future<void> _launchUrl() async {
    final uri = Uri.parse(item.url);
    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
      throw Exception('Could not launch ${item.url}');
    }
  }

  @override
  Widget build(BuildContext context) {
    final dateFormat = DateFormat('dd MMM yyyy, HH:mm');

    return Scaffold(
      backgroundColor: AppTheme.bg,
      appBar: AppBar(
        backgroundColor: const Color(0xF20B0F19),
        surfaceTintColor: Colors.transparent,
        leading: Container(
          margin: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: AppTheme.surface,
            borderRadius: BorderRadius.circular(10),
            border: Border.all(color: AppTheme.border, width: 1),
          ),
          child: IconButton(
            icon: const Icon(Icons.arrow_back, size: 20),
            color: AppTheme.textColor,
            onPressed: () => Navigator.pop(context),
            padding: EdgeInsets.zero,
          ),
        ),
        title: const Text(
          'Detail',
          style: TextStyle(
            fontSize: 17,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Title
            Text(
              item.title,
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w700,
                color: AppTheme.textColor,
                height: 1.4,
              ),
            ),
            const SizedBox(height: 16),

            // Link
            InkWell(
              onTap: _launchUrl,
              borderRadius: BorderRadius.circular(8),
              child: Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 8,
                ),
                decoration: BoxDecoration(
                  color: AppTheme.accent.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(
                    color: AppTheme.accent.withValues(alpha: 0.2),
                    width: 1,
                  ),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(
                      Icons.link,
                      size: 16,
                      color: AppTheme.accent,
                    ),
                    const SizedBox(width: 6),
                    Flexible(
                      child: Text(
                        item.domain,
                        style: const TextStyle(
                          color: AppTheme.accent,
                          fontSize: 14,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),

            // Informasi Section
            _buildSection(
              'INFORMASI',
              Container(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildInfoRow('Domain:', item.domain),
                    const SizedBox(height: 12),
                    _buildInfoRow('Kategori:', item.category),
                    const SizedBox(height: 12),
                    Text(
                      'Dipublikasikan: ${dateFormat.format(item.publishedAt)}',
                      style: const TextStyle(
                        color: AppTheme.muted,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Ringkasan Section
            if (item.summary != null) ...[
              const SizedBox(height: 24),
              _buildSection(
                'RINGKASAN',
                Container(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        item.summary!,
                        style: const TextStyle(
                          color: AppTheme.textColor,
                          fontSize: 15,
                          height: 1.6,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: AppTheme.bg,
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(
                            color: AppTheme.border,
                            width: 1,
                          ),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Mengapa penting:',
                              style: TextStyle(
                                color: AppTheme.textColor,
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              'Staying informed about AI/ML developments helps engineers make better technical decisions, adopt new tools and techniques, and understand the evolving landscape of machine learning.',
                              style: TextStyle(
                                color: AppTheme.muted,
                                fontSize: 14,
                                height: 1.5,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],

            // Analisis Skor Section
            if (item.score != null) ...[
              const SizedBox(height: 24),
              _buildSection(
                'ANALISIS SKOR (SKALA 1-10)',
                Container(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      // Top 3 scores
                      Row(
                        children: [
                          Expanded(
                            child: _buildScoreCard(
                              item.score!.toStringAsFixed(1),
                              'Final',
                              AppTheme.accent,
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: _buildScoreCard(
                              '0.0',
                              'Hot',
                              AppTheme.warning,
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: _buildScoreCard(
                              '0.2',
                              'Relevan',
                              AppTheme.success,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      // Additional scores
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: AppTheme.bg,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Column(
                          children: [
                            _buildScoreRow('Kredibilitas:', '0.5/10'),
                            const SizedBox(height: 12),
                            _buildScoreRow('Kebaruan:', '1.0/10'),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],

            const SizedBox(height: 100),
          ],
        ),
      ),
    );
  }

  Widget _buildSection(String title, Widget content) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: const TextStyle(
            fontSize: 11,
            color: AppTheme.muted,
            fontWeight: FontWeight.w600,
            letterSpacing: 0.5,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          width: double.infinity,
          decoration: BoxDecoration(
            color: AppTheme.surface,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: AppTheme.border, width: 1),
          ),
          child: content,
        ),
      ],
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: AppTheme.textColor,
            fontSize: 15,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            value,
            style: const TextStyle(
              color: AppTheme.muted,
              fontSize: 15,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildScoreCard(String score, String label, Color color) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppTheme.bg,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          Text(
            score,
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.w700,
              color: color,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: const TextStyle(
              fontSize: 12,
              color: AppTheme.muted,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildScoreRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 14,
            color: AppTheme.muted,
          ),
        ),
        Text(
          value,
          style: const TextStyle(
            fontSize: 14,
            color: AppTheme.textColor,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}
