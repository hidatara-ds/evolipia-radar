"use client";

import React from "react";
import { Layers3, TrendingUp, ArrowUpRight } from "lucide-react";

interface ClusterItem {
  ticker: string;
  keyword: string;
  frequency: string;
  growth: string;
  sentiment: "Bullish" | "Neutral" | "High Hype";
  sources: number;
}

const CLUSTERS: ClusterItem[] = [
  { ticker: "$GPT6", keyword: "GPT-6 Models", frequency: "4.2k", growth: "+38%", sentiment: "Bullish", sources: 28 },
  { ticker: "$CLAUDE", keyword: "Claude Code", frequency: "2.9k", growth: "+91%", sentiment: "High Hype", sources: 34 },
  { ticker: "$MCP", keyword: "MCP Protocol", frequency: "2.1k", growth: "+123%", sentiment: "Bullish", sources: 22 },
  { ticker: "$REASON", keyword: "Reasoning Models", frequency: "1.8k", growth: "+64%", sentiment: "Bullish", sources: 19 },
  { ticker: "$ROBOT", keyword: "Embodied Robotics", frequency: "1.4k", growth: "+201%", sentiment: "High Hype", sources: 16 },
  { ticker: "$VECDB", keyword: "Vector Database", frequency: "1.1k", growth: "+15%", sentiment: "Neutral", sources: 14 },
  { ticker: "$AGENTS", keyword: "Agent Frameworks", frequency: "980", growth: "+45%", sentiment: "Bullish", sources: 25 },
  { ticker: "$OPENSRC", keyword: "Open Source Weights", frequency: "870", growth: "+28%", sentiment: "Bullish", sources: 18 },
];

interface TrendingClustersProps {
  onSelectCluster: (keyword: string) => void;
  isDarkMode?: boolean;
}

export const TrendingClusters: React.FC<TrendingClustersProps> = ({
  onSelectCluster,
  isDarkMode = true,
}) => {
  return (
    <div
      className={`rounded-2xl border p-4 sm:p-5 transition-all shadow-md ${
        isDarkMode ? "border-white/10 bg-[#09131D]/90 text-white" : "border-slate-200 bg-white text-slate-900"
      }`}
    >
      {/* Header */}
      <div className="flex items-center justify-between gap-3 mb-3 pb-2.5 border-b border-white/10">
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <Layers3 className="w-4 h-4" />
          </div>
          <h3 className="text-sm font-black tracking-tight uppercase font-mono flex items-center gap-2">
            NLP Keyword Clusters Watchlist
            <span className="text-[10px] text-slate-400 font-normal normal-case">
              NER &amp; Embedding Clusters
            </span>
          </h3>
        </div>

        <span className="text-[10px] font-mono font-bold px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
          8 Active Heatmaps
        </span>
      </div>

      {/* Stock Watchlist Compact Table Grid (2 rows x 4 columns) */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
        {CLUSTERS.map((cluster) => (
          <button
            key={cluster.ticker}
            onClick={() => onSelectCluster(cluster.keyword)}
            className={`p-2.5 rounded-xl border text-left transition-all relative flex flex-col justify-between cursor-pointer group ${
              isDarkMode
                ? "bg-slate-900/60 border-white/10 hover:border-emerald-500/50 hover:bg-emerald-500/10"
                : "bg-slate-50 border-slate-200 hover:border-emerald-500/50 hover:bg-emerald-50"
            }`}
          >
            <div className="flex items-center justify-between gap-1 mb-1">
              <span className="text-[10px] font-mono font-bold text-slate-400">
                {cluster.ticker}
              </span>
              <span className="text-[9px] font-mono font-bold px-1.5 py-0.2 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                {cluster.sentiment}
              </span>
            </div>

            <div className={`text-xs font-black truncate group-hover:text-emerald-400 transition-colors mb-2 ${
              isDarkMode ? "text-white" : "text-slate-900"
            }`}>
              {cluster.keyword}
            </div>

            <div className="flex items-center justify-between text-xs font-mono pt-1.5 border-t border-white/10">
              <span className={isDarkMode ? "text-slate-400" : "text-slate-500"}>
                {cluster.frequency} <span className="text-[10px] opacity-75">({cluster.sources} src)</span>
              </span>
              <span className="font-black text-emerald-400 flex items-center gap-0.5">
                <ArrowUpRight className="w-3 h-3" />
                {cluster.growth}
              </span>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};
