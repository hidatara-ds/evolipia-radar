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
          return items.map((item) => NewsItem.fromJson(item)).toList();
        } else {
          return _getDummyTrending();
        }
      } else {
        return _getDummyTrending();
      }
    } catch (e) {
      return _getDummyTrending();
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
    return [
      NewsItem(
        id: '1',
        title: 'Tony Hoare has died',
        url: 'https://blog.computationalcomplexity.org/2026/03/tony-hoare-1934-2026.html',
        domain: 'blog.computationalcomplexity.org',
        publishedAt: DateTime.now().subtract(const Duration(hours: 2)),
        category: 'tech',
        score: 9.0,
        tldr: 'Computer science pioneer Tony Hoare, inventor of Quicksort algorithm, has passed away.',
        whyItMatters: 'Tony Hoare made fundamental contributions to computer science that are still used today.',
        tags: ['computer-science', 'algorithms'],
      ),
      NewsItem(
        id: '2',
        title: 'Launch HN: RunAnywhere (YC W26) – Faster AI Inference on Apple Silicon',
        url: 'https://github.com/RunanywhereAI/rcli',
        domain: 'github.com',
        publishedAt: DateTime.now().subtract(const Duration(hours: 4)),
        category: 'tech',
        score: 7.0,
        tldr: 'New YC startup optimizes AI inference performance on Apple Silicon chips.',
        whyItMatters: 'Could significantly improve AI performance on Mac devices.',
        tags: ['ai', 'apple', 'performance'],
      ),
      NewsItem(
        id: '3',
        title: 'After outages, Amazon to make senior engineers sign off on AI-assisted changes',
        url: 'https://arstechnica.com/ai/2026/03/after-outages-amazon-to-make-senior-engineers-sign-off-on-ai-assisted-changes',
        domain: 'arstechnica.com',
        publishedAt: DateTime.now().subtract(const Duration(hours: 6)),
        category: 'tech',
        score: 7.0,
        tldr: 'Amazon implements new review process for AI-generated code changes after recent outages.',
        whyItMatters: 'Shows growing concerns about AI code quality and reliability in production systems.',
        tags: ['ai', 'amazon', 'engineering'],
      ),
      NewsItem(
        id: '4',
        title: 'Redox OS has adopted a Certificate of Origin policy and a strict no-LLM policy',
        url: 'https://gitlab.redox-os.org/redox-os/redox/-/blob/master/CONTRIBUTING.md',
        domain: 'gitlab.redox-os.org',
        publishedAt: DateTime.now().subtract(const Duration(hours: 8)),
        category: 'tech',
        score: 6.0,
        tldr: 'Redox OS implements new policies regarding AI-generated contributions.',
        whyItMatters: 'Reflects growing debate about AI in open source development.',
        tags: ['open-source', 'ai', 'policy'],
      ),
      NewsItem(
        id: '5',
        title: 'Yann LeCun raises \$1B to build AI that understands the physical world',
        url: 'https://www.wired.com/story/yann-lecun-raises-dollar1-billion-to-build-ai-that-understands-the-physical-world',
        domain: 'wired.com',
        publishedAt: DateTime.now().subtract(const Duration(hours: 10)),
        category: 'tech',
        score: 8.0,
        tldr: 'AI pioneer secures massive funding for next-generation AI research.',
        whyItMatters: 'Could lead to breakthrough in AI understanding of physical reality.',
        tags: ['ai', 'funding', 'research'],
      ),
    ];
  }
  
  static List<NewsItem> _getDummyTrending() {
    return _getDummyNews().where((item) => item.score! > 7.0).toList();
  }
  
  static List<NewsItem> _getDummySearch(String query) {
    return _getDummyNews()
        .where((item) => 
            item.title.toLowerCase().contains(query.toLowerCase()) ||
            item.tags?.any((tag) => tag.toLowerCase().contains(query.toLowerCase())) == true)
        .toList();
  }
}