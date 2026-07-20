import { ScanOrphans, DeleteOrphan, RequestConfirmation, RemoveShaderCacheBatch } from '../wailsjs/go/main/App';
import { steam } from '../wailsjs/go/models';
import { formatBytes } from './lib/format';

export type SortOption = 'name' | 'size' | 'id';

// Manage the global state and actions of the application
export function createAppState() {
  let orphans = $state<steam.AppInfo[]>([]);
  let errorMessage = $state<string>('');
  let successMessage = $state<string>('');
  let isLoading = $state<boolean>(true);
  let isSteamFound = $state<boolean>(true);
  let sortOption = $state<SortOption>('name');
  let isScanning = false;
  let successTimer: ReturnType<typeof setTimeout> | null = null;

  const totalSize = $derived(orphans.reduce((sum, app) => sum + (app.size || 0), 0));

  const sortedOrphans = $derived([...orphans].sort((a, b) => {
    switch(sortOption) {
      case 'name': return a.name.localeCompare(b.name);
      case 'size': return b.size - a.size;
      case 'id': return a.appID.localeCompare(b.appID);
      default: return 0;
    }
  }))

  function setSort(option: SortOption) {
    sortOption = option;
  }

  function showSuccess(msg: string) {
    if (successTimer) clearTimeout(successTimer);
    successMessage = msg;
    successTimer = setTimeout(() => {
      successMessage = '';
      successTimer = null;
    }, 3000);
  }

  // Scan or re-scan the storage for orphan folders
  async function refreshScan() {
    if (isScanning) return;
    isScanning = true;
    isLoading = true;
    errorMessage = '';
    successMessage = '';

    try {
      const result = await ScanOrphans();
      orphans = result || [];
      isSteamFound = true;
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      // Check if the error is specifically about the missing directory
      if (msg.includes("shadercache directory not found")) {
        isSteamFound = false;
      } else {
        errorMessage = msg;
      }
      orphans = [];
    } finally {
      isLoading = false;
      isScanning = false;
    }
  }

  async function handleDelete(appID: string) {
    const app = orphans.find(a => a.appID === appID);
    const sizeStr = app ? formatBytes(app.size) : "unknown size";

    const confirmed = await RequestConfirmation(
      "Confirm Deletion",
      `Are you sure you want to delete the shadercache for AppID ${appID} (${sizeStr})?`
    );
    
    if (!confirmed) return;

    try {
      isLoading = true;
      errorMessage = '';
      successMessage = '';
      await DeleteOrphan(appID);
      orphans = orphans.filter(a => a.appID !== appID);
      showSuccess("Shadercache deleted successfully!");
    } catch (err) {
      errorMessage = `Deletion failed: ${err instanceof Error ? err.message : String(err)}`;
    } finally {
      isLoading = false;
    }
  }

  async function handleDeleteAll() {
    const confirmed = await RequestConfirmation(
      "Confirm Mass Deletion",
      `Are you sure you want to delete ALL ${orphans.length} orphan shadercaches? This will free up ${formatBytes(totalSize)}. This action cannot be undone.`
    );
    
    if (!confirmed) return;

    try {
      isLoading = true;
      errorMessage = '';
      successMessage = '';

      const idsToDelete = orphans.map(a => a.appID);
      const result = await RemoveShaderCacheBatch(idsToDelete);
      const failedIDs = result ?? [];

      if (failedIDs.length > 0) {
        orphans = orphans.filter(a => failedIDs.includes(a.appID));
        errorMessage = `Mass deletion completed with ${failedIDs.length} errors`;
      } else {
        orphans = [];
        showSuccess("All orphan directories cleared!");
      }
    } catch (err) {
      errorMessage = `Critical failure during mass deletion: ${err instanceof Error ? err.message : String(err)}`;
      await refreshScan();
    } finally {
      isLoading = false;
    }
  }

  // Expose states and methods as read-only properties or direct functions
  return {
    get orphans() { return sortedOrphans; },
    get totalSize() { return totalSize; },
    get sortOption() { return sortOption; },
    get errorMessage() { return errorMessage; },
    get successMessage() { return successMessage; },
    get isLoading() { return isLoading; },
    get isSteamFound() { return isSteamFound; },
    refreshScan,
    handleDelete,
    handleDeleteAll,
    setSort
  };
}