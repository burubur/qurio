import logging
from crawl4ai import AsyncWebCrawler

logger = logging.getLogger(__name__)

async def handle_web_task(url: str, max_depth: int = 1, exclusions: list[str] = None) -> str:
    """
    Crawls a website and returns the markdown content.
    Currently only supports single page crawling properly as crawl4ai's recursive mode might need more config.
    For MVP, we stick to what the library provides easily.
    """
    logger.info(f"Starting crawl for {url} with depth {max_depth}")
    
    # Initialize crawler
    # Note: crawl4ai configuration might change based on version. 
    # Using basic context manager as per documentation.
    async with AsyncWebCrawler(verbose=True) as crawler:
        result = await crawler.arun(url=url)
        
        if not result.success:
            logger.error(f"Crawl failed for {url}: {result.error_message}")
            raise Exception(f"Crawl failed: {result.error_message}")
            
        return result.markdown
