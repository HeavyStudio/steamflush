<script lang="ts">
    import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
    import { APP_VERSION } from "../lib/version";
    import GitHubLogo from "../assets/images/GitHub_Invertocat_White.svg";

    function openGitHub() {
        BrowserOpenURL("https://github.com/HeavyStudio/steamflush");
    }
</script>

<footer class="app-footer">
    <!-- Ajout des mentions légales discrètes à gauche -->
    <div class="footer-left">
        <span>© 2026 Jonathan Alfonso</span>
    </div>

    <!-- Tes éléments actuels à droite -->
    <div class="footer-right">
        <span class="version">v{APP_VERSION}</span>
        <button
            type="button"
            class="github-btn"
            onclick={openGitHub}
            aria-label="GitHub Repository"
        >
            <img src={GitHubLogo} alt="GitHub" width="20" height="20" />
        </button>
    </div>
</footer>

<style>
    .app-footer {
        display: flex;
        justify-content: space-between; /* Pour séparer gauche et droite */
        align-items: center;
        padding: 14px 24px;
        font-size: 0.85rem; /* Légèrement plus petit pour ne pas distraire */
        color: var(--color-text-muted);
        border-top: 1px solid var(--color-border-subtle);
        margin-top: auto;
    }

    .footer-right {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .github-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 40px;
        height: 40px;
        padding: 0;
        box-sizing: border-box;
        background: transparent;
        border: 2px solid transparent;
        border-radius: 6px;
        opacity: 0.8;
        cursor: pointer;
        outline: none !important;
        transition: opacity 0.15s ease, border-color 0.15s ease, box-shadow 0.15s ease;
    }

    .github-btn:hover,
    .github-btn:focus,
    .github-btn:focus-visible {
        opacity: 1;
        width: 40px;
        height: 40px;
        border-color: var(--color-accent);
        border-radius: 6px;
        outline: none !important;
        box-shadow: inset 0 0 0 1px var(--color-accent);
    }
</style>
