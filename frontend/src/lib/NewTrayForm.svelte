<script lang="ts">
  import { createCollection, SHELF_ICONS } from "./actions";

  export let compact: boolean = false;
  export let onClose: () => void;

  let name = "";
  let icon = "📁";

  async function doCreate() {
    if (!name.trim()) return;
    await createCollection(name, icon);
    name = "";
    icon = "📁";
    onClose();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") doCreate();
    if (e.key === "Escape") onClose();
  }
</script>

{#if compact}
  <!-- Inline row layout (BulkBar) -->
  <div class="flex items-center gap-2 flex-wrap">
    <div class="flex gap-0.5">
      {#each SHELF_ICONS.slice(0, 8) as ic}
        <button
          class="w-[22px] h-[22px] flex items-center justify-center rounded cursor-pointer text-[0.7rem] transition-colors border
            {icon === ic ? 'bg-accent-dim border-accent-border' : 'bg-white/[0.04] border-transparent hover:bg-white/10'}"
          on:click={() => (icon = ic)}>{ic}</button
        >
      {/each}
    </div>
    <input
      type="text"
      class="w-28 bg-surface border border-surface-border rounded text-white text-xs px-2 py-1 outline-none font-inherit placeholder:text-white/30 focus:border-accent/50"
      placeholder="tray name…"
      bind:value={name}
      on:keydown={handleKeydown}
    />
    <button
      class="px-2.5 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
      on:click={doCreate}>Create</button
    >
    <button
      class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white/60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
      on:click={onClose}>✕</button
    >
  </div>
{:else}
  <!-- Stacked layout (Sidebar) -->
  <div class="flex flex-col gap-1.5 p-2 bg-white/[0.04] rounded-md mt-1">
    <div class="flex flex-wrap gap-0.5">
      {#each SHELF_ICONS as ic}
        <button
          class="w-[26px] h-[26px] flex items-center justify-center rounded cursor-pointer text-sm transition-colors border
            {icon === ic ? 'bg-accent-dim border-accent-border' : 'bg-white/[0.04] border-transparent hover:bg-white/10'}"
          on:click={() => (icon = ic)}>{ic}</button
        >
      {/each}
    </div>
    <input
      type="text"
      class="w-full bg-transparent border-b border-white/10 text-white text-[0.7rem] py-0.5 outline-none font-inherit placeholder:text-white/25 focus:border-b-accent/50"
      placeholder="tray name…"
      bind:value={name}
      on:keydown={handleKeydown}
    />
    <div class="flex gap-1">
      <button
        class="px-2 py-1 rounded-md text-xs bg-accent-dim border border-accent-border text-white cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
        on:click={doCreate}>Create</button
      >
      <button
        class="px-2 py-1 rounded-md text-xs bg-surface border border-surface-border text-white opacity-60 cursor-pointer font-inherit hover:bg-surface-hover transition-colors"
        on:click={onClose}>Cancel</button
      >
    </div>
  </div>
{/if}
