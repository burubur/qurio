### llms.txt Ingestion Strategy (2026-01-10)
- **The "Invisible Door" Problem:**
    - **Issue:** Standard web crawlers rely on explicit `<a>` tags. `llms.txt` is rarely linked in the HTML footer or navigation, making it "invisible" to standard crawling logic regardless of depth settings.
    - **Fix:** Implemented **Active Probing** in the worker. The crawler now preemptively attempts to fetch `/llms.txt` at the root of the domain (e.g., `https://example.com/llms.txt`) for every web task, merging it into the result stream if found.
- **The Depth/Noise Tradeoff:**
    - **Issue:** To crawl links *inside* a manifest, users had to increase the Global Depth (e.g., to 2). This inadvertently forced the crawler to process "junk" links (About, Contact, Legal) on the main page, wasting resources and diluting the context.
    - **Fix:** Implemented **Virtual Depth** logic in the Backend.
        - When the backend receives a result for `.../llms.txt`, it locally increments the `MaxDepth` by +1 (e.g., treats `MaxDepth: 0` as `1`).
        - **Result:** The crawler respects the strict depth limit for the homepage (ignoring junk) but grants a "Free Pass" to the manifest to unfold its technical links. This allows targeted "Homepage + Manifest" ingestion without over-crawling.