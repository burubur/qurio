### Task 1: Exhaustive Metadata Extraction Testing

**Files:**
- Create: `apps/ingestion-worker/tests/test_metadata_extraction.py`
- Modify: `apps/ingestion-worker/handlers/file.py` (if needed for testability)
- Modify: `apps/ingestion-worker/handlers/web.py` (if needed for testability)

**Requirements:**
- **Acceptance Criteria**
  1. `handlers/file.py` correctly extracts metadata (title, author, date) from `docling` results.
  2. `handlers/web.py` correctly extracts metadata from `crawl4ai` results.
  3. Tests cover edge cases: missing fields, different formats (list vs string authors), callable vs static dates.

- **Functional Requirements**
  1. Handle `docling` v2 metadata structure.
  2. Handle `crawl4ai` extraction results.

- **Non-Functional Requirements**
  1. Use property-based testing or extensive table-driven tests (`pytest.mark.parametrize`).

- **Test Coverage**
  - [Unit] `test_file_metadata_extraction` - Parameterized test for docling metadata.
  - [Unit] `test_web_metadata_extraction` - Parameterized test for web metadata.

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_metadata_extraction.py
import pytest
from datetime import datetime
from unittest.mock import MagicMock
from handlers.file import handle_file_task # Assuming this imports logic or I might need to extract the metadata logic to a pure function for easier testing

# Note: Ideally, extraction logic should be a pure function. 
# If it's embedded in handle_file_task, we should extract it or mock heavily.
# For this plan, assuming we extract `extract_metadata(doc_result)` to make it testable.

@pytest.mark.parametrize("doc_result, expected", [
    (MagicMock(input={"title": "T"}, metadata={"author": "A"}), {"title": "T", "author": "A"}),
    (MagicMock(input={}, metadata=None), {"title": "Untitled", "author": "Unknown"}),
    # Add more cases: list authors, callable date, etc.
])
def test_file_metadata_extraction_logic(doc_result, expected):
    from handlers.file import extract_metadata_from_doc # We need to create this refactor
    assert extract_metadata_from_doc(doc_result) == expected
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_metadata_extraction.py`
Expected: FAIL (ImportError or Assertion Error as function doesn't exist yet)

**Step 3: Write minimal implementation**
```python
# Refactor handlers/file.py to extract metadata logic
def extract_metadata_from_doc(doc_result):
    # Implementation logic moved here
    pass

# Update handle_file_task to use this function
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_metadata_extraction.py`
Expected: PASS


### Task 2: Isolate and Test "Zombie Task" Prevention

**Files:**
- Create: `apps/ingestion-worker/tests/test_zombie_prevention.py`
- Modify: `apps/ingestion-worker/main.py`

**Requirements:**
- **Acceptance Criteria**
  1. When a task is cancelled (e.g., NSQ connection lost), the cleanup phase runs.
  2. `stop_touch` event is set.
  3. No "zombie" processes remain.

- **Functional Requirements**
  1. `process_message` must have a `finally` block or `except CancelledError`.
  2. `touch_loop` must respond to `stop_touch`.

- **Test Coverage**
  - [Integration] `test_task_cancellation_cleanup` - Mock NSQ connection drop/task cancel.

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_zombie_prevention.py
import pytest
import asyncio
from unittest.mock import MagicMock, AsyncMock
from main import process_message

@pytest.mark.asyncio
async def test_process_message_cleanup_on_cancel():
    # Setup mocks
    mock_nsq_msg = MagicMock()
    mock_nsq_msg.body = b'{"url": "http://example.com"}'
    mock_nsq_msg.touch = MagicMock()
    
    stop_event = asyncio.Event()
    
    # Mock handle_web_task to sleep forever to simulate long work
    with pytest.patch("main.handle_web_task", new_callable=AsyncMock) as mock_handler:
        mock_handler.side_effect = asyncio.sleep(10)
        
        # Run process_message in a task
        task = asyncio.create_task(process_message(mock_nsq_msg, {}))
        
        await asyncio.sleep(0.1) # Let it start
        
        # Cancel the task
        task.cancel()
        
        try:
            await task
        except asyncio.CancelledError:
            pass
            
        # Verify cleanup (e.g. check if specific cleanup logs or side effects happened)
        # For now, we assume process_message logic doesn't strictly enforce cleanup yet
        # or we verify a side effect that we will implement.
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_zombie_prevention.py`
Expected: FAIL (if current implementation doesn't clean up properly or if we add a spy that isn't called)

**Step 3: Write minimal implementation**
```python
# main.py
async def process_message(msg, config):
    stop_touch = asyncio.Event()
    # ... start touch loop ...
    try:
        # dispatch task
        await handle_task(...)
    finally:
        stop_touch.set() # Ensure this happens even on CancelledError
        # ... wait for touch loop to finish ...
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_zombie_prevention.py`
Expected: PASS


### Task 3: Deep-Dive into Error Taxonomy

**Files:**
- Create: `apps/ingestion-worker/tests/test_error_taxonomy.py`
- Modify: `apps/ingestion-worker/main.py`
- Modify: `apps/ingestion-worker/handlers/file.py`

**Requirements:**
- **Acceptance Criteria**
  1. `pebble.ProcessExpired` maps to `IngestionError(ERR_TIMEOUT)`.
  2. `TimeoutError` maps to `IngestionError(ERR_TIMEOUT)`.
  3. `fail_payload` includes correlation_id and error code.

- **Functional Requirements**
  1. Update `handle_file_task` to catch generic errors and wrap them.

- **Test Coverage**
  - [Unit] `test_error_mapping` - Verify exception wrapping.
  - [Unit] `test_fail_payload_structure` - Verify JSON output.

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_error_taxonomy.py
import pytest
from pebble import ProcessExpired
from handlers.file import handle_file_task
from exceptions import IngestionError

def test_pebble_timeout_mapping():
    mock_pool = MagicMock()
    mock_future = MagicMock()
    mock_future.result.side_effect = ProcessExpired("Timeout")
    mock_pool.schedule.return_value = mock_future
    
    with pytest.raises(IngestionError) as exc:
        handle_file_task(..., _pool=mock_pool)
    
    assert exc.value.code == "ERR_TIMEOUT"
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_error_taxonomy.py`
Expected: FAIL

**Step 3: Write minimal implementation**
```python
# handlers/file.py
def handle_file_task(...):
    try:
        # ...
    except ProcessExpired:
        raise IngestionError("ERR_TIMEOUT", "Processing timed out")
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_error_taxonomy.py`
Expected: PASS


### Task 4: Mock Third-Party "Logging Leak"

**Files:**
- Create: `apps/ingestion-worker/tests/test_logging_bridge.py`
- Modify: `apps/ingestion-worker/logger.py`

**Requirements:**
- **Acceptance Criteria**
  1. Standard library `logging` messages are captured by `structlog`.
  2. Output is JSON formatted.

- **Functional Requirements**
  1. Configure `structlog` to wrap stdlib.

- **Test Coverage**
  - [Unit] `test_stdlib_logging_capture` - Emit logs via `logging` and verify capture.

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_logging_bridge.py
import logging
import pytest
from logger import setup_logging
import structlog
from io import StringIO

def test_stdlib_logging_capture(capsys):
    setup_logging()
    logger = logging.getLogger("third_party_lib")
    logger.info("leaky message")
    
    captured = capsys.readouterr()
    assert '"event": "leaky message"' in captured.out
    assert '"logger": "third_party_lib"' in captured.out
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_logging_bridge.py`
Expected: FAIL

**Step 3: Write minimal implementation**
```python
# logger.py
def setup_logging():
    # Configure structlog.stdlib.LoggerFactory
    # Add StandardLibHandler to root logger
    pass
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_logging_bridge.py`
Expected: PASS


### Task 5: Validate Semaphore and Concurrency Limits

**Files:**
- Create: `apps/ingestion-worker/tests/test_concurrency.py`
- Modify: `apps/ingestion-worker/main.py` (if needed to expose semaphore)

**Requirements:**
- **Acceptance Criteria**
  1. Verify tasks respect `CONCURRENCY_LIMIT`.
  2. Excess tasks wait.

- **Functional Requirements**
  1. `process_message` uses `asyncio.Semaphore`.

- **Test Coverage**
  - [Integration] `test_semaphore_saturation`.

**Step 1: Write failing test**
```python
# apps/ingestion-worker/tests/test_concurrency.py
import pytest
import asyncio
from main import semaphore, process_message # Need to expose semaphore or inject it

@pytest.mark.asyncio
async def test_semaphore_saturation():
    # Assume CONCURRENCY_LIMIT = 8
    # Start 10 tasks that sleep
    # Assert 2 are pending initially
    pass
```

**Step 2: Verify test fails**
Run: `pytest apps/ingestion-worker/tests/test_concurrency.py`
Expected: FAIL

**Step 3: Write minimal implementation**
```python
# Ensure semaphore is correctly initialized and used in main.py
```

**Step 4: Verify test passes**
Run: `pytest apps/ingestion-worker/tests/test_concurrency.py`
Expected: PASS
