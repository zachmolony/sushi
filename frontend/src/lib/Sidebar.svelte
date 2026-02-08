<script lang="ts">
  import {
    watchFolders,
    allTags,
    filterTag,
    filterTags,
    activeCollectionId,
    activeView,
    collections,
    assets,
  } from "./stores";
  import {
    addFolder,
    removeFolder,
    toggleTagFilter,
    clearTagFilters,
    setCollectionFilter,
    setActiveView,
    createCollection,
    deleteCollectionById,
    SHELF_ICONS,
  } from "./actions";
  import { blenderConnected } from "./stores";
  import type { ViewId } from "./stores";

  let showNewCollection = false;
  let newCollectionName = "";
  let newCollectionIcon = "📁";

  const views: { id: ViewId; label: string; icon: string }[] = [
    { id: "all", label: "All Assets", icon: "📦" },
    { id: "untagged", label: "Untagged", icon: "🏷️" },
    { id: "recent-added", label: "Recently Added", icon: "🆕" },
    { id: "recent-used", label: "Recently Used", icon: "🕐" },
    { id: "favorites", label: "Favorites", icon: "⭐" },
  ];

  function handleNewCollectionKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") doCreateCollection();
    if (e.key === "Escape") showNewCollection = false;
  }

  async function doCreateCollection() {
    await createCollection(newCollectionName, newCollectionIcon);
    newCollectionName = "";
    newCollectionIcon = "📁";
    showNewCollection = false;
  }
</script>

<aside class="sidebar">
  <h2>sushi 🍣</h2>

  <div class="blender-status">
    <div class="status-dot" class:connected={$blenderConnected}></div>
    <span class="status-text">
      {$blenderConnected ? "Blender connected" : "Blender not found"}
    </span>
  </div>

  <div class="sidebar-section">
    <h3>Views</h3>
    <div class="view-list">
      {#each views as view}
        <button
          class="view-item"
          class:active={$activeView === view.id && $activeCollectionId === null && $filterTags.length === 0 && !$filterTag}
          on:click={() => setActiveView(view.id)}
        >
          <span class="view-icon">{view.icon}</span>
          <span class="view-name">{view.label}</span>
        </button>
      {/each}
    </div>
  </div>

  <div class="sidebar-section">
    <h3>Watch Folders</h3>
    <button class="btn btn-sm btn-add" on:click={addFolder}>+ Add Folder</button>
    {#if $watchFolders.length > 0}
      <div class="folder-list">
        {#each $watchFolders as folder}
          <div class="folder-item">
            <span class="folder-path" title={folder.path}>
              {folder.path.split("/").pop()}
            </span>
            <button
              class="folder-remove"
              on:click={() => removeFolder(folder.id)}
              title="Remove this folder">✕</button>
          </div>
        {/each}
      </div>
    {:else}
      <p class="hint">No folders watched yet.</p>
    {/if}
  </div>

  {#if $allTags.length > 0}
    <div class="sidebar-section">
      <h3>Tags</h3>
      <div class="tag-filter-list">
        <button
          class="tag-chip"
          class:active={$filterTag === "" && $filterTags.length === 0 && $activeCollectionId === null}
          on:click={() => clearTagFilters()}>All</button>
        {#each $allTags as tag}
          <button
            class="tag-chip"
            class:active={$filterTags.includes(tag.name) || $filterTag === tag.name}
            on:click={() => toggleTagFilter(tag.name)}>{tag.name}</button>
        {/each}
      </div>
    </div>
  {/if}

  <div class="sidebar-section">
    <h3>Collections</h3>
    <div class="collection-list">
      {#each $collections as col}
        <button
          class="collection-item"
          class:active={$activeCollectionId === col.id}
          on:click={() => setCollectionFilter(col.id)}
        >
          <span class="collection-icon">{col.icon}</span>
          <span class="collection-name">{col.name}</span>
          <span class="collection-count">{col.asset_count}</span>
          <button class="collection-delete" on:click|stopPropagation={() => deleteCollectionById(col.id)} title="Delete collection">✕</button>
        </button>
      {/each}
      {#if showNewCollection}
        <div class="new-collection-form">
          <div class="icon-picker">
            {#each SHELF_ICONS as icon}
              <button
                class="icon-option"
                class:active={newCollectionIcon === icon}
                on:click={() => (newCollectionIcon = icon)}>{icon}</button>
            {/each}
          </div>
          <input
            type="text"
            class="tag-input"
            placeholder="collection name…"
            bind:value={newCollectionName}
            on:keydown={handleNewCollectionKeydown}
          />
          <div class="new-collection-actions">
            <button class="btn btn-sm btn-primary" on:click={doCreateCollection}>Create</button>
            <button class="btn btn-sm btn-muted" on:click={() => (showNewCollection = false)}>Cancel</button>
          </div>
        </div>
      {:else}
        <button class="btn btn-sm btn-add" on:click={() => (showNewCollection = true)}>+ New Collection</button>
      {/if}
    </div>
  </div>

  <div class="sidebar-spacer"></div>
  <div class="sidebar-footer">
    <span class="hint">{$assets.length} asset{$assets.length !== 1 ? "s" : ""}</span>
  </div>
</aside>

<style>
  .sidebar {
    width: 220px;
    min-width: 220px;
    background: rgba(20, 28, 40, 1);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    border-right: 1px solid rgba(255, 255, 255, 0.06);
    overflow-y: auto;
  }
  .sidebar h2 {
    margin: 0;
    font-size: 1.3rem;
    letter-spacing: 0.05em;
  }
  .sidebar-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .sidebar-section h3 {
    margin: 0;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    opacity: 0.5;
  }
  .sidebar-spacer { flex: 1; }
  .sidebar-footer {
    padding-top: 0.5rem;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
  }

  .blender-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: rgba(255, 80, 80, 0.8);
    flex-shrink: 0;
  }
  .status-dot.connected {
    background: rgba(80, 220, 120, 0.9);
  }
  .status-text {
    font-size: 0.75rem;
    opacity: 0.6;
  }

  .btn-add {
    width: 100%;
    text-align: center;
    background: rgba(80, 160, 255, 0.15);
    border-color: rgba(80, 160, 255, 0.3);
  }
  .btn-add:hover { background: rgba(80, 160, 255, 0.25); }
  .folder-list {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }
  .folder-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.25rem;
    padding: 0.3rem 0.5rem;
    background: rgba(255, 255, 255, 0.04);
    border-radius: 4px;
    font-size: 0.75rem;
  }
  .folder-path {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    opacity: 0.7;
  }
  .folder-remove {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.3);
    cursor: pointer;
    font-size: 0.65rem;
    padding: 0 0.2rem;
    flex-shrink: 0;
  }
  .folder-remove:hover { color: rgba(255, 100, 100, 0.9); }

  .tag-filter-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
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
  .tag-chip.active {
    background: rgba(80, 160, 255, 0.25);
    border-color: rgba(80, 160, 255, 0.4);
  }

  .collection-list {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  .collection-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0.5rem;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid transparent;
    border-radius: 5px;
    cursor: pointer;
    color: inherit;
    font-family: inherit;
    font-size: 0.78rem;
    transition: background 0.15s;
    text-align: left;
  }
  .collection-item:hover { background: rgba(255, 255, 255, 0.07); }
  .collection-item.active {
    background: rgba(80, 160, 255, 0.15);
    border-color: rgba(80, 160, 255, 0.3);
  }
  .collection-icon { font-size: 0.9rem; flex-shrink: 0; }
  .collection-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .collection-count { font-size: 0.65rem; opacity: 0.35; flex-shrink: 0; }
  .collection-delete {
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.2);
    cursor: pointer;
    font-size: 0.55rem;
    padding: 0 0.15rem;
    flex-shrink: 0;
    opacity: 0;
    transition: opacity 0.15s, color 0.15s;
  }
  .collection-item:hover .collection-delete { opacity: 1; }
  .collection-delete:hover { color: rgba(255, 100, 100, 0.9); }

  .view-list {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }
  .view-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0.5rem;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid transparent;
    border-radius: 5px;
    cursor: pointer;
    color: inherit;
    font-family: inherit;
    font-size: 0.78rem;
    transition: background 0.15s;
    text-align: left;
  }
  .view-item:hover { background: rgba(255, 255, 255, 0.07); }
  .view-item.active {
    background: rgba(80, 160, 255, 0.15);
    border-color: rgba(80, 160, 255, 0.3);
  }
  .view-icon { font-size: 0.85rem; flex-shrink: 0; }
  .view-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .new-collection-form {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    padding: 0.5rem;
    background: rgba(255, 255, 255, 0.04);
    border-radius: 6px;
    margin-top: 0.25rem;
  }
  .icon-picker { display: flex; flex-wrap: wrap; gap: 0.15rem; }
  .icon-option {
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid transparent;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: background 0.15s;
  }
  .icon-option:hover { background: rgba(255, 255, 255, 0.1); }
  .icon-option.active {
    background: rgba(80, 160, 255, 0.2);
    border-color: rgba(80, 160, 255, 0.4);
  }
  .new-collection-actions { display: flex; gap: 0.3rem; }

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

  .hint {
    font-size: 0.75rem;
    opacity: 0.4;
    line-height: 1.4;
  }
</style>
