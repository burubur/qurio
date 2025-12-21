# Implementation Details

## Backend Architecture (Hexagonal/Clean)

### Core Logic (Internal)
- **Worker:** `internal/worker` - Handles ingestion (Fetch -> Embed -> Store).
- **Retrieval:** `internal/retrieval` - Semantic search logic.
- **Vector:** `internal/vector` - Schema management.
- **Config:** `internal/config` - Env-based configuration.

### Adapters (Infrastructure)
- **Gemini:** `internal/adapter/gemini` - Implements `Embedder` using Google GenAI.
- **Docling:** `internal/adapter/docling` - Implements `Fetcher` via HTTP call to Python service.
- **Weaviate:** `internal/adapter/weaviate` - Implements `VectorStore` (Store/Search).
- **Postgres:** `features/source/repo.go` - Implements `Repository`.

### Wiring
- `main.go` wires all adapters to services using dependency injection.
- Conditional initialization: Gemini and NSQ only start if config is present.

## Frontend Architecture
- **Tech:** Vue 3 + TypeScript + Vite.
- **Testing:** Vitest + Vue Test Utils.
- **Components:** Feature-based (`features/sources`).
