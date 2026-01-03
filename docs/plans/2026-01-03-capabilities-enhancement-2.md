---
name: technical-constitution
description: Generates technical implementation plans and architectural strategies that enforce the Project Constitution.
---

# Implementation Plan - Capabilities Enhancement (Part 2)

**Date:** 2026-01-03
**Status:** Planned
**Feature:** Capabilities Enhancement (Contextual Embeddings, Advanced Ingestion, MCP Upgrade)
**Sequence:** 2
**Reference:** `docs/2026-01-02-missing-implementation.md`, `docs/2026-01-02-prd.md`

## 1. Requirements Analysis

### Scope
Implementation of missing critical features from the PRD: Contextual Embeddings (Source Name, Breadcrumbs), Advanced Ingestion (API detection), Retrieval Improvements (Sorting, Filtering), and MCP Tool Upgrades (`qurio_search` rename, `qurio_fetch_page`).

### Gap Analysis
*   **Contextual Embeddings:**
    *   **Missing Data:** `Source.Name` in DB and Struct.
    *   **Missing Data:** `Breadcrumbs` in ingestion payload.
    *   **Missing Logic:** Context string construction in `result_consumer.go` using these fields.
*   **Ingestion Logic:**
    *   **Missing Logic:** `ChunkTypeAPI` classification in `chunker.go`.
*   **Retrieval:**
    *   **Missing Logic:** Sorting by `chunkIndex` in `GetChunksByURL`.
    *   **Missing Logic:** Robust filtering in `Search`.
*   **MCP Tools:**
    *   **Incorrect Name:** `search` instead of `qurio_search`.
    *   **Missing Tool:** `qurio_fetch_page`.
    *   **Incomplete Schema:** Tool descriptions and filter arguments.

### Exclusions
*   Front-end updates (strictly backend/worker/mcp).
*   New crawler engine (using existing `crawl4ai` setup with enhancements).

## 2. Knowledge Enrichment

*   **Weaviate Sorting/Filtering:**
    *   Sort: Use `WithSort` builder with `path=["chunkIndex"]` and `order=Asc`.
    *   Filter: Use `WithWhere` with `Operator: Equal` for strict filtering of `type` and `language`.
*   **MCP Schema:**
    *   Follow `Model Context Protocol` standards for tool definition (name, description, inputSchema).
*   **Breadcrumbs:**
    *   Derive from URL path segments in `web.py` as a robust fallback (e.g., `docs/core/quickstart` -> `core > quickstart`).

## 3. Implementation Tasks

### Task 1: Add Source Name to Domain & DB (Backend)

**Files:**
-   Create: `apps/backend/migrations/000011_add_source_name.up.sql`
-   Create: `apps/backend/migrations/000011_add_source_name.down.sql`
-   Modify: `apps/backend/features/source/source.go`
-   Modify: `apps/backend/features/source/repo.go`
-   Test: `apps/backend/features/source/repo_test.go` (if exists) or `apps/backend/features/source/source_test.go`

**Requirements:**
-   **AC-01:** `sources` table has a `name` column (TEXT, NULLABLE or DEFAULT '').
-   **AC-02:** `Source` struct has `Name` field.
-   **AC-03:** Repository `Save`, `Get`, `List` methods handle the `Name` field.
-   **FR-01:** Allow users (or system) to set a human-readable name for a source (e.g., "React Docs").

**Step 1: Write failing test (if applicable for repo)**
*Add a test case in `repo_test.go` that attempts to save and retrieve a Source with a Name.*

**Step 2: Write minimal implementation**
1.  Create migration files.
    ```sql
    ALTER TABLE sources ADD COLUMN name TEXT DEFAULT '';
    ```
2.  Update `Source` struct.
3.  Update `PostgresRepo` methods to scan/insert `name`.

**Step 3: Verify test passes**
*Run repo tests.*

### Task 2: Enhance Ingestion Worker (Breadcrumbs & Path)

**Files:**
-   Modify: `apps/ingestion-worker/handlers/web.py`
-   Test: `apps/ingestion-worker/tests/test_handlers.py`

**Requirements:**
-   **AC-01:** Worker payload includes a `path` field derived from the URL (or scraped breadcrumbs).
-   **AC-02:** `path` format is `segment > segment` (e.g., "features > mcp").
-   **AC-03:** `title` extraction is robust (fallback to URL segment if regex fails).

**Step 1: Write failing test**
*In `test_handlers.py`, assert that the returned result from `handle_web_task` contains a `path` field.*

**Step 2: Write minimal implementation**
1.  In `handle_web_task`, parse `result.url`.
2.  Split path by `/`, filter empty/common segments, join with ` > `.
3.  Add `path` to the returned dictionary.

**Step 3: Verify test passes**
*Run `pytest apps/ingestion-worker/tests/test_handlers.py`.*

### Task 3: Enhance Backend Ingestion (API Classification & Embeddings)

**Files:**
-   Modify: `apps/backend/internal/text/chunker.go`
-   Modify: `apps/backend/internal/worker/result_consumer.go`
-   Test: `apps/backend/internal/text/chunker_test.go`
-   Test: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
-   **AC-01:** `Chunker` identifies chunks containing "swagger", "openapi", "route", "endpoint" as `ChunkTypeAPI`.
-   **AC-02:** `ResultConsumer` constructs embedding context string: `Title: ... 
Source: ... 
Path: ... 
Type: ...`.
-   **AC-03:** `ResultConsumer` fetches `Source.Name` using `source_id` (via `SourceFetcher` or Repo).

**Step 1: Write failing test**
*In `chunker_test.go`, add a test case with API-like text and assert `ChunkType` is `api`.
In `result_consumer_test.go`, check `contextualString` format.*

**Step 2: Write minimal implementation**
1.  Update `ChunkMarkdown` to check for API keywords/patterns.
2.  Update `ResultConsumer.HandleMessage` to:
    -   Extract `path` from payload.
    -   Fetch `Source` to get `Name`.
    -   Prepend header to the text before embedding.

**Step 3: Verify test passes**
*Run `go test ./apps/backend/internal/text/...` and `go test ./apps/backend/internal/worker/...`.*

### Task 4: Retrieval Improvements (Sort & Filter)

**Files:**
-   Modify: `apps/backend/internal/adapter/weaviate/store.go`
-   Test: `apps/backend/internal/adapter/weaviate/store_test.go` (if integration tests exist) or create manual check script.

**Requirements:**
-   **AC-01:** `GetChunksByURL` returns chunks sorted by `chunk_index` ASC.
-   **AC-02:** `Search` accepts `type` and `language` in `options.Filter` and applies them as `Where` clauses.

**Step 1: Write failing test**
*Mock Weaviate client or use integration test to assert sorting/filtering.*

**Step 2: Write minimal implementation**
1.  In `GetChunksByURL`, add `.WithSort(graphql.Sort{Path: []string{"chunk_index"}, Order: graphql.Asc})`.
2.  In `Search`, map `options.Filter` to `filters.Where()` clauses (Equal operator).

**Step 3: Verify test passes**

### Task 5: Upgrade MCP Tools

**Files:**
-   Modify: `apps/backend/features/mcp/handler.go`
-   Test: `apps/backend/features/mcp/handler_test.go`

**Requirements:**
-   **AC-01:** Rename tool `search` to `qurio_search`.
-   **AC-02:** Implement tool `qurio_fetch_page` (calls `GetChunksByURL` and concatenates).
-   **AC-03:** Update tool descriptions with "Strategy Guide" from PRD.
-   **AC-04:** Update `inputSchema` for `qurio_search` to support `filter` object.

**Step 1: Write failing test**
*In `handler_test.go`, assert `ListTools` returns `qurio_search` and `qurio_fetch_page`.*

**Step 2: Write minimal implementation**
1.  Rename `search` definition.
2.  Add `qurio_fetch_page` definition and handler.
    -   Handler calls `retrievalService.GetFullPage` (which calls `GetChunksByURL`).
    -   Concatenates chunk content.
3.  Update descriptions strings.

**Step 3: Verify test passes**
*Run `go test ./apps/backend/features/mcp/...`.*
