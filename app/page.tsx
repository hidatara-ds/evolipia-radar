"use client";

import { useEffect, useState, useCallback } from "react";
import {
  AlertCircle,
  ArrowRight,
  Bell,
  BrainCircuit,
  CheckCircle,
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
  Mail,
  Moon,
  RefreshCw,
  Search,
  Settings,
  Shield,
  Sparkles,
  Sun,
  TrendingUp,
  User,
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
  const [isDarkMode, setIsDarkMode] = useState(true);
  const [activeTab, setActiveTab] = useState("Signals");
  const [emailInput, setEmailInput] = useState("");
  const [subscribed, setSubscribed] = useState(false);

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

  const queryParamsKey = JSON.stringify(filterHook.queryParams);

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
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryParamsKey]);

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

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    if (emailInput.trim()) {
      setSubscribed(true);
      setEmailInput("");
    }
  };

  return (
    <main className={`min-h-screen font-sans transition-colors duration-300 ${
      isDarkMode
        ? "bg-[#050A0F] text-slate-200 selection:bg-emerald-500/30"
        : "bg-slate-50 text-slate-900 selection:bg-emerald-500/20"
    }`}>
      {/* Background Decorative Ambient Glows */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden z-0">
        {isDarkMode ? (
          <>
            <div className="absolute left-[-14rem] top-[-12rem] h-[34rem] w-[34rem] rounded-full bg-emerald-500/10 blur-[130px]" />
            <div className="absolute right-[-12rem] top-[18rem] h-[30rem] w-[30rem] rounded-full bg-cyan-500/10 blur-[130px]" />
            <div className="absolute bottom-[-16rem] left-1/3 h-[28rem] w-[28rem] rounded-full bg-indigo-500/10 blur-[140px]" />
          </>
        ) : (
          <>
            <div className="absolute left-[-14rem] top-[-12rem] h-[34rem] w-[34rem] rounded-full bg-emerald-400/15 blur-[120px]" />
            <div className="absolute right-[-12rem] top-[18rem] h-[30rem] w-[30rem] rounded-full bg-cyan-400/15 blur-[120px]" />
          </>
        )}
      </div>

      {/* 1. Top Navigation Bar */}
      <header className={`sticky top-0 z-[60] backdrop-blur-2xl transition-colors border-b ${
        isDarkMode
          ? "border-white/10 bg-[#050A0F]/80 text-white"
          : "border-slate-200 bg-white/80 text-slate-900 shadow-sm"
      }`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-3 sm:py-4 flex items-center justify-between gap-4">
          {/* Brand & Logo */}
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-3">
              <div className="relative shrink-0">
                <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-cyan-500 rounded-2xl blur opacity-30" />
                <img
                  src="/assets/icon.webp"
                  alt="Evolipia logo"
                  loading="lazy"
                  className="relative w-10 h-10 sm:w-11 sm:h-11 rounded-2xl border border-white/10 shadow-lg object-cover"
                />
              </div>
              <div>
                <h1 className="text-lg sm:text-xl font-black tracking-tight flex items-center gap-2">
                  Evolipia
                  <DataFreshness lastCrawledAt={lastCrawledAt} />
                </h1>
                <p className={`text-xs font-medium hidden sm:block ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
                  AI Intelligence Platform
                </p>
              </div>
            </div>

            {/* Navigation Links */}
            <nav className="hidden lg:flex items-center gap-1 bg-slate-500/10 p-1 rounded-2xl border border-white/5">
              {["Signals", "Analytics", "Sources", "Briefings"].map((tab) => {
                const isActive = activeTab === tab;
                return (
                  <button
                    key={tab}
                    onClick={() => setActiveTab(tab)}
                    className={`px-4 py-1.5 rounded-xl text-xs font-bold transition-all ${
                      isActive
                        ? "bg-emerald-500 text-slate-950 shadow-md shadow-emerald-500/20"
                        : isDarkMode
                        ? "text-slate-400 hover:text-slate-100 hover:bg-white/5"
                        : "text-slate-600 hover:text-slate-900 hover:bg-slate-200/60"
                    }`}
                  >
                    {tab}
                  </button>
                );
              })}
            </nav>
          </div>

          {/* Right Action Icons & User Avatar */}
          <div className="flex items-center gap-2 sm:gap-3">
            {/* Status Pill */}
            <div className={`hidden sm:inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold border ${
              isDarkMode ? "bg-slate-900 border-slate-800 text-slate-300" : "bg-slate-100 border-slate-200 text-slate-700"
            }`}>
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              <span>Last crawl: Just now</span>
            </div>

            {/* Light / Dark Mode Toggle */}
            <button
              onClick={() => setIsDarkMode(!isDarkMode)}
              className={`p-2.5 rounded-2xl border transition-all ${
                isDarkMode
                  ? "bg-slate-900 border-slate-800 text-amber-400 hover:bg-slate-800"
                  : "bg-slate-100 border-slate-200 text-slate-700 hover:bg-slate-200"
              }`}
              title={isDarkMode ? "Switch to Light Mode" : "Switch to Dark Mode"}
            >
              {isDarkMode ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4 text-slate-700" />}
            </button>

            {/* Notification Bell */}
            <button
              className={`relative p-2.5 rounded-2xl border transition-all ${
                isDarkMode
                  ? "bg-slate-900 border-slate-800 text-slate-400 hover:text-white hover:bg-slate-800"
                  : "bg-slate-100 border-slate-200 text-slate-600 hover:text-slate-900 hover:bg-slate-200"
              }`}
              title="Notifications"
            >
              <Bell className="w-4 h-4" />
              <span className="absolute top-2 right-2 w-2 h-2 rounded-full bg-emerald-400" />
            </button>

            {/* Settings Gear */}
            <button
              className={`p-2.5 rounded-2xl border transition-all ${
                isDarkMode
                  ? "bg-slate-900 border-slate-800 text-slate-400 hover:text-white hover:bg-slate-800"
                  : "bg-slate-100 border-slate-200 text-slate-600 hover:text-slate-900 hover:bg-slate-200"
              }`}
              title="System settings"
            >
              <Settings className="w-4 h-4" />
            </button>

            {/* User Profile Avatar */}
            <div className="flex items-center gap-2 pl-2 border-l border-slate-500/20">
              <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-gradient-to-tr from-emerald-500 to-cyan-500 p-0.5 shadow-md flex items-center justify-center cursor-pointer">
                <div className={`w-full h-full rounded-full flex items-center justify-center font-bold text-xs ${
                  isDarkMode ? "bg-slate-950 text-white" : "bg-white text-slate-900"
                }`}>
                  JD
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 py-6 sm:py-8 space-y-6 pb-20">
        
        {/* Research Modes Toolbar & Main Hero Banner */}
        <section className="grid grid-cols-1 xl:grid-cols-[minmax(0,1.5fr)_minmax(20rem,0.8fr)] gap-5">
          
          {/* Hero Banner: Live AI Signal Briefing + Agent Evoli Mascot */}
          <div className={`relative overflow-hidden rounded-[2rem] border p-5 sm:p-7 shadow-2xl flex flex-col justify-between space-y-6 ${
            isDarkMode
              ? "border-white/10 bg-slate-950/70"
              : "border-slate-200 bg-white"
          }`}>
            <div className={`absolute inset-0 pointer-events-none ${
              isDarkMode
                ? "bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_35%),radial-gradient(circle_at_bottom_right,rgba(59,130,246,0.16),transparent_35%)]"
                : "bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.08),transparent_35%),radial-gradient(circle_at_bottom_right,rgba(59,130,246,0.08),transparent_35%)]"
            }`} />
            
            {/* Top Row: Heading & Agent Evoli Mascot */}
            <div className="relative flex flex-col md:flex-row items-center justify-between gap-6">
              {/* Left Column: Heading & Description */}
              <div className="space-y-4 flex-1 min-w-0">
                <div className={`inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-bold ${
                  isDarkMode
                    ? "border-emerald-400/25 bg-emerald-400/10 text-emerald-300"
                    : "border-emerald-500/30 bg-emerald-50 text-emerald-800"
                }`}>
                  <Sparkles className="w-3.5 h-3.5" />
                  Live AI Signal Briefing
                </div>

                <div className="space-y-2">
                  <h2 className={`text-3xl sm:text-4xl lg:text-5xl font-black leading-tight tracking-tight ${
                    isDarkMode ? "text-white" : "text-slate-900"
                  }`}>
                    What is moving in the AI ecosystem?
                  </h2>
                  <p className={`max-w-2xl text-xs sm:text-sm leading-relaxed ${
                    isDarkMode ? "text-slate-400" : "text-slate-600"
                  }`}>
                    Evolipia Radar frames incoming articles into early intelligence signals, combining crawl freshness,
                    relevance, source diversity, and emerging themes into a daily research command center.
                  </p>
                </div>
              </div>

              {/* Right Column: Agent Evoli Mascot */}
              <div className="relative shrink-0 flex flex-col items-center md:items-end w-full md:w-auto max-w-[240px]">
                <div className="absolute inset-x-4 bottom-4 h-16 rounded-full bg-emerald-400/20 blur-2xl pointer-events-none" />
                <img
                  src="/assets/maskot1.webp"
                  alt="Agent Evoli research assistant"
                  fetchPriority="high"
                  loading="eager"
                  className="relative w-40 sm:w-48 lg:w-56 h-auto object-contain drop-shadow-[0_10px_20px_rgba(16,185,129,0.25)]"
                  style={{
                    maskImage: "linear-gradient(to bottom, black 75%, transparent 100%)",
                    WebkitMaskImage: "linear-gradient(to bottom, black 75%, transparent 100%)",
                  }}
                />
                <div className={`relative -mt-4 w-full rounded-2xl border p-3 backdrop-blur-md shadow-xl ${
                  isDarkMode ? "border-emerald-400/25 bg-black/80" : "border-emerald-500/30 bg-white/90"
                }`}>
                  <div className="flex items-center justify-between gap-2">
                    <div>
                      <p className="text-[10px] font-bold uppercase tracking-[0.24em] text-emerald-500">Agent Evoli</p>
                      <p className={`mt-0.5 text-xs line-clamp-1 font-medium ${isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
                        {isCrawling ? progressState?.message || "Collecting signals..." : "Standing by for next crawl"}
                      </p>
                    </div>
                    <div className={`h-2.5 w-2.5 rounded-full shrink-0 ${isCrawling ? "bg-amber-400 animate-pulse" : "bg-emerald-400"}`} />
                  </div>
                </div>
              </div>
            </div>

            {/* Running Text Ticker Banner */}
            <div className={`relative w-full rounded-2xl border backdrop-blur-md overflow-hidden p-1.5 flex items-center shadow-xl ${
              isDarkMode ? "border-emerald-500/30 bg-black/60" : "border-emerald-500/30 bg-slate-100"
            }`}>
              <div className="shrink-0 flex items-center gap-2 bg-gradient-to-r from-emerald-600 to-cyan-600 px-3 py-1.5 rounded-xl text-white text-[11px] font-black tracking-wider uppercase shadow-md mr-3 z-10">
                <span className="w-2 h-2 rounded-full bg-white animate-ping" />
                TRENDING SIGNALS
              </div>

              <div className="relative overflow-hidden w-full flex items-center py-1">
                <div className={`animate-marquee whitespace-nowrap flex items-center gap-8 text-xs font-semibold ${
                  isDarkMode ? "text-slate-200" : "text-slate-800"
                }`}>
                  {items.map((item, idx) => (
                    <a
                      key={`ticker-1-${item.id || idx}`}
                      href={item.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-2 hover:text-emerald-500 transition-colors group cursor-pointer"
                    >
                      <span className="px-2 py-0.5 rounded bg-emerald-500/15 border border-emerald-500/30 text-emerald-600 font-mono text-[10px] font-bold">
                        {item.source_name || item.domain || "GLOBAL"}
                      </span>
                      <span className="font-bold group-hover:underline">{item.title}</span>
                      {item.relevance_score ? (
                        <span className="text-emerald-500 font-mono text-[11px]">
                          [{item.relevance_score}%]
                        </span>
                      ) : null}
                      <span className="text-slate-400 mx-2">•</span>
                    </a>
                  ))}

                  {/* Duplicate list for seamless infinite loop */}
                  {items.map((item, idx) => (
                    <a
                      key={`ticker-2-${item.id || idx}`}
                      href={item.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-2 hover:text-emerald-500 transition-colors group cursor-pointer"
                    >
                      <span className="px-2 py-0.5 rounded bg-emerald-500/15 border border-emerald-500/30 text-emerald-600 font-mono text-[10px] font-bold">
                        {item.source_name || item.domain || "GLOBAL"}
                      </span>
                      <span className="font-bold group-hover:underline">{item.title}</span>
                      {item.relevance_score ? (
                        <span className="text-emerald-500 font-mono text-[11px]">
                          [{item.relevance_score}%]
                        </span>
                      ) : null}
                      <span className="text-slate-400 mx-2">•</span>
                    </a>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Consolidated Operations & Signal Health Sidebar */}
          <aside className={`rounded-[2rem] border p-5 shadow-2xl flex flex-col justify-between ${
            isDarkMode ? "border-white/10 bg-slate-950/70 text-white" : "border-slate-200 bg-white text-slate-900"
          }`}>
            <div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Operations</p>
                  <h3 className="mt-1 text-xl font-black">Signal Health</h3>
                </div>
                <Gauge className="h-6 w-6 text-emerald-500" />
              </div>

              {/* Progress Indicators Flow */}
              <div className="mt-5 space-y-3">
                <BriefingStat isDarkMode={isDarkMode} label="Fresh signals" value={paginationInfo.filteredCount} detail="matching current filters" icon={<Zap className="h-4 w-4" />} />
                <BriefingStat isDarkMode={isDarkMode} label="High confidence" value={highConfidenceCount} detail="relevance ≥ 75%" icon={<Shield className="h-4 w-4" />} />
                <BriefingStat isDarkMode={isDarkMode} label="Source diversity" value={sourceDiversity || "—"} detail="unique sources in view" icon={<Globe className="h-4 w-4" />} />
              </div>
            </div>

            <button
              onClick={startManualCrawl}
              disabled={isCrawling}
              className="mt-5 w-full rounded-2xl bg-emerald-500 px-4 py-3.5 text-sm font-black text-slate-950 shadow-lg shadow-emerald-500/20 transition-all hover:bg-emerald-400 disabled:cursor-not-allowed disabled:bg-slate-400 disabled:text-slate-700 focus:outline-none focus:ring-2 focus:ring-emerald-400/60"
            >
              {isCrawling ? "Crawl running..." : "Trigger Signal Crawl"}
            </button>
          </aside>
        </section>

        {/* Crawl Progress Notification Component */}
        <CrawlProgress
          progressState={progressState}
          isCrawling={isCrawling}
          onStartManualCrawl={startManualCrawl}
          toastMessage={toastMessage}
          onClearToast={clearToast}
        />

        {/* Consolidated Key Stats & Line Graph Cards */}
        <section className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
          <MetricCard
            isDarkMode={isDarkMode}
            label="Signals Processed"
            value={metrics?.articles_processed ?? 12}
            detail="crawler throughput"
            icon={<FileText className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard
            isDarkMode={isDarkMode}
            label="Active Signal View"
            value={paginationInfo.filteredCount}
            detail="after filters"
            icon={<Eye className="w-5 h-5" />}
            loading={loading}
          />
          <MetricCard
            isDarkMode={isDarkMode}
            label="Total Knowledge Items"
            value={paginationInfo.totalCount}
            detail="current compatibility dataset"
            icon={<BrainCircuit className="w-5 h-5" />}
            highlight
            loading={loading}
          />
          <MetricCard
            isDarkMode={isDarkMode}
            label="Cluster Momentum"
            value={metrics?.avg_cluster_score?.toFixed(1) ?? "8.4"}
            detail="average cluster score"
            icon={<TrendingUp className="w-5 h-5" />}
            suffix="/10"
            showGraph
            loading={loading}
          />
        </section>

        {/* Research Modes Toolbar Section */}
        <section className={`rounded-[1.75rem] border p-4 sm:p-5 ${
          isDarkMode ? "border-white/10 bg-slate-950/70" : "border-slate-200 bg-white shadow-sm"
        }`}>
          <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Research modes</p>
              <h3 className={`mt-1 text-lg font-black ${isDarkMode ? "text-white" : "text-slate-900"}`}>
                Start with a signal preset
              </h3>
            </div>
            <div className="flex flex-wrap gap-2">
              <QuickPreset isDarkMode={isDarkMode} label="Today" onClick={() => applyQuickPreset("today")} icon={<Clock className="h-4 w-4" />} />
              <QuickPreset isDarkMode={isDarkMode} label="High Relevance" onClick={() => applyQuickPreset("high")} icon={<Flame className="h-4 w-4" />} />
              <QuickPreset isDarkMode={isDarkMode} label="Research" onClick={() => applyQuickPreset("research")} icon={<Layers3 className="h-4 w-4" />} />
              <QuickPreset isDarkMode={isDarkMode} label="Developer" onClick={() => applyQuickPreset("developer")} icon={<Search className="h-4 w-4" />} />
              <QuickPreset isDarkMode={isDarkMode} label="Agents" onClick={() => applyQuickPreset("agents")} icon={<Sparkles className="h-4 w-4" />} />
            </div>
          </div>
        </section>

        {/* Central Filter Panel */}
        <FilterBar filterHook={filterHook} isDarkMode={isDarkMode} />

        {/* Signal Stream / Article List */}
        <section className="space-y-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Signal stream</p>
              <h3 className={`text-2xl font-black ${isDarkMode ? "text-white" : "text-slate-900"}`}>
                Latest intelligence signals
              </h3>
            </div>
            <p className="text-xs text-slate-500">
              Showing {items.length} of {paginationInfo.filteredCount} matching signals
            </p>
          </div>

          {newsLoading ? (
            <div className={`py-20 text-center flex flex-col items-center gap-3 rounded-[1.75rem] border ${
              isDarkMode ? "border-white/10 bg-slate-950/70 text-slate-400" : "border-slate-200 bg-white text-slate-500"
            }`}>
              <RefreshCw className="w-8 h-8 animate-spin text-emerald-500" />
              <span>Loading filtered signals...</span>
            </div>
          ) : error ? (
            <div className="py-12 bg-rose-950/30 border border-rose-800/40 rounded-[1.75rem] text-center text-rose-300 p-6">
              <AlertCircle className="w-8 h-8 mx-auto mb-2 text-rose-400" />
              <p className="font-semibold">{error}</p>
            </div>
          ) : items.length === 0 ? (
            <div className={`py-16 text-center rounded-[1.75rem] border ${
              isDarkMode ? "bg-slate-900/40 border-slate-800 text-slate-400" : "bg-slate-100 border-slate-200 text-slate-600"
            }`}>
              <Sparkles className="w-8 h-8 mx-auto mb-2 text-slate-400" />
              <p className="font-medium">No signals match your active filters.</p>
              <button
                onClick={filterHook.resetFilters}
                className="mt-3 text-xs font-bold text-emerald-500 hover:underline focus:outline-none focus:ring-2 focus:ring-emerald-400/60 rounded"
              >
                Reset Filters
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
              {/* Article Cards Grid */}
              <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-4">
                {items.map((item) => (
                  <SignalCard key={item.id} item={item} isDarkMode={isDarkMode} />
                ))}
              </div>

              {/* Briefing Queue Sidebar */}
              <aside className="hidden lg:block space-y-4">
                <div className={`sticky top-24 rounded-[1.75rem] border p-5 ${
                  isDarkMode ? "border-white/10 bg-slate-950/75" : "border-slate-200 bg-white shadow-md"
                }`}>
                  <p className="text-xs font-bold uppercase tracking-[0.24em] text-slate-500">Briefing queue</p>
                  <h4 className={`mt-1 text-lg font-black ${isDarkMode ? "text-white" : "text-slate-900"}`}>
                    Top signals to inspect
                  </h4>
                  <div className="mt-4 space-y-3">
                    {topSignals.map((item, index) => (
                      <a
                        key={item.id}
                        href={item.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className={`block rounded-2xl border p-3 transition-all hover:border-emerald-500/40 hover:bg-emerald-500/5 ${
                          isDarkMode ? "border-white/10 bg-white/[0.03]" : "border-slate-200 bg-slate-50"
                        }`}
                      >
                        <p className="text-[10px] font-black uppercase tracking-[0.22em] text-emerald-500">
                          Priority 0{index + 1}
                        </p>
                        <p className={`mt-1 line-clamp-2 text-sm font-bold ${isDarkMode ? "text-slate-100" : "text-slate-900"}`}>
                          {item.title}
                        </p>
                        <p className="mt-2 text-xs text-slate-500">{item.source_name || item.domain}</p>
                      </a>
                    ))}
                  </div>
                </div>
              </aside>
            </div>
          )}
        </section>

        {/* Pagination Section */}
        {paginationInfo.totalPages > 1 && (
          <div className={`flex flex-col sm:flex-row items-center justify-between gap-4 mt-8 pt-4 border-t text-xs ${
            isDarkMode ? "border-slate-800" : "border-slate-200"
          }`}>
            <span className={isDarkMode ? "text-slate-400 font-medium" : "text-slate-600 font-medium"}>
              Showing page <span className={`font-bold ${isDarkMode ? "text-slate-200" : "text-slate-900"}`}>{filterHook.page}</span> of{" "}
              <span className={`font-bold ${isDarkMode ? "text-slate-200" : "text-slate-900"}`}>{paginationInfo.totalPages}</span> ({paginationInfo.filteredCount} total signals)
            </span>

            <div className="flex items-center gap-1.5">
              <button
                disabled={filterHook.page <= 1}
                onClick={() => filterHook.setPage(filterHook.page - 1)}
                className={`px-3 py-1.5 rounded-xl border flex items-center gap-1 font-semibold transition-all disabled:opacity-40 disabled:cursor-not-allowed ${
                  isDarkMode ? "bg-slate-900 border-slate-800 text-slate-300 hover:bg-slate-800" : "bg-slate-100 border-slate-200 text-slate-700 hover:bg-slate-200"
                }`}
                aria-label="Previous page"
              >
                <ChevronLeft className="w-4 h-4" />
                <span>Prev</span>
              </button>

              {Array.from({ length: Math.min(5, paginationInfo.totalPages) }, (_, i) => {
                let pNum = i + 1;
                if (paginationInfo.totalPages > 5 && filterHook.page > 3) {
                  pNum = filterHook.page - 2 + i;
                  if (pNum > paginationInfo.totalPages) {
                    pNum = paginationInfo.totalPages - (4 - i);
                  }
                }
                const isActive = filterHook.page === pNum;
                return (
                  <button
                    key={pNum}
                    onClick={() => filterHook.setPage(pNum)}
                    className={`w-8 h-8 rounded-xl font-bold flex items-center justify-center transition-all ${
                      isActive
                        ? "bg-emerald-500 text-slate-950 shadow-md shadow-emerald-500/20 font-black"
                        : isDarkMode
                        ? "bg-slate-900 border border-slate-800 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                        : "bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 hover:bg-slate-200"
                    }`}
                  >
                    {pNum}
                  </button>
                );
              })}

              <button
                disabled={filterHook.page >= paginationInfo.totalPages}
                onClick={() => filterHook.setPage(filterHook.page + 1)}
                className={`px-3 py-1.5 rounded-xl border flex items-center gap-1 font-semibold transition-all disabled:opacity-40 disabled:cursor-not-allowed ${
                  isDarkMode ? "bg-slate-900 border-slate-800 text-slate-300 hover:bg-slate-800" : "bg-slate-100 border-slate-200 text-slate-700 hover:bg-slate-200"
                }`}
                aria-label="Next page"
              >
                <span>Next</span>
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}

        {/* 6. Consolidated CTA Section (Dark-patterned CTA card) */}
        <section className="relative overflow-hidden rounded-[2.5rem] border border-white/10 bg-gradient-to-br from-slate-950 via-slate-900 to-black p-8 sm:p-12 text-white shadow-2xl mt-16">
          <div className="absolute -right-20 -top-20 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl pointer-events-none" />
          <div className="absolute -left-20 -bottom-20 h-64 w-64 rounded-full bg-cyan-500/10 blur-3xl pointer-events-none" />

          <div className="relative max-w-3xl mx-auto text-center space-y-6">
            <div className="inline-flex items-center gap-2 rounded-full border border-emerald-400/30 bg-emerald-400/10 px-4 py-1.5 text-xs font-bold text-emerald-300">
              <Sparkles className="w-4 h-4" />
              Stay Ahead Of The Curve
            </div>

            <h3 className="text-3xl sm:text-4xl lg:text-5xl font-black tracking-tight text-white">
              Get real-time AI signal alerts
            </h3>

            <p className="text-sm sm:text-base text-slate-300 max-w-xl mx-auto leading-relaxed">
              Subscribe free to receive curated daily briefings on breaking AI models, research papers, and developer tools directly in your inbox.
            </p>

            {subscribed ? (
              <div className="inline-flex items-center gap-2 p-4 rounded-2xl bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 text-sm font-bold shadow-lg">
                <CheckCircle className="w-5 h-5 text-emerald-400" />
                Thank you for subscribing! You will receive daily AI signal briefings.
              </div>
            ) : (
              <form onSubmit={handleSubscribe} className="flex flex-col sm:flex-row items-center gap-3 max-w-lg mx-auto pt-2">
                <div className="relative flex-1 w-full">
                  <Mail className="w-4 h-4 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
                  <input
                    type="email"
                    required
                    placeholder="Enter your email address..."
                    value={emailInput}
                    onChange={(e) => setEmailInput(e.target.value)}
                    className="w-full pl-11 pr-4 h-12 rounded-2xl bg-black/60 border border-white/15 text-sm text-white placeholder-slate-400 focus:outline-none focus:border-emerald-400 focus:ring-2 focus:ring-emerald-400/20 shadow-inner"
                  />
                </div>
                <button
                  type="submit"
                  className="w-full sm:w-auto h-12 px-6 rounded-2xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm shadow-xl shadow-emerald-500/25 transition-all shrink-0 active:scale-95"
                >
                  Subscribe Free
                </button>
              </form>
            )}
          </div>
        </section>

      </div>

      {/* 7. Unified Footer */}
      <footer className={`border-t py-12 transition-colors ${
        isDarkMode ? "border-white/10 bg-slate-950 text-slate-400" : "border-slate-200 bg-slate-100 text-slate-600"
      }`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-3">
            <img
              src="/assets/icon.webp"
              alt="Evolipia logo"
              className="w-8 h-8 rounded-xl border border-white/10 object-cover"
            />
            <div>
              <span className={`text-base font-black tracking-tight ${isDarkMode ? "text-white" : "text-slate-900"}`}>
                Evolipia
              </span>
              <span className="text-xs text-slate-500 block">AI Intelligence Platform</span>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-6 text-xs font-semibold">
            <a href="#" className="hover:text-emerald-500 transition-colors">Documentation</a>
            <a href="#" className="hover:text-emerald-500 transition-colors">API</a>
            <a href="#" className="hover:text-emerald-500 transition-colors">Privacy</a>
            <a href="#" className="hover:text-emerald-500 transition-colors">Terms</a>
          </div>

          <p className="text-xs text-slate-500">
            © 2026 Evolipia Radar. All rights reserved.
          </p>
        </div>
      </footer>
    </main>
  );
}

function getSignalPhase(score: number) {
  if (score >= 85) return { label: "Accelerating", style: "bg-emerald-400/10 text-emerald-500 border-emerald-400/30" };
  if (score >= 60) return { label: "Emerging", style: "bg-cyan-400/10 text-cyan-500 border-cyan-400/30" };
  return { label: "Watch", style: "bg-amber-400/10 text-amber-500 border-amber-400/30" };
}

function formatCategory(category?: string) {
  if (!category) return "General AI";
  return category
    .split(/[-_\s]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function SignalCard({ item, isDarkMode }: { item: NewsItem; isDarkMode: boolean }) {
  const relevance = item.relevance_score ?? Math.round((item.scaled_score || 8.5) * 10);
  const phase = getSignalPhase(relevance);
  const source = item.source_name || item.domain || "Unknown source";

  return (
    <article className={`group flex min-h-[17rem] flex-col justify-between rounded-[1.5rem] border p-5 shadow-lg transition-all hover:-translate-y-0.5 ${
      isDarkMode
        ? "border-white/10 bg-slate-950/70 hover:border-emerald-400/35 hover:bg-slate-900/85 text-slate-100"
        : "border-slate-200 bg-white hover:border-emerald-500/40 hover:shadow-xl text-slate-900"
    }`}>
      <div>
        <div className="mb-3 flex items-start justify-between gap-3 text-xs">
          <div className="min-w-0">
            <span className={`inline-flex max-w-full items-center gap-1.5 rounded-full border px-2.5 py-1 font-bold ${
              isDarkMode ? "border-white/10 bg-white/[0.04] text-slate-300" : "border-slate-200 bg-slate-100 text-slate-700"
            }`}>
              <Globe className="h-3.5 w-3.5 shrink-0 text-slate-400" />
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
          className={`flex items-start gap-2 text-base font-black leading-snug transition-colors line-clamp-2 hover:text-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-400/60 rounded ${
            isDarkMode ? "text-slate-100" : "text-slate-900"
          }`}
        >
          {item.title}
          <ExternalLink className="mt-1 h-4 w-4 shrink-0 opacity-0 transition-opacity group-hover:opacity-100 text-emerald-500" />
        </a>

        {(item.raw_excerpt || item.tldr) && (
          <p className={`mt-3 text-sm leading-relaxed line-clamp-3 ${
            isDarkMode ? "text-slate-400" : "text-slate-600"
          }`}>
            {item.raw_excerpt || item.tldr}
          </p>
        )}

        {item.why_it_matters && (
          <div className={`mt-3 p-2.5 rounded-xl border text-xs leading-relaxed ${
            isDarkMode ? "bg-emerald-950/20 border-emerald-500/20 text-emerald-200" : "bg-emerald-50 border-emerald-200 text-emerald-900"
          }`}>
            <span className="font-bold text-emerald-500 block mb-0.5">Why it matters:</span>
            {item.why_it_matters}
          </div>
        )}

        {item.tags && item.tags.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-1.5">
            {item.tags.map((tag, idx) => (
              <span key={idx} className={`px-2 py-0.5 rounded-md text-[10px] font-mono border ${
                isDarkMode ? "bg-slate-900 text-slate-400 border-slate-800" : "bg-slate-100 text-slate-600 border-slate-200"
              }`}>
                #{tag}
              </span>
            ))}
          </div>
        )}
      </div>

      <div className="mt-5 space-y-4">
        <div>
          <div className="mb-1.5 flex items-center justify-between text-xs">
            <span className="font-semibold text-slate-400">Signal score proxy</span>
            <span className="font-mono font-black text-emerald-500">{relevance}%</span>
          </div>
          <div className={`h-2 rounded-full overflow-hidden ${isDarkMode ? "bg-slate-800" : "bg-slate-200"}`}>
            <div
              className="h-full rounded-full bg-gradient-to-r from-emerald-400 via-cyan-400 to-blue-400"
              style={{ width: `${Math.max(8, Math.min(100, relevance))}%` }}
            />
          </div>
        </div>

        <div className={`flex flex-wrap items-center justify-between gap-3 border-t pt-3 text-xs ${
          isDarkMode ? "border-white/10 text-slate-400" : "border-slate-200 text-slate-500"
        }`}>
          <span className="inline-flex items-center gap-1.5">
            <Clock className="h-3.5 w-3.5" />
            {new Date(item.published_at).toLocaleDateString()}
          </span>
          <span className={`rounded-full border px-2.5 py-1 font-bold ${
            isDarkMode ? "border-white/10 bg-white/[0.04] text-slate-300" : "border-slate-200 bg-slate-100 text-slate-700"
          }`}>
            {formatCategory(item.category)}
          </span>
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 font-black text-emerald-500 hover:text-emerald-400"
          >
            Inspect <ArrowRight className="h-3.5 w-3.5" />
          </a>
        </div>
      </div>
    </article>
  );
}

function BriefingStat({ label, value, detail, icon, isDarkMode }: { label: string; value: string | number; detail: string; icon: React.ReactNode; isDarkMode: boolean }) {
  return (
    <div className={`rounded-2xl border p-4 ${
      isDarkMode ? "border-white/10 bg-white/[0.03]" : "border-slate-200 bg-slate-50"
    }`}>
      <div className="flex items-center justify-between gap-3">
        <div>
          <p className="text-xs font-semibold text-slate-500">{label}</p>
          <p className={`mt-1 text-2xl font-black ${isDarkMode ? "text-white" : "text-slate-900"}`}>{value}</p>
        </div>
        <div className="rounded-xl bg-emerald-500/10 p-2 text-emerald-500">{icon}</div>
      </div>
      <p className="mt-2 text-xs text-slate-500">{detail}</p>
    </div>
  );
}

function QuickPreset({ label, onClick, icon, isDarkMode }: { label: string; onClick: () => void; icon: React.ReactNode; isDarkMode: boolean }) {
  return (
    <button
      onClick={onClick}
      className={`inline-flex items-center gap-2 rounded-full border px-3.5 py-2 text-xs font-black transition-all focus:outline-none focus:ring-2 focus:ring-emerald-400/60 ${
        isDarkMode
          ? "border-white/10 bg-slate-950/70 text-slate-300 hover:border-emerald-400/40 hover:bg-emerald-400/10 hover:text-emerald-200"
          : "border-slate-200 bg-slate-100 text-slate-700 hover:border-emerald-500/40 hover:bg-emerald-50 hover:text-emerald-800"
      }`}
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
  showGraph = false,
  loading = false,
  isDarkMode = true,
}: {
  label: string;
  value: string | number;
  icon: React.ReactNode;
  detail: string;
  suffix?: string;
  highlight?: boolean;
  showGraph?: boolean;
  loading?: boolean;
  isDarkMode?: boolean;
}) {
  return (
    <div className={`relative group p-4 sm:p-5 rounded-[1.5rem] border transition-all duration-300 overflow-hidden ${
      highlight
        ? isDarkMode
          ? "bg-gradient-to-br from-emerald-500/15 via-cyan-500/10 to-transparent border-emerald-500/30 shadow-lg"
          : "bg-emerald-50/60 border-emerald-200 shadow-sm"
        : isDarkMode
        ? "bg-white/[0.025] border-white/10 hover:border-white/15"
        : "bg-white border-slate-200 hover:border-slate-300 shadow-sm"
    }`}>
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs font-bold uppercase tracking-[0.16em] text-slate-500">{label}</span>
        <div className={`p-2 rounded-xl ${highlight ? "bg-emerald-500/20 text-emerald-500" : "bg-emerald-500/10 text-emerald-500"}`}>
          {icon}
        </div>
      </div>

      <div className="flex items-baseline justify-between gap-1">
        <div className="flex items-baseline gap-1">
          {loading ? (
            <div className="h-8 w-16 bg-slate-400/20 rounded animate-pulse" />
          ) : (
            <span className={`text-2xl sm:text-3xl font-black tracking-tight ${isDarkMode ? "text-white" : "text-slate-900"}`}>
              {value}
            </span>
          )}
          {suffix && <span className="text-xs font-semibold text-slate-500">{suffix}</span>}
        </div>

        {/* Sparkline Graph visual for Cluster Momentum */}
        {showGraph && (
          <div className="w-16 h-8 shrink-0">
            <svg viewBox="0 0 60 25" className="w-full h-full stroke-emerald-500 fill-none stroke-[2.5] stroke-linecap-round stroke-linejoin-round">
              <path d="M 0 18 Q 15 15, 25 10 T 45 6 T 60 2" />
            </svg>
          </div>
        )}
      </div>

      <p className="mt-2 text-xs text-slate-500">{detail}</p>
    </div>
  );
}
