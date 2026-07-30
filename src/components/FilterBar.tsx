"use client";

import React, { useState } from "react";
import { useFilters, FilterPreset } from "../hooks/useFilters";
import {
  Search,
  Calendar,
  SlidersHorizontal,
  Bookmark,
  RotateCcw,
  Plus,
  Trash2,
  Check,
  ChevronDown,
  Filter,
  Sparkles,
  ArrowUpDown,
  Layers,
  CheckCircle2,
  Clock,
  AlertTriangle,
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

const STATUS_OPTIONS = [
  { id: "all", label: "All Items", icon: Layers },
  { id: "verified", label: "Verified", icon: CheckCircle2 },
  { id: "pending", label: "Pending", icon: Clock },
  { id: "failed", label: "Failed", icon: AlertTriangle },
];

const SORT_OPTIONS = [
  { id: "date_desc", label: "Date (Newest)", sortBy: "date", sortOrder: "desc" },
  { id: "date_asc", label: "Date (Oldest)", sortBy: "date", sortOrder: "asc" },
  { id: "relevance_desc", label: "Relevance Score", sortBy: "relevance", sortOrder: "desc" },
  { id: "impact_desc", label: "Impact Score", sortBy: "impact", sortOrder: "desc" },
];

export const FilterBar: React.FC<{ filterHook: ReturnType<typeof useFilters> }> = ({ filterHook }) => {
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
    status,
    setStatus,
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
    <div className="w-full bg-slate-900/90 border border-slate-800 backdrop-blur-xl rounded-2xl p-5 sm:p-6 shadow-2xl mb-8 text-slate-100 space-y-5 transition-all">
      {/* Search & Sort Controls Header */}
      <div className="flex flex-col lg:flex-row items-stretch lg:items-center justify-between gap-4">
        {/* Search Input Box */}
        <div className="relative flex-1 min-w-[260px] w-full">
          <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          <input
            type="text"
            value={search}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearch(e.target.value)}
            placeholder="Search articles, research papers, topics..."
            className="w-full pl-10 pr-10 py-2.5 bg-slate-950/80 border border-slate-800/90 rounded-xl text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/30 transition-all"
          />
          {search && (
            <button
              onClick={() => setSearch("")}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300 text-xs px-1 py-0.5 rounded bg-slate-900"
            >
              Clear
            </button>
          )}
        </div>

        {/* Sort Selector & Reset */}
        <div className="flex items-center gap-3 shrink-0">
          <div className="relative flex items-center bg-slate-950/80 border border-slate-800/90 rounded-xl px-3 py-2">
            <ArrowUpDown className="w-3.5 h-3.5 text-indigo-400 mr-2 shrink-0" />
            <select
              value={`${sortBy}_${sortOrder}`}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                const opt = SORT_OPTIONS.find(o => o.id === e.target.value);
                if (opt) {
                  setSortBy(opt.sortBy);
                  setSortOrder(opt.sortOrder);
                }
              }}
              className="bg-transparent text-xs font-semibold text-slate-200 focus:outline-none cursor-pointer pr-4 appearance-none"
            >
              {SORT_OPTIONS.map(opt => (
                <option key={opt.id} value={opt.id} className="bg-slate-900 text-slate-100 py-1">
                  {opt.label}
                </option>
              ))}
            </select>
            <ChevronDown className="w-3 h-3 text-slate-400 pointer-events-none absolute right-2.5" />
          </div>

          <button
            onClick={resetFilters}
            className="px-3 py-2 rounded-xl bg-slate-800/80 hover:bg-slate-700/80 text-slate-300 text-xs font-semibold flex items-center gap-1.5 border border-slate-700/50 transition-all active:scale-95"
            title="Reset Filters"
          >
            <RotateCcw className="w-3.5 h-3.5 text-slate-400" />
            <span>Reset</span>
          </button>
        </div>
      </div>

      {/* Main Filter Section */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5 pt-4 border-t border-slate-800/70">
        
        {/* Radio Box Segment: Date Range */}
        <div className="space-y-2">
          <label className="text-xs font-bold text-slate-300 flex items-center gap-1.5">
            <Calendar className="w-3.5 h-3.5 text-indigo-400" /> Timeframe
          </label>
          <div className="flex items-center gap-1 p-1 bg-slate-950/80 border border-slate-800 rounded-xl overflow-x-auto scrollbar-hide">
            {DATE_RANGE_OPTIONS.map(opt => {
              const isActive = dateRange === opt.id;
              return (
                <button
                  key={opt.id}
                  type="button"
                  onClick={() => handleDateRangeChange(opt.id)}
                  className={`px-2.5 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-all flex-1 text-center ${
                    isActive
                      ? "bg-indigo-600 text-white shadow-sm shadow-indigo-500/30"
                      : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
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
                className="w-1/2 bg-slate-950 border border-slate-800 text-xs text-slate-200 p-2 rounded-lg focus:outline-none focus:border-indigo-500"
              />
              <input
                type="date"
                value={dateTo}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDateTo(e.target.value)}
                className="w-1/2 bg-slate-950 border border-slate-800 text-xs text-slate-200 p-2 rounded-lg focus:outline-none focus:border-indigo-500"
              />
            </div>
          )}
        </div>

        {/* Source Dropdown Selector */}
        <div className="space-y-2 relative">
          <label className="text-xs font-bold text-slate-300 flex items-center gap-1.5">
            <Filter className="w-3.5 h-3.5 text-indigo-400" /> Data Sources
          </label>
          <button
            type="button"
            onClick={() => setShowSourceDropdown(!showSourceDropdown)}
            className="w-full bg-slate-950/80 border border-slate-800 text-xs font-medium text-slate-200 p-2.5 rounded-xl flex items-center justify-between focus:outline-none focus:border-indigo-500 transition-all hover:bg-slate-900"
          >
            <span className="truncate">
              {selectedSources.length === 0
                ? "All Sources Enabled"
                : `${selectedSources.length} source${selectedSources.length > 1 ? "s" : ""} active`}
            </span>
            <ChevronDown className={`w-3.5 h-3.5 text-slate-400 transition-transform ${showSourceDropdown ? "rotate-180" : ""}`} />
          </button>

          {showSourceDropdown && (
            <div className="absolute left-0 right-0 top-full mt-2 bg-slate-950 border border-slate-800 rounded-xl p-2.5 z-30 shadow-2xl max-h-56 overflow-y-auto space-y-1">
              <div className="flex items-center justify-between px-2 py-1 text-[11px] font-semibold text-slate-400 border-b border-slate-800 mb-1">
                <span>Select Sources</span>
                {selectedSources.length > 0 && (
                  <button
                    onClick={() => setSelectedSources([])}
                    className="text-indigo-400 hover:underline"
                  >
                    Reset
                  </button>
                )}
              </div>
              {AVAILABLE_SOURCES.map(src => {
                const checked = selectedSources.includes(src);
                return (
                  <label
                    key={src}
                    className={`flex items-center justify-between text-xs px-2.5 py-1.5 rounded-lg cursor-pointer transition-all ${
                      checked ? "bg-indigo-950/60 text-indigo-200 font-semibold" : "text-slate-300 hover:bg-slate-900"
                    }`}
                  >
                    <div className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={checked}
                        onChange={() => toggleSource(src)}
                        className="rounded border-slate-700 bg-slate-900 text-indigo-500 focus:ring-0 cursor-pointer"
                      />
                      <span>{src}</span>
                    </div>
                    {checked && <Check className="w-3.5 h-3.5 text-indigo-400" />}
                  </label>
                );
              })}
            </div>
          )}
        </div>

        {/* Min Relevance Slider */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <label className="text-xs font-bold text-slate-300 flex items-center gap-1.5">
              <Sparkles className="w-3.5 h-3.5 text-amber-400" /> Min Relevance
            </label>
            <span className="text-xs font-mono font-bold text-amber-300 bg-amber-950/60 border border-amber-800/40 px-2 py-0.5 rounded-md">
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
            className="w-full accent-indigo-500 bg-slate-950 h-2 rounded-lg cursor-pointer"
          />
          <div className="flex justify-between text-[10px] text-slate-500 font-mono">
            <button onClick={() => setMinRelevance(0)} className="hover:text-slate-300">0% (All)</button>
            <button onClick={() => setMinRelevance(50)} className="hover:text-slate-300">50%</button>
            <button onClick={() => setMinRelevance(80)} className="hover:text-slate-300">80% (High)</button>
          </div>
        </div>

        {/* Radio Box Segment: Verification Status */}
        <div className="space-y-2">
          <label className="text-xs font-bold text-slate-300 flex items-center gap-1.5">
            <SlidersHorizontal className="w-3.5 h-3.5 text-indigo-400" /> Verification Status
          </label>
          <div className="grid grid-cols-2 gap-1 p-1 bg-slate-950/80 border border-slate-800 rounded-xl">
            {STATUS_OPTIONS.map(opt => {
              const isActive = status === opt.id;
              const IconComponent = opt.icon;
              return (
                <button
                  key={opt.id}
                  type="button"
                  onClick={() => setStatus(opt.id)}
                  className={`px-2 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center justify-center gap-1.5 ${
                    isActive
                      ? "bg-indigo-600 text-white shadow-sm shadow-indigo-500/30"
                      : "text-slate-400 hover:text-slate-200 hover:bg-slate-900"
                  }`}
                >
                  <IconComponent className="w-3 h-3" />
                  <span className="truncate">{opt.label}</span>
                </button>
              );
            })}
          </div>
        </div>

      </div>

      {/* Category Pills */}
      <div className="pt-3 border-t border-slate-800/60 flex flex-wrap items-center gap-2">
        <span className="text-xs font-bold text-slate-400 mr-1">Categories:</span>
        {AVAILABLE_CATEGORIES.map(cat => {
          const isSelected = selectedCategories.includes(cat);
          return (
            <button
              key={cat}
              onClick={() => toggleCategory(cat)}
              className={`px-3 py-1 rounded-full text-xs font-semibold transition-all flex items-center gap-1.5 ${
                isSelected
                  ? "bg-indigo-600 text-white shadow-md shadow-indigo-500/20"
                  : "bg-slate-950 text-slate-400 hover:text-slate-200 border border-slate-800 hover:border-slate-700"
              }`}
            >
              {isSelected && <Check className="w-3 h-3" />}
              <span className="capitalize">{cat}</span>
            </button>
          );
        })}
      </div>

      {/* Saved Presets Toolbar */}
      <div className="pt-3 border-t border-slate-800/60 flex flex-wrap items-center justify-between gap-3 text-xs">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="font-bold text-slate-400 flex items-center gap-1">
            <Bookmark className="w-3.5 h-3.5 text-indigo-400" /> Presets:
          </span>

          {savedPresets.length === 0 ? (
            <span className="text-slate-600 italic">No saved filter presets yet</span>
          ) : (
            savedPresets.map((preset: FilterPreset) => (
              <div
                key={preset.id}
                className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-slate-800/80 border border-slate-700/60 text-slate-200 hover:bg-slate-700/80 transition-all"
              >
                <button onClick={() => loadPreset(preset)} className="font-semibold hover:underline">
                  {preset.name}
                </button>
                <button
                  onClick={() => deletePreset(preset.id)}
                  className="text-slate-400 hover:text-rose-400 ml-1 p-0.5"
                  title="Delete Preset"
                >
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
            ))
          )}
        </div>

        {/* Save Current Preset Button */}
        <div>
          {!showSavePreset ? (
            <button
              onClick={() => setShowSavePreset(true)}
              className="px-3 py-1 rounded-lg bg-indigo-950/80 border border-indigo-700/60 text-indigo-300 hover:bg-indigo-900 font-semibold flex items-center gap-1.5 transition-all"
            >
              <Plus className="w-3.5 h-3.5" /> Save Filter Preset
            </button>
          ) : (
            <form onSubmit={handleSavePreset} className="flex items-center gap-1.5">
              <input
                type="text"
                placeholder="Preset Name..."
                value={presetNameInput}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPresetNameInput(e.target.value)}
                className="bg-slate-950 border border-slate-700 text-slate-100 text-xs px-2.5 py-1 rounded-lg focus:outline-none focus:border-indigo-500"
                autoFocus
              />
              <button
                type="submit"
                className="px-2.5 py-1 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-semibold transition-all"
              >
                Save
              </button>
              <button
                type="button"
                onClick={() => setShowSavePreset(false)}
                className="px-2.5 py-1 bg-slate-800 text-slate-400 hover:text-slate-200 rounded-lg text-xs"
              >
                Cancel
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
