"use client";

import React from "react";
import { CrawlProgressEvent } from "../hooks/useCrawlProgress";
import { CheckCircle2, AlertCircle, Loader2, Play } from "lucide-react";

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
  const currentStep = progressState?.step || (isCrawling ? 1 : 0);
  const progressPct = progressState?.progress ?? (isCrawling ? 10 : 0);

  return (
    <div className={`w-full rounded-2xl p-5 sm:p-6 mb-6 transition-all border ${
      isDarkMode
        ? "bg-slate-950/70 border-white/10 text-slate-100 shadow-2xl"
        : "bg-white border-slate-200 text-slate-900 shadow-lg"
    }`}>
      {/* Header & Controls */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
        <div>
          <h3 className={`text-lg font-black flex items-center gap-2 ${
            isDarkMode ? "text-white" : "text-slate-900"
          }`}>
            {isCrawling ? (
              <Loader2 className="w-5 h-5 text-emerald-400 animate-spin" />
            ) : (
              <CheckCircle2 className="w-5 h-5 text-emerald-500" />
            )}
            Auto-Crawl Status & Real-time Progress
          </h3>
          <p className={`text-xs mt-0.5 ${isDarkMode ? "text-slate-400" : "text-slate-500"}`}>
            {isCrawling
              ? progressState?.message || "Crawler is actively running..."
              : "Auto-scheduler runs every 6 hours automatically."}
          </p>
        </div>

        <button
          onClick={onStartManualCrawl}
          disabled={isCrawling}
          className={`px-4 py-2.5 rounded-xl text-xs font-black flex items-center justify-center gap-2 transition-all ${
            isCrawling
              ? isDarkMode
                ? "bg-slate-800 text-slate-500 cursor-not-allowed border border-slate-700"
                : "bg-slate-200 text-slate-400 cursor-not-allowed border border-slate-300"
              : "bg-emerald-500 hover:bg-emerald-400 text-slate-950 shadow-md shadow-emerald-500/20 active:scale-95 cursor-pointer"
          }`}
        >
          {isCrawling ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              Crawling In Progress...
            </>
          ) : (
            <>
              <Play className="w-4 h-4 fill-current" />
              Trigger Manual Crawl
            </>
          )}
        </button>
      </div>

      {/* Progress Bar & Details */}
      {isCrawling && (
        <div className={`space-y-3 mt-4 pt-4 border-t ${isDarkMode ? "border-white/10" : "border-slate-200"}`}>
          <div className="flex justify-between text-xs font-semibold">
            <span className={isDarkMode ? "text-slate-300" : "text-slate-700"}>
              {progressState?.current_source
                ? `Current Source: ${progressState.current_source}`
                : progressState?.message || "Processing..."}
            </span>
            <span className="font-mono text-emerald-500">{progressPct}%</span>
          </div>

          {/* Progress Bar Container */}
          <div className={`w-full h-2.5 rounded-full overflow-hidden ${isDarkMode ? "bg-slate-800" : "bg-slate-200"}`}>
            <div
              className="h-full bg-gradient-to-r from-emerald-500 via-teal-400 to-cyan-400 transition-all duration-300 ease-out"
              style={{ width: `${Math.max(5, Math.min(100, progressPct))}%` }}
            />
          </div>

          {/* Estimated Time Remaining */}
          {progressState?.estimated_remaining_secs !== undefined && (
            <div className="text-right text-xs text-emerald-500 font-mono font-bold">
              Estimated time remaining: ~{progressState.estimated_remaining_secs}s
            </div>
          )}
        </div>
      )}

      {/* Step Indicator Stepper */}
      <div className={`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-2 mt-4 pt-4 border-t ${
        isDarkMode ? "border-white/10" : "border-slate-200"
      }`}>
        {STEPS.map((stepLabel, idx) => {
          const stepNum = idx + 1;
          const isActive = currentStep === stepNum;
          const isDone = currentStep > stepNum || (!isCrawling && progressState?.is_complete);

          return (
            <div
              key={idx}
              className={`p-2.5 rounded-xl text-xs font-medium border transition-all ${
                isActive
                  ? isDarkMode
                    ? "bg-emerald-950/60 border-emerald-500 text-emerald-200 shadow-sm"
                    : "bg-emerald-50 border-emerald-500 text-emerald-900 font-bold shadow-sm"
                  : isDone
                  ? isDarkMode
                    ? "bg-slate-900 border-emerald-500/40 text-emerald-400"
                    : "bg-emerald-50/50 border-emerald-300 text-emerald-800"
                  : isDarkMode
                  ? "bg-slate-900/40 border-slate-800 text-slate-500"
                  : "bg-slate-100/80 border-slate-200 text-slate-500"
              }`}
            >
              <div className="flex items-center gap-1.5 mb-1">
                {isDone ? (
                  <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500 shrink-0" />
                ) : isActive ? (
                  <Loader2 className="w-3.5 h-3.5 text-emerald-400 animate-spin shrink-0" />
                ) : (
                  <span className={`w-4 h-4 rounded-full border text-[10px] flex items-center justify-center font-mono shrink-0 ${
                    isDarkMode ? "border-slate-700 text-slate-400" : "border-slate-300 text-slate-500"
                  }`}>
                    {stepNum}
                  </span>
                )}
                <span className="font-bold text-[11px] tracking-wide uppercase">Step {stepNum}</span>
              </div>
              <p className={`truncate text-xs font-medium ${
                isDarkMode ? "text-slate-300" : "text-slate-700"
              }`}>{stepLabel}</p>
            </div>
          );
        })}
      </div>

      {/* Toast Notification Banner */}
      {toastMessage && (
        <div
          className={`mt-4 p-3 rounded-xl border text-xs font-medium flex items-center justify-between ${
            toastMessage.type === "success"
              ? isDarkMode
                ? "bg-emerald-950/80 border-emerald-500/50 text-emerald-200"
                : "bg-emerald-50 border-emerald-300 text-emerald-900"
              : isDarkMode
              ? "bg-rose-950/80 border-rose-500/50 text-rose-200"
              : "bg-rose-50 border-rose-300 text-rose-900"
          }`}
        >
          <div className="flex items-center gap-2">
            {toastMessage.type === "success" ? (
              <CheckCircle2 className="w-4 h-4 text-emerald-500 shrink-0" />
            ) : (
              <AlertCircle className="w-4 h-4 text-rose-500 shrink-0" />
            )}
            <span>{toastMessage.text}</span>
          </div>
          <button
            onClick={onClearToast}
            className={`text-xs px-2 py-1 rounded font-bold ${
              isDarkMode ? "text-slate-400 hover:text-slate-200 bg-slate-900/60" : "text-slate-600 hover:text-slate-900 bg-slate-200/80"
            }`}
          >
            Dismiss
          </button>
        </div>
      )}
    </div>
  );
};
