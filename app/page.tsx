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
  Sparkles,
  Settings,
  X,
  Shield,
  Key,
  MessageSquare,
  Flame,
  ArrowUpDown,
  Bot,
  Globe,
  Server,
  Lightbulb
} from "lucide-react";

// API Response Types
interface NewsItem {
  id: string;
  title: string;
  url: string;
  domain: string;
  published_at: string;
  category: string;
  score: number;       // 1-10 scale (already converted by backend)
  raw_score: number;   // 0.0-1.0 internal
  heat_level: string;  // "hot" | "rising" | "signal" | "low"
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

interface SettingsState {
  x_api_key: string;
  threads_api_key: string;
  openrouter_api_key: string;
}

// Topic filter configuration — aligned with real AI verticals
const TOPICS = [
  { id: "all", label: "All", icon: <TrendingUp className="w-4 h-4" /> },
  { id: "llm", label: "LLM", icon: <BrainCircuit className="w-4 h-4" /> },
  { id: "agents", label: "Agents", icon: <Bot className="w-4 h-4" /> },
  { id: "vision", label: "Vision", icon: <Zap className="w-4 h-4" /> },
  { id: "open-source", label: "Open Source", icon: <Globe className="w-4 h-4" /> },
  { id: "infra", label: "Infra", icon: <Server className="w-4 h-4" /> },
  { id: "robotics", label: "Robotics", icon: <Activity className="w-4 h-4" /> },
  { id: "security", label: "Security", icon: <Shield className="w-4 h-4" /> },
] as const;

const SORT_OPTIONS = [
  { id: "", label: "Trending" },
  { id: "newest", label: "Newest" },
  { id: "oldest", label: "Oldest" },
] as const;

export default function Dashboard() {
  const [news, setNews] = useState<NewsItem[]>([]);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [newsLoading, setNewsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTopic, setSelectedTopic] = useState<string>("all");
  const [sortMode, setSortMode] = useState<string>("");
  const [triggering, setTriggering] = useState(false);
  const [toast, setToast] = useState<{message: string, type: 'info' | 'success' | 'error'} | null>(null);
  const [showSettings, setShowSettings] = useState(false);
  const [settings, setSettings] = useState<SettingsState>({
    x_api_key: "",
    threads_api_key: "",
    openrouter_api_key: "",
  });

  const [baseUrl, setBaseUrl] = useState("");

  const buildUrl = (base: string, path: string) => {
    const b = (base || "").replace(/\/+$/g, "");
    const p = path.startsWith("/") ? path : `/${path}`;
    return b ? `${b}${p}` : p;
  };

  useEffect(() => {
    const url = (process.env.NEXT_PUBLIC_API_BASE_URL || "").replace(/\/+$/g, "");
    setBaseUrl(url);
    fetchMetrics(url);
    fetchNews(url, selectedTopic, sortMode);
    fetchSettings(url);
    
    const interval = setInterval(() => {
      fetchMetrics(url);
    }, 15000);
    
    return () => clearInterval(interval);
  }, [selectedTopic, sortMode]);

  const fetchMetrics = async (url: string) => {
    try {
      const res = await fetch(buildUrl(url, "/metrics"));
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

  const fetchNews = async (url: string, topic: string, sort: string) => {
    setNewsLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams();
      if (topic !== "all") params.set("topic", topic);
      if (sort) params.set("sort", sort);
      const qs = params.toString() ? `?${params.toString()}` : "";
      
      const res = await fetch(`${buildUrl(url, "/api/news")}${qs}`);
      if (!res.ok) throw new Error("Failed to fetch news");
      const data: NewsResponse = await res.json();
      if (data.success && data.data) {
        setNews(data.data.items || []);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load news");
    } finally {
      setNewsLoading(false);
    }
  };

  const fetchSettings = async (url: string) => {
    try {
      const res = await fetch(buildUrl(url, "/v1/settings"));
      if (res.ok) {
        const data = await res.json();
        setSettings({
          x_api_key: data.x_api_key || "",
          threads_api_key: data.threads_api_key || "",
          openrouter_api_key: data.openrouter_api_key || "",
        });
      }
    } catch (e) {
      console.error("Failed to fetch settings", e);
    }
  };

  const saveSettings = async () => {
    try {
      const res = await fetch(buildUrl(baseUrl, "/v1/settings"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(settings),
      });
      if (res.ok) {
        setToast({ message: "Settings saved!", type: 'success' });
        setTimeout(() => setToast(null), 3000);
        setShowSettings(false);
      }
    } catch (e) {
      setToast({ message: "Failed to save settings.", type: 'error' });
    }
  };

  const handleTrigger = async () => {
    setTriggering(true);
    setToast({ message: "Evoli is starting a new crawl cycle...", type: 'info' });
    try {
      const res = await fetch(buildUrl(baseUrl, "/v2/crawl/trigger"), { method: "POST" });
      if (res.ok) {
        setToast({ message: "Crawl complete! Insights are being clustered.", type: 'success' });
        fetchMetrics(baseUrl);
        fetchNews(baseUrl, selectedTopic, sortMode);
      } else {
        setToast({ message: "Failed to trigger crawl.", type: 'error' });
      }
    } catch (e) {
      setToast({ message: "Network error.", type: 'error' });
    } finally {
      setTriggering(false);
      setTimeout(() => setToast(null), 5000);
    }
  };

  return (
    <main className="min-h-screen bg-[#050A0F] text-slate-200 font-sans selection:bg-emerald-500/30">
      {/* Dynamic Background */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-emerald-500/10 blur-[120px] rounded-full" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-blue-500/10 blur-[120px] rounded-full" />
      </div>

      {/* Header */}
      <header className="sticky top-0 z-[60] border-b border-white/5 bg-black/40 backdrop-blur-2xl">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-3 sm:py-4 flex items-center justify-between">
          <div className="flex items-center gap-3 sm:gap-4">
            <div className="relative group">
              <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-teal-500 rounded-xl blur opacity-25 group-hover:opacity-50 transition duration-1000 group-hover:duration-200" />
              <img 
                src="/assets/icon.png" 
                alt="Logo" 
                className="relative w-9 h-9 sm:w-11 sm:h-11 rounded-xl shadow-2xl transition-transform hover:scale-105"
              />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-bold tracking-tight text-white flex items-center gap-2">
                Evolipia Radar
                <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-mono border border-emerald-500/20 uppercase tracking-widest">
                  Active
                </span>
              </h1>
              <p className="text-xs text-slate-500 font-medium hidden sm:block">Autonomous Research Engine</p>
            </div>
          </div>

          <div className="flex items-center gap-2 sm:gap-3">
            <button
              onClick={() => setShowSettings(true)}
              className="p-2 sm:p-2.5 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 hover:border-white/20 transition-all text-slate-400 hover:text-white"
              title="System Settings"
            >
              <Settings className="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
            <button
              onClick={handleTrigger}
              disabled={triggering}
              className="relative group overflow-hidden flex items-center gap-2 px-4 sm:px-6 py-2 sm:py-2.5 bg-white text-black font-bold rounded-xl transition-all active:scale-95 disabled:opacity-50 text-sm sm:text-base"
            >
              <RefreshCw className={`w-4 h-4 ${triggering ? 'animate-spin' : ''}`} />
              <span className="hidden sm:inline">{triggering ? 'Processing...' : 'Run Cycle'}</span>
              <span className="sm:hidden">{triggering ? '...' : 'Run'}</span>
              <div className="absolute inset-0 bg-gradient-to-r from-emerald-400/20 to-teal-400/20 opacity-0 group-hover:opacity-100 transition-opacity" />
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section with Mascot */}
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 py-8 sm:py-10">
        <div className="flex flex-col lg:flex-row items-center gap-8 lg:gap-16">
          <div className="flex-1 space-y-5 text-center lg:text-left">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-semibold">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Agent Evoli monitoring 12 sources</span>
            </div>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-white leading-[1.1] tracking-tight">
              Predict the future of <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 via-teal-400 to-blue-500">
                AI Innovation.
              </span>
            </h2>
            <p className="text-base sm:text-lg text-slate-400 max-w-2xl leading-relaxed">
              Real-time semantic clustering of global research signals. 
              Evolipia filters the noise to bring you what actually moves markets.
            </p>
          </div>
          
          <div className="relative lg:w-1/3 flex justify-center">
            <div className="absolute inset-0 bg-emerald-500/10 blur-[100px] rounded-full animate-pulse" />
            <div className="relative group">
              <div className="absolute -inset-1 bg-emerald-500/20 rounded-full blur opacity-25 group-hover:opacity-50 transition duration-1000" />
              <img 
                src="/assets/maskot1.png" 
                alt="Evoli Mascot" 
                className="w-40 h-40 lg:w-56 lg:h-56 rounded-full object-cover relative border-4 border-emerald-500/20 bg-emerald-500/5 shadow-[0_0_50px_rgba(16,185,129,0.2)] animate-float"
              />
              <div className="absolute bottom-3 left-1/2 -translate-x-1/2 px-4 py-1.5 bg-black/80 backdrop-blur-md rounded-full border border-emerald-500/20 text-[10px] font-black uppercase tracking-[0.3em] text-emerald-500 opacity-0 group-hover:opacity-100 transition-all">
                Agent Evoli
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 pb-20">
        {/* Metrics Grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 mb-10 sm:mb-12">
          <MetricCard 
            label="Sources Crawled" 
            value={metrics?.articles_processed || 0} 
            icon={<FileText className="w-5 h-5" />}
          />
          <MetricCard 
            label="AI Analyzed" 
            value={metrics?.filtered_articles || 0} 
            icon={<Shield className="w-5 h-5" />}
          />
          <MetricCard 
            label="Summaries" 
            value={metrics?.clusters || 0} 
            icon={<BrainCircuit className="w-5 h-5" />}
            highlight
          />
          <MetricCard 
            label="Avg Score" 
            value={metrics?.avg_cluster_score?.toFixed(1) || "0.0"} 
            icon={<TrendingUp className="w-5 h-5" />}
            suffix="/10"
          />
        </div>

        {/* Intelligence Feed */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Feed */}
          <div className="lg:col-span-2 space-y-6">
            {/* Filter Bar + Sort */}
            <div className="space-y-3 border-b border-white/5 pb-4">
              <div className="flex items-center justify-between">
                <h3 className="text-xl sm:text-2xl font-bold text-white flex items-center gap-2">
                  <Sparkles className="w-5 h-5 sm:w-6 sm:h-6 text-emerald-400" />
                  Latest Insights
                </h3>
                
                {/* Sort Toggle */}
                <div className="flex items-center gap-1 bg-white/5 rounded-xl p-1 border border-white/5">
                  {SORT_OPTIONS.map(opt => (
                    <button
                      key={opt.id}
                      onClick={() => setSortMode(opt.id)}
                      className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all ${
                        sortMode === opt.id
                        ? 'bg-emerald-500 text-black' 
                        : 'text-slate-400 hover:text-white'
                      }`}
                    >
                      {opt.label}
                    </button>
                  ))}
                </div>
              </div>
              
              {/* Topic Filters — Scrollable */}
              <div className="flex items-center gap-2 overflow-x-auto scrollbar-hide pb-1">
                {TOPICS.map(topic => (
                  <button
                    key={topic.id}
                    onClick={() => setSelectedTopic(topic.id)}
                    className={`flex-shrink-0 px-4 py-2 rounded-xl text-sm font-semibold transition-all flex items-center gap-2 ${
                      selectedTopic === topic.id 
                      ? 'bg-emerald-500 text-black shadow-lg shadow-emerald-500/20' 
                      : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white border border-white/5'
                    }`}
                  >
                    {topic.icon}
                    <span>{topic.label}</span>
                  </button>
                ))}
              </div>
            </div>

            {newsLoading ? (
              <div className="space-y-4">
                {[1, 2, 3].map(i => <SkeletonCard key={i} />)}
              </div>
            ) : error ? (
              <div className="p-8 sm:p-12 text-center bg-rose-500/5 border border-rose-500/20 rounded-2xl">
                <AlertCircle className="w-10 h-10 text-rose-500 mx-auto mb-3" />
                <h4 className="text-lg font-bold text-white mb-2">Sync Interrupted</h4>
                <p className="text-slate-400 mb-4 text-sm">{error}</p>
                <button 
                  onClick={() => fetchNews(baseUrl, selectedTopic, sortMode)}
                  className="px-6 py-2 bg-rose-500 text-white font-bold rounded-xl hover:bg-rose-600 transition-colors text-sm"
                >
                  Reconnect
                </button>
              </div>
            ) : news.length === 0 ? (
              <div className="py-16 text-center border-2 border-dashed border-white/5 rounded-2xl">
                <div className="w-16 h-16 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Database className="w-8 h-8 text-slate-700" />
                </div>
                <h4 className="text-lg font-bold text-white mb-2">Feed Empty</h4>
                <p className="text-slate-500 max-w-sm mx-auto text-sm">
                  Run a crawl cycle to ingest fresh AI research signals.
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-4">
                {news.map((item, idx) => (
                  <NewsCard key={item.id} item={item} index={idx} />
                ))}
              </div>
            )}
          </div>

          {/* Sidebar - Emerging Trends */}
          <aside className="space-y-6">
            <div className="p-5 sm:p-6 bg-black/40 border border-white/5 rounded-2xl backdrop-blur-xl">
              <h3 className="text-base sm:text-lg font-bold text-white mb-5 flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-emerald-400" />
                Top Trending
              </h3>
              
              {!metrics?.top_cluster_titles || metrics.top_cluster_titles.length === 0 ? (
                <div className="py-10 text-center space-y-3">
                  <BrainCircuit className="w-8 h-8 text-slate-800 mx-auto" />
                  <p className="text-sm text-slate-600">No active trends yet.</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {metrics.top_cluster_titles.map((title, i) => (
                    <div key={i} className="flex gap-3 group cursor-pointer">
                      <span className="text-emerald-500/40 font-mono text-sm group-hover:text-emerald-400 transition-colors pt-0.5">{String(i+1).padStart(2, '0')}</span>
                      <div>
                        <h4 className="text-sm font-semibold text-slate-300 group-hover:text-white transition-colors leading-snug line-clamp-2">
                          {title}
                        </h4>
                        <div className="flex items-center gap-2 mt-1">
                          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                          <span className="text-[10px] text-slate-500 font-bold uppercase tracking-wider">Emerging</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="p-6 sm:p-8 bg-gradient-to-br from-emerald-600 to-teal-700 rounded-2xl text-white shadow-2xl shadow-emerald-500/10 relative overflow-hidden group">
              <div className="absolute top-[-20%] right-[-20%] w-[60%] h-[60%] bg-white/10 blur-3xl rounded-full transition-transform group-hover:scale-125 duration-700" />
              <div className="relative z-10 space-y-3">
                <h4 className="text-lg sm:text-xl font-black">Join 5,000+ Researchers</h4>
                <p className="text-sm text-white/80 font-medium leading-relaxed">
                  Daily intelligence alerts for emerging LLM and CV breakthroughs.
                </p>
                <button className="w-full py-3 bg-white text-emerald-700 font-extrabold rounded-xl shadow-xl hover:shadow-2xl transition-all active:scale-95 text-sm">
                  Subscribe for Free
                </button>
              </div>
            </div>
          </aside>
        </div>
      </div>

      {/* Settings Modal */}
      {showSettings && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-12">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setShowSettings(false)} />
          <div className="relative w-full max-w-xl bg-[#0A1118] border border-white/10 rounded-2xl shadow-2xl p-6 sm:p-10 overflow-y-auto max-h-[90vh]">
            <div className="absolute top-0 right-0 p-4">
              <button 
                onClick={() => setShowSettings(false)}
                className="p-2 rounded-xl bg-white/5 text-slate-400 hover:text-white transition-colors"
                title="Close"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="flex items-center gap-3 mb-8">
              <div className="p-3 bg-emerald-500/10 border border-emerald-500/20 rounded-xl">
                <Settings className="w-6 h-6 text-emerald-400" />
              </div>
              <div>
                <h3 className="text-2xl font-black text-white">System Core</h3>
                <p className="text-slate-500 font-medium text-sm">Manage your intelligence keys</p>
              </div>
            </div>

            <div className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <Key className="w-4 h-4" /> OpenRouter API Key
                </label>
                <input 
                  type="password"
                  value={settings.openrouter_api_key}
                  onChange={(e) => setSettings({...settings, openrouter_api_key: e.target.value})}
                  placeholder="sk-or-v1-..."
                  className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
                <p className="text-[10px] text-slate-500 ml-1">Required for AI clustering and summarization.</p>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <MessageSquare className="w-4 h-4" /> X (Twitter) API Key
                </label>
                <input 
                  type="password"
                  value={settings.x_api_key}
                  onChange={(e) => setSettings({...settings, x_api_key: e.target.value})}
                  placeholder="Enter X API Key..."
                  className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <Activity className="w-4 h-4" /> Threads API Key
                </label>
                <input 
                  type="password"
                  value={settings.threads_api_key}
                  onChange={(e) => setSettings({...settings, threads_api_key: e.target.value})}
                  placeholder="Enter Threads API Key..."
                  className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
              </div>

              <div className="pt-4">
                <button 
                  onClick={saveSettings}
                  className="w-full py-4 bg-white text-black font-black text-base rounded-xl hover:bg-emerald-400 transition-all shadow-xl active:scale-95"
                >
                  Authorize System Keys
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Toast Notification */}
      {toast && (
        <div className="fixed bottom-6 right-6 z-[200] animate-in">
          <div className={`
            px-5 py-3 rounded-xl shadow-2xl border backdrop-blur-xl flex items-center gap-3 text-sm
            ${toast.type === 'success' ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400' : 
              toast.type === 'error' ? 'bg-rose-500/10 border-rose-500/20 text-rose-400' :
              'bg-blue-500/10 border-blue-500/20 text-blue-400'}
          `}>
             {toast.type === 'success' ? <Sparkles className="w-4 h-4" /> : 
              toast.type === 'error' ? <AlertCircle className="w-4 h-4" /> : 
              <RefreshCw className="w-4 h-4 animate-spin" />}
             <span className="font-bold">{toast.message}</span>
          </div>
        </div>
      )}

      <style jsx global>{`
        @keyframes float {
          0%, 100% { transform: translateY(0px); }
          50% { transform: translateY(-20px); }
        }
        .animate-float {
          animation: float 6s ease-in-out infinite;
        }
      `}</style>
    </main>
  );
}

// ============================================================================
// Components
// ============================================================================

function MetricCard({ label, value, icon, highlight, suffix }: {
  label: string; value: number | string; icon: React.ReactNode; highlight?: boolean; suffix?: string;
}) {
  return (
    <div className={`p-4 sm:p-6 rounded-2xl bg-gradient-to-br border transition-all hover:scale-[1.02] cursor-default group ${
      highlight 
      ? 'from-emerald-500/10 to-emerald-500/[0.02] border-emerald-500/20 text-emerald-400' 
      : 'from-white/5 to-white/[0.02] border-white/5 text-slate-400 group-hover:border-emerald-500/20'
    }`}>
      <div className="flex items-start justify-between mb-4 sm:mb-6">
        <div className="p-2 sm:p-3 bg-black/40 border border-white/5 rounded-xl group-hover:border-white/20 transition-all">
          {icon}
        </div>
      </div>
      <div className="space-y-0.5">
        <p className="text-2xl sm:text-3xl font-black text-white tracking-tighter">
          {value}{suffix && <span className="text-sm sm:text-base font-bold text-slate-500">{suffix}</span>}
        </p>
        <p className="text-[10px] sm:text-xs font-bold text-slate-500 uppercase tracking-widest">{label}</p>
      </div>
    </div>
  );
}

function HeatBadge({ level, score }: { level: string; score: number }) {
  const config: Record<string, { icon: React.ReactNode; bg: string; text: string; label: string }> = {
    hot:    { icon: <Flame className="w-3 h-3" />, bg: "bg-orange-500/10 border-orange-500/20", text: "text-orange-400", label: "Hot" },
    rising: { icon: <TrendingUp className="w-3 h-3" />, bg: "bg-amber-500/10 border-amber-500/20", text: "text-amber-400", label: "Rising" },
    signal: { icon: <Lightbulb className="w-3 h-3" />, bg: "bg-blue-500/10 border-blue-500/20", text: "text-blue-400", label: "Signal" },
    low:    { icon: <Activity className="w-3 h-3" />, bg: "bg-slate-500/10 border-slate-500/20", text: "text-slate-400", label: "New" },
  };
  const c = config[level] || config.low;

  return (
    <div className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border text-[11px] font-bold uppercase tracking-wider ${c.bg} ${c.text}`}>
      {c.icon}
      <span>{c.label}</span>
      <span className="opacity-60">{score.toFixed(1)}</span>
    </div>
  );
}

function NewsCard({ item, index }: { item: NewsItem, index: number }) {
  const [isExpanded, setIsExpanded] = useState(false);
  
  // Ensure score is on 1-10 scale (backend now sends it correctly)
  const displayScore = item.score >= 0 && item.score <= 1.1 
    ? (item.score * 9) + 1  // Fallback conversion if backend sent raw
    : item.score;           // Already 1-10 scale
  
  const heatLevel = item.heat_level || (displayScore >= 7 ? "hot" : displayScore >= 5 ? "rising" : displayScore >= 3 ? "signal" : "low");

  return (
    <div className="group relative">
      <div className="absolute -inset-[1px] bg-gradient-to-r from-emerald-500/50 to-blue-500/50 rounded-2xl blur-sm opacity-0 group-hover:opacity-10 transition-opacity" />
      <div className="relative p-4 sm:p-6 bg-black/40 border border-white/5 rounded-2xl hover:border-white/10 transition-all backdrop-blur-xl">
        <div className="flex items-start gap-4">
          <div className="hidden sm:flex flex-shrink-0 w-10 h-10 rounded-xl bg-white/5 border border-white/10 items-center justify-center font-black text-base text-slate-600 group-hover:text-emerald-500/40 group-hover:border-emerald-500/20 transition-all">
            {index + 1}
          </div>
          
          <div className="flex-1 space-y-3">
            <div className="flex flex-wrap items-center gap-2">
              <HeatBadge level={heatLevel} score={displayScore} />
              <span className="text-[11px] font-bold text-slate-500 flex items-center gap-1 uppercase tracking-widest">
                <Clock className="w-3 h-3" /> {new Date(item.published_at).toLocaleDateString()}
              </span>
              <span className="text-[11px] text-slate-600">{item.domain}</span>
            </div>

            <h4 className="text-lg sm:text-xl font-black text-white leading-tight group-hover:text-emerald-400 transition-colors">
              {item.title}
            </h4>

            {item.tldr && (
              <div className="space-y-2">
                <div className={`text-slate-400 text-sm leading-relaxed ${!isExpanded && 'line-clamp-2'}`}>
                  {item.tldr}
                </div>
                {item.why_it_matters && isExpanded && (
                  <div className="pt-3 border-t border-white/5">
                    <h5 className="text-xs font-black uppercase tracking-[0.2em] text-blue-400 mb-1.5">Why It Matters</h5>
                    <p className="text-slate-400 text-sm leading-relaxed">{item.why_it_matters}</p>
                  </div>
                )}
                <button 
                  onClick={() => setIsExpanded(!isExpanded)}
                  className="text-xs font-black text-emerald-500 uppercase tracking-widest hover:text-emerald-400 transition-colors"
                >
                  {isExpanded ? 'Collapse' : 'Expand Insight'}
                </button>
              </div>
            )}

            <div className="flex items-center justify-between pt-3">
              <div className="flex flex-wrap gap-1.5">
                {item.tags?.slice(0, 3).map(tag => (
                  <span key={tag} className="px-2.5 py-0.5 bg-white/5 border border-white/10 rounded-lg text-[10px] font-bold uppercase tracking-widest text-slate-400 group-hover:border-emerald-500/20 group-hover:text-emerald-500/80 transition-all">
                    {tag}
                  </span>
                ))}
              </div>
              
              <a 
                href={item.url} 
                target="_blank" 
                className="flex-shrink-0 flex items-center gap-1.5 px-4 py-1.5 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 font-bold text-xs rounded-lg hover:bg-emerald-500 hover:text-black transition-all"
              >
                Source
                <ExternalLink className="w-3.5 h-3.5" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function SkeletonCard() {
  return (
    <div className="p-4 sm:p-6 bg-white/5 border border-white/10 rounded-2xl animate-pulse">
      <div className="flex gap-4">
        <div className="hidden sm:block w-10 h-10 bg-white/5 rounded-xl" />
        <div className="flex-1 space-y-3">
          <div className="w-24 h-4 bg-white/5 rounded" />
          <div className="w-full h-6 bg-white/5 rounded" />
          <div className="w-full h-14 bg-white/5 rounded" />
        </div>
      </div>
    </div>
  );
}
