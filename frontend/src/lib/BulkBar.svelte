<script lang="ts">
  import {
    selectedAssetIds,
    filteredAssets,
    collections,
    blenderConnected,
  } from "./stores";
  import {
    selectAllVisible,
    clearSelection,
    bulkTag,
    bulkAddToCollection,
    bulkSendToBlender,
  } from "./actions";

  let bulkTagInput = "";

  function handleBulkTagKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      bulkTag(bulkTagInput);
      bulkTagInput = "";
    }
    if (e.key === "Escape") clearSelection();
  }
</script>

<div class="bulk-bar">
  <span class="bulk-count">{$selectedAssetIds.size} selected</span>
  <button class="btn btn-sm" on:click={selectAllVisible}>
    Select all ({$filteredAssets.length})
  </button>
  <div class="bulk-tag-wrap">
    <input
      type="text"
      class="bulk-tag-input"
      placeholder="tag all…"
      bind:value={bulkTagInput}
      on:keydown={handleBulkTagKeydown}
    />
    {#if bulkTagInput.trim()}
      <button class="btn btn-sm btn-primary" on:click={() => { bulkTag(bulkTagInput); bulkTagInput = ""; }}>Tag</button>
    {/if}
  </div>
  {#if $collections.length > 0}
    <select
      class="bulk-collection-select"
      on:change={(e) => {
        const val = parseInt(e.currentTarget.value);
        if (val) bulkAddToCollection(val);
        e.currentTarget.value = "";
      }}
    >
      <option value="">+ collection…</option>
      {#each $collections as col}
        <option value={col.id}>{col.icon} {col.name}</option>
      {/each}
    </select>
  {/if}
  {#if $blenderConnected}
    <button class="btn btn-sm btn-primary" on:click={bulkSendToBlender}>🚀 Send to Blender</button>
  {/if}
  <button class="btn btn-sm btn-muted" on:click={clearSelection}>✕</button>
</div>

<style>
  .bulk-bar {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: rgba(80, 160, 255, 0.1);
    border: 1px solid rgba(80, 160, 255, 0.2);
    border-radius: 8px;
    margin-bottom: 0.75rem;
    flex-wrap: wrap;
  }
  .bulk-count {
    font-size: 0.78rem;
    font-weight: 600;
    color: rgba(140, 200, 255, 0.9);
    white-space: nowrap;
  }
  .bulk-tag-wrap { display: flex; align-items: center; gap: 0.3rem; }
  .bulk-tag-input {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 5px;
    color: white;
    font-size: 0.75rem;
    padding: 0.25rem 0.5rem;
    outline: none;
    font-family: inherit;
    width: 100px;
  }
  .bulk-tag-input::placeholder { color: rgba(255, 255, 255, 0.3); }
  .bulk-tag-input:focus { border-color: rgba(80, 160, 255, 0.5); }
  .bulk-collection-select {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 5px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.75rem;
    padding: 0.25rem 0.4rem;
    outline: none;
    font-family: inherit;
    cursor: pointer;
  }
  .bulk-collection-select option { background: #1a1a2e; color: white; }

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
  .btn-sm { padding: 0.3rem 0.6rem; font-size: 0.75rem; }
  .btn-muted { opacity: 0.6; }
  .btn-primary {
    background: rgba(80, 160, 255, 0.25);
    border-color: rgba(80, 160, 255, 0.4);
  }
  .btn-primary:hover { background: rgba(80, 160, 255, 0.35); }
</style>
