"use client";

import React from "react";
import { Sparkles, Building2, Cpu, Globe, ArrowUpRight, Flame, Layers } from "lucide-react";
import { NewsItem } from "@/src/api/client";

interface IntelligenceSidebarProps {
  topSignals: NewsItem[];
  onSelectTopic: (keyword: string) => void;
  isDarkMode?: boolean;
}

const MOST_DISCUSSED_MODELS = [
  { name: "Claude 3.7 Sonnet", org: "Anthropic", growth: "+114%" },
  { name: "GPT-6 Experimental", org: "OpenAI", growth: "+92%" },
  { name: "Gemini 2.5 Pro", org: "Google", growth: "+48%" },
  { name: "DeepSeek V3", org: "DeepSeek", growth: "+37%" },
];

const MOST_DISCUSSED_COMPANIES = [
  { name: "Anthropic", mentions: "1,842", sentiment: "Bullish" },
  { name: "OpenAI", mentions: "2,410", sentiment: "Bullish" },
  { name: "Google DeepMind", mentions: "1,290", sentiment: "Bullish" },
  { name: "Meta AI", mentions: "950", sentiment: "Neutral" },
];

const RESEARCH_OPPORTUNITIES = [
  "Long-context retrieval & reasoning benchmark trade-offs",
  "Agent tool-use latency optimization via MCP",
  "Small model fine-tuning with synthetic reasoning traces",
];

export const IntelligenceSidebar: React.FC<IntelligenceSidebarProps> = ({
  topSignals,
  onSelectTopic,
  isDarkMode = true,
}) => {
  return (
    <div className="space-y-5">
      {/* Today's Brief & Priority Queue */}
      <div
        className={`rounded-[1.75rem] border p-5 transition-all ${
          isDarkMode ? "border-white/10 bg-[#09131D] text-white" : "border-slate-200 bg-white text-slate-900 shadow-sm"
        }`}
      >
        <div className="flex items-center justify-between mb-4">
          <div>
            <p className="text-[10px] font-black uppercase tracking-[0.2em] text-emerald-400 font-mono">
              Priority Queue
            </p>
            <h4 className="text-base font-black tracking-tight mt-0.5">Top Signals To Read</h4>
          </div>
          <Sparkles className="w-4 h-4 text-emerald-400" />
        </div>

        <div className="space-y-2.5">
          {topSignals.slice(0, 3).map((item, index) => (
            <a
              key={item.id}
              href={item.url}
              target="_blank"
              rel="noopener noreferrer"
              className={`block p-3 rounded-xl border transition-all hover:border-emerald-500/40 hover:bg-emerald-500/5 ${
                isDarkMode ? "border-white/10 bg-slate-900/40" : "border-slate-200 bg-slate-50"
              }`}
            >
              <div className="flex items-center justify-between text-[10px] font-mono font-bold text-emerald-400 mb-1">
                <span>0{index + 1} • {item.source_name || item.domain || "SOURCE"}</span>
                {item.relevance_score ? <span>{item.relevance_score}% Relevance</span> : null}
              </div>
              <p className={`line-clamp-2 text-xs font-bold leading-snug ${isDarkMode ? "text-slate-200" : "text-slate-800"}`}>
                {item.title}
              </p>
            </a>
          ))}
        </div>
      </div>

      {/* Most Discussed Models */}
      <div
        className={`rounded-[1.75rem] border p-5 transition-all ${
          isDarkMode ? "border-white/10 bg-[#09131D] text-white" : "border-slate-200 bg-white text-slate-900 shadow-sm"
        }`}
      >
        <div className="flex items-center gap-2 mb-3.5">
          <Cpu className="w-4 h-4 text-emerald-400" />
          <h4 className="text-sm font-black tracking-tight">Most Discussed Models</h4>
        </div>

        <div className="space-y-2">
          {MOST_DISCUSSED_MODELS.map((model) => (
            <button
              key={model.name}
              onClick={() => onSelectTopic(model.name)}
              className={`w-full p-2.5 rounded-xl border text-left flex items-center justify-between transition-all cursor-pointer ${
                isDarkMode
                  ? "border-white/10 bg-slate-900/30 hover:border-emerald-500/40 hover:bg-emerald-500/5"
                  : "border-slate-200 bg-slate-50 hover:border-emerald-500/40 hover:bg-emerald-50"
              }`}
            >
              <div>
                <p className={`text-xs font-bold ${isDarkMode ? "text-slate-200" : "text-slate-800"}`}>
                  {model.name}
                </p>
                <span className="text-[10px] text-slate-500 font-mono">{model.org}</span>
              </div>
              <span className="text-xs font-mono font-bold text-emerald-400">{model.growth}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Research Opportunities */}
      <div
        className={`rounded-[1.75rem] border p-5 transition-all ${
          isDarkMode ? "border-white/10 bg-[#09131D] text-white" : "border-slate-200 bg-white text-slate-900 shadow-sm"
        }`}
      >
        <div className="flex items-center gap-2 mb-3">
          <Layers className="w-4 h-4 text-emerald-400" />
          <h4 className="text-sm font-black tracking-tight">Research Opportunities</h4>
        </div>

        <ul className="space-y-2 text-xs font-medium">
          {RESEARCH_OPPORTUNITIES.map((opp, i) => (
            <li
              key={i}
              className={`p-2.5 rounded-xl border flex items-start gap-2 ${
                isDarkMode ? "border-white/10 bg-slate-900/30 text-slate-300" : "border-slate-200 bg-slate-50 text-slate-700"
              }`}
            >
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 shrink-0 mt-1.5" />
              <span>{opp}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};
