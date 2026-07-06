# Rules for GitHub Report Generator Development

Here are the guidelines, architecture, and UI/UX design rules for the GitHub monthly commit report generator application.

---

## 🏛️ Project Architecture

The application is built using **Wails (v2)**, which binds a **Go backend** with a **Svelte + Vite frontend**.

### 1. Go Backend (`app.go`, `main.go`)
- **Config Management**: Persist configuration (GitHub Token, Username, client details) locally in the user config directory (`%APPDATA%/github-report-generator/config.json`) using `LoadConfig` and `SaveConfig`.
- **API Fetching**: Query GitHub Search APIs (`/search/commits` and `/search/issues` for Pull Requests) to load user activity for a specific date range.
- **Wails Bindings**: Expose backend operations through Wails binding functions on the `App` struct. Never make direct fetch calls to GitHub API from the frontend; always channel them through the Go backend to ensure secure header handling and avoid CORS issues.

### 2. Svelte Frontend (`frontend/src/`)
- **State Management**: Keep user configurations and fetched activities in reactive stores or high-level variables in `App.svelte` to pass down to subcomponents.
- **Component Breakdown**:
  - [Settings.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/Settings.svelte): Manage user token, username, rates, and client details.
  - [ActivityFetch.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/ActivityFetch.svelte): Handle date selections, execute fetching, and display/toggle selected items.
  - [ReportBuilder.svelte](file:///c:/Users/eduar/Documents/relatorio/frontend/src/components/ReportBuilder.svelte): Preview, group, and format selected activities into a markdown/text report for the user's boss.

---

## 🎨 UI/UX Design System & Rules

The application must feel state-of-the-art, premium, and professional.

### 1. Color Palette & Dark Mode
- **Background**: Deep obsidian/charcoal dark theme (`#0d1117` or `#0b0f19`) to align with modern developer tools (GitHub-esque).
- **Cards & Elements**: Slightly lighter gray/navy surfaces (`#161b22`, `#1f2937`) with subtle semi-transparent borders.
- **Accents**: Neon green (`#238636` / `#10b981`) for positive actions and success states; Cyan/Blue (`#58a6ff` / `#00d2ff`) for interactive components, tabs, and loading spinners.
- **Text**: Off-white (`#c9d1d9`) for body text, pure white (`#ffffff`) for titles, and muted gray (`#8b949e`) for dates, tags, and secondary descriptions.

### 2. Typography
- Use modern sans-serif typography (e.g., **Inter** or **Outfit**) via Google Fonts, rather than standard system defaults.
- Maintain clear vertical hierarchy with appropriate font weights and line heights.

### 3. Micro-animations & Transitions
- Add smooth transitions on interactive components (e.g., hover effects on buttons, checkbox selections, and list items).
- Implement a custom, smooth loading spinner or progress bar during GitHub API fetches.
- Include subtle transitions (`fade` / `slide` from Svelte) when switching between dashboard tabs or expanding repository groups.

### 4. Report Previewer & Printing Usability
- The report generation screen must show a live markdown preview.
- Include a **"Copiar para Área de Transferência"** (Copy to Clipboard) button that triggers a brief success animation.
- Include a **"Imprimir Relatório"** (Print Report) button that triggers native print (`window.print()`).
- **Critical CSS Print Rules**: Use a `@media print` query in the CSS to hide navigation bars, sidebar tabs, settings panels, and control buttons. When printing, the page should show *only* the formatted report text on a clean, light-colored background.

---

## ⚠️ Security & Reliability Rules
- **Token Security**: Never log or print the GitHub Personal Access Token in the app logs or UI console.
- **Rate Limit Handling**: Implement graceful error feedback on the UI when GitHub API rate limits are hit or credentials are invalid.
- **Robust Parsing**: Handle multi-line commit messages gracefully (usually the first line is the summary, and the rest is body description).

---

## ⚡ Performance, Bundle Size & Troubleshooting Rules

### 1. Bundle Chunk Size Control
- **Never let main bundle chunks grow excessively**: Large dependencies (e.g., pdf generators, charting libraries) must **never** be statically imported at the top of a Svelte component or JavaScript file.
- **Always use dynamic imports** (`await import(...)`) inside event handlers or specific lifecycle hooks to force Vite/Rollup code-splitting, keeping the initial payload small.

### 2. Handling `@current_problems` & Automated Skill Generation
- When the user mentions `@current_problems`, prioritize analyzing the compilation/linting issues immediately, apply correct code fixes, and verify that the build succeeds.
- At the end of the solution, **always create a new troubleshooting skill** in the workspace customizations root under `.agents/skills/troubleshooting_<category>/SKILL.md` (e.g., `troubleshooting_window_types`, `troubleshooting_date_arithmetic`).
- The `SKILL.md` file must explain:
  - The compiler/linter error pattern.
  - The root cause.
  - The exact resolution pattern.
- Divide skills into specific categories to save tokens during searches by future agents (do not create one large generic troubleshooting skill).

---

## 🚀 GitHub Versioning & Release Workflow Rules

When compiling new versions or publishing releases of this application, always adhere to the following workflow:

1. **Version Control**:
   - Ensure the version number is declared in [wails.json](file:///c:/Users/eduar/Documents/relatorio/wails.json) via the `"version"` field (e.g., `"version": "1.0.0"`). If the field is missing, initialize it.
   - Commit all pending code changes and push to the remote repository before tagging.

2. **NSIS Installer Compilation**:
   - Run the production Wails build using the Windows installer flag: `wails build -nsis`.
   - If the system does not recognize the `makensis` command, dynamically append the NSIS directory to the environment `PATH` in the execution shell (e.g., `$env:Path += ";C:\Program Files (x86)\NSIS"` in PowerShell).
   - Ensure both the executable [relatorio.exe](file:///c:/Users/eduar/Documents/relatorio/build/bin/relatorio.exe) and the installer [relatorio-amd64-installer.exe](file:///c:/Users/eduar/Documents/relatorio/build/bin/relatorio-amd64-installer.exe) are generated in the `build/bin/` folder.

3. **Tagging & GitHub Release**:
   - Create a corresponding git tag matching the version (e.g., `git tag v1.0.0`).
   - Push the tag to remote: `git push origin v1.0.0`.
   - Guide or assist the user in creating a GitHub Release for that tag on GitHub and upload the compiled `relatorio-amd64-installer.exe` as a release asset.

