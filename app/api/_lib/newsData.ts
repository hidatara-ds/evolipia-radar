import { readFile } from "node:fs/promises";
import path from "node:path";

export type NewsItem = {
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
};

export type NewsData = {
  items: NewsItem[];
  total_count?: number;
  last_updated?: string;
};

export async function loadNewsData(): Promise<NewsData> {
  const candidates = [
    path.join(process.cwd(), "data", "news.json"),
    path.join(process.cwd(), "api", "news.json"),
  ];

  let lastErr: unknown;
  for (const p of candidates) {
    try {
      const raw = await readFile(p, "utf8");
      return JSON.parse(raw) as NewsData;
    } catch (e) {
      lastErr = e;
    }
  }

  throw lastErr instanceof Error ? lastErr : new Error("Failed to load news.json");
}

export function normalizeTopic(topic: string | null): string | null {
  if (!topic) return null;
  const t = topic.trim();
  if (!t || t.toLowerCase() === "all") return null;
  return t;
}

