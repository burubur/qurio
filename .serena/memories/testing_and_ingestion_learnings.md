
## Performance & Reliability Optimization (Jan 21 2026)

**Issue:**
High concurrency ingestion (50 threads x 10 replicas) caused severe resource exhaustion, `Page.goto` timeouts, and NSQ connection drops ("Stream is closed"). The worker was spawning a new browser instance for *every* URL, and blocking the event loop during heavy parsing.

**Architectural Changes:**
1.  **Persistent Browser Instance:** Refactored `apps/ingestion-worker/main.py` to initialize a single global `AsyncWebCrawler` (browser) at startup. This instance is passed to `handle_web_task` and reused for all requests, eliminating the overhead of launching/killing browser processes per page.
2.  **Removed Redundant Checks:** Removed the per-page probe for `llms.txt`.
3.  **Increased Timeouts:** 
    - `CRAWLER_PAGE_TIMEOUT`: Increased from 60s to 120s to handle heavy documentation pages under load.
    - `NSQ_HEARTBEAT_INTERVAL`: Increased from 30s to 60s in `nsq.Reader` to prevent connection drops when the Python event loop is blocked by CPU-intensive parsing (e.g., large single-page documentation).

**Concurrency Clarification:**
- **Scraping Concurrency:** Controlled by `NSQ_MAX_IN_FLIGHT` (per worker replica).
- **Backend Processing:** Controlled by `INGESTION_CONCURRENCY` (per backend replica).
- These are independent variables and should be tuned separately.
