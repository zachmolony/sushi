<script lang="ts">
  import {
    filteredAssets,
    activeView,
    activeCollectionId,
    activeFolderPath,
    filterTags,
    excludeTags,
    searchQuery,
  } from "./stores";
  import AssetCard from "./AssetCard.svelte";

  let recentExpanded = false;

  const ONE_WEEK = 7 * 24 * 60 * 60 * 1000;

  $: showRecentSection =
    $activeView === "all" &&
    !$activeCollectionId &&
    !$activeFolderPath &&
    !$filterTags.length &&
    !$excludeTags.length &&
    !$searchQuery.trim();

  $: cutoff = Date.now() - ONE_WEEK;
  $: recentAssets = showRecentSection
    ? $filteredAssets.filter((a) => new Date(a.modified_at).getTime() > cutoff)
    : [];
  $: recentIds = new Set(recentAssets.map((a) => a.id));
  $: restAssets = showRecentSection
    ? $filteredAssets.filter((a) => !recentIds.has(a.id))
    : $filteredAssets;
  $: showRecent = recentAssets.length > 0;
</script>

{#if $filteredAssets.length === 0}
  <div class="flex flex-col items-center justify-center h-full opacity-70">
    <div class="text-5xl mb-4">🔍</div>
    <h2 class="text-base font-semibold m-0 mb-1">No matches</h2>
    <p class="text-sm opacity-60 my-1">No assets match the current filter.</p>
  </div>
{:else}
  {#if showRecent}
    <div class="mb-6">
      <div class="flex items-center gap-2 mb-2">
        <span class="text-[0.7rem] uppercase tracking-wider opacity-40 font-semibold">Recently Added</span>
        <span class="text-[0.6rem] opacity-25">{recentAssets.length}</span>
        {#if recentAssets.length > 12}
          <button
            class="text-[0.65rem] px-1.5 py-0.5 rounded bg-white/[0.06] border border-white/[0.08] text-white/40 hover:text-white/70 hover:bg-white/10 cursor-pointer font-inherit transition-colors ml-auto"
            on:click={() => (recentExpanded = !recentExpanded)}
          >
            {recentExpanded ? "▲ Collapse" : "▼ Show all"}
          </button>
        {/if}
      </div>
      <div
        class="grid grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-3 overflow-hidden transition-all duration-300"
        style={!recentExpanded && recentAssets.length > 12 ? "max-height: 310px" : ""}
      >
        {#each recentAssets as asset, i}
          <AssetCard {asset} index={i} />
        {/each}
      </div>
      {#if restAssets.length > 0}
        <div class="border-b border-white/[0.06] mt-5 mb-2"></div>
        <span class="text-[0.7rem] uppercase tracking-wider opacity-40 font-semibold mb-2 block">All Assets</span>
      {/if}
    </div>
  {/if}

  {#if restAssets.length > 0 || !showRecent}
    <div class="grid grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-3">
      {#each restAssets as asset, i}
        <AssetCard {asset} index={i + (showRecent ? recentAssets.length : 0)} />
      {/each}
    </div>
  {/if}
{/if}
