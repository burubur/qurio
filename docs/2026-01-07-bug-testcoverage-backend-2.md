1. Strategy for Low-Coverage Packages
For the packages you identified—adapter/, retrieval, middleware/, and worker/—the project sources suggest a mix of unit and integration testing to seal the "unsealed rivets" in the codebase.
A. Internal Adapters (qurio/apps/backend/internal/adapter/)
• The Issue: These are currently "black boxes". Unit tests often mock the clients (like Weaviate or Gemini), so you aren't testing the actual GraphQL query building or API response unmarshaling.
• Recommendation: Integration-Heavy Mix.
    ◦ Use Testcontainers for Weaviate to verify that complex filters (like sourceId and url) match the actual schema types.
    ◦ Use Go’s httptest.NewServer to simulate network failures (503 Service Unavailable) and GraphQL errors for the Weaviate and Gemini adapters.
B. Retrieval Service (qurio/apps/backend/internal/retrieval)
• The Issue: This package orchestrates hybrid search and reranking. Coverage is likely low on the reranking preservation and metadata mapping paths.
• Recommendation: Unit Tests.
    ◦ You can reach high coverage here using Table-Driven Unit Tests. Mock the VectorStore and Reranker to return various result counts and verify that metadata (Author, CreatedAt) is correctly preserved through the pipeline.
C. Middleware (qurio/apps/backend/internal/middleware/)
• The Issue: This is simple logic, but critical for the "trace chain".
• Recommendation: Unit Tests.
    ◦ Standard HTTP handler testing is sufficient here. Verify that the CorrelationID middleware correctly injects headers and populates the request context. This requires zero external infrastructure.
D. Ingestion Worker (qurio/apps/backend/internal/worker/)
• The Issue: The ResultConsumer is the most critical logic hub, but it currently lacks tests for "Poison Pill" messages (malformed JSON) and context timeouts.
• Recommendation: Hybrid Hardening.
    ◦ Unit Tests: Add tests for corrupted NSQ bodies to ensure the worker acknowledges and drops "poison pills" rather than crashing.
    ◦ Integration Tests: Use the "Full Flow Simulation" scoped in the plan to verify that after a message is processed, the real PostgreSQL status moves to completed and chunks exist in the real Weaviate [Plan Phase 3].
2. Improving "Pure Business Logic" (Rule 2)
The sources note a violation of Rule 2 (Pure Business Logic) in the worker. To improve testability and coverage without containers, you should extract logic from I/O.
• Example: The link discovery logic was extracted from the ResultConsumer into a pure function in link_discovery.go. This allowed 100% coverage of complex exclusion regex and depth-matching logic through fast unit tests, without needing a database.