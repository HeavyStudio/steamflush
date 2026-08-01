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
        max-width: 420px;
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        align-items: center;
        box-shadow: 0 16px 32px rgba(0, 0, 0, 0.5);
        overflow: hidden;
        color: #a9b1d6;
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
        color: #9ece6a;
        font-size: 1.5rem;
        font-weight: bold;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 1rem;
        border: 1px solid #9ece6a;
    }

    .message {
        font-size: 1.05rem;
        color: #c0caf5;
        margin: 0 0 0.5rem 0;
        font-weight: 500;
    }

    .timer-text {
        font-size: 0.8rem;
        color: #565f89;
        font-family: "JetBrains Mono", monospace;
    }

    .modal-footer {
        margin-top: 1.25rem;
        width: 100%;
        display: flex;
        justify-content: center;
    }

    .btn-ok {
        background-color: #7aa2f7;
        color: #15161e;
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
        background-color: #bb9af7;
        outline: 2px solid #7dcfff;
        outline-offset: 2px;
    }
</style>
