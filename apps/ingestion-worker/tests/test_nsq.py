import pytest
from unittest.mock import MagicMock, patch
import json
import main  # import the module

@pytest.mark.asyncio
async def test_process_message_success():
    # Arrange
    message = MagicMock()
    message.body = json.dumps({"type": "web", "url": "http://example.com", "id": "123"}).encode('utf-8')
    message.finish = MagicMock()
    message.requeue = MagicMock()

    # Mock handlers
    with patch('main.handle_web_task', new_callable=MagicMock) as mock_web_task:
        # Make it awaitable
        async def async_mock(*args, **kwargs):
            return "content"
        mock_web_task.side_effect = async_mock
        
        # Mock producer
        mock_producer = MagicMock()
        main.producer = mock_producer
        
        # Act
        await main.process_message(message)

        # Assert
        message.finish.assert_called_once()
        message.requeue.assert_not_called()
        mock_producer.pub.assert_called()

@pytest.mark.asyncio
async def test_process_message_failure():
    # Arrange
    message = MagicMock()
    message.body = b"invalid json"
    message.finish = MagicMock()
    message.requeue = MagicMock()

    # Act
    await main.process_message(message)

    # Assert
    message.finish.assert_not_called()
    message.requeue.assert_called_once()