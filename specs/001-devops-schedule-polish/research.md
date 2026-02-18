# Research: DevOps Infrastructure & Schedule Polish

**Feature**: 001-devops-schedule-polish  
**Date**: 2026-02-17

## R1: CI/CD Tool Selection

**Decision**: GitHub Actions  
**Rationale**: Native integration with GitHub (no external service). Free tier (2000 min/month) sufficient. Supports Go natively with `setup-go`.  
**Alternatives**: GitLab CI (requires migration), CircleCI (extra config), Drone (self-hosted).

## R2: Release Tool Selection

**Decision**: GoReleaser  
**Rationale**: De facto standard for Go CLI releases. Cross-compiles, generates checksums, creates GitHub Releases. Single YAML config.  
**Alternatives**: `go build` scripts (manual), `nfpm` (overkill for CLI).

## R3: Branching Strategy

**Decision**: GitHub Flow (main + feature branches)  
**Rationale**: Simple, one production branch. PRs required. Fits small team / solo dev.  
**Alternatives**: Git Flow (too complex for MVP), Trunk-Based (no PR review).

## R4: UACloud Event Color Mapping

**Decision**: Map color hex codes to event types. Need to inspect real API response to get actual color values.  
**Rationale**: UACloud uses FullCalendar which assigns colors per event type. We inspect the JSON response and map them.  
**Action**: Log unique colors from real API response to build mapping table.

## R5: Week Filtering Strategy

**Decision**: Client-side filter in `ScheduleService.GetScheduleForWeek`. Compare event `Start` against Monday 00:00 – Sunday 23:59 of target week.  
**Rationale**: The API endpoint accepts `start`/`end` params but returns events beyond the range. Client-side filtering is safest.  
**Alternatives**: Trust API params only (unreliable based on testing).

## R6: Deduplication Strategy

**Decision**: Deduplicate by composite key: `Title + Start + End`. Keep first occurrence.  
**Rationale**: UACloud returns identical events multiple times (likely one per group/subgroup). Students only care about unique slots.  
**Alternatives**: Deduplicate by ID (unreliable, IDs may be empty).
