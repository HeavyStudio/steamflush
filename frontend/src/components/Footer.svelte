<script lang="ts">
    import { APP_VERSION } from "../lib/version";
    import GitHubLogo from "../assets/images/GitHub_Invertocat_White.svg";
</script>

<footer class="app-footer">
    <!-- Ajout des mentions légales discrètes à gauche -->
    <div class="footer-left">
        <span>© 2026 Jonathan Alfonso | Heavy Studio</span>
    </div>

    <!-- Tes éléments actuels à droite -->
    <div class="footer-right">
        <span class="version">v{APP_VERSION}</span>
        <a
            href="https://github.com/HeavyStudio/steamflush"
            target="_blank"
            rel="noopener"
        >
            <img src={GitHubLogo} alt="GitHub" width="20" height="20" />
        </a>
    </div>
</footer>

<style>
    .app-footer {
        display: flex;
        justify-content: space-between; /* Pour séparer gauche et droite */
        align-items: center;
        padding: 14px 24px;
        font-size: 0.85rem; /* Légèrement plus petit pour ne pas distraire */
        color: #7885b5;
        border-top: 1px solid #2f334d;
        margin-top: auto;
    }

    .footer-right {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    a {
        display: flex;
        align-items: center;
    }
</style>
