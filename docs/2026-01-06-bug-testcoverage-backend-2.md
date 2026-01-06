The current backend unit test coverage stands at 55.7% of statements. While core logic packages like internal/text and internal/settings have reached high coverage (>90%), the overall report is significantly skewed by the 0% coverage of main.go, which contains the primary application wiring and orchestration logic.
To push the backend coverage into the 95% region, the following strategic actions must be taken:
1. Complete the Decoupling of main.go
The primary reason for the low coverage is that main.go currently handles infrastructure pings, migrations, and service initialization. Because these depend on real connections, they are rarely covered by unit tests.
• Action: Move all infrastructure setup (DB pings, Weaviate schema checks, NSQ producer connection logic) out of main.go and into the internal/app package.
• Goal: main.go should only contain a simple call to an app.Run() function. This allows you to test the entire wiring lifecycle in app_test.go by passing in mocked dependencies to the app.New() constructor.
2. Deep Testing of the ResultConsumer
The ResultConsumer is the most critical logic hub in the backend, currently sitting at 70% coverage. Reaching 95% requires testing all negative paths and "micro-logic" decisions within the HandleMessage loop:
• Test Payload Variations: Create tests for corrupted JSON payloads, missing correlation_id, and every combination of task_type.
• Error Path Injection: Use mocks for the PageManager, TaskPublisher, and Embedder that explicitly return errors to verify that the consumer handles failures without crashing or losing data.
• Context Timeouts: Specifically test that the 60-second timeouts for embedding and storage are correctly triggered and handled.
3. Exhaustive Table-Driven Tests for MCP Tools
The MCP handler.go is a large file with complex JSON-RPC unmarshaling and multiple tool definitions (qurio_search, qurio_read_page, etc.).
• Action: Implement table-driven tests for the processRequest function.
• Coverage Targets: Each tool needs tests for valid arguments, missing required fields (e.g., missing url in qurio_read_page), and internal service errors. This ensures the massive switch block and unmarshaling logic are fully exercised.
4. Close the Gap in Feature Handlers
The features/source (64%) and features/job (59%) handlers have significant uncovered blocks, primarily in their HTTP response logic.
• HTTP Status Codes: Write tests to trigger every h.writeError call in the handlers.
• MIME/Size Validation: For the Upload handler, explicitly test files exceeding the 50MB limit and unsupported file types to cover the early return paths.
5. Standardize Adapter Testing
Adapters for Weaviate and Gemini often lack coverage for retry logic and network edge cases.
• Mock Initialization: Use the established pattern of httptest.NewServer to simulate Weaviate's /v1/meta and /v1/graphql endpoints, allowing you to test the adapter's behavior when the database returns specific GraphQL errors.

--------------------------------------------------------------------------------
Analogy for Coverage Improvement: Reaching 95% coverage is like waterproofing a ship. Right now, your ship has a strong hull in some places (Core Logic), but the engine room (main.go) is completely open to the elements, and there are dozens of tiny, unsealed rivets (Error Paths) throughout the deck. To sail safely in high-pressure environments (Production), you must seal every tiny gap where an error could leak in, and ensure that even if one room floods (Service Failure), your bulkheads (Error Handling) are tested and ready to hold.