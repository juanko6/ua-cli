# Tasks: Automatic Cookie Login

**Input**: Design documents from `/specs/001-auto-cookie-login/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Tests**: Not explicitly requested in spec. Test tasks included only for the proxy adapter (critical component).

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Exact file paths included in descriptions

---

## Phase 1: Setup

**Purpose**: Create new files and port interface needed by all stories

- [x] T001 [P] Create `CookieCapturer` port interface in `internal/domain/auth/capturer.go`
- [x] T002 [P] Create success HTML page constant in `internal/adapters/auth/proxy_page.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Implement the reverse proxy adapter — core infrastructure for US1

**⚠️ CRITICAL**: US1 cannot begin integration until this phase is complete

- [x] T003 Implement `ProxyAdapter` (reverse proxy cookie capturer) in `internal/adapters/auth/proxy.go` — includes: listener with port fallback, `httputil.ReverseProxy` with `ModifyResponse` to intercept `Set-Cookie` headers, CAS `service` URL rewriting, cookie channel, 2-minute timeout, graceful shutdown
- [x] T004 Write unit tests for `ProxyAdapter` in `internal/adapters/auth/proxy_test.go` — test cookie interception from mock Set-Cookie headers, timeout behavior, port fallback

**Checkpoint**: Proxy adapter ready — user story implementation can begin

---

## Phase 3: User Story 1 — One-Command Automatic Login (Priority: P1) 🎯 MVP

**Goal**: Running `ua login` opens the browser, user logs in via CAS, cookies are captured and stored automatically.

**Independent Test**: Run `ua login` → complete CAS login in browser → CLI shows "✓ Login successful" → `ua login --status` shows "✓ Session active"

### Implementation for User Story 1

- [x] T005 [US1] Modify `cmd/ua-cli/login.go` — refactor `promptCookie()` to use `ProxyAdapter.Capture()` as default flow when `--cookie` flag is NOT set. Wire: create proxy → call `Capture(ctx)` → pass cookie to `svc.StoreCookie()`. Update terminal messages with Lipgloss styling (spinner/progress indicator while waiting).
- [x] T006 [US1] Update browser opening in `cmd/ua-cli/login.go` — change URL from `https://cvnet.cpd.ua.es/uaCloud` to `http://localhost:<port>/uaCloud` (proxy URL) to route through the local proxy.
- [x] T007 [US1] Update success HTML page served by proxy in `internal/adapters/auth/proxy_page.go` — styled page with UA branding saying "✓ Login successful! You can close this tab."

**Checkpoint**: User Story 1 fully functional — `ua login` works end-to-end with automatic cookie capture

---

## Phase 4: User Story 2 — Manual Cookie Fallback (Priority: P2)

**Goal**: `ua login --cookie` continues to work as manual fallback for headless/restricted environments.

**Independent Test**: Run `ua login --cookie` → paste valid cookie → CLI shows "✓ Login successful"

### Implementation for User Story 2

- [x] T008 [US2] Verify and document manual fallback in `cmd/ua-cli/login.go` — ensure `--cookie` flag path is preserved unchanged. Add help text explaining when to use manual mode ("Use this if the automatic browser login doesn't work").

**Checkpoint**: Both automatic and manual login work independently

---

## Phase 5: User Story 3 — Login Timeout Protection (Priority: P3)

**Goal**: Login process auto-cancels after 2 minutes if user doesn't complete authentication.

**Independent Test**: Run `ua login` → do NOT complete browser login → CLI shows timeout error after 2 minutes and exits cleanly

### Implementation for User Story 3

- [x] T009 [US3] Implement timeout handling in `cmd/ua-cli/login.go` — wrap proxy `Capture()` call with `context.WithTimeout(ctx, 2*time.Minute)`. Display countdown or elapsed time with Lipgloss. On timeout: display clear error message, clean up proxy gracefully.
- [x] T010 [US3] Handle edge case: browser fails to open in `cmd/ua-cli/login.go` — if `OpenInBrowser()` returns error, display the proxy URL manually and instruct user to open it. Don't abort the flow.

**Checkpoint**: All user stories independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final cleanup and validation

- [x] T011 [P] Update `go.mod` / `go.sum` if new dependencies added — run `go mod tidy` in project root
- [x] T012 [P] Verify all existing tests pass — run `go test ./...` from project root
- [ ] T013 Run `quickstart.md` validation — follow the quickstart steps end-to-end to verify the feature works as documented
- [ ] T014 Commit all changes to branch `001-auto-cookie-login`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — T001 and T002 can start immediately and in parallel
- **Foundational (Phase 2)**: Depends on T001 (port interface) — BLOCKS user stories
- **US1 (Phase 3)**: Depends on Phase 2 completion (T003, T004)
- **US2 (Phase 4)**: No dependencies on other stories — can start after Phase 2 (or even earlier since it's mostly verification)
- **US3 (Phase 5)**: Depends on T005 (login.go refactored) — timeout wraps the proxy flow
- **Polish (Phase 6)**: Depends on all stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 2 — no cross-story dependencies
- **US2 (P2)**: Independent — existing functionality, just verify and document
- **US3 (P3)**: Depends on US1's T005 (needs proxy flow to wrap with timeout)

### Within Each User Story

- Port interface before adapter implementation
- Adapter before CLI wiring
- CLI wiring before UX polish

### Parallel Opportunities

- T001 and T002 can run in parallel (different files)
- T011 and T012 can run in parallel (different concerns)
- US2 (T008) can start in parallel with US1 (different code paths)

---

## Parallel Example: Phase 1

```bash
# Launch both setup tasks together:
Task T001: "Create CookieCapturer port in internal/domain/auth/capturer.go"
Task T002: "Create success HTML page in internal/adapters/auth/proxy_page.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001, T002)
2. Complete Phase 2: Foundational (T003, T004)
3. Complete Phase 3: User Story 1 (T005, T006, T007)
4. **STOP and VALIDATE**: Test `ua login` end-to-end
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Proxy adapter ready
2. Add US1 → `ua login` works automatically → **MVP!**
3. Add US2 → Manual fallback verified → Robust
4. Add US3 → Timeout protection → Production-ready
5. Polish → Tests, cleanup, commit

---

## Notes

- Total tasks: **14**
- Tasks per story: US1=3, US2=1, US3=2, Setup=2, Foundational=2, Polish=4
- Parallel opportunities: Phase 1 (2 tasks), Phase 6 (2 tasks), US2 alongside US1
- MVP scope: Phases 1-3 (T001–T007) = 7 tasks
- All task format validated: checkbox ✅, ID ✅, [P]/[Story] labels ✅, file paths ✅
