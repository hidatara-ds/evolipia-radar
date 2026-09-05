import fs from "fs";
import path from "path";
import { Pool } from "pg";

export interface CrawledArticle {
  id: string;
  title: string;
  url: string;
  domain: string;
  published_at: string;
  category: string;
  score: number;
  tldr: string;
  why_it_matters: string;
  tags: string[];
  impact?: number;
  engineering_value?: number;
  reasoning?: string;
}

export interface NewsDataFile {
  items: CrawledArticle[];
  last_updated: string;
  total_count: number;
}

// Normalize URLs to avoid duplicate entries with query tracking params
export function normalizeUrl(rawUrl: string): string {
  try {
    const parsed = new URL(rawUrl.trim());
    parsed.hash = "";
    // Remove analytics tracking parameters
    const paramsToDelete: string[] = [];
    parsed.searchParams.forEach((_, key) => {
      if (
        key.startsWith("utm_") ||
        key === "ref" ||
        key === "source" ||
        key === "fbclid" ||
        key === "gclid"
      ) {
        paramsToDelete.push(key);
      }
    });
    paramsToDelete.forEach((k) => parsed.searchParams.delete(k));

    let clean = parsed.toString();
    if (clean.endsWith("/") && parsed.pathname !== "/") {
      clean = clean.slice(0, -1);
    }
    return clean;
  } catch {
    return rawUrl.trim();
  }
}

export function extractDomain(rawUrl: string): string {
  try {
    const parsed = new URL(rawUrl);
    return parsed.hostname.replace(/^www\./, "");
  } catch {
    return "external";
  }
}

function cleanHtmlText(raw: string): string {
  if (!raw) return "";
  return raw
    .replace(/<!\[CDATA\[([\s\S]*?)\]\]>/g, "$1")
    .replace(/<[^>]+>/g, " ")
    .replace(/&quot;/g, '"')
    .replace(/&apos;/g, "'")
    .replace(/&#39;/g, "'")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/\s+/g, " ")
    .trim();
}

function inferCategoryAndTags(title: string, excerpt: string, domain: string): { category: string; tags: string[] } {
  const combined = `${title} ${excerpt} ${domain}`.toLowerCase();
  const tagsSet = new Set<string>();

  let category = "llm";

  if (combined.includes("arxiv") || combined.includes("paper") || combined.includes("benchmark") || combined.includes("survey")) {
    tagsSet.add("research");
    category = "research";
  }
  if (combined.includes("agent") || combined.includes("autonomous") || combined.includes("workflow") || combined.includes("autogpt") || combined.includes("swarm")) {
    tagsSet.add("agents");
    category = "agents";
  }
  if (combined.includes("vision") || combined.includes("dall-e") || combined.includes("diffusion") || combined.includes("image") || combined.includes("video generation") || combined.includes("flux")) {
    tagsSet.add("vision");
    if (category === "llm") category = "vision";
  }
  if (combined.includes("robot") || combined.includes("robotics") || combined.includes("atlas") || combined.includes("optimus") || combined.includes("embodied")) {
    tagsSet.add("robotics");
    category = "robotics";
  }
  if (combined.includes("open-source") || combined.includes("open source") || combined.includes("github") || combined.includes("huggingface") || combined.includes("weights")) {
    tagsSet.add("open-source");
    if (category === "llm") category = "open-source";
  }
  if (combined.includes("infra") || combined.includes("postgres") || combined.includes("rag") || combined.includes("vector") || combined.includes("cuda") || combined.includes("inference") || combined.includes("vllm") || combined.includes("latency")) {
    tagsSet.add("infra");
    if (category === "llm") category = "infra";
  }
  if (combined.includes("security") || combined.includes("safety") || combined.includes("jailbreak") || combined.includes("alignment") || combined.includes("eval")) {
    tagsSet.add("security");
    if (category === "llm") category = "security";
  }

  // Ensure primary category is tagged
  tagsSet.add(category);
  if (combined.includes("llm") || combined.includes("gpt") || combined.includes("claude") || combined.includes("gemini") || combined.includes("deepseek")) {
    tagsSet.add("llm");
  }

  return { category, tags: Array.from(tagsSet) };
}

function generateUuid(): string {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

// Fetch RSS XML feed and parse items using lightweight parser
async function fetchRssFeed(feedUrl: string, sourceName: string, defaultDomain: string): Promise<CrawledArticle[]> {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 6000);

    const res = await fetch(feedUrl, {
      signal: controller.signal,
      headers: {
        "User-Agent": "Mozilla/5.0 (compatible; EvolipiaRadar/1.0; +https://evolipia.ai)",
        Accept: "application/rss+xml, application/xml, text/xml;q=0.9, */*;q=0.8",
      },
    });
    clearTimeout(timeout);

    if (!res.ok) return [];
    const xml = await res.text();

    const articles: CrawledArticle[] = [];
    const itemRegex = /<item>([\s\S]*?)<\/item>/gi;
    let match: RegExpExecArray | null;

    while ((match = itemRegex.exec(xml)) !== null && articles.length < 12) {
      const itemBlock = match[1];

      const titleMatch = itemBlock.match(/<title>([\s\S]*?)<\/title>/i);
      const linkMatch = itemBlock.match(/<link>([\s\S]*?)<\/link>/i);
      const descMatch = itemBlock.match(/<description>([\s\S]*?)<\/description>/i);
      const pubDateMatch = itemBlock.match(/<pubDate>([\s\S]*?)<\/pubDate>/i);

      if (!titleMatch || !linkMatch) continue;

      const rawTitle = cleanHtmlText(titleMatch[1]);
      let rawLink = cleanHtmlText(linkMatch[1]);
      const rawDesc = cleanHtmlText(descMatch ? descMatch[1] : "");

      if (!rawTitle || !rawLink) continue;

      // Unescape or clean link if needed
      if (rawLink.startsWith("http") === false) {
        continue;
      }

      const normalized = normalizeUrl(rawLink);
      const domain = extractDomain(normalized) || defaultDomain;

      let publishedAt = new Date().toISOString();
      if (pubDateMatch) {
        const parsedDate = new Date(pubDateMatch[1].trim());
        if (!isNaN(parsedDate.getTime())) {
          publishedAt = parsedDate.toISOString();
        }
      }

      const { category, tags } = inferCategoryAndTags(rawTitle, rawDesc, domain);

      // Score calculation
      const baseScore = 0.85;
      const scoreVariance = (rawTitle.length % 12) * 0.01;
      const finalScore = Math.min(0.97, Math.max(0.78, Number((baseScore + scoreVariance).toFixed(2))));

      articles.push({
        id: generateUuid(),
        title: rawTitle,
        url: normalized,
        domain,
        published_at: publishedAt,
        category,
        score: finalScore,
        tldr: rawDesc.length > 20 ? (rawDesc.length > 200 ? rawDesc.slice(0, 197) + "..." : rawDesc) : `${rawTitle}. Latest update from ${sourceName}.`,
        why_it_matters: `High signal AI development regarding ${category} and practical deployments.`,
        tags,
        impact: Math.round(finalScore * 10),
        engineering_value: Math.round(finalScore * 10),
        reasoning: `Extracted from ${sourceName} live RSS stream.`,
      });
    }

    return articles;
  } catch (e) {
    console.warn(`[Crawler] Failed to fetch RSS feed ${feedUrl}:`, e);
    return [];
  }
}

// Fetch trending stories from Algolia HackerNews API
async function fetchHnTrending(): Promise<CrawledArticle[]> {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 6000);

    const res = await fetch("https://hn.algolia.com/api/v1/search?tags=front_page&hitsPerPage=20", {
      signal: controller.signal,
      headers: {
        "User-Agent": "Mozilla/5.0 (compatible; EvolipiaRadar/1.0)",
        Accept: "application/json",
      },
    });
    clearTimeout(timeout);

    if (!res.ok) return [];
    const data = await res.json();
    const hits = data.hits || [];

    const articles: CrawledArticle[] = [];
    for (const hit of hits) {
      if (!hit.title || !hit.url) continue;

      const title = hit.title.trim();
      const normUrl = normalizeUrl(hit.url);
      const domain = extractDomain(normUrl);

      // Filter for AI/Tech relevant items or high point stories
      const text = `${title} ${hit.story_text || ""}`.toLowerCase();
      const isAiRelevant =
        text.includes("ai") ||
        text.includes("llm") ||
        text.includes("gpt") ||
        text.includes("claude") ||
        text.includes("model") ||
        text.includes("agent") ||
        text.includes("robot") ||
        text.includes("open source") ||
        text.includes("database") ||
        text.includes("compiler") ||
        hit.points > 120;

      if (!isAiRelevant) continue;

      let publishedAt = new Date().toISOString();
      if (hit.created_at) {
        const d = new Date(hit.created_at);
        if (!isNaN(d.getTime())) {
          publishedAt = d.toISOString();
        }
      }

      const excerpt = cleanHtmlText(hit.story_text || "");
      const { category, tags } = inferCategoryAndTags(title, excerpt, domain);
      const score = Math.min(0.96, Math.max(0.80, Number((0.82 + Math.min(hit.points, 500) / 3000).toFixed(2))));

      articles.push({
        id: generateUuid(),
        title,
        url: normUrl,
        domain,
        published_at: publishedAt,
        category,
        score,
        tldr: excerpt.length > 20 ? (excerpt.length > 200 ? excerpt.slice(0, 197) + "..." : excerpt) : `${title} trending on Hacker News with ${hit.points || 50}+ points.`,
        why_it_matters: "Trending developer discussions and critical feedback from tech practitioners.",
        tags: Array.from(new Set([...tags, "hacker-news"])),
        impact: Math.round(score * 10),
        engineering_value: Math.round(score * 10),
        reasoning: `Front-page Hacker News discussion with ${hit.points} points.`,
      });

      if (articles.length >= 10) break;
    }

    return articles;
  } catch (e) {
    console.warn("[Crawler] Failed to fetch HN trending:", e);
    return [];
  }
}

// Load existing articles from disk
export function loadExistingNews(): { items: CrawledArticle[]; filePath: string } {
  const possiblePaths = [
    path.join(process.cwd(), "data", "news.json"),
    path.join(process.cwd(), "api", "news.json"),
  ];

  for (const p of possiblePaths) {
    try {
      if (fs.existsSync(p)) {
        const raw = fs.readFileSync(p, "utf-8");
        const parsed = JSON.parse(raw);
        if (Array.isArray(parsed.items) && parsed.items.length > 0) {
          return { items: parsed.items, filePath: p };
        }
      }
    } catch {}
  }

  return { items: [], filePath: possiblePaths[0] };
}

// Persist articles to both data/news.json and api/news.json
export function saveNewsFiles(items: CrawledArticle[]): void {
  const newsData: NewsDataFile = {
    items,
    last_updated: new Date().toISOString(),
    total_count: items.length,
  };

  const payload = JSON.stringify(newsData, null, 2);

  const targets = [
    path.join(process.cwd(), "data", "news.json"),
    path.join(process.cwd(), "api", "news.json"),
  ];

  targets.forEach((targetPath) => {
    try {
      const dir = path.dirname(targetPath);
      if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
      }
      fs.writeFileSync(targetPath, payload, "utf-8");
    } catch (err) {
      console.error(`[Crawler] Failed to write to ${targetPath}:`, err);
    }
  });
}

// Optional Neon DB / PostgreSQL upsert
async function upsertToDatabase(newItems: CrawledArticle[]): Promise<void> {
  const dbUrl = process.env.DATABASE_URL;
  if (!dbUrl || newItems.length === 0) return;

  let pool: Pool | null = null;
  try {
    pool = new Pool({
      connectionString: dbUrl,
      ssl: dbUrl.includes("sslmode=disable") ? false : { rejectUnauthorized: false },
      connectionTimeoutMillis: 5000,
    });

    for (const item of newItems) {
      try {
        // 1. Get or create source
        const srcRes = await pool.query(
          `INSERT INTO sources (name, type, category, url, enabled, status)
           VALUES ($1, 'rss', $2, $3, true, 'active')
           ON CONFLICT (url) DO UPDATE SET updated_at = now()
           RETURNING id`,
          [item.domain, item.category, `https://${item.domain}`]
        );
        const sourceId = srcRes.rows[0]?.id;

        // 2. Insert item
        const itemRes = await pool.query(
          `INSERT INTO items (source_id, title, url, domain, category, published_at, raw_excerpt, content_hash)
           VALUES ($1, $2, $3, $4, $5, $6, $7, MD5($3))
           ON CONFLICT (content_hash) DO UPDATE SET
             title = EXCLUDED.title,
             category = EXCLUDED.category
           RETURNING id`,
          [sourceId, item.title, item.url, item.domain, item.category, new Date(item.published_at), item.tldr]
        );
        const itemId = itemRes.rows[0]?.id;

        if (itemId) {
          // 3. Upsert score
          await pool.query(
            `INSERT INTO scores (item_id, relevance, impact, engineering_value, final)
             VALUES ($1, $2, $3, $4, $2)
             ON CONFLICT (item_id) DO UPDATE SET
               relevance = EXCLUDED.relevance,
               final = EXCLUDED.final,
               computed_at = now()`,
            [itemId, item.score, (item.impact || 85) / 100, (item.engineering_value || 85) / 100]
          );

          // 4. Upsert summary
          await pool.query(
            `INSERT INTO summaries (item_id, tldr, why_it_matters, tags, method)
             VALUES ($1, $2, $3, $4::jsonb, 'crawler')
             ON CONFLICT (item_id) DO UPDATE SET
               tldr = EXCLUDED.tldr,
               why_it_matters = EXCLUDED.why_it_matters,
               tags = EXCLUDED.tags`,
            [itemId, item.tldr, item.why_it_matters, JSON.stringify(item.tags)]
          );
        }
      } catch (rowErr) {
        // Continue on individual upsert errors
      }
    }
  } catch (dbErr) {
    console.warn("[Crawler] Database sync skipped or failed:", dbErr);
  } finally {
    if (pool) {
      await pool.end().catch(() => {});
    }
  }
}

// Master execution: crawls feeds, deduplicates, merges with existing items and saves
export async function executeLiveCrawlAndMerge(): Promise<{
  success: boolean;
  discovered: number;
  newMerged: number;
  totalCount: number;
  lastUpdated: string;
}> {
  console.log("[Crawler] Starting live crawl cycle across multi-sources...");

  // 1. Concurrently fetch candidate signals from active feeds
  const [hnTrending, arxivFeed, techCrunchFeed, hnRss] = await Promise.all([
    fetchHnTrending(),
    fetchRssFeed("https://arxiv.org/rss/cs.AI", "ArXiv AI", "arxiv.org"),
    fetchRssFeed("https://techcrunch.com/category/artificial-intelligence/feed/", "TechCrunch AI", "techcrunch.com"),
    fetchRssFeed("https://news.ycombinator.com/rss", "Hacker News", "news.ycombinator.com"),
  ]);

  const rawDiscovered: CrawledArticle[] = [
    ...hnTrending,
    ...arxivFeed,
    ...techCrunchFeed,
    ...hnRss,
  ];

  console.log(`[Crawler] Discovered ${rawDiscovered.length} raw articles from live feeds.`);

  // 2. Load existing articles from disk to merge with
  const { items: existingItems } = loadExistingNews();
  console.log(`[Crawler] Loaded ${existingItems.length} existing articles to preserve and merge.`);

  // 3. Build lookup sets for deduplication by URL and Title
  const seenUrls = new Set<string>();
  const seenTitles = new Set<string>();

  existingItems.forEach((it) => {
    if (it.url) seenUrls.add(normalizeUrl(it.url).toLowerCase());
    if (it.title) seenTitles.add(it.title.trim().toLowerCase());
  });

  // 4. Filter only genuinely new articles
  const genuinelyNew: CrawledArticle[] = [];
  for (const cand of rawDiscovered) {
    const norm = normalizeUrl(cand.url).toLowerCase();
    const titleKey = cand.title.trim().toLowerCase();

    if (seenUrls.has(norm) || seenTitles.has(titleKey)) {
      continue;
    }

    seenUrls.add(norm);
    seenTitles.add(titleKey);
    genuinelyNew.push(cand);
  }

  console.log(`[Crawler] Found ${genuinelyNew.length} new unique articles to merge.`);

  // 5. Merge new items at the top while preserving historical items (e.g. from 30 7 2026)
  const mergedList: CrawledArticle[] = [...genuinelyNew, ...existingItems];

  // Keep list capped at 150 items so file size and latency stay fast
  const trimmedList = mergedList.slice(0, 150);

  // 6. Save back to data/news.json and api/news.json
  saveNewsFiles(trimmedList);

  // 7. If database available, asynchronously sync new items
  if (genuinelyNew.length > 0) {
    upsertToDatabase(genuinelyNew).catch((e) => console.warn("[Crawler] Async DB upsert error:", e));
  }

  const nowIso = new Date().toISOString();

  return {
    success: true,
    discovered: rawDiscovered.length,
    newMerged: genuinelyNew.length,
    totalCount: trimmedList.length,
    lastUpdated: nowIso,
  };
}
