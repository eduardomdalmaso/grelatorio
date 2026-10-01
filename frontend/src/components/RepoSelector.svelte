<script>
  import { onMount } from 'svelte';
  import { FetchUserRepositories } from '../../wailsjs/go/main/App.js';

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined;

  /** @type {{ selectedRepos?: string[], disabled?: boolean }} */
  let { selectedRepos = $bindable(['all']), disabled = false } = $props();

  /** @type {Array<{id: number, name: string, full_name: string, private: boolean, description?: string}>} */
  let repos = $state([]);
  let loadingRepos = $state(false);
  let isOpen = $state(false);
  let searchTerm = $state('');
  /** @type {HTMLElement | null} */
  let dropdownRef = $state(null);

  onMount(() => {
    if (isWailsAvailable) {
      loadRepos();
    }

    /** @param {MouseEvent} e */
    function handleClickOutside(e) {
      if (dropdownRef && !dropdownRef.contains(/** @type {Node} */ (e.target))) {
        isOpen = false;
      }
    }

    /** @param {KeyboardEvent} e */
    function handleKeyDown(e) {
      if (e.key === 'Escape' && isOpen) {
        isOpen = false;
      }
    }

    document.addEventListener('click', handleClickOutside);
    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('click', handleClickOutside);
      document.removeEventListener('keydown', handleKeyDown);
    };
  });

  async function loadRepos() {
    loadingRepos = true;
    try {
      const list = await FetchUserRepositories();
      repos = list || [];
    } catch (err) {
      console.error('Erro ao listar repositórios:', err);
    } finally {
      loadingRepos = false;
    }
  }

  let filteredRepos = $derived(
    repos.filter(r => {
      if (!searchTerm) return true;
      const term = searchTerm.toLowerCase();
      return (
        r.full_name.toLowerCase().includes(term) ||
        (r.description && r.description.toLowerCase().includes(term))
      );
    })
  );

  let isAllSelected = $derived(
    selectedRepos.includes('all') || selectedRepos.length === 0
  );

  /** @param {string} repoFullName */
  function toggleRepo(repoFullName) {
    if (repoFullName === 'all') {
      selectedRepos = ['all'];
      return;
    }

    let next = selectedRepos.filter(r => r !== 'all');
    if (next.includes(repoFullName)) {
      next = next.filter(r => r !== repoFullName);
      if (next.length === 0) {
        next = ['all'];
      }
    } else {
      next.push(repoFullName);
    }
    selectedRepos = next;
  }

  function selectAll() {
    selectedRepos = ['all'];
  }

  function clearAll() {
    selectedRepos = ['all'];
  }

  let displayText = $derived.by(() => {
    if (isAllSelected) return '🌐 Todos os Repositórios';
    if (selectedRepos.length === 1) {
      const found = repos.find(r => r.full_name === selectedRepos[0]);
      return found ? (found.private ? '🔒 ' : '📂 ') + found.full_name : selectedRepos[0];
    }
    return `🗂️ ${selectedRepos.length} repositórios selecionados`;
  });
</script>

<div class="repo-selector-container" bind:this={dropdownRef}>
  <div class="repo-label-row">
    <label for="repo-trigger-btn">Repositório Alvo</label>
    {#if loadingRepos}
      <span class="loading-label">(Carregando...)</span>
    {/if}
  </div>

  <button
    id="repo-trigger-btn"
    type="button"
    class="repo-trigger-btn form-control"
    class:active={isOpen}
    onclick={() => !disabled && (isOpen = !isOpen)}
    {disabled}
  >
    <span class="repo-display-text">{displayText}</span>
    <div class="trigger-meta">
      {#if !isAllSelected && selectedRepos.length > 1}
        <span class="repo-count-pill">{selectedRepos.length}</span>
      {/if}
      <span class="chevron" class:open={isOpen}>▼</span>
    </div>
  </button>

  {#if isOpen}
    <div class="repo-dropdown">
      <div class="dropdown-header">
        <input
          type="text"
          class="form-control search-input"
          placeholder="🔍 Filtrar repositórios..."
          bind:value={searchTerm}
          onclick={e => e.stopPropagation()}
        />
        <div class="dropdown-quick-actions">
          <button type="button" class="quick-btn" onclick={selectAll}>
            Selecionar Todos
          </button>
          <button type="button" class="quick-btn" onclick={clearAll}>
            Restaurar
          </button>
        </div>
      </div>

      <div class="repo-list">
        <label class="repo-item" class:selected={isAllSelected}>
          <input
            type="checkbox"
            checked={isAllSelected}
            onchange={() => toggleRepo('all')}
          />
          <span class="repo-name">🌐 <strong>Todos os Repositórios</strong></span>
        </label>

        {#each filteredRepos as repo (repo.id)}
          {@const isChecked = !isAllSelected && selectedRepos.includes(repo.full_name)}
          <label class="repo-item" class:selected={isChecked}>
            <input
              type="checkbox"
              checked={isChecked}
              onchange={() => toggleRepo(repo.full_name)}
            />
            <span class="repo-icon">{repo.private ? '🔒' : '📂'}</span>
            <span class="repo-name" title={repo.description || repo.full_name}>
              {repo.full_name}
            </span>
          </label>
        {/each}

        {#if filteredRepos.length === 0 && repos.length > 0}
          <div class="no-results">Nenhum repositório encontrado</div>
        {/if}
      </div>

      <div class="dropdown-footer">
        <span>{repos.length} repositórios disponíveis</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .repo-selector-container {
    position: relative;
    width: 100%;
  }

  .repo-label-row {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    margin-bottom: 0.5rem;
  }

  .repo-label-row label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted, #8b949e);
  }

  .loading-label {
    font-size: 0.75rem;
    color: var(--accent-light, #58a6ff);
  }

  .repo-trigger-btn {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    background-color: var(--bg-card, #161b22);
    border: 1px solid var(--border-color, #30363d);
    color: var(--text-main, #c9d1d9);
    padding: 0.6rem 0.85rem;
    border-radius: var(--radius-sm, 6px);
    transition: all 0.2s ease;
    text-align: left;
    height: 42px;
  }

  .repo-trigger-btn:hover:not(:disabled),
  .repo-trigger-btn.active {
    border-color: var(--accent-light, #58a6ff);
    box-shadow: 0 0 0 1px var(--accent-light, #58a6ff);
  }

  .repo-trigger-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .repo-display-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.9rem;
  }

  .trigger-meta {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-left: 0.5rem;
    flex-shrink: 0;
  }

  .repo-count-pill {
    background-color: var(--accent, #6366f1);
    color: #ffffff;
    font-size: 0.75rem;
    font-weight: 700;
    padding: 0.1rem 0.45rem;
    border-radius: 10px;
  }

  .chevron {
    font-size: 0.65rem;
    color: var(--text-muted, #8b949e);
    transition: transform 0.2s ease;
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .repo-dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    background-color: var(--bg-surface, #161b22);
    border: 1px solid var(--border-color, #30363d);
    border-radius: var(--radius-md, 8px);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
    z-index: 1000;
    display: flex;
    flex-direction: column;
    max-height: 380px;
    backdrop-filter: blur(8px);
  }

  .dropdown-header {
    padding: 0.6rem;
    border-bottom: 1px solid var(--border-color, #30363d);
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .search-input {
    width: 100%;
    padding: 0.4rem 0.6rem;
    font-size: 0.85rem;
    background-color: var(--bg-main, #0d1117);
    border: 1px solid var(--border-color, #30363d);
    border-radius: var(--radius-sm, 6px);
  }

  .dropdown-quick-actions {
    display: flex;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .quick-btn {
    background: none;
    border: none;
    color: var(--accent-light, #58a6ff);
    font-size: 0.75rem;
    cursor: pointer;
    padding: 0.15rem 0.3rem;
    border-radius: 4px;
    transition: background-color 0.15s ease;
  }

  .quick-btn:hover {
    background-color: rgba(88, 166, 255, 0.1);
  }

  .repo-list {
    overflow-y: auto;
    max-height: 240px;
    padding: 0.35rem 0;
  }

  .repo-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.45rem 0.75rem;
    cursor: pointer;
    font-size: 0.85rem;
    color: var(--text-main, #c9d1d9);
    transition: background-color 0.15s ease;
    user-select: none;
  }

  .repo-item:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }

  .repo-item.selected {
    background-color: rgba(99, 102, 241, 0.12);
  }

  .repo-item input[type="checkbox"] {
    cursor: pointer;
    accent-color: var(--accent, #6366f1);
  }

  .repo-icon {
    font-size: 0.85rem;
  }

  .repo-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .no-results {
    padding: 1rem;
    text-align: center;
    font-size: 0.85rem;
    color: var(--text-muted, #8b949e);
  }

  .dropdown-footer {
    padding: 0.4rem 0.75rem;
    border-top: 1px solid var(--border-color, #30363d);
    font-size: 0.75rem;
    color: var(--text-muted, #8b949e);
    text-align: right;
  }
</style>
