# Implementation Plan - MVP Part 2.2: Fixes & Frontend Design System

**Scope:**
1.  **Frontend:** Integrate shadcn-vue, Tailwind CSS, and refactor existing components to use the new design system.
2.  **Backend:** Fix logging inconsistencies (slog), strictly enforce JSON error envelopes, and ensure timeout configuration for external clients.

**References:**
- `docs/2025-12-22-bugs-inconsistencies.md`
- `apps/frontend` (Codebase Investigator analysis)
- `apps/backend` (Codebase Investigator analysis)

---

### Task 1: Install Tailwind CSS & PostCSS

**Files:**
- Modify: `apps/frontend/package.json`
- Create: `apps/frontend/postcss.config.js`
- Create: `apps/frontend/tailwind.config.js`
- Modify: `apps/frontend/src/style.css`

**Requirements:**
- **Acceptance Criteria**
  1. `npm install` runs successfully with new devDependencies.
  2. Tailwind directives (`@tailwind base;` etc.) are present in `style.css`.
  3. Tailwind config file exists.

- **Functional Requirements**
  1. Enable utility-first CSS framework (Tailwind) for styling.

- **Non-Functional Requirements**
  None for this task.

- **Test Coverage**
  - [Manual] Run `npm run dev` and verify no build errors.

**Step 1: Write failing test**
*Skipped (Infrastructure Setup)* - Verification via build command.

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Run shell command in `apps/frontend`:
    ```bash
    npm install -D tailwindcss@3 autoprefixer postcss
    npx tailwindcss init -p
    ```
2.  Update `tailwind.config.js` content matches `src/**/*.{vue,js,ts,jsx,tsx}`.
3.  Add directives to `src/style.css`.

**Step 4: Verify test passes**
Run: `cd apps/frontend && npm run build`
Expected: Success.

---

### Task 2: Configure Path Aliases

**Files:**
- Modify: `apps/frontend/tsconfig.app.json`
- Modify: `apps/frontend/vite.config.ts`

**Requirements:**
- **Acceptance Criteria**
  1. Importing from `@/components` works.

- **Functional Requirements**
  1. Map `@` to `./src`.

- **Non-Functional Requirements**
  Standardize imports.

- **Test Coverage**
  - [Manual] Build verification.

**Step 1: Write failing test**
*Skipped (Configuration)*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Update `tsconfig.app.json`: `compilerOptions.paths` = `{"@/*": ["./src/*"]}`.
2.  Update `vite.config.ts`: `resolve.alias` = `{"@": path.resolve(__dirname, "./src")}`. (Import `path` module).

**Step 4: Verify test passes**
Run: `cd apps/frontend && npm run build`
Expected: Success.

---

### Task 3: Initialize shadcn-vue

**Files:**
- Create: `apps/frontend/components.json`
- Create: `apps/frontend/src/lib/utils.ts` (or `utils/cn.ts` depending on config)

**Requirements:**
- **Acceptance Criteria**
  1. `components.json` exists with correct configuration.
  2. `cn` utility function exists.

- **Functional Requirements**
  1. Bootstrap shadcn-vue configuration.

- **Non-Functional Requirements**
  Use `slate` as base color (default).

- **Test Coverage**
  - [Manual] File existence check.

**Step 1: Write failing test**
*Skipped*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Run shell command in `apps/frontend`:
    ```bash
    npm install -D typescript
    npx shadcn-vue@latest init -d
    ```
    (Using `-d` for defaults to avoid interactivity).

**Step 4: Verify test passes**
Run: `ls apps/frontend/components.json`
Expected: File found.

---

### Task 4: Add Base Components

**Files:**
- Create: `apps/frontend/src/components/ui/button/Button.vue`
- Create: `apps/frontend/src/components/ui/input/Input.vue`
- Create: `apps/frontend/src/components/ui/badge/Badge.vue`
- Create: `apps/frontend/src/components/ui/card/Card.vue`
- Create: `apps/frontend/src/components/ui/select/Select.vue`
- Modify: `apps/frontend/src/components/ui/form` (if needed by others, but starting with basics)

**Requirements:**
- **Acceptance Criteria**
  1. UI components exist in `src/components/ui`.

- **Functional Requirements**
  1. Install Button, Input, Badge, Card, Select.

- **Non-Functional Requirements**
  None.

- **Test Coverage**
  - [Manual] File existence check.

**Step 1: Write failing test**
*Skipped*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Run shell command in `apps/frontend`:
    ```bash
    npx shadcn-vue@latest add button input badge card select -y
    ```

**Step 4: Verify test passes**
Run: `ls apps/frontend/src/components/ui/button/Button.vue`
Expected: File found.

---

### Task 5: Backend Logging (slog) - Fix Violations

**Files:**
- Modify: `apps/backend/features/source/handler.go`
- Modify: `apps/backend/internal/worker/ingest.go`
- Test: `apps/backend/features/source/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `fmt.Printf` and `log.Printf` replaced with `slog.Info` or `slog.Error`.
  2. Structured logging used (key-value pairs).

- **Functional Requirements**
  1. Log errors with "error" key.

- **Non-Functional Requirements**
  Compliance with Technical Constitution.

- **Test Coverage**
  - [Unit] Verify handlers still function (logging changes shouldn't break logic).

**Step 1: Write failing test**
*Skipped (Refactoring logging)* - rely on existing tests.

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Import `log/slog`.
2.  Replace `fmt.Printf("Error: %v", err)` with `slog.Error("operation failed", "error", err)`.
3.  Replace `log.Printf(...)` with `slog.Info(...)`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/source/...`
Expected: PASS.

---

### Task 6: Backend Missing Logging

**Files:**
- Modify: `apps/backend/features/mcp/handler.go`
- Modify: `apps/backend/internal/settings/handler.go`

**Requirements:**
- **Acceptance Criteria**
  1. Handlers log "request received" at start.
  2. Handlers log "request completed" or "request failed" at end.

- **Functional Requirements**
  1. Include `method` and `path` in start logs.

- **Non-Functional Requirements**
  Compliance with Technical Constitution.

- **Test Coverage**
  - [Unit] Verify handlers execution.

**Step 1: Write failing test**
*Skipped*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Add `slog.Info("request received", ...)` at handler start.
2.  Add `slog.Info("request completed", ...)` before successful return.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/...`
Expected: PASS.

---

### Task 7: JSON Error Handling & Correlation ID

**Files:**
- Modify: `apps/backend/internal/settings/handler.go`
- Modify: `apps/backend/features/source/handler.go`

**Requirements:**
- **Acceptance Criteria**
  1. No `http.Error()` calls (which return text/plain).
  2. Error responses use JSON format: `{ "error": { "code": "...", "message": "..." }, "correlationId": "..." }`.
  3. Correlation ID generated (or retrieved from context) and included.

- **Functional Requirements**
  1. Define a helper function `writeError(w, code, message, correlationID)` or similar if not exists (check shared utils).

- **Non-Functional Requirements**
  Compliance with Technical Constitution.

- **Test Coverage**
  - [Unit] Test error scenarios return JSON.

**Step 1: Write failing test**
Modify `settings/service_test.go` or equivalent to check error content-type is `application/json`.

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/settings/...`
Expected: Fail (currently returns text).

**Step 3: Write minimal implementation**
1.  Create/Use a standard error response struct.
2.  Generate UUID for correlation ID if missing.
3.  Replace `http.Error` with `json.NewEncoder(w).Encode(errorResponse)`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/settings/...`
Expected: PASS.

---

### Task 8: Resource Management (Timeouts)

**Files:**
- Modify: `apps/backend/internal/adapter/docling/client.go`

**Requirements:**
- **Acceptance Criteria**
  1. `http.Client` is initialized with a `Timeout`.

- **Functional Requirements**
  1. Set timeout to 30s (default safe value).

- **Non-Functional Requirements**
  Prevent resource exhaustion.

- **Test Coverage**
  - [Manual] Code review or integration test with delay (complex). Manual verification preferred for config change.

**Step 1: Write failing test**
*Skipped*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Change `&http.Client{}` to `&http.Client{Timeout: 30 * time.Second}`.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/adapter/docling/...`
Expected: PASS.

---

### Task 9: Frontend Refactor - SourceForm

**Files:**
- Modify: `apps/frontend/src/features/sources/SourceForm.vue`

**Requirements:**
- **Acceptance Criteria**
  1. Uses `<Button>` and `<Input>` from `@/components/ui`.
  2. No native `<button>` or `<input>`.

- **Functional Requirements**
  1. Same v-model binding behavior.

- **Test Coverage**
  - [Unit] `SourceForm.spec.ts`.

**Step 1: Write failing test**
Update `SourceForm.spec.ts` to look for shadcn component classes or props if applicable, or just ensure existing tests pass after refactor.

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Import `Button` and `Input`.
2.  Replace templates.

**Step 4: Verify test passes**
Run: `npm run test:unit apps/frontend/src/features/sources/SourceForm.spec.ts`
Expected: PASS.

---

### Task 10: Frontend Refactor - SourceList

**Files:**
- Modify: `apps/frontend/src/features/sources/SourceList.vue`
- Modify: `apps/frontend/src/components/ui/StatusBadge.vue` (Migrate this first or inline it)

**Requirements:**
- **Acceptance Criteria**
  1. `SourceList` uses `<Card>` for items.
  2. `StatusBadge` uses `<Badge>`.

**Step 1: Write failing test**
*Skipped*

**Step 2: Verify test fails**
*Skipped*

**Step 3: Write minimal implementation**
1.  Update `StatusBadge.vue` to wrap shadcn `<Badge>`.
2.  Update `SourceList.vue` to use `<Card>` structure (Header, Content, Footer).

**Step 4: Verify test passes**
Run: `npm run test:unit`
Expected: PASS.
