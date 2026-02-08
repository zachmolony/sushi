<script lang="ts">
  import {
    selectedAssetIds,
    filteredAssets,
    collections,
    blenderConnected,
    tagsWithCounts,
  } from "./stores";
  import {
    selectAllVisible,
    clearSelection,
    bulkTag,
    bulkAddToCollection,
    bulkSendToBlender,
    bulkSetFavorite,
    createCollection,
    SUGGESTED_TAGS,
    SHELF_ICONS,
  } from "./actions";

  let bulkTagInput = "";
  let showSuggestions = false;
  let showNewCollection = false;
  let newCollectionName = "";
  let newCollectionIcon = "📁";

  // Build suggestions: most-used tags first, then prebuilt tags not yet in the list
  $: existingTagNames = $tagsWithCounts.map((t) => t.name);
  $: topTags = $tagsWithCounts.slice(0, 8).map((t) => t.name);
  $: prebuiltFiltered = SUGGESTED_TAGS.filter((t) => !existingTagNames.includes(t));
  $: allSuggestions = [...topTags, ...prebuiltFiltered.slice(0, Math.max(0, 12 - topTags.length))];
  $: filteredSuggestions = bulkTagInput.trim()
    ? allSuggestions.filter((t) => t.includes(bulkTagInput.trim().toLowerCase()))
    : allSuggestions;

  function handleBulkTagKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      applyBulkTag(bulkTagInput);
    }
    if (e.key === "Escape") {
      showSuggestions = false;
    }
  }

  function applyBulkTag(tag: string) {
    bulkTag(tag);
    bulkTagInput = "";
    showSuggestions = false;
  }

  async function doCreateCollection() {
    if (!newCollectionName.trim()) return;
    await createCollection(newCollectionName, newCollectionIcon);
    newCollectionName = "";
    newCollectionIcon = "📁";
    showNewCollection = false;
  }
</script>

<div class="bulk-bar">
  <div class="bulk-row">
    <span class="bulk-count">{$selectedAssetIds.size} selected</span>
    <button class="btn btn-sm" on:click={selectAllVisible}>
      Select all ({$filteredAssets.length})
    </button>

    <div class="bulk-divider"></div>

    <div class="bulk-tag-wrap">
      <input
        type="text"
        class="bulk-tag-input"
        placeholder="tag all…"
        bind:value={bulkTagInput}
        on:keydown={handleBulkTagKeydown}
        on:focus={() => (showSuggestions = true)}
      />
      {#if bulkTagInput.trim()}
        <button class="btn btn-sm btn-primary" on:click={() => applyBulkTag(bulkTagInput)}>Tag</button>
      {/if}
    </div>

    {#if $collections.length > 0}
      <select
        class="bulk-collection-select"
        on:change={(e) => {
          const target = e.currentTarget;
          const val = parseInt(target.value);
          if (val) bulkAddToCollection(val);
          target.value = "";
        }}
      >
        <option value="">+ tray…</option>
        {#each $collections as col}
          <option value={col.id}>{col.icon} {col.name}</option>
        {/each}
      </select>
    {/if}

    {#if !showNewCollection}
      <button class="btn btn-sm" on:click={() => (showNewCollection = true)} title="Create new collection">+ New Tray</button>
    {/if}

    <div class="bulk-divider"></div>

    <button class="btn btn-sm" on:click={() => bulkSetFavorite(true)} title="Favorite selected">⭐</button>
    <button class="btn btn-sm btn-muted" on:click={() => bulkSetFavorite(false)} title="Unfavorite selected">☆</button>

    {#if $blenderConnected}
      <button class="btn btn-sm btn-primary" on:click={bulkSendToBlender}>🚀 Send to Blender</button>
    {/if}

    <div class="bulk-spacer"></div>
    <button class="btn btn-sm btn-muted" on:click={clearSelection}>✕ Clear</button>
  </div>

  {#if showSuggestions && filteredSuggestions.length > 0}
    <div class="bulk-suggestions">
      <span class="suggestions-label">Suggestions:</span>
      {#each filteredSuggestions as tag}
        <button class="suggestion-chip" on:click={() => applyBulkTag(tag)}>{tag}</button>
      {/each}
    </div>
  {/if}

  {#if showNewCollection}
    <div class="bulk-new-collection">
      <div class="icon-picker-mini">
        {#each SHELF_ICONS.slice(0, 8) as icon}
          <button
            class="icon-option-mini"
            class:active={newCollectionIcon === icon}
            on:click={() => (newCollectionIcon = icon)}>{icon}</button>
        {/each}
      </div>
      <input
        type="text"
        class="bulk-tag-input"
        placeholder="collection name…"
        bind:value={newCollectionName}
        on:keydown={(e) => { if (e.key === "Enter") doCreateCollection(); if (e.key === "Escape") showNewCollection = false; }}
      />
      <button class="btn btn-sm btn-primary" on:click={doCreateCollection}>Create</button>
      <button class="btn btn-sm btn-muted" on:click={() => (showNewCollection = false)}>✕</button>
    </div>
  {/if}
</div>

<style>
  .bulk-bar {
    position: sticky;
    bottom: 0;
    z-index: 50;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    padding: 0.6rem 0.75rem;
    background: rgba(20, 30, 50, 0.95);
    border: 1px solid rgba(80, 160, 255, 0.25);
    border-radius: 10px;
    margin-top: 0.75rem;
    backdrop-filter: blur(12px);
  }
  .bulk-row {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    flex-wrap: wrap;
  }
  .bulk-count {
    font-size: 0.78rem;
    font-weight: 600;
    color: rgba(140, 200, 255, 0.9);
    white-space: nowrap;
  }
  .bulk-divider {
    width: 1px;
    height: 18px;
    background: rgba(255, 255, 255, 0.08);
    flex-shrink: 0;
  }
  .bulk-spacer { flex: 1; }
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

  .bulk-suggestions {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    flex-wrap: wrap;
  }
  .suggestions-label {
    font-size: 0.65rem;
    opacity: 0.35;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    white-space: nowrap;
  }
  .suggestion-chip {
    display: inline-flex;
    align-items: center;
    padding: 0.12rem 0.45rem;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 99px;
    font-size: 0.65rem;
    color: rgba(255, 255, 255, 0.55);
    cursor: pointer;
    font-family: inherit;
    transition: all 0.15s;
  }
  .suggestion-chip:hover {
    background: rgba(80, 160, 255, 0.2);
    border-color: rgba(80, 160, 255, 0.35);
    color: rgba(180, 220, 255, 0.9);
  }

  .bulk-new-collection {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    flex-wrap: wrap;
  }
  .icon-picker-mini { display: flex; gap: 0.1rem; }
  .icon-option-mini {
    width: 22px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid transparent;
    border-radius: 3px;
    cursor: pointer;
    font-size: 0.7rem;
    transition: background 0.15s;
  }
  .icon-option-mini:hover { background: rgba(255, 255, 255, 0.1); }
  .icon-option-mini.active {
    background: rgba(80, 160, 255, 0.2);
    border-color: rgba(80, 160, 255, 0.4);
  }

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
    white-space: nowrap;
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
