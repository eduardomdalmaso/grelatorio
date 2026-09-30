<script>
  let {
    filteredActivities = [],
    selectedActivities = $bindable([]),
    onCheckboxChange,
    onToggleAll,
    onStartEdit,
    onSaveEdit,
    onCancelEdit
  } = $props();

  /**
   * @param {string} dateStr
   */
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

<table class="activity-table">
  <thead>
    <tr>
      <th style="width: 40px; text-align: center;">
        <input
          type="checkbox"
          onchange={onToggleAll}
          checked={filteredActivities.length > 0 && filteredActivities.every(a => a.included)}
        />
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
          <input
            type="checkbox"
            bind:checked={act.included}
            onchange={() => onCheckboxChange(act)}
          />
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
              <input
                type="text"
                class="form-control"
                style="width: 100%; margin-bottom: 0.5rem;"
                bind:value={act.tempTitle}
              />
              <textarea
                class="form-control"
                style="width: 100%; height: 60px; font-size: 0.85rem;"
                bind:value={act.tempDescription}
              ></textarea>
            </div>
          {:else}
            <div class="activity-content">
              <div class="activity-title">{act.title}</div>
              {#if act.description}
                <div class="activity-desc">{act.description}</div>
              {/if}
              <div class="activity-ref">
                Ref: <a href={act.url} target="_blank" rel="noreferrer" style="color: var(--accent-light);">{act.ref}</a>
              </div>
            </div>
          {/if}
        </td>
        <td style="font-size: 0.8rem; color: var(--text-muted); vertical-align: middle;">
          {formatDate(act.date)}
        </td>
        <td style="text-align: right; vertical-align: middle;">
          {#if act.editing}
            <div style="display: flex; gap: 0.5rem; justify-content: flex-end;">
              <button class="btn-icon btn-icon-success" title="Salvar" onclick={() => onSaveEdit(act)}>✓</button>
              <button class="btn-icon btn-icon-error" title="Cancelar" onclick={() => onCancelEdit(act)}>✕</button>
            </div>
          {:else}
            <button
              class="btn btn-secondary"
              style="padding: 0.4rem 0.8rem; font-size: 0.8rem;"
              onclick={() => onStartEdit(act)}
            >
              Editar
            </button>
          {/if}
        </td>
      </tr>
    {/each}
  </tbody>
</table>

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
</style>
