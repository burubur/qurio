# Missing Implementation Report: 2026-01-02 PRD vs Codebase

**Date:** 2026-01-02
**Status:** Deep Dive Gap Analysis
**Reference:** `docs/2026-01-02-prd.md`

## 1. Contextual Embeddings (Critical)
**File:** `apps/backend/internal/worker/result_consumer.go`

*   **Missing "Source Name" & "Path":**
    *   **Requirement:** FR-06 format: `Source: <Source Name>` and `Path: <Breadcrumbs>`.
    *   **Current State:** The code uses `URL` in place of `Path`. `Source Name` is completely missing from the Go `Source` struct and database schema.
    *   **Impact:** Context header is incomplete. "Source" usually refers to the repository or site name (e.g., "React Docs"), which is distinct from the base URL.
    *   **Fix:**
        1.  Add `Name` column to `sources` table and Go struct.
        2.  Update `result_consumer.go` to fetch Source Name by ID.
        3.  Update worker to extract Breadcrumbs (if possible) or derive `Path` from URL path segments.

## 2. Ingestion Logic
**File:** `apps/backend/internal/text/chunker.go` & `apps/ingestion-worker/handlers/web.py`

*   **Missing API Classification:**
    *   **Requirement:** FR-02 requires classifying chunks as `api`.
    *   **Current State:** Logic missing in `ChunkMarkdown`.
    *   **Fix:** Map `graphql`, `proto`, `thrift`, `openapi`, `swagger` to `ChunkTypeAPI`.

*   **Missing Breadcrumbs Extraction:**
    *   **Requirement:** `Path: <Breadcrumbs>` for embedding.
    *   **Current State:** Python worker extracts `title` but not `breadcrumbs`.
    *   **Fix:** Update `web.py` to extract breadcrumbs (e.g., from generic schema.org metadata or URL path analysis) and include in payload.

## 3. Backend / Storage
**File:** `apps/backend/internal/adapter/weaviate/store.go`

*   **Missing Sorting in Fetch (FR-11):**
    *   **Current State:** `GetChunksByURL` returns unsorted chunks.
    *   **Fix:** Sort by `chunkIndex` (integer) ascending.

*   **Weak Filtering (FR-10):**
    *   **Current State:** Only supports exact string match.
    *   **Fix:** While PRD doesn't explicitly demand array/OR logic, it's safer to ensure the filter map handling is robust. Current implementation is acceptable for the strict MVP but brittle.

## 4. MCP Tools Experience (AX)
**File:** `apps/backend/features/mcp/handler.go`

*   **Tool Descriptions & Guide:**
    *   **Requirement:** FR-13 & Section 2.4.
    *   **Current State:** Generic descriptions.
    *   **Fix:** Copy-paste the "Strategy Guide" and "Examples" from PRD.

*   **Schema Mismatches:**
    *   Argument name: `filters` (Code) vs `filter` (PRD).
    *   Enum validation missing for `type`.
    *   Search results missing `URL`.
    *   Fetch page output: Has decorative headers that break Markdown.