<script lang="ts">
  import { toastVisible, toastMessage } from "./stores";
</script>

{#if $toastVisible}
  <div class="toast">{$toastMessage}</div>
{/if}

<style>
  .toast {
    position: fixed;
    bottom: 1.5rem;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(100, 180, 255, 0.9);
    color: #111;
    padding: 0.5rem 1.25rem;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    z-index: 100;
    animation: toast-in 0.2s ease;
  }
  @keyframes toast-in {
    from { opacity: 0; transform: translateX(-50%) translateY(10px); }
    to { opacity: 1; transform: translateX(-50%) translateY(0); }
  }
</style>
