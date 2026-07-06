---
name: Troubleshooting Activity Hours Calculation
description: Logic template to calculate elapsed hours for daily developer activities (commits/PRs) clamped inside standard working hours (7h to 17h).
---

# Troubleshooting Activity Hours Calculation

This skill documents how to calculate realistic hours for a list of daily developer activities based on their timestamps, clamped to a standard 10-hour working day (07:00 to 17:00).

## 🔴 The Problem
Using range-based division (e.g. dividing total days * 10h by the total number of commits) yields static, unrealistic values (e.g. `1.7h` for every commit), which doesn't reflect actual time spent per task or the chronology of the developer's work.

## 🟢 The Solution
Implement chronological subtraction between commit times on the same day, clamped to a 7:00 - 17:00 workday interval:

1. **Group** activities by their local date.
2. **Sort** activities of each day chronologically.
3. For the **first** activity of the day, compute the difference between the commit hour (clamped between 7.0 and 17.0) and the workday start (7.0).
4. For **subsequent** activities of the day, compute the difference between the current commit hour (clamped) and the previous commit hour (clamped).
5. Enforce a **minimum** block (e.g., `0.5` hours/30 minutes) to ensure small tasks are accounted for.

### Implementation Pattern (JS)
```javascript
// Group activities by date
const dateGroups = {};
selectedActivities.forEach(act => {
  const d = new Date(act.date);
  const dateKey = d.toISOString().split('T')[0];
  if (!dateGroups[dateKey]) dateGroups[dateKey] = [];
  dateGroups[dateKey].push(act);
});

// Calculate hours per group
Object.values(dateGroups).forEach(group => {
  group.sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
  
  let prevClampedHour = 7.0; // 07:00 start
  
  group.forEach(act => {
    const d = new Date(act.date);
    const hour = d.getHours() + d.getMinutes() / 60;
    const clampedHour = Math.min(17.0, Math.max(7.0, hour));
    
    let diff = clampedHour - prevClampedHour;
    if (diff < 0.5) diff = 0.5; // Minimum time
    
    act.hours = parseFloat(diff.toFixed(1));
    prevClampedHour = clampedHour;
  });
});
```
This yields accurate, chronological billing records for each item in the report.
