"use client";

import React from "react";
import { TrendingUp, Flame, ArrowUpRight, TrendingDown } from "lucide-react";

interface AIPulseKeywordsProps {
  activeKeyword: string | null;
  onSelectKeyword: (keyword: string) => void;
  isDarkMode?: boolean;
}

export interface PulseTopic {
  ticker: string;
  name: string;
  mentions: string;
  growth: string;
  isPositive: boolean;
  category: string;
  volume: string;
}

export const STOCK_TICKERS: PulseTopic[] = [
  { ticker: "$GPT6", name: "GPT-6 Models", mentions: "4,281", growth: "+38%", isPositive: true, category: "MODEL", volume: "High" },
  { ticker: "$CLAUDE", name: "Claude Code", mentions: "2,942", growth: "+91%", isPositive: true, category: "AGENT", volume: "Surging" },
  { ticker: "$MCP", name: "MCP Standard", mentions: "2,110", growth: "+123%", isPositive: true, category: "PROTOCOL", volume: "Exploding" },
  { ticker: "$GEMINI", name: "Gemini CLI", mentions: "1,487", growth: "+51%", isPositive: true, category: "CLI", volume: "Active" },
  { ticker: "$OPENSRC", name: "Open Source AI", mentions: "1,334", growth: "+17%", isPositive: true, category: "ECOSYSTEM", volume: "Steady" },
  { ticker: "$ROBOT", name: "Robotics & HW", mentions: "978", growth: "+201%", isPositive: true, category: "EMBODIED", volume: "Exploding" },
];

export const AIPulseKeywords: React.FC<AIPulseKeywordsProps> = ({
  activeKeyword,
  onSelectKeyword,
  isDarkMode = true,
}) => {
  return (
    <div
      className={`rounded-2xl border p-4 sm:p-5 transition-all shadow-md ${
        isDarkMode
          ? "border-white/10 bg-[#09131D]/90 text-white"
          : "border-slate-200 bg-white text-slate-900"
      }`}
    >
      {/* Header Bar */}
      <div className="flex items-center justify-between gap-3 mb-3 pb-2.5 border-b border-white/10">
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <Flame className="w-4 h-4" />
          </div>
          <div>
            <h3 className="text-sm font-black tracking-tight uppercase font-mono flex items-center gap-2">
              AI Stock Pulse Watchlist
              <span className="text-[10px] font-normal text-emerald-400 normal-case px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20">
                24h Market Velocity
              </span>
            </h3>
          </div>
        </div>

        {activeKeyword && (
          <button
            onClick={() => onSelectKeyword("")}
            className="text-xs font-bold font-mono text-emerald-400 hover:underline cursor-pointer"
          >
            Reset Filter ({activeKeyword})
          </button>
        )}
      </div>

      {/* Stock Ticker Grid (Compact 3-column / 6-column watchlist) */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
        {STOCK_TICKERS.map((stock) => {
          const isActive = activeKeyword?.toLowerCase() === stock.name.toLowerCase() || activeKeyword?.toLowerCase() === stock.ticker.toLowerCase();
          return (
            <button
              key={stock.ticker}
              onClick={() => onSelectKeyword(isActive ? "" : stock.name)}
              className={`p-2.5 rounded-xl border text-left transition-all relative cursor-pointer flex flex-col justify-between ${
                isActive
                  ? "bg-emerald-500 border-emerald-400 text-slate-950 shadow-md scale-[1.02]"
                  : isDarkMode
                  ? "bg-slate-900/60 border-white/10 hover:border-emerald-500/50 hover:bg-emerald-500/10"
                  : "bg-slate-50 border-slate-200 hover:border-emerald-500/50 hover:bg-emerald-50"
              }`}
            >
              <div className="flex items-center justify-between gap-1 mb-1">
                <span className={`text-[10px] font-mono font-black ${isActive ? "text-slate-950" : "text-slate-400"}`}>
                  {stock.ticker}
                </span>
                <span
                  className={`text-[9px] font-mono font-bold px-1.5 py-0.2 rounded ${
                    isActive
                      ? "bg-slate-950/20 text-slate-950"
                      : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                  }`}
                >
                  {stock.category}
                </span>
              </div>

              <div className={`text-xs font-black truncate mb-1.5 ${isActive ? "text-slate-950" : isDarkMode ? "text-white" : "text-slate-900"}`}>
                {stock.name}
              </div>

              <div className="flex items-baseline justify-between font-mono pt-1.5 border-t border-white/10">
                <span className={`text-xs font-bold ${isActive ? "text-slate-900" : isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
                  {stock.mentions}
                </span>
                <span className={`text-xs font-black flex items-center gap-0.5 ${isActive ? "text-slate-950" : "text-emerald-400"}`}>
                  <ArrowUpRight className="w-3 h-3" />
                  {stock.growth}
                </span>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
};
