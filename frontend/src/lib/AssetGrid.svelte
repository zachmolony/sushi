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
  <div class="empty-state">
    <div class="empty-icon">🔍</div>
    <h2>No matches</h2>
    <p>No assets match the current filter.</p>
  </div>
{:else}
  <div class="asset-grid">
    {#each $filteredAssets as asset, i}
      <button
        class="asset-card"
        class:selected={$selectedAsset?.id === asset.id}
        class:multi-selected={$selectedAssetIds.has(asset.id)}
        on:click={(e) => handleAssetClick(asset, i, e)}
        on:mouseenter={() => hoverAssetId.set(asset.id)}
        on:mouseleave={() => hoverAssetId.set(null)}
      >
        <div class="asset-thumb">
          {#if $thumbnailCache[asset.id]}
            <img
              src={$thumbnailCache[asset.id]}
              alt={asset.filename}
              class="thumb-img"
            />
          {:else}
            <span class="asset-ext">.glb</span>
          {/if}
          {#if $hoverAssetId === asset.id && $collections.length > 0 && !$showBulkActions}
            <div class="hover-actions">
              {#each $collections.slice(0, 4) as col}
                <button
                  class="hover-action-btn"
                  title="Add to {col.name}"
                  on:click={(e) => hoverAddToCollection(asset.id, col.id, e)}
                >{col.icon}</button>
              {/each}
            </div>
          {/if}
        </div>
        <div class="asset-name" title={asset.absolute_path}>
          {#if asset.favorited === 1}<span class="fav-star">★</span>{/if}
          {asset.filename}
        </div>
      </button>
    {/each}
  </div>
{/if}

<style>
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

  .asset-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 0.75rem;
  }
  .asset-card {
    background: rgba(255, 255, 255, 0.03);
    border: 2px solid transparent;
    border-radius: 8px;
    padding: 0;
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
    text-align: center;
    color: inherit;
    font-family: inherit;
    overflow: hidden;
  }
  .asset-card:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.08);
  }
  .asset-card.selected {
    border-color: rgba(100, 180, 255, 0.6);
    background: rgba(100, 180, 255, 0.06);
  }
  .asset-card.multi-selected {
    border-color: rgba(80, 200, 120, 0.6);
    background: rgba(80, 200, 120, 0.08);
  }
  .asset-card.multi-selected:hover {
    border-color: rgba(80, 200, 120, 0.8);
  }
  .asset-thumb {
    height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.2);
    overflow: hidden;
    position: relative;
  }
  .thumb-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .asset-ext {
    font-size: 1.2rem;
    opacity: 0.2;
    font-weight: 700;
  }
  .asset-name {
    padding: 0.4rem 0.5rem;
    font-size: 0.75rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .fav-star {
    color: rgba(255, 200, 50, 0.8);
    font-size: 0.65rem;
    margin-right: 0.15rem;
  }

  .hover-actions {
    position: absolute;
    bottom: 4px;
    right: 4px;
    display: flex;
    gap: 2px;
    z-index: 10;
  }
  .hover-action-btn {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.7);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.7rem;
    padding: 0;
    transition: all 0.15s;
    color: white;
  }
  .hover-action-btn:hover {
    background: rgba(80, 160, 255, 0.5);
    border-color: rgba(80, 160, 255, 0.6);
    transform: scale(1.1);
  }
</style>
