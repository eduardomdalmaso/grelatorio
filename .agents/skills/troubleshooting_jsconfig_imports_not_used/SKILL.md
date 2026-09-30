---
name: Troubleshooting jsconfig importsNotUsedAsValues
description: Resolution pattern for 'Option importsNotUsedAsValues has been removed. Use verbatimModuleSyntax instead' in jsconfig.json/tsconfig.json.
---

# Troubleshooting jsconfig.json: importsNotUsedAsValues Error

## 1. Error Pattern
In IDE problems / TypeScript compiler diagnostics:
```text
Option 'importsNotUsedAsValues' has been removed. Please remove it from your configuration.
  Use 'verbatimModuleSyntax' instead.
```

## 2. Root Cause
In TypeScript 5.0+, the compiler option `"importsNotUsedAsValues"` was deprecated and subsequently removed. Modern TypeScript configurations require replacing it with `"verbatimModuleSyntax": true` to control how type-only imports and exports are handled during module emission.

## 3. Resolution Pattern
In `jsconfig.json` or `tsconfig.json`, replace `"importsNotUsedAsValues"` with `"verbatimModuleSyntax": true`:

```json
{
  "compilerOptions": {
    "moduleResolution": "Node",
    "target": "ESNext",
    "module": "ESNext",
    "verbatimModuleSyntax": true,
    "isolatedModules": true,
    ...
  }
}
```
