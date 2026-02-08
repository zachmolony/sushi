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
    activeFolderPath,
  } from "./stores";
  import {
    addFolder,
    removeFolder,
    toggleTagFilter,
    clearTagFilters,
    setCollectionFilter,
    setActiveView,
    deleteCollectionById,
    clearFolderFilter,
    regenerateAllThumbnails,
  } from "./actions";
  import { blenderConnected } from "./stores";
  import type { ViewId } from "./stores";
  import FolderTree from "./FolderTree.svelte";
  import NewTrayForm from "./NewTrayForm.svelte";

  let showNewCollection = false;

  // Build folder tree from asset paths
  interface FolderNode {
    name: string;
    fullPath: string;
    children: FolderNode[];
    assetCount: number;
    watchFolderId?: number;
  }

  function buildFolderTree(
    watchFoldersList: { id: number; path: string }[],
    assetList: { absolute_path: string }[],
  ): FolderNode[] {
    const roots: FolderNode[] = [];

    for (const wf of watchFoldersList) {
      const root: FolderNode = {
        name: wf.path.split("/").pop() || wf.path,
        fullPath: wf.path,
        children: [],
        assetCount: 0,
        watchFolderId: wf.id,
      };

      // Get all unique subdirectories for this watch folder
      const subDirs = new Set<string>();
      for (const asset of assetList) {
        if (asset.absolute_path.startsWith(wf.path + "/")) {
          const relative = asset.absolute_path.slice(wf.path.length + 1);
          const parts = relative.split("/");
          // Build all intermediate paths
          for (let i = 1; i < parts.length; i++) {
            subDirs.add(parts.slice(0, i).join("/"));
          }
        }
      }

      // Count assets directly in root
      root.assetCount = assetList.filter((a) => {
        if (!a.absolute_path.startsWith(wf.path + "/")) return false;
        const relative = a.absolute_path.slice(wf.path.length + 1);
        return !relative.includes("/");
      }).length;

      // Build tree from sorted subdirs
      const sortedDirs = Array.from(subDirs).sort();
      const nodeMap: Record<string, FolderNode> = {};

      for (const dir of sortedDirs) {
        const parts = dir.split("/");
        const folderName = parts[parts.length - 1];
        const fullPath = wf.path + "/" + dir;
        const dirAssetCount = assetList.filter((a) => {
          if (!a.absolute_path.startsWith(fullPath + "/")) return false;
          const after = a.absolute_path.slice(fullPath.length + 1);
          return !after.includes("/");
        }).length;

        const node: FolderNode = {
          name: folderName,
          fullPath,
          children: [],
          assetCount: dirAssetCount,
        };
        nodeMap[dir] = node;

        if (parts.length === 1) {
          root.children.push(node);
        } else {
          const parentDir = parts.slice(0, -1).join("/");
          if (nodeMap[parentDir]) {
            nodeMap[parentDir].children.push(node);
          }
        }
      }

      roots.push(root);
    }

    return roots;
  }

  $: folderTree = buildFolderTree($watchFolders, $assets);

  // Track which folders are expanded
  let expandedFolders: Set<string> = new Set();

  function toggleExpand(path: string) {
    if (expandedFolders.has(path)) {
      expandedFolders.delete(path);
    } else {
      expandedFolders.add(path);
    }
    expandedFolders = expandedFolders; // trigger reactivity
  }

  const views: { id: ViewId; label: string; icon: string }[] = [
    { id: "all", label: "All Assets", icon: "📦" },
    { id: "untagged", label: "Untagged", icon: "🏷️" },
    { id: "recent-added", label: "Recently Added", icon: "🆕" },
    { id: "recent-used", label: "Recently Used", icon: "🕐" },
    { id: "favorites", label: "Favorites", icon: "⭐" },
  ];
</script>

<aside
  class="w-[220px] min-w-[220px] bg-base-800 p-5 flex flex-col gap-5 border-r border-white/[0.06] overflow-y-auto"
>
  <h2 class="m-0 text-xl tracking-wide">sushi 🍣</h2>

  <!-- Blender status -->
  <div class="flex items-center gap-2">
    <div
      class="w-2 h-2 rounded-full shrink-0 {$blenderConnected
        ? 'bg-green-400'
        : 'bg-red-400'}"
    ></div>
    <span class="text-xs opacity-60">
      {$blenderConnected ? "Blender connected" : "Blender not found"}
    </span>
  </div>

  <!-- Views -->
  <div class="flex flex-col gap-2">
    <h3
      class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold"
    >
      Views
    </h3>
    <div class="flex flex-col gap-0.5">
      {#each views as view}
        <button
          class="flex items-center gap-2 px-2 py-1.5 rounded text-[0.78rem] cursor-pointer font-inherit text-left border transition-colors
            {$activeView === view.id &&
          $activeCollectionId === null &&
          $filterTags.length === 0 &&
          !$filterTag
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

  <!-- Folder Browser -->
  <div class="flex flex-col gap-2">
    <h3
      class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold"
    >
      Folders
    </h3>
    <button
      class="w-full text-center px-2 py-1.5 rounded-md text-xs bg-accent-glow border border-accent-border/75 text-white cursor-pointer font-inherit hover:bg-accent-dim transition-colors"
      on:click={addFolder}>+ Add Folder</button
    >
    {#if folderTree.length > 0}
      <div class="flex flex-col gap-0">
        <FolderTree
          nodes={folderTree}
          depth={0}
          {expandedFolders}
          onToggleExpand={toggleExpand}
          onRemoveFolder={removeFolder}
        />
      </div>
      {#if $activeFolderPath}
        <button
          class="text-[0.65rem] text-white/40 hover:text-white/70 bg-transparent border-none cursor-pointer font-inherit px-1 py-0.5 text-left transition-colors"
          on:click={clearFolderFilter}>✕ clear folder filter</button
        >
      {/if}
    {:else}
      <p class="text-xs opacity-40 leading-relaxed">No folders watched yet.</p>
    {/if}
  </div>

  <!-- Tags -->
  {#if $allTags.length > 0}
    <div class="flex flex-col gap-2">
      <h3
        class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold"
      >
        Tags
      </h3>
      <div class="flex flex-wrap gap-1">
        <button
          class="inline-flex items-center px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-colors
            {$filterTag === '' &&
          $filterTags.length === 0 &&
          $activeCollectionId === null
            ? 'bg-accent-dim border-accent-border text-white'
            : 'bg-surface border-surface-border text-white hover:bg-surface-hover'}"
          on:click={() => clearTagFilters()}>All</button
        >
        {#each $allTags as tag}
          <button
            class="inline-flex items-center px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-colors
              {$filterTags.includes(tag.name) || $filterTag === tag.name
              ? 'bg-accent-dim border-accent-border text-white'
              : 'bg-surface border-surface-border text-white hover:bg-surface-hover'}"
            on:click={() => toggleTagFilter(tag.name)}>{tag.name}</button
          >
        {/each}
      </div>
    </div>
  {/if}

  <!-- Trays -->
  <div class="flex flex-col gap-2">
    <h3
      class="m-0 text-[0.7rem] uppercase tracking-wider opacity-50 font-semibold"
    >
      Trays
    </h3>
    <div class="flex flex-col gap-0.5">
      {#each $collections as col}
        <div
          class="flex items-center gap-1.5 px-2 py-1.5 rounded text-[0.78rem] cursor-pointer font-inherit text-left border transition-colors group
            {$activeCollectionId === col.id
            ? 'bg-accent-glow border-accent-border text-white'
            : 'bg-surface-dim border-transparent text-inherit hover:bg-white/[0.07]'}"
          on:click={() => setCollectionFilter(col.id)}
          on:keydown={(e) => {
            if (e.key === "Enter") setCollectionFilter(col.id);
          }}
          role="button"
          tabindex="0"
        >
          <span class="text-[0.9rem] shrink-0">{col.icon}</span>
          <span class="flex-1 truncate">{col.name}</span>
          <span class="text-[0.65rem] opacity-35 shrink-0"
            >{col.asset_count}</span
          >
          <button
            class="bg-transparent border-none text-white/20 hover:text-red-400 cursor-pointer text-[0.55rem] p-0 px-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            on:click|stopPropagation={() => deleteCollectionById(col.id)}
            title="Delete tray">✕</button
          >
        </div>
      {/each}
      {#if showNewCollection}
        <NewTrayForm compact={false} onClose={() => (showNewCollection = false)} />
      {:else}
        <button
          class="w-full text-center px-2 py-1.5 rounded-md text-xs bg-accent-glow border border-accent-border/75 text-white cursor-pointer font-inherit hover:bg-accent-dim transition-colors"
          on:click={() => (showNewCollection = true)}>+ New Tray</button
        >
      {/if}
    </div>
  </div>

  <div class="flex-1"></div>
  <div class="pt-2 border-t border-white/[0.06] flex flex-col gap-2">
    <span class="text-xs opacity-40"
      >{$assets.length} asset{$assets.length !== 1 ? "s" : ""}</span
    >
    <button
      class="text-[0.65rem] text-white/30 hover:text-white/60 bg-transparent border-none cursor-pointer font-inherit p-0 text-left transition-colors"
      on:click={regenerateAllThumbnails}
      title="Clear all cached thumbnails and regenerate them with current settings"
      >🔄 Regenerate thumbnails</button
    >
  </div>
</aside>
