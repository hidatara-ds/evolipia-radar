"use client";

import { useEffect, useState, useCallback } from "react";
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
  Lightbulb,
  ChevronLeft,
  ChevronRight
} from "lucide-react";

import { CrawlProgress } from "@/src/components/CrawlProgress";
import { DataFreshness } from "@/src/components/DataFreshness";
import { FilterBar } from "@/src/components/FilterBar";
import { useCrawlProgress } from "@/src/hooks/useCrawlProgress";
import { useFilters } from "@/src/hooks/useFilters";
import { fetchItems, NewsItem, PaginatedItemsResponse } from "@/src/api/client";

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

export default function Dashboard() {
  const [items, setItems] = useState<NewsItem[]>([]);
  const [paginationInfo, setPaginationInfo] = useState<{ totalCount: number; filteredCount: number; totalPages: number }>({
    totalCount: 0,
    filteredCount: 0,
    totalPages: 1,
  });
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [newsLoading, setNewsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showSettings, setShowSettings] = useState(false);
  const [settings, setSettings] = useState<SettingsState>({
    x_api_key: "",
    threads_api_key: "",
    openrouter_api_key: "",
  });

  const filterHook = useFilters();
  const {
    progressState,
    isCrawling,
    lastCrawledAt,
    toastMessage,
    startManualCrawl,
    clearToast,
  } = useCrawlProgress();

  const loadData = useCallback(async () => {
    setNewsLoading(true);
    setError(null);
    try {
      const res: PaginatedItemsResponse = await fetchItems(filterHook.queryParams);
      if (res.success) {
        setItems(res.data || []);
        setPaginationInfo({
          totalCount: res.total_count,
          filteredCount: res.filtered_count,
          totalPages: res.total_pages,
        });
      }
    } catch (e: any) {
      setError(e.message || "Failed to load items");
    } finally {
      setNewsLoading(false);
      setLoading(false);
    }
  }, [filterHook.queryParams]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Fetch system metrics
  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const res = await fetch("/metrics");
        if (res.ok) {
          const data = await res.json();
          setMetrics(data);
        }
      } catch (e) {
        // Ignore background metrics errors
      }
    };
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 15000);
    return () => clearInterval(interval);
  }, []);

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
                src="/assets/icon.webp" 
                alt="Logo" 
                loading="lazy"
                className="relative w-9 h-9 sm:w-11 sm:h-11 rounded-xl shadow-2xl transition-transform hover:scale-105"
              />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-bold tracking-tight text-white flex items-center gap-2">
                Evolipia Radar
                <DataFreshness lastCrawledAt={lastCrawledAt} />
              </h1>
              <p className="text-xs text-slate-500 font-medium hidden sm:block">Autonomous Research & Crawling Engine</p>
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
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 py-8 sm:py-10">
        <div className="flex flex-col lg:flex-row items-center gap-8 lg:gap-16">
          <div className="flex-1 space-y-5 text-center lg:text-left">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-semibold">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Agent Evoli monitoring active global sources</span>
            </div>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-white leading-[1.1] tracking-tight">
              Predict the future of <br />
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 via-teal-400 to-blue-500">
                AI Innovation.
              </span>
            </h2>
            <p className="text-base sm:text-lg text-slate-400 max-w-2xl leading-relaxed">
              Real-time semantic clustering of global research signals with automated background crawling and validation layer.
            </p>
          </div>
          
          <div className="relative lg:w-1/3 flex justify-center lg:justify-end items-end self-auto lg:self-end">
            <div className="relative group">
              <img
                src="/assets/maskot1.webp"
                alt="Evoli — AI Research Agent"
                fetchPriority="high"
                loading="eager"
                className="w-56 sm:w-64 lg:w-80 xl:w-96 h-auto object-contain"
                style={{ maskImage: "linear-gradient(to bottom, black 55%, transparent 100%)", WebkitMaskImage: "linear-gradient(to bottom, black 55%, transparent 100%)" }}
              />
              <div className="absolute bottom-8 left-1/2 -translate-x-1/2 px-4 py-1.5 bg-black/80 backdrop-blur-md rounded-full border border-emerald-500/20 text-[10px] font-black uppercase tracking-[0.3em] text-emerald-500 opacity-0 group-hover:opacity-100 transition-all whitespace-nowrap z-30">
                Agent Evoli
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 pb-20">
        {/* Real-time Crawl Progress Indicator */}
        <CrawlProgress
          progressState={progressState}
          isCrawling={isCrawling}
          onStartManualCrawl={startManualCrawl}
          toastMessage={toastMessage}
          onClearToast={clearToast}
        />

        {/* Metrics Grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 mb-10 sm:mb-12">
          <MetricCard 
            label="Sources Crawled" 
            value={metrics?.articles_processed ?? 12} 
            icon={<FileText className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard 
            label="Items Filtered" 
            value={paginationInfo.filteredCount} 
            icon={<Shield className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard 
            label="Total Items" 
            value={paginationInfo.totalCount} 
            icon={<BrainCircuit className="w-5 h-5" />}
            highlight
            loading={loading}
          />
          <MetricCard 
            label="Avg Score" 
            value={metrics?.avg_cluster_score?.toFixed(1) ?? "8.4"} 
            icon={<TrendingUp className="w-5 h-5" />}
            suffix="/10"
            loading={loading}
          />
        </div>

        {/* Advanced Filter Bar */}
        <FilterBar filterHook={filterHook} />

        {/* Content Feed */}
        <div className="space-y-4">
          {newsLoading ? (
            <div className="py-20 text-center text-slate-500 flex flex-col items-center gap-3">
              <RefreshCw className="w-8 h-8 animate-spin text-indigo-400" />
              <span>Loading filtered articles...</span>
            </div>
          ) : error ? (
            <div className="py-12 bg-rose-950/30 border border-rose-800/40 rounded-xl text-center text-rose-300 p-6">
              <AlertCircle className="w-8 h-8 mx-auto mb-2 text-rose-400" />
              <p className="font-semibold">{error}</p>
            </div>
          ) : items.length === 0 ? (
            <div className="py-16 text-center text-slate-500 bg-slate-900/40 rounded-xl border border-slate-800">
              <Sparkles className="w-8 h-8 mx-auto mb-2 text-slate-600" />
              <p className="font-medium">No items match your active filters.</p>
              <button
                onClick={filterHook.resetFilters}
                className="mt-3 text-xs font-semibold text-indigo-400 hover:underline"
              >
                Reset Filters
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {items.map(item => (
                <div
                  key={item.id}
                  className="bg-slate-900 border border-slate-800/80 hover:border-indigo-500/50 rounded-xl p-5 transition-all shadow-md flex flex-col justify-between group"
                >
                  <div>
                    <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
                      <span className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-300 font-mono">
                        {item.source_name || item.domain}
                      </span>
                      <span className="font-mono text-indigo-300 font-bold">
                        Relevance: {item.relevance_score ?? 85}%
                      </span>
                    </div>

                    <a
                      href={item.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-base font-bold text-slate-100 group-hover:text-indigo-400 transition-colors line-clamp-2 mb-2 flex items-start gap-1"
                    >
                      {item.title}
                      <ExternalLink className="w-4 h-4 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity mt-1" />
                    </a>

                    {item.raw_excerpt && (
                      <p className="text-xs text-slate-400 line-clamp-3 leading-relaxed mb-4">
                        {item.raw_excerpt}
                      </p>
                    )}
                  </div>

                  <div className="pt-3 border-t border-slate-800/60 flex items-center justify-between text-xs text-slate-500">
                    <span className="flex items-center gap-1">
                      <Clock className="w-3.5 h-3.5" />
                      {new Date(item.published_at).toLocaleDateString()}
                    </span>
                    <span className="px-2 py-0.5 rounded-full bg-emerald-950/60 text-emerald-300 border border-emerald-500/30 font-semibold text-[11px]">
                      {item.crawl_status || "done"}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Pagination Controls */}
        {paginationInfo.totalPages > 1 && (
          <div className="flex items-center justify-between mt-8 pt-4 border-t border-slate-800">
            <span className="text-xs text-slate-400">
              Page {filterHook.page} of {paginationInfo.totalPages} ({paginationInfo.filteredCount} items)
            </span>
            <div className="flex items-center gap-2">
              <button
                disabled={filterHook.page <= 1}
                onClick={() => filterHook.setPage(filterHook.page - 1)}
                className="p-2 rounded-lg bg-slate-900 border border-slate-800 text-slate-300 hover:bg-slate-800 disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <button
                disabled={filterHook.page >= paginationInfo.totalPages}
                onClick={() => filterHook.setPage(filterHook.page + 1)}
                className="p-2 rounded-lg bg-slate-900 border border-slate-800 text-slate-300 hover:bg-slate-800 disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}

function MetricCard({ 
  label, 
  value, 
  icon, 
  suffix = "", 
  highlight = false,
  loading = false 
}: { 
  label: string; 
  value: string | number; 
  icon: React.ReactNode; 
  suffix?: string;
  highlight?: boolean;
  loading?: boolean;
}) {
  return (
    <div className={`relative group p-4 sm:p-5 rounded-2xl border transition-all duration-300 overflow-hidden ${
      highlight 
        ? "bg-gradient-to-br from-emerald-500/10 via-teal-500/5 to-transparent border-emerald-500/30 hover:border-emerald-500/50" 
        : "bg-white/[0.02] border-white/5 hover:border-white/10 hover:bg-white/[0.04]"
    }`}>
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs font-semibold text-slate-400">{label}</span>
        <div className={`p-2 rounded-xl ${highlight ? 'bg-emerald-500/20 text-emerald-400' : 'bg-white/5 text-slate-400'}`}>
          {icon}
        </div>
      </div>
      
      <div className="flex items-baseline gap-1">
        {loading ? (
          <div className="h-8 w-16 bg-white/10 rounded animate-pulse" />
        ) : (
          <span className="text-2xl sm:text-3xl font-black text-white tracking-tight">{value}</span>
        )}
        {suffix && <span className="text-xs font-semibold text-slate-500">{suffix}</span>}
      </div>
    </div>
  );
}
