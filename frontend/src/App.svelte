<script>
  import { onMount } from 'svelte';
  import { LoadConfig } from '../wailsjs/go/main/App.js';
  import ActivityFetch from './components/ActivityFetch.svelte';
  import ReportBuilder from './components/ReportBuilder.svelte';
  import Settings from './components/Settings.svelte';

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined && (/** @type {any} */(window)).go.main !== undefined;

  let currentView = $state('activity');
  /** @type {any[]} */
  let selectedActivities = $state([]);
  let isConfigLoaded = $state(false);
  let isLoggedIn = $state(false);
  /** @type {number | null} */
  let daysRemaining = $state(null);
  let isSidebarCollapsed = $state(false);

  $effect(() => {
    const el = document.getElementById('app');
    if (el) {
      if (isSidebarCollapsed) {
        el.classList.add('sidebar-collapsed');
      } else {
        el.classList.remove('sidebar-collapsed');
      }
    }
  });

  /**
   * @param {string} expStr
   */
  function calculateDaysRemaining(expStr) {
    if (!expStr) return null;
    // Converts "YYYY-MM-DD HH:MM:SS UTC" to ISO Format "YYYY-MM-DDTHH:MM:SSZ"
    const isoStr = expStr.replace(' ', 'T').replace(' UTC', 'Z');
    const expDate = new Date(isoStr);
    const today = new Date();
    const diffTime = expDate.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }

  async function checkCredentials() {
    if (!isWailsAvailable) {
      isLoggedIn = false;
      isConfigLoaded = true;
      return;
    }
    try {
      const config = await LoadConfig();
      if (config && config.token && config.username) {
        isLoggedIn = true;
        currentView = 'activity';
        daysRemaining = calculateDaysRemaining(config.token_expiration);
      } else {
        isLoggedIn = false;
        daysRemaining = null;
      }
    } catch (err) {
      isLoggedIn = false;
      daysRemaining = null;
    } finally {
      isConfigLoaded = true;
    }
  }

  onMount(() => {
    checkCredentials();
  });

  async function handleSaveSuccess() {
    await checkCredentials();
  }

  function handleLogout() {
    isLoggedIn = false;
    daysRemaining = null;
  }

  /**
   * @param {string} view
   */
  function setView(view) {
    if (isLoggedIn) {
      currentView = view;
    }
  }
</script>

{#if !isConfigLoaded}
  <div class="loading-screen">
    <div class="spinner"></div>
    <p>Carregando configurações...</p>
  </div>
{:else if !isLoggedIn}
  {#if !isWailsAvailable}
    <div class="browser-warning-banner">
      ⚠️ <strong>Ambiente Web Detectado</strong>: O backend Go não está disponível no navegador. Utilize a janela do aplicativo desktop Wails criada pelo comando <code>wails dev</code>.
    </div>
  {/if}
  <div class="login-container" class:with-banner={!isWailsAvailable}>
    <div class="login-card">
      <div class="login-header">
        <div class="login-logo-icon">G</div>
        <h2 class="gradient-text">Conectar ao GitHub</h2>
        <p>Para iniciar, insira suas credenciais do GitHub. Elas serão salvas localmente e com segurança no seu dispositivo.</p>
      </div>
      <Settings onSaveSuccess={handleSaveSuccess} isLoginMode={true} />
    </div>
  </div>
{:else}
  {#if !isWailsAvailable}
    <div class="browser-warning-banner">
      ⚠️ <strong>Ambiente Web Detectado</strong>: O backend Go não está disponível no navegador. Utilize a janela do aplicativo desktop Wails criada pelo comando <code>wails dev</code>.
    </div>
  {/if}
  <div class="sidebar">
    <div class="sidebar-brand">
      <div class="sidebar-logo-icon">G</div>
      <span class="gradient-text" style="font-weight: 700;">GitReport</span>
    </div>

    <ul class="nav-links">
      <li class="nav-item">
        <button
          class="nav-btn"
          class:active={currentView === 'activity'}
          onclick={() => setView('activity')}
        >
          <span style="font-size: 1.1rem;">🔍</span> Buscar Atividades
        </button>
      </li>
      <li class="nav-item">
        <button
          class="nav-btn"
          class:active={currentView === 'report'}
          onclick={() => setView('report')}
        >
          <span style="font-size: 1.1rem;">📄</span> Gerador de PDF
        </button>
      </li>
      <li class="nav-item">
        <button
          class="nav-btn"
          class:active={currentView === 'settings'}
          onclick={() => setView('settings')}
        >
          <span style="font-size: 1.1rem;">⚙️</span> Configurações
        </button>
      </li>
    </ul>

    <div class="sidebar-footer">
      <div>GitReport Generator v1.0</div>
      {#if daysRemaining !== null}
        <div class="token-expiry" class:warning={daysRemaining <= 7}>
          {#if daysRemaining > 0}
            🔑 Token: {daysRemaining} {daysRemaining === 1 ? 'dia restante' : 'dias restantes'}
          {:else}
            ⚠️ Token expirado
          {/if}
        </div>
      {/if}
    </div>
  </div>

  <div class="workspace">
    <div class="workspace-header">
      <div style="display: flex; align-items: center; gap: 1rem;">
        <button class="toggle-sidebar-btn" onclick={() => isSidebarCollapsed = !isSidebarCollapsed} title={isSidebarCollapsed ? "Mostrar Menu" : "Esconder Menu"}>
          ☰
        </button>
        <h1 class="workspace-title">
          {#if currentView === 'activity'}
            Busca de Atividades no GitHub
          {:else}
            {#if currentView === 'report'}
              Visualização & Exportação de PDF
            {:else}
              Configurações da Conta
            {/if}
          {/if}
        </h1>
      </div>
      <div class="workspace-meta" style="font-size: 0.85rem; color: var(--text-muted); display: flex; align-items: center; gap: 1rem;">
        {#if selectedActivities.length > 0}
          <span style="background-color: rgba(99, 102, 241, 0.1); padding: 0.4rem 0.8rem; border-radius: 20px; border: 1px solid rgba(99, 102, 241, 0.2);">
            {selectedActivities.length} itens no Relatório
          </span>
        {/if}
      </div>
    </div>

    <div class="workspace-content">
      {#if currentView === 'activity'}
        <ActivityFetch
          bind:selectedActivities={selectedActivities}
          goToReport={() => setView('report')}
        />
      {:else}
        {#if currentView === 'report'}
          <ReportBuilder bind:selectedActivities={selectedActivities} />
        {:else}
          <Settings onLogout={handleLogout} />
        {/if}
      {/if}
    </div>
  </div>
{/if}

<style>
  .loading-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    width: 100vw;
    background-color: #0b0f19;
    color: #c9d1d9;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(88, 166, 255, 0.1);
    border-radius: 50%;
    border-top-color: #58a6ff;
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .browser-warning-banner {
    background-color: #7f1d1d;
    border-bottom: 1px solid #b91c1c;
    color: #fca5a5;
    padding: 0.75rem 1rem;
    text-align: center;
    font-size: 0.9rem;
    font-weight: 500;
    width: 100%;
    box-sizing: border-box;
    position: sticky;
    top: 0;
    z-index: 9999;
  }

  .login-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    width: 100vw;
    background: radial-gradient(circle at center, #0f172a 0%, #020617 100%);
    padding: 1.5rem;
    box-sizing: border-box;
  }

  .login-container.with-banner {
    min-height: calc(100vh - 40px);
  }

  .login-card {
    background-color: #0d1117;
    border: 1px solid #30363d;
    border-radius: 12px;
    padding: 2.5rem;
    max-width: 500px;
    width: 100%;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5), 0 8px 10px -6px rgba(0, 0, 0, 0.5);
  }

  .login-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .login-logo-icon {
    font-size: 2rem;
    font-weight: 800;
    color: #ffffff;
    background: linear-gradient(135deg, #00d2ff 0%, #0066ff 100%);
    width: 50px;
    height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    margin: 0 auto 1rem auto;
    box-shadow: 0 4px 12px rgba(0, 210, 255, 0.3);
  }

  .login-header h2 {
    font-size: 1.75rem;
    margin-bottom: 0.5rem;
    font-weight: 700;
  }

  .login-header p {
    color: #8b949e;
    font-size: 0.9rem;
    line-height: 1.4;
  }

  .token-expiry {
    font-size: 0.75rem;
    color: #8b949e;
    margin-top: 0.5rem;
    background-color: rgba(255, 255, 255, 0.03);
    padding: 0.35rem 0.6rem;
    border-radius: 6px;
    display: inline-block;
    border: 1px solid rgba(255, 255, 255, 0.05);
    font-weight: 500;
  }

  .token-expiry.warning {
    color: #f87171;
    background-color: rgba(239, 68, 68, 0.1);
    border-color: rgba(239, 68, 68, 0.25);
  }

  .toggle-sidebar-btn {
    background: transparent;
    border: 1px solid var(--border-color);
    color: var(--text-main);
    border-radius: var(--radius-md);
    width: 38px;
    height: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 1.1rem;
    transition: all 0.2s ease;
  }

  .toggle-sidebar-btn:hover {
    background-color: var(--bg-hover);
    border-color: var(--text-muted);
  }
</style>


