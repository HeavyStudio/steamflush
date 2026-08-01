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
                    <button
                        type="button"
                        class="history-card"
                        tabindex="0"
                        onclick={() => toggleExpand(index)}
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
                    </button>
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
        background: rgba(15, 16, 25, 0.8);
        backdrop-filter: blur(6px);
        pointer-events: none;
    }

    .modal-container {
        position: relative;
        z-index: 1001;
        background-color: #1f2335;
        border: 1px solid #414868;
        width: 100%;
        max-width: 600px;
        max-height: 85vh;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        box-shadow: 0 16px 32px rgba(0, 0, 0, 0.5);
        overflow: hidden;
        color: #a9b1d6;
        pointer-events: auto;
    }

    .modal-header {
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid #292e42;
        background-color: #1a1b26;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .modal-header h2 {
        font-family: "Rajdhani", sans-serif;
        color: #7aa2f7;
        font-size: 1.4rem;
        margin: 0;
        text-transform: uppercase;
        letter-spacing: 0.1em;
    }

    .history-subtitle {
        font-size: 0.875rem;
        color: #787c99;
        margin: 0.25rem 0 0 0;
    }

    .highlight {
        font-family: "JetBrains Mono", monospace;
        color: #7dcfff;
        font-weight: 700;
    }

    .close-btn {
        background: transparent;
        border: none;
        color: #787c99;
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
        color: #c0caf5;
        background-color: #292e42;
        outline: 2px solid #7aa2f7;
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
        color: #787c99;
        padding: 2.5rem 0;
    }

    .history-card {
        width: 100%;
        background-color: #24283b;
        border: 1px solid #414868;
        border-radius: 6px;
        padding: 0.875rem 1rem;
        text-align: left;
        cursor: pointer;
        transition:
            border-color 0.2s,
            background-color 0.2s;
        color: inherit;
        font-family: inherit;
    }

    .history-card:hover,
    .history-card:focus,
    .history-card:focus-visible {
        border-color: #7aa2f7;
        background-color: #292e42;
        outline: 2px solid #7aa2f7;
        outline-offset: -2px;
    }

    .card-summary {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .date {
        font-family: "JetBrains Mono", monospace;
        font-size: 0.75rem;
        color: #565f89;
    }

    .title {
        font-size: 0.95rem;
        font-weight: 600;
        color: #c0caf5;
        margin-top: 0.2rem;
    }

    .right {
        text-align: right;
    }

    .size {
        font-family: "JetBrains Mono", monospace;
        font-size: 0.95rem;
        font-weight: 700;
        color: #7dcfff;
    }

    .toggle-text {
        font-size: 0.75rem;
        color: #787c99;
        margin-top: 0.2rem;
    }

    .card-details {
        margin-top: 0.75rem;
        padding-top: 0.75rem;
        border-top: 1px solid #292e42;
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
    }

    .detail-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 0.85rem;
        background-color: #1a1b26;
        padding: 0.4rem 0.6rem;
        border-radius: 4px;
        color: #a9b1d6;
    }

    .item-name {
        font-weight: 500;
    }

    .app-id {
        font-family: "JetBrains Mono", monospace;
        color: #565f89;
        font-size: 0.75rem;
    }

    .item-size {
        font-family: "JetBrains Mono", monospace;
        color: #7dcfff;
        font-size: 0.8rem;
    }

    .modal-footer {
        padding: 1rem 1.25rem;
        border-top: 1px solid #292e42;
        background-color: #1a1b26;
        display: flex;
        justify-content: flex-end;
    }

    .btn-close {
        background-color: #414868;
        color: #c0caf5;
    }

    .btn-close:hover:not(:disabled),
    .btn-close:focus,
    .btn-close:focus-visible {
        background-color: #565f89;
        color: #ffffff;
        outline: 2px solid #7aa2f7;
    }
</style>
