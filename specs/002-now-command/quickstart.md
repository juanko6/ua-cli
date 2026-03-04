# Quickstart: ua now Command

**Feature**: 002-now-command | **Date**: 2026-03-04

## Description

The `ua now` command provides instant visibility into your schedule, bypassing the overhead of looking at a full week table. 

It checks the current system time and finds the single most relevant class for right now or coming up next today.

## Prerequisites

- You must have logged in (`ua login`) at least once so the API can fetch the schedule.
- The system time must be accurate.

## Usage

### Default (TUI Card)

```bash
ua now
```

**Expected Outputs:**
- **Scenario A (Before class):** A styled card indicating "● PRÓXIMA", the class name, time, and room.
- **Scenario B (In class, > 30m left):** A styled card indicating "● AHORA MISMO", the class name, time, and room.
- **Scenario C (In class, < 30m left):** A styled card indicating "● PRÓXIMA" for the *next* chronological class.
- **Scenario D (End of day):** "🎉 No tienes más clases por hoy. ¡A descansar!"

### JSON Output (Robots/Scripts)

```bash
ua now --json
```

**Expected Output:**
- An array containing a single `schedule.Event` JSON object, or an empty object `{}` if no classes remain.
