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

const SORT_OPTIONS = [
  { id: "date_desc", label: "Date (Newest First)", sortBy: "date", sortOrder: "desc" },
  { id: "date_asc", label: "Date (Oldest First)", sortBy: "date", sortOrder: "asc" },
  { id: "relevance_desc", label: "Relevance Score (High to Low)", sortBy: "relevance", sortOrder: "desc" },
  { id: "credibility_desc", label: "Source Credibility", sortBy: "credibility", sortOrder: "desc" },
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

  return (
    <div className="w-full bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-xl mb-6 text-slate-100 space-y-4">
      {/* Top Bar: Search & Sort */}
      <div className="flex flex-col md:flex-row items-center justify-between gap-4">
        {/* Search Input */}
        <div className="relative w-full md:w-96">
          <Search className="w-4 h-4 absolute left-3 top-3 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearch(e.target.value)}
            placeholder="Search title, content, or domain..."
            className="w-full pl-9 pr-4 py-2 bg-slate-950 border border-slate-800 rounded-lg text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-indigo-500 transition-colors"
          />
        </div>

        {/* Sort Selector & Preset Controls */}
        <div className="flex items-center gap-3 w-full md:w-auto justify-between md:justify-end">
          <div className="flex items-center gap-2">
            <ArrowUpDown className="w-4 h-4 text-slate-400" />
            <select
              value={`${sortBy}_${sortOrder}`}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                const opt = SORT_OPTIONS.find(o => o.id === e.target.value);
                if (opt) {
                  setSortBy(opt.sortBy);
                  setSortOrder(opt.sortOrder);
                }
              }}
              className="bg-slate-950 border border-slate-800 text-xs font-semibold text-slate-200 px-3 py-2 rounded-lg focus:outline-none focus:border-indigo-500 cursor-pointer"
            >
              {SORT_OPTIONS.map(opt => (
                <option key={opt.id} value={opt.id}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          <button
            onClick={resetFilters}
            className="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold flex items-center gap-1.5 transition-all"
            title="Reset Filters"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            <span>Reset</span>
          </button>
        </div>
      </div>

      {/* Main Filter Controls Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 pt-4 border-t border-slate-800">
        {/* Date Range Selector */}
        <div>
          <label className="block text-xs font-semibold text-slate-400 mb-1.5 flex items-center gap-1.5">
            <Calendar className="w-3.5 h-3.5 text-indigo-400" /> Date Range
          </label>
          <select
            value={dateRange}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
              const val = e.target.value;
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
              }
            }}
            className="w-full bg-slate-950 border border-slate-800 text-xs text-slate-200 p-2 rounded-lg focus:outline-none focus:border-indigo-500"
          >
            <option value="today">Today</option>
            <option value="7d">Last 7 Days</option>
            <option value="30d">Last 30 Days</option>
            <option value="custom">Custom Range</option>
          </select>

          {dateRange === "custom" && (
            <div className="flex gap-2 mt-2">
              <input
                type="date"
                value={dateFrom}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDateFrom(e.target.value)}
                className="w-1/2 bg-slate-950 border border-slate-800 text-xs text-slate-200 p-1.5 rounded-md"
              />
              <input
                type="date"
                value={dateTo}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDateTo(e.target.value)}
                className="w-1/2 bg-slate-950 border border-slate-800 text-xs text-slate-200 p-1.5 rounded-md"
              />
            </div>
          )}
        </div>

        {/* Source Multi-select Dropdown */}
        <div className="relative">
          <label className="block text-xs font-semibold text-slate-400 mb-1.5 flex items-center gap-1.5">
            <Filter className="w-3.5 h-3.5 text-indigo-400" /> Sources ({selectedSources.length})
          </label>
          <button
            onClick={() => setShowSourceDropdown(!showSourceDropdown)}
            className="w-full bg-slate-950 border border-slate-800 text-xs text-slate-200 p-2 rounded-lg flex items-center justify-between focus:outline-none"
          >
            <span className="truncate">
              {selectedSources.length === 0
                ? "All Sources"
                : `${selectedSources.length} selected`}
            </span>
            <ChevronDown className="w-3.5 h-3.5 text-slate-400" />
          </button>

          {showSourceDropdown && (
            <div className="absolute left-0 right-0 top-full mt-1 bg-slate-950 border border-slate-800 rounded-lg p-2 z-20 shadow-xl max-h-48 overflow-y-auto space-y-1">
              {AVAILABLE_SOURCES.map(src => (
                <label
                  key={src}
                  className="flex items-center gap-2 text-xs text-slate-300 hover:bg-slate-900 p-1.5 rounded cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={selectedSources.includes(src)}
                    onChange={() => toggleSource(src)}
                    className="rounded border-slate-700 bg-slate-900 text-indigo-600 focus:ring-0"
                  />
                  <span>{src}</span>
                </label>
              ))}
            </div>
          )}
        </div>

        {/* Relevance Score Slider */}
        <div>
          <div className="flex justify-between items-center mb-1.5">
            <label className="text-xs font-semibold text-slate-400 flex items-center gap-1.5">
              <Sparkles className="w-3.5 h-3.5 text-amber-400" /> Min Relevance
            </label>
            <span className="text-xs font-mono text-amber-300 font-bold">{minRelevance}%</span>
          </div>
          <input
            type="range"
            min="0"
            max="100"
            step="5"
            value={minRelevance}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setMinRelevance(Number(e.target.value))}
            className="w-full accent-indigo-500 cursor-pointer"
          />
        </div>

        {/* Status Selector */}
        <div>
          <label className="block text-xs font-semibold text-slate-400 mb-1.5 flex items-center gap-1.5">
            <SlidersHorizontal className="w-3.5 h-3.5 text-indigo-400" /> Verification Status
          </label>
          <select
            value={status}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => setStatus(e.target.value)}
            className="w-full bg-slate-950 border border-slate-800 text-xs text-slate-200 p-2 rounded-lg focus:outline-none focus:border-indigo-500"
          >
            <option value="all">All Items</option>
            <option value="verified">Verified (Done)</option>
            <option value="pending">Pending</option>
            <option value="failed">Failed</option>
          </select>
        </div>
      </div>

      {/* Category Pills */}
      <div className="pt-2 flex flex-wrap items-center gap-2">
        <span className="text-xs font-semibold text-slate-400 mr-1">Categories:</span>
        {AVAILABLE_CATEGORIES.map(cat => {
          const isSelected = selectedCategories.includes(cat);
          return (
            <button
              key={cat}
              onClick={() => toggleCategory(cat)}
              className={`px-2.5 py-1 rounded-full text-xs font-semibold transition-all flex items-center gap-1 ${
                isSelected
                  ? "bg-indigo-600 text-white shadow-md shadow-indigo-500/20"
                  : "bg-slate-950 text-slate-400 hover:text-slate-200 border border-slate-800"
              }`}
            >
              {isSelected && <Check className="w-3 h-3" />}
              {cat}
            </button>
          );
        })}
      </div>

      {/* Saved Filter Presets Toolbar */}
      <div className="pt-3 border-t border-slate-800/60 flex flex-wrap items-center justify-between gap-3 text-xs">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="font-semibold text-slate-400 flex items-center gap-1">
            <Bookmark className="w-3.5 h-3.5 text-indigo-400" /> Saved Filters:
          </span>

          {savedPresets.length === 0 ? (
            <span className="text-slate-600 italic">No saved presets yet</span>
          ) : (
            savedPresets.map((preset: FilterPreset) => (
              <div
                key={preset.id}
                className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md bg-slate-800 border border-slate-700 text-slate-200 hover:bg-slate-700"
              >
                <button onClick={() => loadPreset(preset)} className="font-medium hover:underline">
                  {preset.name}
                </button>
                <button
                  onClick={() => deletePreset(preset.id)}
                  className="text-slate-500 hover:text-rose-400 ml-1"
                >
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
            ))
          )}
        </div>

        {/* Save Current Preset Button & Form */}
        <div className="relative">
          {!showSavePreset ? (
            <button
              onClick={() => setShowSavePreset(true)}
              className="px-2.5 py-1 rounded bg-indigo-950 border border-indigo-700/60 text-indigo-300 hover:bg-indigo-900 font-semibold flex items-center gap-1"
            >
              <Plus className="w-3 h-3" /> Save Preset
            </button>
          ) : (
            <form onSubmit={handleSavePreset} className="flex items-center gap-1.5">
              <input
                type="text"
                placeholder="Preset Name..."
                value={presetNameInput}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPresetNameInput(e.target.value)}
                className="bg-slate-950 border border-slate-700 text-slate-100 text-xs px-2 py-1 rounded focus:outline-none"
                autoFocus
              />
              <button
                type="submit"
                className="px-2 py-1 bg-indigo-600 text-white rounded text-xs font-semibold"
              >
                Save
              </button>
              <button
                type="button"
                onClick={() => setShowSavePreset(false)}
                className="px-2 py-1 bg-slate-800 text-slate-400 rounded text-xs"
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
