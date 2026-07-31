"use client";

import React, { useState } from "react";
import { useFilters, FilterPreset } from "../hooks/useFilters";
import {
  Search,
  Calendar,
  Bookmark,
  RotateCcw,
  Plus,
  Trash2,
  Check,
  ChevronDown,
  Filter,
  Sparkles,
  ArrowUpDown,
  X,
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
  { id: "custom", label: "Custom" },
];

const SORT_OPTIONS = [
  { id: "date_desc", label: "Date (Newest)", sortBy: "date", sortOrder: "desc" },
  { id: "date_asc", label: "Date (Oldest)", sortBy: "date", sortOrder: "asc" },
  { id: "relevance_desc", label: "Relevance Score", sortBy: "relevance", sortOrder: "desc" },
  { id: "impact_desc", label: "Impact Score", sortBy: "impact", sortOrder: "desc" },
];

export const FilterBar: React.FC<{ filterHook: ReturnType<typeof useFilters>; isDarkMode?: boolean }> = ({ filterHook, isDarkMode = true }) => {
  const {
    search,
    setSearch,
    dateRange,
    setDateRange,
    dateFrom,
    setDateFrom,
    dateTo,
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
    savedPresets,
    savePreset,
    loadPreset,
    deletePreset,
    resetFilters,
  } = filterHook;

  const [presetNameInput, setPresetNameInput] = useState<string>("");
  const [showSavePreset, setShowSavePreset] = useState<boolean>(false);
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

  const handleSavePreset = (e: React.FormEvent) => {
    e.preventDefault();
    if (presetNameInput.trim()) {
      savePreset(presetNameInput.trim());
      setPresetNameInput("");
      setShowSavePreset(false);
    }
  };

  const handleDateRangeChange = (val: string) => {
    setDateRange(val);
    const now = new Date();
    if (val === "today") {
      setDateFrom(now.toISOString().split("T")[0]);
      setDateTo(now.toISOString().split("T")[0]);
    } else if (val === "7d") {
      const d = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
      setDateFrom(d.toISOString().split("T")[0]);
      setDateTo("");
    } else if (val === "30d") {
      const d = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
      setDateFrom(d.toISOString().split("T")[0]);
      setDateTo("");
    } else if (val === "all") {
      setDateFrom("");
      setDateTo("");
    }
  };

  return (
    <div className={`w-full backdrop-blur-xl rounded-2xl p-5 sm:p-6 mb-8 transition-all space-y-5 border ${
      isDarkMode
        ? "bg-slate-900/90 border-slate-800 text-slate-100 shadow-2xl"
        : "bg-white/95 border-slate-200 text-slate-900 shadow-xl"
    }`}>
      {/* Search & Sort Controls Header */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 sm:gap-4">
        {/* Search Input Box */}
        <div className="relative flex-1 min-w-[260px] w-full group">
          <div className="absolute left-3.5 top-1/2 -translate-y-1/2 flex items-center pointer-events-none text-slate-400 group-focus-within:text-emerald-500 transition-colors z-10">
            <Search className="w-4 h-4" />
          </div>
          <input
            type="text"
            value={search}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearch(e.target.value)}
            placeholder="Search articles, research papers, topics..."
            className={`w-full pl-11 pr-24 h-11 rounded-xl text-sm transition-all focus:outline-none focus:ring-2 focus:ring-emerald-500/30 ${
              isDarkMode
                ? "bg-slate-950/90 border border-slate-800 text-slate-100 placeholder-slate-500 focus:border-emerald-500/80 shadow-inner"
                : "bg-slate-100/90 border border-slate-200 text-slate-900 placeholder-slate-400 focus:border-emerald-500 focus:bg-white shadow-sm"
            }`}
          />
          <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1.5 z-10">
            {search ? (
              <button
                type="button"
                onClick={() => setSearch("")}
                className={`p-1 rounded-lg transition-all ${
                  isDarkMode ? "text-slate-400 hover:text-slate-100 hover:bg-slate-800" : "text-slate-500 hover:text-slate-900 hover:bg-slate-200"
                }`}
                title="Clear search"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            ) : (
              <span className={`hidden sm:inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-mono border ${
                isDarkMode ? "bg-slate-900 border-slate-800 text-slate-500" : "bg-slate-200/80 border-slate-300 text-slate-600"
              }`}>
                Search
              </span>
            )}
          </div>
        </div>

        {/* Sort Selector & Reset */}
        <div className="flex items-center gap-3 shrink-0">
          <div className={`relative flex items-center h-11 border rounded-xl px-3.5 pr-8 transition-all ${
            isDarkMode ? "bg-slate-950/90 border-slate-800 hover:border-slate-700" : "bg-slate-100 border-slate-200 hover:border-slate-300"
          }`}>
            <ArrowUpDown className="w-3.5 h-3.5 text-emerald-500 mr-2 shrink-0" />
            <select
              value={`${sortBy}_${sortOrder}`}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                const opt = SORT_OPTIONS.find(o => o.id === e.target.value);
                if (opt) {
                  setSortBy(opt.sortBy);
                  setSortOrder(opt.sortOrder);
                }
              }}
              className={`bg-transparent text-xs font-bold focus:outline-none cursor-pointer appearance-none truncate pr-2 ${
                isDarkMode ? "text-slate-200" : "text-slate-800"
              }`}
            >
              {SORT_OPTIONS.map(opt => (
                <option key={opt.id} value={opt.id} className={isDarkMode ? "bg-slate-900 text-slate-100 py-1" : "bg-white text-slate-900 py-1"}>
                  {opt.label}
                </option>
              ))}
            </select>
            <ChevronDown className="w-3.5 h-3.5 text-slate-400 pointer-events-none absolute right-3" />
          </div>

          <button
            onClick={resetFilters}
            className={`h-11 px-4 rounded-xl text-xs font-semibold flex items-center gap-2 border shadow-sm transition-all active:scale-95 ${
              isDarkMode
                ? "bg-slate-800/80 hover:bg-slate-700/80 text-slate-200 border-slate-700/60"
                : "bg-slate-200 hover:bg-slate-300/80 text-slate-800 border-slate-300"
            }`}
            title="Reset Filters"
          >
            <RotateCcw className="w-3.5 h-3.5 text-emerald-500" />
            <span>Reset</span>
          </button>
        </div>
      </div>

      {/* Main Filter Section (3 Spread-out columns without Verification Status) */}
      <div className={`grid grid-cols-1 md:grid-cols-3 gap-6 pt-4 border-t ${
        isDarkMode ? "border-slate-800/70" : "border-slate-200"
      }`}>
        
        {/* Radio Box Segment: Date Range */}
        <div className="space-y-2">
          <label className={`text-xs font-bold flex items-center gap-1.5 ${isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
            <Calendar className="w-3.5 h-3.5 text-emerald-500" /> Timeframe
          </label>
          <div className={`flex items-center gap-1 p-1 border rounded-xl overflow-x-auto scrollbar-hide ${
            isDarkMode ? "bg-slate-950/80 border-slate-800" : "bg-slate-100 border-slate-200"
          }`}>
            {DATE_RANGE_OPTIONS.map(opt => {
              const isActive = dateRange === opt.id;
              return (
                <button
                  key={opt.id}
                  type="button"
                  onClick={() => handleDateRangeChange(opt.id)}
                  className={`px-2.5 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-all flex-1 text-center ${
                    isActive
                      ? "bg-emerald-500 text-slate-950 font-black shadow-md shadow-emerald-500/20"
                      : isDarkMode
                      ? "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
                      : "text-slate-600 hover:text-slate-900 hover:bg-slate-200"
                  }`}
                >
                  {opt.label}
                </button>
              );
            })}
          </div>

          {dateRange === "custom" && (
            <div className="flex gap-2 mt-2">
              <input
                type="date"
                value={dateFrom}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDateFrom(e.target.value)}
                className={`w-1/2 border text-xs p-2 rounded-lg focus:outline-none focus:border-emerald-500 ${
                  isDarkMode ? "bg-slate-950 border-slate-800 text-slate-200" : "bg-slate-50 border-slate-300 text-slate-900"
                }`}
              />
              <input
                type="date"
                value={dateTo}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDateTo(e.target.value)}
                className={`w-1/2 border text-xs p-2 rounded-lg focus:outline-none focus:border-emerald-500 ${
                  isDarkMode ? "bg-slate-950 border-slate-800 text-slate-200" : "bg-slate-50 border-slate-300 text-slate-900"
                }`}
              />
            </div>
          )}
        </div>

        {/* Source Dropdown Selector */}
        <div className="space-y-2 relative">
          <label className={`text-xs font-bold flex items-center gap-1.5 ${isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
            <Filter className="w-3.5 h-3.5 text-emerald-500" /> Data Sources
          </label>
          <button
            type="button"
            onClick={() => setShowSourceDropdown(!showSourceDropdown)}
            className={`w-full border text-xs font-medium p-2.5 rounded-xl flex items-center justify-between focus:outline-none focus:border-emerald-500 transition-all ${
              isDarkMode ? "bg-slate-950/80 border-slate-800 text-slate-200 hover:bg-slate-900" : "bg-slate-100 border-slate-200 text-slate-800 hover:bg-slate-200/80"
            }`}
          >
            <span className="truncate">
              {selectedSources.length === 0
                ? "All Sources Enabled"
                : `${selectedSources.length} source${selectedSources.length > 1 ? "s" : ""} active`}
            </span>
            <ChevronDown className={`w-3.5 h-3.5 text-slate-400 transition-transform ${showSourceDropdown ? "rotate-180" : ""}`} />
          </button>

          {showSourceDropdown && (
            <div className={`absolute left-0 min-w-[240px] w-full top-full mt-2 backdrop-blur-2xl rounded-xl p-2.5 z-50 shadow-2xl max-h-64 overflow-y-auto space-y-1 border ${
              isDarkMode ? "bg-slate-950/98 border-slate-700 text-slate-200" : "bg-white border-slate-300 text-slate-900"
            }`}>
              <div className={`flex items-center justify-between px-2 py-1.5 text-[11px] font-bold border-b mb-1 ${
                isDarkMode ? "text-slate-400 border-slate-800" : "text-slate-600 border-slate-200"
              }`}>
                <span>Select Sources</span>
                {selectedSources.length > 0 ? (
                  <button
                    type="button"
                    onClick={() => setSelectedSources([])}
                    className="text-emerald-500 hover:text-emerald-400 font-semibold"
                  >
                    Reset (All)
                  </button>
                ) : (
                  <span className="text-[10px] text-slate-500 font-normal">All active</span>
                )}
              </div>
              {AVAILABLE_SOURCES.map(src => {
                const checked = selectedSources.includes(src);
                return (
                  <label
                    key={src}
                    className={`flex items-center justify-between text-xs px-2.5 py-2 rounded-lg cursor-pointer transition-all ${
                      checked
                        ? isDarkMode
                          ? "bg-emerald-950/60 text-emerald-300 font-semibold border border-emerald-800/40"
                          : "bg-emerald-50 text-emerald-800 font-semibold border border-emerald-200"
                        : isDarkMode
                        ? "text-slate-300 hover:bg-slate-900 border border-transparent"
                        : "text-slate-700 hover:bg-slate-100 border border-transparent"
                    }`}
                  >
                    <div className="flex items-center gap-2.5">
                      <input
                        type="checkbox"
                        checked={checked}
                        onChange={() => toggleSource(src)}
                        className="w-4 h-4 rounded border-slate-400 bg-slate-100 text-emerald-500 focus:ring-0 cursor-pointer accent-emerald-500"
                      />
                      <span>{src}</span>
                    </div>
                    {checked && <Check className="w-3.5 h-3.5 text-emerald-500 shrink-0" />}
                  </label>
                );
              })}
            </div>
          )}
        </div>

        {/* Min Relevance Slider with Green Accent Handle */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <label className={`text-xs font-bold flex items-center gap-1.5 ${isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
              <Sparkles className="w-3.5 h-3.5 text-emerald-500" /> Min Relevance
            </label>
            <span className={`text-xs font-mono font-bold px-2 py-0.5 rounded-md border ${
              isDarkMode ? "text-emerald-300 bg-emerald-950/60 border-emerald-800/40" : "text-emerald-800 bg-emerald-50 border-emerald-200"
            }`}>
              {minRelevance}%
            </span>
          </div>
          <input
            type="range"
            min="0"
            max="100"
            step="5"
            value={minRelevance}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setMinRelevance(Number(e.target.value))}
            className={`w-full accent-emerald-500 h-2 rounded-lg cursor-pointer ${
              isDarkMode ? "bg-slate-950" : "bg-slate-200"
            }`}
          />
          <div className={`flex justify-between text-[10px] font-mono ${isDarkMode ? "text-slate-500" : "text-slate-400"}`}>
            <button onClick={() => setMinRelevance(0)} className={isDarkMode ? "hover:text-slate-300" : "hover:text-slate-800"}>0% (All)</button>
            <button onClick={() => setMinRelevance(50)} className={isDarkMode ? "hover:text-slate-300" : "hover:text-slate-800"}>50%</button>
            <button onClick={() => setMinRelevance(80)} className={isDarkMode ? "hover:text-slate-300" : "hover:text-slate-800"}>80% (High)</button>
          </div>
        </div>

      </div>

      {/* Category Pills */}
      <div className={`pt-3 border-t flex flex-wrap items-center gap-2 ${
        isDarkMode ? "border-slate-800/60" : "border-slate-200"
      }`}>
        <span className={`text-xs font-bold mr-1 ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>Categories:</span>
        {AVAILABLE_CATEGORIES.map(cat => {
          const isSelected = selectedCategories.includes(cat);
          return (
            <button
              key={cat}
              onClick={() => toggleCategory(cat)}
              className={`px-3 py-1 rounded-full text-xs font-semibold transition-all flex items-center gap-1.5 ${
                isSelected
                  ? "bg-emerald-500 text-slate-950 font-black shadow-md shadow-emerald-500/20"
                  : isDarkMode
                  ? "bg-slate-950 text-slate-400 hover:text-slate-200 border border-slate-800 hover:border-slate-700"
                  : "bg-slate-100 text-slate-600 hover:text-slate-900 border border-slate-200 hover:border-slate-300"
              }`}
            >
              {isSelected && <Check className="w-3 h-3" />}
              <span className="capitalize">{cat}</span>
            </button>
          );
        })}
      </div>

      {/* Presets Management Bar */}
      <div className={`pt-3 border-t flex flex-wrap items-center justify-between gap-3 text-xs ${
        isDarkMode ? "border-slate-800/60 text-slate-400" : "border-slate-200 text-slate-500"
      }`}>
        <div className="flex flex-wrap items-center gap-2">
          <span className="font-bold flex items-center gap-1">
            <Bookmark className="w-3.5 h-3.5 text-emerald-500" /> Presets:
          </span>
          {savedPresets.length === 0 ? (
            <span className="italic text-[11px]">No saved presets yet</span>
          ) : (
            savedPresets.map((preset: FilterPreset) => (
              <div
                key={preset.id}
                className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border text-[11px] font-semibold transition-all ${
                  isDarkMode ? "bg-slate-950 border-slate-800 text-slate-300" : "bg-slate-100 border-slate-200 text-slate-700"
                }`}
              >
                <button
                  type="button"
                  onClick={() => loadPreset(preset)}
                  className="hover:text-emerald-500 transition-colors"
                >
                  {preset.name}
                </button>
                <button
                  type="button"
                  onClick={() => deletePreset(preset.id)}
                  className="text-slate-500 hover:text-rose-400 transition-colors"
                  title="Delete preset"
                >
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
            ))
          )}
        </div>

        {showSavePreset ? (
          <form onSubmit={handleSavePreset} className="flex items-center gap-2">
            <input
              type="text"
              placeholder="Preset name..."
              value={presetNameInput}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPresetNameInput(e.target.value)}
              className={`px-2.5 py-1 text-xs rounded-lg border focus:outline-none focus:border-emerald-500 ${
                isDarkMode ? "bg-slate-950 border-slate-800 text-slate-200" : "bg-white border-slate-300 text-slate-900"
              }`}
              autoFocus
            />
            <button
              type="submit"
              className="px-2.5 py-1 rounded-lg bg-emerald-500 text-slate-950 font-bold text-xs hover:bg-emerald-400"
            >
              Save
            </button>
            <button
              type="button"
              onClick={() => setShowSavePreset(false)}
              className="text-slate-400 hover:text-slate-200 text-xs"
            >
              Cancel
            </button>
          </form>
        ) : (
          <button
            type="button"
            onClick={() => setShowSavePreset(true)}
            className="inline-flex items-center gap-1 text-xs font-semibold text-emerald-500 hover:text-emerald-400"
          >
            <Plus className="w-3.5 h-3.5" /> Save Current Preset
          </button>
        )}
      </div>
    </div>
  );
};
