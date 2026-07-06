---
name: Troubleshooting Date Arithmetic Types
description: Standard resolutions for TypeScript/linter errors when performing arithmetic operations on Date objects.
---

# Troubleshooting Date Arithmetic Types

This skill explains how to resolve the TypeScript compiler/linter error:
`The left-hand side/right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.`

## 🔴 The Problem
Subtracting two Date objects directly to find the difference in milliseconds:
```javascript
const expDate = new Date(isoStr);
const today = new Date();
const diffTime = expDate - today;
```
causes TypeScript linter warnings because it expects numeric values for mathematical operators like `-`.

## 🟢 The Solution
Use the `.getTime()` method on both Date instances. This converts them to numeric epoch milliseconds, satisfying the type system:
```javascript
const expDate = new Date(isoStr);
const today = new Date();
const diffTime = expDate.getTime() - today.getTime();
```
This is type-safe and avoids compilation warnings.
