<script lang="ts">
  import { toastVisible, toastMessage } from "./stores";
</script>

{#if $toastVisible}
  <div class="fixed bottom-6 left-1/2 -translate-x-1/2 bg-accent/90 text-gray-900 px-5 py-2 rounded-lg text-sm font-semibold z-[100] animate-[toast-in_0.2s_ease]">
    {$toastMessage}
  </div>
{/if}

<style>
  @keyframes toast-in {
    from { opacity: 0; transform: translateX(-50%) translateY(10px); }
    to { opacity: 1; transform: translateX(-50%) translateY(0); }
  }
</style>
