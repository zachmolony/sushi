import { get } from "svelte/store";
import {
  OpenFileInFolder,
  SendToBlender,
  PingBlender,
  AddWatchFolder,
  GetWatchFolders,
  RemoveWatchFolder,
  GetAssets,
  GetAssetsByTag,
  GetAssetsByTags,
  GetAllTags,
  GetTagsWithCounts,
  GetTagsForAsset,
  AddTagToAsset,
  RemoveTagFromAsset,
  BulkTagAssets,
  BulkAddToCollection,
  SaveThumbnail,
  SavePolyCount,
  CreateCollection,
  GetCollections,
  DeleteCollection,
  AddToCollection,
  RemoveFromCollection,
  GetCollectionAssets,
  GetCollectionsForAsset,
  GetFileServerURL,
  ToggleFavorite,
  MarkAssetUsed,
  GetUntaggedAssets,
  GetFavoritedAssets,
  GetRecentlyAddedAssets,
  GetRecentlyUsedAssets,
  BulkSetFavorite,
} from "../../wailsjs/go/main/App.js";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime.js";
import { renderThumbnail } from "./thumbnails";
import type { Asset } from "./types";
import {
  assets,
  watchFolders,
  allTags,
  tagsWithCounts,
  collections,
  loading,
  filterTag,
  filterTags,
  activeCollectionId,
  activeView,
  sortField,
  sortDirection,
  selectedAsset,
  selectedAssetTags,
  selectedAssetCollections,
  selectedAssetIds,
  lastClickedIndex,
  showBulkActions,
  thumbnailCache,
  fileServerBase,
  blenderConnected,
  toastMessage,
  toastVisible,
  displayedAssets,
  filteredAssets,
} from "./stores";
import type { ViewId, SortField, SortDirection } from "./stores";

// --- Toast ---

export function showToast(msg: string) {
  toastMessage.set(msg);
  toastVisible.set(true);
  setTimeout(() => {
    toastVisible.set(false);
  }, 2000);
}

// --- Thumbnails ---

export async function generateMissingThumbnails() {
  let base = get(fileServerBase);
  if (!base) {
    try {
      base = await GetFileServerURL();
      fileServerBase.set(base);
      console.log("[sushi] File server URL:", base);
    } catch (e) {
      console.error("[sushi] Failed to get file server URL:", e);
      return;
    }
  }
  if (!base) {
    console.warn("[sushi] No file server available, skipping thumbnails");
    return;
  }

  const currentAssets = get(assets);
  let generated = 0;
  let cached = 0;
  let failed = 0;
  const cache = { ...get(thumbnailCache) };

  for (const asset of currentAssets) {
    if (asset.thumbnail && asset.poly_count > 0) {
      cache[asset.id] = asset.thumbnail;
      cached++;
      continue;
    }

    // Has thumbnail but no poly count — need to re-parse for count
    const needsPolyOnly = !!asset.thumbnail && asset.poly_count === 0;
    if (needsPolyOnly) {
      cache[asset.id] = asset.thumbnail;
      cached++;
    }

    try {
      const url = `${base}/localfile/?path=${encodeURIComponent(asset.absolute_path)}`;
      const result = await renderThumbnail(url);
      if (result) {
        if (!needsPolyOnly) {
          cache[asset.id] = result.dataUrl;
          SaveThumbnail(asset.id, result.dataUrl).catch(() => {});
          generated++;
        }
        if (result.polyCount > 0) {
          SavePolyCount(asset.id, result.polyCount).catch(() => {});
        }
      } else if (!needsPolyOnly) {
        failed++;
      }
    } catch (e) {
      console.warn(`[sushi] Thumbnail failed for ${asset.filename}:`, e);
      failed++;
    }
  }
  thumbnailCache.set(cache);
  console.log(
    `[sushi] Thumbnails: ${cached} cached, ${generated} generated, ${failed} failed`
  );
}

// --- Init ---

export async function loadData() {
  loading.set(true);
  try {
    const [a, f, t, c, tc] = await Promise.all([
      GetAssets(),
      GetWatchFolders(),
      GetAllTags(),
      GetCollections(),
      GetTagsWithCounts(),
    ]);
    assets.set(a || []);
    displayedAssets.set(a || []);
    watchFolders.set(f || []);
    allTags.set(t || []);
    collections.set(c || []);
    tagsWithCounts.set(tc || []);
  } catch (e) {
    console.error("Failed to load data:", e);
  }
  loading.set(false);
  generateMissingThumbnails();
}

export function startBlenderPolling() {
  function check() {
    PingBlender()
      .then((status) => blenderConnected.set(status.connected))
      .catch(() => blenderConnected.set(false));
  }
  check();
  return setInterval(check, 5000);
}

// --- Filtering ---

export async function applyFilter() {
  loading.set(true);
  const colId = get(activeCollectionId);
  const fTags = get(filterTags);
  const fTag = get(filterTag);
  const view = get(activeView);
  try {
    if (colId !== null) {
      const a = await GetCollectionAssets(colId);
      displayedAssets.set(a || []);
    } else if (fTags.length > 0) {
      const a = await GetAssetsByTags(fTags);
      displayedAssets.set(a || []);
    } else if (fTag) {
      const a = await GetAssetsByTag(fTag);
      displayedAssets.set(a || []);
    } else {
      // Apply smart view
      let a: Asset[];
      switch (view) {
        case "untagged":
          a = await GetUntaggedAssets();
          break;
        case "favorites":
          a = await GetFavoritedAssets();
          break;
        case "recent-added":
          a = await GetRecentlyAddedAssets();
          break;
        case "recent-used":
          a = await GetRecentlyUsedAssets();
          break;
        default:
          a = await GetAssets();
          break;
      }
      displayedAssets.set(a || []);
      if (view === "all") assets.set(a || []);
    }
  } catch (e) {
    console.error("Filter failed:", e);
  }
  loading.set(false);
  generateMissingThumbnails();
}

export function toggleTagFilter(tag: string) {
  activeCollectionId.set(null);
  activeView.set("all");
  filterTag.set("");
  const current = get(filterTags);
  const idx = current.indexOf(tag);
  if (idx >= 0) {
    filterTags.set(current.filter((t) => t !== tag));
  } else {
    filterTags.set([...current, tag]);
  }
  applyFilter();
}

export function clearTagFilters() {
  filterTags.set([]);
  filterTag.set("");
  activeCollectionId.set(null);
  activeView.set("all");
  applyFilter();
}

export function setTagFilter(tag: string) {
  activeCollectionId.set(null);
  activeView.set("all");
  filterTags.set([]);
  filterTag.update((current) => (current === tag ? "" : tag));
  applyFilter();
}

export function setCollectionFilter(id: number | null) {
  filterTag.set("");
  filterTags.set([]);
  activeView.set("all");
  activeCollectionId.update((current) => (current === id ? null : id));
  applyFilter();
}

// --- Views ---

export function setActiveView(view: ViewId) {
  filterTag.set("");
  filterTags.set([]);
  activeCollectionId.set(null);
  activeView.set(view);
  applyFilter();
}

// --- Sorting ---

export function setSort(field: SortField) {
  const currentField = get(sortField);
  if (currentField === field) {
    // Toggle direction
    sortDirection.update((d) => (d === "asc" ? "desc" : "asc"));
  } else {
    sortField.set(field);
    sortDirection.set(field === "name" ? "asc" : "desc");
  }
}

// --- Multi-select ---

export function handleAssetClick(asset: Asset, index: number, e: MouseEvent) {
  const ids = get(selectedAssetIds);
  const fa = get(filteredAssets);

  if (e.ctrlKey || e.metaKey) {
    // Ctrl/Cmd+click: toggle this item in selection
    const next = new Set(ids);
    if (next.has(asset.id)) {
      next.delete(asset.id);
    } else {
      next.add(asset.id);
    }
    selectedAssetIds.set(next);
    lastClickedIndex.set(index);
    showBulkActions.set(next.size > 0);
  } else if (e.shiftKey && get(lastClickedIndex) >= 0) {
    // Shift+click: range select
    const start = Math.min(get(lastClickedIndex), index);
    const end = Math.max(get(lastClickedIndex), index);
    const next = new Set(ids);
    for (let i = start; i <= end; i++) {
      if (fa[i]) next.add(fa[i].id);
    }
    selectedAssetIds.set(next);
    showBulkActions.set(next.size > 0);
  } else if (ids.size > 0) {
    // Already in bulk-select mode: toggle this item (don't nuke the selection)
    const next = new Set(ids);
    if (next.has(asset.id)) {
      next.delete(asset.id);
    } else {
      next.add(asset.id);
    }
    selectedAssetIds.set(next);
    lastClickedIndex.set(index);
    showBulkActions.set(next.size > 0);
    // Also show detail panel for the clicked asset
    selectAsset(asset);
  } else {
    // No selection active: just open detail panel
    lastClickedIndex.set(index);
    selectAsset(asset);
  }
}

export function selectAllVisible() {
  const fa = get(filteredAssets);
  const next = new Set(get(selectedAssetIds));
  for (const a of fa) next.add(a.id);
  selectedAssetIds.set(next);
  showBulkActions.set(true);
}

export function clearSelection() {
  selectedAssetIds.set(new Set());
  showBulkActions.set(false);
}

export function toggleAssetSelection(assetId: number) {
  const ids = get(selectedAssetIds);
  const next = new Set(ids);
  if (next.has(assetId)) {
    next.delete(assetId);
  } else {
    next.add(assetId);
  }
  selectedAssetIds.set(next);
  showBulkActions.set(next.size > 0);
}

// --- Bulk actions ---

export async function bulkTag(tagInput: string) {
  const tagName = tagInput.trim().toLowerCase();
  const ids = get(selectedAssetIds);
  if (!tagName || ids.size === 0) return;
  try {
    await BulkTagAssets([...ids], tagName);
    showToast(`Tagged ${ids.size} assets with "${tagName}"`);
    const [t, tc] = await Promise.all([GetAllTags(), GetTagsWithCounts()]);
    allTags.set(t || []);
    tagsWithCounts.set(tc || []);
    const sa = get(selectedAsset);
    if (sa && ids.has(sa.id)) {
      const tags = await GetTagsForAsset(sa.id);
      selectedAssetTags.set(tags || []);
    }
  } catch (e) {
    showToast("Failed to bulk tag");
  }
}

export async function bulkAddToCollection(collectionId: number) {
  const ids = get(selectedAssetIds);
  if (ids.size === 0) return;
  try {
    await BulkAddToCollection(collectionId, [...ids]);
    const c = await GetCollections();
    collections.set(c || []);
    showToast(`Added ${ids.size} assets to collection`);
    const sa = get(selectedAsset);
    if (sa && ids.has(sa.id)) {
      const cols = await GetCollectionsForAsset(sa.id);
      selectedAssetCollections.set(cols || []);
    }
  } catch (e) {
    showToast("Failed to add to collection");
  }
}

export async function bulkSendToBlender() {
  const ids = get(selectedAssetIds);
  if (ids.size === 0) return;
  const fa = get(filteredAssets);
  const selected = fa.filter((a) => ids.has(a.id));
  const paths = selected.map((a) => a.absolute_path);
  try {
    const status = await SendToBlender(paths);
    if (status.connected && !status.error) {
      showToast(`Sent ${paths.length} assets → Blender`);
      // Mark all sent assets as used
      for (const a of selected) {
        MarkAssetUsed(a.id).catch(() => {});
      }
    } else {
      showToast(status.error || "Could not reach Blender");
    }
    blenderConnected.set(status.connected);
  } catch (e) {
    showToast("Failed to send to Blender");
  }
}

export async function bulkSetFavorite(favorited: boolean) {
  const ids = get(selectedAssetIds);
  if (ids.size === 0) return;
  try {
    await BulkSetFavorite([...ids], favorited);
    showToast(favorited ? `Favorited ${ids.size} assets` : `Unfavorited ${ids.size} assets`);
    await applyFilter();
  } catch (e) {
    showToast("Failed to update favorites");
  }
}

// --- Hover quick-add ---

export async function hoverAddToCollection(
  assetId: number,
  collectionId: number,
  e: MouseEvent
) {
  e.stopPropagation();
  try {
    await AddToCollection(collectionId, assetId);
    const c = await GetCollections();
    collections.set(c || []);
    showToast("Added to collection");
  } catch (e) {
    showToast("Already in collection");
  }
}

// --- Watch Folders ---

export async function addFolder() {
  try {
    const updatedAssets = await AddWatchFolder();
    assets.set(updatedAssets || []);
    displayedAssets.set(updatedAssets || []);
    const folders = await GetWatchFolders();
    watchFolders.set(folders || []);
    const a = get(assets);
    if (a.length > 0) showToast(`Found ${a.length} assets`);
    generateMissingThumbnails();
  } catch (e) {
    showToast("Failed to add folder");
  }
}

export async function removeFolder(id: number) {
  try {
    await RemoveWatchFolder(id);
    watchFolders.update((f) => f.filter((folder) => folder.id !== id));
    const a = await GetAssets();
    assets.set(a || []);
    displayedAssets.set(a || []);
    const sa = get(selectedAsset);
    if (sa && !(a || []).find((asset) => asset.id === sa.id)) {
      selectedAsset.set(null);
      selectedAssetTags.set([]);
    }
    showToast("Folder removed");
  } catch (e) {
    showToast("Failed to remove folder");
  }
}

// --- Asset detail ---

export async function selectAsset(asset: Asset) {
  selectedAsset.set(asset);
  try {
    const [tags, cols] = await Promise.all([
      GetTagsForAsset(asset.id),
      GetCollectionsForAsset(asset.id),
    ]);
    selectedAssetTags.set(tags || []);
    selectedAssetCollections.set(cols || []);
  } catch (e) {
    selectedAssetTags.set([]);
    selectedAssetCollections.set([]);
  }
}

export async function copyPath() {
  const sa = get(selectedAsset);
  if (!sa) return;
  await ClipboardSetText(sa.absolute_path);
  showToast("Path copied!");
}

export async function sendToBlender() {
  const sa = get(selectedAsset);
  if (!sa) return;
  try {
    const status = await SendToBlender([sa.absolute_path]);
    if (status.connected && !status.error) {
      showToast(`Sent → ${sa.filename}`);
      // Mark as used
      MarkAssetUsed(sa.id).catch(() => {});
    } else {
      showToast(status.error || "Could not reach Blender");
    }
    blenderConnected.set(status.connected);
  } catch (e) {
    showToast("Failed to send to Blender");
  }
}

export async function showInFolder() {
  const sa = get(selectedAsset);
  if (!sa) return;
  await OpenFileInFolder(sa.absolute_path);
}

// --- Tagging (single asset) ---

export async function addTag(tagName: string) {
  const sa = get(selectedAsset);
  if (!sa || !tagName.trim()) return;
  try {
    const tags = await AddTagToAsset(sa.id, tagName.trim().toLowerCase());
    selectedAssetTags.set(tags || []);
    const [t, tc] = await Promise.all([GetAllTags(), GetTagsWithCounts()]);
    allTags.set(t || []);
    tagsWithCounts.set(tc || []);
  } catch (e) {
    showToast("Failed to add tag");
  }
}

export async function removeTag(tagId: number) {
  const sa = get(selectedAsset);
  if (!sa) return;
  try {
    const tags = await RemoveTagFromAsset(sa.id, tagId);
    selectedAssetTags.set(tags || []);
    const [t, tc] = await Promise.all([GetAllTags(), GetTagsWithCounts()]);
    allTags.set(t || []);
    tagsWithCounts.set(tc || []);
  } catch (e) {
    showToast("Failed to remove tag");
  }
}

// --- Favorites ---

export async function toggleFavorite(assetId: number) {
  try {
    const isFav = await ToggleFavorite(assetId);
    showToast(isFav ? "Added to favorites" : "Removed from favorites");
    // Refresh the asset list
    const a = await GetAssets();
    assets.set(a || []);
    await applyFilter();
  } catch (e) {
    showToast("Failed to toggle favorite");
  }
}

// --- Collections ---

export const SHELF_ICONS = [
  "📁", "📺", "🎮", "⭐", "🗑️", "🔧",
  "🎨", "🏠", "🚗", "🌿", "💀", "🍣",
];

export const SUGGESTED_TAGS = [
  "lowpoly", "psx", "vehicle", "clutter", "environment",
  "character", "animated", "prop", "weapon", "building",
  "nature", "furniture", "sci-fi", "fantasy", "modular",
];

export async function createCollection(name: string, icon: string) {
  if (!name.trim()) return;
  try {
    await CreateCollection(name.trim(), icon);
    const c = await GetCollections();
    collections.set(c || []);
    showToast("Collection created");
  } catch (e) {
    showToast("Failed to create collection");
  }
}

export async function deleteCollectionById(id: number) {
  try {
    await DeleteCollection(id);
    if (get(activeCollectionId) === id) {
      activeCollectionId.set(null);
      applyFilter();
    }
    const c = await GetCollections();
    collections.set(c || []);
    showToast("Collection deleted");
  } catch (e) {
    showToast("Failed to delete collection");
  }
}

export async function addSelectedToCollection(collectionId: number) {
  const sa = get(selectedAsset);
  if (!sa) return;
  try {
    await AddToCollection(collectionId, sa.id);
    const [cols, allCols] = await Promise.all([
      GetCollectionsForAsset(sa.id),
      GetCollections(),
    ]);
    selectedAssetCollections.set(cols || []);
    collections.set(allCols || []);
    showToast("Added to collection");
  } catch (e) {
    showToast("Failed to add to collection");
  }
}

export async function removeSelectedFromCollection(collectionId: number) {
  const sa = get(selectedAsset);
  if (!sa) return;
  try {
    await RemoveFromCollection(collectionId, sa.id);
    const [cols, allCols] = await Promise.all([
      GetCollectionsForAsset(sa.id),
      GetCollections(),
    ]);
    selectedAssetCollections.set(cols || []);
    collections.set(allCols || []);
    if (get(activeCollectionId) === collectionId) {
      applyFilter();
    }
    showToast("Removed from collection");
  } catch (e) {
    showToast("Failed to remove from collection");
  }
}

// --- Utilities ---

export function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function viewLabel(): string {
  const colId = get(activeCollectionId);
  const cols = get(collections);
  const fTags = get(filterTags);
  const fTag = get(filterTag);
  const view = get(activeView);

  if (colId !== null) {
    const col = cols.find((c) => c.id === colId);
    return col ? `${col.icon} ${col.name}` : "Collection";
  }
  if (fTags.length > 0) return fTags.map((t) => `#${t}`).join(" + ");
  if (fTag) return `#${fTag}`;

  switch (view) {
    case "untagged": return "🏷️ Untagged";
    case "favorites": return "⭐ Favorites";
    case "recent-added": return "🆕 Recently Added";
    case "recent-used": return "🕐 Recently Used";
    default: return "All assets";
  }
}
