<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import Sidebar from "./lib/Sidebar.svelte";
  import AssetGrid from "./lib/AssetGrid.svelte";
  import BulkBar from "./lib/BulkBar.svelte";
  import DetailPanel from "./lib/DetailPanel.svelte";
  import Toast from "./lib/Toast.svelte";
  import Toolbar from "./lib/Toolbar.svelte";
  import TagBar from "./lib/TagBar.svelte";
  import EmptyState from "./lib/EmptyState.svelte";
  import {
    loading,
    assets,
    activeCollectionId,
    filterTag,
    filterTags,
    showBulkActions,
    selectedAssetIds,
  } from "./lib/stores";
  import {
    loadData,
    startBlenderPolling,
    addFolder,
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

<div class="flex h-screen overflow-hidden text-left">
  <Sidebar />

  <main class="flex-1 overflow-y-auto p-6 relative">
    {#if $loading}
      <EmptyState icon="⏳" title="Loading..." />
    {:else if $assets.length === 0 && $activeCollectionId === null && !$filterTag && $filterTags.length === 0}
      <EmptyState
        icon="📦"
        title="No assets yet"
        subtitle="Add a watch folder to start indexing .glb files."
        actionLabel="+ Add Watch Folder"
        onAction={addFolder}
      />
    {:else}
      <Toolbar />
      <TagBar />
      <AssetGrid />

      {#if $showBulkActions && $selectedAssetIds.size > 0}
        <BulkBar />
      {/if}
    {/if}
  </main>

  <DetailPanel />
  <Toast />
</div>
