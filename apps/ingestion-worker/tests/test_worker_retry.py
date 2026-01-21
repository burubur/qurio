import pytest
from unittest.mock import MagicMock, patch, AsyncMock
import asyncio
# We mock main imports to avoid side effects if main() runs
# But we need process_message.
# Assuming main.py can be imported without running main() because of if __name__ == "__main__":
from main import process_message

MAX_RETRIES = 3

@pytest.mark.asyncio
async def test_requeue_on_timeout():
    mock_msg = MagicMock()
    mock_msg.body = b'{"type": "web", "url": "http://fail.com", "id": "1"}'
    mock_msg.attempts = 1
    
    # Simulate TimeoutError in handle_web_task
    with patch('main.handle_web_task', side_effect=asyncio.TimeoutError("Timeout")):
        with patch('main.producer') as mock_producer, \
             patch('main.WORKER_SEMAPHORE', asyncio.Semaphore(1)):
            await process_message(mock_msg)
            
            # Should NOT finish, should requeue
            mock_msg.finish.assert_not_called()
            mock_msg.requeue.assert_called()
            
            # Check delay (attempt * 30 -> 1 * 30 = 30)
            # Assuming requeue(delay=...) or similar
            args, kwargs = mock_msg.requeue.call_args
            # Verify delay is passed (either as arg 0 or kwarg 'delay')
            # Adjust based on pynsq API findings, but for now assuming delay in seconds/ms
            assert kwargs.get('delay') == 30 or (args and args[0] == 30) or kwargs.get('delay') == 30000

@pytest.mark.asyncio
async def test_fail_on_max_retries():
    mock_msg = MagicMock()
    mock_msg.body = b'{"type": "web", "url": "http://fail.com", "id": "1"}'
    mock_msg.attempts = 4 # Max is 3
    
    with patch('main.handle_web_task', side_effect=asyncio.TimeoutError("Timeout")):
        with patch('main.producer') as mock_producer, \
             patch('main.WORKER_SEMAPHORE', asyncio.Semaphore(1)):
            await process_message(mock_msg)
            
            # Should finish and publish failure
            mock_msg.finish.assert_called()
            mock_msg.requeue.assert_not_called()
            mock_producer.pub.assert_called() # Publish failure
