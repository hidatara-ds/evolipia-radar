import { useState, useEffect, useCallback } from "react";
import { ItemQueryParams } from "../api/client";

export interface FilterPreset {
  id: string;
  name: string;
  filters: Partial<ItemQueryParams>;
}

const LOCAL_STORAGE_KEY = "evolipia_radar_saved_filters";

export function useFilters() {
  const [search, setSearch] = useState<string>("");
  const [debouncedSearch, setDebouncedSearch] = useState<string>("");
  const [dateRange, setDateRange] = useState<string>("7d");
  const [dateFrom, setDateFrom] = useState<string>("");
  const [dateTo, setDateTo] = useState<string>("");
  const [selectedSources, setSelectedSources] = useState<string[]>([]);
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [minRelevance, setMinRelevance] = useState<number>(30);
  const [status, setStatus] = useState<string>("all");
  const [sortBy, setSortBy] = useState<string>("date");
  const [sortOrder, setSortOrder] = useState<string>("desc");
  const [page, setPage] = useState<number>(1);
  const [savedPresets, setSavedPresets] = useState<FilterPreset[]>([]);

  // 300ms Debounce search input
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(search);
      setPage(1);
    }, 300);
    return () => clearTimeout(handler);
  }, [search]);

  // Load saved presets from localStorage
  useEffect(() => {
    if (typeof window !== "undefined") {
      try {
        const stored = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (stored) {
          setSavedPresets(JSON.parse(stored));
        }
      } catch (e) {
        console.error("Failed to load saved filter presets", e);
      }
    }
  }, []);

  // Sync state FROM URL query params on initial load
  useEffect(() => {
    if (typeof window !== "undefined") {
      const params = new URLSearchParams(window.location.search);
      if (params.has("search")) setSearch(params.get("search") || "");
      if (params.has("date_range")) setDateRange(params.get("date_range") || "7d");
      if (params.has("date_from")) setDateFrom(params.get("date_from") || "");
      if (params.has("date_to")) setDateTo(params.get("date_to") || "");
      if (params.has("min_relevance")) setMinRelevance(Number(params.get("min_relevance")) || 30);
      if (params.has("status")) setStatus(params.get("status") || "all");
      if (params.has("sort_by")) setSortBy(params.get("sort_by") || "date");
      if (params.has("sort_order")) setSortOrder(params.get("sort_order") || "desc");
      if (params.has("page")) setPage(Number(params.get("page")) || 1);

      const srcList = params.getAll("sources[]");
      if (srcList.length > 0) setSelectedSources(srcList);

      const catList = params.getAll("categories[]");
      if (catList.length > 0) setSelectedCategories(catList);
    }
  }, []);

  // Sync state TO URL query params when filters change
  const syncToURL = useCallback(() => {
    if (typeof window === "undefined") return;

    const params = new URLSearchParams();
    if (debouncedSearch) params.set("search", debouncedSearch);
    if (dateRange) params.set("date_range", dateRange);
    if (dateFrom) params.set("date_from", dateFrom);
    if (dateTo) params.set("date_to", dateTo);
    if (minRelevance > 0) params.set("min_relevance", minRelevance.toString());
    if (status !== "all") params.set("status", status);
    if (sortBy) params.set("sort_by", sortBy);
    if (sortOrder) params.set("sort_order", sortOrder);
    if (page > 1) params.set("page", page.toString());

    selectedSources.forEach((s: string) => params.append("sources[]", s));
    selectedCategories.forEach((c: string) => params.append("categories[]", c));

    const newUrl = `${window.location.pathname}?${params.toString()}`;
    window.history.replaceState(null, "", newUrl);
  }, [debouncedSearch, dateRange, dateFrom, dateTo, minRelevance, status, sortBy, sortOrder, page, selectedSources, selectedCategories]);

  useEffect(() => {
    syncToURL();
  }, [syncToURL]);

  const savePreset = (name: string) => {
    const newPreset: FilterPreset = {
      id: Date.now().toString(),
      name,
      filters: {
        search: debouncedSearch,
        date_from: dateFrom,
        date_to: dateTo,
        sources: selectedSources,
        categories: selectedCategories,
        min_relevance: minRelevance,
        status,
        sort_by: sortBy,
        sort_order: sortOrder,
      },
    };
    const updated = [...savedPresets, newPreset];
    setSavedPresets(updated);
    if (typeof window !== "undefined") {
      localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(updated));
    }
  };

  const loadPreset = (preset: FilterPreset) => {
    const f = preset.filters;
    if (f.search !== undefined) setSearch(f.search);
    if (f.date_from !== undefined) setDateFrom(f.date_from);
    if (f.date_to !== undefined) setDateTo(f.date_to);
    if (f.sources !== undefined) setSelectedSources(f.sources);
    if (f.categories !== undefined) setSelectedCategories(f.categories);
    if (f.min_relevance !== undefined) setMinRelevance(f.min_relevance);
    if (f.status !== undefined) setStatus(f.status);
    if (f.sort_by !== undefined) setSortBy(f.sort_by);
    if (f.sort_order !== undefined) setSortOrder(f.sort_order);
    setPage(1);
  };

  const deletePreset = (id: string) => {
    const updated = savedPresets.filter((p: FilterPreset) => p.id !== id);
    setSavedPresets(updated);
    if (typeof window !== "undefined") {
      localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(updated));
    }
  };

  const resetFilters = () => {
    setSearch("");
    setDebouncedSearch("");
    setDateRange("7d");
    setDateFrom("");
    setDateTo("");
    setSelectedSources([]);
    setSelectedCategories([]);
    setMinRelevance(30);
    setStatus("all");
    setSortBy("date");
    setSortOrder("desc");
    setPage(1);
  };

  const queryParams: ItemQueryParams = {
    search: debouncedSearch,
    date_from: dateFrom,
    date_to: dateTo,
    sources: selectedSources,
    categories: selectedCategories,
    min_relevance: minRelevance,
    status,
    sort_by: sortBy,
    sort_order: sortOrder,
    page,
    limit: 20,
  };

  return {
    search,
    setSearch,
    debouncedSearch,
    dateRange,
    setDateRange,
    dateFrom,
    setDateFrom,
    dateTo,
    setDateTo,
    selectedSources,
    setSelectedSources,
    selectedCategories,
    setSelectedCategories,
    minRelevance,
    setMinRelevance,
    status,
    setStatus,
    sortBy,
    setSortBy,
    sortOrder,
    setSortOrder,
    page,
    setPage,
    savedPresets,
    savePreset,
    loadPreset,
    deletePreset,
    resetFilters,
    queryParams,
  };
}
