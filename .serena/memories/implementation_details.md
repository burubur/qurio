# Implementation Update (2025-12-21)

Completed MVP Refinement tasks:

1.  **Ingestion/Chunking:**
    -   Added `apps/backend/internal/text` package with `Chunk` function (512 tokens, 50 overlap).
    -   Updated `IngestHandler` to process multiple chunks per source URL.
    -   Updated `Chunk` struct to include `ChunkIndex`.

2.  **Crawling:**
    -   Updated `Crawler` to discover links from `/sitemap.xml` and `/llms.txt`.
    -   Discovery runs as a pre-crawl seed expansion.

3.  **Retrieval:**
    -   Updated `Weaviate.Store` to support Hybrid Search with configurable Alpha.
    -   Updated `Retrieval.Service` to default Alpha to 0.5.
    -   Added `Jina` Reranker support via `adapter/reranker`.

4.  **MCP:**
    -   Verified and tested `search` tool JSON-RPC handler.

All unit/integration tests passed.
