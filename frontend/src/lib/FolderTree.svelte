<script lang="ts">
  import { activeFolderPath } from "./stores";
  import { setFolderFilter } from "./actions";

  interface FolderNode {
    name: string;
    fullPath: string;
    children: FolderNode[];
    assetCount: number;
    watchFolderId?: number;
  }

  export let nodes: FolderNode[] = [];
  export let depth: number = 0;
  export let expandedFolders: Set<string>;
  export let onToggleExpand: (path: string) => void;
  export let onRemoveFolder: ((id: number) => void) | null = null;
</script>

{#each nodes as node}
  <div
    class="flex flex-col"
    style="padding-left: {depth > 0 ? '0.75rem' : '0'}"
  >
    <div class="flex items-center gap-0.5 group">
      {#if node.children.length > 0}
        <button
          class="bg-transparent border-none text-white/30 text-[0.6rem] p-0 w-4 h-4 flex items-center justify-center cursor-pointer shrink-0 hover:text-white/60"
          on:click|stopPropagation={() => onToggleExpand(node.fullPath)}
          >{expandedFolders.has(node.fullPath) ? "▼" : "▶"}</button
        >
      {:else}
        <span class="w-4"></span>
      {/if}
      <button
        class="flex-1 flex items-center gap-1.5 px-1.5 py-1 rounded text-[0.73rem] cursor-pointer font-inherit text-left border transition-colors truncate
          {$activeFolderPath === node.fullPath
          ? 'bg-accent-glow border-accent-border text-white'
          : 'bg-transparent border-transparent text-inherit hover:bg-white/[0.06]'}"
        on:click={() => setFolderFilter(node.fullPath)}
        title={node.fullPath}
      >
        <span class="shrink-0 text-[0.75rem]">{depth === 0 ? "📂" : "📁"}</span>
        <span class="truncate flex-1">{node.name}</span>
        {#if node.assetCount > 0}
          <span class="text-[0.6rem] opacity-30 shrink-0"
            >{node.assetCount}</span
          >
        {/if}
      </button>
      {#if onRemoveFolder && node.watchFolderId != null}
        <button
          class="bg-transparent border-none text-white/20 hover:text-red-400 cursor-pointer text-[0.55rem] p-0 px-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
          on:click|stopPropagation={() => onRemoveFolder(node.watchFolderId)}
          title="Remove watch folder">✕</button
        >
      {/if}
    </div>
    {#if expandedFolders.has(node.fullPath) && node.children.length > 0}
      <svelte:self
        nodes={node.children}
        depth={depth + 1}
        {expandedFolders}
        {onToggleExpand}
        onRemoveFolder={null}
      />
    {/if}
  </div>
{/each}
