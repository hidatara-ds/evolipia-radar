import { NextRequest, NextResponse } from "next/server";
import { executeLiveCrawlAndMerge } from "@/src/lib/crawler";

export async function POST(_request: NextRequest) {
  const backendUrl = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  
  // 1. Try forwarding to Go backend if running
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 2000);

    const goRes = await fetch(`${backendUrl}/api/crawl`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      signal: controller.signal,
    });
    clearTimeout(timeoutId);

    if (goRes.ok) {
      const data = await goRes.json();
      return NextResponse.json(data, { status: 200 });
    }
  } catch (_e) {
    // Go backend not running, proceed to internal TypeScript crawler
  }

  // 2. Execute internal live feed crawl and smart merge
  try {
    const result = await executeLiveCrawlAndMerge();
    return NextResponse.json({
      success: true,
      status: "completed",
      message: `Crawl completed! Discovered ${result.discovered} candidates, merged ${result.newMerged} new signals.`,
      stats: result,
      timestamp: result.lastUpdated,
    });
  } catch (err: any) {
    console.error("[Crawl Route] Live crawl failed:", err);
    return NextResponse.json(
      {
        success: false,
        status: "failed",
        error: err.message || "Failed to execute live crawler",
      },
      { status: 500 }
    );
  }
}

export async function GET(request: NextRequest) {
  return POST(request);
}

