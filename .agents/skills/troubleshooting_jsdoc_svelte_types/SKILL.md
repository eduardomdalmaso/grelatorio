---
name: Troubleshooting JSDoc Svelte Types
description: Standard guidelines to resolve implicit 'any' and state assignment type checking errors in Svelte components using JSDoc.
---

# Troubleshooting JSDoc Svelte Types

This skill outlines resolutions for implicit type errors in Svelte/JavaScript code:
- `Variable 'selectedActivities' implicitly has an 'any[]' type.`
- `Parameter 'view' implicitly has an 'any' type.`
- `Type 'number | null' is not assignable to type 'null'.`

## 🔴 The Problem
In JavaScript-based Svelte projects, type checking systems (like the VS Code Svelte extension) run static checks. 
If variables or function parameters do not have explicit types, or if a state variable is initialized with `null` or `[]` and later assigned a different type, the system complains about implicit types and assignment mismatches.

## 🟢 The Solution
Add JSDoc comments to declare explicit types:

### 1. Variables and Arrays
```javascript
/** @type {any[]} */
let selectedActivities = $state([]);
```

### 2. Null Initializations
```javascript
/** @type {number | null} */
let daysRemaining = $state(null);
```

### 3. Function Parameters
```javascript
/**
 * @param {string} view
 */
function setView(view) {
  currentView = view;
}
```
This enables static type verification while keeping the project in clean, pure JavaScript.
