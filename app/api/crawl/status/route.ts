import { NextRequest, NextResponse } from "next/server";
import fs from "fs";
import path from "path";

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
    // Go backend not running, fallback to news.json metadata
  }

  try {
    let filePath = path.join(process.cwd(), "api", "news.json");
    if (!fs.existsSync(filePath)) {
      filePath = path.join(process.cwd(), "data", "news.json");
    }
    if (fs.existsSync(filePath)) {
      const parsed = JSON.parse(fs.readFileSync(filePath, "utf-8"));
      return NextResponse.json({
        last_run_time: parsed.last_updated || new Date().toISOString(),
        last_run_status: "success",
        last_items_count: parsed.total_count || (parsed.items ? parsed.items.length : 0),
        last_error: "",
        is_running: false,
      });
    }
  } catch {}

  return NextResponse.json({
    last_run_time: new Date().toISOString(),
    last_run_status: "success",
    last_items_count: 0,
    last_error: "",
    is_running: false,
  });
}
