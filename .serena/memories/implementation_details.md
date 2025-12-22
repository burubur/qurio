# Implementation Details

## Dynamic Settings & Configuration
-   **Architecture:** Database-driven configuration for runtime updates.
-   **Table:** `settings` (Singleton ID=1).
-   **Fields:** `rerank_provider`, `rerank_api_key`, `gemini_api_key`.
-   **Adapters:**
    -   `DynamicEmbedder`: Wraps Gemini client, re-initializes on key change (or per request check).
    -   `DynamicClient` (Reranker): Switches provider/key dynamically.
-   **Removal:** Environment variables `GEMINI_KEY`, `RERANK_*` removed from `config.go` and `docker-compose.yml`.

## Source Management
-   **Deduplication:**
    -   **URL Hash:** Checked at creation (prevent duplicate URLs).
    -   **Body Hash:** Calculated (`sha256`) during ingestion, stored in `body_hash` (content change detection).
-   **Lifecycle:**
    -   **Soft Delete:** Sets `deleted_at` timestamp. API filters these out.
    -   **Re-Sync:** Triggers ingestion event with existing ID. Note: Currently may duplicate vector chunks (known limitation).

## Frontend Architecture
-   **Framework:** Vue 3 + TypeScript + Vite.
-   **State Management:** Pinia (`source.store.ts`, `settings.store.ts`).
-   **Routing:** `vue-router` with history mode.
-   **Styling:** Custom CSS variables (Brand: Void Black, Cognitive Blue).
-   **Layout:** `AppLayout` with fixed `Sidebar`.
-   **Icons:** `lucide-vue-next`.

## Backend Architecture
-   **Pattern:** Feature-based (`features/source`), Internal-based (`internal/settings`, `internal/worker`).
-   **Service/Repo:** Interface-based dependency injection.
-   **Worker:** NSQ consumer for ingestion tasks.
