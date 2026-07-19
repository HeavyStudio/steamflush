<script lang="ts">
  interface Props {
    name: string;
    id: string;
    size: number;
    disabled: boolean;
    onDelete: (id: string) => void;
  }

  let { name, id, size, disabled, onDelete }: Props = $props();
  
  // Track image load failure to trigger fallback icon
  let imageError = $state(false);
  let imageUrl = $derived(`https://cdn.akamai.steamstatic.com/steam/apps/${id}/header.jpg`);

  // Get directory's size
  function formatBytes(bytes: number) {
    if (bytes == 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
</script>

<li>
  <div class="card-info">
    {#if !imageError}
      <!-- Steam header image, fallback to folder icon on error -->
      <img 
        src={imageUrl} 
        alt={name} 
        onerror={() => imageError = true} 
      />
    {:else}
      <span class="folder-icon">📁</span>
    {/if}

    <div class="details">
      <div class="size-title">
        <span class="size-label">{formatBytes(size)}</span>
        <span class="game-name">{name || "Unknown Application"}</span>
      </div>
      <span class="appid">AppID: {id}</span>
    </div>
  </div>
  <button onclick={() => onDelete(id)} {disabled}>
    Delete
  </button>
</li>

<style>
  li {
    background-color: #1b2838;
    padding: 10px;
    margin-bottom: 4px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 1px solid #2d4059;
  }

  li:hover {
    background-color: #26384f;
  }

  .card-info {
    display: flex;
    align-items: center;
    gap: 12px;
    height: 45px;
    flex-shrink: 0;
  }

  .details {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    height: 100%;
    min-width: 0;
    flex: 1;
  }

  img, .folder-icon {
    width: 100px;
    height: 46px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 2px;
    object-fit: contain;
    flex-shrink: 0;
  }

  .folder-icon { 
    font-size: 1.2rem;
    background-color: #0f1922;
  }

  .game-name {
    font-weight: 500;
    color: #e4e4e4;
    font-size: 0.95rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    width: 100%;
    text-align: left;
  }

  .appid {
    font-family: sans-serif;
    font-size: 0.75rem;
    color: #6691b0;
  }

  .size-label {
    border: 1px solid #6691b0;
    padding: 0px 6px;
    border-radius: 2px;
    font-size: 0.70rem;
    font-family: sans-serif;
    color: #c6d4df;
    margin-right: 8px;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .size-title {
    display: flex;
    flex-direction: row;
    align-items: center;
    width: fit-content;
    min-width: 0;
  }

  button {
    background-color: #3d4450;
    color: #c6d4df;
    border: 1px solid #2d363e;
    padding: 4px 12px;
    cursor: pointer;
    font-size: 0.85rem;
    transition: background 0.2s;
    flex-shrink: 0;
  }

  button:hover:not(:disabled) {
    background-color: #4e5968;
    color: #ffffff;
  }

  button:disabled {
    background-color: #2a3138;
    color: #4a5a6a;
    cursor: not-allowed;
  }
</style>