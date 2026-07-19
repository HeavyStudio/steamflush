<script lang="ts">
  // Import global styles
  import './App.css';
  
  // Import components
  import Header from './components/Header.svelte';
  import StatusMessage from './components/StatusMessage.svelte';
  import OrphanCard from './components/OrphanCard.svelte';

  // Import and initialize state logic
  import { createAppState } from './App.svelte.ts';
  const appState = createAppState();

  function formatBytes(bytes: number) {
    if (bytes == 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
</script>

<Header />

<main>
  <section class="content">
    <StatusMessage 
      isLoading={appState.isLoading} 
      errorMessage={appState.errorMessage} 
      isEmpty={appState.orphans.length === 0} 
      onRefresh={appState.refreshScan} 
    />

    {#if !appState.isLoading && !appState.errorMessage && appState.orphans.length > 0}
      <div class="results">
        <div class="results-header">
          <h2>Detected Orphan AppIDs ({appState.orphans.length})</h2>
          <div class="actions-wrapper">
            <button class="btn-refresh" onclick={appState.refreshScan} disabled={appState.isLoading}>
              Scan
            </button>
            <button class="btn-danger-all" onclick={appState.handleDeleteAll} disabled={appState.isLoading}>
              Delete All
            </button>
          </div>
        </div>
        
        <p class="warning-text">These games are no longer installed but still occupy {formatBytes(appState.totalSize)}:</p>
        <ul>
          {#each appState.orphans as app (app.appID)}
            <OrphanCard
              name={app.name}
              id={app.appID}
              size={app.size}
              disabled={appState.isLoading}
              onDelete={appState.handleDelete}
            />
          {/each}
        </ul>
      </div>
    {/if}
  </section>
</main>