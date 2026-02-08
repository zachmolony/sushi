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

<aside class="w-[220px] min-w-[220px] bg-base-800 p-5 flex flex-col gap-5 border-r border-white/[0.06] overflow-y-auto">
  <h2 class="m-0 text-xl tracking-wide">sushi 🍣</h2>

  <!-- Blender status -->
  <div class="flex items-center gap-2">
    <div class="w-2 h-2 rounded-full shrink-0 {$blenderConnected ? 'bg-green-400' : 'bg-red-400'}"></div>
    <span class="text-xs opacity-60">
      {$blenderConnected ? "Blender connected" : "Blender not found"}
    </span>
  </div>

  <!-- Views -->
  <div class="flex flex-col gap-2">
    <h3 class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold">Views</h3>
    <div class="flex flex-col gap-0.5">
      {#each views as view}
        <button
          class="flex items-center gap-2 px-2 py-1.5 rounded text-[0.78rem] cursor-pointer font-inherit text-left border transition-colors
            {$activeView === view.id && $activeCollectionId === null && $filterTags.length === 0 && !$filterTag
              ? 'bg-accent-glow border-accent-border text-white'
              : 'bg-surface-dim border-transparent text-inherit hover:bg-white/[0.07]'}"
          on:click={() => setActiveView(view.id)}
        >
          <span class="text-[0.85rem] shrink-0">{view.icon}</span>
          <span class="flex-1 truncate">{view.label}</span>
        </button>
      {/each}
    </div>
  </div>

  <!-- Watch Folders -->
  <div class="flex flex-col gap-2">
    <h3 class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold">Watch Folders</h3>
    <button
      class="w-full text-center px-2 py-1.5 rounded-md text-xs bg-accent-glow border border-accent-border/75 text-white cursor-pointer font-inherit hover:bg-accent-dim transition-colors"
      on:click={addFolder}
    >+ Add Folder</button>
    {#if $watchFolders.length > 0}
      <div class="flex flex-col gap-0.5">
        {#each $watchFolders as folder}
          <div class="flex items-center justify-between gap-1 px-2 py-1 bg-white/[0.04] rounded text-xs group">
            <span class="truncate flex-1 opacity-70" title={folder.path}>
              {folder.path.split("/").pop()}
            </span>
            <button
              class="bg-transparent border-none text-white/30 hover:text-red-400 cursor-pointer text-[0.65rem] p-0 px-0.5 shrink-0"
              on:click={() => removeFolder(folder.id)}
              title="Remove this folder"
            >✕</button>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-xs opacity-40 leading-relaxed">No folders watched yet.</p>
    {/if}
  </div>

  <!-- Tags -->
  {#if $allTags.length > 0}
    <div class="flex flex-col gap-2">
      <h3 class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold">Tags</h3>
      <div class="flex flex-wrap gap-1">
        <button
          class="inline-flex items-center px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-colors
            {$filterTag === '' && $filterTags.length === 0 && $activeCollectionId === null
              ? 'bg-accent-dim border-accent-border text-white'
              : 'bg-surface border-surface-border text-white hover:bg-surface-hover'}"
          on:click={() => clearTagFilters()}
        >All</button>
        {#each $allTags as tag}
          <button
            class="inline-flex items-center px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-colors
              {$filterTags.includes(tag.name) || $filterTag === tag.name
                ? 'bg-accent-dim border-accent-border text-white'
                : 'bg-surface border-surface-border text-white hover:bg-surface-hover'}"
            on:click={() => toggleTagFilter(tag.name)}
          >{tag.name}</button>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Collections -->
  <div class="flex flex-col gap-2">
    <h3 class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold">Collections</h3>
    <div class="flex flex-col gap-0.5">
      {#each $collections as col}
        <div
          class="flex items-center gap-1.5 px-2 py-1.5 rounded text-[0.78rem] cursor-pointer font-inherit text-left border transition-colors group
            {$activeCollectionId === col.id
              ? 'bg-accent-glow border-accent-border text-white'
              : 'bg-surface-dim border-transparent text-inherit hover:bg-white/[0.07]'}"
          on:click={() => setCollectionFilter(col.id)}
          on:keydown={(e) => { if (e.key === 'Enter') setCollectionFilter(col.id); }}
          role="button"
          tabindex="0"
        >
          <span class="text-[0.9rem] shrink-0">{col.icon}</span>
          <span class="flex-1 truncate">{col.name}</span>
          <span class="text-[0.65rem] opacity-35 shrink-0">{col.asset_count}</span>
          <button
            class="bg-transparent border-none text-white/20 hover:text-red-400 cursor-pointer text-[0.55rem] p-0 px-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            on:click|stopPropagation={() => deleteCollectionById(col.id)}
            title="Delete collection"
          >✕</button>
        </div>
      {/each}
      {#if showNewCollection}
        <div class="flex flex-col gap-1.5 p-2 bg-white/[0.04] rounded-md mt-1">
          <div class="flex flex-wrap gap-0.5">
            {#each SHELF_ICONS as icon}
              <button
                class="w-[26px] h-[26px] flex items-center justify-center rounded cursor-pointer text-sm transition-colors border
                  {newCollectionIcon === icon
                    ? 'bg-accent-dim border-accent-border'
                    : 'bg-white/[0.04] border-transparent hover:bg-white/10'}"
                on:click={() => (newCollectionIcon = icon)}
              >{icon}</button>
            {/each}
          </div>
          <input
            type="text"
            class="w-full bg-transparent border-b border-white/10 text-white text-[0.7rem] py-0.5 outline-none font-inherit placeholder:text-white/25 focus:border-b-accent/50"
            placeholder="collection name…"
            bind:value={newCollectionName}
            on:keydown={handleNewCollectionKeydown}
          />
          <div class="flex gap-1">
            <button
              class="px-2 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
              on:click={doCreateCollection}
            >Create</button>
            <button
              class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white opacity-60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
              on:click={() => (showNewCollection = false)}
            >Cancel</button>
          </div>
        </div>
      {:else}
        <button
          class="w-full text-center px-2 py-1.5 rounded-md text-xs bg-accent-glow border border-accent-border/75 text-white cursor-pointer font-inherit hover:bg-accent-dim transition-colors"
          on:click={() => (showNewCollection = true)}
        >+ New Collection</button>
      {/if}
    </div>
  </div>

  <div class="flex-1"></div>
  <div class="pt-2 border-t border-white/[0.06]">
    <span class="text-xs opacity-40">{$assets.length} asset{$assets.length !== 1 ? "s" : ""}</span>
  </div>
</aside>
