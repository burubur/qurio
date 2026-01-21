## facing a bunch of issues and failure when scraping pinia.vuejs.org seems like the origin of the errors are the same, here are 3 random samples for you to check further
it started from 100 failure ingestion, I retry 1 by one, down to 70, then retry again down to 40, retry again down to 8, and retry again, finally all is well

### Example 1
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://pinia.vuejs.org/api/pinia/functions/getActivePinia.html", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

### Example 2
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://pinia.vuejs.org/cookbook/", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

### Example 3
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://pinia.vuejs.org/zh/core-concepts", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

## I tried to ingest VueJS documentation also getting similar error

### Example 1
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://vuejs.org/rules-essential", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

### Example 2
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://vuejs.org/guide/quick-start.md", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

### Example 3
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://vuejs.org/api/component-instance.md", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its

### Example 4
Crawl failed: Unexpected error in _crawl_web at line 718 in _crawl_web (../usr/local/lib/python3.12/site-packages/crawl4ai/async_crawler_strategy.py): Error: Failed on navigating ACS-GOTO: Page.goto: Timeout 60000ms exceeded. Call log: - navigating to "https://vuejs.org/api/composition-api-setup.html", waiting until "domcontentloaded" Code context: 713 tag="GOTO", 714 params={"url": url}, 715 ) 716 response = None 717 else: 718 → raise RuntimeError(f"Failed on navigating ACS-GOTO:\n{str(e)}") 719 720 # ────────────────────────────────────────────────────────────── 721 # Walk the redirect chain. Playwright returns only the last 722 # hop, so we trace the `request.redirected_from` links until the 723 # first response that differs from the final one and surface its


## for extra information here is the worker scaling configuration in the .env file

INGESTION_CONCURRENCY=50
INGESTION_WORKER_WEB_REPLICAS=10
INGESTION_WORKER_FILE_REPLICAS=3
BACKEND_WORKER_REPLICAS=5