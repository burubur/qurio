### Task 1: Rename `qurio_fetch_page` and Upgrade Description

**Files:**
- Modify: `apps/backend/features/mcp/handler.go`
- Modify: `apps/backend/features/mcp/handler_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. `qurio_fetch_page` is renamed to `qurio_read_page`.
  2. The tool description is updated to be a detailed "system specification".
  3. **Target Description:**
     ```text
     Deep Reading / Full Context tool. Retrieves the *entire* content of a specific page or document by its URL. Use this when a search result snippet is truncated or insufficient, or when you need to read a full guide/tutorial. Crucial: Always prefer this over guessing content if the search result is incomplete.

     USAGE EXAMPLE:
     read_page(url="https://docs.stripe.com/webhooks/signatures")
     ```

- **Test Coverage**
  - Update `TestToolsCall_FetchPage` to `TestToolsCall_ReadPage` and verify it works with the new name.
  - Verify `tools/list` returns the new description.

**Step 1: Write failing test**
In `apps/backend/features/mcp/handler_test.go`:
1. Rename `TestToolsCall_FetchPage` to `TestToolsCall_ReadPage`.
2. Update calls to use `qurio_read_page`.

**Step 2: Verify test fails**
Run: `go test ./apps/backend/features/mcp/...`

**Step 3: Write minimal implementation**
In `apps/backend/features/mcp/handler.go`:
1. Rename `qurio_fetch_page` string to `qurio_read_page`.
2. Replace `Description` with the new text.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/...`

---

### Task 2: Upgrade Agent UX for Discovery Tools (`qurio_list_sources`, `qurio_list_pages`)

**Files:**
- Modify: `apps/backend/features/mcp/handler.go`

**Requirements:**
- **Acceptance Criteria**
  1. `qurio_list_sources` description is updated to guide the agent on *discovery* and includes a usage example.
     - **Target Description:**
       ```text
       Discovery tool. Lists all available documentation sets (sources) currently indexed. Use this at the start of a session to understand what documentation is available.

       USAGE EXAMPLE:
       qurio_list_sources()
       ```

  2. `qurio_list_pages` description is updated to guide the agent on *navigation* and includes a usage example.
     - **Target Description:**
       ```text
       Navigation tool. Lists all individual pages/documents within a specific source. Use this to find the exact URL of a document when a search query is too broad or to browse the table of contents.

       USAGE EXAMPLE:
       qurio_list_pages(source_id="src_stripe_api")
       ```
       *(Note: `source_id` can be found in `qurio_list_sources` output or in `qurio_search` results)*

- **Test Coverage**
  - [Unit] `TestToolsList_ReturnsQurioTools` (existing) - Manual verification that descriptions are updated.

**Step 1: Write failing test**
N/A

**Step 2: Verify test fails**
N/A

**Step 3: Write minimal implementation**
In `apps/backend/features/mcp/handler.go`:
1. Update `qurio_list_sources` Description.
2. Update `qurio_list_pages` Description.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/...`

---

### Task 3: Upgrade Agent UX for `qurio_search`

**Files:**
- Modify: `apps/backend/features/mcp/handler.go`

**Requirements:**
- **Acceptance Criteria**
  1. `qurio_search` description is refined to position it as a **"Search & Exploration"** tool.
  2. The description MUST include the full **Argument Guide** and **Usage Examples**.
  3. **Target Description:**
     ```text
     Search & Exploration tool. Performs a hybrid search (Keyword + Vector). Use this for specific questions, finding code snippets, or exploring topics across known sources.

     ARGUMENT GUIDE:

     [Alpha: Hybrid Search Balance]
     - 0.0 (Keyword): Use for Error Codes ("0x8004"), IDs ("550e8400"), or unique strings.
     - 0.3 (Mostly Keyword): Use for specific function names ("handle_web_task") where exact match matters but context helps.
     - 0.5 (Hybrid - Default): Safe bet for general queries like "database configuration".
     - 1.0 (Vector): Use for conceptual "How do I..." questions (e.g. "stop server" matches "shutdown").

     [Limit: Result Count]
     - Default: 10
     - Recommended: 5-15 (Prevent context bloat)
     - Max: 50

     [Filters: Metadata Filtering]
     - type: Filter by content type (e.g., "code", "prose", "api", "config").
     - language: Filter by language (e.g., "go", "python", "json").

     USAGE EXAMPLES:
     - Specific: search(query="webhook signature", alpha=0.3)
     - Conceptual: search(query="how to handle errors", alpha=1.0)
     - Filtered: search(query="User struct", filters={"type": "code", "language": "go"})
     ```
  4. **CRITICAL:** Update the `qurio_search` output formatting to explicitly include `SourceID` in the returned text block. This allows agents to pivot from a search result to `qurio_list_pages(source_id)`.
     - *Format:* `SourceID: [id]\n` (below Type/Language).

- **Test Coverage**
  - [Unit] `TestToolsList_ReturnsQurioTools`.
  - [Unit] `TestHandleMessage_ContextPropagation` (indirectly tests search output structure if we inspect the response, but main validation is visual/structural).

**Step 1: Write failing test**
N/A (Output formatting changes are often easier to verify by inspection or integration test, but we can rely on existing tests passing).

**Step 2: Verify test fails**
N/A

**Step 3: Write minimal implementation**
In `apps/backend/features/mcp/handler.go`:
1. Replace `qurio_search` Description with the new text.
2. In the `qurio_search` output generation loop:
   - Extract `sourceId` from metadata (if available) or ensure `SearchResult` contains it.
   - *Note:* `SearchResult` struct has `Metadata map[string]interface{}`. We need to check if `sourceId` is populated there. If not, `internal/retrieval/service.go` or `store.go` might need checking, but for now assuming it's in metadata is reasonable for this plan.
   - Add `SourceID: %s` to the formatted string.

**Step 4: Verify test passes**
Run: `go test ./apps/backend/features/mcp/...`
