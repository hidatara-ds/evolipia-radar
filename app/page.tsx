"use client";

import { useEffect, useState } from "react";
import { Activity, Database, FileText, BrainCircuit, RefreshCw, Zap } from "lucide-react";

interface Metrics {
  articles_processed: number;
  filtered_articles: number;
  api_hits: number;
  clusters: number;
  avg_cluster_score: number;
  top_cluster_titles: string[] | null;
}

export default function Dashboard() {
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [triggering, setTriggering] = useState(false);
  const [lastTrigger, setLastTrigger] = useState<Date | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const [baseUrl, setBaseUrl] = useState("");

  useEffect(() => {
    // If NEXT_PUBLIC_API_BASE_URL is set, use it. Otherwise assume local dev relative path.
    const url = process.env.NEXT_PUBLIC_API_BASE_URL || "";
    setBaseUrl(url);
    fetchMetrics(url);
    
    // Auto-refresh every 30s
    const interval = setInterval(() => fetchMetrics(url), 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchMetrics = async (url: string) => {
    try {
      const res = await fetch(`${url}/metrics`);
      if (res.ok) {
        const data = await res.json();
        setMetrics(data);
      }
    } catch (e) {
      console.error("Failed to fetch metrics", e);
    } finally {
      setLoading(false);
    }
  };

  const handleTrigger = async () => {
    setTriggering(true);
    setToast("Crawling...");
    try {
      const res = await fetch(`${baseUrl}/v2/crawl/trigger`, { method: "POST" });
      if (res.ok) {
        const data = await res.json();
        setLastTrigger(new Date());
        setToast(`Success! Found: ${data.stats?.discovered || 0}`);
        fetchMetrics(baseUrl);
      } else {
        setToast("Failed to trigger cycle.");
      }
    } catch (e) {
      setToast("Network error.");
    } finally {
      setTriggering(false);
      setTimeout(() => setToast(null), 3000);
    }
  };

  return (
    <main className="max-w-6xl mx-auto p-6 md:p-12">
      <header className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-12 border-b border-gray-800 pb-6">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <Zap className="text-yellow-400 w-8 h-8" />
            Evolipia Radar
          </h1>
          <p className="text-gray-400 mt-2">Zero-Cost Autonomous Intelligence Engine</p>
        </div>
        
        <div className="flex items-center gap-4">
          {toast && (
            <span className="text-sm px-3 py-1 bg-gray-800 border border-gray-700 rounded-md text-emerald-400 animate-pulse transition-all">
              {toast}
            </span>
          )}
          <button
            onClick={handleTrigger}
            disabled={triggering}
            className="flex items-center gap-2 px-5 py-2.5 bg-white text-black font-semibold rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
          >
            <RefreshCw className={`w-4 h-4 ${triggering ? 'animate-spin' : ''}`} />
            Run Crawl
          </button>
        </div>
      </header>

      {loading ? (
        <div className="animate-pulse flex flex-col gap-6">
          <div className="h-32 bg-gray-900 rounded-xl w-full"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div className="h-40 bg-gray-900 rounded-xl"></div>
            <div className="h-40 bg-gray-900 rounded-xl"></div>
            <div className="h-40 bg-gray-900 rounded-xl"></div>
            <div className="h-40 bg-gray-900 rounded-xl"></div>
          </div>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
            <MetricCard
              icon={<FileText className="text-blue-400 w-5 h-5" />}
              label="Sources Processed"
              value={metrics?.articles_processed || 0}
            />
            <MetricCard
              icon={<Database className="text-rose-400 w-5 h-5" />}
              label="Filtered (Noise)"
              value={metrics?.filtered_articles || 0}
            />
            <MetricCard
              icon={<BrainCircuit className="text-emerald-400 w-5 h-5" />}
              label="Active Clusters"
              value={metrics?.clusters || 0}
            />
            <MetricCard
              icon={<Activity className="text-amber-400 w-5 h-5" />}
              label="Avg Cluster Score"
              value={(metrics?.avg_cluster_score || 0).toFixed(1)}
            />
          </div>

          <section className="bg-gray-900/50 border border-gray-800 rounded-xl p-8 shadow-xl">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold flex items-center gap-2">
			    Activity Log
			  </h2>
              {lastTrigger && (
                <span className="text-sm text-gray-400 font-mono">
                  Last Sync: {lastTrigger.toLocaleTimeString()}
                </span>
              )}
            </div>

            {(metrics?.clusters || 0) === 0 ? (
              <div className="text-center py-16 px-4 border border-dashed border-gray-800 rounded-lg bg-black/20">
                <BrainCircuit className="w-16 h-16 text-gray-700 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-300">Semantic Engine Idle</h3>
                <p className="text-gray-500 mt-2 max-w-lg mx-auto leading-relaxed">
                  AI clustering is not currently active (Pending Phase 5 launch), or the system is operating in DRY_RUN mode. Trigger a crawl to populate the database.
                </p>
              </div>
            ) : (
              <div className="space-y-4">
                {metrics?.top_cluster_titles?.map((title, i) => (
                  <div key={i} className="p-4 bg-black/40 rounded-lg border border-gray-800 flex items-center gap-4 hover:border-gray-700 transition-colors">
                    <span className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-500/10 text-blue-400 flex items-center justify-center font-bold text-sm">
                      {i + 1}
                    </span>
                    <span className="text-gray-200 font-medium">{title}</span>
                  </div>
                ))}
              </div>
            )}
            
            <div className="mt-8 border-t border-gray-800 pt-6 text-center">
              <span className="inline-block px-3 py-1 bg-indigo-500/10 text-indigo-400 text-xs font-medium tracking-wide uppercase rounded-full">
                Trending Clusters (Coming Soon in Phase 5)
              </span>
            </div>
          </section>

          <footer className="mt-16 text-center text-sm text-gray-600 font-mono">
            Phase 4.5 • Powered by Vercel Serverless
          </footer>
        </>
      )}
    </main>
  );
}

function MetricCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: string | number }) {
  return (
    <div className="bg-gray-900/40 border border-gray-800 p-6 rounded-xl flex flex-col gap-4 shadow-sm hover:border-gray-700 transition-colors">
      <div className="flex items-center gap-3">
        <div className="p-2.5 bg-black rounded-lg border border-gray-800">
          {icon}
        </div>
        <h3 className="text-sm font-medium text-gray-400 tracking-wide">{label}</h3>
      </div>
      <p className="text-4xl font-semibold tracking-tight">{value}</p>
    </div>
  );
}
