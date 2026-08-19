"use client";

import React, { useState } from "react";
import { useFilters } from "../hooks/useFilters";
import {
  Search,
  RotateCcw,
  ChevronDown,
  Filter,
  ArrowUpDown,
  Sparkles,
} from "lucide-react";

const AVAILABLE_SOURCES = [
  "Hacker News",
  "ArXiv AI",
  "Reddit MachineLearning",
  "TechCrunch AI",
  "GitHub Trending",
  "Twitter / X",
];

const AVAILABLE_CATEGORIES = [
  "llm",
  "agents",
  "vision",
  "open-source",
  "infra",
  "robotics",
  "security",
];

const DATE_RANGE_OPTIONS = [
  { id: "all", label: "All Time" },
  { id: "today", label: "Today" },
  { id: "7d", label: "7 Days" },
  { id: "30d", label: "30 Days" },
];

const SORT_OPTIONS = [
  { id: "date_desc", label: "Newest Date", sortBy: "date", sortOrder: "desc" },
  { id: "date_asc", label: "Oldest Date", sortBy: "date", sortOrder: "asc" },
  { id: "relevance_desc", label: "Highest Relevance", sortBy: "relevance", sortOrder: "desc" },
  { id: "impact_desc", label: "Highest Impact", sortBy: "impact", sortOrder: "desc" },
];

export const FilterBar: React.FC<{ filterHook: ReturnType<typeof useFilters>; isDarkMode?: boolean }> = ({
  filterHook,
  isDarkMode = true,
}) => {
  const {
    search,
    setSearch,
    dateRange,
    setDateRange,
    setDateFrom,
    setDateTo,
    selectedSources,
    setSelectedSources,
    selectedCategories,
    setSelectedCategories,
    minRelevance,
    setMinRelevance,
    sortBy,
    setSortBy,
    sortOrder,
    setSortOrder,
    resetFilters,
  } = filterHook;

  const [showSourceDropdown, setShowSourceDropdown] = useState<boolean>(false);

  const toggleSource = (source: string) => {
    if (selectedSources.includes(source)) {
      setSelectedSources(selectedSources.filter((s: string) => s !== source));
    } else {
      setSelectedSources([...selectedSources, source]);
    }
  };

  const toggleCategory = (cat: string) => {
    if (selectedCategories.includes(cat)) {
      setSelectedCategories(selectedCategories.filter((c: string) => c !== cat));
    } else {
      setSelectedCategories([...selectedCategories, cat]);
    }
  };

  const activeSortId = `${sortBy}_${sortOrder}`;

  return (
    <div
      className={`rounded-2xl border p-4 transition-all shadow-md space-y-3.5 ${
        isDarkMode
          ? "border-white/10 bg-[#09131D]/90 text-slate-100"
          : "border-slate-200 bg-white text-slate-900"
      }`}
    >
      {/* Top Search & Primary Toolbar Row */}
      <div className="flex flex-col sm:flex-row items-center gap-2.5">
        {/* Search Bar Input */}
        <div className="relative flex-1 w-full">
          <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-emerald-400 pointer-events-none" />
          <input
            type="text"
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              filterHook.setPage(1);
            }}
            placeholder="Search articles, research papers, topics..."
            className={`w-full pl-10 pr-4 h-10 rounded-xl text-xs sm:text-sm font-medium outline-none border transition-all ${
              isDarkMode
                ? "bg-slate-900/80 border-white/10 text-white placeholder-slate-400 focus:border-emerald-400 focus:ring-1 focus:ring-emerald-400/30"
                : "bg-slate-50 border-slate-200 text-slate-900 placeholder-slate-400 focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500/20"
            }`}
          />
        </div>

        {/* Sort Select & Reset Button */}
        <div className="flex items-center gap-2 w-full sm:w-auto shrink-0">
          <div className="relative flex-1 sm:flex-none">
            <select
              value={activeSortId}
              onChange={(e) => {
                const opt = SORT_OPTIONS.find((s) => s.id === e.target.value);
                if (opt) {
                  setSortBy(opt.sortBy);
                  setSortOrder(opt.sortOrder);
                  filterHook.setPage(1);
                }
              }}
              className={`w-full h-10 pl-3 pr-8 rounded-xl text-xs font-bold appearance-none outline-none border cursor-pointer ${
                isDarkMode
                  ? "bg-slate-900/80 border-white/10 text-slate-200 hover:border-emerald-400/40"
                  : "bg-slate-50 border-slate-200 text-slate-800 hover:border-slate-300"
              }`}
            >
              {SORT_OPTIONS.map((opt) => (
                <option key={opt.id} value={opt.id} className={isDarkMode ? "bg-slate-900 text-white" : "bg-white text-slate-900"}>
                  Sort: {opt.label}
                </option>
              ))}
            </select>
            <ArrowUpDown className="w-3.5 h-3.5 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400" />
          </div>

          <button
            onClick={resetFilters}
            className={`h-10 px-3.5 rounded-xl border text-xs font-bold flex items-center gap-1.5 transition-all shrink-0 cursor-pointer ${
              isDarkMode
                ? "bg-slate-900/80 border-white/10 text-slate-400 hover:text-emerald-400 hover:border-emerald-400/40"
                : "bg-slate-50 border-slate-200 text-slate-600 hover:text-emerald-600 hover:border-emerald-500/40"
            }`}
            title="Reset Filters"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            <span>Reset</span>
          </button>
        </div>
      </div>

      {/* Second Row: Timeframe Pills & Source Dropdown */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-2 border-t border-white/10">
        {/* Timeframe Pills */}
        <div className="flex items-center gap-1 w-full sm:w-auto overflow-x-auto scrollbar-hide py-0.5">
          <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-slate-400 mr-1.5 shrink-0">
            Timeframe:
          </span>
          {DATE_RANGE_OPTIONS.map((opt) => {
            const isActive = dateRange === opt.id;
            return (
              <button
                key={opt.id}
                onClick={() => {
                  setDateRange(opt.id);
                  setDateFrom("");
                  setDateTo("");
                  filterHook.setPage(1);
                }}
                className={`px-3 py-1.5 rounded-xl text-xs font-bold transition-all shrink-0 cursor-pointer ${
                  isActive
                    ? "bg-emerald-500 text-slate-950 shadow-md font-black"
                    : isDarkMode
                    ? "bg-slate-900/50 border border-white/10 text-slate-400 hover:text-white hover:bg-slate-800"
                    : "bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 hover:bg-slate-200"
                }`}
              >
                {opt.label}
              </button>
            );
          })}
        </div>

        {/* Multi-Source Selector Dropdown */}
        <div className="relative w-full sm:w-auto shrink-0">
          <button
            onClick={() => setShowSourceDropdown(!showSourceDropdown)}
            className={`w-full sm:w-auto px-3.5 py-1.5 rounded-xl border text-xs font-bold flex items-center justify-between gap-2 cursor-pointer transition-all ${
              selectedSources.length > 0
                ? "border-emerald-500/50 bg-emerald-500/10 text-emerald-400"
                : isDarkMode
                ? "border-white/10 bg-slate-900/50 text-slate-300 hover:border-slate-700"
                : "border-slate-200 bg-slate-100 text-slate-700 hover:bg-slate-200"
            }`}
          >
            <div className="flex items-center gap-1.5">
              <Filter className="w-3.5 h-3.5 text-emerald-400" />
              <span>
                {selectedSources.length === 0
                  ? "All Data Sources"
                  : `Sources (${selectedSources.length})`}
              </span>
            </div>
            <ChevronDown className="w-3.5 h-3.5 opacity-70" />
          </button>

          {showSourceDropdown && (
            <div
              className={`absolute right-0 top-full mt-2 w-56 rounded-xl border p-2 shadow-2xl z-40 transition-all ${
                isDarkMode
                  ? "bg-[#09131D] border-white/15 text-white"
                  : "bg-white border-slate-200 text-slate-900 shadow-xl"
              }`}
            >
              <div className="text-[10px] font-mono font-bold uppercase text-slate-400 px-2 py-1 border-b border-white/10 mb-1">
                Filter by Source
              </div>
              <div className="space-y-1 max-h-48 overflow-y-auto scrollbar-hide">
                {AVAILABLE_SOURCES.map((src) => {
                  const isChecked = selectedSources.includes(src);
                  return (
                    <button
                      key={src}
                      onClick={() => {
                        toggleSource(src);
                        filterHook.setPage(1);
                      }}
                      className={`w-full px-2.5 py-1.5 rounded-lg text-xs font-medium text-left flex items-center justify-between transition-colors cursor-pointer ${
                        isChecked
                          ? "bg-emerald-500/20 text-emerald-300 font-bold"
                          : isDarkMode
                          ? "hover:bg-slate-800 text-slate-300"
                          : "hover:bg-slate-100 text-slate-700"
                      }`}
                    >
                      <span>{src}</span>
                      {isChecked && <span className="text-emerald-400 font-black">✓</span>}
                    </button>
                  );
                })}
              </div>
              {selectedSources.length > 0 && (
                <button
                  onClick={() => {
                    setSelectedSources([]);
                    filterHook.setPage(1);
                  }}
                  className="w-full text-center text-[10px] font-mono font-bold text-rose-400 hover:underline pt-1.5 mt-1 border-t border-white/10"
                >
                  Clear Source Filters
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Third Row: Min Relevance Slider & Category Tag Chips */}
      <div className="pt-2 border-t border-white/10 space-y-2">
        {/* Min Relevance Slider */}
        <div className="flex items-center justify-between text-xs font-mono">
          <span className="text-slate-400 flex items-center gap-1">
            <Sparkles className="w-3.5 h-3.5 text-emerald-400" />
            Min Relevance Threshold:
          </span>
          <span className="font-bold text-emerald-400 font-mono px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20">
            {minRelevance}%
          </span>
        </div>

        <input
          type="range"
          min={0}
          max={90}
          step={5}
          value={minRelevance}
          onChange={(e) => {
            setMinRelevance(Number(e.target.value));
            filterHook.setPage(1);
          }}
          className="w-full h-1.5 rounded-lg bg-slate-800 appearance-none cursor-pointer accent-emerald-500"
        />

        {/* Category Tag Chips */}
        <div className="flex flex-wrap items-center gap-1.5 pt-1">
          <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-slate-400 mr-1">
            Category:
          </span>
          {AVAILABLE_CATEGORIES.map((cat) => {
            const isCatActive = selectedCategories.includes(cat);
            return (
              <button
                key={cat}
                onClick={() => {
                  toggleCategory(cat);
                  filterHook.setPage(1);
                }}
                className={`px-2.5 py-1 rounded-lg text-[11px] font-mono font-bold transition-all cursor-pointer border ${
                  isCatActive
                    ? "bg-emerald-500 border-emerald-400 text-slate-950 shadow-sm"
                    : isDarkMode
                    ? "bg-slate-900/60 border-white/10 text-slate-400 hover:border-emerald-500/40 hover:text-slate-200"
                    : "bg-slate-100 border-slate-200 text-slate-600 hover:border-emerald-500/40 hover:text-slate-900"
                }`}
              >
                #{cat}
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
};
