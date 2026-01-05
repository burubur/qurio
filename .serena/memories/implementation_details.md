Plan updated with clean formatting and explicit target descriptions.
`qurio_list_pages` uses `source_id`.
`qurio_read_page` uses `url`.
Full argument guide included for `qurio_search`.

## Ingestion Worker
- `handle_file_task` now returns `list[dict]` (standardized with `handle_web_task`) to simplify `main.py` logic.
- Metadata (Title, Author, PageCount) is extracted via Docling v2 and passed to backend.