<script lang="ts">
  // Import global styles
  import "./App.css";

  // Import components
  import Header from "./components/Header.svelte";
  import StatusMessage from "./components/StatusMessage.svelte";
  import OrphanCard from "./components/OrphanCard.svelte";
  import { formatBytes } from "./lib/format.ts";

  // Import and initialize state logic
  import { createAppState } from "./App.svelte.ts";
  import Footer from "./components/Footer.svelte";

  const appState = createAppState();

  $effect(() => {
    appState.refreshScan();
  });
</script>

<div class="app-container">
  <Header />

  <main>
    <section class="content">
      <StatusMessage
        isLoading={appState.isLoading}
        errorMessage={appState.errorMessage}
        successMessage={appState.successMessage}
        isSteamFound={appState.isSteamFound}
        isEmpty={appState.orphans.length === 0}
        onRefresh={appState.refreshScan}
      />

      {#if !appState.isLoading && !appState.errorMessage && appState.orphans.length > 0}
        <div class="results">
          <div class="results-header">
            <h2>Detected Orphan AppIDs ({appState.orphans.length})</h2>

            <div class="toolbar">
              <button
                class:active={appState.sortOption === "name"}
                onclick={() => appState.setSort("name")}
              >
                Name
              </button>
              <button
                class:active={appState.sortOption === "size"}
                onclick={() => appState.setSort("size")}
              >
                Size
              </button>
              <button
                class:active={appState.sortOption === "id"}
                onclick={() => appState.setSort("id")}
              >
                ID
              </button>
            </div>

            <div class="actions-wrapper">
              <button
                class="btn-refresh"
                onclick={appState.refreshScan}
                disabled={appState.isLoading}
              >
                Scan
              </button>
              <button
                class="btn-danger-all"
                onclick={appState.handleDeleteAll}
                disabled={appState.isLoading}
              >
                Delete All
              </button>
            </div>
          </div>

          <p class="warning-text">
            These games are no longer installed but still occupy {formatBytes(
              appState.totalSize,
            )}:
          </p>
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

  <Footer />
</div>
