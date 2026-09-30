# Rules for GitHub Report Generator Development

Here are the guidelines, architecture, and UI/UX design rules for the GitHub monthly commit report generator application.

---

## 🏛️ Project Architecture & DDD - Hexagonal Standards

The application is built using **Wails (v2)**, binding a **Go backend** with a **Svelte + Vite frontend**, following **Domain-Driven Design (DDD) and Hexagonal Architecture (Ports & Adapters)**.

### 1. Go Backend Architecture (`internal/`)
The Go codebase must strictly respect the layer boundaries:
- **`internal/core/domain/`**: Pure business domain models, entities, and validations (`Config`, `GithubActivity`, `Repository`). No external I/O or frameworks.
- **`internal/core/ports/`**: Interfaces defining input and output boundaries (`ConfigRepository`, `GithubService`).
- **`internal/core/usecases/`**: Application use cases orchestrating domain entities and ports (`ConfigUseCase`, `ActivityUseCase`).
- **`internal/adapters/secondary/`**: Infrastructure implementations:
  - `storage/`: Local configuration persistence (`FileConfigRepository`).
  - `github/`: HTTP client communication with GitHub REST APIs (`GithubAPIClient`).
- **`app.go` & `main.go`**: Primary Adapter (Wails binding façade) exposing clean methods to the frontend.

### 2. Svelte Frontend (`frontend/src/`)
- **Component Breakdown & Modularity**:
  - [Settings.svelte](file:///c:/Users/hades/Documents/grelatorio/frontend/src/components/Settings.svelte): Manage user token, username, rates, and client details.
  - [ActivityFetch.svelte](file:///c:/Users/hades/Documents/grelatorio/frontend/src/components/ActivityFetch.svelte): Handle date selections, repository targeting, and orchestrating search.
  - [RepoSelector.svelte](file:///c:/Users/hades/Documents/grelatorio/frontend/src/components/RepoSelector.svelte): Dropdown selector for targeting specific user repositories or all repositories.
  - [ActivityTable.svelte](file:///c:/Users/hades/Documents/grelatorio/frontend/src/components/ActivityTable.svelte): Table with checkboxes, badges, and inline editing.
  - [ReportBuilder.svelte](file:///c:/Users/hades/Documents/grelatorio/frontend/src/components/ReportBuilder.svelte): Preview, group, and format selected activities into a markdown/text report and PDF.

---

## 📏 Agent Rule: File Size & Modularity Limits (Anti-Monolith)

To prevent code degradation and keep token efficiency high:
1. **Maximum File Size**:
   - Hard limit of **300 lines of code** per file (absolute max 400 lines for complex Svelte views with embedded CSS).
   - If a file approaches 300 lines, the agent **MUST** decompose it into smaller sub-components, domain entities, or utility services.
2. **Single Responsibility Principle (SRP)**:
   - One primary responsibility per component/file.
   - Do not mix UI presentation, state manipulation, and external API mappings in a single file.
3. **No Direct External Calls from Frontend**:
   - Always route all GitHub and file system operations through the Go backend ports and adapters.

---

## 🎨 UI/UX Design System & Rules

The application must feel state-of-the-art, premium, and professional.

### 1. Color Palette & Dark Mode
- **Background**: Deep obsidian/charcoal dark theme (`#0d1117` or `#0b0f19`).
- **Cards & Elements**: Slightly lighter gray/navy surfaces (`#161b22`, `#1f2937`) with subtle semi-transparent borders.
- **Accents**: Neon green (`#238636` / `#10b981`) for positive actions; Cyan/Blue (`#58a6ff` / `#00d2ff`) for interactive components, tabs, and loading spinners.
- **Text**: Off-white (`#c9d1d9`) for body text, pure white (`#ffffff`) for titles, and muted gray (`#8b949e`) for dates, tags, and secondary descriptions.

### 2. Typography & Micro-animations
- Sans-serif modern typography via Google Fonts (**Inter** or **Outfit**).
- Smooth transitions on interactive elements (hover effects, repo chip toggles, row selections).
- Smooth custom loading spinners during GitHub API fetches.

### 3. Report Previewer & Printing Usability
- Live markdown / printable preview.
- **"Copiar para Área de Transferência"** (Copy to Clipboard) with feedback animation.
- **"Imprimir Relatório"** / Export PDF.
- `@media print` query in CSS to hide navigation bars, sidebar tabs, settings panels, and control buttons during printing.

---

## ⚡ Performance, Bundle Size & Troubleshooting Rules

### 1. Dynamic Imports
- Large dependencies (e.g. `html2pdf.js`) must **always** use dynamic imports (`await import(...)`) to maintain small initial bundle sizes.

### 2. Handling `@current_problems`
- Analyze and fix compiler/linter issues immediately and verify builds pass.
- Create or update troubleshooting skills under `.agents/skills/` when solving recurring environmental or typing issues.

---

## 🚀 GitHub Versioning & Release Workflow Rules

1. **Version Declaration**: Ensure version is set in [wails.json](file:///c:/Users/hades/Documents/grelatorio/wails.json).
2. **NSIS Installer**: Compile via `wails build -nsis` (appending NSIS to `$env:Path` if needed).
3. **Tagging & Release**: Tag release (e.g. `git tag v1.0.0`), push, and attach `relatorio-amd64-installer.exe` to GitHub Releases.
