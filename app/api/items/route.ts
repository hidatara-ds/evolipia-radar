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

function getEffectiveSourceName(rawSourceName?: string, domain?: string): string {
  const dom = (domain || "").toLowerCase().trim();
  const src = (rawSourceName || "").toLowerCase().trim();

  // 1. Check specific domains first to accurately map article source badge
  if (dom.includes("ycombinator.com") || dom.includes("news.ycombinator")) {
    return "Hacker News";
  }
  if (dom.includes("arxiv")) {
    return "ArXiv AI";
  }
  if (dom.includes("techcrunch")) {
    return "TechCrunch AI";
  }
  if (dom.includes("reddit")) {
    return "Reddit MachineLearning";
  }
  if (dom.includes("twitter") || dom.includes("x.com")) {
    return "Twitter / X";
  }
  if (dom.includes("github")) {
    return "GitHub Trending";
  }
  if (dom.includes("openai.com")) {
    return "OpenAI Blog";
  }
  if (dom.includes("deepmind")) {
    return "DeepMind Research";
  }
  if (dom.includes("stability.ai")) {
    return "Stability AI";
  }
  if (dom.includes("anthropic.com")) {
    return "Anthropic";
  }

  // 2. Fall back to raw source name matching if domain isn't a specific known domain
  if (src.includes("hacker") || src.includes("ycombinator")) {
    return "Hacker News";
  }
  if (src.includes("arxiv")) {
    return "ArXiv AI";
  }
  if (src.includes("techcrunch")) {
    return "TechCrunch AI";
  }
  if (src.includes("reddit")) {
    return "Reddit MachineLearning";
  }
  if (src.includes("twitter") || src.includes("x.com")) {
    return "Twitter / X";
  }
  if (src.includes("github")) {
    return "GitHub Trending";
  }
  if (src.includes("openai")) {
    return "OpenAI Blog";
  }
  if (src.includes("deepmind")) {
    return "DeepMind Research";
  }
  if (src.includes("stability")) {
    return "Stability AI";
  }
  if (src.includes("anthropic")) {
    return "Anthropic";
  }

  if (rawSourceName && rawSourceName !== "Global Source" && rawSourceName !== "Unknown" && rawSourceName !== "RSSAgent") {
    return rawSourceName;
  }
  if (domain) {
    return domain;
  }
  return "Global Source";
}

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);

  // 1. Try proxying to local Go backend if available
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

  // 2. Direct Neon DB / PostgreSQL connection using exact Neon schema
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
        whereClauses.push(`(LOWER(i.title) LIKE $${paramIdx} OR LOWER(i.domain) LIKE $${paramIdx} OR LOWER(COALESCE(i.raw_excerpt, '')) LIKE $${paramIdx} OR LOWER(COALESCE(sm.tldr, '')) LIKE $${paramIdx})`);
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
        whereClauses.push(`ROUND(COALESCE(sc.relevance, sc.final, 0.85) * (CASE WHEN COALESCE(sc.relevance, sc.final, 0.85) <= 1.0 THEN 100 ELSE 1 END)) >= $${paramIdx}`);
        queryParams.push(minRelevance);
        paramIdx++;
      }

      if (sources.length > 0) {
        const sourceOrClauses: string[] = [];
        sources.forEach(src => {
          if (src.includes("hacker") || src.includes("ycombinator")) {
            sourceOrClauses.push(`(LOWER(s.name) LIKE '%hacker%' OR LOWER(s.name) LIKE '%ycombinator%' OR LOWER(i.domain) LIKE '%ycombinator%')`);
          } else if (src.includes("arxiv")) {
            sourceOrClauses.push(`(LOWER(s.name) LIKE '%arxiv%' OR LOWER(i.domain) LIKE '%arxiv%')`);
          } else if (src.includes("techcrunch")) {
            sourceOrClauses.push(`(LOWER(s.name) LIKE '%techcrunch%' OR LOWER(i.domain) LIKE '%techcrunch%')`);
          } else if (src.includes("reddit")) {
            sourceOrClauses.push(`(LOWER(s.name) LIKE '%reddit%' OR LOWER(i.domain) LIKE '%reddit%')`);
          } else if (src.includes("twitter") || src.includes("x")) {
            sourceOrClauses.push(`(LOWER(s.name) LIKE '%twitter%' OR LOWER(i.domain) LIKE '%twitter%' OR LOWER(i.domain) LIKE '%x.com%')`);
          } else if (src.includes("github")) {
            sourceOrClauses.push(`((LOWER(s.name) LIKE '%github%' OR LOWER(i.domain) LIKE '%github%') AND LOWER(COALESCE(s.name, '')) NOT LIKE '%hacker%' AND LOWER(i.domain) NOT LIKE '%ycombinator%')`);
          } else {
            sourceOrClauses.push(`(LOWER(s.name) LIKE $${paramIdx} OR LOWER(i.domain) LIKE $${paramIdx})`);
            queryParams.push(`%${src}%`);
            paramIdx++;
          }
        });
        if (sourceOrClauses.length > 0) {
          whereClauses.push(`(${sourceOrClauses.join(" OR ")})`);
        }
      }

      if (categories.length > 0) {
        const catOrClauses: string[] = [];
        categories.forEach(c => {
          const lowerC = c.toLowerCase().trim();
          let catCondition = `(LOWER(i.category) LIKE $${paramIdx} OR LOWER(COALESCE(i.title, '')) LIKE $${paramIdx} OR LOWER(COALESCE(i.raw_excerpt, '')) LIKE $${paramIdx} OR LOWER(COALESCE(sm.tldr, '')) LIKE $${paramIdx} OR EXISTS (SELECT 1 FROM jsonb_array_elements_text(COALESCE(sm.tags, '[]'::jsonb)) elem WHERE LOWER(elem) LIKE $${paramIdx}))`;

          if (lowerC === "llm") {
            catCondition += ` OR LOWER(i.category) IN ('research', 'news', 'general', 'models')`;
          } else if (lowerC === "open-source") {
            catCondition += ` OR LOWER(i.category) IN ('tools', 'github', 'models') OR LOWER(i.domain) LIKE '%github%'`;
          } else if (lowerC === "infra") {
            catCondition += ` OR LOWER(i.category) IN ('tools', 'infra', 'benchmarks')`;
          } else if (lowerC === "agents") {
            catCondition += ` OR LOWER(i.category) IN ('research', 'tools', 'news')`;
          }

          catOrClauses.push(`(${catCondition})`);
          queryParams.push(`%${lowerC}%`);
          paramIdx++;
        });
        if (catOrClauses.length > 0) {
          whereClauses.push(`(${catOrClauses.join(" OR ")})`);
        }
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
        LEFT JOIN summaries sm ON i.id = sm.item_id
        ${whereSQL}
      `;
      const filteredCountRes = await dbPool.query(filteredCountQuery, queryParams);
      const filteredCount = parseInt(filteredCountRes.rows[0]?.count || "0", 10);

      // Sorting
      let orderCol = "i.published_at";
      if (sortBy === "relevance") {
        orderCol = "COALESCE(sc.relevance, sc.final, 0.85)";
      } else if (sortBy === "impact") {
        orderCol = "COALESCE(sc.impact, 0.85)";
      }
      const dir = sortOrder.toUpperCase() === "ASC" ? "ASC" : "DESC";

      const offset = (page - 1) * limit;
      const dataQuery = `
        SELECT 
          i.id, 
          i.source_id, 
          COALESCE(s.name, 'Global Source') as source_name, 
          i.title, 
          i.url, 
          i.published_at, 
          i.domain, 
          i.category, 
          i.raw_excerpt, 
          'done' as crawl_status, 
          NULL as crawl_error,
          ROUND(COALESCE(sc.relevance, sc.final, 0.85) * (CASE WHEN COALESCE(sc.relevance, sc.final, 0.85) <= 1.0 THEN 100 ELSE 1 END)) as relevance_score,
          COALESCE(sc.impact, 0.85) as impact,
          COALESCE(sc.engineering_value, 0.85) as engineering_value,
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

      const items = dataRes.rows.map(row => {
        const rawRel = Number(row.relevance_score);
        const relScore = !isNaN(rawRel) && rawRel > 0 ? Math.round(rawRel) : 85;
        const rawImp = Number(row.impact);
        const impactScore = !isNaN(rawImp) ? (rawImp <= 1.0 ? Math.round(rawImp * 100) : Math.round(rawImp)) : 85;

        const effSource = getEffectiveSourceName(row.source_name, row.domain);

        return {
          id: row.id,
          source_id: row.source_id,
          source_name: effSource,
          title: row.title,
          url: row.url,
          published_at: row.published_at ? new Date(row.published_at).toISOString() : new Date().toISOString(),
          domain: row.domain,
          category: row.category || "llm",
          raw_excerpt: row.raw_excerpt || row.tldr || row.title,
          tldr: row.tldr,
          why_it_matters: row.why_it_matters,
          tags: Array.isArray(row.tags) ? row.tags : (typeof row.tags === "string" ? JSON.parse(row.tags) : []),
          crawl_status: "done",
          crawl_error: null,
          relevance_score: relScore,
          scaled_score: Number((relScore / 10).toFixed(1)),
          created_at: row.created_at ? new Date(row.created_at).toISOString() : new Date().toISOString(),
          impact: impactScore,
          engineering_value: Number(row.engineering_value) || 85,
          reasoning: row.reasoning,
        };
      });

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
      const effSource = getEffectiveSourceName(item.domain === "github.com" ? "GitHub Trending" : item.domain.includes("arxiv") ? "ArXiv AI" : (item as any).source_name || "", item.domain);

      return {
        id: item.id,
        source_id: item.id,
        source_name: effSource,
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
      filtered = filtered.filter(it => {
        const itemSrc = (it.source_name || "").toLowerCase();
        return sources.some(src => itemSrc.includes(src) || src.includes(itemSrc));
      });
    }

    if (categories.length > 0) {
      filtered = filtered.filter(it =>
        categories.some(c => {
          const lowerC = c.toLowerCase().trim();
          const cat = (it.category || "").toLowerCase();
          const title = (it.title || "").toLowerCase();
          const excerpt = (it.raw_excerpt || it.tldr || "").toLowerCase();
          const tags = (it.tags || []).map((t: string) => t.toLowerCase());

          if (cat.includes(lowerC) || title.includes(lowerC) || excerpt.includes(lowerC) || tags.some((t: string) => t.includes(lowerC))) {
            return true;
          }
          if (lowerC === "llm" && (cat === "research" || cat === "news" || cat === "general" || cat === "models")) return true;
          if (lowerC === "open-source" && (cat === "tools" || cat === "github" || (it.domain || "").includes("github"))) return true;
          if (lowerC === "infra" && (cat === "tools" || cat === "infra" || cat === "benchmarks")) return true;
          if (lowerC === "agents" && (cat === "research" || cat === "tools" || cat === "news")) return true;
          return false;
        })
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
