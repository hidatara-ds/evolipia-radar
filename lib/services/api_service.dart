import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/news_item.dart';
import 'notification_service.dart';

class ApiService {
  // Use your deployed API URL here
  // For local testing: 'http://localhost:3000'
  // For Vercel: 'https://your-app.vercel.app'
  static const String baseUrl = 'https://evolipia-radar.vercel.app';
  
  static Future<List<NewsItem>> getNews({String? topic, String? date}) async {
    try {
      final uri = Uri.parse('$baseUrl/api/news').replace(
        queryParameters: {
          if (topic != null) 'topic': topic,
          if (date != null) 'date': date,
        },
      );
      
      final response = await http.get(uri);
      
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final items = data['data']['items'] as List;
          return items.map((item) => NewsItem.fromJson(item)).toList();
        } else {
          // If API fails, return dummy data for now
          return _getDummyNews();
        }
      } else {
        // If API fails, return dummy data for now
        return _getDummyNews();
      }
    } catch (e) {
      // If API fails, return dummy data for now
      return _getDummyNews();
    }
  }
  
  static Future<List<NewsItem>> getTrending() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/api/trending'));
      
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final items = data['data']['items'] as List;
          final trendingItems = items.map((item) => NewsItem.fromJson(item)).toList();
          
          // Check for high-scoring items and send notifications
          await _checkForTrendingNotifications(trendingItems);
          
          return trendingItems;
        } else {
          final dummyItems = _getDummyTrending();
          await _checkForTrendingNotifications(dummyItems);
          return dummyItems;
        }
      } else {
        final dummyItems = _getDummyTrending();
        await _checkForTrendingNotifications(dummyItems);
        return dummyItems;
      }
    } catch (e) {
      final dummyItems = _getDummyTrending();
      await _checkForTrendingNotifications(dummyItems);
      return dummyItems;
    }
  }
  
  static Future<void> _checkForTrendingNotifications(List<NewsItem> items) async {
    final notificationsEnabled = await NotificationService.areNotificationsEnabled();
    if (!notificationsEnabled) return;
    
    // Send notification for items with score >= 8
    for (final item in items) {
      if (item.score != null && item.score! >= 8) {
        await NotificationService.showTrendingNotification(item);
        break; // Only send one notification to avoid spam
      }
    }
  }
  
  static Future<List<NewsItem>> search(String query) async {
    try {
      final uri = Uri.parse('$baseUrl/api/search').replace(
        queryParameters: {'q': query},
      );
      
      final response = await http.get(uri);
      
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final items = data['data']['items'] as List;
          return items.map((item) => NewsItem.fromJson(item)).toList();
        } else {
          return _getDummySearch(query);
        }
      } else {
        return _getDummySearch(query);
      }
    } catch (e) {
      return _getDummySearch(query);
    }
  }
  
  // Dummy data for fallback
  static List<NewsItem> _getDummyNews() {
    final now = DateTime.now();
    return [
      NewsItem(
        id: '1',
        title: 'OpenAI releases GPT-5 with breakthrough reasoning capabilities',
        url: 'https://openai.com/blog/gpt-5-announcement',
        domain: 'openai.com',
        publishedAt: now.subtract(const Duration(minutes: 30)),
        category: 'tech',
        score: 10,
        tldr: 'OpenAI unveils GPT-5 with revolutionary reasoning and multimodal capabilities.',
        whyItMatters: 'This could be the biggest leap in AI since ChatGPT, affecting every industry.',
        tags: ['ai', 'openai', 'breakthrough'],
      ),
      NewsItem(
        id: '2',
        title: 'Meta announces new AR glasses with neural interface',
        url: 'https://about.meta.com/news/ar-neural-interface',
        domain: 'meta.com',
        publishedAt: now.subtract(const Duration(hours: 1)),
        category: 'tech',
        score: 9,
        tldr: 'Meta reveals AR glasses that can be controlled by thought patterns.',
        whyItMatters: 'Could revolutionize human-computer interaction and accessibility.',
        tags: ['ar', 'meta', 'neural-interface'],
      ),
      NewsItem(
        id: '3',
        title: 'Google announces Quantum Computer breakthrough for cryptography',
        url: 'https://blog.google/technology/quantum-computing/breakthrough-2026',
        domain: 'google.com',
        publishedAt: now.subtract(const Duration(hours: 2)),
        category: 'tech',
        score: 8,
        tldr: 'Google achieves quantum supremacy in cryptographic applications.',
        whyItMatters: 'Will require complete overhaul of current encryption methods.',
        tags: ['quantum', 'google', 'cryptography'],
      ),
      NewsItem(
        id: '4',
        title: 'Apple unveils M5 chip with 50% performance boost',
        url: 'https://www.apple.com/newsroom/2026/03/m5-chip-announcement',
        domain: 'apple.com',
        publishedAt: now.subtract(const Duration(hours: 4)),
        category: 'tech',
        score: 7,
        tldr: 'Apple M5 chip delivers unprecedented performance for mobile devices.',
        whyItMatters: 'Sets new standard for mobile computing and AI processing.',
        tags: ['apple', 'm5', 'performance'],
      ),
      NewsItem(
        id: '5',
        title: 'Tesla achieves full self-driving milestone in urban environments',
        url: 'https://www.tesla.com/blog/fsd-urban-milestone',
        domain: 'tesla.com',
        publishedAt: now.subtract(const Duration(hours: 6)),
        category: 'tech',
        score: 8,
        tldr: 'Tesla FSD successfully navigates complex urban scenarios without intervention.',
        whyItMatters: 'Major step toward fully autonomous vehicles becoming mainstream.',
        tags: ['tesla', 'autonomous', 'fsd'],
      ),
    ];
  }
  
  static List<NewsItem> _getDummyTrending() {
    return _getDummyNews().where((item) => item.score! > 7).toList();
  }
  
  static List<NewsItem> _getDummySearch(String query) {
    return _getDummyNews()
        .where((item) => 
            item.title.toLowerCase().contains(query.toLowerCase()) ||
            item.tags?.any((tag) => tag.toLowerCase().contains(query.toLowerCase())) == true)
        .toList();
  }
}