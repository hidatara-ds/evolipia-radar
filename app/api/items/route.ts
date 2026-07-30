import { NextRequest, NextResponse } from "next/server";
import { Pool } from "pg";
import fs from "fs";
import path from "path";

interface RawNewsItem {
  id: string;
  title: string;
  url: string;
  domain: string;
  published_at: string;
  category: string;
  score?: number;
  raw_score?: number;
  tldr?: string;
  why_it_matters?: string;
  tags?: string[];
  novelty?: number;
  impact?: number;
  engineering_value?: number;
  reasoning?: string;
}

let pool: Pool | null = null;

function getDBPool(): Pool | null {
  const dbUrl = process.env.DATABASE_URL;
  if (!dbUrl) return null;

  if (!pool) {
    pool = new Pool({
      connectionString: dbUrl,
      ssl: dbUrl.includes("sslmode=disable")
        ? false
        : { rejectUnauthorized: false },
      max: 5,
      idleTimeoutMillis: 30000,
      connectionTimeoutMillis: 5000,
    });
  }
  return pool;
}

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);

  // 1. Try proxying to Go backend if running on localhost:8080
  const backendUrl = process.env.BACKEND_URL || "http://localhost:8080";
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 1200);

    const goApiUrl = `${backendUrl}/api/items?${searchParams.toString()}`;
    const goRes = await fetch(goApiUrl, { signal: controller.signal });
    clearTimeout(timeoutId);

    if (goRes.ok) {
      const data = await goRes.json();
      return NextResponse.json(data);
    }
  } catch (_e) {
    // Go backend not running, proceed to direct Neon DB / PostgreSQL query
  }

  // 2. Direct Neon DB / PostgreSQL connection using DATABASE_URL
  const dbPool = getDBPool();
  if (dbPool) {
    try {
      const search = (searchParams.get("search") || "").toLowerCase().trim();
      const dateFrom = searchParams.get("date_from");
      const dateTo = searchParams.get("date_to");
      const minRelevance = Number(searchParams.get("min_relevance")) || 0;
      const status = searchParams.get("status");
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

      const whereClauses: string[] = [];
      const queryParams: any[] = [];
      let paramIdx = 1;

      if (search) {
        whereClauses.push(`(LOWER(i.title) LIKE $${paramIdx} OR LOWER(i.domain) LIKE $${paramIdx} OR LOWER(COALESCE(i.raw_excerpt, '')) LIKE $${paramIdx})`);
        queryParams.push(`%${search}%`);
        paramIdx++;
      }

      if (dateFrom) {
        whereClauses.push(`i.published_at >= $${paramIdx}`);
        queryParams.push(new Date(dateFrom));
        paramIdx++;
      }

      if (dateTo) {
        const toDate = new Date(dateTo);
        toDate.setDate(toDate.getDate() + 1);
        whereClauses.push(`i.published_at <= $${paramIdx}`);
        queryParams.push(toDate);
        paramIdx++;
      }

      if (minRelevance > 0) {
        whereClauses.push(`COALESCE(i.relevance_score, ROUND(COALESCE(sc.final, 0.85) * 100)) >= $${paramIdx}`);
        queryParams.push(minRelevance);
        paramIdx++;
      }

      if (status && status !== "all") {
        whereClauses.push(`COALESCE(i.crawl_status, 'done') = $${paramIdx}`);
        queryParams.push(status);
        paramIdx++;
      }

      if (sources.length > 0) {
        const srcPlaceholders = sources.map((_, idx) => `$${paramIdx + idx}`).join(", ");
        whereClauses.push(`(LOWER(s.name) IN (${srcPlaceholders}) OR LOWER(i.domain) IN (${srcPlaceholders}))`);
        sources.forEach(s => queryParams.push(s));
        paramIdx += sources.length;
      }

      if (categories.length > 0) {
        const catPlaceholders = categories.map((_, idx) => `$${paramIdx + idx}`).join(", ");
        whereClauses.push(`LOWER(i.category) IN (${catPlaceholders})`);
        categories.forEach(c => queryParams.push(c));
        paramIdx += categories.length;
      }

      const whereSQL = whereClauses.length > 0 ? `WHERE ${whereClauses.join(" AND ")}` : "";

      // Total count in items table
      const totalCountRes = await dbPool.query("SELECT COUNT(*) FROM items");
      const totalCount = parseInt(totalCountRes.rows[0]?.count || "0", 10);

      // Filtered count
      const filteredCountQuery = `
        SELECT COUNT(*) FROM items i 
        LEFT JOIN sources s ON i.source_id = s.id 
        LEFT JOIN scores sc ON i.id = sc.item_id 
        ${whereSQL}
      `;
      const filteredCountRes = await dbPool.query(filteredCountQuery, queryParams);
      const filteredCount = parseInt(filteredCountRes.rows[0]?.count || "0", 10);

      // Sorting
      let orderCol = "i.published_at";
      if (sortBy === "relevance") {
        orderCol = "COALESCE(i.relevance_score, ROUND(COALESCE(sc.final, 0.85) * 100))";
      } else if (sortBy === "impact") {
        orderCol = "COALESCE(sc.impact, 8.5)";
      }
      const dir = sortOrder.toUpperCase() === "ASC" ? "ASC" : "DESC";

      const offset = (page - 1) * limit;
      const dataQuery = `
        SELECT 
          i.id, i.source_id, COALESCE(s.name, 'Global Source') as source_name, i.title, i.url, 
          i.published_at, i.domain, i.category, i.raw_excerpt, 
          COALESCE(i.crawl_status, 'done') as crawl_status, i.crawl_error,
          COALESCE(i.relevance_score, ROUND(COALESCE(sc.final, 0.85) * 100)) as relevance_score,
          COALESCE(sc.impact, 8.5) as impact,
          COALESCE(sc.engineering_value, 8.5) as engineering_value,
          COALESCE(sc.reasoning, '') as reasoning,
          COALESCE(sm.tldr, i.raw_excerpt, i.title) as tldr,
          COALESCE(sm.why_it_matters, '') as why_it_matters,
          COALESCE(sm.tags, '[]'::jsonb) as tags,
          i.created_at
        FROM items i
        LEFT JOIN sources s ON i.source_id = s.id
        LEFT JOIN scores sc ON i.id = sc.item_id
        LEFT JOIN summaries sm ON i.id = sm.item_id
        ${whereSQL}
        ORDER BY ${orderCol} ${dir}
        LIMIT $${paramIdx} OFFSET $${paramIdx + 1}
      `;

      const dataRes = await dbPool.query(dataQuery, [...queryParams, limit, offset]);

      const items = dataRes.rows.map(row => ({
        id: row.id,
        source_id: row.source_id,
        source_name: row.source_name,
        title: row.title,
        url: row.url,
        published_at: row.published_at ? new Date(row.published_at).toISOString() : new Date().toISOString(),
        domain: row.domain,
        category: row.category || "llm",
        raw_excerpt: row.raw_excerpt || row.tldr || row.title,
        tldr: row.tldr,
        why_it_matters: row.why_it_matters,
        tags: Array.isArray(row.tags) ? row.tags : (typeof row.tags === "string" ? JSON.parse(row.tags) : []),
        crawl_status: row.crawl_status,
        crawl_error: row.crawl_error,
        relevance_score: Number(row.relevance_score) || 85,
        scaled_score: Number(((Number(row.relevance_score) || 85) / 10).toFixed(1)),
        created_at: row.created_at ? new Date(row.created_at).toISOString() : new Date().toISOString(),
        impact: Number(row.impact) || 8.5,
        engineering_value: Number(row.engineering_value) || 8.5,
        reasoning: row.reasoning,
      }));

      const totalPages = Math.max(1, Math.ceil(filteredCount / limit));

      return NextResponse.json({
        success: true,
        data: items,
        total_count: totalCount,
        filtered_count: filteredCount,
        page,
        total_pages: totalPages,
        last_updated: new Date().toISOString(),
      });
    } catch (err: any) {
      console.error("Direct Neon DB Query Failed:", err);
      // Fallthrough to JSON fallback below
    }
  }

  // 3. Fallback to reading api/news.json or data/news.json
  try {
    let filePath = path.join(process.cwd(), "api", "news.json");
    if (!fs.existsSync(filePath)) {
      filePath = path.join(process.cwd(), "data", "news.json");
    }

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

    let filtered = items.map(item => {
      const scaledScore = Math.round((item.score || item.raw_score || 0.8) * (item.score && item.score > 1 ? 1 : 100));
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
        scaled_score: Number((scaledScore / 10).toFixed(1)),
        created_at: item.published_at,
        impact: item.impact ? Math.round(item.impact * 10) : 85,
        engineering_value: item.engineering_value ? Math.round(item.engineering_value * 10) : 85,
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
