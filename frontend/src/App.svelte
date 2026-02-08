<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import Sidebar from "./lib/Sidebar.svelte";
  import AssetGrid from "./lib/AssetGrid.svelte";
  import BulkBar from "./lib/BulkBar.svelte";
  import DetailPanel from "./lib/DetailPanel.svelte";
  import Toast from "./lib/Toast.svelte";
  import {
    loading,
    assets,
    activeCollectionId,
    filterTag,
    filterTags,
    searchQuery,
    filteredAssets,
    tagsWithCounts,
    showBulkActions,
    selectedAssetIds,
    sortField,
    sortDirection,
  } from "./lib/stores";
  import {
    loadData,
    startBlenderPolling,
    addFolder,
    toggleTagFilter,
    clearTagFilters,
    viewLabel,
    clearSelection,
    setSort,
  } from "./lib/actions";
  import type { SortField } from "./lib/stores";

  function handleSortChange(e: Event) {
    const val = (e.currentTarget as HTMLSelectElement).value;
    setSort(val as SortField);
  }

  let blenderInterval: ReturnType<typeof setInterval>;

  onMount(() => {
    loadData();
    blenderInterval = startBlenderPolling();
  });

  onDestroy(() => {
    if (blenderInterval) clearInterval(blenderInterval);
  });

  // Esc key clears selection
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape" && $showBulkActions) {
      clearSelection();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="flex h-screen overflow-hidden text-left">
  <Sidebar />

  <main class="flex-1 overflow-y-auto p-6 relative">
    {#if $loading}
      <div class="flex flex-col items-center justify-center h-full opacity-70">
        <div class="text-5xl mb-4">⏳</div>
        <h2 class="text-base font-semibold m-0">Loading...</h2>
      </div>
    {:else if $assets.length === 0 && $activeCollectionId === null && !$filterTag && $filterTags.length === 0}
      <div class="flex flex-col items-center justify-center h-full opacity-70">
        <div class="text-5xl mb-4">📦</div>
        <h2 class="text-base font-semibold m-0 mb-1">No assets yet</h2>
        <p class="text-sm opacity-60 my-1">Add a watch folder to start indexing .glb files.</p>
        <button
          class="mt-4 px-4 py-2 rounded-md bg-accent-dim border border-accent-border text-white text-sm cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
          on:click={addFolder}
        >+ Add Watch Folder</button>
      </div>
    {:else}
      <div class="flex items-center gap-3 mb-4">
        <h2 class="m-0 text-sm whitespace-nowrap opacity-70 font-semibold">{viewLabel()}</h2>
        <input
          type="text"
          class="flex-1 max-w-[260px] bg-surface border border-surface-border rounded-md text-white text-xs px-3 py-1.5 outline-none font-inherit placeholder:text-white/25 focus:border-accent/40 transition-colors"
          placeholder="Search files…"
          bind:value={$searchQuery}
        />
        <div class="flex items-center gap-1">
          <select
            class="bg-surface border border-surface-border rounded text-white/70 text-[0.72rem] px-2 py-1 outline-none font-inherit cursor-pointer"
            value={$sortField}
            on:change={handleSortChange}
          >
            <option value="name">Name</option>
            <option value="date-added">Date Added</option>
            <option value="file-modified">Modified</option>
            <option value="file-size">Size</option>
          </select>
          <button
            class="bg-surface border border-surface-border rounded text-white/60 text-sm px-1.5 py-0.5 cursor-pointer font-inherit hover:bg-surface-hover transition-colors leading-none"
            on:click={() => setSort($sortField)}
            title="Toggle sort direction"
          >{$sortDirection === "asc" ? "↑" : "↓"}</button>
        </div>
        <span class="text-[0.7rem] opacity-35 whitespace-nowrap">
          {$filteredAssets.length} result{$filteredAssets.length !== 1 ? "s" : ""}
        </span>
      </div>

      {#if $tagsWithCounts.length > 0}
        <div class="flex flex-wrap gap-1.5 mb-3 items-center">
          {#each $tagsWithCounts as tag}
            <button
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-all
                {$filterTags.includes(tag.name)
                  ? 'bg-accent-dim border-accent-border text-blue-200'
                  : 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white/85'}"
              on:click={() => toggleTagFilter(tag.name)}
            >
              {tag.name}
              <span class="text-[0.6rem] opacity-40">{tag.count}</span>
            </button>
          {/each}
          {#if $filterTags.length > 0}
            <button
              class="bg-transparent border-none text-red-400/60 hover:text-red-400/90 cursor-pointer text-[0.65rem] font-inherit px-1.5 py-0.5"
              on:click={clearTagFilters}
            >✕ clear</button>
          {/if}
        </div>
      {/if}

      <AssetGrid />

      {#if $showBulkActions && $selectedAssetIds.size > 0}
        <BulkBar />
      {/if}
    {/if}
  </main>

  <DetailPanel />
  <Toast />
</div>
