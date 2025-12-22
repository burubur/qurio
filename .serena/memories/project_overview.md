# Project Overview

**Status:** MVP Part 2 Complete (2025-12-22)

## Recent Progress
-   **MVP Part 2 Completed:**
    -   **Dynamic Configuration:** Moved Gemini and Reranker keys from `.env` to a singleton PostgreSQL table (`settings`).
    -   **Source Management:** Implemented Soft Delete, Re-Sync, and Content Hashing (deduplication).
    -   **Frontend Overhaul:** Implemented "Cyber-Librarian" brand identity (Void Black/Cognitive Blue), added Sidebar navigation, and `vue-router`.
    -   **Refactoring:** Removed static config dependencies for AI adapters; implemented `DynamicEmbedder` and `DynamicReranker`.
    -   **MVP Part 2.2 (Fixes & Design System):**
        -   **Frontend:** Integrated Tailwind CSS, PostCSS, and shadcn-vue. Refactored Source components.
        -   **Backend:** Standardized logging (`slog`) and error handling (JSON envelopes). Enforced timeouts.

## Upcoming Features (Planned - MVP Part 3)
-   **Advanced Retrieval:** Metadata filtering, Date-based ranking.
-   **Observability:** Dashboard for ingestion stats.
-   **Agent Integration:** MCP Server improvements.

## Active Plan
-   `docs/plans/2025-12-22-qurio-mvp-part2-2.md` (Completed)
