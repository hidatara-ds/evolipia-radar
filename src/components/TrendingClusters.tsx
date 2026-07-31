"use client";

import React from "react";
import { Layers3, Sparkles, TrendingUp, BarChart2, Hash, ArrowUpRight } from "lucide-react";

interface ClusterItem {
  keyword: string;
  frequency: string;
  growth: string;
  momentum: string;
  sentiment: "Bullish" | "Neutral" | "High Hype";
  sources: number;
  score: number;
}

const CLUSTERS: ClusterItem[] = [
  { keyword: "GPT-6", frequency: "4.2k", growth: "+38%", momentum: "Exploding", sentiment: "Bullish", sources: 28, score: 98 },
  { keyword: "Claude Code", frequency: "2.9k", growth: "+91%", momentum: "Surging", sentiment: "High Hype", sources: 34, score: 96 },
  { keyword: "MCP Protocol", frequency: "2.1k", growth: "+123%", momentum: "Breakthrough", sentiment: "Bullish", sources: 22, score: 94 },
  { keyword: "Reasoning Models", frequency: "1.8k", growth: "+64%", momentum: "Rising", sentiment: "Bullish", sources: 19, score: 91 },
  { keyword: "Embodied Robotics", frequency: "1.4k", growth: "+201%", momentum: "Exploding", sentiment: "High Hype", sources: 16, score: 89 },
  { keyword: "Vector Database", frequency: "1.1k", growth: "+15%", momentum: "Steady", sentiment: "Neutral", sources: 14, score: 82 },
  { keyword: "Agent Frameworks", frequency: "980", growth: "+45%", momentum: "Surging", sentiment: "Bullish", sources: 25, score: 88 },
  { keyword: "Open Source Weights", frequency: "870", growth: "+28%", momentum: "Steady", sentiment: "Bullish", sources: 18, score: 85 },
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
      className={`rounded-[2rem] border p-6 shadow-xl transition-all ${
        isDarkMode ? "border-white/10 bg-[#09131D]/80 text-white" : "border-slate-200 bg-white text-slate-900"
      }`}
    >
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-5">
        <div className="flex items-center gap-2.5">
          <div className="p-2 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
            <Layers3 className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-lg font-black tracking-tight">NLP Keyword Clusters &amp; Topic Heatmap</h3>
            <p className={`text-xs ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
              Extracted via Named Entity Recognition (NER), TF-IDF, &amp; Semantic Embedding Similarity
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 text-xs font-mono">
          <span className="px-2.5 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 font-bold">
            8 Clusters Active
          </span>
        </div>
      </div>

      {/* Cluster Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5">
        {CLUSTERS.map((cluster) => (
          <button
            key={cluster.keyword}
            onClick={() => onSelectCluster(cluster.keyword)}
            className={`group p-4 rounded-2xl border text-left transition-all relative flex flex-col justify-between cursor-pointer ${
              isDarkMode
                ? "bg-slate-900/40 border-white/10 hover:border-emerald-500/50 hover:bg-emerald-500/5"
                : "bg-slate-50 border-slate-200 hover:border-emerald-500/50 hover:bg-emerald-50"
            }`}
          >
            <div>
              <div className="flex items-center justify-between gap-2 mb-2">
                <span className="text-[10px] font-mono font-bold px-2 py-0.5 rounded bg-emerald-500/15 border border-emerald-500/30 text-emerald-400">
                  {cluster.sentiment}
                </span>
                <span className="text-xs font-mono font-bold text-emerald-400 flex items-center gap-0.5">
                  <TrendingUp className="w-3 h-3" /> {cluster.growth}
                </span>
              </div>

              <h4 className={`text-base font-black tracking-tight group-hover:text-emerald-400 transition-colors ${
                isDarkMode ? "text-white" : "text-slate-900"
              }`}>
                {cluster.keyword}
              </h4>
            </div>

            <div className="mt-4 pt-3 border-t border-white/10 space-y-2">
              <div className="flex items-center justify-between text-xs font-mono">
                <span className={isDarkMode ? "text-slate-400" : "text-slate-500"}>Mentions</span>
                <span className="font-bold">{cluster.frequency}</span>
              </div>
              <div className="flex items-center justify-between text-xs font-mono">
                <span className={isDarkMode ? "text-slate-400" : "text-slate-500"}>Sources</span>
                <span className="font-bold">{cluster.sources} active</span>
              </div>

              {/* Progress bar visual for cluster momentum */}
              <div className="w-full h-1.5 rounded-full bg-slate-800 overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-emerald-500 to-cyan-400 rounded-full"
                  style={{ width: `${cluster.score}%` }}
                />
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};
