<script lang="ts">
  interface Props {
    isLoading: boolean;
    errorMessage: string;
    isSteamFound: boolean; // Ajout de la prop
    isEmpty: boolean;
    onRefresh: () => void;
  }

  let { isLoading, errorMessage, isSteamFound, isEmpty, onRefresh }: Props = $props();
</script>

{#if isLoading}
  <p class="status">Scanning your storage...</p>
{:else if errorMessage}
  <div class="status-container">
    <p class="status error">Error: {errorMessage}</p>
    <button class="btn-refresh" onclick={onRefresh}>Retry Scan</button>
  </div>
{:else if !isSteamFound} <!-- Nouvelle condition prioritaire -->
  <div class="status-container">
    <p class="status error">Steam installation not found.</p>
    <p>Please ensure Steam is installed and configured correctly.</p>
  </div>
{:else if isEmpty}
  <div class="status-container">
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

  .error { color: #f7768e; }
  .success { color: #9ece6a; }

  .btn-refresh {
    background-color: #7aa2f7;
    color: #1a1b26;
    border: none;
    padding: 6px 14px;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
    transition: background-color 0.2s ease;
  }

  .btn-refresh:hover {
    background-color: #6593f5;
  }
</style>