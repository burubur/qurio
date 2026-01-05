# Implementation Report: Bug Fixes & Inconsistencies (Jan 5, 2026)

## Standardized Ingestion Worker Handlers
- Refactored `handle_file_task` in `apps/ingestion-worker` to return `list[dict]`, matching `handle_web_task` signature.
- Removed brittle manual list wrapping in `main.py`.
- Verified with updated unit tests (`test_file_handlers.py`).

## MCP Tool Usability
- `qurio_search` output now explicitly instructs agents to use `qurio_read_page` for full content.
- Added regression test `TestQurioSearchInstruction`.

## Metadata Exposure
- Promoted metadata fields (`Author`, `CreatedAt`, `PageCount`, `Language`, `Type`, `SourceID`, `URL`) to top-level fields in `SearchResult` struct.
- Updated Weaviate adapter to populate these fields from GraphQL response.
- Refactored MCP handler to use strong typing instead of `Metadata` map lookups.
- Verified with integration test `TestStore_Search_PopulatesMetadata`.
