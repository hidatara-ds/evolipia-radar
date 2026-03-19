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
  MessageSquare
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

interface SettingsState {
  x_api_key: string;
  threads_api_key: string;
  openrouter_api_key: string;
}

// Topic filter configuration
const TOPICS = [
  { id: "all", label: "All Insights", icon: <TrendingUp className="w-4 h-4" /> },
  { id: "llm", label: "LLM", icon: <BrainCircuit className="w-4 h-4" /> },
  { id: "vision", label: "Vision", icon: <Zap className="w-4 h-4" /> },
  { id: "data", label: "Data", icon: <Database className="w-4 h-4" /> },
  { id: "robotics", label: "Robotics", icon: <Zap className="w-4 h-4" /> },
  { id: "credits", label: "Free Credits", icon: <Sparkles className="w-4 h-4" /> },
  { id: "ide", label: "IDE", icon: <FileText className="w-4 h-4" /> },
  { id: "security", label: "Security", icon: <Shield className="w-4 h-4" /> },
] as const;

export default function Dashboard() {
  const [news, setNews] = useState<NewsItem[]>([]);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [newsLoading, setNewsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTopic, setSelectedTopic] = useState<string>("all");
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
    fetchNews(url, selectedTopic);
    fetchSettings(url);
    
    const interval = setInterval(() => {
      fetchMetrics(url);
    }, 15000);
    
    return () => clearInterval(interval);
  }, [selectedTopic]);

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

  const fetchNews = async (url: string, topic: string) => {
    setNewsLoading(true);
    setError(null);
    try {
      const topicParam = topic !== "all" ? `?topic=${topic}` : "";
      const res = await fetch(`${buildUrl(url, "/api/news")}${topicParam}`);
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
        fetchNews(baseUrl, selectedTopic);
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
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="relative group">
              <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-teal-500 rounded-xl blur opacity-25 group-hover:opacity-50 transition duration-1000 group-hover:duration-200" />
              <img 
                src="/assets/icon.png" 
                alt="Logo" 
                className="relative w-11 h-11 rounded-xl shadow-2xl transition-transform hover:scale-105"
              />
            </div>
            <div>
              <h1 className="text-xl font-bold tracking-tight text-white flex items-center gap-2">
                Evolipia Radar
                <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-mono border border-emerald-500/20 uppercase tracking-widest">
                  Active
                </span>
              </h1>
              <p className="text-xs text-slate-500 font-medium">Autonomous Research Engine</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowSettings(true)}
              className="p-2.5 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 hover:border-white/20 transition-all text-slate-400 hover:text-white"
              title="System Settings"
            >
              <Settings className="w-5 h-5" />
            </button>
            <button
              onClick={handleTrigger}
              disabled={triggering}
              className="relative group overflow-hidden flex items-center gap-2.5 px-6 py-2.5 bg-white text-black font-bold rounded-xl transition-all active:scale-95 disabled:opacity-50"
            >
              <RefreshCw className={`w-4 h-4 ${triggering ? 'animate-spin' : ''}`} />
              <span>{triggering ? 'Processing...' : 'Run Cycle'}</span>
              <div className="absolute inset-0 bg-gradient-to-r from-emerald-400/20 to-teal-400/20 opacity-0 group-hover:opacity-100 transition-opacity" />
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section with Mascot */}
      <div className="relative max-w-7xl mx-auto px-6 py-12">
        <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-24">
          <div className="flex-1 space-y-6 text-center lg:text-left">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-semibold mb-2">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Agent Evoli is currently monitoring 12 sources</span>
            </div>
            <h2 className="text-4xl lg:text-6xl font-black text-white leading-[1.1] tracking-tight">
              Predict the future of <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 via-teal-400 to-blue-500">
                AI Innovation.
              </span>
            </h2>
            <p className="text-lg text-slate-400 max-w-2xl leading-relaxed">
              Real-time semantic clustering of global research signals. 
              Evolipia filters the noise to bring you the signal that actually moves markets.
            </p>
          </div>
          
          <div className="relative lg:w-1/3 flex justify-center">
            <div className="absolute inset-0 bg-emerald-500/10 blur-[100px] rounded-full animate-pulse" />
            <div className="relative group">
              <div className="absolute -inset-1 bg-emerald-500/20 rounded-full blur opacity-25 group-hover:opacity-50 transition duration-1000" />
              <img 
                src="/assets/maskot1.png" 
                alt="Evoli Mascot" 
                className="w-48 h-48 lg:w-64 lg:h-64 rounded-full object-cover relative border-4 border-emerald-500/20 bg-emerald-500/5 shadow-[0_0_50px_rgba(16,185,129,0.2)] animate-float"
              />
              <div className="absolute bottom-4 left-1/2 -translate-x-1/2 px-4 py-1.5 bg-black/80 backdrop-blur-md rounded-full border border-emerald-500/20 text-[10px] font-black uppercase tracking-[0.3em] text-emerald-500 opacity-0 group-hover:opacity-100 transition-all">
                Agent Evoli
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 pb-24">
        {/* Metrics Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-20">
          <PremiumMetricCard 
            label="Sources Crawled" 
            value={metrics?.articles_processed || 0} 
            icon={<FileText className="w-5 h-5" />}
            color="standard" 
          />
          <PremiumMetricCard 
            label="Noise Suppressed" 
            value={metrics?.filtered_articles || 0} 
            icon={<Shield className="w-5 h-5" />}
            color="standard" 
          />
          <PremiumMetricCard 
            label="Intelligence Clusters" 
            value={metrics?.clusters || 0} 
            icon={<BrainCircuit className="w-5 h-5" />}
            color="active" 
            trend="+12% today"
          />
          <PremiumMetricCard 
            label="Global Impact Score" 
            value={(metrics?.avg_cluster_score || 0).toFixed(1)} 
            icon={<TrendingUp className="w-5 h-5" />}
            color="standard" 
          />
        </div>

        {/* Intelligence Feed */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
          {/* Main Feed */}
          <div className="lg:col-span-2 space-y-8">
            <div className="flex items-center justify-between border-b border-white/5 pb-4">
              <h3 className="text-2xl font-bold text-white flex items-center gap-3">
                <Sparkles className="w-6 h-6 text-emerald-400" />
                Latest Insights
              </h3>
              
              <div className="flex items-center gap-2">
                {TOPICS.map(topic => (
                  <button
                    key={topic.id}
                    onClick={() => setSelectedTopic(topic.id)}
                    className={`px-4 py-2 rounded-xl text-sm font-semibold transition-all flex items-center gap-2 ${
                      selectedTopic === topic.id 
                      ? 'bg-emerald-500 text-black shadow-lg shadow-emerald-500/20' 
                      : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white border border-white/5'
                    }`}
                  >
                    {topic.icon}
                    <span className="hidden sm:inline">{topic.label}</span>
                  </button>
                ))}
              </div>
            </div>

            {newsLoading ? (
              <div className="space-y-6">
                {[1, 2, 3].map(i => <PremiumSkeleton key={i} />)}
              </div>
            ) : error ? (
              <div className="p-12 text-center bg-rose-500/5 border border-rose-500/20 rounded-3xl">
                <AlertCircle className="w-12 h-12 text-rose-500 mx-auto mb-4" />
                <h4 className="text-xl font-bold text-white mb-2">Sync Interrupted</h4>
                <p className="text-slate-400 mb-6">{error}</p>
                <button 
                  onClick={() => fetchNews(baseUrl, selectedTopic)}
                  className="px-6 py-2 bg-rose-500 text-white font-bold rounded-xl hover:bg-rose-600 transition-colors"
                >
                  Reconnect
                </button>
              </div>
            ) : news.length === 0 ? (
              <div className="py-20 text-center border-2 border-dashed border-white/5 rounded-3xl">
                <div className="w-20 h-20 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6">
                  <Database className="w-10 h-10 text-slate-700" />
                </div>
                <h4 className="text-xl font-bold text-white mb-2">Feed Vacant</h4>
                <p className="text-slate-500 max-w-sm mx-auto">
                  Run a system cycle to ingest new research signals from the edge.
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-6">
                {news.map((item, idx) => (
                  <PremiumNewsCard key={item.id} item={item} index={idx} />
                ))}
              </div>
            )}
          </div>

          {/* Sidebar - Emerging Trends */}
          <aside className="space-y-8">
            <div className="p-6 bg-black/40 border border-white/5 rounded-3xl backdrop-blur-xl">
              <h3 className="text-lg font-bold text-white mb-6 flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-emerald-400" />
                Trending Clusters
              </h3>
              
              {!metrics?.top_cluster_titles ? (
                <div className="py-12 text-center space-y-4">
                  <BrainCircuit className="w-10 h-10 text-slate-800 mx-auto" />
                  <p className="text-sm text-slate-600">No active clusters detected.</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {metrics.top_cluster_titles.map((title, i) => (
                    <div key={i} className="flex gap-4 group cursor-pointer">
                      <span className="text-emerald-500/40 font-mono text-sm group-hover:text-emerald-400 transition-colors pt-1">0{i+1}</span>
                      <div>
                        <h4 className="text-sm font-semibold text-slate-300 group-hover:text-white transition-colors leading-snug line-clamp-2">
                          {title}
                        </h4>
                        <div className="flex items-center gap-2 mt-1.5">
                          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                          <span className="text-[10px] text-slate-500 font-bold uppercase tracking-wider">Emerging Signal</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="p-8 bg-gradient-to-br from-emerald-600 to-teal-700 rounded-3xl text-white shadow-2xl shadow-emerald-500/10 relative overflow-hidden group">
              <div className="absolute top-[-20%] right-[-20%] w-[60%] h-[60%] bg-white/10 blur-3xl rounded-full transition-transform group-hover:scale-125 duration-700" />
              <div className="relative z-10 space-y-4">
                <h4 className="text-xl font-black">Join over 5,000 Researchers</h4>
                <p className="text-sm text-white/80 font-medium leading-relaxed">
                  Get daily intelligence alerts for emerging LLM and CV breakthroughs.
                </p>
                <button className="w-full py-3 bg-white text-emerald-700 font-extrabold rounded-xl shadow-xl hover:shadow-2xl transition-all active:scale-95">
                  Subscribe for Free
                </button>
              </div>
            </div>
          </aside>
        </div>
      </div>

      {/* Settings Modal */}
      {showSettings && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-6 sm:p-12">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setShowSettings(false)} />
          <div className="relative w-full max-w-xl bg-[#0A1118] border border-white/10 rounded-[2.5rem] shadow-2xl p-8 sm:p-12 overflow-hidden overflow-y-auto max-h-[90vh]">
            <div className="absolute top-0 right-0 p-6">
              <button 
                onClick={() => setShowSettings(false)}
                className="p-2 rounded-xl bg-white/5 text-slate-400 hover:text-white transition-colors"
                title="Close"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="flex items-center gap-4 mb-10">
              <div className="p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl">
                <Settings className="w-8 h-8 text-emerald-400" />
              </div>
              <div>
                <h3 className="text-3xl font-black text-white">System Core</h3>
                <p className="text-slate-500 font-medium">Manage your secret intelligence keys</p>
              </div>
            </div>

            <div className="space-y-8">
              <div className="space-y-3">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <Key className="w-4 h-4" /> OpenRouter API Key
                </label>
                <input 
                  type="password"
                  value={settings.openrouter_api_key}
                  onChange={(e) => setSettings({...settings, openrouter_api_key: e.target.value})}
                  placeholder="sk-or-v1-..."
                  className="w-full bg-white/5 border border-white/10 rounded-2xl px-6 py-4 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
                <p className="text-[10px] text-slate-500 ml-1">Necessary for AI clustering and summarization.</p>
              </div>

              <div className="space-y-3">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <MessageSquare className="w-4 h-4" /> X (Twitter) API Key
                </label>
                <input 
                  type="password"
                  value={settings.x_api_key}
                  onChange={(e) => setSettings({...settings, x_api_key: e.target.value})}
                  placeholder="Enter X API Key..."
                  className="w-full bg-white/5 border border-white/10 rounded-2xl px-6 py-4 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
              </div>

              <div className="space-y-3">
                <label className="text-sm font-bold text-slate-400 ml-1 flex items-center gap-2">
                  <Activity className="w-4 h-4" /> Threads API Key
                </label>
                <input 
                  type="password"
                  value={settings.threads_api_key}
                  onChange={(e) => setSettings({...settings, threads_api_key: e.target.value})}
                  placeholder="Enter Threads API Key..."
                  className="w-full bg-white/5 border border-white/10 rounded-2xl px-6 py-4 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500/50 transition-all font-mono text-sm"
                />
              </div>

              <div className="pt-6">
                <button 
                  onClick={saveSettings}
                  className="w-full py-5 bg-white text-black font-black text-lg rounded-3xl hover:bg-emerald-400 transition-all shadow-xl active:scale-95"
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
        <div className="fixed bottom-10 right-10 z-[200] animate-in fade-in slide-in-from-bottom-5">
          <div className={`
            px-6 py-4 rounded-2xl shadow-2xl border backdrop-blur-xl flex items-center gap-4
            ${toast.type === 'success' ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400' : 
              toast.type === 'error' ? 'bg-rose-500/10 border-rose-500/20 text-rose-400' :
              'bg-blue-500/10 border-blue-500/20 text-blue-400'}
          `}>
             {toast.type === 'success' ? <Sparkles className="w-5 h-5" /> : 
              toast.type === 'error' ? <AlertCircle className="w-5 h-5" /> : 
              <RefreshCw className="w-5 h-5 animate-spin" />}
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

function PremiumMetricCard({ label, value, icon, color, trend }: any) {
  const colors: any = {
    standard: "from-white/5 to-white/[0.02] border-white/5 text-slate-400 group-hover:border-emerald-500/20",
    active: "from-emerald-500/10 to-emerald-500/[0.02] border-emerald-500/20 text-emerald-400 group-hover:border-emerald-500/40",
  };

  return (
    <div className={`p-6 rounded-[2rem] bg-gradient-to-br border border-white/5 transition-all hover:scale-[1.02] active:scale-95 cursor-default group ${colors[color]}`}>
      <div className="flex items-start justify-between mb-8">
        <div className={`p-3 bg-black/40 border border-white/5 rounded-2xl group-hover:border-white/20 transition-all`}>
          {icon}
        </div>
        {trend && <span className="text-[10px] font-black uppercase tracking-widest">{trend}</span>}
      </div>
      <div className="space-y-1">
        <p className="text-4xl font-black text-white tracking-tighter">{value}</p>
        <p className="text-xs font-bold text-slate-500 uppercase tracking-widest">{label}</p>
      </div>
    </div>
  );
}

function PremiumNewsCard({ item, index }: { item: NewsItem, index: number }) {
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div className="group relative">
      <div className="absolute -inset-[1px] bg-gradient-to-r from-emerald-500/50 to-blue-500/50 rounded-[2.5rem] blur-sm opacity-0 group-hover:opacity-10 transition-opacity" />
      <div className="relative p-8 bg-black/40 border border-white/5 rounded-[2.5rem] hover:border-white/10 transition-all backdrop-blur-xl">
        <div className="flex items-start gap-6">
          <div className="flex-shrink-0 w-14 h-14 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center font-black text-xl text-slate-600 group-hover:text-emerald-500/40 group-hover:border-emerald-500/20 transition-all">
            {index + 1}
          </div>
          
          <div className="flex-1 space-y-4">
            <div className="flex items-center gap-3">
              <span className="text-[11px] font-black uppercase tracking-[0.2em] text-emerald-500/80 bg-emerald-500/5 px-2 py-0.5 rounded border border-emerald-500/10">
                Validated Node
              </span>
              <span className="text-[11px] font-bold text-slate-500 flex items-center gap-1.5 uppercase tracking-widest">
                <Clock className="w-3 h-3" /> {new Date(item.published_at).toLocaleDateString()}
              </span>
            </div>

            <h4 className="text-2xl font-black text-white leading-tight group-hover:text-emerald-400 transition-colors">
              {item.title}
            </h4>

            {item.tldr && (
              <div className="space-y-3">
                <div className={`text-slate-400 text-base leading-relaxed ${!isExpanded && 'line-clamp-2'}`}>
                  {item.tldr}
                </div>
                {item.why_it_matters && isExpanded && (
                  <div className="pt-4 border-t border-white/5">
                    <h5 className="text-xs font-black uppercase tracking-[0.2em] text-blue-400 mb-2 mt-2">Critical Impact</h5>
                    <p className="text-slate-400 text-sm leading-relaxed">{item.why_it_matters}</p>
                  </div>
                )}
                <button 
                  onClick={() => setIsExpanded(!isExpanded)}
                  className="text-xs font-black text-emerald-500 uppercase tracking-widest hover:text-emerald-400 transition-colors"
                >
                  {isExpanded ? 'Collapse Insight' : 'Expand AI Insight'}
                </button>
              </div>
            )}

            <div className="flex items-center justify-between pt-4">
              <div className="flex flex-wrap gap-2">
                {item.tags?.slice(0, 3).map(tag => (
                  <span key={tag} className="px-3 py-1 bg-white/5 border border-white/10 rounded-xl text-[10px] font-black uppercase tracking-widest text-slate-400 group-hover:border-emerald-500/20 group-hover:text-emerald-500/80 transition-all">
                    {tag}
                  </span>
                ))}
              </div>
              
              <a 
                href={item.url} 
                target="_blank" 
                className="flex items-center gap-2 px-6 py-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 font-bold text-sm rounded-xl hover:bg-emerald-500 hover:text-black transition-all"
              >
                Access Source
                <ExternalLink className="w-4 h-4" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function PremiumSkeleton() {
  return (
    <div className="p-8 bg-white/5 border border-white/10 rounded-[2.5rem] animate-pulse">
      <div className="flex gap-6">
        <div className="w-14 h-14 bg-white/5 rounded-2xl" />
        <div className="flex-1 space-y-4">
          <div className="w-32 h-4 bg-white/5 rounded" />
          <div className="w-full h-8 bg-white/5 rounded" />
          <div className="w-full h-20 bg-white/5 rounded" />
        </div>
      </div>
    </div>
  );
}
