"use client";

import React, { useState, useEffect } from "react";
import { Search, Sparkles, X, FileText, ArrowUpRight } from "lucide-react";
import { NewsItem } from "@/src/api/client";

interface CommandPaletteProps {
  isOpen: boolean;
  onClose: () => void;
  items: NewsItem[];
  onSelectTopic: (keyword: string) => void;
  isDarkMode?: boolean;
}

const QUICK_ENTITIES = [
  { name: "GPT-6", category: "Models", count: "4,281 mentions" },
  { name: "Claude Code", category: "Agents", count: "2,942 mentions" },
  { name: "MCP Standard", category: "Protocols", count: "2,110 mentions" },
  { name: "Gemini CLI", category: "Tools", count: "1,487 mentions" },
  { name: "Open Source AI", category: "Ecosystem", count: "1,334 mentions" },
  { name: "Robotics & Hardware", category: "Embodied AI", count: "978 mentions" },
];

export const CommandPalette: React.FC<CommandPaletteProps> = ({
  isOpen,
  onClose,
  items,
  onSelectTopic,
  isDarkMode = true,
}) => {
  const [query, setQuery] = useState("");

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        if (isOpen) onClose();
        else setQuery("");
      }
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const filteredItems = query.trim()
    ? items.filter(
        (item) =>
          item.title.toLowerCase().includes(query.toLowerCase()) ||
          (item.summary && item.summary.toLowerCase().includes(query.toLowerCase())) ||
          (item.source_name && item.source_name.toLowerCase().includes(query.toLowerCase()))
      ).slice(0, 5)
    : items.slice(0, 4);

  const filteredEntities = query.trim()
    ? QUICK_ENTITIES.filter((e) => e.name.toLowerCase().includes(query.toLowerCase()))
    : QUICK_ENTITIES;

  return (
    <div className="fixed inset-0 z-[100] flex items-start justify-center pt-16 sm:pt-24 px-4 bg-black/60 backdrop-blur-md animate-in">
      <div
        className={`w-full max-w-2xl rounded-2xl border shadow-2xl overflow-hidden flex flex-col transition-all ${
          isDarkMode
            ? "bg-[#09131D] border-white/15 text-slate-100 shadow-emerald-950/20"
            : "bg-white border-slate-200 text-slate-900 shadow-2xl"
        }`}
      >
        {/* Search Input Bar */}
        <div className="flex items-center px-4 py-3.5 border-b border-white/10 relative">
          <Search className="w-5 h-5 text-emerald-400 shrink-0 mr-3" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Type a topic, model, company, or search intelligence signals..."
            className={`w-full bg-transparent text-sm sm:text-base outline-none font-medium placeholder:text-slate-500 ${
              isDarkMode ? "text-white" : "text-slate-900"
            }`}
            autoFocus
          />
          <button
            onClick={onClose}
            className={`p-1.5 rounded-lg border transition-all ${
              isDarkMode
                ? "bg-slate-800/80 border-slate-700 text-slate-400 hover:text-white"
                : "bg-slate-100 border-slate-200 text-slate-500 hover:text-slate-900"
            }`}
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Results Body */}
        <div className="max-h-[60vh] overflow-y-auto p-4 space-y-6 scrollbar-hide">
          {/* Quick Entities & Topics */}
          <div>
            <p className="text-[10px] font-black uppercase tracking-[0.2em] text-emerald-400 mb-2.5 flex items-center gap-1.5">
              <Sparkles className="w-3.5 h-3.5" />
              Exploding Topics & Entities
            </p>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
              {filteredEntities.map((entity) => (
                <button
                  key={entity.name}
                  onClick={() => {
                    onSelectTopic(entity.name);
                    onClose();
                  }}
                  className={`p-2.5 rounded-xl border text-left transition-all flex items-center justify-between group cursor-pointer ${
                    isDarkMode
                      ? "bg-slate-900/60 border-white/10 hover:border-emerald-500/50 hover:bg-emerald-500/10"
                      : "bg-slate-50 border-slate-200 hover:border-emerald-500/50 hover:bg-emerald-50"
                  }`}
                >
                  <div className="min-w-0 pr-2">
                    <p className={`text-xs font-bold truncate ${isDarkMode ? "text-slate-200" : "text-slate-800"}`}>
                      {entity.name}
                    </p>
                    <span className="text-[10px] text-slate-500 font-mono block mt-0.5">{entity.count}</span>
                  </div>
                  <ArrowUpRight className="w-3.5 h-3.5 text-emerald-500 opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
                </button>
              ))}
            </div>
          </div>

          {/* Signals Matching Search */}
          <div>
            <p className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mb-2.5 flex items-center gap-1.5">
              <FileText className="w-3.5 h-3.5 text-slate-400" />
              Recent Intelligence Signals ({filteredItems.length})
            </p>
            <div className="space-y-2">
              {filteredItems.map((item) => (
                <a
                  key={item.id}
                  href={item.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={onClose}
                  className={`p-3 rounded-xl border block transition-all hover:translate-x-1 ${
                    isDarkMode
                      ? "bg-slate-900/40 border-white/10 hover:border-emerald-500/40 hover:bg-emerald-500/5"
                      : "bg-white border-slate-200 hover:border-emerald-500/40 hover:bg-emerald-50"
                  }`}
                >
                  <div className="flex items-center justify-between gap-2">
                    <span className="text-[10px] font-mono font-bold text-emerald-400 px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20">
                      {item.source_name || item.domain || "SIGNAL"}
                    </span>
                    {item.relevance_score ? (
                      <span className="text-[10px] font-bold text-emerald-400 font-mono">
                        {item.relevance_score}% Impact
                      </span>
                    ) : null}
                  </div>
                  <h4 className={`text-xs sm:text-sm font-bold mt-1.5 line-clamp-1 ${isDarkMode ? "text-slate-100" : "text-slate-900"}`}>
                    {item.title}
                  </h4>
                </a>
              ))}
            </div>
          </div>
        </div>

        {/* Footer shortcuts */}
        <div className={`px-4 py-2.5 border-t text-[11px] font-mono flex items-center justify-between ${
          isDarkMode ? "border-white/10 bg-slate-950/80 text-slate-400" : "border-slate-200 bg-slate-50 text-slate-500"
        }`}>
          <div className="flex items-center gap-3">
            <span><kbd className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">ESC</kbd> Close</span>
            <span><kbd className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">⌘K</kbd> Toggle</span>
          </div>
          <span className="text-emerald-400 font-semibold">Evolipia Intelligence Search</span>
        </div>
      </div>
    </div>
  );
};
