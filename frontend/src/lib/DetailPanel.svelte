<script lang="ts">
  import {
    selectedAsset,
    selectedAssetTags,
    selectedAssetCollections,
    collections,
    thumbnailCache,
    blenderConnected,
    tagsWithCounts,
  } from "./stores";
  import {
    sendToBlender,
    copyPath,
    showInFolder,
    addTag,
    removeTag,
    addSelectedToCollection,
    removeSelectedFromCollection,
    toggleFavorite,
    formatSize,
    SUGGESTED_TAGS,
  } from "./actions";

  let newTagInput = "";
  let showTagSuggestions = false;

  // Suggestions: existing tags not already on the asset, then prebuilt
  $: assetTagNames = $selectedAssetTags.map((t) => t.name);
  $: existingNotApplied = $tagsWithCounts.filter((t) => !assetTagNames.includes(t.name)).map((t) => t.name).slice(0, 6);
  $: prebuiltNotApplied = SUGGESTED_TAGS.filter((t) => !assetTagNames.includes(t) && !existingNotApplied.includes(t));
  $: tagSuggestions = [...existingNotApplied, ...prebuiltNotApplied.slice(0, Math.max(0, 8 - existingNotApplied.length))];
  $: filteredTagSuggestions = newTagInput.trim()
    ? tagSuggestions.filter((t) => t.includes(newTagInput.trim().toLowerCase()))
    : tagSuggestions;

  function handleTagKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      addTag(newTagInput);
      newTagInput = "";
    }
    if (e.key === "Escape") {
      showTagSuggestions = false;
    }
  }

  function applyTagSuggestion(tag: string) {
    addTag(tag);
    newTagInput = "";
    showTagSuggestions = false;
  }
</script>

{#if $selectedAsset}
  <aside class="detail-panel">
    <div class="detail-preview">
      {#if $thumbnailCache[$selectedAsset.id]}
        <img
          src={$thumbnailCache[$selectedAsset.id]}
          alt={$selectedAsset.filename}
          class="detail-thumb-img"
        />
      {:else}
        <div class="preview-placeholder">
          <span class="asset-ext">.glb</span>
        </div>
      {/if}
    </div>

    <div class="detail-title-row">
      <h3>{$selectedAsset.filename}</h3>
      <button
        class="fav-btn"
        class:is-fav={$selectedAsset.favorited === 1}
        on:click={() => toggleFavorite($selectedAsset.id)}
        title={$selectedAsset.favorited === 1 ? "Remove from favorites" : "Add to favorites"}
      >
        {$selectedAsset.favorited === 1 ? "★" : "☆"}
      </button>
    </div>

    <div class="detail-meta">
      <div class="meta-row">
        <span class="meta-label">Path</span>
        <code class="meta-value">{$selectedAsset.absolute_path}</code>
      </div>
      <div class="meta-row">
        <span class="meta-label">Size</span>
        <span class="meta-value">{formatSize($selectedAsset.file_size)}</span>
      </div>
      <div class="meta-row">
        <span class="meta-label">Modified</span>
        <span class="meta-value">{new Date($selectedAsset.modified_at).toLocaleDateString()}</span>
      </div>
    </div>

    <div class="detail-tags">
      <span class="meta-label">Tags</span>
      <div class="tag-list">
        {#each $selectedAssetTags as tag}
          <span class="tag-chip tag-removable">
            {tag.name}
            <button class="tag-x" on:click={() => removeTag(tag.id)}>✕</button>
          </span>
        {/each}
        <div class="tag-input-wrap">
          <input
            type="text"
            class="tag-input"
            placeholder="add tag…"
            bind:value={newTagInput}
            on:keydown={handleTagKeydown}
            on:focus={() => (showTagSuggestions = true)}
          />
        </div>
      </div>
      {#if showTagSuggestions && filteredTagSuggestions.length > 0}
        <div class="tag-suggestions">
          {#each filteredTagSuggestions as tag}
            <button class="suggestion-chip" on:click={() => applyTagSuggestion(tag)}>{tag}</button>
          {/each}
        </div>
      {/if}
    </div>

    <div class="detail-collections">
      <span class="meta-label">Collections</span>
      <div class="collection-chips">
        {#each $selectedAssetCollections as col}
          <span class="tag-chip tag-removable">
            {col.icon} {col.name}
            <button class="tag-x" on:click={() => removeSelectedFromCollection(col.id)}>✕</button>
          </span>
        {/each}
        {#if $collections.filter((c) => !$selectedAssetCollections.find((sc) => sc.id === c.id)).length > 0}
          <select
            class="collection-select"
            on:change={(e) => {
              const val = parseInt(e.currentTarget.value);
              if (val) addSelectedToCollection(val);
              e.currentTarget.value = "";
            }}
          >
            <option value="">+ add to…</option>
            {#each $collections.filter((c) => !$selectedAssetCollections.find((sc) => sc.id === c.id)) as col}
              <option value={col.id}>{col.icon} {col.name}</option>
            {/each}
          </select>
        {/if}
      </div>
    </div>

    <div class="detail-actions">
      <button
        class="btn btn-primary"
        on:click={sendToBlender}
        title={$blenderConnected ? "Import into Blender" : "Blender addon not running"}
        disabled={!$blenderConnected}>🚀 Send to Blender</button>
      <button class="btn" on:click={copyPath}>📋 Copy Path</button>
      <button class="btn btn-muted" on:click={showInFolder}>📂 Show in Folder</button>
    </div>
  </aside>
{/if}

<style>
  .detail-panel {
    width: 300px;
    min-width: 300px;
    background: rgba(20, 28, 40, 1);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    border-left: 1px solid rgba(255, 255, 255, 0.06);
    overflow-y: auto;
  }
  .detail-preview {
    width: 100%;
    aspect-ratio: 1;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .detail-thumb-img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
  .detail-panel h3 {
    margin: 0;
    font-size: 1rem;
    word-break: break-all;
  }
  .detail-title-row {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
  }
  .detail-title-row h3 { flex: 1; }
  .fav-btn {
    background: none;
    border: none;
    font-size: 1.3rem;
    cursor: pointer;
    color: rgba(255, 255, 255, 0.3);
    padding: 0;
    line-height: 1;
    transition: color 0.15s, transform 0.15s;
    flex-shrink: 0;
  }
  .fav-btn:hover { color: rgba(255, 200, 50, 0.8); transform: scale(1.15); }
  .fav-btn.is-fav { color: rgba(255, 200, 50, 0.9); }
  .detail-meta {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }
  .meta-row {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }
  .meta-label {
    font-size: 0.65rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    opacity: 0.4;
  }
  .meta-value {
    font-size: 0.78rem;
    word-break: break-all;
    line-height: 1.4;
    opacity: 0.8;
  }
  .preview-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0.3;
  }
  .asset-ext {
    font-size: 1.2rem;
    opacity: 0.2;
    font-weight: 700;
  }

  .detail-tags, .detail-collections {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }
  .tag-list, .collection-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    align-items: center;
  }
  .tag-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.2rem 0.5rem;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 99px;
    font-size: 0.7rem;
    color: white;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }
  .tag-chip:hover { background: rgba(255, 255, 255, 0.14); }
  .tag-removable { cursor: default; }
  .tag-x {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    font-size: 0.6rem;
    padding: 0;
    line-height: 1;
  }
  .tag-x:hover { color: rgba(255, 100, 100, 0.9); }
  .tag-input-wrap { flex: 1; min-width: 60px; }
  .tag-input {
    width: 100%;
    background: transparent;
    border: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
    font-size: 0.7rem;
    padding: 0.2rem 0;
    outline: none;
    font-family: inherit;
  }
  .tag-input::placeholder { color: rgba(255, 255, 255, 0.25); }
  .tag-input:focus { border-bottom-color: rgba(80, 160, 255, 0.5); }

  .tag-suggestions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.2rem;
  }
  .suggestion-chip {
    display: inline-flex;
    align-items: center;
    padding: 0.1rem 0.4rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 99px;
    font-size: 0.6rem;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    font-family: inherit;
    transition: all 0.15s;
  }
  .suggestion-chip:hover {
    background: rgba(80, 160, 255, 0.18);
    border-color: rgba(80, 160, 255, 0.3);
    color: rgba(180, 220, 255, 0.85);
  }

  .collection-select {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 99px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.7rem;
    padding: 0.2rem 0.4rem;
    outline: none;
    font-family: inherit;
    cursor: pointer;
  }
  .collection-select option { background: #1a1a2e; color: white; }

  .detail-actions {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-top: auto;
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
  }
  .btn:hover { background: rgba(255, 255, 255, 0.14); }
  .btn-muted { opacity: 0.6; }
  .btn-primary {
    background: rgba(80, 160, 255, 0.25);
    border-color: rgba(80, 160, 255, 0.4);
  }
  .btn-primary:hover { background: rgba(80, 160, 255, 0.35); }
  .btn-primary:disabled { opacity: 0.35; cursor: not-allowed; }
</style>
