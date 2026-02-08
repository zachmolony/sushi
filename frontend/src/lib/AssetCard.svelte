<script lang="ts">
  import {
    selectedAsset,
    selectedAssetIds,
    showBulkActions,
    thumbnailCache,
    collections,
    hoverAssetId,
  } from "./stores";
  import {
    handleAssetClick,
    hoverAddToCollection,
    toggleAssetSelection,
    formatPoly,
  } from "./actions";

  export let asset: any;
  export let index: number;
</script>

<button
  class="bg-surface-dim border-2 rounded-lg p-0 cursor-pointer transition-all text-center text-inherit font-inherit overflow-hidden
    {$selectedAssetIds.has(asset.id)
    ? 'border-green-400/60 bg-green-400/[0.08] hover:border-green-400/80'
    : $selectedAsset?.id === asset.id
      ? 'border-blue-400/60 bg-blue-400/[0.06]'
      : 'border-transparent hover:bg-white/[0.06] hover:border-white/[0.08]'}"
  on:click={(e) => handleAssetClick(asset, index, e)}
  on:mouseenter={() => hoverAssetId.set(asset.id)}
  on:mouseleave={() => hoverAssetId.set(null)}
>
  <div
    class="h-[120px] flex items-center justify-center bg-black/20 overflow-hidden relative"
  >
    {#if $thumbnailCache[asset.id]}
      <img
        src={$thumbnailCache[asset.id]}
        alt={asset.filename}
        class="w-full h-full object-cover"
      />
    {:else}
      <span class="text-xl opacity-20 font-bold">.glb</span>
    {/if}
    {#if $selectedAssetIds.has(asset.id)}
      <button
        class="absolute top-1.5 left-1.5 w-5 h-5 rounded bg-green-500 flex items-center justify-center text-white text-[0.7rem] font-bold shadow-md z-20 cursor-pointer border-none p-0 hover:bg-green-400 transition-colors"
        on:click|stopPropagation={() => toggleAssetSelection(asset.id)}
        title="Deselect">✓</button
      >
    {:else if $hoverAssetId === asset.id || $showBulkActions}
      <button
        class="absolute top-1.5 left-1.5 w-5 h-5 rounded bg-black/50 border-2 border-white/30 flex items-center justify-center z-20 cursor-pointer p-0 hover:border-green-400/70 hover:bg-black/70 transition-colors backdrop-blur-sm"
        on:click|stopPropagation={() => toggleAssetSelection(asset.id)}
        title="Select"
      ></button>
    {/if}
    {#if asset.poly_count > 0}
      <div
        class="absolute top-1.5 right-1.5 px-1.5 py-0.5 rounded bg-black/60 text-[0.6rem] text-white/80 font-mono leading-none backdrop-blur-sm"
      >
        {formatPoly(asset.poly_count)} △
      </div>
    {/if}
    {#if $hoverAssetId === asset.id && $collections.length > 0 && !$showBulkActions}
      <div class="absolute bottom-1 right-1 flex gap-0.5 z-10">
        {#each $collections.slice(0, 4) as col}
          <button
            class="w-6 h-6 flex items-center justify-center bg-black/70 border border-white/15 rounded cursor-pointer text-[0.7rem] p-0 transition-all text-white hover:bg-accent/50 hover:border-accent/60 hover:scale-110"
            title="Add to {col.name}"
            on:click={(e) => hoverAddToCollection(asset.id, col.id, e)}
            >{col.icon}</button
          >
        {/each}
      </div>
    {/if}
  </div>
  <div
    class="px-2 py-1.5 text-xs whitespace-nowrap overflow-hidden text-ellipsis"
  >
    {#if asset.favorited === 1}<span class="text-fav text-[0.65rem] mr-0.5"
        >★</span
      >{/if}
    {asset.filename}
  </div>
</button>
