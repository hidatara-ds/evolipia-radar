import { NextResponse } from "next/server";
import { loadNewsData } from "../_lib/newsData";

export const runtime = "nodejs";

export async function GET(req: Request) {
  try {
    const url = new URL(req.url);
    const q = (url.searchParams.get("q") || "").trim();
    if (!q) {
      return NextResponse.json({ success: false, error: "Query parameter 'q' is required" }, { status: 400 });
    }

    const qLower = q.toLowerCase();
    const data = await loadNewsData();
    const items = Array.isArray(data.items) ? data.items : [];

    const results = items
      .filter((i) => {
        if ((i.title || "").toLowerCase().includes(qLower)) return true;
        return (i.tags || []).some((t) => t.toLowerCase().includes(qLower));
      })
      .slice(0, 20);

    return NextResponse.json({
      success: true,
      data: { items: results, total_count: results.length, query: q },
    });
  } catch (e) {
    return NextResponse.json(
      { success: false, error: e instanceof Error ? e.message : "Failed to search" },
      { status: 500 },
    );
  }
}

