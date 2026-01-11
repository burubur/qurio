import sys
from unittest.mock import MagicMock, AsyncMock, patch
import pytest

@pytest.fixture
def mock_crawl4ai_env():
    # Setup mocks for crawl4ai
    mock_crawl4ai = MagicMock()
    sys.modules["crawl4ai"] = mock_crawl4ai
    sys.modules["crawl4ai.content_filter_strategy"] = MagicMock()
    sys.modules["crawl4ai.markdown_generation_strategy"] = MagicMock()
    
    # Force reload of handlers.web to pick up these mocks
    if 'handlers.web' in sys.modules:
        del sys.modules['handlers.web']
    import handlers.web
    
    yield handlers.web
    
    # Cleanup
    if 'handlers.web' in sys.modules:
        del sys.modules['handlers.web']

@pytest.mark.asyncio
async def test_llms_txt_uses_default_generator(mock_crawl4ai_env):
    handlers_web = mock_crawl4ai_env
    handle_web_task = handlers_web.handle_web_task
    
    url = "https://example.com/llms.txt"
    mock_result = MagicMock()
    mock_result.success = True
    mock_result.markdown = "content"
    mock_result.url = url
    mock_result.links = {'internal': []}

    mock_crawler = MagicMock()
    async def fake_arun(url, config=None):
        return mock_result
    mock_crawler.arun.side_effect = fake_arun
        
    mock_crawler_cm = AsyncMock()
    mock_crawler_cm.__aenter__.return_value = mock_crawler
    mock_crawler_cm.__aexit__.return_value = None
    
    mock_factory = MagicMock(return_value=mock_crawler_cm)
    
    # Patch DefaultMarkdownGenerator IN the reloaded module
    with patch.object(handlers_web, 'DefaultMarkdownGenerator') as MockGen:
        generator_instance = MagicMock(name="generator_instance")
        MockGen.return_value = generator_instance
        
        await handle_web_task(url, crawler_factory=mock_factory)
        
        # Verify via CrawlerRunConfig calls
        MockCrawlerRunConfig = sys.modules['crawl4ai'].CrawlerRunConfig
        
        calls = MockCrawlerRunConfig.call_args_list
        found = False
        for call in calls:
            if call.kwargs.get('markdown_generator') == generator_instance:
                found = True
                break
        assert found, "CrawlerRunConfig should be initialized with DefaultMarkdownGenerator instance"

@pytest.mark.asyncio
async def test_standard_page_uses_llm_filter(mock_crawl4ai_env):
    handlers_web = mock_crawl4ai_env
    handle_web_task = handlers_web.handle_web_task
    
    url = "https://example.com/page"
    mock_result = MagicMock()
    mock_result.success = True
    mock_result.markdown = "content"
    mock_result.url = url
    mock_result.links = {'internal': []}

    mock_crawler = MagicMock()
    async def fake_arun(url, config=None):
        if "llms.txt" in url:
             m_res = MagicMock()
             m_res.success = False
             return m_res
        return mock_result

    mock_crawler.arun.side_effect = fake_arun
    
    mock_crawler_cm = AsyncMock()
    mock_crawler_cm.__aenter__.return_value = mock_crawler
    mock_crawler_cm.__aexit__.return_value = None
    
    mock_factory = MagicMock(return_value=mock_crawler_cm)
    
    with patch.object(handlers_web, 'DefaultMarkdownGenerator') as MockGen:
        await handle_web_task(url, crawler_factory=mock_factory)
        
        calls = MockGen.call_args_list
        assert len(calls) > 0
        
        has_filter = False
        for call in calls:
            if 'content_filter' in call.kwargs:
                has_filter = True
                break
        assert has_filter, "DefaultMarkdownGenerator should be called with content_filter for standard pages"