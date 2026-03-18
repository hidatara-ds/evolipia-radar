import { NextResponse } from "next/server";
import { loadNewsData } from "../_lib/newsData";

export const runtime = "nodejs";

export async function GET() {
  try {
    const data = await loadNewsData();
    const items = Array.isArray(data.items) ? data.items : [];
    const cutoff = Date.now() - 2 * 60 * 60 * 1000;

    const trending = items
      .filter((i) => {
        const t = Date.parse(i.published_at);
        return Number.isFinite(t) && t >= cutoff && (i.score || 0) > 0.5;
      })
      .slice(0, 20);

    return NextResponse.json({
      success: true,
      data: { items: trending, total_count: trending.length },
    });
  } catch (e) {
    return NextResponse.json(
      { success: false, error: e instanceof Error ? e.message : "Failed to load trending" },
      { status: 500 },
    );
  }
}

