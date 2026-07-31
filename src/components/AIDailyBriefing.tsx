"use client";

import React from "react";
import { Sparkles, Clock, ArrowRight, Zap, Bookmark } from "lucide-react";

interface AIDailyBriefingProps {
  isDarkMode?: boolean;
}

export const AIDailyBriefing: React.FC<AIDailyBriefingProps> = ({ isDarkMode = true }) => {
  const highlights = [
    { text: "GPT-6 architectural speculation and benchmark leaks dominate global social media discussions.", topic: "Models" },
    { text: "Claude Code terminal agent adoption accelerates across senior engineering teams (+142% WoW).", topic: "Agents" },
    { text: "Model Context Protocol (MCP) officially endorsed by major developer tools ecosystem.", topic: "Protocols" },
    { text: "Robotics and embodied AI venture funding saw a 3x uptick following breakthrough vision-action models.", topic: "Robotics" },
  ];

  return (
    <div
      className={`relative overflow-hidden rounded-[2rem] border p-6 sm:p-8 shadow-2xl transition-all ${
        isDarkMode
          ? "border-emerald-500/20 bg-gradient-to-b from-[#0B1924] to-[#050A0F] text-white"
          : "border-emerald-500/30 bg-gradient-to-b from-white to-slate-50 text-slate-900 shadow-xl"
      }`}
    >
      {/* Background ambient glow */}
      <div className="absolute top-0 right-0 -mt-10 -mr-10 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl pointer-events-none" />

      {/* Top Meta Bar */}
      <div className="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-white/10 mb-6">
        <div className="flex items-center gap-2.5">
          <div className="flex items-center justify-center w-7 h-7 rounded-xl bg-emerald-500/20 border border-emerald-500/30 text-emerald-400">
            <Sparkles className="w-4 h-4" />
          </div>
          <div>
            <h3 className="text-xs font-black uppercase tracking-[0.25em] text-emerald-400">AI Daily Briefing</h3>
            <p className={`text-[11px] font-medium ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
              Curated by Evolipia NLP Intelligence Engine
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3 text-xs font-mono">
          <span className={`flex items-center gap-1 px-2.5 py-1 rounded-full border ${
            isDarkMode ? "bg-slate-900/80 border-white/10 text-slate-300" : "bg-slate-100 border-slate-200 text-slate-700"
          }`}>
            <Clock className="w-3 h-3 text-emerald-400" /> 2 min read
          </span>
          <span className="text-emerald-400 font-semibold">• Generated 4 minutes ago</span>
        </div>
      </div>

      {/* Main Content Title */}
      <div className="space-y-3 mb-6">
        <h2 className={`text-2xl sm:text-3xl font-black tracking-tight ${isDarkMode ? "text-white" : "text-slate-900"}`}>
          Today&apos;s AI Landscape
        </h2>
        <p className={`text-sm leading-relaxed max-w-3xl ${isDarkMode ? "text-slate-300" : "text-slate-600"}`}>
          Key takeaways distilled from thousands of incoming papers, social threads, developer commit logs, and industry breaking announcements over the past 24 hours.
        </p>
      </div>

      {/* Executive Highlights Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mb-6">
        {highlights.map((item, idx) => (
          <div
            key={idx}
            className={`p-4 rounded-2xl border transition-all flex items-start gap-3 ${
              isDarkMode
                ? "bg-slate-900/40 border-white/10 hover:border-emerald-500/40 hover:bg-emerald-500/5"
                : "bg-white border-slate-200 hover:border-emerald-500/40 hover:bg-emerald-50"
            }`}
          >
            <span className="flex shrink-0 items-center justify-center w-6 h-6 rounded-full bg-emerald-500/20 text-emerald-400 text-xs font-black font-mono">
              {idx + 1}
            </span>
            <div className="space-y-1">
              <span className="inline-block px-2 py-0.5 text-[9px] font-black uppercase tracking-wider font-mono rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                {item.topic}
              </span>
              <p className={`text-xs sm:text-sm font-medium leading-snug ${isDarkMode ? "text-slate-200" : "text-slate-800"}`}>
                {item.text}
              </p>
            </div>
          </div>
        ))}
      </div>

      {/* Footer Action */}
      <div className="flex items-center justify-between pt-4 border-t border-white/10">
        <div className="flex items-center gap-2">
          <Zap className="w-4 h-4 text-emerald-400 animate-pulse" />
          <span className={`text-xs font-semibold ${isDarkMode ? "text-slate-400" : "text-slate-600"}`}>
            98.4% Confidence Rating across 42 sources
          </span>
        </div>
      </div>
    </div>
  );
};
