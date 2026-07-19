import { ScanOrphans, DeleteOrphan, RequestConfirmation } from '../wailsjs/go/main/App';
import { steam } from '../wailsjs/go/models';

// Manage the global state and actions of the application
export function createAppState() {
  let orphans = $state<steam.AppInfo[]>([]);
  let errorMessage = $state<string>('');
  let isLoading = $state<boolean>(true);
  let isSteamFound = $state<boolean>(true);

  // Scan or re-scan the storage for orphan folders
  async function refreshScan() {
    try {
      isLoading = true;
      errorMessage = '';
      
      const result = await ScanOrphans();
      // Secure assignment: fallback to an empty array if Wails/Go returns null
      orphans = result || []
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
    }
  }

  // Handle the deletion of a specific shadercache folder
  async function handleDelete(appID: string) {
    const confirmed = await RequestConfirmation(
      "Confirm Deletion",
      `Are you sure you want to delete the shadercache for AppID ${appID}?`
    );
    
    if (!confirmed) return;

    try {
      isLoading = true;
      errorMessage = '';
      await DeleteOrphan(appID);
      orphans = orphans.filter(a => a.appID !== appID);
    } catch (err) {
      errorMessage = `Deletion failed: ${err instanceof Error ? err.message : String(err)}`;
    } finally {
      isLoading = false;
    }
  }

  // Handle the deletion of all detected orphan folders
  async function handleDeleteAll() {
    const confirmed = await RequestConfirmation(
      "Confirm Mass Deletion",
      `Are you sure you want to delete ALL ${orphans.length} orphan shadercaches? This action cannot be undone.`
    );
    
    if (!confirmed) return;

    try {
      isLoading = true;
      errorMessage = '';

      for (const id of orphans) {
        await DeleteOrphan(id.appID);
      }

      orphans = [];
    } catch (err) {
      errorMessage = `Mass deletion failed: ${err instanceof Error ? err.message : String(err)}`;
      await refreshScan();
    } finally {
      isLoading = false;
    }
  }

  // Execute an initial scan on state creation
  refreshScan();

  // Expose states and methods as read-only properties or direct functions
  return {
    get orphans() { return orphans; },
    get errorMessage() { return errorMessage; },
    get isLoading() { return isLoading; },
    get isSteamFound() { return isSteamFound; },
    refreshScan,
    handleDelete,
    handleDeleteAll
  };
}