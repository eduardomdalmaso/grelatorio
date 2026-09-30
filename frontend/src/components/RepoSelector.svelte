<script>
  import { onMount } from 'svelte';
  import { FetchUserRepositories } from '../../wailsjs/go/main/App.js';

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined;

  let { selectedRepo = $bindable('all'), disabled = false } = $props();

  /** @type {Array<{id: number, name: string, full_name: string, private: boolean}>} */
  let repos = $state([]);
  let loadingRepos = $state(false);

  onMount(async () => {
    if (!isWailsAvailable) return;
    loadingRepos = true;
    try {
      const list = await FetchUserRepositories();
      repos = list || [];
    } catch (err) {
      console.error('Erro ao listar repositórios:', err);
    } finally {
      loadingRepos = false;
    }
  });
</script>

<div class="form-group repo-selector-group">
  <label for="repoSelect">
    Repositório Alvo
    {#if loadingRepos}
      <span class="loading-label">(Carregando...)</span>
    {/if}
  </label>
  <select
    id="repoSelect"
    class="form-control repo-select"
    bind:value={selectedRepo}
    {disabled}
  >
    <option value="all">🌐 Todos os Repositórios</option>
    {#each repos as repo}
      <option value={repo.full_name}>
        {repo.private ? '🔒' : '📂'} {repo.full_name}
      </option>
    {/each}
  </select>
</div>

<style>
  .repo-selector-group {
    margin-bottom: 0;
  }

  .loading-label {
    font-size: 0.75rem;
    color: var(--accent-light, #58a6ff);
    font-weight: normal;
    margin-left: 0.25rem;
  }

  .repo-select {
    cursor: pointer;
    background-color: var(--bg-card, #161b22);
  }
</style>
