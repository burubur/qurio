# Implementation Details - Chunk Pagination

**Date:** 2026-01-08

## Backend
- **Interface Updates:** `ChunkStore` and `VectorStore` interfaces now support `GetChunks(ctx, id, limit, offset)` and `CountChunksBySource(ctx, id)`.
- **Weaviate Adapter:** Implemented pagination using GraphQL `limit`/`offset` and counting using `Aggregate` with `where` filter.
- **Service Layer:** `Service.Get` now retrieves `TotalChunks` separately and supports fetching specific pages of chunks.
- **API Handler:** `GET /sources/{id}` accepts `limit`, `offset`, and `exclude_chunks` query parameters.

## Frontend
- **Store:** `getSource` supports pagination params. Added `fetchChunks` for incremental loading and `pollSourceStatus` for lightweight polling.
- **UI:** `SourceDetailView` implements a "Load More" button and uses optimized polling (metadata only) to prevent overwriting user's scroll state.
