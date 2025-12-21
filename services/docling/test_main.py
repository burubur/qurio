from unittest.mock import MagicMock
import sys
import os

# Mock docling before importing main
mock_docling = MagicMock()
sys.modules["docling"] = mock_docling
sys.modules["docling.document_converter"] = MagicMock()

from fastapi.testclient import TestClient
from main import app, converter # Import converter to mock its method

client = TestClient(app)

def test_process_document():
    # Mock the converter's convert method
    mock_result = MagicMock()
    mock_result.document.export_to_markdown.return_value = "processed text"
    converter.convert.return_value = mock_result

    # Create a dummy file
    with open("test.txt", "w") as f:
        f.write("Hello World")
    
    try:
        with open("test.txt", "rb") as f:
            files = {'file': ('test.txt', f, 'text/plain')}
            response = client.post("/process", files=files)
            
        assert response.status_code == 200
        json_response = response.json()
        assert "text" in json_response
        assert json_response["text"] == "processed text"
    finally:
        if os.path.exists("test.txt"):
            os.remove("test.txt")
