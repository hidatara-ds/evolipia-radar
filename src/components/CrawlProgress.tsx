"use client";

import React, { useState } from "react";
import { CrawlProgressEvent } from "../hooks/useCrawlProgress";
import { CheckCircle2, AlertCircle, Loader2, Play, ChevronDown, ChevronUp, Terminal, Shield } from "lucide-react";

interface CrawlProgressProps {
  progressState: CrawlProgressEvent | null;
  isCrawling: boolean;
  onStartManualCrawl: () => void;
  toastMessage: { type: "success" | "error"; text: string } | null;
  onClearToast: () => void;
  isDarkMode?: boolean;
}

const STEPS = [
  "Initializing crawler...",
  "Scanning sources...",
  "Parsing content...",
  "Validating data...",
  "Saving to database...",
  "Done!",
];

export const CrawlProgress: React.FC<CrawlProgressProps> = ({
  progressState,
  isCrawling,
  onStartManualCrawl,
  toastMessage,
  onClearToast,
  isDarkMode = true,
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const currentStep = progressState?.step || (isCrawling ? 1 : 0);
  const progressPct = progressState?.progress ?? (isCrawling ? 10 : 0);

  return (
    <div
      className={`w-full rounded-2xl transition-all border ${
        isDarkMode
          ? "bg-[#09131D]/60 border-white/10 text-slate-100"
          : "bg-slate-50 border-slate-200 text-slate-900 shadow-sm"
      }`}
    >
      {/* Compact Header Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-4">
        <div className="flex items-center gap-3">
          <div className={`p-2 rounded-xl border ${
            isDarkMode ? "bg-slate-900 border-white/10 text-emerald-400" : "bg-white border-slate-200 text-emerald-600"
          }`}>
            <Terminal className="w-4 h-4" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h4 className="text-xs font-black uppercase tracking-wider">
                Developer &amp; Crawler Status
              </h4>
              <span className={`text-[10px] font-mono px-2 py-0.5 rounded-full border ${
                isCrawling
                  ? "bg-amber-500/20 border-amber-500/30 text-amber-400 font-bold animate-pulse"
                  : "bg-emerald-500/10 border-emerald-500/20 text-emerald-400 font-semibold"
              }`}>
                {isCrawling ? "Active Crawl" : "Background Sync Idle"}
              </span>
            </div>
            <p className={`text-[11px] mt-0.5 ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
              {isCrawling
                ? progressState?.message || "Collecting and embedding incoming signals..."
                : "Continuous NLP extraction running on 6h schedule."}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={onStartManualCrawl}
            disabled={isCrawling}
            className={`px-3 py-1.5 rounded-xl text-xs font-bold flex items-center gap-1.5 transition-all cursor-pointer ${
              isCrawling
                ? "bg-slate-800 text-slate-500 cursor-not-allowed border border-slate-700"
                : "bg-emerald-500 hover:bg-emerald-400 text-slate-950 shadow-md shadow-emerald-500/20 active:scale-95"
            }`}
          >
            {isCrawling ? (
              <>
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
                Crawling...
              </>
            ) : (
              <>
                <Play className="w-3.5 h-3.5 fill-current" />
                Trigger Crawl
              </>
            )}
          </button>

          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className={`p-1.5 rounded-xl border transition-all cursor-pointer ${
              isDarkMode ? "bg-slate-900 border-white/10 text-slate-400 hover:text-white" : "bg-white border-slate-200 text-slate-600 hover:text-slate-900"
            }`}
            title={isExpanded ? "Collapse Details" : "Expand Technical Monitor"}
          >
            {isExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
          </button>
        </div>
      </div>

      {/* Expanded Technical Inspector */}
      {(isExpanded || isCrawling) && (
        <div className={`p-4 border-t space-y-4 ${isDarkMode ? "border-white/10 bg-slate-950/40" : "border-slate-200 bg-white"}`}>
          {isCrawling && (
            <div className="space-y-2">
              <div className="flex justify-between text-xs font-mono font-semibold">
                <span className={isDarkMode ? "text-slate-300" : "text-slate-700"}>
                  {progressState?.current_source
                    ? `Ingesting source: ${progressState.current_source}`
                    : progressState?.message || "Parsing..."}
                </span>
                <span className="text-emerald-400 font-bold">{progressPct}%</span>
              </div>
              <div className="w-full h-2 rounded-full overflow-hidden bg-slate-800">
                <div
                  className="h-full bg-gradient-to-r from-emerald-500 to-cyan-400 transition-all duration-300"
                  style={{ width: `${Math.max(5, Math.min(100, progressPct))}%` }}
                />
              </div>
            </div>
          )}

          {/* Stepper details */}
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-2">
            {STEPS.map((stepLabel, idx) => {
              const stepNum = idx + 1;
              const isActive = currentStep === stepNum;
              const isDone = currentStep > stepNum || (!isCrawling && progressState?.is_complete);

              return (
                <div
                  key={idx}
                  className={`p-2 rounded-xl text-[11px] font-medium border ${
                    isActive
                      ? "bg-emerald-950/60 border-emerald-500 text-emerald-200"
                      : isDone
                      ? "bg-slate-900/60 border-emerald-500/40 text-emerald-400"
                      : "bg-slate-900/20 border-slate-800 text-slate-500"
                  }`}
                >
                  <div className="flex items-center gap-1 font-mono text-[10px] font-bold mb-0.5">
                    {isDone ? (
                      <CheckCircle2 className="w-3 h-3 text-emerald-400" />
                    ) : isActive ? (
                      <Loader2 className="w-3 h-3 text-emerald-400 animate-spin" />
                    ) : (
                      <span>Step {stepNum}</span>
                    )}
                  </div>
                  <p className="truncate font-sans text-xs">{stepLabel}</p>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Toast Notification */}
      {toastMessage && (
        <div className={`p-3 border-t text-xs font-medium flex items-center justify-between ${
          toastMessage.type === "success" ? "bg-emerald-950/80 text-emerald-200" : "bg-rose-950/80 text-rose-200"
        }`}>
          <div className="flex items-center gap-2">
            {toastMessage.type === "success" ? <CheckCircle2 className="w-4 h-4 text-emerald-400" /> : <AlertCircle className="w-4 h-4 text-rose-400" />}
            <span>{toastMessage.text}</span>
          </div>
          <button onClick={onClearToast} className="text-xs font-bold text-slate-400 hover:text-white">
            Dismiss
          </button>
        </div>
      )}
    </div>
  );
};
