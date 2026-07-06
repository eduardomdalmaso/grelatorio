<script>
  import { onMount } from 'svelte';
  import { FetchGithubActivity } from '../../wailsjs/go/main/App.js';

  const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined && (/** @type {any} */(window)).go.main !== undefined;

  // Props using Svelte 5 Runes
  let { selectedActivities = $bindable([]), goToReport } = $props();

  let startDate = $state('');
  let endDate = $state('');
  let activities = $state([]); // The raw list of fetched items
  let loading = $state(false);
  let errorMsg = $state('');
  
  // Search query for UI filtering
  let searchQuery = $state('');

  // Repository filter states
  let activeRepoFilters = $state([]);
  let uniqueRepos = $derived([...new Set(activities.map(act => act.repo))]);

  onMount(() => {
    // Set default dates: startDate is 7 days ago, endDate is today
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
      const result = await FetchGithubActivity(startDate, endDate);
      // Map result to include an 'included' flag for selection and support edits
      activities = (result || []).map(act => {
        // Check if this item is already in selectedActivities
        const isSelected = selectedActivities.some(s => s.ref === act.ref && s.repo === act.repo);
        return {
          ...act,
          included: isSelected,
          editing: false,
          tempTitle: act.title,
          tempDescription: act.description
        };
      });
      
      // Initialize activeRepoFilters to all unique repos
      activeRepoFilters = [...new Set(activities.map(act => act.repo))];

      if (activities.length === 0) {
        errorMsg = 'Nenhuma atividade encontrada neste período.';
      }
    } catch (err) {
      errorMsg = err.message || 'Falha ao buscar atividades.';
      activities = [];
      activeRepoFilters = [];
    } finally {
      loading = false;
    }
  }

  // Filter activities based on search query and repo selections
  let filteredActivities = $derived(
    activities.filter(act => {
      const matchesRepo = activeRepoFilters.length === 0 || activeRepoFilters.includes(act.repo);
      if (!matchesRepo) return false;

      const query = searchQuery.toLowerCase();
      return (
        act.title.toLowerCase().includes(query) ||
        act.repo.toLowerCase().includes(query) ||
        act.description.toLowerCase().includes(query) ||
        act.type.toLowerCase().includes(query)
      );
    })
  );

  // Sync selected changes back to parent
  function handleCheckboxChange(activity) {
    if (activity.included) {
      // Add to selected list
      if (!selectedActivities.some(s => s.ref === activity.ref && s.repo === activity.repo)) {
        selectedActivities = [...selectedActivities, activity];
      }
    } else {
      // Remove from selected list
      selectedActivities = selectedActivities.filter(s => !(s.ref === activity.ref && s.repo === activity.repo));
    }
  }

  function toggleAll(event) {
    const checked = event.target.checked;
    activities = activities.map(act => {
      act.included = checked;
      return act;
    });

    if (checked) {
      // Add all currently filtered activities that aren't already added
      const toAdd = filteredActivities.filter(
        fa => !selectedActivities.some(sa => sa.ref === fa.ref && sa.repo === fa.repo)
      );
      selectedActivities = [...selectedActivities, ...toAdd];
    } else {
      // Remove all currently filtered activities
      selectedActivities = selectedActivities.filter(
        sa => !filteredActivities.some(fa => fa.ref === sa.ref && fa.repo === sa.repo)
      );
    }
  }

  // Edit methods
  function startEdit(act) {
    act.editing = true;
    act.tempTitle = act.title;
    act.tempDescription = act.description;
  }

  function saveEdit(act) {
    act.title = act.tempTitle;
    act.description = act.tempDescription;
    act.editing = false;
    
    // Also update in selectedActivities if included
    if (act.included) {
      selectedActivities = selectedActivities.map(s => {
        if (s.ref === act.ref && s.repo === act.repo) {
          return { ...s, title: act.title, description: act.description };
        }
        return s;
      });
    }
  }

  function cancelEdit(act) {
    act.editing = false;
  }

  function formatDate(dateStr) {
    if (!dateStr) return '';
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('pt-BR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (_) {
      return dateStr;
    }
  }
</script>

<div class="card">
  <h2>Buscar Atividades do GitHub</h2>
  <p style="color: var(--text-muted); margin-bottom: 1.5rem; font-size: 0.9rem;">
    Escolha o período para recuperar seus Commits e Pull Requests. Você poderá revisar, renomear e escolher quais itens farão parte do relatório.
  </p>

  <div class="grid-3" style="align-items: flex-end; gap: 1rem;">
    <div class="form-group" style="margin-bottom: 0;">
      <label for="startDate">Data de Início</label>
      <input id="startDate" type="date" class="form-control" bind:value={startDate} />
    </div>

    <div class="form-group" style="margin-bottom: 0;">
      <label for="endDate">Data de Fim</label>
      <input id="endDate" type="date" class="form-control" bind:value={endDate} />
    </div>

    <button class="btn btn-primary" style="height: 42px;" onclick={fetchActivity} disabled={loading}>
      {#if loading}
        Buscando...
      {:else}
        Buscar Atividade
      {/if}
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
    <div style="display: flex; flex-direction: column; gap: 1rem; margin-bottom: 1.5rem; border-bottom: 1px solid var(--border-color); padding-bottom: 1.25rem;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem;">
        <div style="display: flex; align-items: center; gap: 1rem; flex: 1; min-width: 250px;">
          <input 
            type="text" 
            placeholder="Filtrar resultados por repo, título..." 
            class="form-control" 
            style="width: 100%; max-width: 300px;" 
            bind:value={searchQuery}
          />
          <span style="font-size: 0.85rem; color: var(--text-muted);">
            {filteredActivities.length} itens encontrados ({selectedActivities.length} selecionados)
          </span>
        </div>

        <div style="display: flex; gap: 0.75rem;">
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
        <div class="repo-filter-container">
          <span class="filter-label">Filtrar por Repositório:</span>
          <div class="repo-chips">
            <button 
              class="chip-btn" 
              class:active={activeRepoFilters.length === uniqueRepos.length}
              onclick={() => {
                if (activeRepoFilters.length === uniqueRepos.length) {
                  activeRepoFilters = [];
                } else {
                  activeRepoFilters = [...uniqueRepos];
                }
              }}
            >
              Todos ({uniqueRepos.length})
            </button>
            {#each uniqueRepos as repo}
              <button 
                class="chip-btn" 
                class:active={activeRepoFilters.includes(repo)}
                onclick={() => {
                  if (activeRepoFilters.includes(repo)) {
                    activeRepoFilters = activeRepoFilters.filter(r => r !== repo);
                  } else {
                    activeRepoFilters = [...activeRepoFilters, repo];
                  }
                }}
              >
                {repo.split('/').pop()}
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <table class="activity-table">
      <thead>
        <tr>
          <th style="width: 40px; text-align: center;">
            <input type="checkbox" onchange={toggleAll} checked={filteredActivities.every(a => a.included) && filteredActivities.length > 0} />
          </th>
          <th style="width: 120px;">Tipo</th>
          <th style="width: 200px;">Repositório</th>
          <th>Descrição da Atividade</th>
          <th style="width: 150px;">Data</th>
          <th style="width: 120px; text-align: right;">Ações</th>
        </tr>
      </thead>
      <tbody>
        {#each filteredActivities as act}
          <tr class:row-selected={act.included}>
            <td style="text-align: center; vertical-align: middle;">
              <input type="checkbox" bind:checked={act.included} onchange={() => handleCheckboxChange(act)} />
            </td>
            <td style="vertical-align: middle;">
              <span class="badge badge-{act.type}">
                {act.type === 'commit' ? 'Commit' : 'Pull Request'}
              </span>
            </td>
            <td style="font-weight: 500; font-size: 0.85rem; color: var(--text-muted); vertical-align: middle;">
              {act.repo}
            </td>
            <td>
              {#if act.editing}
                <div class="inline-edit-box">
                  <input type="text" class="form-control" style="width: 100%; margin-bottom: 0.5rem;" bind:value={act.tempTitle} />
                  <textarea class="form-control" style="width: 100%; height: 60px; font-size: 0.85rem;" bind:value={act.tempDescription}></textarea>
                </div>
              {:else}
                <div class="activity-content">
                  <div class="activity-title">{act.title}</div>
                  {#if act.description}
                    <div class="activity-desc">{act.description}</div>
                  {/if}
                  <div class="activity-ref">Ref: <a href={act.url} target="_blank" style="color: var(--accent-light);">{act.ref}</a></div>
                </div>
              {/if}
            </td>
            <td style="font-size: 0.8rem; color: var(--text-muted); vertical-align: middle;">
              {formatDate(act.date)}
            </td>
            <td style="text-align: right; vertical-align: middle;">
              {#if act.editing}
                <div style="display: flex; gap: 0.5rem; justify-content: flex-end;">
                  <button class="btn-icon btn-icon-success" title="Salvar" onclick={() => saveEdit(act)}>✓</button>
                  <button class="btn-icon btn-icon-error" title="Cancelar" onclick={() => cancelEdit(act)}>✕</button>
                </div>
              {:else}
                <button class="btn btn-secondary" style="padding: 0.4rem 0.8rem; font-size: 0.8rem;" onclick={() => startEdit(act)}>
                  Editar
                </button>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}

<style>
  .activity-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }

  .activity-table th {
    background-color: rgba(255, 255, 255, 0.02);
    border-bottom: 2px solid var(--border-color);
    color: var(--text-muted);
    font-size: 0.85rem;
    font-weight: 600;
    padding: 0.75rem 1rem;
    text-align: left;
  }

  .activity-table td {
    padding: 1rem;
    border-bottom: 1px solid var(--border-color);
    font-size: 0.9rem;
  }

  .activity-table tr:hover {
    background-color: rgba(255, 255, 255, 0.01);
  }

  .row-selected {
    background-color: rgba(99, 102, 241, 0.03) !important;
  }

  .badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
  }

  .badge-commit {
    background-color: rgba(99, 102, 241, 0.15);
    color: var(--accent-light);
  }

  .badge-pull_request {
    background-color: rgba(16, 185, 129, 0.15);
    color: #34d399;
  }

  .activity-content {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .activity-title {
    font-weight: 500;
    color: var(--text-main);
  }

  .activity-desc {
    font-size: 0.8rem;
    color: var(--text-muted);
    white-space: pre-wrap;
    max-height: 80px;
    overflow-y: auto;
  }

  .activity-ref {
    font-size: 0.75rem;
    color: var(--text-muted);
  }

  .inline-edit-box {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .btn-icon {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    color: white;
    font-size: 0.85rem;
  }

  .btn-icon-success {
    background-color: var(--success);
  }

  .btn-icon-error {
    background-color: var(--error);
  }

  /* Repo Chips Filters */
  .repo-filter-container {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .filter-label {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--text-muted);
  }

  .repo-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
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
