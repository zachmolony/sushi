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
  } from "./lib/stores";
  import {
    loadData,
    startBlenderPolling,
    addFolder,
    toggleTagFilter,
    clearTagFilters,
    viewLabel,
    clearSelection,
  } from "./lib/actions";

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

<div class="app">
  <Sidebar />

  <main class="main">
    {#if $loading}
      <div class="empty-state">
        <div class="empty-icon">⏳</div>
        <h2>Loading...</h2>
      </div>
    {:else if $assets.length === 0 && $activeCollectionId === null && !$filterTag && $filterTags.length === 0}
      <div class="empty-state">
        <div class="empty-icon">📦</div>
        <h2>No assets yet</h2>
        <p>Add a watch folder to start indexing .glb files.</p>
        <button class="btn btn-primary" on:click={addFolder} style="margin-top: 1rem;">
          + Add Watch Folder
        </button>
      </div>
    {:else}
      <div class="main-header">
        <h2 class="view-label">{viewLabel()}</h2>
        <input
          type="text"
          class="search-input"
          placeholder="Search files…"
          bind:value={$searchQuery}
        />
        <span class="result-count">
          {$filteredAssets.length} result{$filteredAssets.length !== 1 ? "s" : ""}
        </span>
      </div>

      {#if $tagsWithCounts.length > 0}
        <div class="tag-bar">
          {#each $tagsWithCounts as tag}
            <button
              class="tag-bar-chip"
              class:active={$filterTags.includes(tag.name)}
              on:click={() => toggleTagFilter(tag.name)}
            >
              {tag.name}
              <span class="tag-bar-count">{tag.count}</span>
            </button>
          {/each}
          {#if $filterTags.length > 0}
            <button class="tag-bar-clear" on:click={clearTagFilters}>✕ clear</button>
          {/if}
        </div>
      {/if}

      {#if $showBulkActions && $selectedAssetIds.size > 0}
        <BulkBar />
      {/if}

      <AssetGrid />
    {/if}
  </main>

  <DetailPanel />
  <Toast />
</div>

<style>
  .app {
    display: flex;
    height: 100vh;
    overflow: hidden;
    text-align: left;
  }

  .main {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    opacity: 0.7;
  }
  .empty-icon { font-size: 3rem; margin-bottom: 1rem; }
  .empty-state h2 { margin: 0 0 0.5rem; }
  .empty-state p { margin: 0.25rem 0; }

  .main-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }
  .view-label {
    margin: 0;
    font-size: 1rem;
    white-space: nowrap;
    opacity: 0.7;
  }
  .search-input {
    flex: 1;
    max-width: 260px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 6px;
    color: white;
    font-size: 0.8rem;
    padding: 0.35rem 0.6rem;
    outline: none;
    font-family: inherit;
  }
  .search-input::placeholder { color: rgba(255, 255, 255, 0.25); }
  .search-input:focus { border-color: rgba(80, 160, 255, 0.4); }
  .result-count {
    font-size: 0.7rem;
    opacity: 0.35;
    white-space: nowrap;
  }

  .tag-bar {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    margin-bottom: 0.75rem;
    align-items: center;
  }
  .tag-bar-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.15rem 0.5rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 99px;
    font-size: 0.7rem;
    color: rgba(255, 255, 255, 0.6);
    cursor: pointer;
    font-family: inherit;
    transition: all 0.15s;
  }
  .tag-bar-chip:hover {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.85);
  }
  .tag-bar-chip.active {
    background: rgba(80, 160, 255, 0.25);
    border-color: rgba(80, 160, 255, 0.4);
    color: rgba(180, 220, 255, 1);
  }
  .tag-bar-count { font-size: 0.6rem; opacity: 0.4; }
  .tag-bar-clear {
    background: none;
    border: none;
    color: rgba(255, 100, 100, 0.6);
    cursor: pointer;
    font-size: 0.65rem;
    font-family: inherit;
    padding: 0.15rem 0.4rem;
  }
  .tag-bar-clear:hover { color: rgba(255, 100, 100, 0.9); }

  .btn {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
    border-radius: 6px;
    padding: 0.45rem 0.75rem;
    cursor: pointer;
    font-size: 0.8rem;
    font-family: inherit;
    transition: background 0.15s;
    text-align: left;
  }
  .btn:hover { background: rgba(255, 255, 255, 0.14); }
  .btn-primary {
    background: rgba(80, 160, 255, 0.25);
    border-color: rgba(80, 160, 255, 0.4);
  }
  .btn-primary:hover { background: rgba(80, 160, 255, 0.35); }
</style>
