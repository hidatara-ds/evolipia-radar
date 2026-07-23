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

  const startManualCrawl = async () => {
    try {
      setIsCrawling(true);
      setProgressState({
        step: 1,
        message: "Initializing crawler...",
        progress: 5,
      });
      await triggerManualCrawl();
      setToastMessage({ type: "success", text: "Manual crawl triggered!" });
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
