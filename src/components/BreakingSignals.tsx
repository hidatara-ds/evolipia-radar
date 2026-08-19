"use client";

import React from "react";
import { NewsItem } from "@/src/api/client";
import { Zap, ShieldCheck, Clock, ExternalLink, Bookmark, Share2, Flame, MessageSquare, TrendingUp } from "lucide-react";

interface BreakingSignalsProps {
  items: NewsItem[];
  bookmarkedIds: Set<string | number>;
  onToggleBookmark: (id: string | number) => void;
  isDarkMode?: boolean;
}

export const BreakingSignals: React.FC<BreakingSignalsProps> = ({
  items,
  bookmarkedIds,
  onToggleBookmark,
  isDarkMode = true,
}) => {
  const breakingList = items.length >= 2 ? items.slice(0, 2) : items;

  if (breakingList.length === 0) return null;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="h-2.5 w-2.5 rounded-full bg-emerald-400 animate-ping" />
          <h3 className="text-xs font-black uppercase tracking-[0.25em] text-emerald-400">Breaking Signals</h3>
        </div>
        <span className={`text-xs font-mono font-medium ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
          Highest Impact Signals Today
        </span>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        {breakingList.map((item, index) => {
          const isBookmarked = bookmarkedIds.has(item.id);
          const impactScore = item.relevance_score || (95 - index * 4);
          const whyItMatters =
            item.summary ||
            "This signal directly influences architectural decisions, enterprise LLM adoption, and upcoming engineering workflow patterns.";

          return (
            <div
              key={item.id}
              className={`relative overflow-hidden rounded-[2rem] border p-6 sm:p-7 shadow-2xl transition-all duration-300 flex flex-col justify-between group hover:border-emerald-500/50 ${
                isDarkMode
                  ? "border-emerald-500/25 bg-gradient-to-br from-[#0B1A24] via-[#08121A] to-[#050A0F] text-white"
                  : "border-emerald-500/30 bg-white text-slate-900 shadow-xl"
              }`}
            >
              {/* Subtle green ambient accent */}
              <div className="absolute top-0 right-0 h-40 w-40 bg-emerald-500/10 blur-3xl pointer-events-none group-hover:bg-emerald-500/20 transition-all" />

              <div>
                {/* Meta Top Header */}
                <div className="flex flex-wrap items-center justify-between gap-2 mb-4">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-full bg-emerald-500/15 border border-emerald-500/30 text-emerald-400 font-mono text-[10px] font-black uppercase tracking-wider">
                      {item.source_name || item.domain || "GLOBAL INTEL"}
                    </span>
                    <span className={`text-[11px] font-mono font-medium ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
                      {item.published_at ? new Date(item.published_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : "Just now"}
                    </span>
                  </div>

                  {/* Impact Score Pill */}
                  <div className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-500 text-slate-950 font-black text-xs shadow-md shadow-emerald-500/20 font-mono">
                    <Zap className="w-3.5 h-3.5 fill-current" />
                    <span>{impactScore}% Impact Score</span>
                  </div>
                </div>

                {/* Main Headline */}
                <h3 className={`text-xl sm:text-2xl font-black leading-tight tracking-tight mb-3 group-hover:text-emerald-400 transition-colors ${
                  isDarkMode ? "text-white" : "text-slate-900"
                }`}>
                  <a href={item.url} target="_blank" rel="noopener noreferrer" className="hover:underline">
                    {item.title}
                  </a>
                </h3>

                {/* Why It Matters Callout Box */}
                <div className={`p-4 rounded-2xl border mb-5 ${
                  isDarkMode
                    ? "bg-slate-950/70 border-white/10 text-slate-300"
                    : "bg-slate-50 border-slate-200 text-slate-700"
                }`}>
                  <div className="flex items-center gap-1.5 text-[10px] font-black uppercase tracking-wider text-emerald-400 mb-1 font-mono">
                    <TrendingUp className="w-3.5 h-3.5" /> Why It Matters
                  </div>
                  <p className="text-xs sm:text-sm font-medium leading-relaxed">
                    {whyItMatters}
                  </p>
                </div>
              </div>

              {/* Card Footer & Quick Actions */}
              <div className="pt-4 border-t border-white/10 flex items-center justify-between gap-4 text-xs">
                <div className="flex items-center gap-4 text-slate-400 font-mono text-[11px]">
                  <span className="flex items-center gap-1">
                    <Clock className="w-3.5 h-3.5 text-emerald-400" /> 3 min read
                  </span>
                  <span className="hidden sm:flex items-center gap-1">
                    <MessageSquare className="w-3.5 h-3.5 text-emerald-400" /> 1.4k discussions
                  </span>
                  <span className="flex items-center gap-1 text-emerald-400 font-semibold">
                    <ShieldCheck className="w-3.5 h-3.5" /> High Trust
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    onClick={() => onToggleBookmark(item.id)}
                    className={`p-2 rounded-xl border transition-all cursor-pointer ${
                      isBookmarked
                        ? "bg-emerald-500 border-emerald-400 text-slate-950"
                        : isDarkMode
                        ? "bg-slate-900 border-slate-700 text-slate-300 hover:text-white hover:border-emerald-500/50"
                        : "bg-slate-100 border-slate-200 text-slate-700 hover:bg-slate-200"
                    }`}
                    title={isBookmarked ? "Remove Bookmark" : "Bookmark Signal"}
                  >
                    <Bookmark className="w-4 h-4" />
                  </button>

                  <a
                    href={item.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="px-3 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-xs shadow-md transition-all flex items-center gap-1.5 cursor-pointer"
                  >
                    <span>Read Signal</span>
                    <ExternalLink className="w-3.5 h-3.5" />
                  </a>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
