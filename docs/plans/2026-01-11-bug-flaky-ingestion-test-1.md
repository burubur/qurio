### Task 1: Refactor `handle_web_task` for Dependency Injection

**Files:**
- Modify: `apps/ingestion-worker/handlers/web.py`
- Modify: `apps/ingestion-worker/main.py`

**Requirements:**
- **Acceptance Criteria**
  1. `handle_web_task` accepts an optional `crawler_factory` argument.
  2. Default behavior (production) remains unchanged (uses `AsyncWebCrawler`).
  3. `main.py` passes the default factory (or relies on default argument).
  4. Both crawler instantiations (manifest check and main crawl) use the injected factory.

- **Functional Requirements**
  1. Decouple `AsyncWebCrawler` instantiation from logic.

- **Non-Functional Requirements**
  1. Maintain backward compatibility for existing consumers (except tests which we will fix).

- **Test Coverage**
  - [Unit] Existing tests in `test_web_handlers.py` should still pass (potentially with minor updates).

**Step 1: Write failing test (Conceptual - existing tests pass but we want to enable DI)**
We are refactoring to ENABLE testing, so strict TDD here means "write a test that USES the new interface and fails because it doesn't exist".
Create a temporary test `apps/ingestion-worker/tests/test_refactor_check.py`:
```python
import pytest
from handlers.web import handle_web_task

@pytest.mark.asyncio
async def test_handle_web_task_accepts_factory():
    # This should fail with TypeError: unexpected keyword argument 'crawler_factory'
    try:
        await handle_web_task("http://example.com", crawler_factory=lambda **k: None)
    except TypeError as e:
        if "unexpected keyword argument" in str(e):
            pytest.fail("handle_web_task does not accept crawler_factory")
    except Exception:
        pass # Other errors are expected as we passed a dummy factory
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_refactor_check.py`
Expected: FAIL "handle_web_task does not accept crawler_factory"

**Step 3: Write minimal implementation**
Modify `apps/ingestion-worker/handlers/web.py`:
```python
# ... imports

def default_crawler_factory(config=None, **kwargs):
    return AsyncWebCrawler(config=config, **kwargs)

async def handle_web_task(url: str, api_key: str = None, crawler_factory=default_crawler_factory) -> list[dict]:
    # ...
    # Replace: async with AsyncWebCrawler(verbose=False) as manifest_crawler:
    # With: async with crawler_factory(verbose=False, config=manifest_config) as manifest_crawler:
    
    # ...
    # Replace: async with AsyncWebCrawler(verbose=True) as crawler:
    # With: async with crawler_factory(verbose=True, config=config) as crawler:
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_refactor_check.py`
Expected: PASS (or fail with a different error, but NOT TypeError on argument)
(Delete the temporary test after verification)


### Task 2: Fix `test_manifest_detection.py` using DI

**Files:**
- Modify: `apps/ingestion-worker/tests/test_manifest_detection.py`

**Requirements:**
- **Acceptance Criteria**
  1. `test_detects_and_merges_llms_txt_links` is unskipped.
  2. Test uses `crawler_factory` injection instead of global patching.
  3. Test passes reliably.

- **Test Coverage**
  - [Unit] `test_detects_and_merges_llms_txt_links` covers manifest detection + merge logic.

**Step 1: Write failing test (Unskip)**
Modify `apps/ingestion-worker/tests/test_manifest_detection.py`:
- Remove `@pytest.mark.skip`
- Run: `pytest apps/ingestion-worker/tests/test_manifest_detection.py`
- Expected: FAIL (original TypeError or similar mock issue)

**Step 2: Refactor test to use DI**
Refactor `test_detects_and_merges_llms_txt_links` to:
1. Create a `mock_factory` that returns a `MagicMock` (context manager).
2. The context manager returns `mock_crawler`.
3. Configure `mock_crawler.arun` side effects (first call -> manifest, second call -> main).
4. Call `handle_web_task(url, crawler_factory=mock_factory)`.
5. Remove global `sys.modules` patching hacks if possible (or keep them minimal if imports require it).

**Step 3: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_manifest_detection.py`
Expected: PASS


### Task 3: Refactor `test_web_handlers.py` and `test_llms_txt_bypass.py`

**Files:**
- Modify: `apps/ingestion-worker/tests/test_web_handlers.py`
- Modify: `apps/ingestion-worker/tests/test_llms_txt_bypass.py`

**Requirements:**
- **Acceptance Criteria**
  1. All tests in `test_web_handlers.py` use DI.
  2. `test_llms_txt_bypass.py` uses DI and mocks are verified correctly.
  3. No reliance on fragile `patch('handlers.web.AsyncWebCrawler')`.

- **Test Coverage**
  - [Unit] Regression testing for existing functionality.

**Step 1: Write failing test (None - refactoring)**
Run all tests to confirm current state (they might pass if using default factory, but we want to enforce DI usage in tests).

**Step 2: Refactor `test_web_handlers.py`**
- Replace `with patch('handlers.web.AsyncWebCrawler', ...)` with `mock_factory`.
- Pass `crawler_factory=mock_factory` to `handle_web_task`.

**Step 3: Refactor `test_llms_txt_bypass.py`**
- Replace `with patch('handlers.web.AsyncWebCrawler', ...)` with `mock_factory`.
- Pass `crawler_factory=mock_factory` to `handle_web_task`.

**Step 4: Verify all tests pass**
Run: `pytest apps/ingestion-worker/tests/`
Expected: PASS (100%)
