import { NextResponse } from "next/server";
import { loadNewsData, normalizeTopic } from "../_lib/newsData";

export const runtime = "nodejs";

export async function GET(req: Request) {
  try {
    const url = new URL(req.url);
    const topic = normalizeTopic(url.searchParams.get("topic"));

    const data = await loadNewsData();
    let items = Array.isArray(data.items) ? data.items : [];

    if (topic) {
      const t = topic.toLowerCase();
      items = items.filter((i) => (i.tags || []).some((tag) => tag.toLowerCase() === t));
    }

    return NextResponse.json({
      success: true,
      data: {
        items,
        total_count: items.length,
        last_updated: data.last_updated || new Date().toISOString(),
      },
    });
  } catch (e) {
    return NextResponse.json(
      {
        success: false,
        error: e instanceof Error ? e.message : "Failed to load news",
      },
      { status: 500 },
    );
  }
}

