<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import { spatialNavigation } from "../lib/spatialNavigation";

    let {
        message,
        onClose,
        autoCloseDelay = 15,
    }: {
        message: string;
        onClose: () => void;
        autoCloseDelay?: number;
    } = $props();

    // Initialisation d'un état local modifiable pour le compte à rebours
    let timeLeft = $state(0);
    let okButton = $state<HTMLButtonElement | null>(null);
    let timer: ReturnType<typeof setInterval>;

    onMount(() => {
        // Assignation explicite au montage pour le linter
        timeLeft = autoCloseDelay;

        // Focus automatique sur le bouton OK pour la navigation manette/clavier
        if (okButton) {
            okButton.focus();
        }

        // Timer de fermeture automatique
        timer = setInterval(() => {
            timeLeft -= 1;
            if (timeLeft <= 0) {
                clearInterval(timer);
                onClose();
            }
        }, 1000);
    });

    onDestroy(() => {
        if (timer) clearInterval(timer);
    });

    function handleKeyDown(e: KeyboardEvent) {
        if (e.key === "Escape" || e.key === "Enter" || e.key === " ") {
            onClose();
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) {
            onClose();
        }
    }
</script>

<svelte:window onkeydown={handleKeyDown} />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="modal-wrapper" onclick={handleBackdropClick}>
    <div class="modal-backdrop"></div>

    <div
        class="modal-container"
        use:spatialNavigation={{ enabled: true }}
        role="dialog"
        aria-modal="true"
    >
        <div class="modal-body">
            <div class="icon-success">✓</div>
            <p class="message">{message}</p>
            <span class="timer-text">Closing automatically in {timeLeft}s</span>
        </div>

        <div class="modal-footer">
            <button
                bind:this={okButton}
                type="button"
                class="btn-ok"
                onclick={onClose}
            >
                OK ({timeLeft}s)
            </button>
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
        max-width: 420px;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        align-items: center;
        box-shadow: 0 16px 32px var(--shadow-modal);
        overflow: hidden;
        color: var(--color-text-secondary);
        padding: 1.5rem;
        text-align: center;
        pointer-events: auto;
    }

    .modal-body {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        width: 100%;
        text-align: center;
    }

    .icon-success {
        width: 48px;
        height: 48px;
        background-color: rgba(158, 206, 106, 0.15);
        color: var(--color-success);
        font-size: 1.5rem;
        font-weight: bold;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 1rem;
        border: 1px solid var(--color-success);
    }

    .message {
        font-size: 1.05rem;
        color: var(--color-text);
        margin: 0 0 0.5rem 0;
        font-weight: 500;
    }

    .timer-text {
        font-size: 0.8rem;
        color: var(--color-text-faint);
        font-family: "JetBrains Mono", monospace;
    }

    .modal-footer {
        margin-top: 1.25rem;
        width: 100%;
        display: flex;
        justify-content: center;
    }

    .btn-ok {
        background-color: var(--color-accent);
        color: var(--color-text-on-accent);
        border: none;
        padding: 0.6rem 2rem;
        border-radius: 6px;
        font-weight: 700;
        cursor: pointer;
        transition: all 0.2s;
        min-width: 120px;
    }

    .btn-ok:hover,
    .btn-ok:focus,
    .btn-ok:focus-visible {
        background-color: var(--color-purple);
        outline: 2px solid var(--color-cyan);
        outline-offset: 2px;
    }
</style>
