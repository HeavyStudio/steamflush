<script lang="ts">
    import { tick } from "svelte";
    import type { steam } from "../../wailsjs/go/models";
    import { formatBytes, formatDate } from "../lib/format";
    import { spatialNavigation } from "../lib/spatialNavigation";

    let {
        history = [],
        onClose,
    }: {
        history: steam.CleanRecord[];
        onClose: () => void;
    } = $props();

    let expandedIndex = $state<number | null>(null);
    let modalContainer = $state<HTMLElement | null>(null);

    let totalBytesAllTime = $derived(
        history.reduce((acc, curr) => acc + curr.bytes_freed, 0),
    );

    function toggleExpand(index: number) {
        expandedIndex = expandedIndex === index ? null : index;
    }

    function handleKeyDown(e: KeyboardEvent) {
        if (e.key === "Escape") {
            onClose();
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) {
            onClose();
        }
    }

    function handleMouseEnter(e: MouseEvent) {
        (e.currentTarget as HTMLElement)?.focus();
    }

    // Effet declenché dès que la modale ou history est prêt dans le DOM
    $effect(() => {
        if (modalContainer) {
            tick().then(() => {
                // Cherche d'abord une carte d'historique
                const firstCard =
                    modalContainer?.querySelector<HTMLElement>(".history-card");
                if (firstCard) {
                    firstCard.focus();
                } else {
                    // S'il n'y a pas de cartes (historique vide), on focus le bouton close principal
                    const closeBtn =
                        modalContainer?.querySelector<HTMLElement>(
                            ".btn-close",
                        );
                    closeBtn?.focus();
                }
            });
        }
    });
</script>

<svelte:window onkeydown={handleKeyDown} />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="modal-wrapper" onclick={handleBackdropClick}>
    <div class="modal-backdrop"></div>

    <div
        bind:this={modalContainer}
        class="modal-container"
        use:spatialNavigation={{ enabled: true }}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
    >
        <div class="modal-header">
            <div>
                <h2 id="modal-title">Clean History</h2>
                <p class="history-subtitle">
                    Total recovered space: <span class="highlight"
                        >{formatBytes(totalBytesAllTime)}</span
                    >
                </p>
            </div>
            <button
                type="button"
                class="close-btn"
                onclick={onClose}
                aria-label="Close modal">✕</button
            >
        </div>

        <div class="modal-body">
            {#if history.length === 0}
                <div class="state-text">No cleaning sessions recorded yet.</div>
            {:else}
                {#each history as record, index}
                    <div
                        class="history-card"
                        role="button"
                        tabindex="0"
                        onmouseenter={handleMouseEnter}
                        onclick={() => toggleExpand(index)}
                        onkeydown={(e) => (e.key == 'Enter' || e.key == ' ') && toggleExpand(index)}
                    >
                        <div class="card-summary">
                            <div>
                                <span class="date"
                                    >{formatDate(record.timestamp)}</span
                                >
                                <div class="title">
                                    {record.items_count} orphan{record.items_count >
                                    1
                                        ? "s"
                                        : ""} cleaned
                                </div>
                            </div>
                            <div class="right">
                                <span class="size"
                                    >{formatBytes(record.bytes_freed)}</span
                                >
                                <div class="toggle-text">
                                    {expandedIndex === index
                                        ? "▲ Hide details"
                                        : "▼ View details"}
                                </div>
                            </div>
                        </div>

                        {#if expandedIndex === index}
                            <div class="card-details">
                                {#each record.items as item}
                                    <div class="detail-row">
                                        <span class="item-name">
                                            {item.name}
                                            <small class="app-id"
                                                >({item.app_id})</small
                                            >
                                        </span>
                                        <span class="item-size"
                                            >{formatBytes(item.size)}</span
                                        >
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                {/each}
            {/if}
        </div>

        <div class="modal-footer">
            <button type="button" class="btn-close" onclick={onClose}
                >Close</button
            >
        </div>
    </div>
</div>

<style>
    .modal-wrapper {
        position: fixed;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1rem;
        pointer-events: auto;
    }

    .modal-backdrop {
        position: absolute;
        inset: 0;
        background: var(--scrim);
        backdrop-filter: blur(6px);
        pointer-events: none;
    }

    .modal-container {
        position: relative;
        z-index: 1001;
        background-color: var(--color-surface);
        border: 1px solid var(--color-border);
        width: 100%;
        max-width: 600px;
        max-height: 85vh;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        box-shadow: 0 16px 32px var(--shadow-modal);
        overflow: hidden;
        color: var(--color-text-secondary);
        pointer-events: auto;
    }

    .modal-header {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--color-border-subtle);
        background-color: var(--color-bg);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .modal-header h2 {
        font-family: "Rajdhani", sans-serif;
        color: var(--color-accent);
        font-size: 1.4rem;
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.1em;
    }

    .history-subtitle {
        font-size: 0.875rem;
        color: var(--color-text-muted);
        margin: 0.25rem 0 0 0;
    }

    .highlight {
        font-family: "JetBrains Mono", monospace;
        color: var(--color-cyan);
        font-weight: 700;
    }

    .close-btn {
        background: transparent;
        border: none;
        color: var(--color-text-muted);
        font-size: 1.2rem;
        cursor: pointer;
        padding: 0.25rem 0.5rem;
        min-height: auto;
        border-radius: 4px;
        transition:
            color 0.2s,
            background-color 0.2s;
    }

    .close-btn:hover,
    .close-btn:focus,
    .close-btn:focus-visible {
        color: var(--color-text);
        background-color: var(--color-surface-hover);
        outline: 2px solid var(--color-accent);
    }

    .modal-body {
        padding: 1.25rem;
        overflow-y: auto;
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .state-text {
        text-align: center;
        color: var(--color-text-muted);
        padding: 2.5rem 0;
    }

    .history-card {
        /* 1. Dimensions et boîte strictes */
        box-sizing: border-box !important;
        width: 100%;
        margin: 0;
        padding: 0.875rem 1rem;
        display: flex;
        flex-direction: column;

        /* 2. Style visuel */
        background-color: var(--color-surface-alt);
        border: 2px solid var(--color-border);
        border-radius: 6px;

        /* 3. Bloque tout recalcul de ligne / texte natif aux boutons */
        outline: none !important;
        line-height: 1.2;
        -webkit-font-smoothing: antialiased;

        color: inherit;
        font-family: inherit;
        text-align: left;
        cursor: pointer;

        /* On ne tienne que les couleurs pour éviter tout saut de layout */
        transition:
            background-color 0.15s ease,
            border-color 0.15s ease,
            box-shadow 0.15s ease;
    }

    .history-card:hover,
    .history-card:focus,
    .history-card:focus-visible {
        background-color: var(--color-surface-hover);
        border-color: var(--color-accent);
        outline: none !important;
        box-shadow:
            0 0 0 2px var(--color-accent),
            0 0 12px var(--focus-glow);
    }

    .card-summary {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        /* Verrouille l'alignement vertical des deux blocs d'information */
        min-height: 2.4rem;
    }

    .date {
        font-family: "JetBrains Mono", monospace;
        font-size: 0.75rem;
        color: var(--color-text-faint);
    }

    .title {
        font-size: 0.95rem;
        font-weight: 600;
        color: var(--color-text);
        margin-top: 0.2rem;
    }

    .right {
        text-align: right;
    }

    .size {
        font-family: "JetBrains Mono", monospace;
        font-size: 0.95rem;
        font-weight: 700;
        color: var(--color-cyan);
    }

    .toggle-text {
        font-size: 0.75rem;
        color: var(--color-text-muted);
        margin-top: 0.2rem;
    }

    .card-details {
        margin-top: 0.75rem;
        padding-top: 0.75rem;
        border-top: 1px solid var(--color-border-subtle);
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
    }

    .detail-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 0.85rem;
        background-color: var(--color-bg);
        padding: 0.4rem 0.6rem;
        border-radius: 4px;
        color: var(--color-text-secondary);
    }

    .item-name {
        font-weight: 500;
    }

    .app-id {
        font-family: "JetBrains Mono", monospace;
        color: var(--color-text-faint);
        font-size: 0.75rem;
    }

    .item-size {
        font-family: "JetBrains Mono", monospace;
        color: var(--color-cyan);
        font-size: 0.8rem;
    }

    .modal-footer {
        padding: 1rem 1.25rem;
        border-top: 1px solid var(--color-border-subtle);
        background-color: var(--color-bg);
        display: flex;
        justify-content: flex-end;
    }

    .btn-close {
        background-color: var(--color-border);
        color: var(--color-text);
    }

    .btn-close:hover:not(:disabled),
    .btn-close:focus,
    .btn-close:focus-visible {
        background-color: var(--color-border-strong);
        color: #ffffff;
        outline: 2px solid var(--color-accent);
    }
</style>
