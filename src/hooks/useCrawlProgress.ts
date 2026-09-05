import { useEffect, useState, useCallback } from "react";
import { triggerManualCrawl, fetchCrawlStatus } from "../api/client";

export interface CrawlProgressEvent {
  step: number;
  message: string;
  progress: number;
  current_source?: string;
  total_sources?: number;
  processed_items?: number;
  estimated_remaining_secs?: number;
  is_complete?: boolean;
  has_error?: boolean;
  error?: string;
  timestamp?: string;
}

export function useCrawlProgress() {
  const [progressState, setProgressState] = useState<CrawlProgressEvent | null>(null);
  const [isCrawling, setIsCrawling] = useState<boolean>(false);
  const [lastCrawledAt, setLastCrawledAt] = useState<Date | null>(null);
  const [toastMessage, setToastMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const checkStatus = useCallback(async () => {
    try {
      const status = await fetchCrawlStatus();
      if (status.last_run_time) {
        setLastCrawledAt(new Date(status.last_run_time));
      }
      setIsCrawling(status.is_running);
    } catch {
      // Ignore background errors
    }
  }, []);

  useEffect(() => {
    checkStatus();

    const apiBase = process.env.NEXT_PUBLIC_API_URL || "";
    const sseUrl = `${apiBase}/api/crawl/progress`;
    let eventSource: EventSource | null = null;

    try {
      eventSource = new EventSource(sseUrl);

      eventSource.addEventListener("progress", (event: MessageEvent) => {
        try {
          const data: CrawlProgressEvent = JSON.parse(event.data);
          setProgressState(data);

          if (data.progress > 0 && data.progress < 100) {
            setIsCrawling(true);
          }

          if (data.is_complete) {
            setIsCrawling(false);
            setLastCrawledAt(new Date());
            if (data.has_error) {
              setToastMessage({ type: "error", text: `Crawl failed: ${data.error || "Unknown error"}` });
            } else {
              setToastMessage({ type: "success", text: `Crawl completed! ${data.processed_items || 0} items processed.` });
            }
          }
        } catch (e) {
          console.error("Failed to parse SSE progress data", e);
        }
      });

      eventSource.onerror = () => {
        // EventSource will auto-reconnect
      };
    } catch (e) {
      console.error("SSE initialization failed", e);
    }

    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [checkStatus]);

  const startManualCrawl = async (onComplete?: () => Promise<void> | void) => {
    try {
      setIsCrawling(true);
      setProgressState({
        step: 1,
        message: "Initializing crawler...",
        progress: 10,
      });

      // Step simulation for instant visual feedback during async execution
      const stepPromise = (async () => {
        const steps = [
          { step: 2, message: "Scanning sources (HN, ArXiv, TechCrunch)...", progress: 30 },
          { step: 3, message: "Parsing content & extracting signals...", progress: 55 },
          { step: 4, message: "Validating data & deduplicating...", progress: 75 },
          { step: 5, message: "Merging signals into database & feeds...", progress: 90 },
        ];
        for (let i = 0; i < steps.length; i++) {
          await new Promise((resolve) => setTimeout(resolve, 450));
          setProgressState(steps[i]);
        }
      })();

      const [res] = await Promise.all([triggerManualCrawl(), stepPromise]);

      const newMergedCount = res?.stats?.newMerged ?? (res?.items_count ?? 0);
      const totalCount = res?.stats?.totalCount ?? 0;

      setProgressState({
        step: 6,
        message: "Done!",
        progress: 100,
        is_complete: true,
        processed_items: newMergedCount,
      });

      setIsCrawling(false);
      setLastCrawledAt(new Date());

      const successMsg =
        newMergedCount > 0
          ? `Signal crawl completed! Successfully merged ${newMergedCount} new signals (${totalCount} total).`
          : "Crawl completed! All discovered signals are already up to date.";

      setToastMessage({ type: "success", text: successMsg });

      if (onComplete) {
        try {
          await onComplete();
        } catch (reloadErr) {
          console.error("Failed to reload data after crawl:", reloadErr);
        }
      }
    } catch (err: any) {
      setIsCrawling(false);
      setToastMessage({ type: "error", text: `Failed to start crawl: ${err.message}` });
    }
  };

  const clearToast = () => setToastMessage(null);

  return {
    progressState,
    isCrawling,
    lastCrawledAt,
    toastMessage,
    startManualCrawl,
    clearToast,
    refreshStatus: checkStatus,
  };
}
