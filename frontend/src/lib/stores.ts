import { writable, derived } from "svelte/store";
import type { Asset, WatchFolder, Tag, TagWithCount, Collection } from "./types";

// --- Core data ---
export const assets = writable<Asset[]>([]);
export const watchFolders = writable<WatchFolder[]>([]);
export const allTags = writable<Tag[]>([]);
export const tagsWithCounts = writable<TagWithCount[]>([]);
export const collections = writable<Collection[]>([]);

// --- UI state ---
export const loading = writable(true);
export const searchQuery = writable("");
export const filterTag = writable("");
export const filterTags = writable<string[]>([]);
export const activeCollectionId = writable<number | null>(null);

// --- Views ---
export type ViewId = "all" | "untagged" | "recent-added" | "recent-used" | "favorites";
export const activeView = writable<ViewId>("all");

// --- Sorting ---
export type SortField = "name" | "date-added" | "file-modified" | "file-size";
export type SortDirection = "asc" | "desc";
export const sortField = writable<SortField>("name");
export const sortDirection = writable<SortDirection>("asc");

// --- Selection ---
export const selectedAsset = writable<Asset | null>(null);
export const selectedAssetTags = writable<Tag[]>([]);
export const selectedAssetCollections = writable<Collection[]>([]);
export const selectedAssetIds = writable<Set<number>>(new Set());
export const lastClickedIndex = writable<number>(-1);
export const showBulkActions = writable(false);

// --- Thumbnails ---
export const thumbnailCache = writable<Record<number, string>>({});
export const fileServerBase = writable("");

// --- Blender ---
export const blenderConnected = writable(false);

// --- Toast ---
export const toastMessage = writable("");
export const toastVisible = writable(false);

// --- Hover ---
export const hoverAssetId = writable<number | null>(null);

// --- Display / filtering ---
export const displayedAssets = writable<Asset[]>([]);

// Derived: filteredAssets applies search + sorting on top of displayedAssets
export const filteredAssets = derived(
  [displayedAssets, searchQuery, sortField, sortDirection],
  ([$displayedAssets, $searchQuery, $sortField, $sortDirection]) => {
    let result = $displayedAssets;

    // Search filter
    if ($searchQuery.trim()) {
      const q = $searchQuery.trim().toLowerCase();
      result = result.filter((a) => a.filename.toLowerCase().includes(q));
    }

    // Sort
    const dir = $sortDirection === "asc" ? 1 : -1;
    result = [...result].sort((a, b) => {
      switch ($sortField) {
        case "name":
          return dir * a.filename.localeCompare(b.filename);
        case "date-added":
          return dir * (new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
        case "file-modified":
          return dir * (new Date(a.modified_at).getTime() - new Date(b.modified_at).getTime());
        case "file-size":
          return dir * (a.file_size - b.file_size);
        default:
          return 0;
      }
    });

    return result;
  }
);
