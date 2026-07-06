<script>
  import { onMount } from 'svelte';
  import { LoadConfig, SaveConfig, GetTokenExpiration } from '../../wailsjs/go/main/App.js';

  /** @type {(() => void) | null} */
  export let onSaveSuccess = null;
  /** @type {boolean} */
  export let isLoginMode = false;
  /** @type {(() => void) | null} */
  export let onLogout = null;

  let token = '';
  let username = '';
  let client = '';
  let rate = '';
  let cnpj = '';
  let companyName = '';

  let statusMessage = '';
  let statusType = ''; // 'success' or 'error'

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined && (/** @type {any} */(window)).go.main !== undefined;

  onMount(async () => {
    if (!isWailsAvailable) {
      return;
    }
    try {
      const config = await LoadConfig();
      token = config.token || '';
      username = config.username || '';
      client = config.client || '';
      rate = config.rate || '';
      cnpj = config.cnpj || '';
      companyName = config.company_name || '';
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      showStatus('Falha ao carregar as configurações: ' + message, 'error');
    }
  });

  async function handleSave() {
    if (!isWailsAvailable) {
      showStatus('Erro: O backend Go não está acessível a partir do navegador convencional.', 'error');
      return;
    }

    if (!token || !username) {
      showStatus('Por favor, preencha o Token do GitHub e o Usuário.', 'error');
      return;
    }

    try {
      // Validate token and retrieve expiration from headers
      const expiration = await GetTokenExpiration(token);

      await SaveConfig({
        token,
        username,
        client,
        rate,
        token_expiration: expiration || '',
        cnpj,
        company_name: companyName
      });
      showStatus('Configurações salvas com sucesso!', 'success');
      if (onSaveSuccess) {
        onSaveSuccess();
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      showStatus('Erro ao salvar ou validar token: ' + message, 'error');
    }
  }

  async function handleDisconnect() {
    if (!isWailsAvailable) {
      showStatus('Erro: O backend Go não está acessível a partir do navegador convencional.', 'error');
      return;
    }

    try {
      await SaveConfig({
        token: '',
        username: '',
        client: '',
        rate: '',
        token_expiration: '',
        cnpj: '',
        company_name: ''
      });
      token = '';
      username = '';
      client = '';
      rate = '';
      cnpj = '';
      companyName = '';
      showStatus('Desconectado com sucesso!', 'success');
      if (onLogout) {
        onLogout();
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      showStatus('Erro ao desconectar: ' + message, 'error');
    }
  }

  /**
   * @param {string} msg
   * @param {string} type
   */
  function showStatus(msg, type) {
    statusMessage = msg;
    statusType = type;
    setTimeout(() => {
      statusMessage = '';
    }, 4000);
  }
</script>

<div class={isLoginMode ? '' : 'card'}>
  {#if !isLoginMode}
    <h2>Configurações do GitHub e Relatório</h2>
    <p style="color: var(--text-muted); margin-bottom: 1.5rem; font-size: 0.9rem;">
      Insira suas credenciais do GitHub para permitir a busca de commits e pull requests.
      Suas credenciais são salvas localmente e com segurança no seu dispositivo.
    </p>
  {/if}

  {#if statusMessage}
    <div class="status-msg status-{statusType}">
      {#if statusType === 'success'}
        <span>✓</span>
      {:else}
        <span>⚠</span>
      {/if}
      <span>{statusMessage}</span>
    </div>
  {/if}

  <form onsubmit={(e) => { e.preventDefault(); handleSave(); }}>
    <div class={isLoginMode ? 'grid-1' : 'grid-2'}>
      <div class="form-group">
        <label for="token">Personal Access Token (PAT) do GitHub</label>
        <input
          id="token"
          type="password"
          class="form-control"
          placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
          bind:value={token}
          autocomplete="off"
        />
        <small style="color: var(--text-muted); font-size: 0.75rem; margin-top: 0.25rem;">
          Requer permissão de leitura de repositórios (<code>repo</code>).
        </small>
      </div>

      <div class="form-group">
        <label for="username">Nome de Usuário do GitHub</label>
        <input
          id="username"
          type="text"
          class="form-control"
          placeholder="seu_usuario_github"
          bind:value={username}
        />
      </div>
    </div>

    {#if !isLoginMode}
      <h3 style="margin-top: 1.5rem; margin-bottom: 0.75rem; font-size: 1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem; color: var(--text-main);">Seus Dados Fiscais / Empresa</h3>
      <div class="grid-2">
        <div class="form-group">
          <label for="companyName">Sua Razão Social / Nome Completo</label>
          <input
            id="companyName"
            type="text"
            class="form-control"
            placeholder="Ex: EDUARDO MONTOVANELLI DALMASO"
            bind:value={companyName}
          />
        </div>

        <div class="form-group">
          <label for="cnpj">Seu CNPJ</label>
          <input
            id="cnpj"
            type="text"
            class="form-control"
            placeholder="Ex: 52.340.739/0001-46"
            bind:value={cnpj}
          />
        </div>
      </div>

      <h3 style="margin-top: 1.5rem; margin-bottom: 0.75rem; font-size: 1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem; color: var(--text-main);">Dados do Cliente Padrão</h3>
      <div class="grid-2">
        <div class="form-group">
          <label for="client">Nome da Empresa / Cliente Padrão</label>
          <input
            id="client"
            type="text"
            class="form-control"
            placeholder="Empresa Exemplo Ltda"
            bind:value={client}
          />
        </div>

        <div class="form-group">
          <label for="rate">Valor por Hora Cobrado (R$ ou $/h)</label>
          <input
            id="rate"
            type="text"
            class="form-control"
            placeholder="150"
            bind:value={rate}
          />
        </div>
      </div>
    {/if}

    <div style="margin-top: 1.5rem; display: flex; justify-content: space-between; align-items: center; gap: 1rem; width: 100%;">
      {#if !isLoginMode}
        <button type="button" class="btn btn-danger" onclick={handleDisconnect}>
          Desconectar Conta
        </button>
      {/if}
      <div style="margin-left: auto;">
        <button type="submit" class="btn btn-primary" style={isLoginMode ? 'width: 100%;' : ''}>
          {isLoginMode ? 'Conectar Conta' : 'Salvar Configurações'}
        </button>
      </div>
    </div>
  </form>
</div>

<style>
  .grid-1 {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .btn-danger {
    background-color: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.2);
  }

  .btn-danger:hover {
    background-color: rgba(239, 68, 68, 0.2);
    border-color: #ef4444;
  }
</style>


