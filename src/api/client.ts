// API Client for Evolipia Radar

export interface ItemQueryParams {
  search?: string;
  date_from?: string;
  date_to?: string;
  sources?: string[];
  categories?: string[];
  min_relevance?: number;
  status?: string;
  sort_by?: string;
  sort_order?: string;
  page?: number;
  limit?: number;
}

export interface NewsItem {
  id: string;
  source_id: string;
  source_name?: string;
  title: string;
  url: string;
  published_at: string;
  domain: string;
  category: string;
  raw_excerpt?: string;
  crawl_status?: string;
  crawl_error?: string;
  relevance_score?: number;
  validated_at?: string;
  created_at: string;
  scaled_score?: number;
  tldr?: string;
  summary?: string;
  why_it_matters?: string;
  tags?: string[];
  impact?: number;
  engineering_value?: number;
  reasoning?: string;
}

export interface PaginatedItemsResponse {
  success: boolean;
  data: NewsItem[];
  total_count: number;
  filtered_count: number;
  page: number;
  total_pages: number;
  last_updated: string;
  error?: string;
}

const API_BASE_URL = typeof window !== "undefined" 
  ? (process.env.NEXT_PUBLIC_API_URL || "") 
  : "";

export async function fetchItems(params: ItemQueryParams): Promise<PaginatedItemsResponse> {
  const query = new URLSearchParams();

  if (params.search) query.set("search", params.search);
  if (params.date_from) query.set("date_from", params.date_from);
  if (params.date_to) query.set("date_to", params.date_to);
  if (params.min_relevance !== undefined) query.set("min_relevance", params.min_relevance.toString());
  if (params.status) query.set("status", params.status);
  if (params.sort_by) query.set("sort_by", params.sort_by);
  if (params.sort_order) query.set("sort_order", params.sort_order);
  if (params.page) query.set("page", params.page.toString());
  if (params.limit) query.set("limit", params.limit.toString());

  if (params.sources && params.sources.length > 0) {
    params.sources.forEach(s => query.append("sources[]", s));
  }

  if (params.categories && params.categories.length > 0) {
    params.categories.forEach(c => query.append("categories[]", c));
  }

  const endpoint = `${API_BASE_URL}/api/items?${query.toString()}`;
  const res = await fetch(endpoint);
  
  let jsonResult: PaginatedItemsResponse;

  if (!res.ok) {
    // Fallback to legacy news endpoint if /api/items is not present
    const legacyRes = await fetch(`${API_BASE_URL}/api/news?${query.toString()}`);
    if (!legacyRes.ok) {
      throw new Error(`API error: ${res.statusText}`);
    }
    const legacyJson = await legacyRes.json();
    const rawList: any[] = legacyJson.data?.items || (Array.isArray(legacyJson.data) ? legacyJson.data : []);
    jsonResult = {
      success: legacyJson.success,
      data: rawList,
      total_count: legacyJson.data?.total_count || rawList.length,
      filtered_count: legacyJson.data?.filtered_count || rawList.length,
      page: 1,
      total_pages: 1,
      last_updated: legacyJson.data?.last_updated || new Date().toISOString(),
    };
  } else {
    jsonResult = await res.json();
  }

  // Normalize data items
  if (jsonResult && Array.isArray(jsonResult.data)) {
    jsonResult.data = jsonResult.data.map((item: any) => {
      const src = (item.source_name || "").toLowerCase().trim();
      const dom = (item.domain || "").toLowerCase().trim();
      let effSource = item.source_name;

      if (dom.includes("ycombinator.com") || dom.includes("news.ycombinator")) {
        effSource = "Hacker News";
      } else if (dom.includes("arxiv")) {
        effSource = "ArXiv AI";
      } else if (dom.includes("techcrunch")) {
        effSource = "TechCrunch AI";
      } else if (dom.includes("reddit")) {
        effSource = "Reddit MachineLearning";
      } else if (dom.includes("twitter") || dom.includes("x.com")) {
        effSource = "Twitter / X";
      } else if (dom.includes("github")) {
        effSource = "GitHub Trending";
      } else if (dom.includes("openai.com")) {
        effSource = "OpenAI Blog";
      } else if (dom.includes("deepmind")) {
        effSource = "DeepMind Research";
      } else if (dom.includes("stability.ai")) {
        effSource = "Stability AI";
      } else if (dom.includes("anthropic.com")) {
        effSource = "Anthropic";
      } else if (!effSource || effSource === "Global Source" || effSource === "Unknown" || effSource === "RSSAgent") {
        if (src.includes("hacker") || src.includes("ycombinator")) {
          effSource = "Hacker News";
        } else if (src.includes("arxiv")) {
          effSource = "ArXiv AI";
        } else if (src.includes("techcrunch")) {
          effSource = "TechCrunch AI";
        } else if (src.includes("reddit")) {
          effSource = "Reddit MachineLearning";
        } else if (src.includes("twitter") || src.includes("x.com")) {
          effSource = "Twitter / X";
        } else if (src.includes("github")) {
          effSource = "GitHub Trending";
        } else {
          effSource = item.domain || "Global Source";
        }
      }

      return {
        ...item,
        raw_excerpt: item.raw_excerpt || item.tldr || item.title || "",
        relevance_score: item.relevance_score ?? (item.score ? Math.round(item.score * 100) : (item.raw_score ? Math.round(item.raw_score * 100) : 85)),
        scaled_score: item.scaled_score ?? (item.score ? Number((item.score * 10).toFixed(1)) : 8.5),
        source_name: effSource,
        crawl_status: item.crawl_status || "done",
      };
    });
  }

  return jsonResult;
}

export async function triggerManualCrawl(): Promise<{ success: boolean; message: string }> {
  const res = await fetch(`${API_BASE_URL}/api/crawl`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });

  if (!res.ok) {
    throw new Error(`Failed to trigger crawl: ${res.statusText}`);
  }

  return res.json();
}

export async function fetchCrawlStatus(): Promise<{
  last_run_time: string;
  last_run_status: string;
  last_items_count: number;
  last_error: string;
  is_running: boolean;
}> {
  const res = await fetch(`${API_BASE_URL}/api/crawl/status`);
  if (!res.ok) {
    throw new Error("Failed to fetch crawl status");
  }
  return res.json();
}
