from fastapi import FastAPI, UploadFile, File
from docling.document_converter import DocumentConverter
import tempfile
import os

app = FastAPI()
converter = DocumentConverter()

@app.post("/process")
async def process(file: UploadFile = File(...)):
    # Save temp file
    suffix = ""
    if file.filename:
        suffix = os.path.splitext(file.filename)[1]
        
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        content = await file.read()
        tmp.write(content)
        tmp_path = tmp.name

    try:
        # Run docling
        result = converter.convert(tmp_path)
        markdown = result.document.export_to_markdown()
        return {"text": markdown}
    finally:
        if os.path.exists(tmp_path):
            os.remove(tmp_path)

@app.get("/")
def read_root():
    return {"status": "ok"}