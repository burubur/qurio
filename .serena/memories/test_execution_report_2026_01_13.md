# Test Suite Execution Report

**Date:** 2026-01-13
**Scope:** Full suite (Unit + Integration) for Backend, Frontend, Ingestion Worker.

## Results

### Backend (Go)
- **Command:** `go test -v ./apps/backend/...`
- **Status:** ✅ Passed
- **Coverage:**
  - Unit tests for all packages.
  - Integration tests using `testcontainers-go` (Postgres, Weaviate, NSQ).
  - Verified logic for Source, Job, Worker, MCP, and internal adapters.

### Frontend (Vue/Vitest)
- **Command:** `npm run test -- --run`
- **Status:** ✅ Passed
- **Coverage:**
  - Component tests (SourceForm, SourceList, Settings, UI components).
  - Store tests (Pinia stores for sources, settings, jobs, stats).
  - Updated tests for `SourceForm` (Source Name Refactor).

### Ingestion Worker (Python)
- **Command:** `PYTHONPATH=. pytest` (in venv)
- **Status:** ✅ Passed
- **Coverage:**
  - Handlers (File, Web).
  - Configuration & Logging.
  - Integration tests for NSQ and Metadata extraction.

## Notes
- Frontend tests required updates to match the "Source Naming" refactor (mocking API keys, targeting correct inputs).
- Ingestion Worker tests required `PYTHONPATH=.` to resolve local modules correctly.
