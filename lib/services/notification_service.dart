import 'dart:math';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/news_item.dart';

class NotificationService {
  static final FlutterLocalNotificationsPlugin _notifications = 
      FlutterLocalNotificationsPlugin();
  
  static bool _initialized = false;
  
  static Future<void> initialize() async {
    if (_initialized) return;
    
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );
    
    const initSettings = InitializationSettings(
      android: androidSettings,
      iOS: iosSettings,
    );
    
    await _notifications.initialize(
      initSettings,
      onDidReceiveNotificationResponse: _onNotificationTapped,
    );
    
    // Request permissions for Android 13+
    await _notifications
        .resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()
        ?.requestNotificationsPermission();
    
    _initialized = true;
  }
  
  static void _onNotificationTapped(NotificationResponse response) {
    // Handle notification tap - could navigate to specific news item
    print('Notification tapped: ${response.payload}');
  }
  
  static Future<void> showTrendingNotification(NewsItem item) async {
    await initialize();
    
    const androidDetails = AndroidNotificationDetails(
      'trending_news',
      'Trending News',
      channelDescription: 'Notifications for trending news items',
      importance: Importance.high,
      priority: Priority.high,
      icon: '@mipmap/ic_launcher',
      color: Color(0xFF00D4FF), // AppTheme.accent
      playSound: true,
      enableVibration: true,
    );
    
    const iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );
    
    const details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );
    
    await _notifications.show(
      Random().nextInt(1000), // Random ID
      '🔥 Trending: ${item.title}',
      '⭐ ${item.score!.toInt()}/10 • ${item.domain}',
      details,
      payload: item.id,
    );
  }
  
  static Future<void> showDailyDigestNotification(List<NewsItem> topItems) async {
    await initialize();
    
    const androidDetails = AndroidNotificationDetails(
      'daily_digest',
      'Daily Digest',
      channelDescription: 'Daily news digest notifications',
      importance: Importance.default_,
      priority: Priority.default_,
      icon: '@mipmap/ic_launcher',
      color: Color(0xFF00D4FF),
      playSound: true,
    );
    
    const iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );
    
    const details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );
    
    final topScore = topItems.isNotEmpty ? topItems.first.score!.toInt() : 0;
    
    await _notifications.show(
      999, // Fixed ID for daily digest
      '📰 Daily Digest Ready',
      '${topItems.length} new stories • Top score: $topScore/10',
      details,
      payload: 'daily_digest',
    );
  }
  
  static Future<void> schedulePeriodicChecks() async {
    await initialize();
    
    // Check for trending news every 30 minutes
    // Note: For production, you'd want to use background tasks or Firebase Functions
    // This is a simplified version for demo purposes
  }
  
  static Future<bool> areNotificationsEnabled() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool('notifications_enabled') ?? true;
  }
  
  static Future<void> setNotificationsEnabled(bool enabled) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('notifications_enabled', enabled);
  }
  
  static Future<void> cancelAll() async {
    await _notifications.cancelAll();
  }
}