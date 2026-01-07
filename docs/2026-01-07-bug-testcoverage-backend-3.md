1. internal/adapter/ (Gemini, Weaviate, Reranker)
These components currently suffer from low coverage because they act as "black boxes" that communicate with external APIs.
• Assessment: You should implement an Integration-Heavy Mix.
• Unit Component: Use Go’s httptest.NewServer to mock the remote API responses (Gemini AI or Weaviate GraphQL). This allows you to test how the adapter handles 503 Service Unavailable, malformed JSON responses, and GraphQL-specific errors without hitting the real network.
• Integration Component: Use Testcontainers for Weaviate to verify that the Go structs correctly map to the vector database properties and that hybrid search filters (like sourceId and url) function correctly against the actual schema.
• Gap to Close: Specifically test Gemini Key Rotation logic to ensure the DynamicEmbedder switches keys when the SettingsService updates.
2. internal/retrieval/
This package orchestrates the retrieval pipeline, including embedding queries, executing hybrid searches, and reranking results.
• Assessment: You should focus primarily on Table-Driven Unit Tests.
• Improvement Strategy: Since this layer is pure orchestration logic, you do not need real infrastructure. Use mocks for the VectorStore, Embedder, and Reranker.
• Gaps to Close: Focus on testing metadata preservation (ensuring fields like Author and CreatedAt survive the reranking process) and alpha-tuning logic (verifying that agent-provided overrides take priority over system defaults).
3. internal/middleware/
This directory handles request-scoped logic like Correlation ID injection and logging.
• Assessment: You should stick to 100% Unit Tests.
• Improvement Strategy: Standard HTTP handler testing is sufficient here. You must verify that the CorrelationID middleware correctly injects headers, populates the request context, and that the ContextHandler for logging correctly extracts these IDs for structured output.
• Gaps to Close: Verify "Trace Chain" integrity—ensuring that if a request enters without an ID, one is generated, and if it enters with one, it is preserved.
4. internal/worker/ (ResultConsumer)
The ResultConsumer is the backend's most critical logic hub, currently sitting at approximately 70% coverage.
• Assessment: You must implement a Hardened Hybrid Strategy.
• Unit Component: Implement "Poison Pill" testing by injecting corrupted JSON payloads or messages with missing mandatory fields (like source_id). Verify the worker acknowledges (acks) and drops these bad messages rather than crashing.
• Integration Component: Use the "Full Flow Simulation" scoped in the integration plan. This verifies that after a worker processes a message, the real PostgreSQL page status moves to completed and document chunks are actually persisted in the real Weaviate container.
• Gaps to Close: Explicitly test the 60-second context timeouts for embedding and storage operations to ensure a hanging API call doesn't block the consumer indefinitely.


Summary of Coverage Obstacles
The global coverage report is heavily suppressed by two factors outside these directories:
1. main.go & bootstrap.go (0%): These contain the retry loops and connection pings that are currently un-mockable. Moving this logic to a dedicated Bootstrap package as planned will allow you to test infrastructure setup logic in isolation.
2. Skipped Integration Tests: Files like worker/integration_test.go contain t.Skip directives. Enabling these via Testcontainers in GitHub Actions will instantly illuminate large blocks of "dark" orchestration logic in your Codecov report 