1. Logging Standards Violations
The Technical Constitution mandates "Structured logging only" using the slog standard library for Go, specifically prohibiting string formatting like fmt.Printf. It also requires every operation to log at the start, on success, and on error.
• Inconsistency: The source/handler.go file uses fmt.Printf for error logging.
• Inconsistency: The worker/ingest.go file uses the standard log package (log.Printf) instead of structured slog.
• Inconsistency: Several key handlers, such as the MCP handler and the Settings handler, contain no logging at all, failing the requirement that every request must log its start and completion.
2. Error Handling and Response Formats
The constitution requires all error paths to return a JSON envelope containing specific fields like code, message, and correlationId.
• Inconsistency: The settings/handler.go and source/handler.go files use http.Error(), which typically returns a plain text response instead of the mandated JSON envelope.
• Inconsistency: There is no evidence of correlationId generation or propagation in the provided handler implementations, despite it being a "Critical Constraint".
3. Resource Management (Timeouts)
The constitution states a universal rule to "Timeout ALL I/O operations" to prevent resource exhaustion.
• Inconsistency: While the Crawler and Reranker clients correctly implement a 10s timeout, the docling/client.go uses an unconfigured &http.Client{} which has no default timeout. This represents an inconsistent application of a "Universal Resource Management Rule"