<script lang="ts">
  interface Props {
    isLoading: boolean;
    errorMessage: string;
    isSteamFound: boolean;
    isEmpty: boolean;
    onRefresh: () => void;
  }

  let { isLoading, errorMessage, isSteamFound, isEmpty, onRefresh }: Props = $props();
</script>

{#if isLoading}
  <p class="status" aria-live="polite">Scanning your storage...</p>
{:else if errorMessage}
  <div class="status-container" role="alert" aria-live="assertive">
    <p class="status error">Error: {errorMessage}</p>
    <button class="btn-refresh" onclick={onRefresh}>Retry Scan</button>
  </div>
{:else if !isSteamFound}
  <div class="status-container" role="alert" aria-live="assertive">
    <p class="status error">Steam installation not found.</p>
    <p>Please ensure Steam is installed and configured correctly.</p>
  </div>
{:else if isEmpty}
  <div class="status-container" role="status" aria-live="polite">
    <p class="status success">Everything is clean! No orphan folders found.</p>
    <button class="btn-refresh" onclick={onRefresh}>Scan Device</button>
  </div>
{/if}

<style>
  .status {
    font-size: 1.2rem;
    font-weight: 500;
  }

  .status-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }

  .error { color: var(--color-danger); }
  .success { color: var(--color-success); }

  .btn-refresh {
    background-color: var(--color-accent);
    color: var(--color-text-on-accent);
    border: none;
    padding: 8px 16px;
    min-height: 40px;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
    transition: background-color 0.2s ease;
  }

  .btn-refresh:hover {
    background-color: var(--color-accent-hover);
  }
</style>