### Task 1: Python Worker - `llms.txt` Content Filter Bypass

**Files:**
- Modify: `apps/ingestion-worker/handlers/web.py`
- Test: `apps/ingestion-worker/tests/test_llms_txt_bypass.py`

**Requirements:**
- **Acceptance Criteria**
  1. URLs ending in `/llms.txt` must NOT use `LLMContentFilter`.
  2. URLs ending in `/llms.txt` MUST use `DefaultMarkdownGenerator` (or equivalent raw generator).
  3. Standard URLs must continue to use `LLMContentFilter`.

- **Functional Requirements**
  1. In `handle_web_task`, check if `task.url` ends with `.txt` (specifically `llms.txt` based on context, but `.txt` might be safer for all text files, though user specified `llms.txt`). Let's stick to `llms.txt` or `.txt` if it's a manifest. User said "Manifest Bypass... identify manifest files (ending in .txt)". I will use `.txt` as the trigger.
  2. Configure `CrawlerRunConfig` with `markdown_generator=DefaultMarkdownGenerator()` when triggered.

- **Non-Functional Requirements**
  - None for this task.

- **Test Coverage**
  - [Unit] `handle_web_task` with `https://example.com/llms.txt` -> asserts `CrawlerRunConfig.markdown_generator` is `DefaultMarkdownGenerator`.
  - [Unit] `handle_web_task` with `https://example.com/page` -> asserts `CrawlerRunConfig.markdown_generator` is `LLMContentFilter`.

**Step 1: Write failing test**
```python
import pytest
from unittest.mock import patch, MagicMock
from handlers.web import handle_web_task
from crawl4ai import DefaultMarkdownGenerator

@pytest.mark.asyncio
async def test_llms_txt_uses_default_generator():
    # Arrange
    mock_task = MagicMock()
    mock_task.url = "https://example.com/llms.txt"
    mock_task.source_id = "test_source"
    
    with patch("handlers.web.AsyncWebCrawler") as MockCrawler:
        mock_crawler_instance = MockCrawler.return_value.__aenter__.return_value
        mock_crawler_instance.arun.return_value = MagicMock(markdown="[link](http://test.com)")
        
        # Act
        await handle_web_task(mock_task)
        
        # Assert
        # Check the config passed to arun
        call_args = mock_crawler_instance.arun.call_args
        assert call_args is not None
        config = call_args.kwargs.get('config')
        assert isinstance(config.markdown_generator, DefaultMarkdownGenerator)

@pytest.mark.asyncio
async def test_standard_page_uses_llm_filter():
    # Arrange
    mock_task = MagicMock()
    mock_task.url = "https://example.com/page"
    mock_task.source_id = "test_source"
    
    with patch("handlers.web.AsyncWebCrawler") as MockCrawler:
        mock_crawler_instance = MockCrawler.return_value.__aenter__.return_value
        mock_crawler_instance.arun.return_value = MagicMock(markdown="content")
        
        # Act
        await handle_web_task(mock_task)
        
        # Assert
        call_args = mock_crawler_instance.arun.call_args
        config = call_args.kwargs.get('config')
        # Assuming LLMContentFilter is the default or explicitly checked
        # For this test, just ensuring it's NOT DefaultMarkdownGenerator might be enough if imports are tricky
        # But ideally we check for LLMContentFilter
        from crawl4ai import LLMContentFilter
        assert isinstance(config.markdown_generator, LLMContentFilter)
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_llms_txt_bypass.py`
Expected: FAIL (both likely use LLMContentFilter currently)

**Step 3: Write minimal implementation**
```python
# In apps/ingestion-worker/handlers/web.py

from crawl4ai import DefaultMarkdownGenerator

# Inside handle_web_task...
    # Determine generator based on file type
    if task.url.endswith('.txt'):
        markdown_generator = DefaultMarkdownGenerator()
    else:
        markdown_generator = LLMContentFilter(...) # Existing logic

    run_config = CrawlerRunConfig(
        markdown_generator=markdown_generator,
        # ... other existing params
    )
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_llms_txt_bypass.py`
Expected: PASS

---

### Task 2: Python Worker - Root Manifest Detection

**Files:**
- Modify: `apps/ingestion-worker/handlers/web.py`
- Test: `apps/ingestion-worker/tests/test_manifest_detection.py`

**Requirements:**
- **Acceptance Criteria**
  1. `handle_web_task` must attempt to fetch `llms.txt` at the root of the domain for new sources.
  2. If found, links from `llms.txt` must be included in the result.
  3. Failure to fetch `llms.txt` should be silent/non-blocking for the main task.

- **Functional Requirements**
  1. Parse base URL from `task.url`.
  2. Construct `{base_url}/llms.txt`.
  3. Perform a quick fetch (HEAD or GET) to check existence.
  4. If exists, crawl/parse it (using the logic from Task 1, i.e., no LLM filter).
  5. Merge discovered links into the main result's links.

- **Non-Functional Requirements**
  - Performance: Do not delay main crawl significantly. Run concurrently if possible, or use a short timeout.

- **Test Coverage**
  - [Unit] `handle_web_task` detects `llms.txt` -> returns links from BOTH main page and `llms.txt`.
  - [Unit] `handle_web_task` 404 on `llms.txt` -> returns links from main page only.

**Step 1: Write failing test**
```python
import pytest
from unittest.mock import patch, MagicMock
from handlers.web import handle_web_task

@pytest.mark.asyncio
async def test_detects_and_merges_llms_txt_links():
    # Arrange
    mock_task = MagicMock()
    mock_task.url = "https://example.com/home"
    
    with patch("handlers.web.AsyncWebCrawler") as MockCrawler:
        crawler = MockCrawler.return_value.__aenter__.return_value
        
        # Setup side effects for arun to simulate two calls:
        # 1. Main page
        # 2. llms.txt (or vice versa depending on impl order)
        
        async def mock_arun(url, config):
            if url.endswith("llms.txt"):
                return MagicMock(markdown="[Manifest Link](https://example.com/manifest-dest)")
            return MagicMock(markdown="[Main Link](https://example.com/main-dest)")
            
        crawler.arun.side_effect = mock_arun
        
        # Act
        result = await handle_web_task(mock_task)
        
        # Assert
        assert "https://example.com/manifest-dest" in result.links
        assert "https://example.com/main-dest" in result.links
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_manifest_detection.py`
Expected: FAIL (Mock will likely only return main link)

**Step 3: Write minimal implementation**
```python
# In apps/ingestion-worker/handlers/web.py

# ... inside handle_web_task ...
    links = []
    
    # Pre-flight check for llms.txt (simplified for plan)
    from urllib.parse import urlparse
    parsed = urlparse(task.url)
    base_url = f"{parsed.scheme}://{parsed.netloc}"
    manifest_url = f"{base_url}/llms.txt"

    # Try fetch manifest
    try:
        manifest_config = CrawlerRunConfig(markdown_generator=DefaultMarkdownGenerator())
        manifest_res = await crawler.arun(url=manifest_url, config=manifest_config)
        if manifest_res.success:
            links.extend(extract_links(manifest_res.markdown)) # Pseudo-code extraction
    except Exception:
        pass # Silent fail

    # ... Proceed with main crawl ...
    main_res = await crawler.arun(url=task.url, config=main_config)
    links.extend(extract_links(main_res.markdown))
    
    return Result(links=links, ...)
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_manifest_detection.py`
Expected: PASS

---

### Task 3: Backend - `llms.txt` Virtual Depth Handling

**Files:**
- Modify: `apps/backend/internal/worker/result_consumer.go`
- Test: `apps/backend/internal/worker/result_consumer_test.go`

**Requirements:**
- **Acceptance Criteria**
  1. Links extracted from an `llms.txt` source must be processed even if `currentDepth >= maxDepth`.
  2. Effectively, `llms.txt` should be treated as "depth -1" or allow +1 depth extension.

- **Functional Requirements**
  1. In `HandleMessage` (or where `DiscoverLinks` is called), check `payload.SourceURL`.
  2. If `SourceURL` ends with `/llms.txt`, increment the local `maxDepth` variable passed to `DiscoverLinks` (or the logic checking it) by 1.
  3. Ensure this only affects the *current* processing step, not the persisted `maxDepth` of the source.

- **Non-Functional Requirements**
  - None.

- **Test Coverage**
  - [Unit] `HandleMessage` with `source_url=".../llms.txt"`, `depth=0`, `maxDepth=0`. Expectation: Links are queued (because effective maxDepth becomes 1).
  - [Unit] `HandleMessage` with `source_url=".../page"`, `depth=0`, `maxDepth=0`. Expectation: Links NOT queued (standard behavior).

**Step 1: Write failing test**
```go
// In apps/backend/internal/worker/result_consumer_test.go

func TestHandleMessage_LLMsTxt_BypassesDepth(t *testing.T) {
    // Arrange
    consumer := &ResultConsumer{
        // ... mocks ...
    }
    payload := &TaskResult{
        SourceURL: "https://example.com/llms.txt",
        Links: []string{"https://example.com/found"},
        Depth: 0,
        MaxDepth: 0, 
    }
    
    // Act
    consumer.HandleMessage(context.Background(), payload)
    
    // Assert
    // Check mock queue to see if "https://example.com/found" was enqueued
    // If logic is standard, 0 >= 0, so no discovery -> Queue empty -> FAIL
}
```

**Step 2: Verify test fails**
Run: `go test ./apps/backend/internal/worker/...`
Expected: FAIL

**Step 3: Write minimal implementation**
```go
// In apps/backend/internal/worker/result_consumer.go

func (r *ResultConsumer) HandleMessage(ctx context.Context, payload *TaskResult) error {
    // ...
    effectiveMaxDepth := payload.MaxDepth
    if strings.HasSuffix(payload.SourceURL, "/llms.txt") {
        effectiveMaxDepth++
    }
    
    // Pass effectiveMaxDepth to DiscoverLinks instead of payload.MaxDepth
    newTasks := r.linkDiscoverer.DiscoverLinks(payload.Links, payload.Depth, effectiveMaxDepth)
    // ...
}
```

**Step 4: Verify test passes**
Run: `go test ./apps/backend/internal/worker/...`
Expected: PASS
