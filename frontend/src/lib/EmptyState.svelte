<script lang="ts">
  export let icon: string = "🔍";
  export let title: string;
  export let subtitle: string = "";
  export let actionLabel: string = "";
  export let onAction: (() => void) | null = null;
</script>

<div class="flex flex-col items-center justify-center h-full opacity-70">
  <div class="text-5xl mb-4">{icon}</div>
  <h2 class="text-base font-semibold m-0 mb-1">{title}</h2>
  {#if subtitle}
    <p class="text-sm opacity-60 my-1">{subtitle}</p>
  {/if}
  {#if actionLabel && onAction}
    <button
      class="mt-4 px-4 py-2 rounded-md bg-accent-dim border border-accent-border text-white text-sm cursor-pointer font-inherit hover:bg-accent-hover transition-colors"
      on:click={onAction}>{actionLabel}</button
    >
  {/if}
</div>
