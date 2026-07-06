---
name: Troubleshooting html2pdf Margins Cutoff
description: Resolution pattern for PDF horizontal shifts/clipping errors when using html2pdf.js with pre-sized A4 templates.
---

# Troubleshooting html2pdf Margins Cutoff

This skill explains how to solve the layout alignment bug where the generated PDF is shifted to the right, showing large white margins on the left and cutting off content on the right.

## 🔴 The Problem
When using `html2pdf.js` to convert an HTML element to a PDF, if the target HTML element is already styled with a fixed page width (like A4 size `width: 210mm;`) and has internal margins (via CSS `padding`), setting a non-zero margin in the `html2pdf` options:
```javascript
const opt = {
  margin: 15, // Adds 15mm extra margin to the layout
  jsPDF: { format: 'a4' }
};
```
will shift the rendered image to the right by `15mm`. Since the A4 width limit is `210mm`, the resulting print area will be `210 + 15 = 225mm`, causing the rightmost `15mm` of the page content to be cut off.

## 🟢 The Solution
Set the `margin` parameter inside `html2pdf` options to `0`. This renders the fixed-width A4 element exactly 1:1 onto the A4 PDF page container without offset shifts, preserving the container's native CSS paddings:
```javascript
const opt = {
  margin: 0, // Align exactly 1:1 on the A4 page
  jsPDF: { format: 'a4', orientation: 'portrait' }
};
```
