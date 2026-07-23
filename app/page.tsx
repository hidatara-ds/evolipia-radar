"use client";

import { useEffect, useState, useCallback } from "react";
import {
  AlertCircle,
  ArrowRight,
  BrainCircuit,
  ChevronLeft,
  ChevronRight,
  Clock,
  ExternalLink,
  Eye,
  FileText,
  Flame,
  Gauge,
  Globe,
  Layers3,
  RefreshCw,
  Search,
  Settings,
  Shield,
  Sparkles,
  TrendingUp,
  Zap,
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

const DEFAULT_TOPICS = ["LLM", "Agents", "Open Source", "Infrastructure", "Research", "Security"];

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
      setError(e.message || "Failed to load signals");
    } finally {
      setNewsLoading(false);
      setLoading(false);
    }
  }, [filterHook.queryParams]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const res = await fetch("/metrics");
        if (res.ok) {
          const data = await res.json();
          setMetrics(data);
        }
      } catch (e) {
        // Ignore background metrics errors so the dashboard can still render.
      }
    };
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 15000);
    return () => clearInterval(interval);
  }, []);

  const topSignals = items.slice(0, 3);
  const topTopics = buildTopTopics(items, metrics?.top_cluster_titles);
  const highConfidenceCount = items.filter((item) => (item.relevance_score ?? 0) >= 75).length;
  const sourceDiversity = new Set(items.map((item) => item.source_name || item.domain).filter(Boolean)).size;

  const applyQuickPreset = (preset: "today" | "high" | "research" | "developer" | "agents") => {
    const today = new Date().toISOString().split("T")[0];
    filterHook.setPage(1);

    if (preset === "today") {
      filterHook.setDateRange("today");
      filterHook.setDateFrom(today);
      filterHook.setDateTo(today);
      filterHook.setMinRelevance(30);
      filterHook.setSelectedCategories([]);
      filterHook.setSortBy("date");
      filterHook.setSortOrder("desc");
      return;
    }

    if (preset === "high") {
      filterHook.setDateRange("7d");
      filterHook.setDateFrom("");
      filterHook.setDateTo("");
      filterHook.setMinRelevance(75);
      filterHook.setSelectedCategories([]);
      filterHook.setSortBy("relevance");
      filterHook.setSortOrder("desc");
      return;
    }

    if (preset === "research") {
      filterHook.setDateRange("30d");
      filterHook.setDateFrom("");
      filterHook.setDateTo("");
      filterHook.setMinRelevance(45);
      filterHook.setSelectedCategories(["research", "llm"]);
      filterHook.setSortBy("impact");
      filterHook.setSortOrder("desc");
      return;
    }

    if (preset === "developer") {
      filterHook.setDateRange("7d");
      filterHook.setDateFrom("");
      filterHook.setDateTo("");
      filterHook.setMinRelevance(40);
      filterHook.setSelectedCategories(["open-source", "infra"]);
      filterHook.setSortBy("impact");
      filterHook.setSortOrder("desc");
      return;
    }

    filterHook.setDateRange("7d");
    filterHook.setDateFrom("");
    filterHook.setDateTo("");
    filterHook.setMinRelevance(40);
    filterHook.setSelectedCategories(["agents"]);
    filterHook.setSortBy("relevance");
    filterHook.setSortOrder("desc");
  };

  return (
    <main className="min-h-screen overflow-hidden bg-[#050A0F] text-slate-200 font-sans selection:bg-emerald-500/30">
      <div className="fixed inset-0 pointer-events-none overflow-hidden">
        <div className="absolute left-[-14rem] top-[-12rem] h-[34rem] w-[34rem] rounded-full bg-emerald-500/10 blur-[130px]" />
        <div className="absolute right-[-12rem] top-[18rem] h-[30rem] w-[30rem] rounded-full bg-cyan-500/10 blur-[130px]" />
        <div className="absolute bottom-[-16rem] left-1/3 h-[28rem] w-[28rem] rounded-full bg-indigo-500/10 blur-[140px]" />
      </div>

      <header className="sticky top-0 z-[60] border-b border-white/10 bg-[#050A0F]/80 backdrop-blur-2xl">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-3 sm:py-4 flex items-center justify-between">
          <div className="flex items-center gap-3 sm:gap-4">
            <div className="relative">
              <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-cyan-500 rounded-2xl blur opacity-25" />
              <img
                src="/assets/icon.webp"
                alt="Evolipia Radar logo"
                loading="lazy"
                className="relative w-10 h-10 sm:w-11 sm:h-11 rounded-2xl border border-white/10 shadow-2xl"
              />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-black tracking-tight text-white flex flex-wrap items-center gap-2">
                Evolipia Radar
                <DataFreshness lastCrawledAt={lastCrawledAt} />
              </h1>
              <p className="text-xs text-slate-500 font-medium hidden sm:block">
                AI Intelligence Platform · Signal command center
              </p>
            </div>
          </div>

          <button
            className="p-2.5 rounded-2xl bg-white/5 border border-white/10 text-slate-400 hover:text-white hover:bg-white/10 hover:border-white/20 transition-all focus:outline-none focus:ring-2 focus:ring-emerald-400/60"
            title="System settings"
          >
            <Settings className="w-5 h-5" />
          </button>
        </div>
      </header>

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 py-6 sm:py-8 space-y-6 pb-20">
        <section className="grid grid-cols-1 xl:grid-cols-[minmax(0,1.5fr)_minmax(20rem,0.8fr)] gap-5">
          <div className="relative overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.035] p-5 sm:p-7 shadow-2xl">
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_35%),radial-gradient(circle_at_bottom_right,rgba(59,130,246,0.16),transparent_35%)]" />
            <div className="relative grid grid-cols-1 lg:grid-cols-[1fr_auto] gap-6 items-center">
              <div className="space-y-5">
                <div className="inline-flex items-center gap-2 rounded-full border border-emerald-400/25 bg-emerald-400/10 px-3 py-1 text-xs font-bold text-emerald-300">
                  <Sparkles className="w-3.5 h-3.5" />
                  Live AI Signal Briefing
                </div>

                <div className="space-y-3">
                  <h2 className="text-3xl sm:text-5xl font-black leading-tight tracking-tight text-white">
                    What is moving in the AI ecosystem?
                  </h2>
                  <p className="max-w-2xl text-sm sm:text-base leading-relaxed text-slate-400">
                    Evolipia Radar now frames incoming articles as early intelligence signals, combining crawl freshness,
                    relevance, source diversity, and emerging themes into a daily research command center.
                  </p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  {topTopics.map((topic, idx) => (
                    <div key={topic} className="rounded-2xl border border-white/10 bg-black/25 p-3">
                      <p className="text-[10px] uppercase tracking-[0.25em] text-slate-500">Theme 0{idx + 1}</p>
                      <p className="mt-1 font-bold text-slate-100 line-clamp-1">{topic}</p>
                      <div className="mt-3 h-1.5 rounded-full bg-slate-800 overflow-hidden">
                        <div
                          className="h-full rounded-full bg-gradient-to-r from-emerald-400 to-cyan-400"
                          style={{ width: `${86 - idx * 14}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="relative mx-auto flex w-full max-w-xs flex-col items-center lg:items-end">
                <div className="absolute inset-x-8 bottom-4 h-20 rounded-full bg-emerald-400/20 blur-3xl" />
                <img
                  src="/assets/maskot1.webp"
                  alt="Agent Evoli research assistant"
                  fetchPriority="high"
                  loading="eager"
                  className="relative w-44 sm:w-56 lg:w-64 h-auto object-contain"
                  style={{ maskImage: "linear-gradient(to bottom, black 62%, transparent 100%)", WebkitMaskImage: "linear-gradient(to bottom, black 62%, transparent 100%)" }}
                />
                <div className="relative -mt-6 w-full rounded-2xl border border-emerald-400/20 bg-black/55 p-4 backdrop-blur-md">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="text-xs font-bold uppercase tracking-[0.24em] text-emerald-300">Agent Evoli</p>
                      <p className="mt-1 text-sm text-slate-300">
                        {isCrawling ? progressState?.message || "Collecting new signals..." : "Standing by for the next crawl."}
                      </p>
                    </div>
                    <div className={`h-3 w-3 rounded-full ${isCrawling ? "bg-amber-300 animate-pulse" : "bg-emerald-400"}`} />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <aside className="rounded-[2rem] border border-white/10 bg-slate-950/70 p-5 shadow-2xl">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Operations</p>
                <h3 className="mt-1 text-xl font-black text-white">Signal Health</h3>
              </div>
              <Gauge className="h-6 w-6 text-emerald-300" />
            </div>

            <div className="mt-5 space-y-3">
              <BriefingStat label="Fresh signals" value={paginationInfo.filteredCount} detail="matching current filters" icon={<Zap className="h-4 w-4" />} />
              <BriefingStat label="High confidence" value={highConfidenceCount} detail="relevance ≥ 75%" icon={<Shield className="h-4 w-4" />} />
              <BriefingStat label="Source diversity" value={sourceDiversity || "—"} detail="unique sources in view" icon={<Globe className="h-4 w-4" />} />
            </div>

            <button
              onClick={startManualCrawl}
              disabled={isCrawling}
              className="mt-5 w-full rounded-2xl bg-emerald-500 px-4 py-3 text-sm font-black text-slate-950 transition-all hover:bg-emerald-400 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-300"
            >
              {isCrawling ? "Crawl running..." : "Trigger Signal Crawl"}
            </button>
          </aside>
        </section>

        <CrawlProgress
          progressState={progressState}
          isCrawling={isCrawling}
          onStartManualCrawl={startManualCrawl}
          toastMessage={toastMessage}
          onClearToast={clearToast}
        />

        <section className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
          <MetricCard
            label="Signals Processed"
            value={metrics?.articles_processed ?? 12}
            detail="crawler throughput"
            icon={<FileText className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard
            label="Active Signal View"
            value={paginationInfo.filteredCount}
            detail="after filters"
            icon={<Eye className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard
            label="Total Knowledge Items"
            value={paginationInfo.totalCount}
            detail="current compatibility dataset"
            icon={<BrainCircuit className="w-5 h-5" />}
            highlight
            loading={loading}
          />
          <MetricCard
            label="Cluster Momentum"
            value={metrics?.avg_cluster_score?.toFixed(1) ?? "8.4"}
            detail="average cluster score"
            icon={<TrendingUp className="w-5 h-5" />}
            suffix="/10"
            loading={loading}
          />
        </section>

        <section className="rounded-[1.75rem] border border-white/10 bg-white/[0.03] p-4 sm:p-5">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Research modes</p>
              <h3 className="mt-1 text-lg font-black text-white">Start with a signal preset</h3>
            </div>
            <div className="flex flex-wrap gap-2">
              <QuickPreset label="Today" onClick={() => applyQuickPreset("today")} icon={<Clock className="h-4 w-4" />} />
              <QuickPreset label="High Relevance" onClick={() => applyQuickPreset("high")} icon={<Flame className="h-4 w-4" />} />
              <QuickPreset label="Research" onClick={() => applyQuickPreset("research")} icon={<Layers3 className="h-4 w-4" />} />
              <QuickPreset label="Developer" onClick={() => applyQuickPreset("developer")} icon={<Search className="h-4 w-4" />} />
              <QuickPreset label="Agents" onClick={() => applyQuickPreset("agents")} icon={<Sparkles className="h-4 w-4" />} />
            </div>
          </div>
        </section>

        <FilterBar filterHook={filterHook} />

        <section className="space-y-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Signal stream</p>
              <h3 className="text-2xl font-black text-white">Latest intelligence signals</h3>
            </div>
            <p className="text-xs text-slate-500">
              Showing {items.length} of {paginationInfo.filteredCount} matching signals
            </p>
          </div>

          {newsLoading ? (
            <div className="py-20 text-center text-slate-500 flex flex-col items-center gap-3 rounded-[1.75rem] border border-white/10 bg-white/[0.03]">
              <RefreshCw className="w-8 h-8 animate-spin text-emerald-400" />
              <span>Loading filtered signals...</span>
            </div>
          ) : error ? (
            <div className="py-12 bg-rose-950/30 border border-rose-800/40 rounded-[1.75rem] text-center text-rose-300 p-6">
              <AlertCircle className="w-8 h-8 mx-auto mb-2 text-rose-400" />
              <p className="font-semibold">{error}</p>
            </div>
          ) : items.length === 0 ? (
            <div className="py-16 text-center text-slate-500 bg-slate-900/40 rounded-[1.75rem] border border-slate-800">
              <Sparkles className="w-8 h-8 mx-auto mb-2 text-slate-600" />
              <p className="font-medium">No signals match your active filters.</p>
              <button
                onClick={filterHook.resetFilters}
                className="mt-3 text-xs font-semibold text-emerald-400 hover:underline focus:outline-none focus:ring-2 focus:ring-emerald-400/60 rounded"
              >
                Reset Filters
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
              <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-4">
                {items.map((item) => (
                  <SignalCard key={item.id} item={item} />
                ))}
              </div>

              <aside className="hidden lg:block space-y-4">
                <div className="sticky top-24 rounded-[1.75rem] border border-white/10 bg-slate-950/75 p-5">
                  <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Briefing queue</p>
                  <h4 className="mt-1 text-lg font-black text-white">Top signals to inspect</h4>
                  <div className="mt-4 space-y-3">
                    {topSignals.map((item, index) => (
                      <a
                        key={item.id}
                        href={item.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="block rounded-2xl border border-white/10 bg-white/[0.03] p-3 transition-all hover:border-emerald-400/40 hover:bg-emerald-400/5"
                      >
                        <p className="text-[10px] font-black uppercase tracking-[0.22em] text-emerald-300">
                          Priority 0{index + 1}
                        </p>
                        <p className="mt-1 line-clamp-2 text-sm font-bold text-slate-100">{item.title}</p>
                        <p className="mt-2 text-xs text-slate-500">{item.source_name || item.domain}</p>
                      </a>
                    ))}
                  </div>
                </div>
              </aside>
            </div>
          )}
        </section>

        {paginationInfo.totalPages > 1 && (
          <div className="flex items-center justify-between mt-8 pt-4 border-t border-slate-800">
            <span className="text-xs text-slate-400">
              Page {filterHook.page} of {paginationInfo.totalPages} ({paginationInfo.filteredCount} signals)
            </span>
            <div className="flex items-center gap-2">
              <button
                disabled={filterHook.page <= 1}
                onClick={() => filterHook.setPage(filterHook.page - 1)}
                className="p-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 hover:bg-slate-800 disabled:opacity-40 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-emerald-400/60"
                aria-label="Previous page"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <button
                disabled={filterHook.page >= paginationInfo.totalPages}
                onClick={() => filterHook.setPage(filterHook.page + 1)}
                className="p-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 hover:bg-slate-800 disabled:opacity-40 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-emerald-400/60"
                aria-label="Next page"
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

function buildTopTopics(items: NewsItem[], clusterTitles?: string[] | null) {
  const categoryCounts = items.reduce<Record<string, number>>((acc, item) => {
    const key = formatCategory(item.category);
    acc[key] = (acc[key] || 0) + 1;
    return acc;
  }, {});

  const categories = Object.entries(categoryCounts)
    .sort((a, b) => b[1] - a[1])
    .map(([category]) => category)
    .filter(Boolean);

  const clusterTopics = (clusterTitles || []).slice(0, 3);
  const merged = [...clusterTopics, ...categories, ...DEFAULT_TOPICS];
  return Array.from(new Set(merged)).slice(0, 3);
}

function formatCategory(category?: string) {
  if (!category) return "General AI";
  return category
    .split(/[-_\s]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function getSignalPhase(score: number) {
  if (score >= 85) return { label: "Accelerating", style: "bg-emerald-400/10 text-emerald-300 border-emerald-400/30" };
  if (score >= 60) return { label: "Emerging", style: "bg-cyan-400/10 text-cyan-300 border-cyan-400/30" };
  return { label: "Watch", style: "bg-amber-400/10 text-amber-300 border-amber-400/30" };
}

function SignalCard({ item }: { item: NewsItem }) {
  const relevance = item.relevance_score ?? Math.round((item.scaled_score || 8.5) * 10);
  const phase = getSignalPhase(relevance);
  const source = item.source_name || item.domain || "Unknown source";

  return (
    <article className="group flex min-h-[17rem] flex-col justify-between rounded-[1.5rem] border border-white/10 bg-slate-950/70 p-5 shadow-lg transition-all hover:-translate-y-0.5 hover:border-emerald-400/35 hover:bg-slate-900/85">
      <div>
        <div className="mb-3 flex items-start justify-between gap-3 text-xs">
          <div className="min-w-0">
            <span className="inline-flex max-w-full items-center gap-1.5 rounded-full border border-white/10 bg-white/[0.04] px-2.5 py-1 font-bold text-slate-300">
              <Globe className="h-3.5 w-3.5 shrink-0 text-slate-500" />
              <span className="truncate">{source}</span>
            </span>
          </div>
          <span className={`shrink-0 rounded-full border px-2.5 py-1 text-[11px] font-black ${phase.style}`}>
            {phase.label}
          </span>
        </div>

        <a
          href={item.url}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-start gap-2 text-base font-black leading-snug text-slate-100 transition-colors line-clamp-2 group-hover:text-emerald-300 focus:outline-none focus:ring-2 focus:ring-emerald-400/60 rounded"
        >
          {item.title}
          <ExternalLink className="mt-1 h-4 w-4 shrink-0 opacity-0 transition-opacity group-hover:opacity-100" />
        </a>

        {item.raw_excerpt && (
          <p className="mt-3 text-sm leading-relaxed text-slate-400 line-clamp-3">
            {item.raw_excerpt}
          </p>
        )}
      </div>

      <div className="mt-5 space-y-4">
        <div>
          <div className="mb-1.5 flex items-center justify-between text-xs">
            <span className="font-semibold text-slate-500">Signal score proxy</span>
            <span className="font-mono font-black text-emerald-300">{relevance}%</span>
          </div>
          <div className="h-2 rounded-full bg-slate-800 overflow-hidden">
            <div
              className="h-full rounded-full bg-gradient-to-r from-emerald-400 via-cyan-400 to-blue-400"
              style={{ width: `${Math.max(8, Math.min(100, relevance))}%` }}
            />
          </div>
        </div>

        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-white/10 pt-3 text-xs text-slate-500">
          <span className="inline-flex items-center gap-1.5">
            <Clock className="h-3.5 w-3.5" />
            {new Date(item.published_at).toLocaleDateString()}
          </span>
          <span className="rounded-full border border-white/10 bg-white/[0.04] px-2.5 py-1 font-bold text-slate-300">
            {formatCategory(item.category)}
          </span>
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 font-bold text-emerald-300 hover:text-emerald-200"
          >
            Inspect <ArrowRight className="h-3.5 w-3.5" />
          </a>
        </div>
      </div>
    </article>
  );
}

function BriefingStat({ label, value, detail, icon }: { label: string; value: string | number; detail: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
      <div className="flex items-center justify-between gap-3">
        <div>
          <p className="text-xs font-semibold text-slate-500">{label}</p>
          <p className="mt-1 text-2xl font-black text-white">{value}</p>
        </div>
        <div className="rounded-xl bg-emerald-400/10 p-2 text-emerald-300">{icon}</div>
      </div>
      <p className="mt-2 text-xs text-slate-500">{detail}</p>
    </div>
  );
}

function QuickPreset({ label, onClick, icon }: { label: string; onClick: () => void; icon: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-slate-950/70 px-3 py-2 text-xs font-black text-slate-300 transition-all hover:border-emerald-400/40 hover:bg-emerald-400/10 hover:text-emerald-200 focus:outline-none focus:ring-2 focus:ring-emerald-400/60"
    >
      {icon}
      {label}
    </button>
  );
}

function MetricCard({
  label,
  value,
  icon,
  detail,
  suffix = "",
  highlight = false,
  loading = false,
}: {
  label: string;
  value: string | number;
  icon: React.ReactNode;
  detail: string;
  suffix?: string;
  highlight?: boolean;
  loading?: boolean;
}) {
  return (
    <div className={`relative group p-4 sm:p-5 rounded-[1.5rem] border transition-all duration-300 overflow-hidden ${
      highlight
        ? "bg-gradient-to-br from-emerald-500/12 via-cyan-500/6 to-transparent border-emerald-500/30 hover:border-emerald-500/50"
        : "bg-white/[0.025] border-white/10 hover:border-white/15 hover:bg-white/[0.045]"
    }`}>
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs font-bold uppercase tracking-[0.16em] text-slate-500">{label}</span>
        <div className={`p-2 rounded-xl ${highlight ? "bg-emerald-500/20 text-emerald-300" : "bg-white/5 text-slate-400"}`}>
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
      <p className="mt-2 text-xs text-slate-500">{detail}</p>
    </div>
  );
}
