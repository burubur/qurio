SCRAPE].. ◆ https://vitest.dev/config/coverage#coverage-processingconcurrency
| ✓ | ⏱: 10.38s
[COMPLETE] ● https://vitest.dev/config/coverage#coverage-processingconcurrency
| ✓ | ⏱: 111.61s
2025-12-27 01:05:52,083 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
2025-12-27 01:05:52,084 - WARNING - [nsqd:4150] connection closed
2025-12-27 01:05:52,084 - INFO - [nsqd:4150] attempting to reconnect in 15.00s
2025-12-27 01:05:52,084 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa012' (Stream is closed)
2025-12-27 01:05:52,084 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa013' (Stream is closed)
2025-12-27 01:05:52,084 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa014' (Stream is closed)
2025-12-27 01:05:52,085 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
2025-12-27 01:05:52,085 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:05:52,085 - ERROR - [nsqd:4150:nsqd:4150] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 01:05:52,088 - WARNING - [ingest.task:worker] lookupd http://nsqlookupd:4161/lookup?topic=ingest.task query error: Timeout while connecting

---

SCRAPE].. ◆ https://vitest.dev/llms-full.txt#reporters
| ✓ | ⏱: 45.04s
[COMPLETE] ● https://vitest.dev/llms-full.txt#reporters
| ✓ | ⏱: 46.47s
2025-12-27 01:04:46,319 - INFO - [nsqd:4150] connecting to nsqd
2025-12-27 01:04:46,325 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa012' (Stream is closed)
2025-12-27 01:04:46,325 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa013' (Stream is closed)
2025-12-27 01:04:46,325 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa014' (Stream is closed)
2025-12-27 01:04:46,326 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
2025-12-27 01:04:46,326 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
2025-12-27 01:04:46,327 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
[FETCH]... ↓ https://vitest.dev/llms-full.txt#default-reporter
| ✓ | ⏱: 45.86s

---

SCRAPE].. ◆ https://vitest.dev/llms-full.txt#has
| ✓ | ⏱: 45.58s
[COMPLETE] ● https://vitest.dev/llms-full.txt#has
| ✓ | ⏱: 140.75s
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa009' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa00a' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa00f' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa010' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa011' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa012' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa013' (Stream is closed)
2025-12-27 01:03:59,530 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa014' (Stream is closed)
2025-12-27 01:03:59,530 - WARNING - [nsqd:4150:ingest.task:worker] connection is stale (91.12s), closing
2025-12-27 01:03:59,530 - WARNING - [nsqd:4150:nsqd:4150] connection is stale (91.12s), closing
2025-12-27 01:03:59,531 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:03:59,531 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:03:59,531 - WARNING - [nsqd:4150:ingest.task:worker] connection closed
2025-12-27 01:03:59,531 - WARNING - [nsqd:4150] connection closed
2025-12-27 01:03:59,531 - INFO - [nsqd:4150] attempting to reconnect in 15.00s
2025-12-27 01:03:59,533 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:03:59,534 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 01:03:59,534 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:03:59,534 - ERROR - [nsqd:4150:nsqd:4150] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 01:03:59,536 - INFO - [nsqd:4150:ingest.task:worker] connecting to nsqd
2025-12-27 01:03:59,536 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
2025-12-27 01:03:59,537 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
2025-12-27 01:03:59,538 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#tobecallablewith"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.538569Z"}
2025-12-27 01:03:59,538 - INFO - {"url": "https://vitest.dev/llms-full.txt#tobecallablewith", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.538636Z"}
2025-12-27 01:03:59,612 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#getbylabeltext"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.612567Z"}
2025-12-27 01:03:59,612 - INFO - {"url": "https://vitest.dev/llms-full.txt#getbylabeltext", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.612766Z"}
2025-12-27 01:03:59,687 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#default-reporter"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.687730Z"}
2025-12-27 01:03:59,687 - INFO - {"url": "https://vitest.dev/llms-full.txt#default-reporter", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.687928Z"}
2025-12-27 01:03:59,764 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#vi-unmock"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.764600Z"}
2025-12-27 01:03:59,764 - INFO - {"url": "https://vitest.dev/llms-full.txt#vi-unmock", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.764804Z"}
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:03:59,845 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/coverage#coverage-enabled"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.844989Z"}
2025-12-27 01:03:59,845 - INFO - {"url": "https://vitest.dev/config/coverage#coverage-enabled", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.845177Z"}
2025-12-27 01:03:59,920 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/coverage#coverage-processingconcurrency"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.920532Z"}
2025-12-27 01:03:59,920 - INFO - {"url": "https://vitest.dev/config/coverage#coverage-processingconcurrency", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.920721Z"}
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:03:59,997 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:03:59.997178Z"}
2025-12-27 01:03:59,997 - INFO - {"url": "https://vitest.dev/guide/", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:03:59.997353Z"}
2025-12-27 01:04:00,072 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#isolation-strategy"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:04:00.072550Z"}
2025-12-27 01:04:00,072 - INFO - {"url": "https://vitest.dev/llms-full.txt#isolation-strategy", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:04:00.072732Z"}
2025-12-27 01:04:00,150 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/api/browser/interactivity#userevent-dblclick"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:04:00.150372Z"}
2025-12-27 01:04:00,150 - INFO - {"url": "https://vitest.dev/api/browser/interactivity#userevent-dblclick", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:04:00.150547Z"}
2025-12-27 01:04:00,223 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/test-annotations"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:04:00.223661Z"}
2025-12-27 01:04:00,223 - INFO - {"url": "https://vitest.dev/guide/test-annotations", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:04:00.223855Z"}
2025-12-27 01:04:00,323 - INFO - {"url": "https://vitest.dev/llms-full.txt#invalidatefile", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.323072Z"}
2025-12-27 01:04:00,350 - INFO - {"url": "https://vitest.dev/llms-full.txt#has", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.350354Z"}
2025-12-27 01:04:00,355 - INFO - {"url": "https://vitest.dev/config/coverage#coverage-thresholds-autoupdate", "links_found": 39, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.355049Z"}
2025-12-27 01:04:00,364 - INFO - {"url": "https://vitest.dev/api/", "links_found": 49, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.364850Z"}
2025-12-27 01:04:00,373 - INFO - {"url": "https://vitest.dev/guide/features.html#benchmarking", "links_found": 45, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.373007Z"}
2025-12-27 01:04:00,376 - INFO - {"url": "https://vitest.dev/comparisons", "links_found": 3, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.376826Z"}
2025-12-27 01:04:00,414 - INFO - {"url": "https://vitest.dev/api/advanced/test-suite#children", "links_found": 28, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:04:00.414591Z"}
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:04:00,532 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#invalidatefile", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.532209Z"}
2025-12-27 01:04:00,532 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa010' (Stream is closed)
2025-12-27 01:04:00,533 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/guide/features.html#benchmarking", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.533164Z"}
2025-12-27 01:04:00,533 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa009' (Stream is closed)
2025-12-27 01:04:00,536 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#has", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.536024Z"}
2025-12-27 01:04:00,536 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa00c' (Stream is closed)
2025-12-27 01:04:00,536 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/api/", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.536439Z"}
2025-12-27 01:04:00,536 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa011' (Stream is closed)
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:04:00,548 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/api/advanced/test-suite#children", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.548205Z"}
2025-12-27 01:04:00,548 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa00f' (Stream is closed)
2025-12-27 01:04:00,582 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/coverage#coverage-thresholds-autoupdate", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.582240Z"}
2025-12-27 01:04:00,582 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa00d' (Stream is closed)
2025-12-27 01:04:00,588 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/comparisons", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:04:00.588234Z"}
2025-12-27 01:04:00,588 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa00a' (Stream is closed)
[INIT].... → Crawl4AI 0.7.8
[FETCH]... ↓ https://vitest.dev/llms-full.txt#reporters
| ✓ | ⏱: 0.91s
01:04:01 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

[SCRAPE].. ◆ https://vitest.dev/guide/cli
| ✓ | ⏱: 10.75s
[COMPLETE] ● https://vitest.dev/guide/cli
| ✓ | ⏱: 118.83s
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa009' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa00a' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa000' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa001' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa002' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa004' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa005' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa003' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa006' (Stream is closed)
2025-12-27 01:01:07,826 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa007' (Stream is closed)
2025-12-27 01:01:07,827 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/coverage#coverage-thresholds-autoupdate"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:01:07.827264Z"}
2025-12-27 01:01:07,827 - INFO - {"url": "https://vitest.dev/config/coverage#coverage-thresholds-autoupdate", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:01:07.827377Z"}
[FETCH]... ↓ https://vitest.dev/guide/improving-performance
| ✓ | ⏱: 110.86s
01:01:07 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
2025-12-27 01:01:07,880 - INFO -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
01:01:07 - LiteLLM:INFO: vertex_and_google_ai_studio_gemini.py:886 - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
2025-12-27 01:01:07,881 - INFO - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
01:01:16 - LiteLLM:INFO: utils.py:1331 - Wrapper: Completed Call, calling success_handler
2025-12-27 01:01:16,899 - INFO - Wrapper: Completed Call, calling success_handler
[SCRAPE].. ◆ https://vitest.dev/guide/improving-performance
| ✓ | ⏱: 9.04s
[COMPLETE] ● https://vitest.dev/guide/improving-performance
| ✓ | ⏱: 127.95s
[FETCH]... ↓ https://vitest.dev/api/browser/locators
| ✓ | ⏱: 119.90s
01:01:16 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
2025-12-27 01:01:16,967 - INFO -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
01:01:16 - LiteLLM:INFO: vertex_and_google_ai_studio_gemini.py:886 - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
2025-12-27 01:01:16,968 - INFO - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
01:01:27 - LiteLLM:INFO: utils.py:1331 - Wrapper: Completed Call, calling success_handler
2025-12-27 01:01:27,254 - INFO - Wrapper: Completed Call, calling success_handler
[SCRAPE].. ◆ https://vitest.dev/api/browser/locators
| ✓ | ⏱: 10.35s
[COMPLETE] ● https://vitest.dev/api/browser/locators
| ✓ | ⏱: 138.30s
2025-12-27 01:01:27,258 - INFO - [nsqd:4150] connecting to nsqd
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa009' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa00a' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa000' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa001' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa002' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa004' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa005' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa003' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa006' (Stream is closed)
2025-12-27 01:01:27,298 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa007' (Stream is closed)
2025-12-27 01:01:27,299 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/#maxconcurrency"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:01:27.299280Z"}
2025-12-27 01:01:27,299 - INFO - {"url": "https://vitest.dev/config/#maxconcurrency", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:01:27.299373Z"}
2025-12-27 01:01:27,334 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
2025-12-27 01:01:27,377 - WARNING - [ingest.task:worker] lookupd http://nsqlookupd:4161/lookup?topic=ingest.task query error: Timeout while connecting
[FETCH]... ↓ https://vitest.dev/api/mock#mockreset
| ✓ | ⏱: 130.35s
01:01:27 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
2025-12-27 01:01:27,423 - INFO -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
01:01:27 - LiteLLM:INFO: vertex_and_google_ai_studio_gemini.py:886 - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
2025-12-27 01:01:27,423 - INFO - Warning: Setting temperature < 1.0 for Gemini 3 models (gemini-3-flash-preview) can cause infinite loops, degraded reasoning performance, and failure on complex tasks. Strongly recommended to use temperature = 1.0 (default).
01:01:38 - LiteLLM:INFO: utils.py:1331 - Wrapper: Completed Call, calling success_handler
2025-12-27 01:01:38,357 - INFO - Wrapper: Completed Call, calling success_handler
[SCRAPE].. ◆ https://vitest.dev/api/mock#mockreset
| ✓ | ⏱: 10.98s
[COMPLETE] ● https://vitest.dev/api/mock#mockreset
| ✓ | ⏱: 149.40s
2025-12-27 01:01:38,361 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/api/advanced/test-suite#children"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:01:38.361055Z"}
2025-12-27 01:01:38,361 - INFO - {"url": "https://vitest.dev/api/advanced/test-suite#children", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:01:38.361164Z"}
2025-12-27 01:01:38,397 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa009' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa00a' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa000' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa001' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa002' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa004' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa005' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa003' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa006' (Stream is closed)
2025-12-27 01:01:38,398 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa007' (Stream is closed)
2025-12-27 01:01:38,437 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#invalidatefile"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:01:38.437193Z"}
2025-12-27 01:01:38,437 - INFO - {"url": "https://vitest.dev/llms-full.txt#invalidatefile", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:01:38.437367Z"}
2025-12-27 01:01:38,520 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/api/"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:01:38.520395Z"}
2025-12-27 01:01:38,520 - INFO - {"url": "https://vitest.dev/api/", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:01:38.520563Z"}
2025-12-27 01:01:38,595 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
2025-12-27 01:01:38,625 - INFO - {"url": "https://vitest.dev/llms-full.txt#grouporder", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.625051Z"}
2025-12-27 01:01:38,641 - INFO - {"url": "https://vitest.dev/llms-full.txt#vi-restoreallmocks", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.641175Z"}
2025-12-27 01:01:38,661 - INFO - {"url": "https://vitest.dev/guide/profiling-test-performance.html#code-coverage", "links_found": 16, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.661900Z"}
2025-12-27 01:01:38,663 - INFO - {"url": "https://vitest.dev/config/#threads", "links_found": 3, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.663160Z"}
2025-12-27 01:01:38,672 - INFO - {"url": "https://vitest.dev/guide/cli", "links_found": 249, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.671932Z"}
2025-12-27 01:01:38,680 - INFO - {"url": "https://vitest.dev/config/expect#expect-requireassertions", "links_found": 12, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.680684Z"}
2025-12-27 01:01:38,697 - INFO - {"url": "https://vitest.dev/api/browser/locators", "links_found": 71, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.697766Z"}
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:01:38,738 - INFO - {"url": "https://vitest.dev/api/mock#mockreset", "links_found": 36, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.738267Z"}
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:01:38,766 - INFO - {"url": "https://vitest.dev/guide/improving-performance", "links_found": 20, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.766667Z"}
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:01:38,793 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa001' (Stream is closed)
2025-12-27 01:01:38,794 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/guide/profiling-test-performance.html#code-coverage", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.794250Z"}
2025-12-27 01:01:38,804 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa005' (Stream is closed)
2025-12-27 01:01:38,806 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#vi-restoreallmocks", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.806091Z"}
2025-12-27 01:01:38,810 - INFO - {"url": "https://vitest.dev/config/browser/testerhtmlpath", "links_found": 3, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T01:01:38.810456Z"}
2025-12-27 01:01:38,816 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa007' (Stream is closed)
2025-12-27 01:01:38,816 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/#threads", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.816914Z"}
2025-12-27 01:01:38,824 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa004' (Stream is closed)
2025-12-27 01:01:38,825 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/guide/cli", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.825331Z"}
2025-12-27 01:01:38,833 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa003' (Stream is closed)
2025-12-27 01:01:38,834 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#grouporder", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.834086Z"}
[INIT].... → Crawl4AI 0.7.8
2025-12-27 01:01:38,893 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa009' (Stream is closed)
2025-12-27 01:01:38,894 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/api/browser/locators", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.894926Z"}
2025-12-27 01:01:38,896 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa002' (Stream is closed)
2025-12-27 01:01:38,896 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/guide/improving-performance", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.896946Z"}
2025-12-27 01:01:38,900 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa006' (Stream is closed)
2025-12-27 01:01:38,901 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/api/mock#mockreset", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.901488Z"}
2025-12-27 01:01:38,908 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b1217caa000' (Stream is closed)
2025-12-27 01:01:38,909 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/expect#expect-requireassertions", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.909069Z"}
2025-12-27 01:01:38,919 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa00a' (Stream is closed)
2025-12-27 01:01:38,919 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/browser/testerhtmlpath", "event": "result_published", "level": "info", "timestamp": "2025-12-27T01:01:38.919475Z"}
[INIT].... → Crawl4AI 0.7.8
[FETCH]... ↓ https://vitest.dev/config/#maxworkers
| ✓ | ⏱: 0.56s
01:01:39 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

SCRAPE].. ◆ https://vitest.dev/llms-full.txt#vi-restoreallmocks
| ✓ | ⏱: 43.70s
[COMPLETE] ● https://vitest.dev/llms-full.txt#vi-restoreallmocks
| ✓ | ⏱: 104.50s
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa000' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa001' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa002' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa004' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa005' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa003' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa006' (Stream is closed)
2025-12-27 01:00:53,453 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b1217caa007' (Stream is closed)
2025-12-27 01:00:53,453 - WARNING - [nsqd:4150:ingest.task:worker] connection is stale (108.97s), closing
2025-12-27 01:00:53,453 - WARNING - [nsqd:4150:nsqd:4150] connection is stale (96.32s), closing
2025-12-27 01:00:53,453 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:00:53,454 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:00:53,454 - WARNING - [nsqd:4150:ingest.task:worker] connection closed
2025-12-27 01:00:53,454 - WARNING - [nsqd:4150] connection closed
2025-12-27 01:00:53,454 - INFO - [nsqd:4150] attempting to reconnect in 15.00s
2025-12-27 01:00:53,455 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:00:53,456 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 01:00:53,456 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 01:00:53,456 - ERROR - [nsqd:4150:nsqd:4150] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 01:00:53,457 - INFO - [nsqd:4150:ingest.task:worker] connecting to nsqd
2025-12-27 01:00:53,457 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
2025-12-27 01:00:53,461 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
2025-12-27 01:00:53,462 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/#bail"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:00:53.462208Z"}
2025-12-27 01:00:53,462 - INFO - {"url": "https://vitest.dev/config/#bail", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:00:53.462264Z"}
2025-12-27 01:00:53,528 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/features.html#benchmarking"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:00:53.528893Z"}
2025-12-27 01:00:53,529 - INFO - {"url": "https://vitest.dev/guide/features.html#benchmarking", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:00:53.529078Z"}
2025-12-27 01:00:53,603 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/comparisons"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:00:53.603745Z"}
2025-12-27 01:00:53,603 - INFO - {"url": "https://vitest.dev/comparisons", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:00:53.603905Z"}
2025-12-27 01:00:53,677 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/#maxworkers"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:00:53.677501Z"}
2025-12-27 01:00:53,677 - INFO - {"url": "https://vitest.dev/config/#maxworkers", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:00:53.677703Z"}
2025-12-27 01:00:53,754 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#has"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T01:00:53.754623Z"}
2025-12-27 01:00:53,754 - INFO - {"url": "https://vitest.dev/llms-full.txt#has", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T01:00:53.754853Z"}
[FETCH]... ↓ https://vitest.dev/config/browser/testerhtmlpath
| ✓ | ⏱: 96.79s
01:00:53 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
2025-12-27 01:00:53,801 - INFO -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

SCRAPE].. ◆ https://tanstack.com/ai/latest/docs/contributors
| ✓ | ⏱: 8.02s
[COMPLETE] ● https://tanstack.com/ai/latest/docs/contributors
| ✓ | ⏱: 213.07s
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8ae5fccaa015' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa000' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa002' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa003' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa004' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa006' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa007' (Stream is closed)
2025-12-27 00:59:16,987 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa008' (Stream is closed)
2025-12-27 00:59:16,990 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa008' (Stream is closed)
2025-12-27 00:59:16,993 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa000' (Stream is closed)
2025-12-27 00:59:16,994 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#ontestsuiteresult", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:16.994794Z"}
2025-12-27 00:59:16,994 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#vi-domock", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:16.994933Z"}
2025-12-27 00:59:16,999 - INFO - {"url": "https://vitest.dev/config/onunhandlederror", "links_found": 5, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:16.999874Z"}
2025-12-27 00:59:17,000 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa002' (Stream is closed)
2025-12-27 00:59:17,000 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa006' (Stream is closed)
2025-12-27 00:59:17,000 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/guide/reporters#junit-reporter", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.000942Z"}
2025-12-27 00:59:17,001 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/browser/api#api-port", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.001033Z"}
2025-12-27 00:59:17,026 - INFO - {"url": "https://vitest.dev/config/#mockreset", "links_found": 2, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:17.026901Z"}
2025-12-27 00:59:17,032 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa007' (Stream is closed)
2025-12-27 00:59:17,032 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/api/vi#vi-unstuballglobals", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.032234Z"}
2025-12-27 00:59:17,039 - INFO - {"url": "https://tanstack.com/ai/latest/docs/contributors", "links_found": 193, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:17.039069Z"}
2025-12-27 00:59:17,091 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa003' (Stream is closed)
2025-12-27 00:59:17,091 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/onunhandlederror", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.091916Z"}
2025-12-27 00:59:17,119 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa004' (Stream is closed)
2025-12-27 00:59:17,120 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/config/#mockreset", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.120422Z"}
2025-12-27 00:59:17,130 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8ae5fccaa015' (Stream is closed)
2025-12-27 00:59:17,130 - INFO - {"source_id": "3b8468d0-2901-4cb0-9939-9445194d8743", "url": "https://tanstack.com/ai/latest/docs/contributors", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:17.130402Z"}
[FETCH]... ↓
https://vitest.dev/guide/profiling-test-performance.html#code-coverage
| ✓ | ⏱: 0.61s
00:59:17 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini
2025-12-27 00:59:17,636 - INFO -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

SCRAPE].. ◆ https://vitest.dev/config/onunhandlederror
| ✓ | ⏱: 2.14s
[COMPLETE] ● https://vitest.dev/config/onunhandlederror
| ✓ | ⏱: 204.18s
2025-12-27 00:59:08,617 - INFO - {"url": "https://vitest.dev/llms-full.txt#recordartifact", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.617702Z"}
2025-12-27 00:59:08,636 - INFO - {"url": "https://vitest.dev/llms-full.txt#vi-domock", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.636719Z"}
2025-12-27 00:59:08,665 - INFO - {"url": "https://vitest.dev/llms-full.txt#ontestsuiteresult", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.665525Z"}
2025-12-27 00:59:08,694 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa001' (Stream is closed)
2025-12-27 00:59:08,698 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#tobedisabled", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:08.698057Z"}
2025-12-27 00:59:08,703 - INFO - {"url": "https://vitest.dev/config/browser/api#api-port", "links_found": 4, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.703311Z"}
2025-12-27 00:59:08,704 - INFO - {"url": "https://vitest.dev/guide/reporters#junit-reporter", "links_found": 29, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.704491Z"}
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
[INIT].... → Crawl4AI 0.7.8
2025-12-27 00:59:08,957 - INFO - {"url": "https://vitest.dev/api/vi#vi-unstuballglobals", "links_found": 71, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:08.957345Z"}
2025-12-27 00:59:08,964 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send FIN b'170e8b12178aa005' (Stream is closed)
2025-12-27 00:59:08,965 - INFO - {"source_id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "url": "https://vitest.dev/llms-full.txt#recordartifact", "event": "result_published", "level": "info", "timestamp": "2025-12-27T00:59:08.965549Z"}
[FETCH]... ↓ https://tanstack.com/ai/latest/docs/contributors
| ✓ | ⏱: 204.64s
00:59:09 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

025-12-27 00:59:04,184 - INFO - Wrapper: Completed Call, calling success_handler
[SCRAPE].. ◆ https://vitest.dev/api/vi#vi-unstuballglobals
| ✓ | ⏱: 9.88s
[COMPLETE] ● https://vitest.dev/api/vi#vi-unstuballglobals
| ✓ | ⏱: 199.62s
2025-12-27 00:59:04,188 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/cli"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:59:04.188019Z"}
2025-12-27 00:59:04,188 - INFO - {"url": "https://vitest.dev/guide/cli", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:59:04.188090Z"}
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8ae5fccaa015' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa000' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa001' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa002' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa003' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa004' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa005' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa006' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa007' (Stream is closed)
2025-12-27 00:59:04,219 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa008' (Stream is closed)
2025-12-27 00:59:04,268 - INFO - {"url": "https://vitest.dev/llms-full.txt#tobedisabled", "links_found": 519, "event": "crawl_completed", "level": "info", "timestamp": "2025-12-27T00:59:04.268854Z"}
2025-12-27 00:59:04,269 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/llms-full.txt#vi-restoreallmocks"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:59:04.269807Z"}
2025-12-27 00:59:04,269 - INFO - {"url": "https://vitest.dev/llms-full.txt#vi-restoreallmocks", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:59:04.269864Z"}
2025-12-27 00:59:04,341 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/api/mock#mockreset"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:59:04.341203Z"}
2025-12-27 00:59:04,341 - INFO - {"url": "https://vitest.dev/api/mock#mockreset", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:59:04.341424Z"}
2025-12-27 00:59:04,417 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/#threads"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:59:04.417044Z"}
2025-12-27 00:59:04,417 - INFO - {"url": "https://vitest.dev/config/#threads", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:59:04.417203Z"}
2025-12-27 00:59:04,486 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
[FETCH]... ↓ https://vitest.dev/config/#mockreset
| ✓ | ⏱: 200.06s
00:59:04 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

SCRAPE].. ◆ https://vitest.dev/llms-full.txt#ontestsuiteresult
| ✓ | ⏱: 44.52s
[COMPLETE] ● https://vitest.dev/llms-full.txt#ontestsuiteresult
| ✓ | ⏱: 179.78s
2025-12-27 00:58:44,428 - INFO - [nsqd:4150] connecting to nsqd
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8ae5fccaa015' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa000' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa001' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa002' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa003' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa004' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa005' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa006' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa007' (Stream is closed)
2025-12-27 00:58:44,433 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa008' (Stream is closed)
2025-12-27 00:58:44,435 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
2025-12-27 00:58:44,435 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
2025-12-27 00:58:44,436 - INFO - [nsqd:4150:ingest.task:worker] received heartbeat
2025-12-27 00:58:44,438 - INFO - [nsqd:4150:nsqd:4150] IDENTIFY received {'max_rdy_count': 2500, 'version': '1.3.0', 'max_msg_timeout': 900000, 'msg_timeout': 60000, 'tls_v1': False, 'deflate': False, 'deflate_level': 6, 'max_deflate_level': 6, 'snappy': False, 'sample_rate': 0, 'auth_required': False, 'output_buffer_size': 16384, 'output_buffer_timeout': 250}
2025-12-27 00:58:44,438 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/api/browser/locators"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:58:44.438721Z"}
2025-12-27 00:58:44,438 - INFO - {"url": "https://vitest.dev/api/browser/locators", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:58:44.438824Z"}
2025-12-27 00:58:44,503 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/browser/testerhtmlpath"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:58:44.503281Z"}
2025-12-27 00:58:44,503 - INFO - {"url": "https://vitest.dev/config/browser/testerhtmlpath", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:58:44.503453Z"}
2025-12-27 00:58:44,570 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/config/expect#expect-requireassertions"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:58:44.570860Z"}
2025-12-27 00:58:44,571 - INFO - {"url": "https://vitest.dev/config/expect#expect-requireassertions", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:58:44.571081Z"}
2025-12-27 00:58:44,645 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/profiling-test-performance.html#code-coverage"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:58:44.645285Z"}
2025-12-27 00:58:44,645 - INFO - {"url": "https://vitest.dev/guide/profiling-test-performance.html#code-coverage", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:58:44.645446Z"}
2025-12-27 00:58:44,717 - INFO - {"data": {"correlation_id": "8e503893-2228-49b7-b080-ac4acbaedb22", "depth": 1, "exclusions": [], "gemini_api_key": "AIzaSyAPQkbzk8Fd2CgibxNe_KouxmosW7q3ErM", "id": "9d39f39c-25f3-4cc9-a8ae-f28445cdd0de", "max_depth": 1, "type": "web", "url": "https://vitest.dev/guide/improving-performance"}, "event": "message_received", "level": "info", "timestamp": "2025-12-27T00:58:44.717128Z"}
2025-12-27 00:58:44,717 - INFO - {"url": "https://vitest.dev/guide/improving-performance", "event": "crawl_starting", "level": "info", "timestamp": "2025-12-27T00:58:44.717310Z"}
[FETCH]... ↓ https://vitest.dev/config/browser/api#api-port
| ✓ | ⏱: 180.29s
00:58:44 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini

---

SCRAPE].. ◆ https://vitest.dev/llms-full.txt#vi-domock
| ✓ | ⏱: 45.55s
[COMPLETE] ● https://vitest.dev/llms-full.txt#vi-domock
| ✓ | ⏱: 135.63s
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa001' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa002' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa003' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa004' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa005' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa006' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa007' (Stream is closed)
2025-12-27 00:57:59,885 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: SendError: failed to send TOUCH b'170e8b12178aa008' (Stream is closed)
2025-12-27 00:57:59,885 - WARNING - [nsqd:4150:ingest.task:worker] connection is stale (90.62s), closing
2025-12-27 00:57:59,885 - WARNING - [nsqd:4150:nsqd:4150] connection is stale (90.62s), closing
2025-12-27 00:57:59,885 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 00:57:59,886 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 00:57:59,886 - WARNING - [nsqd:4150:ingest.task:worker] connection closed
2025-12-27 00:57:59,886 - WARNING - [nsqd:4150] connection closed
2025-12-27 00:57:59,886 - INFO - [nsqd:4150] attempting to reconnect in 15.00s
2025-12-27 00:57:59,889 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 00:57:59,889 - ERROR - [nsqd:4150:ingest.task:worker] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 00:57:59,889 - ERROR - uncaught exception in data event
Traceback (most recent call last):
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 289, in _read_body
    self.trigger(event.DATA, conn=self, data=data)
  File "/usr/local/lib/python3.12/site-packages/nsq/event.py", line 84, in trigger
    ev(*args, **kwargs)
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 504, in _on_data
    self.send(protocol.nop())
  File "/usr/local/lib/python3.12/site-packages/nsq/conn.py", line 295, in send
    return self.stream.write(self.encoder.encode(data))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 523, in write
    self._check_closed()
  File "/usr/local/lib/python3.12/site-packages/tornado/iostream.py", line 998, in _check_closed
    raise StreamClosedError(real_error=self.error)
tornado.iostream.StreamClosedError: Stream is closed
2025-12-27 00:57:59,889 - ERROR - [nsqd:4150:nsqd:4150] ERROR: ConnectionClosedError('Stream is closed')
2025-12-27 00:57:59,893 - INFO - [nsqd:4150:ingest.task:worker] connecting to nsqd
2025-12-27 00:57:59,898 - INFO - [nsqd:4150:ingest.task:worker] IDENTIFY sent {'short_id': 'd3fe91d786d3', 'long_id': 'd3fe91d786d3', 'client_id': 'd3fe91d786d3', 'hostname': 'd3fe91d786d3', 'heartbeat_interval': 30000, 'feature_negotiation': True, 'tls_v1': False, 'snappy': False, 'deflate': False, 'deflate_level': 6, 'output_buffer_timeout': 250, 'output_buffer_size': 16384, 'sample_rate': 0, 'user_agent': 'pynsq/0.9.1'}
[FETCH]... ↓ https://vitest.dev/llms-full.txt#ontestsuiteresult
| ✓ | ⏱: 135.25s
00:57:59 - LiteLLM:INFO: utils.py:3476 -
LiteLLM completion() model= gemini-3-flash-preview; provider = gemini