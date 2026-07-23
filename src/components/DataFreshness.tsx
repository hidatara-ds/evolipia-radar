"use client";

import React from "react";
import { Clock } from "lucide-react";

interface DataFreshnessProps {
  lastCrawledAt: Date | null;
}

export const DataFreshness: React.FC<DataFreshnessProps> = ({ lastCrawledAt }) => {
  if (!lastCrawledAt) {
    return (
      <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-slate-800 text-slate-400 border border-slate-700">
        <Clock className="w-3.5 h-3.5 animate-pulse" />
        Last crawled: Unknown
      </div>
    );
  }

  const now = new Date();
  const diffMs = now.getTime() - lastCrawledAt.getTime();
  const diffHours = diffMs / (1000 * 60 * 60);
  const diffMinutes = Math.floor(diffMs / (1000 * 60));

  let badgeStyle = "bg-emerald-950/80 text-emerald-300 border-emerald-500/40";
  let dotStyle = "bg-emerald-400";
  let statusText = `${diffMinutes < 1 ? "Just now" : `${diffMinutes}m ago`}`;

  if (diffHours >= 24) {
    badgeStyle = "bg-rose-950/80 text-rose-300 border-rose-500/40";
    dotStyle = "bg-rose-400";
    statusText = `${Math.floor(diffHours / 24)}d ago`;
  } else if (diffHours >= 6) {
    badgeStyle = "bg-amber-950/80 text-amber-300 border-amber-500/40";
    dotStyle = "bg-amber-400";
    statusText = `${Math.floor(diffHours)}h ago`;
  }

  return (
    <div className={`inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold border ${badgeStyle} shadow-sm`}>
      <span className={`w-2 h-2 rounded-full ${dotStyle} animate-pulse`} />
      <Clock className="w-3.5 h-3.5" />
      <span>Last crawled: {statusText}</span>
    </div>
  );
};
