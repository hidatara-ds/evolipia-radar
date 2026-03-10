import 'package:flutter/material.dart';

class AppTheme {
  // Colors matching web design
  static const Color bg = Color(0xFF0B0F19);
  static const Color surface = Color(0xFF151B2B);
  static const Color surfaceHover = Color(0xFF1E2538);
  static const Color border = Color(0xFF2A3447);
  static const Color textColor = Color(0xFFF8FAFC);
  static const Color muted = Color(0xFF94A3B8);
  static const Color accent = Color(0xFF38BDF8);
  static const Color accentDim = Color(0xFF0284C7);
  static const Color danger = Color(0xFFEF4444);
  static const Color success = Color(0xFF22C55E);
  static const Color warning = Color(0xFFF59E0B);

  static ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      scaffoldBackgroundColor: bg,
      colorScheme: const ColorScheme.dark(
        primary: accent,
        secondary: accentDim,
        surface: surface,
        error: danger,
        onPrimary: bg,
        onSurface: textColor,
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Color(0xF20B0F19),
        elevation: 0,
        centerTitle: false,
        titleTextStyle: TextStyle(
          fontSize: 17,
          fontWeight: FontWeight.w800,
          letterSpacing: 0.2,
        ),
      ),
      cardTheme: CardThemeData(
        color: surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: const BorderSide(color: border, width: 1),
        ),
      ),
      textTheme: const TextTheme(
        bodyLarge: TextStyle(color: textColor),
        bodyMedium: TextStyle(color: textColor),
        bodySmall: TextStyle(color: muted),
      ),
    );
  }
}
