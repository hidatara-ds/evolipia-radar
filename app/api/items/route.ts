import { NextRequest, NextResponse } from "next/server";
import fs from "fs";
import path from "path";

interface RawNewsItem {
  id: string;
  title: string;
  url: string;
  domain: string;
  published_at: string;
  category: string;
  score: number;
  tldr?: string;
  why_it_matters?: string;
  tags?: string[];
  novelty?: number;
  impact?: number;
  engineering_value?: number;
  reasoning?: string;
}

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);

  // Try proxying to local Go backend if available
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 1500);

    const goApiUrl = `http://localhost:8080/api/items?${searchParams.toString()}`;
    const goRes = await fetch(goApiUrl, { signal: controller.signal });
    clearTimeout(timeoutId);

    if (goRes.ok) {
      const data = await goRes.json();
      return NextResponse.json(data);
    }
  } catch (_e) {
    // Go backend not running, proceed to JSON fallback below
  }

  // Fallback to reading api/news.json
  try {
    const filePath = path.join(process.cwd(), "api", "news.json");
    if (!fs.existsSync(filePath)) {
      return NextResponse.json({
        success: true,
        data: [],
        total_count: 0,
        filtered_count: 0,
        page: 1,
        total_pages: 1,
        last_updated: new Date().toISOString(),
      });
    }

    const fileContent = fs.readFileSync(filePath, "utf-8");
    const json = JSON.parse(fileContent);
    let items: RawNewsItem[] = json.items || [];

    // Query params
    const search = (searchParams.get("search") || "").toLowerCase().trim();
    const dateFrom = searchParams.get("date_from");
    const dateTo = searchParams.get("date_to");
    const minRelevance = Number(searchParams.get("min_relevance")) || 0;
    const sortBy = searchParams.get("sort_by") || "date";
    const sortOrder = searchParams.get("sort_order") || "desc";
    const page = Math.max(1, Number(searchParams.get("page")) || 1);
    const limit = Math.max(1, Math.min(100, Number(searchParams.get("limit")) || 20));

    const sources = searchParams.getAll("sources[]").concat(
      searchParams.get("sources") ? (searchParams.get("sources") as string).split(",") : []
    ).map(s => s.toLowerCase().trim()).filter(Boolean);

    const categories = searchParams.getAll("categories[]").concat(
      searchParams.get("categories") ? (searchParams.get("categories") as string).split(",") : []
    ).map(c => c.toLowerCase().trim()).filter(Boolean);

    // Apply filtering
    let filtered = items.map(item => {
      const scaledScore = Math.round((item.score || 0.8) * 100);
      return {
        id: item.id,
        source_id: item.id,
        source_name: item.domain === "github.com" ? "GitHub Trending" : item.domain.includes("arxiv") ? "ArXiv AI" : item.domain,
        title: item.title,
        url: item.url,
        published_at: item.published_at,
        domain: item.domain,
        category: item.category || "llm",
        raw_excerpt: item.tldr || item.title,
        tldr: item.tldr,
        why_it_matters: item.why_it_matters,
        tags: item.tags || [],
        crawl_status: "done",
        relevance_score: scaledScore,
        scaled_score: Number(((item.score || 0.8) * 10).toFixed(1)),
        created_at: item.published_at,
        impact: item.impact ? Math.round(item.impact * 10) : undefined,
        engineering_value: item.engineering_value ? Math.round(item.engineering_value * 10) : undefined,
        reasoning: item.reasoning,
      };
    });

    if (search) {
      filtered = filtered.filter(it => 
        it.title.toLowerCase().includes(search) || 
        it.domain.toLowerCase().includes(search) ||
        (it.raw_excerpt && it.raw_excerpt.toLowerCase().includes(search))
      );
    }

    if (dateFrom) {
      const fromTime = new Date(dateFrom).getTime();
      if (!isNaN(fromTime)) {
        filtered = filtered.filter(it => new Date(it.published_at).getTime() >= fromTime);
      }
    }

    if (dateTo) {
      const toTime = new Date(dateTo).getTime() + 24 * 60 * 60 * 1000;
      if (!isNaN(toTime)) {
        filtered = filtered.filter(it => new Date(it.published_at).getTime() <= toTime);
      }
    }

    if (minRelevance > 0) {
      filtered = filtered.filter(it => (it.relevance_score || 0) >= minRelevance);
    }

    if (sources.length > 0) {
      filtered = filtered.filter(it => 
        sources.some(s => it.source_name.toLowerCase().includes(s) || it.domain.toLowerCase().includes(s))
      );
    }

    if (categories.length > 0) {
      filtered = filtered.filter(it => 
        categories.some(c => 
          it.category.toLowerCase().includes(c) || 
          it.tags.some(t => t.toLowerCase().includes(c))
        )
      );
    }

    // Sorting
    filtered.sort((a, b) => {
      let valA: any = a.published_at;
      let valB: any = b.published_at;

      if (sortBy === "relevance") {
        valA = a.relevance_score;
        valB = b.relevance_score;
      } else if (sortBy === "impact") {
        valA = a.impact || a.relevance_score;
        valB = b.impact || b.relevance_score;
      } else {
        valA = new Date(a.published_at).getTime();
        valB = new Date(b.published_at).getTime();
      }

      if (sortOrder === "asc") {
        return valA > valB ? 1 : valA < valB ? -1 : 0;
      } else {
        return valA < valB ? 1 : valA > valB ? -1 : 0;
      }
    });

    const totalCount = items.length;
    const filteredCount = filtered.length;
    const totalPages = Math.max(1, Math.ceil(filteredCount / limit));
    const paginatedItems = filtered.slice((page - 1) * limit, page * limit);

    return NextResponse.json({
      success: true,
      data: paginatedItems,
      total_count: totalCount,
      filtered_count: filteredCount,
      page,
      total_pages: totalPages,
      last_updated: new Date().toISOString(),
    });
  } catch (err: any) {
    return NextResponse.json(
      { success: false, error: err.message || "Failed to load fallback items" },
      { status: 500 }
    );
  }
}
