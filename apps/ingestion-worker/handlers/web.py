import asyncio
import structlog
import re
from urllib.parse import urljoin, urlparse
from crawl4ai import AsyncWebCrawler, CrawlerRunConfig, CacheMode, LLMConfig
from crawl4ai.content_filter_strategy import LLMContentFilter
from crawl4ai.markdown_generation_strategy import DefaultMarkdownGenerator
from config import settings as app_settings

logger = structlog.get_logger(__name__)

INSTRUCTION = """
    Extract technical content from this software documentation page.
    
    KEEP:
    - All code examples with their comments
    - Function/method signatures and parameters
    - Configuration examples and syntax
    - Technical explanations and concepts
    - Error messages and troubleshooting steps
    - Links to related API documentation
    
    REMOVE:
    - Navigation menus and sidebars
    - Copyright and legal notices
    - Unrelated marketing content
    - "Edit this page" links
    - Cookie banners and consent forms
    
    PRESERVE:
    - Code block language annotations (```go, etc.)
    - Heading hierarchy for context
    - Inline code references
    - Numbered lists for sequential steps
"""

def extract_web_metadata(result, url: str) -> dict:
    """
    Extracts metadata (title, path, links) from a crawl result.
    """
    # Extract internal links
    # Crawl4AI result.links is usually a dictionary with 'internal' and 'external' keys
    # containing lists of dicts (href, text, etc.)
    internal_links = []
    if result.links and 'internal' in result.links:
            for link in result.links['internal']:
                if 'href' in link:
                    internal_links.append(link['href'])
    
    # Additional Regex Extraction for Markdown (e.g. llms.txt)
    if result.markdown:
        markdown_links = re.findall(r'\[.*?\]\((.*?)\)', result.markdown)
        parsed_base = urlparse(url)
        base_domain = parsed_base.netloc
        
        for link in markdown_links:
            # Resolve relative URLs
            full_url = urljoin(url, link)
            # Filter internal
            if urlparse(full_url).netloc == base_domain:
                internal_links.append(full_url)

    # De-duplicate
    internal_links = list(set(internal_links))

    # Extract title (simplistic regex fallback if not in result)
    title = ""
    if result.markdown:
        match = re.search(r'^#\s+(.+)$', result.markdown, re.MULTILINE)
        if match:
            title = match.group(1).strip()
    
    # Extract path (breadcrumbs)
    parsed_url = urlparse(result.url)
    path_segments = [s for s in parsed_url.path.split('/') if s]
    path_str = " > ".join(path_segments)
    
    return {
        "title": title,
        "path": path_str,
        "links": internal_links
    }

async def handle_web_task(url: str, api_key: str = None) -> dict:
    """
    Crawls a single page and returns content and discovered internal links.
    """
    logger.info("crawl_starting", url=url)
        
    # Use passed api_key or fallback to settings
    token = api_key if api_key else app_settings.gemini_api_key
    
    # 1. Manifest Detection (llms.txt)
    results = []
    
    if not url.endswith("llms.txt"):
        try:
            parsed = urlparse(url)
            base_url = f"{parsed.scheme}://{parsed.netloc}"
            manifest_url = f"{base_url}/llms.txt"
            
            # Lightweight check/crawl for manifest
            manifest_config = CrawlerRunConfig(
                markdown_generator=DefaultMarkdownGenerator(),
                cache_mode=CacheMode.ENABLED
            )
            
            async with AsyncWebCrawler(verbose=False) as manifest_crawler:
                manifest_res = await asyncio.wait_for(
                    manifest_crawler.arun(url=manifest_url, config=manifest_config),
                    timeout=10.0 # Short timeout for manifest check
                )
                
                if manifest_res.success and manifest_res.markdown:
                    logger.info("manifest_found", url=manifest_url)
                    manifest_meta = extract_web_metadata(manifest_res, manifest_url)
                    
                    results.append({
                        "url": manifest_res.url,
                        "title": manifest_meta['title'] or "llms.txt",
                        "path": manifest_meta['path'],
                        "content": manifest_res.markdown,
                        "links": manifest_meta['links']
                    })
        except Exception as e:
            # Silent failure for manifest detection is acceptable
            logger.debug("manifest_check_failed", error=str(e))

    # 2. Configure Generator (Bypass LLM for .txt/llms.txt)
    if url.endswith(".txt") or url.endswith("llms.txt"):
        md_generator = DefaultMarkdownGenerator()
        logger.info("llm_bypass_enabled", url=url, reason="text_file")
    else:
        llm_config = LLMConfig(
            provider="gemini/gemini-3-flash-preview", 
            api_token=token,
            temperature=1.0
        )

        llm_filter = LLMContentFilter(
            llm_config=llm_config,
            instruction=INSTRUCTION,
            chunk_token_threshold=8000
        )
        
        md_generator = DefaultMarkdownGenerator(content_filter=llm_filter)

    config = CrawlerRunConfig(
        cache_mode=CacheMode.ENABLED,
        # Remove excluded_tags to ensure links in nav/sidebar are discovered.
        # The LLMContentFilter will handle removing them from the content.
        # excluded_tags=['nav', 'footer', 'aside', 'header'], 
        exclude_external_links=True,
        markdown_generator=md_generator,
        check_robots_txt=True 
    )
    
    # Initialize crawler
    try:
        async with AsyncWebCrawler(verbose=True) as crawler:
            # Single page crawl
            result = await asyncio.wait_for(
                crawler.arun(url=url, config=config),
                timeout=300.0
            )
            
            if not result.success:
                logger.error("crawl_failed", url=url, error=result.error_message)
                raise Exception(f"Crawl failed: {result.error_message}")
                
            meta = extract_web_metadata(result, url)

            logger.info("crawl_completed", url=url, links_found=len(meta['links']), title=meta['title'], path=meta['path'])

            results.append({
                "url": result.url,
                "title": meta['title'],
                "path": meta['path'],
                "content": result.markdown,
                "links": meta['links']
            })
            
            return results

    except asyncio.TimeoutError:
        logger.error("crawl_timeout", url=url)
        raise
    except Exception as e:
        logger.error("crawl_exception", url=url, error=str(e))
        raise
