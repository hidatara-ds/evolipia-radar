import { NextResponse } from "next/server";

export const runtime = "nodejs";

export async function GET() {
  // Local/dev fallback so dashboard doesn't 404.
  return NextResponse.json({
    articles_processed: 0,
    filtered_articles: 0,
    api_hits: 0,
    clusters: 0,
    avg_cluster_score: 0,
    top_cluster_titles: null,
  });
}

