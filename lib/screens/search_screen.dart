import 'package:flutter/material.dart';
import '../theme.dart';

class SearchScreen extends StatelessWidget {
  const SearchScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          TextField(
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
            ),
            style: const TextStyle(color: AppTheme.textColor),
          ),
          const SizedBox(height: 24),
          Expanded(
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    '🔍',
                    style: TextStyle(
                      fontSize: 48,
                      color: AppTheme.muted.withOpacity(0.5),
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
}
