---
name: Troubleshooting html2pdf Page Breaks
description: Guidelines to prevent content (like table rows, signatures, or sections) from being cut in half across pages when generating PDFs with html2pdf.js.
---

# Troubleshooting html2pdf Page Breaks

This skill documents how to resolve the PDF generation bug where table rows (`tr`), signature sections, or indicator cards are cut in half across page transitions.

## 🔴 The Problem
When converting a long HTML document or table into a multi-page PDF using `html2pdf.js`, the rendering engine splits content at strict pixel heights. If page-break rules are not configured, text, table rows, and borders get divided horizontally across the bottom of one page and the top of the next.

## 🟢 The Solution
Combine both CSS properties and `html2pdf.js` page break configurations:

### 1. CSS Page Break Avoidance
Add `page-break-inside: avoid;` and `break-inside: avoid;` to elements that should never be split across pages (such as table rows or specific grid panels):
```css
.pdf-table tr {
  page-break-inside: avoid;
  break-inside: avoid;
}

.pdf-signatures {
  page-break-inside: avoid;
  break-inside: avoid;
}
```

### 2. Configure `html2pdf` Options
Pass a `pagebreak` parameter inside the `html2pdf.js` configuration object. Set the `mode` to recognize both CSS classes and automatic avoidance, and explicitly include selectors to bypass in the `avoid` list:
```javascript
const opt = {
  margin: 0,
  filename: 'report.pdf',
  jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
  pagebreak: { 
    mode: ['avoid-all', 'css'], 
    avoid: ['tr', '.pdf-signatures', '.pdf-checklist-block', '.pdf-indicators-block'] 
  }
};
```
This forces the renderer to shift the entire row or block to the next page instead of clipping it in half.
