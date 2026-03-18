"use client";

import { useEffect, useState } from "react";
import { 
  Activity, 
  Database, 
  FileText, 
  BrainCircuit, 
  RefreshCw, 
  Zap,
  ExternalLink,
  TrendingUp,
  Clock,
  AlertCircle,
  Sparkles
} from "lucide-react";

// API Response Types
interface NewsItem {
  id: string;
  title: string;
  url: string;
  domain: string;
  published_at: string;
  category: string;
  score: number;
  tldr?: string;
  why_it_matters?: string;
  tags: string[];
}

interface NewsResponse {
  success: boolean;
  data?: {
    items: NewsItem[];
    total_count: number;
    last_updated: string;
  };
  error?: string;
}

interface Metrics {
  articles_processed: number;
  filtered_articles: number;
  api_hits: number;
  clusters: number;
  avg_cluster_score: number;
  top_cluster_titles: string[] | null;
}

// Topic filter configuration
const TOPICS = [
  { id: "all", label: "All", color: "gray" },
  { id: "llm", label: "LLM", color: "purple" },
  { id: "vision", label: "Vision", color: "blue" },
  { id: "data", label: "Data", color: "green" },
  { id: "security", label: "Security", color: "red" },
  { id: "rl", label: "RL", color: "yellow" },
  { id: "robotics", label: "Robotics", color: "orange" },
  { id: "ide", label: "IDE", color: "cyan" },
  { id: "free-credits", label: "Free Credits", color: "pink" },
] as const;

export default function Dashboard() {
  const [news, setNews] = useState<NewsItem[]>([]);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [newsLoading, setNewsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTopic, setSelectedTopic] = useState<string>("all");
  const [triggering, setTriggering] = useState(false);
  const [lastTrigger, setLastTrigger] = useState<Date | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const [baseUrl, setBaseUrl] = useState("");

  useEffect(() => {
    const url = process.env.NEXT_PUBLIC_API_BASE_URL || "";
    console.log("🔧 API Base URL:", url || "(empty - using relative path)");
    setBaseUrl(url);
    fetchMetrics(url);
    fetchNews(url, selectedTopic);
    
    // Auto-refresh every 30s
    const interval = setInterval(() => {
      fetchMetrics(url);
      fetchNews(url, selectedTopic);
    }, 30000);
    
    return () => clearInterval(interval);
  }, [selectedTopic]);

  const fetchMetrics = async (url: string) => {
    try {
      const res = await fetch(`${url}/metrics`);
      if (res.ok) {
        const data = await res.json();
        setMetrics(data);
      }
    } catch (e) {
      console.error("Failed to fetch metrics", e);
    } finally {
      setLoading(false);
    }
  };

  const fetchNews = async (url: string, topic: string) => {
    setNewsLoading(true);
    setError(null);
    
    try {
      const topicParam = topic !== "all" ? `?topic=${topic}` : "";
      const res = await fetch(`${url}/api/news${topicParam}`);
      
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`);
      }
      
      const data: NewsResponse = await res.json();
      
      if (data.success && data.data) {
        setNews(data.data.items || []);
      } else {
        throw new Error(data.error || "Failed to load news");
      }
    } catch (e) {
      console.error("Failed to fetch news", e);
      setError(e instanceof Error ? e.message : "Failed to load news");
      setNews([]);
    } finally {
      setNewsLoading(false);
    }
  };

  const handleTrigger = async () => {
    setTriggering(true);
    setToast("Crawling...");
    try {
      const res = await fetch(`${baseUrl}/v2/crawl/trigger`, { method: "POST" });
      if (res.ok) {
        const data = await res.json();
        setLastTrigger(new Date());
        setToast(`Success! Found: ${data.stats?.discovered || 0}`);
        fetchMetrics(baseUrl);
        fetchNews(baseUrl, selectedTopic);
      } else {
        setToast("Failed to trigger cycle.");
      }
    } catch (e) {
      setToast("Network error.");
    } finally {
      setTriggering(false);
      setTimeout(() => setToast(null), 3000);
    }
  };

  const handleRetry = () => {
    fetchNews(baseUrl, selectedTopic);
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-[#050A0F] via-[#0A1118] to-[#050A0F]">
      {/* Header */}
      <header className="border-b border-gray-800/50 bg-black/20 backdrop-blur-xl sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-lg">
                <Sparkles className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
                  Evolipia Radar
                </h1>
                <p className="text-xs text-gray-500">AI Research Intelligence</p>
              </div>
            </div>
            
            <div className="flex items-center gap-3">
              {toast && (
                <span className="text-sm px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-emerald-400 animate-pulse">
                  {toast}
                </span>
              )}
              <button
                onClick={handleTrigger}
                disabled={triggering}
                className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-emerald-500 to-teal-600 text-white font-medium rounded-lg hover:from-emerald-600 hover:to-teal-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-emerald-500/20"
              >
                <RefreshCw className={`w-4 h-4 ${triggering ? 'animate-spin' : ''}`} />
                <span className="hidden sm:inline">Run Crawl</span>
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Metrics Cards */}
        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-900/40 rounded-xl animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <MetricCard
              icon={<FileText className="text-blue-400 w-5 h-5" />}
              label="Sources Processed"
              value={metrics?.articles_processed || 0}
              gradient="from-blue-500/10 to-blue-600/5"
            />
            <MetricCard
              icon={<Database className="text-rose-400 w-5 h-5" />}
              label="Filtered (Noise)"
              value={metrics?.filtered_articles || 0}
              gradient="from-rose-500/10 to-rose-600/5"
            />
            <MetricCard
              icon={<BrainCircuit className="text-emerald-400 w-5 h-5" />}
              label="Active Clusters"
              value={metrics?.clusters || 0}
              gradient="from-emerald-500/10 to-emerald-600/5"
            />
            <MetricCard
              icon={<Activity className="text-amber-400 w-5 h-5" />}
              label="Avg Cluster Score"
              value={(metrics?.avg_cluster_score || 0).toFixed(1)}
              gradient="from-amber-500/10 to-amber-600/5"
            />
          </div>
        )}

        {/* Topic Filter Bar */}
        <div className="mb-6">
          <div className="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-hide">
            {TOPICS.map((topic) => (
              <button
                key={topic.id}
                onClick={() => setSelectedTopic(topic.id)}
                className={`
                  px-4 py-2 rounded-lg font-medium text-sm whitespace-nowrap transition-all
                  ${selectedTopic === topic.id
                    ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/30'
                    : 'bg-gray-900/40 text-gray-400 hover:bg-gray-800/60 hover:text-gray-300 border border-gray-800/50'
                  }
                `}
              >
                {topic.label}
              </button>
            ))}
          </div>
        </div>

        {/* Latest News Section */}
        <section className="mb-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <TrendingUp className="w-6 h-6 text-emerald-400" />
              Latest News
              {!newsLoading && news.length > 0 && (
                <span className="text-sm font-normal text-gray-500 ml-2">
                  ({news.length} articles)
                </span>
              )}
            </h2>
            {lastTrigger && (
              <span className="text-xs text-gray-500 font-mono">
                Last sync: {lastTrigger.toLocaleTimeString()}
              </span>
            )}
          </div>

          {/* Loading State */}
          {newsLoading && (
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <LoadingSkeleton key={i} />
              ))}
            </div>
          )}

          {/* Error State */}
          {!newsLoading && error && (
            <ErrorState error={error} onRetry={handleRetry} />
          )}

          {/* Empty State */}
          {!newsLoading && !error && news.length === 0 && (
            <EmptyState topic={selectedTopic} />
          )}

          {/* News List */}
          {!newsLoading && !error && news.length > 0 && (
            <div className="space-y-4">
              {news.map((item, index) => (
                <NewsCard key={item.id} item={item} index={index} />
              ))}
            </div>
          )}
        </section>

        {/* Activity Log */}
        <section className="bg-gray-900/30 border border-gray-800/50 rounded-xl p-6 backdrop-blur-sm">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <Activity className="w-5 h-5 text-emerald-400" />
            Activity Log
          </h2>

          {(metrics?.clusters || 0) === 0 ? (
            <div className="text-center py-12 px-4 border border-dashed border-gray-800 rounded-lg bg-black/20">
              <BrainCircuit className="w-12 h-12 text-gray-700 mx-auto mb-3" />
              <h3 className="text-base font-medium text-gray-400">Semantic Engine Idle</h3>
              <p className="text-sm text-gray-600 mt-2 max-w-md mx-auto">
                AI clustering pending Phase 5 launch. Trigger a crawl to populate the database.
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {metrics?.top_cluster_titles?.map((title, i) => (
                <div 
                  key={i} 
                  className="p-3 bg-black/40 rounded-lg border border-gray-800/50 flex items-center gap-3 hover:border-gray-700 transition-colors"
                >
                  <span className="flex-shrink-0 w-7 h-7 rounded-full bg-emerald-500/10 text-emerald-400 flex items-center justify-center font-bold text-xs">
                    {i + 1}
                  </span>
                  <span className="text-sm text-gray-300">{title}</span>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Footer */}
        <footer className="mt-12 text-center text-xs text-gray-600 font-mono">
          Phase 4.5 • Powered by Vercel Serverless
        </footer>
      </div>
    </main>
  );
}

function MetricCard({ 
  icon, 
  label, 
  value,
  gradient 
}: { 
  icon: React.ReactNode; 
  label: string; 
  value: string | number;
  gradient: string;
}) {
  return (
    <div className={`bg-gradient-to-br ${gradient} border border-gray-800/50 p-5 rounded-xl backdrop-blur-sm hover:border-gray-700/50 transition-all group`}>
      <div className="flex items-center justify-between mb-3">
        <div className="p-2 bg-black/40 rounded-lg border border-gray-800/50 group-hover:border-gray-700/50 transition-colors">
          {icon}
        </div>
      </div>
      <p className="text-3xl font-bold tracking-tight mb-1">{value}</p>
      <h3 className="text-xs font-medium text-gray-500 uppercase tracking-wide">{label}</h3>
    </div>
  );
}

// News Card Component
function NewsCard({ item, index }: { item: NewsItem; index: number }) {
  const getTimeAgo = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
    
    if (seconds < 60) return `${seconds}s ago`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    return `${Math.floor(seconds / 86400)}d ago`;
  };

  const getTagColor = (tag: string) => {
    const colors: Record<string, string> = {
      llm: "bg-purple-500/10 text-purple-400 border-purple-500/20",
      vision: "bg-blue-500/10 text-blue-400 border-blue-500/20",
      data: "bg-green-500/10 text-green-400 border-green-500/20",
      security: "bg-red-500/10 text-red-400 border-red-500/20",
      rl: "bg-yellow-500/10 text-yellow-400 border-yellow-500/20",
      robotics: "bg-orange-500/10 text-orange-400 border-orange-500/20",
      ide: "bg-cyan-500/10 text-cyan-400 border-cyan-500/20",
      "free-credits": "bg-pink-500/10 text-pink-400 border-pink-500/20",
      general_ai: "bg-gray-500/10 text-gray-400 border-gray-500/20",
      research: "bg-indigo-500/10 text-indigo-400 border-indigo-500/20",
      tools: "bg-teal-500/10 text-teal-400 border-teal-500/20",
    };
    return colors[tag.toLowerCase()] || "bg-gray-500/10 text-gray-400 border-gray-500/20";
  };

  return (
    <article className="group bg-gray-900/30 border border-gray-800/50 rounded-xl p-5 hover:bg-gray-900/50 hover:border-gray-700/50 transition-all backdrop-blur-sm">
      <div className="flex items-start gap-4">
        {/* Rank Badge */}
        <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-500/10 to-teal-600/10 border border-emerald-500/20 flex items-center justify-center">
          <span className="text-sm font-bold text-emerald-400">#{index + 1}</span>
        </div>

        <div className="flex-1 min-w-0">
          {/* Title */}
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="block group/link"
          >
            <h3 className="text-lg font-semibold text-gray-100 group-hover/link:text-emerald-400 transition-colors line-clamp-2 mb-2">
              {item.title}
              <ExternalLink className="inline-block w-4 h-4 ml-1 opacity-0 group-hover/link:opacity-100 transition-opacity" />
            </h3>
          </a>

          {/* Summary */}
          {(item.tldr || item.why_it_matters) && (
            <p className="text-sm text-gray-400 line-clamp-2 mb-3">
              {item.tldr || item.why_it_matters}
            </p>
          )}

          {/* Meta Info */}
          <div className="flex flex-wrap items-center gap-3 text-xs text-gray-500">
            {/* Domain */}
            <span className="flex items-center gap-1">
              <FileText className="w-3 h-3" />
              {item.domain}
            </span>

            {/* Time */}
            <span className="flex items-center gap-1">
              <Clock className="w-3 h-3" />
              {getTimeAgo(item.published_at)}
            </span>

            {/* Score */}
            {item.score > 0 && (
              <span className="flex items-center gap-1 text-emerald-400">
                <TrendingUp className="w-3 h-3" />
                {item.score.toFixed(2)}
              </span>
            )}
          </div>

          {/* Tags */}
          {item.tags && item.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-3">
              {item.tags.slice(0, 5).map((tag, i) => (
                <span
                  key={i}
                  className={`px-2 py-1 rounded-md text-xs font-medium border ${getTagColor(tag)}`}
                >
                  {tag}
                </span>
              ))}
              {item.tags.length > 5 && (
                <span className="px-2 py-1 rounded-md text-xs font-medium bg-gray-800/50 text-gray-500 border border-gray-700/50">
                  +{item.tags.length - 5}
                </span>
              )}
            </div>
          )}
        </div>
      </div>
    </article>
  );
}

// Loading Skeleton Component
function LoadingSkeleton() {
  return (
    <div className="bg-gray-900/30 border border-gray-800/50 rounded-xl p-5 animate-pulse">
      <div className="flex items-start gap-4">
        <div className="w-10 h-10 bg-gray-800 rounded-lg" />
        <div className="flex-1 space-y-3">
          <div className="h-6 bg-gray-800 rounded w-3/4" />
          <div className="h-4 bg-gray-800 rounded w-full" />
          <div className="h-4 bg-gray-800 rounded w-5/6" />
          <div className="flex gap-2">
            <div className="h-6 w-16 bg-gray-800 rounded" />
            <div className="h-6 w-16 bg-gray-800 rounded" />
            <div className="h-6 w-16 bg-gray-800 rounded" />
          </div>
        </div>
      </div>
    </div>
  );
}

// Empty State Component
function EmptyState({ topic }: { topic: string }) {
  return (
    <div className="text-center py-16 px-4 border border-dashed border-gray-800 rounded-xl bg-black/20">
      <FileText className="w-16 h-16 text-gray-700 mx-auto mb-4" />
      <h3 className="text-xl font-semibold text-gray-300 mb-2">No News Found</h3>
      <p className="text-gray-500 max-w-md mx-auto">
        {topic === "all" 
          ? "No articles available yet. Try running a crawl to fetch the latest news."
          : `No articles found for "${topic}". Try selecting a different topic or run a new crawl.`
        }
      </p>
    </div>
  );
}

// Error State Component
function ErrorState({ error, onRetry }: { error: string; onRetry: () => void }) {
  return (
    <div className="text-center py-16 px-4 border border-red-900/20 rounded-xl bg-red-950/10">
      <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
      <h3 className="text-xl font-semibold text-red-400 mb-2">Failed to Load News</h3>
      <p className="text-gray-400 max-w-md mx-auto mb-6">{error}</p>
      <button
        onClick={onRetry}
        className="px-6 py-2.5 bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/20 rounded-lg font-medium transition-colors"
      >
        Try Again
      </button>
    </div>
  );
}
