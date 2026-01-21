import pytest
from config import Settings
from unittest.mock import patch, AsyncMock, MagicMock
from importlib import reload
import sys

def test_settings_timeout_default():
    s = Settings()
    # Default is 60000ms (60s)
    assert s.crawler_page_timeout == 60000

@pytest.mark.asyncio
async def test_web_handler_uses_timeout():
    # Reload handlers.web to ensure we are testing a fresh module instance
    # independent of other tests that might have messed with sys.modules
    if 'handlers.web' in sys.modules:
        import handlers.web
        reload(handlers.web)
    else:
        import handlers.web
    
    from handlers.web import handle_web_task
    
    # Patch settings to return a custom timeout
    with patch('handlers.web.app_settings') as mock_settings:
        mock_settings.crawler_page_timeout = 120000
        mock_settings.gemini_api_key = "fake"
        
        # Patch the crawler factory and config using patch.object
        with patch.object(handlers.web, 'default_crawler_factory') as mock_factory, \
             patch.object(handlers.web, 'CrawlerRunConfig') as MockCrawlerRunConfig:
            
            mock_instance = AsyncMock()
            mock_factory.return_value.__aenter__.return_value = mock_instance
            
            mock_result = MagicMock()
            mock_result.success = True
            mock_result.markdown = "test"
            mock_result.url = "http://example.com"
            mock_result.links = {}
            mock_instance.arun.return_value = mock_result
            
            # We mock asyncio.wait_for to avoid actual waiting if logic uses it
            async def mock_wait_for_impl(awaitable, timeout):
                return await awaitable

            with patch('asyncio.wait_for', side_effect=mock_wait_for_impl) as mock_wait:
                await handle_web_task("http://example.com", crawler_factory=mock_factory)
                
                # Verify CrawlerRunConfig was called with page_timeout=120000
                assert MockCrawlerRunConfig.call_count >= 1
                
                found = False
                for call in MockCrawlerRunConfig.call_args_list:
                    if call.kwargs.get('page_timeout') == 120000:
                        found = True
                        break
                
                assert found, f"CrawlerRunConfig was not called with page_timeout=120000. Calls: {MockCrawlerRunConfig.call_args_list}"