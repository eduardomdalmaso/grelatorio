<script>
  import { onMount } from 'svelte';
  import { LoadConfig } from '../../wailsjs/go/main/App.js';

  // Props using Svelte 5 Runes
  let { selectedActivities = $bindable([]) } = $props();

  let title = $state('Relatório Mensal de Atividades');
  let description = $state('Este documento apresenta o detalhamento das atividades desenvolvidas no período, incluindo commits de código e pull requests finalizados.');
  let client = $state('');
  let developerName = $state('Prestador de Serviço');
  
  // Tax / Company fields
  let companyName = $state('');
  let cnpj = $state('');
  
  // Reference month
  let referenceMonth = $state('');
  
  // Date range defaults
  let periodStart = $state('');
  let periodEnd = $state('');
  
  // Billing fields
  let rate = $state('');
  let hours = $state('0');

  // Month Indicators
  let bugsFixed = $state('0');
  let deploysCount = $state('0');

  // Checklist
  let checkVersioned = $state(true);
  let checkNoRetention = $state(true);
  let checkTested = $state(true);
  let checkNoDependencies = $state(true);

  onMount(async () => {
    try {
      const config = await LoadConfig();
      client = config.client || '';
      rate = config.rate || '';
      developerName = config.username || 'Prestador de Serviço';
      cnpj = config.cnpj || '';
      companyName = config.company_name || '';
      
      // Auto fill start/end dates from activities if present
      if (selectedActivities.length > 0) {
        const timestamps = selectedActivities
          .map(a => new Date(a.date).getTime())
          .filter(t => !isNaN(t));
        if (timestamps.length > 0) {
          const minDate = new Date(Math.min(...timestamps));
          const maxDate = new Date(Math.max(...timestamps));
          periodStart = minDate.toISOString().split('T')[0];
          periodEnd = maxDate.toISOString().split('T')[0];
          
          // Auto fill reference month from periodEnd date
          const months = ['janeiro', 'fevereiro', 'março', 'abril', 'maio', 'junho', 'julho', 'agosto', 'setembro', 'outubro', 'novembro', 'dezembro'];
          referenceMonth = `${months[maxDate.getMonth()]}/${maxDate.getFullYear()}`;
        }
      } else {
        // Fallback to current month if no activities selected yet
        const today = new Date();
        const months = ['janeiro', 'fevereiro', 'março', 'abril', 'maio', 'junho', 'julho', 'agosto', 'setembro', 'outubro', 'novembro', 'dezembro'];
        referenceMonth = `${months[today.getMonth()]}/${today.getFullYear()}`;
      }
    } catch (_) {}
  });

  // Calculate hours dynamically for each activity based on commit timestamps (7h to 17h)
  function calculateHoursFromCommits() {
    /** @type {Record<string, any[]>} */
    const dateGroups = {};
    
    selectedActivities.forEach(act => {
      if (!act.date) return;
      const d = new Date(act.date);
      if (isNaN(d.getTime())) return;
      
      const year = d.getFullYear();
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const day = String(d.getDate()).padStart(2, '0');
      const dateKey = `${year}-${month}-${day}`;
      
      if (!dateGroups[dateKey]) {
        dateGroups[dateKey] = [];
      }
      dateGroups[dateKey].push(act);
    });
    
    Object.values(dateGroups).forEach(group => {
      // Sort chronologically
      group.sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
      
      let prevClampedHour = 7.0; // workday starts at 7h
      
      group.forEach((act) => {
        const d = new Date(act.date);
        const hour = d.getHours() + d.getMinutes() / 60;
        const clampedHour = Math.min(17.0, Math.max(7.0, hour));
        
        let diff = clampedHour - prevClampedHour;
        if (diff < 0.5) {
          diff = 0.5; // Minimum 30 mins
        }
        
        act.hours = parseFloat(diff.toFixed(1));
        prevClampedHour = clampedHour;
      });
    });
    
    // Force Svelte reactivity update
    selectedActivities = [...selectedActivities];
    recalculateTotalHours();
  }

  // Effect to auto-calculate and distribute hours dynamically
  $effect(() => {
    if (selectedActivities.length > 0) {
      const needsHours = selectedActivities.some(act => act.hours === undefined);
      if (needsHours) {
        calculateHoursFromCommits();
      } else {
        recalculateTotalHours();
      }
    } else {
      hours = '0';
    }
  });

  // Effect to auto-count and register bugs and deploys from commit titles
  $effect(() => {
    if (selectedActivities.length > 0) {
      let bugs = 0;
      let deploys = 0;
      
      selectedActivities.forEach(act => {
        const titleLower = (act.title || '').toLowerCase();
        const isBug = /fix|bug|corrigido|resolvido|ajuste|conserto/i.test(titleLower);
        const isDeploy = /deploy|release|prod|publicado|entrega/i.test(titleLower);
        
        if (isBug) bugs++;
        if (isDeploy) deploys++;
      });
      
      bugsFixed = String(bugs);
      deploysCount = String(deploys);
    } else {
      bugsFixed = '0';
      deploysCount = '0';
    }
  });

  function recalculateTotalHours() {
    const total = selectedActivities.reduce((acc, act) => acc + (parseFloat(act.hours) || 0), 0);
    hours = String(parseFloat(total.toFixed(1)));
  }

  function distributeEqually() {
    const totalHours = parseFloat(hours) || 0;
    if (selectedActivities.length > 0 && totalHours > 0) {
      const avg = parseFloat((totalHours / selectedActivities.length).toFixed(1));
      selectedActivities = selectedActivities.map(act => ({
        ...act,
        hours: avg
      }));
    }
  }

  // Group activities by Repository
  let activitiesByRepo = $derived.by(() => {
    /** @type {Record<string, any[]>} */
    const groups = {};
    selectedActivities.forEach(act => {
      if (!groups[act.repo]) {
        groups[act.repo] = [];
      }
      groups[act.repo].push(act);
    });
    return groups;
  });

  async function handleExport() {
    const element = document.getElementById('report-pdf-template');
    if (!element) return;
    
    // Load html2pdf dynamically to split chunk from main bundle
    const { default: html2pdf } = await import('html2pdf.js');
    
    /** @type {any} */
    const opt = {
      margin: 0,
      filename: `Relatorio_Atividades_${client.replace(/[^a-zA-Z0-9]/g, '_')}_${referenceMonth.replace('/', '_')}.pdf`,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { 
        scale: 2, 
        useCORS: true,
        logging: false,
        backgroundColor: '#ffffff'
      },
      jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
      pagebreak: { 
        mode: ['avoid-all', 'css'], 
        avoid: ['tr', '.pdf-signatures', '.pdf-checklist-block', '.pdf-indicators-block'] 
      }
    };

    html2pdf().set(opt).from(element).save();
  }

  /** @param {string} dateStr */
  function formatDateBR(dateStr) {
    if (!dateStr) return '';
    try {
      const parts = dateStr.split('-');
      if (parts.length === 3) {
        return `${parts[2]}/${parts[1]}/${parts[0]}`;
      }
      const d = new Date(dateStr);
      return d.toLocaleDateString('pt-BR');
    } catch (_) {
      return dateStr;
    }
  }

  /** @param {string} dateStr */
  function formatActivityDate(dateStr) {
    if (!dateStr) return '';
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit', year: '2-digit' });
    } catch (_) {
      return dateStr;
    }
  }
</script>

<div class="report-builder-layout">
  <!-- Controls Panel (Left) -->
  <div class="controls-panel">
    <div class="card">
      <h3>Informações Fiscais & Período</h3>
      
      <div class="form-group">
        <label for="rep-comp-name">Sua Razão Social (Empresa)</label>
        <input id="rep-comp-name" type="text" class="form-control" bind:value={companyName} />
      </div>

      <div class="form-group">
        <label for="rep-cnpj">Seu CNPJ</label>
        <input id="rep-cnpj" type="text" class="form-control" bind:value={cnpj} />
      </div>

      <div class="form-group">
        <label for="rep-dev">Responsável Técnico</label>
        <input id="rep-dev" type="text" class="form-control" bind:value={developerName} />
      </div>

      <div class="form-group">
        <label for="rep-ref-month">Mês de Referência</label>
        <input id="rep-ref-month" type="text" class="form-control" placeholder="Ex: março/2026" bind:value={referenceMonth} />
      </div>

      <div class="form-group" style="margin-top: 0.75rem;">
        <label for="rep-start">Período Início</label>
        <input id="rep-start" type="date" class="form-control" bind:value={periodStart} />
      </div>
      <div class="form-group" style="margin-top: 0.75rem;">
        <label for="rep-end">Período Fim</label>
        <input id="rep-end" type="date" class="form-control" bind:value={periodEnd} />
      </div>
    </div>

    <div class="card">
      <h3>Metadados do Relatório</h3>
      <div class="form-group">
        <label for="rep-title">Título do Relatório</label>
        <input id="rep-title" type="text" class="form-control" bind:value={title} />
      </div>
      <div class="form-group">
        <label for="rep-client">Nome do Cliente / Tomador</label>
        <input id="rep-client" type="text" class="form-control" bind:value={client} />
      </div>
      <div class="form-group">
        <label for="rep-desc">Descrição / Escopo Geral</label>
        <textarea id="rep-desc" class="form-control" style="height: 60px;" bind:value={description}></textarea>
      </div>
    </div>

    <div class="card">
      <h3>Indicadores do Mês</h3>
      <div class="grid-2">
        <div class="form-group">
          <label for="rep-bugs">Bugs Corrigidos</label>
          <input id="rep-bugs" type="number" class="form-control" bind:value={bugsFixed} min="0" />
        </div>
        <div class="form-group">
          <label for="rep-deploys">Deploys Realizados</label>
          <input id="rep-deploys" type="number" class="form-control" bind:value={deploysCount} min="0" />
        </div>
      </div>
      
      <div class="form-group" style="margin-top: 0.5rem;">
        <label for="rep-hours">Horas Trabalhadas (Calculado: 10h/dia)</label>
        <input id="rep-hours" type="number" step="0.5" class="form-control" bind:value={hours} />
        <small style="color: var(--text-muted); font-size: 0.75rem;">Calculado automaticamente (7h às 17h, Seg a Dom).</small>
      </div>
    </div>

    {#if selectedActivities.length > 0}
      <div class="card">
        <h3>Detalhamento de Horas por Item</h3>
        <p style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 0.75rem;">
          Altere as horas de cada item. O total de horas do mês será recalculado.
        </p>
        <div class="activities-hours-list">
          {#each selectedActivities as act}
            <div class="activity-hour-row">
              <span class="activity-hour-title" title={act.title}>
                <span class="repo-badge">{act.repo.split('/').pop()}</span> {act.title}
              </span>
              <input 
                type="number" 
                step="0.5" 
                min="0"
                class="form-control hour-input-field" 
                bind:value={act.hours}
                oninput={recalculateTotalHours}
              />
            </div>
          {/each}
        </div>
        <div style="display: flex; gap: 0.5rem; margin-top: 0.75rem;">
          <button class="btn btn-secondary" style="flex: 1; font-size: 0.8rem; padding: 0.5rem;" onclick={distributeEqually} disabled={parseFloat(hours) <= 0}>
            Distribuir Igualmente ({hours}h)
          </button>
          <button class="btn btn-secondary" style="flex: 1; font-size: 0.8rem; padding: 0.5rem;" onclick={calculateHoursFromCommits}>
            Calcular por Commits
          </button>
        </div>
      </div>
    {/if}

    <div class="card">
      <h3>Checklist de Entrega</h3>
      <div class="checkbox-group">
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={checkVersioned} />
          Código versionado no repositório oficial
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={checkNoRetention} />
          Sem retenção de acessos ou credenciais
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={checkTested} />
          Código funcional e testado
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={checkNoDependencies} />
          Sem dependências ocultas
        </label>
      </div>
    </div>

    <button class="btn btn-success btn-block" style="width: 100%; height: 48px; font-size: 1rem; margin-bottom: 2rem;" onclick={handleExport} disabled={selectedActivities.length === 0}>
      ⇩ Exportar PDF Oficial
    </button>
  </div>

  <!-- Live Preview Panel (Right) -->
  <div class="preview-panel">
    <div class="preview-header-bar">
      <span>Visualização do PDF (A4)</span>
      <span style="font-size: 0.8rem; color: var(--text-muted);">{selectedActivities.length} itens incluídos</span>
    </div>

    <div class="pdf-container">
      <!-- The A4 Page container -->
      <div id="report-pdf-template" class="pdf-page">
        <!-- Header -->
        <div class="pdf-header">
          <div class="pdf-header-main">
            <div class="pdf-logo">RELATÓRIO MENSAL DE PRESTAÇÃO DE SERVIÇOS</div>
            <div class="pdf-doc-info">
              <div><strong>Emitido em:</strong> {new Date().toLocaleDateString('pt-BR')}</div>
              <div><strong>Período:</strong> {formatDateBR(periodStart)} a {formatDateBR(periodEnd)}</div>
            </div>
          </div>
          <div class="pdf-header-bar-accent"></div>
        </div>

        <!-- Body -->
        <div class="pdf-body">
          <!-- Metadata Table (Spreadsheet style layout) -->
          <div class="pdf-meta-table">
            <table>
              <tbody>
                <tr>
                  <td class="pdf-meta-label">Empresa:</td>
                  <td class="pdf-meta-value">{companyName}</td>
                  <td class="pdf-meta-label">Mês de Referência:</td>
                  <td class="pdf-meta-value">{referenceMonth}</td>
                </tr>
                <tr>
                  <td class="pdf-meta-label">CNPJ:</td>
                  <td class="pdf-meta-value">{cnpj}</td>
                  <td class="pdf-meta-label">Responsável Técnico:</td>
                  <td class="pdf-meta-value">{developerName}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="pdf-section" style="margin-top: 15px;">
            <h1 class="pdf-title">{title}</h1>
            <p class="pdf-description">{description}</p>
          </div>

          <!-- Activity List -->
          <div class="pdf-section" style="margin-top: 15px;">
            <h2 class="pdf-section-title">Detalhamento das Atividades Realizadas</h2>
            
            {#if Object.keys(activitiesByRepo).length === 0}
              <p style="font-style: italic; color: #666; font-size: 0.9rem;">Nenhuma atividade selecionada para exibição no relatório.</p>
            {:else}
              {#each Object.entries(activitiesByRepo) as [repo, items]}
                <div class="pdf-repo-group">
                  <h3 class="pdf-repo-title">Projeto / Repositório: {repo.split('/').pop()}</h3>
                  <table class="pdf-table">
                    <thead>
                      <tr>
                        <th>Descrição da Atividade</th>
                        <th style="width: 80px; text-align: center;">Data</th>
                        <th style="width: 65px; text-align: center;">Tipo</th>
                        <th style="width: 55px; text-align: center;">Horas</th>
                        <th style="width: 85px; text-align: center;">Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      {#each items as item}
                        <tr>
                          <td>
                            <strong>{item.title}</strong>
                            {#if item.description}
                              <div style="font-size: 0.8rem; color: #666; margin-top: 2px;">{item.description}</div>
                            {/if}
                          </td>
                          <td style="color: #666; font-size: 0.85rem; text-align: center;">{formatActivityDate(item.date)}</td>
                          <td style="font-size: 0.8rem; text-align: center; font-weight: 500;">
                            {item.type === 'commit' ? 'Commit' : 'PR'}
                          </td>
                          <td style="text-align: center; color: #0f172a; font-weight: 500;">
                            {item.hours !== undefined ? item.hours : 0}h
                          </td>
                          <td style="color: #10b981; font-weight: 600; font-size: 0.8rem; text-align: center;">
                            {item.type === 'commit' ? 'Concluído' : 'Merged'}
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              {/each}
            {/if}
          </div>

          <!-- Month Indicators -->
          <div class="pdf-section pdf-indicators-block" style="margin-top: 20px; page-break-inside: avoid;">
            <h2 class="pdf-section-title">Indicadores do Mês</h2>
            <div class="pdf-indicators-grid">
              <div class="pdf-indicator-item">
                <span class="pdf-indicator-label">Total de Horas Trabalhadas:</span>
                <span class="pdf-indicator-val">{hours}h</span>
              </div>
              <div class="pdf-indicator-item">
                <span class="pdf-indicator-label">Total de Tarefas Registradas:</span>
                <span class="pdf-indicator-val">{selectedActivities.length}</span>
              </div>
              <div class="pdf-indicator-item">
                <span class="pdf-indicator-label">Bugs Corrigidos:</span>
                <span class="pdf-indicator-val">{bugsFixed}</span>
              </div>
              <div class="pdf-indicator-item">
                <span class="pdf-indicator-label">Deploys Realizados:</span>
                <span class="pdf-indicator-val">{deploysCount}</span>
              </div>
            </div>
          </div>

          <!-- Terms Checklist -->
          <div class="pdf-section pdf-checklist-block" style="margin-top: 20px; page-break-inside: avoid;">
            <div class="pdf-checklist-title">DECLARAÇÕES E TERMOS DE ENTREGA</div>
            <div class="pdf-checklist-grid">
              <div class="pdf-checklist-item">
                <span class="pdf-checkbox-box">{checkVersioned ? '☑' : '☐'}</span>
                <span>Código versionado no repositório oficial</span>
              </div>
              <div class="pdf-checklist-item">
                <span class="pdf-checkbox-box">{checkNoRetention ? '☑' : '☐'}</span>
                <span>Sem retenção de acessos ou credenciais</span>
              </div>
              <div class="pdf-checklist-item">
                <span class="pdf-checkbox-box">{checkTested ? '☑' : '☐'}</span>
                <span>Código funcional e testado</span>
              </div>
              <div class="pdf-checklist-item">
                <span class="pdf-checkbox-box">{checkNoDependencies ? '☑' : '☐'}</span>
                <span>Sem dependências ocultas</span>
              </div>
            </div>
          </div>

          <!-- Signatures -->
          <div class="pdf-signatures" style="margin-top: 50px; page-break-inside: avoid;">
            <div class="pdf-sig-line">
              <div class="pdf-sig-bar"></div>
              <div>{developerName}</div>
              <div style="font-size: 0.8rem; color: #666;">Responsável Técnico / Desenvolvedor</div>
              <div style="font-size: 0.75rem; color: #94a3b8; margin-top: 4px;">Data: {new Date().toLocaleDateString('pt-BR')}</div>
            </div>
            <div class="pdf-sig-line">
              <div class="pdf-sig-bar"></div>
              <div>{client || 'Assinatura do Cliente'}</div>
              <div style="font-size: 0.8rem; color: #666;">Representante Autorizado / Cliente</div>
              <div style="font-size: 0.75rem; color: #94a3b8; margin-top: 4px;">Data: ____/____/______</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .report-builder-layout {
    display: flex;
    gap: 1rem;
    height: 100%;
    overflow: hidden;
  }

  .controls-panel {
    width: 380px;
    flex-shrink: 0;
    overflow-y: auto;
    padding-right: 0.5rem;
  }

  .controls-panel .card {
    padding: 0.85rem;
    margin-bottom: 0.65rem;
  }

  .controls-panel .grid-2 {
    gap: 0.75rem;
  }

  .controls-panel .form-group {
    margin-bottom: 0.65rem;
  }

  .controls-panel .form-group label {
    font-size: 0.8rem;
  }

  .controls-panel .form-control {
    font-size: 0.85rem;
    padding: 0.6rem 0.8rem;
  }

  .controls-panel h3 {
    font-size: 0.9rem;
    margin-bottom: 0.65rem;
  }

  .activities-hours-list {
    max-height: 250px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding-right: 0.25rem;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    padding: 0.75rem;
    background-color: rgba(0, 0, 0, 0.15);
  }

  .activity-hour-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 0.75rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.03);
    padding-bottom: 0.5rem;
  }

  .activity-hour-row:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }

  .activity-hour-title {
    font-size: 0.8rem;
    color: var(--text-main);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .repo-badge {
    background-color: rgba(99, 102, 241, 0.15);
    color: var(--accent-light);
    font-size: 0.7rem;
    padding: 0.1rem 0.35rem;
    border-radius: 4px;
    font-weight: 600;
  }

  .hour-input-field {
    width: 75px !important;
    height: 30px !important;
    padding: 0.2rem 0.5rem !important;
    text-align: center;
  }

  .preview-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background-color: #1e293b;
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
  }

  .preview-header-bar {
    background-color: var(--bg-panel);
    padding: 0.75rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.9rem;
    font-weight: 500;
  }

  .pdf-container {
    flex: 1;
    overflow-y: auto;
    overflow-x: auto;
    padding: 2rem;
    display: flex;
    justify-content: center;
    align-items: flex-start;
  }


  .checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    color: var(--text-muted);
    cursor: pointer;
  }

  .checkbox-label input {
    accent-color: var(--accent-cyan);
    width: 16px;
    height: 16px;
  }

  /* PDF Styles */
  .pdf-page {
    width: 210mm;
    min-height: 297mm;
    background-color: #ffffff;
    color: #1e293b;
    padding: 20mm;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
    font-family: 'Inter', Arial, sans-serif;
    font-size: 10pt;
    line-height: 1.5;
    position: relative;
    box-sizing: border-box;
    text-align: left;
  }

  .pdf-header {
    margin-bottom: 20px;
  }

  .pdf-header-main {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    padding-bottom: 10px;
  }

  .pdf-logo {
    font-size: 11pt;
    font-weight: 800;
    color: #4f46e5;
    letter-spacing: -0.02em;
  }

  .pdf-doc-info {
    text-align: right;
    font-size: 8.5pt;
    color: #475569;
  }

  .pdf-header-bar-accent {
    height: 4px;
    background: linear-gradient(90deg, #4f46e5 0%, #a855f7 100%);
    border-radius: 2px;
  }

  .pdf-title {
    font-size: 16pt;
    font-weight: 700;
    color: #0f172a;
    margin-bottom: 10px;
    letter-spacing: -0.025em;
  }

  .pdf-description {
    color: #475569;
    font-size: 9pt;
    line-height: 1.4;
    margin-bottom: 15px;
  }

  .pdf-meta-table {
    width: 100%;
    margin-bottom: 20px;
  }

  .pdf-meta-table table {
    width: 100%;
    border-collapse: collapse;
  }

  .pdf-meta-table td {
    padding: 6px 10px;
    border: 1px solid #cbd5e1;
    font-size: 9pt;
    vertical-align: middle;
  }

  .pdf-meta-label {
    background-color: #f1f5f9;
    font-weight: 700;
    color: #475569;
    width: 20%;
  }

  .pdf-meta-value {
    color: #0f172a;
    width: 30%;
    font-weight: 600;
  }

  .pdf-section-title {
    font-size: 11pt;
    font-weight: 700;
    color: #0f172a;
    border-bottom: 2px solid #e2e8f0;
    padding-bottom: 4px;
    margin-bottom: 12px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .pdf-repo-group {
    margin-bottom: 20px;
  }

  .pdf-repo-title {
    font-size: 9.5pt;
    font-weight: 700;
    color: #4f46e5;
    margin-bottom: 6px;
  }

  .pdf-table {
    width: 100%;
    border-collapse: collapse;
  }

  .pdf-table tr {
    page-break-inside: avoid;
    break-inside: avoid;
  }

  .pdf-table th {
    background-color: #f8fafc;
    border: 1px solid #cbd5e1;
    font-size: 8.5pt;
    font-weight: 700;
    color: #475569;
    padding: 6px 8px;
    text-align: left;
  }

  .pdf-table td {
    border: 1px solid #e2e8f0;
    padding: 6px 8px;
    font-size: 8.5pt;
    vertical-align: top;
  }

  .pdf-indicators-block {
    background-color: #f8fafc;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    padding: 15px;
  }

  .pdf-indicators-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px 20px;
  }

  .pdf-indicator-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px dashed #e2e8f0;
    padding-bottom: 4px;
  }

  .pdf-indicator-label {
    font-size: 9pt;
    color: #475569;
    font-weight: 500;
  }

  .pdf-indicator-val {
    font-size: 10pt;
    font-weight: 700;
    color: #0f172a;
  }

  .pdf-checklist-block {
    background-color: #faf5ff;
    border: 1px solid #e9d5ff;
    border-radius: 6px;
    padding: 15px;
  }

  .pdf-checklist-title {
    font-size: 9pt;
    font-weight: 800;
    color: #6b21a8;
    margin-bottom: 10px;
    letter-spacing: 0.05em;
  }

  .pdf-checklist-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px 15px;
  }

  .pdf-checklist-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 8.5pt;
    color: #581c87;
  }

  .pdf-checkbox-box {
    font-size: 12pt;
    font-family: monospace;
    line-height: 1;
    color: #7e22ce;
  }

  .pdf-signatures {
    display: flex;
    justify-content: space-between;
    gap: 40px;
  }

  .pdf-sig-line {
    flex: 1;
    text-align: center;
  }

  .pdf-sig-bar {
    height: 1px;
    background-color: #94a3b8;
    margin-bottom: 8px;
  }
</style>
