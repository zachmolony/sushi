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
    if (e.key === "Escape") showTagSuggestions = false;
  }

  function applyTagSuggestion(tag: string) {
    addTag(tag);
    newTagInput = "";
    showTagSuggestions = false;
  }
</script>

{#if $selectedAsset}
  <aside class="w-[300px] min-w-[300px] bg-base-800 p-5 flex flex-col gap-4 border-l border-white/[0.06] overflow-y-auto">
    <!-- Preview -->
    <div class="w-full aspect-square bg-black/30 rounded-lg overflow-hidden flex items-center justify-center">
      {#if $thumbnailCache[$selectedAsset.id]}
        <img
          src={$thumbnailCache[$selectedAsset.id]}
          alt={$selectedAsset.filename}
          class="w-full h-full object-contain"
        />
      {:else}
        <span class="text-xl opacity-20 font-bold">.glb</span>
      {/if}
    </div>

    <!-- Title + fav -->
    <div class="flex items-start gap-2">
      <h3 class="m-0 text-base break-all flex-1">{$selectedAsset.filename}</h3>
      <button
        class="bg-transparent border-none text-xl cursor-pointer p-0 leading-none shrink-0 transition-all hover:scale-110
          {$selectedAsset.favorited === 1 ? 'text-fav' : 'text-white/30 hover:text-yellow-400/80'}"
        on:click={() => toggleFavorite($selectedAsset.id)}
        title={$selectedAsset.favorited === 1 ? "Remove from favorites" : "Add to favorites"}
      >{$selectedAsset.favorited === 1 ? "★" : "☆"}</button>
    </div>

    <!-- Meta -->
    <div class="flex flex-col gap-1.5">
      <div class="flex flex-col gap-0.5">
        <span class="text-[0.65rem] uppercase tracking-wider opacity-40">Path</span>
        <code class="text-[0.78rem] break-all leading-relaxed opacity-80">{$selectedAsset.absolute_path}</code>
      </div>
      <div class="flex flex-col gap-0.5">
        <span class="text-[0.65rem] uppercase tracking-wider opacity-40">Size</span>
        <span class="text-[0.78rem] opacity-80">{formatSize($selectedAsset.file_size)}</span>
      </div>
      <div class="flex flex-col gap-0.5">
        <span class="text-[0.65rem] uppercase tracking-wider opacity-40">Modified</span>
        <span class="text-[0.78rem] opacity-80">{new Date($selectedAsset.modified_at).toLocaleDateString()}</span>
      </div>
    </div>

    <!-- Tags -->
    <div class="flex flex-col gap-1.5">
      <span class="text-[0.65rem] uppercase tracking-wider opacity-40">Tags</span>
      <div class="flex flex-wrap gap-1 items-center">
        {#each $selectedAssetTags as tag}
          <span class="inline-flex items-center gap-1 px-2 py-0.5 bg-surface border border-surface-border rounded-full text-[0.7rem] text-white">
            {tag.name}
            <button
              class="bg-transparent border-none text-white/40 hover:text-red-400 cursor-pointer text-[0.6rem] p-0 leading-none"
              on:click={() => removeTag(tag.id)}
            >✕</button>
          </span>
        {/each}
        <div class="flex-1 min-w-[60px]">
          <input
            type="text"
            class="w-full bg-transparent border-b border-white/10 text-white text-[0.7rem] py-0.5 outline-none font-inherit placeholder:text-white/25 focus:border-b-accent/50"
            placeholder="add tag…"
            bind:value={newTagInput}
            on:keydown={handleTagKeydown}
            on:focus={() => (showTagSuggestions = true)}
          />
        </div>
      </div>
      {#if showTagSuggestions && filteredTagSuggestions.length > 0}
        <div class="flex flex-wrap gap-1">
          {#each filteredTagSuggestions as tag}
            <button
              class="inline-flex items-center px-1.5 py-0.5 bg-white/5 border border-white/[0.06] rounded-full text-[0.6rem] text-white/40 cursor-pointer font-inherit transition-all hover:bg-accent-dim hover:border-accent-border hover:text-blue-200"
              on:click={() => applyTagSuggestion(tag)}
            >{tag}</button>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Collections -->
    <div class="flex flex-col gap-1.5">
      <span class="text-[0.65rem] uppercase tracking-wider opacity-40">Collections</span>
      <div class="flex flex-wrap gap-1 items-center">
        {#each $selectedAssetCollections as col}
          <span class="inline-flex items-center gap-1 px-2 py-0.5 bg-surface border border-surface-border rounded-full text-[0.7rem] text-white">
            {col.icon} {col.name}
            <button
              class="bg-transparent border-none text-white/40 hover:text-red-400 cursor-pointer text-[0.6rem] p-0 leading-none"
              on:click={() => removeSelectedFromCollection(col.id)}
            >✕</button>
          </span>
        {/each}
        {#if $collections.filter((c) => !$selectedAssetCollections.find((sc) => sc.id === c.id)).length > 0}
          <select
            class="bg-white/[0.06] border border-surface-border rounded-full text-white/50 text-[0.7rem] px-2 py-0.5 outline-none font-inherit cursor-pointer"
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

    <!-- Actions -->
    <div class="flex flex-col gap-1.5 mt-auto">
      <button
        class="px-3 py-2 rounded-md text-sm bg-accent-dim border border-accent-border text-white text-left cursor-pointer font-inherit hover:bg-accent-hover transition-colors disabled:opacity-35 disabled:cursor-not-allowed"
        on:click={sendToBlender}
        title={$blenderConnected ? "Import into Blender" : "Blender addon not running"}
        disabled={!$blenderConnected}
      >🚀 Send to Blender</button>
      <button
        class="px-3 py-2 rounded-md text-sm bg-surface border border-surface-border text-white text-left cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
        on:click={copyPath}
      >📋 Copy Path</button>
      <button
        class="px-3 py-2 rounded-md text-sm bg-surface border border-surface-border text-white text-left cursor-pointer font-inherit hover:bg-surface-hover transition-colors opacity-60"
        on:click={showInFolder}
      >📂 Show in Folder</button>
    </div>
  </aside>
{/if}