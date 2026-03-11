import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/news_item.dart';

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
          throw Exception(data['error'] ?? 'Failed to load news');
        }
      } else {
        throw Exception('HTTP ${response.statusCode}: ${response.reasonPhrase}');
      }
    } catch (e) {
      throw Exception('Failed to load news: $e');
    }
  }
  
  static Future<List<NewsItem>> getTrending() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/api/trending'));
      
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final items = data['data']['items'] as List;
          return items.map((item) => NewsItem.fromJson(item)).toList();
        } else {
          throw Exception(data['error'] ?? 'Failed to load trending');
        }
      } else {
        throw Exception('HTTP ${response.statusCode}: ${response.reasonPhrase}');
      }
    } catch (e) {
      throw Exception('Failed to load trending: $e');
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
          throw Exception(data['error'] ?? 'Failed to search');
        }
      } else {
        throw Exception('HTTP ${response.statusCode}: ${response.reasonPhrase}');
      }
    } catch (e) {
      throw Exception('Failed to search: $e');
    }
  }
}