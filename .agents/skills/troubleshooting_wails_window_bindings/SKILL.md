---
name: Troubleshooting Wails Window Bindings
description: Guidelines for resolving linter/type errors related to Wails 'window.go' bindings in Svelte/JS.
---

# Troubleshooting Wails Window Bindings

This skill documents how to resolve the TypeScript compiler/linter error:
`Property 'go' does not exist on type 'Window & typeof globalThis'.`

## 🔴 The Problem
When coding in JavaScript or TypeScript within Svelte components in a Wails project, referencing `window.go` directly:
```javascript
const isWailsAvailable = typeof window !== 'undefined' && window.go !== undefined;
```
triggers type checking errors because `go` is dynamically injected by Wails at runtime and does not exist in standard browser `Window` definitions.

## 🟢 The Solution
Cast the `window` reference dynamically to `any` using JSDoc. This satisfies the linter/type checker without changing runtime behavior:
```javascript
const isWailsAvailable = typeof window !== 'undefined' && (/** @type {any} */(window)).go !== undefined;
```

For subsequent checks of nested properties:
```javascript
const isWailsAvailable = typeof window !== 'undefined' && 
  (/** @type {any} */(window)).go !== undefined && 
  (/** @type {any} */(window)).go.main !== undefined;
```
