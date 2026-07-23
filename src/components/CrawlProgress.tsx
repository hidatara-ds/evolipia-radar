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
}) => {
  const currentStep = progressState?.step || (isCrawling ? 1 : 0);
  const progressPct = progressState?.progress ?? (isCrawling ? 10 : 0);

  return (
    <div className="w-full bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg mb-6 text-slate-100">
      {/* Header & Controls */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
        <div>
          <h3 className="text-lg font-bold text-slate-50 flex items-center gap-2">
            {isCrawling ? (
              <Loader2 className="w-5 h-5 text-indigo-400 animate-spin" />
            ) : (
              <CheckCircle2 className="w-5 h-5 text-emerald-400" />
            )}
            Auto-Crawl Status & Real-time Progress
          </h3>
          <p className="text-xs text-slate-400">
            {isCrawling
              ? progressState?.message || "Crawler is actively running..."
              : "Auto-scheduler runs every 6 hours automatically."}
          </p>
        </div>

        <button
          onClick={onStartManualCrawl}
          disabled={isCrawling}
          className={`px-4 py-2 rounded-lg text-sm font-semibold flex items-center justify-center gap-2 transition-all ${
            isCrawling
              ? "bg-slate-800 text-slate-500 cursor-not-allowed border border-slate-700"
              : "bg-indigo-600 hover:bg-indigo-500 text-white shadow-md shadow-indigo-500/20 active:scale-95"
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
        <div className="space-y-3 mt-4 pt-4 border-t border-slate-800">
          <div className="flex justify-between text-xs text-slate-300 font-medium">
            <span>
              {progressState?.current_source
                ? `Current Source: ${progressState.current_source}`
                : progressState?.message || "Processing..."}
            </span>
            <span>{progressPct}%</span>
          </div>

          {/* Progress Bar Container */}
          <div className="w-full h-2.5 bg-slate-800 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-indigo-500 via-purple-500 to-emerald-400 transition-all duration-300 ease-out"
              style={{ width: `${Math.max(5, Math.min(100, progressPct))}%` }}
            />
          </div>

          {/* Estimated Time Remaining */}
          {progressState?.estimated_remaining_secs !== undefined && (
            <div className="text-right text-xs text-indigo-300 font-mono">
              Estimated time remaining: ~{progressState.estimated_remaining_secs}s
            </div>
          )}
        </div>
      )}

      {/* Step Indicator Stepper */}
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-2 mt-4 pt-4 border-t border-slate-800/60">
        {STEPS.map((stepLabel, idx) => {
          const stepNum = idx + 1;
          const isActive = currentStep === stepNum;
          const isDone = currentStep > stepNum || (!isCrawling && progressState?.is_complete);

          return (
            <div
              key={idx}
              className={`p-2 rounded-lg text-xs font-medium border transition-all ${
                isActive
                  ? "bg-indigo-950/80 border-indigo-500 text-indigo-200 shadow-sm"
                  : isDone
                  ? "bg-slate-900 border-emerald-500/40 text-emerald-300"
                  : "bg-slate-950/50 border-slate-800 text-slate-500"
              }`}
            >
              <div className="flex items-center gap-1.5 mb-1">
                {isDone ? (
                  <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />
                ) : isActive ? (
                  <Loader2 className="w-3.5 h-3.5 text-indigo-400 animate-spin" />
                ) : (
                  <span className="w-3.5 h-3.5 rounded-full border border-slate-700 text-[10px] flex items-center justify-center font-mono">
                    {stepNum}
                  </span>
                )}
                <span className="font-semibold">{stepNum}. Step</span>
              </div>
              <p className="truncate text-[11px] opacity-90">{stepLabel}</p>
            </div>
          );
        })}
      </div>

      {/* Toast Notification Banner */}
      {toastMessage && (
        <div
          className={`mt-4 p-3 rounded-lg border text-xs font-medium flex items-center justify-between ${
            toastMessage.type === "success"
              ? "bg-emerald-950/80 border-emerald-500/50 text-emerald-200"
              : "bg-rose-950/80 border-rose-500/50 text-rose-200"
          }`}
        >
          <div className="flex items-center gap-2">
            {toastMessage.type === "success" ? (
              <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
            ) : (
              <AlertCircle className="w-4 h-4 text-rose-400 shrink-0" />
            )}
            <span>{toastMessage.text}</span>
          </div>
          <button
            onClick={onClearToast}
            className="text-slate-400 hover:text-slate-200 text-xs px-2 py-0.5 rounded bg-slate-900/60"
          >
            Dismiss
          </button>
        </div>
      )}
    </div>
  );
};
