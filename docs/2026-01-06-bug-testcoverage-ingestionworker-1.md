To achieve 95% unit test coverage for the ingestion-worker (Python), we must address the remaining gaps in the single-page handlers and the asynchronous orchestration logic in main.py. While the core handlers have been standardized to return list[dict], the following strategic actions are required to eliminate the "untested corners" of the worker:
1. Exhaustive Metadata Extraction Testing
The handlers/file.py and handlers/web.py contain complex conditional logic for extracting document titles, authors, and creation dates with multiple fallback layers.
• Action: Implement property-based testing or extensive table-driven tests in test_file_handlers.py and test_web_handlers.py to cover every branch.
• Target Scenarios: Test behavior when metadata fields are None, when authors are a list vs. a string, and specifically when the creation date is a callable method versus a static value.
2. Isolate and Test the "Zombie Task" Prevention Logic
The main.py utilizes a touch_loop that explicitly cancels processing tasks if the NSQ connection is lost to prevent "zombie" tasks.
• Action: Expand test_worker_reliability.py to verify the cleanup phase of a cancelled task.
• Goal: Ensure that when current_task.cancel() is called, the finally block in process_message correctly sets the stop_touch event and awaits the task completion without leaking resources.
3. Deep-Dive into Error Taxonomy (4xx vs. 5xx)
The worker now uses a specific IngestionError taxonomy (e.g., ERR_ENCRYPTED, ERR_TIMEOUT).
• Action: Write tests that explicitly mock the pebble.ProcessPool to return ProcessExpired or TimeoutError to ensure the handle_file_task wraps these in the correct IngestionError codes.
• Traceability: Verify that the fail_payload generated in main.py correctly includes the correlation_id and the standardized error code.
4. Mock the Third-Party "Logging Leak"
A known inconsistency is that third-party libraries like tornado and pynsq can leak unstructured raw text into the structured log stream.
• Action: Add a test in test_logger.py that utilizes a standard library logger to emit a message and asserts that the structlog bridge correctly captures and renders it as machine-parsable JSON.
5. Validate Semaphore and Concurrency Limits
The worker uses an asyncio.Semaphore and max_workers to prevent OOM errors on small instances.
• Action: Implement tests that saturate the CONCURRENCY_LIMIT (currently set to 8) and verify that subsequent tasks are delayed rather than discarded.
• Robustness: The test must handle varying CPU configurations in the CI environment, as noted in the recent Test Fixes Report.

--------------------------------------------------------------------------------
Analogy for Worker Coverage: Reaching 95% coverage on the worker is like testing an automated assembly line. It’s not enough to know that the machines can build a car (Successful Ingestion); you must also verify that if the power flickers (NSQ Connection Drop), the machines stop instantly rather than flailing wildly (Zombie Tasks). You must also ensure that the sensors can read labels written in different handwriting (Metadata Fallbacks) and that every error siren sounds a specific, recognizable code rather than a generic alarm.