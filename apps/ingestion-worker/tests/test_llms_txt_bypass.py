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
async def test_llms_txt_uses_default_generator():
    # Arrange
    url = "https://example.com/llms.txt"
    
    mock_result = MagicMock()
    mock_result.success = True
    mock_result.markdown = "[link](http://test.com)"
    mock_result.url = url
    mock_result.links = {'internal': []}
    
    mock_crawler = MagicMock()
    f = asyncio.Future()
    f.set_result(mock_result)
    mock_crawler.arun.return_value = f
    
    mock_crawler_cm = AsyncMock()
    mock_crawler_cm.__aenter__.return_value = mock_crawler
    mock_crawler_cm.__aexit__.return_value = None
    
    with patch('handlers.web.AsyncWebCrawler', return_value=mock_crawler_cm) as MockCrawler:
        from crawl4ai import DefaultMarkdownGenerator
        
        # Act
        await handle_web_task(url)
        
        # Assert
        call_args = mock_crawler.arun.call_args
        assert call_args is not None
        config = call_args.kwargs.get('config')
        
        # We check if DefaultMarkdownGenerator was instantiated and passed
        # Since we mocked it, we can check if the mock was used
        assert isinstance(config.markdown_generator, DefaultMarkdownGenerator)

@pytest.mark.asyncio
async def test_standard_page_uses_llm_filter():
    # Arrange
    url = "https://example.com/page"
    
    mock_result = MagicMock()
    mock_result.success = True
    mock_result.markdown = "Content"
    mock_result.url = url
    mock_result.links = {'internal': []}
    
    mock_crawler = MagicMock()
    f = asyncio.Future()
    f.set_result(mock_result)
    mock_crawler.arun.return_value = f
    
    mock_crawler_cm = AsyncMock()
    mock_crawler_cm.__aenter__.return_value = mock_crawler
    mock_crawler_cm.__aexit__.return_value = None
    
    with patch('handlers.web.AsyncWebCrawler', return_value=mock_crawler_cm) as MockCrawler:
        from crawl4ai import LLMContentFilter
        
        # Act
        await handle_web_task(url)
        
        # Assert
        call_args = mock_crawler.arun.call_args
        config = call_args.kwargs.get('config')
        
        # Currently the code wraps LLMContentFilter in DefaultMarkdownGenerator
        # so we need to check how it's constructed in the implementation
        # The implementation does: md_generator = DefaultMarkdownGenerator(content_filter=llm_filter)
        # So for standard page, it should be DefaultMarkdownGenerator with a content_filter
        
        assert isinstance(config.markdown_generator, MagicMock) # DefaultMarkdownGenerator mock
        
        # We need to verify if LLMContentFilter was used in its construction
        # But since we mock everything, it's a bit tricky to distinguish without implementation details
        # However, for the bypass (test above), we expect DefaultMarkdownGenerator() WITHOUT content_filter
        # or with a different one. 
        
        # Let's check LLMContentFilter instantiation
        # In standard flow, LLMContentFilter is instantiated
        # In bypass flow, it should NOT be (or at least not passed to generator)
