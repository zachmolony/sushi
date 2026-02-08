<script lang="ts">
  import {
    tagsWithCounts,
    filterTags,
    excludeTags,
  } from "./stores";
  import {
    toggleTagFilter,
    toggleExcludeTag,
    clearTagFilters,
  } from "./actions";
</script>

{#if $tagsWithCounts.length > 0}
  <div class="flex flex-wrap gap-1.5 mb-3 items-center">
    {#each $tagsWithCounts as tag}
      <button
        class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[0.7rem] border cursor-pointer font-inherit transition-all
          {$excludeTags.includes(tag.name)
          ? 'bg-red-950/40 border-red-800/40 text-red-400/70 line-through'
          : $filterTags.includes(tag.name)
            ? 'bg-accent-dim border-accent-border text-blue-200'
            : 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white/85'}"
        on:click={() => toggleTagFilter(tag.name)}
        on:contextmenu|preventDefault={() => toggleExcludeTag(tag.name)}
        title="Click to include · Right-click to exclude"
      >
        {#if $excludeTags.includes(tag.name)}
          <span class="text-[0.6rem] opacity-60 no-underline" style="text-decoration: none">⊘</span>
        {/if}
        {tag.name}
        <span class="text-[0.6rem] opacity-40">{tag.count}</span>
      </button>
    {/each}
    {#if $filterTags.length > 0 || $excludeTags.length > 0}
      <button
        class="bg-transparent border-none text-red-400/60 hover:text-red-400/90 cursor-pointer text-[0.65rem] font-inherit px-1.5 py-0.5"
        on:click={clearTagFilters}>✕ clear</button
      >
    {/if}
  </div>
{/if}
