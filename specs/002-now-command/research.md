# Research: ua now Command

**Feature**: 002-now-command | **Date**: 2026-03-04

## R1: Temporal Filtering Logic

**Decision**: Implementing custom threshold-based date parsing in standard Go, avoiding external date libraries.

**Rationale**:
- The project aims for zero unnecessary dependencies. Go's `time` package is highly robust.
- The logic requires comparing the system's `time.Now()` against the existing `schedule.Event` `Start` and `End` properties.
- **The 30-minute rule:**
  - If `Now < Start`: The class is in the future. Select it as "up next" if it's the first one found.
  - If `Start <= Now < End`: The class is currently ongoing.
    - Calculate `TimeLeft = End - Now`.
    - If `TimeLeft <= 30m`: Skip this class and select the chronologically *next* available class for the day, assuming the student is mostly done and wants to know where to go next.
    - If `TimeLeft > 30m`: Show the ongoing class.

**Alternative Considered**:
- Only showing the "current" class strictly within the time bounds. Rejected because once a student is near the end of a class, they almost exclusively care about what comes next.

## R2: Data Fetching and Performance

**Decision**: Reuse the `ScheduleService.GetScheduleForWeek` adapter but filter locally.

**Rationale**:
- The UACloud API endpoints provide data in week-long blocks minimum. Asking for a single exact day does not reduce network payload meaningfully due to UA's backend structure.
- Fetching the week block leverages the existing local `.ua-cli/cache/schedule.json`. If the user has fetched their schedule recently, `ua now` will execute locally in ~5ms.
- To handle the filter: we extract `e.Start.Date()` and compare it to `time.Now().Date()`. We only process events where the year, month, and day match precisely.

## R3: UI Output Format (Lipgloss)

**Decision**: Display a styled "Card" instead of the standard ascii table for TTY sessions.

**Rationale**:
- A table is great for multi-row data (`ua schedule`), but terrible UX for a single emphasized item.
- Using `lipgloss.RoundedBorder()` with color coding (Green for "Ongoing", Blue for "Upcoming") provides a modern, glanceable widget perfectly aligned with the project's "terminal-first UX" philosophy.
