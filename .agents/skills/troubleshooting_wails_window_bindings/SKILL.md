---
name: Troubleshooting Wails Window Bindings
description: Guidelines for resolving linter/type errors related to Wails 'window.go' bindings in Svelte/JS.
---

# Troubleshooting Wails Window Bindings

This skill documents how to resolve the TypeScript compiler/linter errors:
- `Property 'go' does not exist on type 'Window & typeof globalThis'.`
- `Parameter 'argX' implicitly has an 'any' type.` in `wailsjs/go/main/App.js` with `// @ts-check`.

## 🔴 The Problem
When coding in JavaScript or TypeScript within Svelte components or Wails generated JS wrappers, referencing `window.go` directly or leaving function parameters unannotated under `// @ts-check`:
```javascript
const isWailsAvailable = typeof window !== 'undefined' && window.go !== undefined;
```
triggers type checking errors because `go` is dynamically injected by Wails at runtime and does not exist in standard browser `Window` definitions.

## 🟢 The Solution
1. Cast the `window` reference dynamically to `any` using JSDoc:
```javascript
const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined;
```

2. In `wailsjs/go/main/App.js`, annotate wrapper functions with JSDoc and cast `window` to `any`:
```javascript
import {domain} from '../models';

/**
 * @param {string} arg1
 * @param {string} arg2
 * @param {string} [arg3]
 * @returns {Promise<Array<domain.GithubActivity>>}
 */
export function FetchGithubActivity(arg1, arg2, arg3 = '') {
  return (/** @type {any} */ (window))['go']['main']['App']['FetchGithubActivity'](arg1, arg2, arg3 || '');
}
```
