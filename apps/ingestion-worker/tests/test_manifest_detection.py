import sys
from unittest.mock import MagicMock
import types

# Create a mock package for crawl4ai if not already present
if "crawl4ai" not in sys.modules:
    crawl4ai = types.ModuleType("crawl4ai")
    sys.modules["crawl4ai"] = crawl4ai

    # Mock submodules
    content_filter_strategy = types.ModuleType("crawl4ai.content_filter_strategy")
    sys.modules["crawl4ai.content_filter_strategy"] = content_filter_strategy
    crawl4ai.content_filter_strategy = content_filter_strategy

    markdown_generation_strategy = types.ModuleType("crawl4ai.markdown_generation_strategy")
    sys.modules["crawl4ai.markdown_generation_strategy"] = markdown_generation_strategy
    crawl4ai.markdown_generation_strategy = markdown_generation_strategy

    # Populate with mocks
    crawl4ai.AsyncWebCrawler = MagicMock()
    crawl4ai.CrawlerRunConfig = MagicMock()
    crawl4ai.CacheMode = MagicMock()
    crawl4ai.LLMConfig = MagicMock()

    content_filter_strategy.PruningContentFilter = MagicMock()
    content_filter_strategy.LLMContentFilter = MagicMock()

    markdown_generation_strategy.DefaultMarkdownGenerator = MagicMock()

import pytest
from unittest.mock import MagicMock, AsyncMock, patch, ANY
import asyncio
from handlers.web import handle_web_task

@pytest.mark.asyncio
async def test_detects_and_merges_llms_txt_links():
    # Arrange
    url = "https://example.com/home"
    
    mock_crawler = MagicMock()
    mock_crawler_cm = AsyncMock()
    mock_crawler_cm.__aenter__.return_value = mock_crawler
    mock_crawler_cm.__aexit__.return_value = None
    
    # Setup side effects for arun to simulate two calls:
    async def mock_arun(url, config):
        res = MagicMock()
        res.success = True
        res.links = {}
        if url.endswith("llms.txt"):
            res.markdown = "[Manifest Link](https://example.com/manifest-dest)"
            # Simulate extract_web_metadata result for links
            # But wait, handle_web_task calls extract_web_metadata
            # We need to ensure extract_web_metadata works or mock it
        else:
            res.markdown = "[Main Link](https://example.com/main-dest)"
        res.url = url
        return res
        
    mock_crawler.arun.side_effect = mock_arun
    
    # We also need to mock extract_web_metadata to return links we expect
    with patch('handlers.web.extract_web_metadata') as mock_extract:
        def extract_side_effect(result, url):
            if "manifest-dest" in result.markdown:
                return {"links": ["https://example.com/manifest-dest"], "title": "Manifest", "path": "/llms.txt"}
            return {"links": ["https://example.com/main-dest"], "title": "Main", "path": "/home"}
        
        mock_extract.side_effect = extract_side_effect
        
        with patch('handlers.web.AsyncWebCrawler', return_value=mock_crawler_cm):
            # Act
            result = await handle_web_task(url)
            
            # Assert
            assert isinstance(result, list)
            # Expect TWO results now: one for manifest, one for main page
            assert len(result) == 2
            
            # Check Result 1 (Manifest - usually appended first if found, or second? 
            # Looking at code: manifest appended FIRST, then main page appended.
            
            manifest_res = result[0]
            assert "llms.txt" in manifest_res["url"]
            assert "https://example.com/manifest-dest" in manifest_res["links"]
            
            main_res = result[1]
            assert "home" in main_res["url"]
            assert "https://example.com/main-dest" in main_res["links"]

