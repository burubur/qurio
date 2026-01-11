# Ingestion Worker Testing Strategy Update (Jan 2026)

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

**Example Test Usage:**
```python
mock_factory = MagicMock()
mock_factory.return_value.__aenter__.return_value = mock_crawler
await handle_web_task(url, crawler_factory=mock_factory)
```

## Sys.modules Mocking & Test Isolation
When mocking external dependencies globally using `sys.modules` (e.g. `crawl4ai`), tests running in the same process (pytest) can suffer from pollution if the System Under Test (SUT) module holds references to old mocks.

**Best Practice:**
- Use a fixture to setup `sys.modules` mocks and **reload** the SUT module (`importlib.reload` or `del sys.modules[...]` + `import`).
- When using `patch` on attributes of a module that imports from a mocked module, ensure the SUT module is the one currently in `sys.modules`.
- If `patch` seems to fail (target not called), it often means the test function is holding a reference to an old version of the function/module, or the module was reloaded by another test and `patch` is patching the new one while your test uses the old one (or vice versa).
- Explicitly importing the SUT module inside the test function or `with patch` block can help verify/ensure you are patching the correct object.

## llms.txt Handling & Testing
Special handling is implemented for `llms.txt` and `.txt` files to optimize ingestion and respect standard manifests.

**Crawling Differences:**
- **LLM Bypass:** URLs ending in `.txt` or `llms.txt` automatically bypass the `LLMContentFilter`.
  - **Reason:** These files are assumed to be dense, structured technical content (manifests, documentation in text format) that does not require "instruction-based" extraction or noise removal (navbars, footers).
  - **Mechanism:** `handle_web_task` initializes `DefaultMarkdownGenerator` directly instead of passing it a `LLMContentFilter`.
- **Manifest Detection:** For standard URLs, the worker actively probes for `domain.com/llms.txt`. If found, it returns *two* results (manifest + main page).

**Testing Approach:**
- **Verification of Bypass:** Since `crawl4ai` classes are mocked globally, we cannot use `isinstance` checks.
- **Config Inspection:** Tests verify the bypass by inspecting the `crawler.arun(..., config=...)` call arguments.
  - **Bypass Case:** Assert `config.markdown_generator == DefaultMarkdownGenerator.return_value` (the exact mock instance).
  - **Standard Case:** Assert `config.markdown_generator` was initialized with `content_filter=...`.
- **Mock Identity:** When testing this logic, ensure the test and the handler are looking at the *same* `DefaultMarkdownGenerator` mock class. Using `patch.object(handlers_web, 'DefaultMarkdownGenerator')` on the reloaded module is the most reliable way to intercept the class instantiation used by the code.