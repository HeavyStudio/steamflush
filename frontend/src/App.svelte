<script lang="ts">
  // Import global styles
  import "./App.css";

  // Import components
  import Header from "./components/Header.svelte";
  import StatusMessage from "./components/StatusMessage.svelte";
  import OrphanCard from "./components/OrphanCard.svelte";
  import Footer from "./components/Footer.svelte";
  import HistoryModal from "./components/HistoryModal.svelte";
  import NotificationModal from "./components/NotificationModal.svelte";
  import { formatBytes } from "./lib/format.ts";

  // Import spatial navigation action
  import { spatialNavigation } from "./lib/spatialNavigation.ts";

  // Import and initialize state logic
  import { createAppState } from "./App.svelte.ts";

  const appState = createAppState();
  let showHistoryModal = $state(false);

  $effect(() => {
    appState.refreshScan();
    appState.refreshHistory();
  });

  function toggleHistory() {
    showHistoryModal = !showHistoryModal;
  }

  function closeNotification() {
    appState.clearSuccessMessage?.();
  }
</script>

<div class="app-container">
  <Header />

  <main
    use:spatialNavigation={{
      enabled:
        !appState.isLoading && !showHistoryModal && !appState.successMessage,
    }}
  >
    <section class="content">
      <StatusMessage
        isLoading={appState.isLoading}
        errorMessage={appState.errorMessage}
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

    <div class="history-trigger">
      <button type="button" class="history-btn" onclick={toggleHistory}>
        Show History
      </button>
    </div>

    {#if showHistoryModal}
      <HistoryModal
        history={appState.history}
        onClose={() => (showHistoryModal = false)}
      />
    {/if}

    {#if appState.successMessage}
      <NotificationModal
        message={appState.successMessage}
        onClose={closeNotification}
        autoCloseDelay={15}
      />
    {/if}
  </main>

  <Footer />
</div>
