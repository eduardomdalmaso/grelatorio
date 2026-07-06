---
name: GitHub Commit Monthly Reporter developer
description: Assists with developing, updating, and testing the GitHub Commit monthly report generator app, including managing Wails backend, Svelte frontend, and GitHub integration.
---

# GitHub Commit Monthly Reporter Development Skill

This skill contains specialized guidelines and execution steps for managing, updating, and testing the GitHub Commit monthly report generator application.

## 🛠️ Typical Development Tasks

### 1. Modifying Backend Functionality
- **File**: [app.go](file:///c:/Users/eduar/Documents/relatorio/app.go)
- **Instructions**:
  - Keep the Wails App struct thin and modular.
  - Implement helper functions for external HTTP requests to keep `FetchGithubActivity` clean.
  - Make sure all parameters passed from Svelte are validated in Go (e.g. check for empty dates or missing tokens).
  - When exposing new methods, ensure they are capitalized so they are exported and visible to Wails bindings generation.

### 2. Updating the UI Components
- **Files**:
  - Main Layout: [App.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/App.svelte)
  - Settings Pane: [Settings.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/Settings.svelte)
  - Fetch Controller: [ActivityFetch.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/ActivityFetch.svelte)
  - Layout Exporter: [ReportBuilder.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/ReportBuilder.svelte)
- **Instructions**:
  - Keep styles scoped within the `<style>` tags of the Svelte component where possible, or add global utilities to [style.css](file:///c:/Users/eduar/Documents/relatorio/frontend/src/style.css).
  - Avoid generic/unstyled inputs; use customized dark-mode borders, focus outlines (`outline: 2px solid var(--accent-cyan)`), and transition timers.

### 3. Local Development Command Reference
Use these commands to start/restart the environment:

- **Run Dev Server with PM2 (Recommended)**:
  Runs both frontend & backend concurrently in the background.
  ```powershell
  pm2 start ecosystem.config.js
  ```

- **Run Dev Server directly in Terminal**:
  Builds the backend bindings and starts the WebView2 app with live hot-reloading:
  ```powershell
  wails dev
  ```

- **Rebuild the Frontend Packages**:
  Run this if `vite` or other frontend dependencies complain on Windows:
  ```powershell
  cd frontend
  npm install --legacy-peer-deps
  ```

## 📋 Monthly Report Structure Requirements
The generated report format should typically follow this markdown layout for the manager/boss:

```markdown
# Relatório Mensal de Atividades - [Mês/Ano]
**Colaborador:** [Nome do Usuário]
**Cliente/Projeto:** [Nome do Cliente]

## 🛠️ Resumo de Entregas
- **Commits Totais:** [Quantidade]
- **Pull Requests Criados:** [Quantidade]

## 📝 Detalhamento de Commits
### [Nome do Repositório 1]
- **[Hash] - [Título do Commit]** ([Data])
  > [Descrição opcional do commit]

### [Nome do Repositório 2]
- **[Hash] - [Título do Commit]** ([Data])

## 🔀 Pull Requests
- **[Número] - [Título do PR]** (Link: [URL])
```
