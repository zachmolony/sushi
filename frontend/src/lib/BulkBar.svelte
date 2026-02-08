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
    bulkDeleteAssets,
    createCollection,
    SUGGESTED_TAGS,
    SHELF_ICONS,
  } from "./actions";

  let bulkTagInput = "";
  let showSuggestions = false;
  let showNewCollection = false;
  let newCollectionName = "";
  let newCollectionIcon = "📁";
  let confirmBulkDelete = false;

  // Build suggestions: most-used tags first, then prebuilt tags not yet in the list
  $: existingTagNames = $tagsWithCounts.map((t) => t.name);
  $: topTags = $tagsWithCounts.slice(0, 8).map((t) => t.name);
  $: prebuiltFiltered = SUGGESTED_TAGS.filter(
    (t) => !existingTagNames.includes(t),
  );
  $: allSuggestions = [
    ...topTags,
    ...prebuiltFiltered.slice(0, Math.max(0, 12 - topTags.length)),
  ];
  $: filteredSuggestions = bulkTagInput.trim()
    ? allSuggestions.filter((t) =>
        t.includes(bulkTagInput.trim().toLowerCase()),
      )
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

<div
  class="sticky bottom-0 z-50 flex flex-col gap-2 px-4 py-3 mt-3 bg-base-800/95 border border-accent/25 rounded-xl backdrop-blur-xl"
>
  <!-- Main row -->
  <div class="flex items-center gap-2 flex-wrap">
    <span
      class="text-[0.78rem] font-semibold text-blue-300/90 whitespace-nowrap"
    >
      {$selectedAssetIds.size} selected
    </span>
    <button
      class="px-2.5 py-1 rounded-md text-xs bg-surface border border-surface-border text-white cursor-pointer font-inherit hover:bg-surface-hover transition-colors whitespace-nowrap"
      on:click={selectAllVisible}>Select all ({$filteredAssets.length})</button
    >

    <div class="w-px h-[18px] bg-white/[0.08] shrink-0"></div>

    <!-- Tag input -->
    <div class="flex items-center gap-1">
      <input
        type="text"
        class="w-24 bg-surface border border-surface-border rounded text-white text-xs px-2 py-1 outline-none font-inherit placeholder:text-white/30 focus:border-accent/50 transition-colors"
        placeholder="tag all…"
        bind:value={bulkTagInput}
        on:keydown={handleBulkTagKeydown}
        on:focus={() => (showSuggestions = true)}
      />
      {#if bulkTagInput.trim()}
        <button
          class="px-2.5 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
          on:click={() => applyBulkTag(bulkTagInput)}>Tag</button
        >
      {/if}
    </div>

    <!-- Collection picker -->
    {#if $collections.length > 0}
      <select
        class="bg-surface border border-surface-border rounded text-white/70 text-xs px-2 py-1 outline-none font-inherit cursor-pointer"
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
      <button
        class="px-2.5 py-1 rounded-md text-xs bg-surface border border-surface-border text-white cursor-pointer font-inherit hover:bg-surface-hover transition-colors whitespace-nowrap"
        on:click={() => (showNewCollection = true)}
        title="Create new collection">+ New Tray</button
      >
    {/if}

    <div class="w-px h-[18px] bg-white/[0.08] shrink-0"></div>

    <!-- Favorites -->
    <button
      class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
      on:click={() => bulkSetFavorite(true)}
      title="Favorite selected">⭐</button
    >
    <button
      class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white/60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
      on:click={() => bulkSetFavorite(false)}
      title="Unfavorite selected">☆</button
    >

    {#if $blenderConnected}
      <button
        class="px-2.5 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors whitespace-nowrap"
        on:click={bulkSendToBlender}>🚀 Send to Blender</button
      >
    {/if}

    <div class="flex-1"></div>
    {#if !confirmBulkDelete}
      <button
        class="px-2.5 py-1 rounded-md text-xs bg-surface border border-red-900/40 text-red-400/70 cursor-pointer font-inherit hover:bg-red-950/40 hover:text-red-400 transition-colors whitespace-nowrap"
        on:click={() => (confirmBulkDelete = true)}
        title="Delete selected assets from library">🗑 Delete</button
      >
    {:else}
      <button
        class="px-2.5 py-1 rounded-md text-xs bg-red-900/40 border border-red-700/50 text-red-300 cursor-pointer font-inherit hover:bg-red-800/50 transition-colors whitespace-nowrap"
        on:click={() => {
          confirmBulkDelete = false;
          bulkDeleteAssets();
        }}>Confirm delete {$selectedAssetIds.size}?</button
      >
      <button
        class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white/60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
        on:click={() => (confirmBulkDelete = false)}>No</button
      >
    {/if}
    <button
      class="px-2.5 py-1 rounded-md text-xs bg-surface border border-surface-border text-white/60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
      on:click={clearSelection}>✕ Clear</button
    >
  </div>

  <!-- Tag suggestions -->
  {#if showSuggestions && filteredSuggestions.length > 0}
    <div class="flex items-center gap-1.5 flex-wrap">
      <span
        class="text-[0.65rem] uppercase tracking-wide opacity-35 whitespace-nowrap"
        >Suggestions:</span
      >
      {#each filteredSuggestions as tag}
        <button
          class="inline-flex items-center px-2 py-0.5 bg-white/[0.06] border border-white/[0.08] rounded-full text-[0.65rem] text-white/55 cursor-pointer font-inherit transition-all hover:bg-accent-dim hover:border-accent-border hover:text-blue-200"
          on:click={() => applyBulkTag(tag)}>{tag}</button
        >
      {/each}
    </div>
  {/if}

  <!-- New collection inline form -->
  {#if showNewCollection}
    <div class="flex items-center gap-2 flex-wrap">
      <div class="flex gap-0.5">
        {#each SHELF_ICONS.slice(0, 8) as icon}
          <button
            class="w-[22px] h-[22px] flex items-center justify-center rounded cursor-pointer text-[0.7rem] transition-colors border
              {newCollectionIcon === icon
              ? 'bg-accent-dim border-accent-border'
              : 'bg-white/[0.04] border-transparent hover:bg-white/10'}"
            on:click={() => (newCollectionIcon = icon)}>{icon}</button
          >
        {/each}
      </div>
      <input
        type="text"
        class="w-28 bg-surface border border-surface-border rounded text-white text-xs px-2 py-1 outline-none font-inherit placeholder:text-white/30 focus:border-accent/50"
        placeholder="collection name…"
        bind:value={newCollectionName}
        on:keydown={(e) => {
          if (e.key === "Enter") doCreateCollection();
          if (e.key === "Escape") showNewCollection = false;
        }}
      />
      <button
        class="px-2.5 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
        on:click={doCreateCollection}>Create</button
      >
      <button
        class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white/60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
        on:click={() => (showNewCollection = false)}>✕</button
      >
    </div>
  {/if}
</div>
