<script lang="ts">
  import { formatBytes } from "../lib/format";

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
</script>

<li>
  <div class="card-info">
    {#if !imageError}
      <!-- Steam header image, fallback to folder icon on error -->
      <img 
        src={imageUrl} 
        alt={name} 
        loading="lazy"
        decoding="async"
        onerror={() => imageError = true} 
      />
    {:else}
      <span class="folder-icon">📁</span>
    {/if}

    <div class="details">
      <div class="size-title">
        <span class="size-label">{formatBytes(size)}</span>
        <span class="game-name">{name || "Unknown Game"}</span>
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
    background-color: var(--color-surface-alt);
    padding: 10px 12px;
    margin-bottom: 6px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 1px solid var(--color-border);
    border-radius: 6px;
  }

  li:hover {
    background-color: var(--color-surface-hover);
  }

  .card-info {
    display: flex;
    align-items: center;
    gap: 12px;
    min-height: 46px;
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

  img {
    background-color: var(--color-bg-inset);
  }

  .folder-icon {
    font-size: 1.2rem;
    background-color: var(--color-bg-inset);
  }

  .game-name {
    font-weight: 500;
    color: var(--color-text);
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
    color: var(--color-text-muted);
  }

  .size-label {
    border: 1px solid var(--color-border-strong);
    padding: 0px 6px;
    border-radius: 2px;
    font-size: 0.70rem;
    font-family: sans-serif;
    color: var(--color-text);
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
    background-color: var(--color-border);
    color: var(--color-text);
    border: 1px solid var(--color-border-subtle);
    border-radius: 4px;
    padding: 0 16px;
    min-height: 44px;
    min-width: 150px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 600;
    transition: background 0.2s;
    flex-shrink: 0;
  }

  button:hover:not(:disabled) {
    background-color: var(--color-border-strong);
    color: #ffffff;
    border-color: var(--color-accent);
    outline: 2px solid var(--color-accent);
  }

  button:disabled {
    background-color: var(--color-surface-hover);
    color: var(--color-text-faint);
    cursor: not-allowed;
  }
</style>