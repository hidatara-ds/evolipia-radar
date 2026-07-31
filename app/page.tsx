"use client";

import { useEffect, useState, useCallback } from "react";
import {
  AlertCircle,
  ArrowRight,
  Bookmark,
  BrainCircuit,
  CheckCircle,
  ChevronLeft,
  ChevronRight,
  Clock,
  ExternalLink,
  Eye,
  FileText,
  Flame,
  Globe,
  Layers3,
  Mail,
  Moon,
  RefreshCw,
  Search,
  Share2,
  Shield,
  Sparkles,
  Sun,
  TrendingUp,
  Zap,
} from "lucide-react";

import { CrawlProgress } from "@/src/components/CrawlProgress";
import { FilterBar } from "@/src/components/FilterBar";
import { CommandPalette } from "@/src/components/CommandPalette";
import { AIDailyBriefing } from "@/src/components/AIDailyBriefing";
import { AIPulseKeywords } from "@/src/components/AIPulseKeywords";
import { BreakingSignals } from "@/src/components/BreakingSignals";
import { TrendingClusters } from "@/src/components/TrendingClusters";
import { IntelligenceSidebar } from "@/src/components/IntelligenceSidebar";

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

function AnimatedCount({ value, duration = 1400 }: { value: number | string; duration?: number }) {
  const rawStr = String(value).replace(/[^0-9.]/g, "");
  const targetNumber = parseFloat(rawStr);
  const isFloat = String(value).includes(".");
  const [current, setCurrent] = useState(0);

  useEffect(() => {
    if (isNaN(targetNumber) || targetNumber === 0) {
      setCurrent(targetNumber || 0);
      return;
    }

    let startTimestamp: number | null = null;
    const step = (timestamp: number) => {
      if (!startTimestamp) startTimestamp = timestamp;
      const progress = Math.min((timestamp - startTimestamp) / duration, 1);
      const easeOut = 1 - Math.pow(1 - progress, 3);
      setCurrent(targetNumber * easeOut);

      if (progress < 1) {
        window.requestAnimationFrame(step);
      }
    };

    window.requestAnimationFrame(step);
  }, [targetNumber, duration]);

  if (isNaN(targetNumber)) return <span>{value}</span>;
  if (isFloat) return <span>{current.toFixed(1)}</span>;
  return <span>{Math.round(current).toLocaleString()}</span>;
}

export default function Dashboard() {
  const [isDarkMode, setIsDarkMode] = useState(true);
  const [emailInput, setEmailInput] = useState("");
  const [subscribed, setSubscribed] = useState(false);
  const [isCmdPaletteOpen, setIsCmdPaletteOpen] = useState(false);

  // Bookmarking System
  const [bookmarkedIds, setBookmarkedIds] = useState<Set<string | number>>(new Set());
  const [showOnlyBookmarked, setShowOnlyBookmarked] = useState(false);

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
  const { progressState, isCrawling, toastMessage, startManualCrawl, clearToast } = useCrawlProgress();

  // Load bookmarks from local storage
  useEffect(() => {
    try {
      const saved = localStorage.getItem("evolipia_bookmarks");
      if (saved) {
        setBookmarkedIds(new Set(JSON.parse(saved)));
      }
    } catch (e) {
      // Ignore fallback
    }
  }, []);

  const toggleBookmark = (id: string | number) => {
    setBookmarkedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      try {
        localStorage.setItem("evolipia_bookmarks", JSON.stringify(Array.from(next)));
      } catch (e) {}
      return next;
    });
  };

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
  }, [queryParamsKey, filterHook.queryParams]);

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
        // Ignore background metrics errors
      }
    };
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 15000);
    return () => clearInterval(interval);
  }, []);

  const displayedItems = showOnlyBookmarked
    ? items.filter((item) => bookmarkedIds.has(item.id))
    : items;

  const topSignals = items.slice(0, 3);

  const handleTopicSelect = (keyword: string) => {
    filterHook.setSearch(keyword);
    filterHook.setPage(1);
  };

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    if (emailInput.trim()) {
      setSubscribed(true);
      setEmailInput("");
    }
  };

  const scrollToSubscribe = () => {
    document.getElementById("subscribe-section")?.scrollIntoView({ behavior: "smooth" });
  };

  return (
    <main
      className={`min-h-screen font-sans transition-colors duration-300 relative ${
        isDarkMode
          ? "bg-[#050A0F] text-[#F8FAFC] selection:bg-emerald-500/30"
          : "bg-slate-50 text-slate-900 selection:bg-emerald-500/20"
      }`}
    >
      {/* Command Palette Modal */}
      <CommandPalette
        isOpen={isCmdPaletteOpen}
        onClose={() => setIsCmdPaletteOpen(false)}
        items={items}
        onSelectTopic={handleTopicSelect}
        isDarkMode={isDarkMode}
      />

      {/* Background Ambient Glows */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden z-0">
        {isDarkMode ? (
          <>
            <div className="absolute left-[-14rem] top-[-12rem] h-[34rem] w-[34rem] rounded-full bg-emerald-500/10 blur-[140px]" />
            <div className="absolute right-[-12rem] top-[18rem] h-[30rem] w-[30rem] rounded-full bg-cyan-500/10 blur-[140px]" />
          </>
        ) : (
          <>
            <div className="absolute left-[-14rem] top-[-12rem] h-[34rem] w-[34rem] rounded-full bg-emerald-400/15 blur-[120px]" />
            <div className="absolute right-[-12rem] top-[18rem] h-[30rem] w-[30rem] rounded-full bg-cyan-400/15 blur-[120px]" />
          </>
        )}
      </div>

      {/* STICKY EDITORIAL TOP NAVIGATION BAR */}
      <header className={`sticky top-0 z-50 border-b backdrop-blur-xl transition-all ${
        isDarkMode ? "border-white/10 bg-[#050A0F]/85" : "border-slate-200 bg-white/85 shadow-sm"
      }`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between gap-4">
          {/* Logo & Brand */}
          <div className="flex items-center gap-3">
            <div className="relative shrink-0">
              <img
                src="/assets/icon.webp"
                alt="Evolipia logo"
                className="w-8 h-8 rounded-xl border border-white/10 object-cover shadow-sm"
              />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h1 className="text-base font-black tracking-tight leading-none">Evolipia</h1>
                <span className="px-2 py-0.5 rounded-md bg-emerald-500/15 border border-emerald-500/30 text-emerald-400 font-mono text-[9px] font-black uppercase tracking-widest">
                  Intelligence
                </span>
              </div>
              <span className="text-[10px] text-slate-400 font-medium hidden sm:block">
                Bloomberg Terminal for AI
              </span>
            </div>
          </div>

          {/* Center Search Bar Trigger (Cmd+K) */}
          <button
            onClick={() => setIsCmdPaletteOpen(true)}
            className={`flex-1 max-w-md hidden md:flex items-center justify-between px-3.5 py-2 rounded-xl border transition-all text-xs cursor-pointer ${
              isDarkMode
                ? "bg-slate-900/60 border-white/10 text-slate-400 hover:border-emerald-500/40 hover:bg-slate-900"
                : "bg-slate-100 border-slate-200 text-slate-500 hover:bg-slate-200"
            }`}
          >
            <div className="flex items-center gap-2">
              <Search className="w-4 h-4 text-emerald-400" />
              <span>Search topics, models, papers...</span>
            </div>
            <kbd className="px-2 py-0.5 rounded bg-slate-800 text-slate-300 text-[10px] font-mono border border-slate-700">
              ⌘K
            </kbd>
          </button>

          {/* Right Header Actions */}
          <div className="flex items-center gap-2 sm:gap-3">
            <button
              onClick={() => setShowOnlyBookmarked(!showOnlyBookmarked)}
              className={`px-3 py-2 rounded-xl border text-xs font-bold flex items-center gap-1.5 transition-all cursor-pointer ${
                showOnlyBookmarked
                  ? "bg-emerald-500 border-emerald-400 text-slate-950 shadow-md shadow-emerald-500/20"
                  : isDarkMode
                  ? "bg-slate-900/80 border-white/10 text-slate-300 hover:border-emerald-500/40"
                  : "bg-slate-100 border-slate-200 text-slate-700 hover:bg-slate-200"
              }`}
            >
              <Bookmark className="w-4 h-4" />
              <span className="hidden sm:inline">Bookmarks</span>
              {bookmarkedIds.size > 0 && (
                <span className="px-1.5 py-0.2 rounded-full bg-slate-950/30 text-[10px] font-mono font-black">
                  {bookmarkedIds.size}
                </span>
              )}
            </button>

            <button
              onClick={() => setIsDarkMode(!isDarkMode)}
              className={`p-2.5 rounded-xl border transition-all cursor-pointer ${
                isDarkMode
                  ? "bg-slate-900/80 border-white/10 text-amber-400 hover:bg-slate-900"
                  : "bg-white border-slate-200 text-slate-700 hover:bg-slate-100"
              }`}
              title={isDarkMode ? "Switch to Light Mode" : "Switch to Dark Mode"}
            >
              {isDarkMode ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
            </button>
          </div>
        </div>
      </header>

      {/* Floating Subscribe FAB */}
      <div className="fixed bottom-6 right-6 z-50">
        <button
          onClick={scrollToSubscribe}
          className="px-5 py-3 rounded-full bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm shadow-2xl shadow-emerald-500/40 transition-all flex items-center gap-2 hover:scale-105 active:scale-95 cursor-pointer border border-emerald-300/40"
        >
          <Mail className="w-4 h-4" />
          <span>Subscribe</span>
        </button>
      </div>

      {/* Main Container */}
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 pt-6 pb-20 space-y-8">
        {/* Ticker Marquee */}
        <div className={`relative w-full rounded-2xl border backdrop-blur-md overflow-hidden p-2 flex items-center shadow-lg ${
          isDarkMode ? "border-emerald-500/25 bg-[#09131D]/80" : "border-emerald-500/30 bg-slate-100"
        }`}>
          <div className="shrink-0 flex items-center gap-2 bg-emerald-500 px-3 py-1 rounded-xl text-slate-950 text-[10px] font-black tracking-wider uppercase shadow-md mr-3 z-10 font-mono">
            <span className="w-2 h-2 rounded-full bg-slate-950 animate-ping" />
            LIVE BREAKING
          </div>

          <div className="relative overflow-hidden w-full flex items-center py-0.5">
            <div className={`animate-marquee whitespace-nowrap flex items-center gap-8 text-xs font-semibold ${
              isDarkMode ? "text-slate-200" : "text-slate-800"
            }`}>
              {items.map((item, idx) => (
                <a
                  key={`ticker-1-${item.id || idx}`}
                  href={item.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-2 hover:text-emerald-400 transition-colors group cursor-pointer"
                >
                  <span className="px-2 py-0.5 rounded bg-emerald-500/15 border border-emerald-500/30 text-emerald-400 font-mono text-[10px] font-bold">
                    {item.source_name || item.domain || "GLOBAL"}
                  </span>
                  <span className="font-bold group-hover:underline">{item.title}</span>
                  <span className="text-slate-500 mx-2">•</span>
                </a>
              ))}
            </div>
          </div>
        </div>

        {/* 1. AI Daily Briefing */}
        <AIDailyBriefing isDarkMode={isDarkMode} />

        {/* 2. Breaking Signals */}
        <BreakingSignals
          items={items}
          bookmarkedIds={bookmarkedIds}
          onToggleBookmark={toggleBookmark}
          isDarkMode={isDarkMode}
        />

        {/* 3. AI Pulse (Exploding Keywords Counter) */}
        <AIPulseKeywords
          activeKeyword={filterHook.search}
          onSelectKeyword={handleTopicSelect}
          isDarkMode={isDarkMode}
        />

        {/* 4. NLP Keyword Clusters & Heatmap */}
        <TrendingClusters
          onSelectCluster={handleTopicSelect}
          isDarkMode={isDarkMode}
        />

        {/* 5. Central Filter Bar */}
        <FilterBar filterHook={filterHook} isDarkMode={isDarkMode} />

        {/* 6. Main Editorial Intelligence Stream + Sidebar */}
        <section className="space-y-4">
          <div className="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-[10px] font-black uppercase tracking-[0.2em] text-emerald-400 font-mono">
                Signal Feed
              </p>
              <h3 className={`text-2xl font-black ${isDarkMode ? "text-white" : "text-slate-900"}`}>
                {showOnlyBookmarked ? "Bookmarked Intelligence Signals" : "Latest Intelligence Signals"}
              </h3>
            </div>
            <p className="text-xs text-slate-400 font-mono">
              Showing {displayedItems.length} of {paginationInfo.filteredCount} signals
            </p>
          </div>

          {newsLoading ? (
            <div className={`py-20 text-center flex flex-col items-center gap-3 rounded-[2rem] border ${
              isDarkMode ? "border-white/10 bg-[#09131D] text-slate-400" : "border-slate-200 bg-white text-slate-500"
            }`}>
              <RefreshCw className="w-8 h-8 animate-spin text-emerald-400" />
              <span className="font-semibold text-sm">Loading intelligence stream...</span>
            </div>
          ) : error ? (
            <div className="py-12 bg-rose-950/30 border border-rose-800/40 rounded-[2rem] text-center text-rose-300 p-6">
              <AlertCircle className="w-8 h-8 mx-auto mb-2 text-rose-400" />
              <p className="font-semibold">{error}</p>
            </div>
          ) : displayedItems.length === 0 ? (
            <div className={`py-16 text-center rounded-[2rem] border ${
              isDarkMode ? "bg-slate-900/40 border-slate-800 text-slate-400" : "bg-slate-100 border-slate-200 text-slate-600"
            }`}>
              <Sparkles className="w-8 h-8 mx-auto mb-2 text-slate-400" />
              <p className="font-medium">No intelligence signals match your active filter.</p>
              <button
                onClick={() => {
                  filterHook.resetFilters();
                  setShowOnlyBookmarked(false);
                }}
                className="mt-3 text-xs font-bold text-emerald-400 hover:underline cursor-pointer"
              >
                Reset Filters
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Editorial Feed Grid */}
              <div className="lg:col-span-2 space-y-4">
                {displayedItems.map((item) => {
                  const isBookmarked = bookmarkedIds.has(item.id);
                  return (
                    <article
                      key={item.id}
                      className={`group rounded-[1.75rem] border p-6 shadow-xl transition-all hover:border-emerald-500/40 ${
                        isDarkMode
                          ? "border-white/10 bg-[#09131D]/90 text-slate-100 hover:bg-[#0B1824]"
                          : "border-slate-200 bg-white text-slate-900 hover:shadow-2xl"
                      }`}
                    >
                      <div className="flex items-center justify-between gap-2 mb-3">
                        <div className="flex items-center gap-2">
                          <span className="px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 font-mono text-[10px] font-bold">
                            {item.source_name || item.domain || "GLOBAL"}
                          </span>
                          <span className="text-[11px] text-slate-400 font-mono">
                            {new Date(item.published_at).toLocaleDateString()}
                          </span>
                        </div>

                        {item.relevance_score ? (
                          <span className="px-2 py-0.5 rounded bg-emerald-500 text-slate-950 font-mono text-[10px] font-black">
                            {item.relevance_score}% Impact
                          </span>
                        ) : null}
                      </div>

                      <h4 className={`text-lg sm:text-xl font-black leading-snug tracking-tight mb-2 group-hover:text-emerald-400 transition-colors ${
                        isDarkMode ? "text-white" : "text-slate-900"
                      }`}>
                        <a href={item.url} target="_blank" rel="noopener noreferrer" className="hover:underline">
                          {item.title}
                        </a>
                      </h4>

                      {(item.raw_excerpt || item.summary || item.tldr) && (
                        <p className={`text-xs sm:text-sm leading-relaxed line-clamp-3 mb-4 ${
                          isDarkMode ? "text-slate-300" : "text-slate-600"
                        }`}>
                          {item.raw_excerpt || item.summary || item.tldr}
                        </p>
                      )}

                      {item.why_it_matters && (
                        <div className={`p-3 rounded-xl border text-xs leading-relaxed mb-4 ${
                          isDarkMode ? "bg-emerald-950/20 border-emerald-500/20 text-emerald-200" : "bg-emerald-50 border-emerald-200 text-emerald-900"
                        }`}>
                          <span className="font-bold text-emerald-400 block mb-0.5">Why it matters:</span>
                          {item.why_it_matters}
                        </div>
                      )}

                      {/* Card Footer Quick Actions */}
                      <div className="pt-3 border-t border-white/10 flex items-center justify-between gap-4 text-xs">
                        <div className="flex items-center gap-2">
                          <button
                            onClick={() => toggleBookmark(item.id)}
                            className={`p-1.5 rounded-lg border transition-all cursor-pointer ${
                              isBookmarked
                                ? "bg-emerald-500 border-emerald-400 text-slate-950"
                                : isDarkMode
                                ? "bg-slate-900 border-slate-700 text-slate-400 hover:text-white"
                                : "bg-slate-100 border-slate-200 text-slate-600 hover:bg-slate-200"
                            }`}
                            title={isBookmarked ? "Bookmarked" : "Bookmark"}
                          >
                            <Bookmark className="w-3.5 h-3.5" />
                          </button>

                          <a
                            href={item.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-xs font-bold text-emerald-400 hover:underline inline-flex items-center gap-1"
                          >
                            Read Full <ArrowRight className="w-3.5 h-3.5" />
                          </a>
                        </div>
                      </div>
                    </article>
                  );
                })}
              </div>

              {/* Intelligence Sidebar */}
              <aside className="space-y-4">
                <IntelligenceSidebar
                  topSignals={topSignals}
                  onSelectTopic={handleTopicSelect}
                  isDarkMode={isDarkMode}
                />
              </aside>
            </div>
          )}
        </section>

        {/* 7. Collapsible Developer & Technical Monitoring Section */}
        <section className="pt-6">
          <CrawlProgress
            progressState={progressState}
            isCrawling={isCrawling}
            onStartManualCrawl={startManualCrawl}
            toastMessage={toastMessage}
            onClearToast={clearToast}
            isDarkMode={isDarkMode}
          />
        </section>

        {/* 8. Newsletter CTA */}
        <section id="subscribe-section" className="relative overflow-hidden rounded-[2.5rem] border border-white/10 bg-gradient-to-br from-[#09131D] via-slate-950 to-black p-8 sm:p-12 text-white shadow-2xl mt-12 scroll-mt-28">
          <div className="absolute -right-20 -top-20 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl pointer-events-none" />
          <div className="relative max-w-3xl mx-auto text-center space-y-6">
            <div className="inline-flex items-center gap-2 rounded-full border border-emerald-400/30 bg-emerald-400/10 px-4 py-1.5 text-xs font-bold text-emerald-300">
              <Sparkles className="w-4 h-4" />
              Daily AI Intelligence Briefing
            </div>

            <h3 className="text-3xl sm:text-4xl lg:text-5xl font-black tracking-tight text-white">
              Never miss a breakthrough signal
            </h3>

            <p className="text-sm text-slate-300 max-w-xl mx-auto leading-relaxed">
              Join thousands of engineers, founders, researchers, and AI enthusiasts receiving curated daily briefings.
            </p>

            {subscribed ? (
              <div className="inline-flex items-center gap-2 p-4 rounded-2xl bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 text-sm font-bold shadow-lg">
                <CheckCircle className="w-5 h-5 text-emerald-400" />
                Subscribed! Welcome to Evolipia Radar Daily.
              </div>
            ) : (
              <form onSubmit={handleSubscribe} className="flex flex-col sm:flex-row items-center gap-3 max-w-lg mx-auto pt-2">
                <div className="relative flex-1 w-full">
                  <Mail className="w-4 h-4 absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
                  <input
                    type="email"
                    required
                    placeholder="Enter your work email..."
                    value={emailInput}
                    onChange={(e) => setEmailInput(e.target.value)}
                    className="w-full pl-11 pr-4 h-12 rounded-2xl bg-black/60 border border-white/15 text-sm text-white placeholder-slate-400 focus:outline-none focus:border-emerald-400 shadow-inner"
                  />
                </div>
                <button
                  type="submit"
                  className="w-full sm:w-auto h-12 px-6 rounded-2xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm shadow-xl shadow-emerald-500/25 transition-all shrink-0 active:scale-95 cursor-pointer"
                >
                  Subscribe Free
                </button>
              </form>
            )}
          </div>
        </section>
      </div>

      {/* Footer */}
      <footer className={`border-t py-12 transition-colors ${
        isDarkMode ? "border-white/10 bg-[#050A0F] text-slate-400" : "border-slate-200 bg-slate-100 text-slate-600"
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
                Evolipia Radar
              </span>
              <span className="text-xs text-slate-500 block">AI Intelligence Platform</span>
            </div>
          </div>

          <p className="text-xs text-slate-500">
            © 2026 Evolipia Radar. Built for engineers, founders, researchers &amp; investors.
          </p>
        </div>
      </footer>
    </main>
  );
}
