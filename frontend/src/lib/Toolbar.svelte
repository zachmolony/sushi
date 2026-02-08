<script lang="ts">
  import {
    searchQuery,
    filteredAssets,
    sortField,
    sortDirection,
  } from "./stores";
  import {
    viewLabel,
    selectAllVisible,
    setSort,
  } from "./actions";
  import type { SortField } from "./stores";

  function handleSortChange(e: Event) {
    const val = (e.currentTarget as HTMLSelectElement).value;
    setSort(val as SortField);
  }
</script>

<div class="flex items-center gap-3 mb-4">
  <h2 class="m-0 text-sm whitespace-nowrap opacity-70 font-semibold">
    {viewLabel()}
  </h2>
  <input
    type="text"
    class="flex-1 max-w-[260px] bg-surface border border-surface-border rounded-md text-white text-xs px-3 py-1.5 outline-none font-inherit placeholder:text-white/25 focus:border-accent/40 transition-colors"
    placeholder="Search files…"
    bind:value={$searchQuery}
  />
  <div class="flex items-center gap-1">
    <select
      class="appearance-none bg-surface border border-surface-border rounded text-white/70 text-[0.72rem] px-2 py-1 outline-none font-inherit cursor-pointer"
      value={$sortField}
      on:change={handleSortChange}
    >
      <option value="name">Name</option>
      <option value="date-added">Date Added</option>
      <option value="file-modified">Modified</option>
      <option value="file-size">Size</option>
      <option value="poly-count">Polys</option>
    </select>
    <button
      class="rounded text-white/60 text-sm px-1.5 py-1 cursor-pointer font-inherit hover:bg-surface-hover transition-colors leading-none"
      on:click={() => setSort($sortField)}
      title="Toggle sort direction"
      >{$sortDirection === "asc" ? "↑" : "↓"}</button
    >
  </div>
  <span class="text-[0.7rem] opacity-35 whitespace-nowrap">
    {$filteredAssets.length} result{$filteredAssets.length !== 1 ? "s" : ""}
  </span>
  <button
    class="text-[0.68rem] px-1.5 py-0.5 rounded bg-white/[0.06] border border-white/[0.08] text-white/40 hover:text-white/70 hover:bg-white/10 cursor-pointer font-inherit transition-colors whitespace-nowrap"
    on:click={selectAllVisible}
    title="Select all visible assets (for bulk tagging / trays)"
    >☑ Select all</button
  >
</div>
