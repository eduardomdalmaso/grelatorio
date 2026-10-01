<script>
  import { onMount } from 'svelte';
  import { FetchGithubActivity } from '../../wailsjs/go/main/App.js';
  import RepoSelector from './RepoSelector.svelte';
  import ActivityTable from './ActivityTable.svelte';

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined;

  let { selectedActivities = $bindable([]), goToReport } = $props();

  let startDate = $state('');
  let endDate = $state('');
  /** @type {string[]} */
  let selectedTargetRepos = $state(['all']);
  /** @type {any[]} */
  let activities = $state([]);
  let loading = $state(false);
  let errorMsg = $state('');
  let searchQuery = $state('');
  /** @type {string[]} */
  let activeRepoFilters = $state([]);

  let uniqueRepos = $derived([...new Set(activities.map(act => act.repo))]);

  onMount(() => {
    const today = new Date();
    const lastWeek = new Date();
    lastWeek.setDate(today.getDate() - 7);
    endDate = today.toISOString().split('T')[0];
    startDate = lastWeek.toISOString().split('T')[0];
  });

  async function fetchActivity() {
    if (!isWailsAvailable) {
      errorMsg = 'O backend Go não está acessível a partir do navegador convencional.';
      return;
    }
    if (!startDate || !endDate) {
      errorMsg = 'Por favor, selecione as datas de início e fim.';
      return;
    }

    loading = true;
    errorMsg = '';

    try {
      const repoFilter = selectedTargetRepos.includes('all') || selectedTargetRepos.length === 0
        ? ''
        : selectedTargetRepos.join(',');
      const result = await FetchGithubActivity(startDate, endDate, repoFilter);
      activities = (result || []).map(act => ({
        ...act,
        included: selectedActivities.some(s => s.ref === act.ref && s.repo === act.repo),
        editing: false,
        tempTitle: act.title,
        tempDescription: act.description
      }));

      activeRepoFilters = [...new Set(activities.map(act => act.repo))];
      if (activities.length === 0) {
        errorMsg = 'Nenhuma atividade encontrada para os filtros selecionados.';
      }
    } catch (err) {
      errorMsg = (/** @type {Error} */ (err)).message || 'Falha ao buscar atividades.';
      activities = [];
      activeRepoFilters = [];
    } finally {
      loading = false;
    }
  }

  let filteredActivities = $derived(
    activities.filter(act => {
      const matchesRepo = activeRepoFilters.length === 0 || activeRepoFilters.includes(act.repo);
      if (!matchesRepo) return false;
      const q = searchQuery.toLowerCase();
      return (
        act.title.toLowerCase().includes(q) ||
        act.repo.toLowerCase().includes(q) ||
        act.description.toLowerCase().includes(q) ||
        act.type.toLowerCase().includes(q)
      );
    })
  );

  /**
   * @param {any} activity
   */
  function handleCheckboxChange(activity) {
    if (activity.included) {
      if (!selectedActivities.some(s => s.ref === activity.ref && s.repo === activity.repo)) {
        selectedActivities = [...selectedActivities, activity];
      }
    } else {
      selectedActivities = selectedActivities.filter(s => !(s.ref === activity.ref && s.repo === activity.repo));
    }
  }

  /**
   * @param {Event} event
   */
  function toggleAll(event) {
    const checked = (/** @type {HTMLInputElement} */ (event.target)).checked;
    activities = activities.map(act => {
      act.included = checked;
      return act;
    });

    if (checked) {
      const toAdd = filteredActivities.filter(
        fa => !selectedActivities.some(sa => sa.ref === fa.ref && sa.repo === fa.repo)
      );
      selectedActivities = [...selectedActivities, ...toAdd];
    } else {
      selectedActivities = selectedActivities.filter(
        sa => !filteredActivities.some(fa => fa.ref === sa.ref && fa.repo === sa.repo)
      );
    }
  }

  /**
   * @param {any} act
   */
  function startEdit(act) {
    act.editing = true;
    act.tempTitle = act.title;
    act.tempDescription = act.description;
  }

  /**
   * @param {any} act
   */
  function saveEdit(act) {
    act.title = act.tempTitle;
    act.description = act.tempDescription;
    act.editing = false;
    if (act.included) {
      selectedActivities = selectedActivities.map(s => {
        if (s.ref === act.ref && s.repo === act.repo) {
          return { ...s, title: act.title, description: act.description };
        }
        return s;
      });
    }
  }

  /**
   * @param {any} act
   */
  function cancelEdit(act) {
    act.editing = false;
  }
</script>

<div class="card">
  <h2>Buscar Atividades do GitHub</h2>
  <p style="color: var(--text-muted); margin-bottom: 1.5rem; font-size: 0.9rem;">
    Escolha o período e o repositório para recuperar seus Commits e Pull Requests.
  </p>

  <div class="grid-form">
    <div class="form-group" style="margin-bottom: 0;">
      <label for="startDate">Data de Início</label>
      <input id="startDate" type="date" class="form-control" bind:value={startDate} />
    </div>

    <div class="form-group" style="margin-bottom: 0;">
      <label for="endDate">Data de Fim</label>
      <input id="endDate" type="date" class="form-control" bind:value={endDate} />
    </div>

    <RepoSelector bind:selectedRepos={selectedTargetRepos} disabled={loading} />

    <button class="btn btn-primary btn-fetch" onclick={fetchActivity} disabled={loading}>
      {loading ? 'Buscando...' : '🔍 Buscar Atividade'}
    </button>
  </div>
</div>

{#if errorMsg}
  <div class="status-msg status-error">
    <span>⚠</span>
    <span>{errorMsg}</span>
  </div>
{/if}

{#if activities.length > 0}
  <div class="card">
    <div class="filter-header">
      <div class="search-bar">
        <input
          type="text"
          placeholder="Filtrar por título, repo..."
          class="form-control"
          style="width: 100%; max-width: 280px;"
          bind:value={searchQuery}
        />
        <span class="count-badge">
          {filteredActivities.length} encontrados ({selectedActivities.length} selecionados)
        </span>
      </div>

      <div class="action-buttons">
        <button class="btn btn-secondary" onclick={() => {
          activities = activities.map(a => ({ ...a, included: false }));
          selectedActivities = [];
        }}>
          Limpar Seleção
        </button>
        <button class="btn btn-success" onclick={goToReport} disabled={selectedActivities.length === 0}>
          Criar Relatório ({selectedActivities.length}) →
        </button>
      </div>
    </div>

    {#if uniqueRepos.length > 1}
      <div class="repo-chips">
        <span class="filter-label">Filtrar:</span>
        <button
          class="chip-btn"
          class:active={activeRepoFilters.length === uniqueRepos.length}
          onclick={() => activeRepoFilters = activeRepoFilters.length === uniqueRepos.length ? [] : [...uniqueRepos]}
        >
          Todos ({uniqueRepos.length})
        </button>
        {#each uniqueRepos as repo}
          <button
            class="chip-btn"
            class:active={activeRepoFilters.includes(repo)}
            onclick={() => {
              activeRepoFilters = activeRepoFilters.includes(repo)
                ? activeRepoFilters.filter(r => r !== repo)
                : [...activeRepoFilters, repo];
            }}
          >
            {repo.split('/').pop()}
          </button>
        {/each}
      </div>
    {/if}

    <ActivityTable
      {filteredActivities}
      bind:selectedActivities={selectedActivities}
      onCheckboxChange={handleCheckboxChange}
      onToggleAll={toggleAll}
      onStartEdit={startEdit}
      onSaveEdit={saveEdit}
      onCancelEdit={cancelEdit}
    />
  </div>
{/if}

<style>
  .grid-form {
    display: grid;
    grid-template-columns: 1fr 1fr 1.5fr auto;
    align-items: flex-end;
    gap: 1rem;
  }

  @media (max-width: 900px) {
    .grid-form {
      grid-template-columns: 1fr;
    }
  }

  .btn-fetch {
    height: 42px;
    white-space: nowrap;
    padding: 0 1.25rem;
  }

  .filter-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 1rem;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 1rem;
    margin-bottom: 1rem;
  }

  .search-bar {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
    min-width: 260px;
  }

  .count-badge {
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .action-buttons {
    display: flex;
    gap: 0.75rem;
  }

  .repo-chips {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .filter-label {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--text-muted);
  }

  .chip-btn {
    background-color: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--border-color);
    color: var(--text-muted);
    padding: 0.35rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .chip-btn:hover {
    background-color: rgba(255, 255, 255, 0.08);
    color: var(--text-main);
  }

  .chip-btn.active {
    background-color: rgba(99, 102, 241, 0.15);
    border-color: var(--accent);
    color: var(--accent-light);
  }
</style>
