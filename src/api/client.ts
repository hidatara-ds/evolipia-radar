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
  
  if (!res.ok) {
    // Fallback to legacy news endpoint if /api/items is not present
    const legacyRes = await fetch(`${API_BASE_URL}/api/news?${query.toString()}`);
    if (!legacyRes.ok) {
      throw new Error(`API error: ${res.statusText}`);
    }
    const legacyJson = await legacyRes.json();
    return {
      success: legacyJson.success,
      data: legacyJson.data?.items || [],
      total_count: legacyJson.data?.total_count || 0,
      filtered_count: legacyJson.data?.items?.length || 0,
      page: 1,
      total_pages: 1,
      last_updated: legacyJson.data?.last_updated || new Date().toISOString(),
    };
  }

  return res.json();
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
