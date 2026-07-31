"use client";

import React from "react";
import { TrendingUp, Flame, Tag, ArrowUpRight } from "lucide-react";

interface AIPulseKeywordsProps {
  activeKeyword: string | null;
  onSelectKeyword: (keyword: string) => void;
  isDarkMode?: boolean;
}

export interface PulseTopic {
  name: string;
  mentions: string;
  growth: string;
  isExploding?: boolean;
  category: string;
}

export const PULSE_TOPICS: PulseTopic[] = [
  { name: "GPT-6", mentions: "4,281", growth: "+38%", isExploding: false, category: "Model" },
  { name: "Claude Code", mentions: "2,942", growth: "+91%", isExploding: true, category: "Agent" },
  { name: "MCP", mentions: "2,110", growth: "+123%", isExploding: true, category: "Protocol" },
  { name: "Gemini CLI", mentions: "1,487", growth: "+51%", isExploding: false, category: "CLI" },
  { name: "Open Source", mentions: "1,334", growth: "+17%", isExploding: false, category: "Ecosystem" },
  { name: "Robotics", mentions: "978", growth: "+201%", isExploding: true, category: "Hardware" },
];

export const AIPulseKeywords: React.FC<AIPulseKeywordsProps> = ({
  activeKeyword,
  onSelectKeyword,
  isDarkMode = true,
}) => {
  return (
    <div
      className={`rounded-[2rem] border p-6 shadow-xl transition-all ${
        isDarkMode
          ? "border-white/10 bg-[#09131D]/90 text-white"
          : "border-slate-200 bg-white text-slate-900"
      }`}
    >
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-5">
        <div className="flex items-center gap-2.5">
          <div className="p-2 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
            <Flame className="w-5 h-5" />
          </div>
          <div>
            <h3 className="text-lg font-black tracking-tight">AI Pulse</h3>
            <p className={`text-xs ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
              Real-time keyword mention volume &amp; 24h momentum velocity
            </p>
          </div>
        </div>

        {activeKeyword && (
          <button
            onClick={() => onSelectKeyword("")}
            className="text-xs font-bold font-mono text-emerald-400 hover:underline self-start sm:self-auto cursor-pointer"
          >
            Clear Filter ({activeKeyword})
          </button>
        )}
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
        {PULSE_TOPICS.map((topic) => {
          const isActive = activeKeyword?.toLowerCase() === topic.name.toLowerCase();
          return (
            <button
              key={topic.name}
              onClick={() => onSelectKeyword(isActive ? "" : topic.name)}
              className={`group p-4 rounded-2xl border text-left transition-all relative overflow-hidden cursor-pointer ${
                isActive
                  ? "bg-emerald-500 border-emerald-400 text-slate-950 shadow-lg shadow-emerald-500/30 scale-[1.02]"
                  : isDarkMode
                  ? "bg-slate-900/50 border-white/10 hover:border-emerald-500/40 hover:bg-emerald-500/5"
                  : "bg-slate-50 border-slate-200 hover:border-emerald-500/40 hover:bg-emerald-50"
              }`}
            >
              {topic.isExploding && !isActive && (
                <span className="absolute top-2 right-2 flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                </span>
              )}

              <div className="flex items-center justify-between gap-1 mb-2">
                <span
                  className={`text-[9px] font-mono font-bold uppercase tracking-wider px-1.5 py-0.5 rounded ${
                    isActive
                      ? "bg-slate-950/20 text-slate-950"
                      : "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                  }`}
                >
                  {topic.category}
                </span>
                <ArrowUpRight
                  className={`w-3.5 h-3.5 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5 ${
                    isActive ? "text-slate-950" : "text-slate-400 group-hover:text-emerald-400"
                  }`}
                />
              </div>

              <h4 className={`text-base font-black tracking-tight ${isActive ? "text-slate-950" : isDarkMode ? "text-white" : "text-slate-900"}`}>
                {topic.name}
              </h4>

              <div className="mt-3 flex items-baseline justify-between gap-1 font-mono">
                <span className={`text-xs font-bold ${isActive ? "text-slate-900" : isDarkMode ? "text-slate-300" : "text-slate-700"}`}>
                  {topic.mentions}
                </span>
                <span className={`text-xs font-black ${isActive ? "text-slate-950" : "text-emerald-400"}`}>
                  {topic.growth}
                </span>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
};
