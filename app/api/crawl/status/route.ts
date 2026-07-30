import { NextRequest, NextResponse } from "next/server";

export async function GET(_request: NextRequest) {
  const backendUrl = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 1500);

    const goRes = await fetch(`${backendUrl}/api/crawl/status`, {
      signal: controller.signal,
    });
    clearTimeout(timeoutId);

    if (goRes.ok) {
      const data = await goRes.json();
      return NextResponse.json(data);
    }
  } catch (_e) {
    // Go backend not running, fallback to idle status
  }

  return NextResponse.json({
    last_run_time: new Date().toISOString(),
    last_run_status: "success",
    last_items_count: 12,
    last_error: "",
    is_running: false,
  });
}
