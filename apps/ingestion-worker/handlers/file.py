import logging
import asyncio
from concurrent.futures import ThreadPoolExecutor
from docling.document_converter import DocumentConverter

logger = logging.getLogger(__name__)

# Initialize converter globally to reuse resources
converter = DocumentConverter()
executor = ThreadPoolExecutor(max_workers=2)

async def handle_file_task(file_path: str) -> str:
    """
    Converts a document to markdown using Docling.
    Executes blocking code in a thread pool.
    """
    logger.info(f"Starting conversion for {file_path}")
    
    loop = asyncio.get_running_loop()
    
    try:
        # Run synchronous convert method in thread pool
        result = await loop.run_in_executor(
            executor,
            converter.convert,
            file_path
        )
        
        return result.document.export_to_markdown()
        
    except Exception as e:
        logger.error(f"Conversion failed for {file_path}: {e}")
        raise Exception(f"Conversion failed: {e}")
