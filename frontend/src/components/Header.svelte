<script lang="ts">
  import testLogoTr from '../assets/images/logo.png';
</script>

<header>
  <div class="logo-group">
    <img src={testLogoTr} alt="SteamFlush Otter Logo TR" class="header-logo" />
    <div class="title-container">
      <h1>STEAMFLUSH</h1>
      <p class="subtitle">CLEAN WHAT STEAM LEAVES BEHIND</p>
    </div>
  </div>
</header>
<hr />

<style>
  header {
    background-color: #1a1b26;
    padding: 1.5rem 1rem;
    display: flex;
    justify-content: center;
    align-items: center;
    border-bottom: 1px solid #414868;
  }

  .logo-group {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .header-logo {
    width: 100px;
    height: 100px;
  }

  .title-container {
    text-align: left;
  }

  h1 {
    font-family: 'Rajdhani', sans-serif;
    color: #ffffff;
    margin: 0;
    font-size: 2.5rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  .subtitle {
    color: #a9b1d6;
    margin: -0.25rem 0 0 0;
    font-size: 0.85rem;
    font-weight: 500;
    letter-spacing: 0.05em;
    opacity: 0.8;
  }
</style>