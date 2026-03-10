import 'package:flutter/material.dart';
import '../theme.dart';

class TrendingScreen extends StatelessWidget {
  const TrendingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text('🔥', style: TextStyle(fontSize: 48)),
          SizedBox(height: 16),
          Text(
            'Sedang Naik',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w700,
              color: AppTheme.textColor,
            ),
          ),
          SizedBox(height: 8),
          Text(
            'Fitur coming soon',
            style: TextStyle(color: AppTheme.muted),
          ),
        ],
      ),
    );
  }
}
