# Ingestion Worker Testing Strategy & Learnings

## Dependency Injection for Crawl4AI
To avoid flaky tests caused by mocking `AsyncWebCrawler` (which is instantiated multiple times in `handle_web_task`), we have adopted a Dependency Injection (DI) pattern.

**Pattern:**
- `handle_web_task` now accepts an optional `crawler_factory` argument.
- Default factory: `def default_crawler_factory(config=None, **kwargs): return AsyncWebCrawler(config=config, **kwargs)`
- Tests: Pass a `mock_factory` (MagicMock) that returns a controlled Mock Context Manager.

**Benefits:**
- Eliminates the need to `patch` `handlers.web.AsyncWebCrawler` globally.
- Prevents mock leakage between multiple instantiations in the same function scope.
- Makes tests more robust and readable.

## Sys.modules Mocking & Test Isolation
When mocking external dependencies globally using `sys.modules` (e.g. `crawl4ai`), tests running in the same process (pytest) can suffer from pollution if the System Under Test (SUT) module holds references to old mocks.

**Best Practice:**
- Use a fixture to setup `sys.modules` mocks and **reload** the SUT module (`importlib.reload` or `del sys.modules[...]` + `import`).
- **Specific Case (2026-01-21):** A flaky test `test_web_handler_uses_timeout` was caused by `test_web_handlers.py` deleting `handlers.web` from `sys.modules`. The fix involved reloading `handlers.web` inside the test function and using `patch.object` to ensure the patch targets the correct module instance.

## llms.txt Handling & Testing
Special handling is implemented for `llms.txt` and `.txt` files to optimize ingestion and respect standard manifests.

**Crawling Differences:**
- **LLM Bypass:** URLs ending in `.txt` or `llms.txt` automatically bypass the `LLMContentFilter`.
- **Manifest Detection:** For standard URLs, the worker actively probes for `domain.com/llms.txt`.

**Testing Approach:**
- **Verification of Bypass:** Tests verify the bypass by inspecting the `crawler.arun(..., config=...)` call arguments.
- **Mock Identity:** Ensure the test and the handler are looking at the *same* `DefaultMarkdownGenerator` mock class.

## Operational Learnings (Jan 2026)
- **Crawl Timeouts:** `crawl4ai` defaults to 60s, which is often too short for documentation sites under load. This must be configurable.
- **Concurrency Rate Limits:** High concurrency (50+) with multiple replicas triggers rate limits on target sites. Smart retries with backoff are essential.

## Verification Log
- **2026-01-21:** Full suite (Unit + Integration) PASSED for Backend (Go), Frontend (Vue/Vitest), and Ingestion Worker (Python/Pytest).
