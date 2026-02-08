<script lang="ts">
  import {
    filteredAssets,
    selectedAsset,
    selectedAssetIds,
    showBulkActions,
    thumbnailCache,
    collections,
    hoverAssetId,
  } from "./stores";
  import { handleAssetClick, hoverAddToCollection } from "./actions";
</script>

{#if $filteredAssets.length === 0}
  <div class="flex flex-col items-center justify-center h-full opacity-70">
    <div class="text-5xl mb-4">🔍</div>
    <h2 class="text-base font-semibold m-0 mb-1">No matches</h2>
    <p class="text-sm opacity-60 my-1">No assets match the current filter.</p>
  </div>
{:else}
  <div class="grid grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-3">
    {#each $filteredAssets as asset, i}
      <button
        class="bg-surface-dim border-2 rounded-lg p-0 cursor-pointer transition-all text-center text-inherit font-inherit overflow-hidden
          {$selectedAssetIds.has(asset.id)
            ? 'border-green-400/60 bg-green-400/[0.08] hover:border-green-400/80'
            : $selectedAsset?.id === asset.id
              ? 'border-blue-400/60 bg-blue-400/[0.06]'
              : 'border-transparent hover:bg-white/[0.06] hover:border-white/[0.08]'}"
        on:click={(e) => handleAssetClick(asset, i, e)}
        on:mouseenter={() => hoverAssetId.set(asset.id)}
        on:mouseleave={() => hoverAssetId.set(null)}
      >
        <div class="h-[120px] flex items-center justify-center bg-black/20 overflow-hidden relative">
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
            <div class="absolute top-1.5 left-1.5 w-5 h-5 rounded-full bg-green-500 flex items-center justify-center text-white text-[0.6rem] font-bold shadow-md">✓</div>
          {/if}
          {#if $hoverAssetId === asset.id && $collections.length > 0 && !$showBulkActions}
            <div class="absolute bottom-1 right-1 flex gap-0.5 z-10">
              {#each $collections.slice(0, 4) as col}
                <button
                  class="w-6 h-6 flex items-center justify-center bg-black/70 border border-white/15 rounded cursor-pointer text-[0.7rem] p-0 transition-all text-white hover:bg-accent/50 hover:border-accent/60 hover:scale-110"
                  title="Add to {col.name}"
                  on:click={(e) => hoverAddToCollection(asset.id, col.id, e)}
                >{col.icon}</button>
              {/each}
            </div>
          {/if}
        </div>
        <div class="px-2 py-1.5 text-xs whitespace-nowrap overflow-hidden text-ellipsis">
          {#if asset.favorited === 1}<span class="text-fav text-[0.65rem] mr-0.5">★</span>{/if}
          {asset.filename}
        </div>
      </button>
    {/each}
  </div>
{/if}
